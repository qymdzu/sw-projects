package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"smart-learning/internal/middleware"
	"smart-learning/internal/service"
)

// fakeSettingSvc 是 service.SettingService 的最小化 fake 实现。
//
// 不依赖业务 model / repository，纯单元测试 handler 层：
//   - 入参校验
//   - ctx user_id 提取
//   - 错误映射
//   - 响应格式
type fakeSettingSvc struct {
	saveFn      func(ctx context.Context, uid uuid.UUID, in service.SaveSettingRequest) (*service.SettingDTO, error)
	getActiveFn func(ctx context.Context, uid uuid.UUID) (*service.SettingDTO, error)
	listFn      func(ctx context.Context, uid uuid.UUID) ([]*service.SettingSummary, error)
	activateFn  func(ctx context.Context, uid uuid.UUID, p string) error
	deleteFn    func(ctx context.Context, uid uuid.UUID, p string) error
}

func (f *fakeSettingSvc) Save(ctx context.Context, uid uuid.UUID, in service.SaveSettingRequest) (*service.SettingDTO, error) {
	if f.saveFn == nil {
		return nil, errors.New("saveFn not stubbed")
	}
	return f.saveFn(ctx, uid, in)
}
func (f *fakeSettingSvc) GetByProvider(ctx context.Context, uid uuid.UUID, p string) (*service.SettingDTO, error) {
	return nil, nil
}
func (f *fakeSettingSvc) GetActive(ctx context.Context, uid uuid.UUID) (*service.SettingDTO, error) {
	if f.getActiveFn == nil {
		return nil, errors.New("getActiveFn not stubbed")
	}
	return f.getActiveFn(ctx, uid)
}
func (f *fakeSettingSvc) List(ctx context.Context, uid uuid.UUID) ([]*service.SettingSummary, error) {
	if f.listFn == nil {
		return nil, nil
	}
	return f.listFn(ctx, uid)
}
func (f *fakeSettingSvc) Activate(ctx context.Context, uid uuid.UUID, p string) error {
	if f.activateFn == nil {
		return errors.New("activateFn not stubbed")
	}
	return f.activateFn(ctx, uid, p)
}
func (f *fakeSettingSvc) Delete(ctx context.Context, uid uuid.UUID, p string) error {
	if f.deleteFn == nil {
		return errors.New("deleteFn not stubbed")
	}
	return f.deleteFn(ctx, uid, p)
}

func init() {
	gin.SetMode(gin.TestMode)
}

// buildRouter 构造一个最小 gin.Engine，注入 user_id 后挂载 setting 路由。
func buildRouter(t *testing.T, fake *fakeSettingSvc) *gin.Engine {
	t.Helper()
	h := NewSettingHandler(fake)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.CtxUserID, uuid.New().String())
		c.Next()
	})
	r.POST("/settings/model", h.CreateOrUpdate)
	r.GET("/settings/model", h.GetActive)
	r.PUT("/settings/model", h.Update)
	r.DELETE("/settings/model", h.Delete)
	return r
}

// TC-FE-H01: CreateOrUpdate 成功路径 — API Key 必须返回掩码 ***
func TestHandler_CreateOrUpdate_OK_KeyMasked(t *testing.T) {
	fake := &fakeSettingSvc{
		saveFn: func(_ context.Context, _ uuid.UUID, in service.SaveSettingRequest) (*service.SettingDTO, error) {
			return &service.SettingDTO{
				Provider:    in.Provider,
				APIEndpoint: in.APIEndpoint,
				APIKey:      in.APIKey, // service 层返回明文
				Model:       in.Model,
				IsDefault:   true,
				UpdatedAt:   time.Now(),
			}, nil
		},
	}
	r := buildRouter(t, fake)
	body := bytes.NewReader([]byte(`{"provider":"openai","api_endpoint":"https://api.openai.com/v1","api_key":"sk-test-1234567890","model":"gpt-4o-mini"}`))
	req := httptest.NewRequest(http.MethodPost, "/settings/model", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("code=%d", resp.Code)
	}
	if resp.Data["api_key"] != "***" {
		t.Fatalf("api_key must be masked as ***, got %v", resp.Data["api_key"])
	}
}

