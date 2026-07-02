package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"smart-learning/internal/handler"
	"smart-learning/internal/middleware"
	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/internal/service"
	"smart-learning/pkg/jwt"
)

// ----- 本地 mock for handler tests -----

type hMockUserRepo struct {
	users   map[uuid.UUID]*model.User
	byPhone map[string]*model.User
	byEmail map[string]*model.User
}

func newHMockUserRepo() *hMockUserRepo {
	return &hMockUserRepo{
		users:   make(map[uuid.UUID]*model.User),
		byPhone: make(map[string]*model.User),
		byEmail: make(map[string]*model.User),
	}
}

func (m *hMockUserRepo) Create(_ context.Context, u *model.User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	m.users[u.ID] = u
	if u.Phone != nil {
		m.byPhone[*u.Phone] = u
	}
	if u.Email != nil {
		m.byEmail[*u.Email] = u
	}
	return nil
}

func (m *hMockUserRepo) GetByID(_ context.Context, id uuid.UUID) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errNotFound
	}
	return u, nil
}

func (m *hMockUserRepo) GetByPhone(_ context.Context, phone string) (*model.User, error) {
	u, ok := m.byPhone[phone]
	if !ok {
		return nil, errNotFound
	}
	return u, nil
}

func (m *hMockUserRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return nil, errNotFound
	}
	return u, nil
}

func (m *hMockUserRepo) Update(_ context.Context, u *model.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *hMockUserRepo) UpdatePassword(_ context.Context, id uuid.UUID, newHash string) error {
	if u, ok := m.users[id]; ok {
		u.PasswordHash = newHash
		return nil
	}
	return errNotFound
}

func (m *hMockUserRepo) UpdateAvatar(_ context.Context, id uuid.UUID, url string) error {
	if u, ok := m.users[id]; ok {
		u.AvatarURL = &url
		return nil
	}
	return errNotFound
}

var errNotFound = repository.ErrNotFound

// ----- 路由构造 -----

func setupAuthRouter(t *testing.T) (*gin.Engine, *hMockUserRepo) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := newHMockUserRepo()
	mgr := jwt.NewManager("test-secret-1234567890", time.Hour, 24*time.Hour)
	svc := service.NewAuthService(repo, mgr)
	h := handler.NewAuthHandler(svc)
	api := r.Group("/api/v1/auth")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
		api.POST("/refresh", h.Refresh)
	}
	return r, repo
}

func setupUserRouter(t *testing.T) (*gin.Engine, *hMockUserRepo) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := newHMockUserRepo()
	userSvc := service.NewUserService(repo)
	mgr := jwt.NewManager("test-secret-1234567890", time.Hour, 24*time.Hour)
	authSvc := service.NewAuthService(repo, mgr)

	hUser := handler.NewUserHandler(userSvc)
	hAuth := handler.NewAuthHandler(authSvc)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", hAuth.Register)
		secured := api.Group("")
		secured.Use(middleware.JWTAuth(mgr))
		secured.GET("/users/me", hUser.GetMe)
		secured.PUT("/users/me", hUser.UpdateMe)
		secured.PUT("/users/me/password", hUser.ChangePassword)
		secured.POST("/users/me/avatar", hUser.UpdateAvatar)
	}
	return r, repo
}

// ----- 测试 -----

func TestAuthHandler_Register_Success(t *testing.T) {
	r, _ := setupAuthRouter(t)
	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(0), resp["code"])
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
	user := data["user"].(map[string]interface{})
	assert.Equal(t, "138****8000", user["phone"])
}

