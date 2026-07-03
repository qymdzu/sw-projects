// Package model 定义 GORM 数据模型，与 docs/design/数据模型设计.md 1:1 映射。

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ModelSetting 模型配置表（本次新增 — Phase B）。
//
// 业务背景：每个用户可以私有配置多个 LLM 提供商（openai / anthropic / qwen / deepseek / ollama / custom），
// 但同一时刻只能有 1 个被标记为 is_default=true 的活跃配置，用于 AI 出题、推荐等场景。
//
// 表设计要点：
//   - user_id + provider 组合唯一
//   - api_key_ciphertext 是加密后的密文（GCM 模式），api_key_nonce 是 12 字节 nonce
//   - config JSONB 字段透传额外 LLM 参数（temperature、top_p、max_tokens 等）
//   - is_default 用于标记用户当前活跃配置（同一 user 只能 1 条 true）
//
// 加密方式详见 backend/pkg/crypto/aesgcm.go，
// 加密密钥派生路径为 SHA256(JWT_SECRET + ".model-setting-salt.v1")。
type ModelSetting struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_ms_user_provider,priority:1" json:"user_id"`
	Provider         string         `gorm:"type:varchar(32);not null;uniqueIndex:idx_ms_user_provider,priority:2" json:"provider"`
	APIEndpoint      string         `gorm:"type:varchar(500);not null" json:"api_endpoint"`
	APIKeyCiphertext []byte         `gorm:"type:bytea;not null" json:"-"`
	APIKeyNonce      []byte         `gorm:"type:bytea;not null" json:"-"`
	Model            string         `gorm:"type:varchar(100);not null" json:"model"`
	Config           datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"extra_config,omitempty"`
	IsDefault        bool           `gorm:"type:boolean;not null;default:false;index:idx_ms_user_default" json:"is_default"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// TableName 指定 GORM 表名。
func (ModelSetting) TableName() string { return "model_settings" }

// Provider 枚举常量（服务层与前端共享）。
const (
	ProviderOpenAI    = "openai"
	ProviderAnthropic = "anthropic"
	ProviderQwen      = "qwen"
	ProviderDeepSeek  = "deepseek"
	ProviderOllama    = "ollama"
	ProviderCustom    = "custom"
)

// SupportedProviders 所有受支持的 Provider 列表（用于服务层校验）。
var SupportedProviders = []string{
	ProviderOpenAI, ProviderAnthropic, ProviderQwen,
	ProviderDeepSeek, ProviderOllama, ProviderCustom,
}

// IsSupportedProvider 判断 provider 是否在白名单枚举内。
func IsSupportedProvider(p string) bool {
	for _, sp := range SupportedProviders {
		if sp == p {
			return true
		}
	}
	return false
}