// Package service 封装业务用例。

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/pkg/crypto"
)

// Service 层错误定义。
//
// 命名风格与既有 ErrAccountConflict / ErrAccountNotFound / ErrPasswordInvalid 保持一致。
var (
	ErrSettingNotFound     = errors.New("model setting not found")
	ErrSettingParamInvalid = errors.New("invalid model setting")
	ErrUnsupportedProvider = errors.New("unsupported provider")
	ErrEndpointUnsafe      = errors.New("unsafe endpoint")
	ErrEndpointInvalid     = errors.New("invalid endpoint url")
	ErrInternalCrypto      = errors.New("crypto failure")
)

// SaveSettingRequest 业务入参（API Key 以明文传入，由 service 加密）。
type SaveSettingRequest struct {
	Provider    string
	APIEndpoint string
	APIKey      string
	Model       string
	ExtraConfig map[string]any
}

// SettingDTO 业务出参（API Key 是否明文由调用方决定）。
type SettingDTO struct {
	Provider     string
	APIEndpoint  string
	APIKey       string // 明文（仅 GetByProvider / GetActive / Save 返回）
	Model        string
	ExtraConfig  map[string]any
	IsDefault    bool
	UpdatedAt    time.Time
}

// SettingSummary 用于列表场景的脱敏 DTO。
type SettingSummary struct {
	Provider     string
	APIEndpoint  string
	APIKeyMasked string
	Model        string
	IsDefault    bool
	UpdatedAt    time.Time
}

// SettingService 是模型配置业务服务接口。
type SettingService interface {
	Save(ctx context.Context, userID uuid.UUID, in SaveSettingRequest) (*SettingDTO, error)
	GetByProvider(ctx context.Context, userID uuid.UUID, provider string) (*SettingDTO, error)
	GetActive(ctx context.Context, userID uuid.UUID) (*SettingDTO, error)
	List(ctx context.Context, userID uuid.UUID) ([]*SettingSummary, error)
	Activate(ctx context.Context, userID uuid.UUID, provider string) error
	Delete(ctx context.Context, userID uuid.UUID, provider string) error
}

type settingSvc struct {
	repo   repository.SettingRepository
	cipher *crypto.AESGCM
}

// NewSettingService 构造 SettingService。
func NewSettingService(repo repository.SettingRepository, cipher *crypto.AESGCM) SettingService {
	return &settingSvc{repo: repo, cipher: cipher}
}

// Save 保存或更新一条配置，并自动设为 default。
func (s *settingSvc) Save(ctx context.Context, userID uuid.UUID, in SaveSettingRequest) (*SettingDTO, error) {
	if err := s.validate(in); err != nil {
		return nil, err
	}
	ct, nonce, err := s.cipher.Encrypt([]byte(in.APIKey))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalCrypto, err)
	}
	cfgBytes, err := json.Marshal(nonNilMap(in.ExtraConfig))
	if err != nil {
		return nil, fmt.Errorf("%w: bad extra_config json: %v", ErrSettingParamInvalid, err)
	}
	m := &model.ModelSetting{
		ID:               uuid.New(),
		UserID:           userID,
		Provider:         in.Provider,
		APIEndpoint:      in.APIEndpoint,
		APIKeyCiphertext: ct,
		APIKeyNonce:      nonce,
		Model:            in.Model,
		Config:           datatypes.JSON(cfgBytes),
		IsDefault:        true,
	}
	// Upsert：以 (user_id, provider) 为冲突键，更新全部字段
	if err := s.repo.Upsert(ctx, m); err != nil {
		return nil, err
	}
	// 然后切换 default（事务化）。Upsert 已把 is_default 写入 true，但为了应对老记录，
	// 仍显式调一次 SetActive，确保用户级别只有一条 default。
	if err := s.repo.SetActive(ctx, userID, in.Provider); err != nil {
		return nil, err
	}
	return s.toDTO(m, in.APIKey), nil
}

// GetByProvider 取指定 provider 的配置（API Key 明文返回）。
func (s *settingSvc) GetByProvider(ctx context.Context, userID uuid.UUID, provider string) (*SettingDTO, error) {
	if !model.IsSupportedProvider(provider) {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, provider)
	}
	m, err := s.repo.GetByUserAndProvider(ctx, userID, provider)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSettingNotFound
		}
		return nil, err
	}
	plain, err := s.cipher.Decrypt(m.APIKeyCiphertext, m.APIKeyNonce)
	if err != nil {
		return nil, ErrInternalCrypto
	}
	return s.toDTO(m, string(plain)), nil
}

// GetActive 取用户当前 default 配置（API Key 明文返回）。
func (s *settingSvc) GetActive(ctx context.Context, userID uuid.UUID) (*SettingDTO, error) {
	m, err := s.repo.GetActiveByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSettingNotFound
		}
		return nil, err
	}
	plain, err := s.cipher.Decrypt(m.APIKeyCiphertext, m.APIKeyNonce)
	if err != nil {
		return nil, ErrInternalCrypto
	}
	return s.toDTO(m, string(plain)), nil
}