// TC-FE-H02: CreateOrUpdate 参数缺失 → 400
func TestHandler_CreateOrUpdate_BadJSON(t *testing.T) {
	fake := &fakeSettingSvc{}
	r := buildRouter(t, fake)
	req := httptest.NewRequest(http.MethodPost, "/settings/model", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
}

// TC-FE-H03: CreateOrUpdate service 返回 ErrUnsupportedProvider → 400
func TestHandler_CreateOrUpdate_UnsupportedProvider(t *testing.T) {
	fake := &fakeSettingSvc{
		saveFn: func(_ context.Context, _ uuid.UUID, _ service.SaveSettingRequest) (*service.SettingDTO, error) {
			return nil, service.ErrUnsupportedProvider
		},
	}
	r := buildRouter(t, fake)
	body := bytes.NewReader([]byte(`{"provider":"openai","api_endpoint":"https://api.openai.com/v1","api_key":"sk-test-1234567890","model":"gpt-4o-mini"}`))
	req := httptest.NewRequest(http.MethodPost, "/settings/model", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

// TC-FE-H04: CreateOrUpdate service 返回 ErrInternalCrypto → 500 + 通用消息
func TestHandler_CreateOrUpdate_CryptoErr_NoLeak(t *testing.T) {
	fake := &fakeSettingSvc{
		saveFn: func(_ context.Context, _ uuid.UUID, _ service.SaveSettingRequest) (*service.SettingDTO, error) {
			return nil, service.ErrInternalCrypto
		},
	}
	r := buildRouter(t, fake)
	body := bytes.NewReader([]byte(`{"provider":"openai","api_endpoint":"https://api.openai.com/v1","api_key":"sk-test-1234567890","model":"gpt-4o-mini"}`))
	req := httptest.NewRequest(http.MethodPost, "/settings/model", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("want 500, got %d", rec.Code)
	}
	// 响应体不应包含 "crypto" 之类的内部细节
	if bytes.Contains(rec.Body.Bytes(), []byte("crypto")) {
		t.Fatalf("response leaks crypto detail: %s", rec.Body.String())
	}
}

// TC-FE-H05: GetActive 找不到 → 200 + null
func TestHandler_GetActive_NotFound(t *testing.T) {
	fake := &fakeSettingSvc{
		getActiveFn: func(_ context.Context, _ uuid.UUID) (*service.SettingDTO, error) {
			return nil, service.ErrSettingNotFound
		},
	}
	r := buildRouter(t, fake)
	req := httptest.NewRequest(http.MethodGet, "/settings/model", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	// body.data.setting 应为 null
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"setting":null`)) {
		t.Fatalf("want setting:null in body, got %s", rec.Body.String())
	}
}

// TC-FE-H06: Delete 缺 provider 参数 → 400
func TestHandler_Delete_NoProvider(t *testing.T) {
	fake := &fakeSettingSvc{
		deleteFn: func(_ context.Context, _ uuid.UUID, _ string) error { return nil },
	}
	r := buildRouter(t, fake)
	req := httptest.NewRequest(http.MethodDelete, "/settings/model", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
}

// TC-FE-H07: Delete 成功 → 204
func TestHandler_Delete_Ok(t *testing.T) {
	fake := &fakeSettingSvc{
		deleteFn: func(_ context.Context, _ uuid.UUID, p string) error {
			if p == "" {
				return errors.New("provider empty")
			}
			return nil
		},
	}
	r := buildRouter(t, fake)
	req := httptest.NewRequest(http.MethodDelete, "/settings/model?provider=openai", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", rec.Code)
	}
}

// TC-FE-H08: Activate (PUT) 成功
func TestHandler_Activate_Ok(t *testing.T) {
	fake := &fakeSettingSvc{
		activateFn: func(_ context.Context, _ uuid.UUID, p string) error {
			if p != "anthropic" {
				return errors.New("stub: only anthropic allowed")
			}
			return nil
		},
	}
	r := buildRouter(t, fake)
	body := bytes.NewReader([]byte(`{"provider":"anthropic"}`))
	req := httptest.NewRequest(http.MethodPut, "/settings/model", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}

// TC-FE-H09: Activate 找不到目标 → 404
func TestHandler_Activate_NotFound(t *testing.T) {
	fake := &fakeSettingSvc{
		activateFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return service.ErrSettingNotFound
		},
	}
	r := buildRouter(t, fake)
	body := bytes.NewReader([]byte(`{"provider":"qwen"}`))
	req := httptest.NewRequest(http.MethodPut, "/settings/model", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}