func TestAuthHandler_Register_BadRequest(t *testing.T) {
	r, _ := setupAuthRouter(t)
	body := `{"name":"","phone":"","password":""}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_Conflict(t *testing.T) {
	r, _ := setupAuthRouter(t)
	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req1 := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	require.Equal(t, http.StatusCreated, w1.Code)

	body2 := `{"name":"李四","phone":"13800138000","password":"Abc12345!"}`
	req2 := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestAuthHandler_Login_AccountNotFound(t *testing.T) {
	r, _ := setupAuthRouter(t)
	body := `{"account":"13900000000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	r, _ := setupAuthRouter(t)
	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	body2 := `{"account":"13800138000","password":"WrongPass1!"}`
	req2 := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	r, _ := setupAuthRouter(t)
	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	body2 := `{"account":"13800138000","password":"Abc12345!"}`
	req2 := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestAuthHandler_Refresh_Success(t *testing.T) {
	r, _ := setupAuthRouter(t)
	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	refreshToken := data["refresh_token"].(string)

	body2, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
	req2 := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestAuthHandler_Refresh_Invalid(t *testing.T) {
	r, _ := setupAuthRouter(t)
	body2 := `{"refresh_token":"invalid.token.value"}`
	req2 := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestUserHandler_GetMe_Unauthorized(t *testing.T) {
	r, _ := setupUserRouter(t)
	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_GetMe_Success(t *testing.T) {
	r, _ := setupUserRouter(t)

	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	accessToken := resp["data"].(map[string]interface{})["access_token"].(string)

	req2 := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	req2.Header.Set("Authorization", "Bearer "+accessToken)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestUserHandler_UpdateMe_Success(t *testing.T) {
	r, _ := setupUserRouter(t)

	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	accessToken := resp["data"].(map[string]interface{})["access_token"].(string)

	body2 := `{"name":"张三（新）"}`
	req2 := httptest.NewRequest("PUT", "/api/v1/users/me", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+accessToken)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestUserHandler_ChangePassword_Success(t *testing.T) {
	r, _ := setupUserRouter(t)

	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	accessToken := resp["data"].(map[string]interface{})["access_token"].(string)

	body2 := `{"old_password":"Abc12345!","new_password":"NewPass123!"}`
	req2 := httptest.NewRequest("PUT", "/api/v1/users/me/password", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+accessToken)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestUserHandler_ChangePassword_WrongOld(t *testing.T) {
	r, _ := setupUserRouter(t)

	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	accessToken := resp["data"].(map[string]interface{})["access_token"].(string)

	body2 := `{"old_password":"WrongOld1!","new_password":"NewPass123!"}`
	req2 := httptest.NewRequest("PUT", "/api/v1/users/me/password", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+accessToken)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestUserHandler_UpdateAvatar(t *testing.T) {
	r, _ := setupUserRouter(t)

	body := `{"name":"张三","phone":"13800138000","password":"Abc12345!"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	accessToken := resp["data"].(map[string]interface{})["access_token"].(string)

	body2 := `{"avatar_url":"https://cdn.example.com/a.jpg"}`
	req2 := httptest.NewRequest("POST", "/api/v1/users/me/avatar", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+accessToken)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

// ----- 中间件测试 -----

func TestJWTAuth_NoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mgr := jwt.NewManager("test-secret", time.Hour, 24*time.Hour)
	r.GET("/x", middleware.JWTAuth(mgr), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req := httptest.NewRequest("GET", "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mgr := jwt.NewManager("test-secret", time.Hour, 24*time.Hour)
	r.GET("/x", middleware.JWTAuth(mgr), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req := httptest.NewRequest("GET", "/x", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_BadHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mgr := jwt.NewManager("test-secret", time.Hour, 24*time.Hour)
	r.GET("/x", middleware.JWTAuth(mgr), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req := httptest.NewRequest("GET", "/x", nil)
	req.Header.Set("Authorization", "NotBearer xyz")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_Expired(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := jwt.NewManager("test-secret", time.Millisecond, 24*time.Hour)
	pair, err := mgr.GenerateTokenPair(uuid.NewString(), "student")
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	r := gin.New()
	r.GET("/x", middleware.JWTAuth(mgr), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req := httptest.NewRequest("GET", "/x", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireRole_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := jwt.NewManager("test-secret", time.Hour, 24*time.Hour)
	pair, err := mgr.GenerateTokenPair(uuid.NewString(), "student")
	require.NoError(t, err)

	r := gin.New()
	r.GET("/admin", middleware.JWTAuth(mgr), middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRole_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mgr := jwt.NewManager("test-secret", time.Hour, 24*time.Hour)
	pair, err := mgr.GenerateTokenPair(uuid.NewString(), "admin")
	require.NoError(t, err)

	r := gin.New()
	r.GET("/admin", middleware.JWTAuth(mgr), middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORS_Preflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.CORS())
	r.POST("/x", func(c *gin.Context) { c.JSON(200, nil) })
	req := httptest.NewRequest("OPTIONS", "/x", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}