// List 列出用户所有配置（API Key 仅返回掩码）。
func (s *settingSvc) List(ctx context.Context, userID uuid.UUID) ([]*SettingSummary, error) {
	ms, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*SettingSummary, 0, len(ms))
	for _, m := range ms {
		plain, _ := s.cipher.Decrypt(m.APIKeyCiphertext, m.APIKeyNonce)
		out = append(out, &SettingSummary{
			Provider:     m.Provider,
			APIEndpoint:  m.APIEndpoint,
			APIKeyMasked: maskAPIKey(string(plain)),
			Model:        m.Model,
			IsDefault:    m.IsDefault,
			UpdatedAt:    m.UpdatedAt,
		})
	}
	return out, nil
}

// Activate 切换用户的 default 到指定 provider（要求目标已存在）。
func (s *settingSvc) Activate(ctx context.Context, userID uuid.UUID, provider string) error {
	if !model.IsSupportedProvider(provider) {
		return fmt.Errorf("%w: %s", ErrUnsupportedProvider, provider)
	}
	if err := s.repo.SetActive(ctx, userID, provider); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSettingNotFound
		}
		return err
	}
	return nil
}

// Delete 删除一条配置；如果删的是 default，则把剩余第一条提升为 default。
// 校验：用户必须有至少 1 条配置才允许操作（防止"删完"）。
func (s *settingSvc) Delete(ctx context.Context, userID uuid.UUID, provider string) error {
	if !model.IsSupportedProvider(provider) {
		return fmt.Errorf("%w: %s", ErrUnsupportedProvider, provider)
	}
	target, err := s.repo.GetByUserAndProvider(ctx, userID, provider)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSettingNotFound
		}
		return err
	}
	if err := s.repo.DeleteByUserAndProvider(ctx, userID, provider); err != nil {
		return err
	}
	// 若删的是默认，自动回落到最近一条
	if target.IsDefault {
		all, err := s.repo.ListByUser(ctx, userID)
		if err != nil {
			return err
		}
		if len(all) > 0 {
			if err := s.repo.SetActive(ctx, userID, all[0].Provider); err != nil {
				return err
			}
		}
	}
	return nil
}

// validate 四件套：provider 枚举、URL 合法且 scheme 正确、API Key 长度、Model 非空。
//
// Ollama 允许 http://localhost / http://127.0.0.1（便于本地部署）。
// 其他 Provider 一律要求 https://。
func (s *settingSvc) validate(in SaveSettingRequest) error {
	if !model.IsSupportedProvider(in.Provider) {
		return fmt.Errorf("%w: %s", ErrUnsupportedProvider, in.Provider)
	}
	if strings.TrimSpace(in.APIEndpoint) == "" {
		return fmt.Errorf("%w: api_endpoint required", ErrEndpointInvalid)
	}
	u, err := url.Parse(in.APIEndpoint)
	if err != nil || u.Host == "" || u.Scheme == "" {
		return fmt.Errorf("%w: %s", ErrEndpointInvalid, in.APIEndpoint)
	}
	if in.Provider != model.ProviderOllama {
		if u.Scheme != "https" {
			return fmt.Errorf("%w: https required for %s", ErrEndpointInvalid, in.Provider)
		}
	} else {
		if u.Scheme != "http" && u.Scheme != "https" {
			return fmt.Errorf("%w: http(s) required for ollama", ErrEndpointInvalid)
		}
	}
	if len(in.APIKey) < 8 {
		return fmt.Errorf("%w: api_key too short (>=8)", ErrSettingParamInvalid)
	}
	if strings.TrimSpace(in.Model) == "" {
		return fmt.Errorf("%w: model required", ErrSettingParamInvalid)
	}
	return nil
}

func (s *settingSvc) toDTO(m *model.ModelSetting, plainKey string) *SettingDTO {
	cfg := map[string]any{}
	_ = json.Unmarshal(m.Config, &cfg)
	return &SettingDTO{
		Provider:    m.Provider,
		APIEndpoint: m.APIEndpoint,
		APIKey:      plainKey,
		Model:       m.Model,
		ExtraConfig: cfg,
		IsDefault:   m.IsDefault,
		UpdatedAt:   m.UpdatedAt,
	}
}

// maskAPIKey 返回前 4 后 4，中间 ***；不足 8 时全打码。
func maskAPIKey(k string) string {
	if len(k) < 8 {
		return strings.Repeat("*", len(k))
	}
	return k[:4] + "***" + k[len(k)-4:]
}

// nonNilMap 把 nil map 归一化为空 map，便于下游 json.Marshal 输出 "{}"。
func nonNilMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	return m
}