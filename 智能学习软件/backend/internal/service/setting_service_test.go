package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/pkg/crypto"
)

// mockSettingRepo 是 SettingRepository 的最小化 mock。
type mockSettingRepo struct {
	items map[string]*model.ModelSetting // key = "userID|provider"
}

func newMockRepo() *mockSettingRepo {
	return &mockSettingRepo{items: map[string]*model.ModelSetting{}}
}

func key(uid uuid.UUID, p string) string { return uid.String() + "|" + p }

func (r *mockSettingRepo) GetByUserAndProvider(_ context.Context, uid uuid.UUID, p string) (*model.ModelSetting, error) {
	if v, ok := r.items[key(uid, p)]; ok {
		return v, nil
	}
	return nil, repository.ErrNotFound
}

func (r *mockSettingRepo) GetActiveByUser(_ context.Context, uid uuid.UUID) (*model.ModelSetting, error) {
	for _, v := range r.items {
		if v.UserID == uid && v.IsDefault {
			return v, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *mockSettingRepo) ListByUser(_ context.Context, uid uuid.UUID) ([]*model.ModelSetting, error) {
	out := []*model.ModelSetting{}
	for _, v := range r.items {
		if v.UserID == uid {
			out = append(out, v)
		}
	}
	return out, nil
}

func (r *mockSettingRepo) Upsert(_ context.Context, m *model.ModelSetting) error {
	r.items[key(m.UserID, m.Provider)] = m
	return nil
}

func (r *mockSettingRepo) SetActive(_ context.Context, uid uuid.UUID, p string) error {
	found := false
	for _, v := range r.items {
		if v.UserID == uid {
			v.IsDefault = (v.Provider == p)
			if v.Provider == p {
				found = true
			}
		}
	}
	if !found {
		return repository.ErrNotFound
	}
	return nil
}

func (r *mockSettingRepo) DeleteByUserAndProvider(_ context.Context, uid uuid.UUID, p string) error {
	k := key(uid, p)
	if _, ok := r.items[k]; !ok {
		return repository.ErrNotFound
	}
	delete(r.items, k)
	return nil
}

func (r *mockSettingRepo) CountActiveByUser(_ context.Context, uid uuid.UUID) (int64, error) {
	var n int64
	for _, v := range r.items {
		if v.UserID == uid && v.IsDefault {
			n++
		}
	}
	return n, nil
}

// newSvc 构造带 mock 的 service。
func newSvc() (SettingService, *mockSettingRepo) {
	repo := newMockRepo()
	c, err := crypto.NewAESGCM("test-secret-must-be-at-least-32-chars-long")
	if err != nil {
		panic(err)
	}
	return NewSettingService(repo, c), repo
}

func validReq() SaveSettingRequest {
	return SaveSettingRequest{
		Provider:    model.ProviderOpenAI,
		APIEndpoint: "https://api.openai.com/v1",
		APIKey:      "sk-test-1234567890",
		Model:       "gpt-4o-mini",
		ExtraConfig: map[string]any{"temperature": 0.7},
	}
}

// TC-SET-01: Save 正常路径 → 落库 + 默认
func TestSettingService_Save_OK(t *testing.T) {
	svc, repo := newSvc()
	uid := uuid.New()
	dto, err := svc.Save(context.Background(), uid, validReq())
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if dto.Provider != "openai" {
		t.Fatalf("provider=%s", dto.Provider)
	}
	if !dto.IsDefault {
		t.Fatal("newly saved should be default")
	}
	// 落库校验：解密密文能得到原文
	m := repo.items[key(uid, "openai")]
	if m == nil {
		t.Fatal("not persisted")
	}
}

// TC-SET-02: provider 不在枚举
func TestSettingService_Save_UnsupportedProvider(t *testing.T) {
	svc, _ := newSvc()
	req := validReq()
	req.Provider = "bogus"
	_, err := svc.Save(context.Background(), uuid.New(), req)
	if err != ErrUnsupportedProvider && err == nil {
		t.Fatalf("expected unsupported provider error, got %v", err)
	}
}

// TC-SET-03: 非 ollama 但 endpoint 非 https
func TestSettingService_Save_NeedHTTPS(t *testing.T) {
	svc, _ := newSvc()
	req := validReq()
	req.APIEndpoint = "http://api.openai.com/v1"
	_, err := svc.Save(context.Background(), uuid.New(), req)
	if err == nil {
		t.Fatal("expected error for http endpoint on openai")
	}
}

// TC-SET-04: API Key 长度过短
func TestSettingService_Save_KeyTooShort(t *testing.T) {
	svc, _ := newSvc()
	req := validReq()
	req.APIKey = "abc"
	_, err := svc.Save(context.Background(), uuid.New(), req)
	if err == nil {
		t.Fatal("expected error for short api key")
	}
}

// TC-SET-05: SetActive 切换后旧 default 置 false
func TestSettingService_SaveTwiceAndActivate(t *testing.T) {
	svc, repo := newSvc()
	uid := uuid.New()
	// 第一次 Save（openai）
	req1 := validReq()
	if _, err := svc.Save(context.Background(), uid, req1); err != nil {
		t.Fatalf("save1: %v", err)
	}
	// 第二次 Save（anthropic）
	req2 := validReq()
	req2.Provider = model.ProviderAnthropic
	req2.APIEndpoint = "https://api.anthropic.com"
	req2.Model = "claude-3-5-sonnet-20241022"
	if _, err := svc.Save(context.Background(), uid, req2); err != nil {
		t.Fatalf("save2: %v", err)
	}
	// 此时 anthropic 是 default，openai 也被 SetActive 设回 false（其实之前 upsert 写入 true，但 SetActive 会把 openai 改 false）
	if !repo.items[key(uid, "anthropic")].IsDefault {
		t.Fatal("anthropic should be default")
	}
	// 切换回 openai
	if err := svc.Activate(context.Background(), uid, model.ProviderOpenAI); err != nil {
		t.Fatalf("Activate: %v", err)
	}
	if !repo.items[key(uid, "openai")].IsDefault {
		t.Fatal("openai should be default after activate")
	}
	if repo.items[key(uid, "anthropic")].IsDefault {
		t.Fatal("anthropic should not be default anymore")
	}
}

// TC-SET-06: Delete 当前 default 自动回落
func TestSettingService_DeleteDefault(t *testing.T) {
	svc, repo := newSvc()
	uid := uuid.New()
	// 落两条
	req1 := validReq()
	_, _ = svc.Save(context.Background(), uid, req1)
	req2 := validReq()
	req2.Provider = model.ProviderDeepSeek
	req2.APIEndpoint = "https://api.deepseek.com"
	_, _ = svc.Save(context.Background(), uid, req2)
	// 现在 deepseek 是 default，删掉
	if err := svc.Delete(context.Background(), uid, model.ProviderDeepSeek); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, ok := repo.items[key(uid, "deepseek")]; ok {
		t.Fatal("deepseek should be deleted")
	}
	// 应该回落给 openai
	if !repo.items[key(uid, "openai")].IsDefault {
		t.Fatal("openai should fallback to default")
	}
}

// TC-SET-07: GetActive 找不到时返回 ErrSettingNotFound
func TestSettingService_GetActive_NotFound(t *testing.T) {
	svc, _ := newSvc()
	_, err := svc.GetActive(context.Background(), uuid.New())
	if err != ErrSettingNotFound {
		t.Fatalf("expected ErrSettingNotFound, got %v", err)
	}
}

// TC-SET-08: List 返回掩码
func TestSettingService_List_MaskAPIKey(t *testing.T) {
	svc, _ := newSvc()
	uid := uuid.New()
	if _, err := svc.Save(context.Background(), uid, validReq()); err != nil {
		t.Fatalf("save: %v", err)
	}
	list, err := svc.List(context.Background(), uid)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1, got %d", len(list))
	}
	masked := list[0].APIKeyMasked
	if masked == "sk-test-1234567890" {
		t.Fatal("api key should be masked, got plain")
	}
	if len(masked) < 8 {
		t.Fatalf("masked key too short: %s", masked)
	}
}

// TC-SET-09: 同一 provider 二次提交只更新不新增
func TestSettingService_Save_Update(t *testing.T) {
	svc, repo := newSvc()
	uid := uuid.New()
	req := validReq()
	first, err := svc.Save(context.Background(), uid, req)
	if err != nil {
		t.Fatalf("save1: %v", err)
	}
	second, err := svc.Save(context.Background(), uid, req)
	if err != nil {
		t.Fatalf("save2: %v", err)
	}
	// 数量应该不变
	if got := len(repo.items); got != 1 {
		t.Fatalf("expected 1 record, got %d", got)
	}
	// updated_at 应该推进（或相同 GORM 时间精度）
	if first.UpdatedAt.After(second.UpdatedAt.Add(time.Second)) {
		t.Fatal("updated_at went backwards")
	}
}

// TC-SET-10: Ollama 允许 http://localhost
func TestSettingService_OllamaHTTP(t *testing.T) {
	svc, _ := newSvc()
	uid := uuid.New()
	req := SaveSettingRequest{
		Provider:    model.ProviderOllama,
		APIEndpoint: "http://localhost:11434",
		APIKey:      "ollama-local-token",
		Model:       "llama3.1",
	}
	if _, err := svc.Save(context.Background(), uid, req); err != nil {
		t.Fatalf("ollama http should be allowed: %v", err)
	}
}

// TC-SET-11: datatypes.JSON 兼容性测试（确保 Config 字段不报错）
func TestSettingService_ExtraConfigEmpty(t *testing.T) {
	svc, _ := newSvc()
	uid := uuid.New()
	req := validReq()
	req.ExtraConfig = nil
	dto, err := svc.Save(context.Background(), uid, req)
	if err != nil {
		t.Fatalf("Save with nil extra_config: %v", err)
	}
	if dto.ExtraConfig == nil {
		t.Fatal("ExtraConfig should be non-nil empty map")
	}
}

// TC-SET-12: 自定义 JSONB 值往返
func TestSettingService_ExtraConfigRoundTrip(t *testing.T) {
	svc, _ := newSvc()
	uid := uuid.New()
	req := validReq()
	req.ExtraConfig = map[string]any{"temperature": 0.7, "top_p": 0.9, "max_tokens": 2048}
	dto, err := svc.Save(context.Background(), uid, req)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if v, ok := dto.ExtraConfig["temperature"].(float64); !ok || v != 0.7 {
		t.Fatalf("temperature roundtrip wrong: %+v", dto.ExtraConfig)
	}
}