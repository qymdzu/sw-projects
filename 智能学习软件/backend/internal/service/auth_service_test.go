package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/internal/service"
	"smart-learning/pkg/jwt"
)

// mockUserRepo 是 UserRepository 的内存 mock。
type mockUserRepo struct {
	users map[uuid.UUID]*model.User
	byPhone map[string]*model.User
	byEmail map[string]*model.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:   make(map[uuid.UUID]*model.User),
		byPhone: make(map[string]*model.User),
		byEmail: make(map[string]*model.User),
	}
}

func (m *mockUserRepo) Create(_ context.Context, u *model.User) error {
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

func (m *mockUserRepo) GetByID(_ context.Context, id uuid.UUID) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByPhone(_ context.Context, phone string) (*model.User, error) {
	u, ok := m.byPhone[phone]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) Update(_ context.Context, u *model.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) UpdatePassword(_ context.Context, id uuid.UUID, newHash string) error {
	if u, ok := m.users[id]; ok {
		u.PasswordHash = newHash
		return nil
	}
	return repository.ErrNotFound
}

func (m *mockUserRepo) UpdateAvatar(_ context.Context, id uuid.UUID, url string) error {
	if u, ok := m.users[id]; ok {
		u.AvatarURL = &url
		return nil
	}
	return repository.ErrNotFound
}

func newAuthSvc() (service.AuthService, *mockUserRepo) {
	repo := newMockUserRepo()
	mgr := jwt.NewManager("test-secret-1234567890", time.Hour, 24*time.Hour)
	return service.NewAuthService(repo, mgr), repo
}

func TestAuthService_Register_Success(t *testing.T) {
	svc, _ := newAuthSvc()
	resp, err := svc.Register(context.Background(), service.RegisterRequest{
		Name:     "张三",
		Phone:    "13800138000",
		Password: "Abc12345!",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, "张三", resp.User.Name)
	assert.Equal(t, "student", resp.User.Role)
}

func TestAuthService_Register_DuplicatePhone(t *testing.T) {
	svc, _ := newAuthSvc()
	_, err := svc.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)
	_, err = svc.Register(context.Background(), service.RegisterRequest{
		Name: "李四", Phone: "13800138000", Password: "Abc12345!",
	})
	assert.ErrorIs(t, err, service.ErrAccountConflict)
}

func TestAuthService_Register_InvalidPhone(t *testing.T) {
	svc, _ := newAuthSvc()
	_, err := svc.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "12345", Password: "Abc12345!",
	})
	assert.Error(t, err)
}

func TestAuthService_Register_WeakPassword(t *testing.T) {
	svc, _ := newAuthSvc()
	_, err := svc.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "weak",
	})
	assert.Error(t, err)
}

func TestAuthService_Register_NoContact(t *testing.T) {
	svc, _ := newAuthSvc()
	_, err := svc.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Password: "Abc12345!",
	})
	assert.Error(t, err)
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, _ := newAuthSvc()
	_, err := svc.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)

	resp, err := svc.Login(context.Background(), service.LoginRequest{
		Account: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc, _ := newAuthSvc()
	_, err := svc.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)
	_, err = svc.Login(context.Background(), service.LoginRequest{
		Account: "13800138000", Password: "WrongPass1!",
	})
	assert.ErrorIs(t, err, service.ErrPasswordInvalid)
}

func TestAuthService_Login_AccountNotFound(t *testing.T) {
	svc, _ := newAuthSvc()
	_, err := svc.Login(context.Background(), service.LoginRequest{
		Account: "13800138000", Password: "Abc12345!",
	})
	assert.ErrorIs(t, err, service.ErrAccountNotFound)
}

func TestAuthService_Login_ByEmail(t *testing.T) {
	svc, _ := newAuthSvc()
	_, err := svc.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Email: "test@example.com", Password: "Abc12345!",
	})
	require.NoError(t, err)
	resp, err := svc.Login(context.Background(), service.LoginRequest{
		Account: "test@example.com", Password: "Abc12345!",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
}

func TestAuthService_Refresh(t *testing.T) {
	svc, _ := newAuthSvc()
	resp, err := svc.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)

	newResp, err := svc.Refresh(context.Background(), resp.RefreshToken)
	require.NoError(t, err)
	assert.NotEqual(t, resp.AccessToken, newResp.AccessToken)
}

func TestAuthService_Refresh_InvalidToken(t *testing.T) {
	svc, _ := newAuthSvc()
	_, err := svc.Refresh(context.Background(), "invalid-token")
	assert.Error(t, err)
}

func TestMaskUser(t *testing.T) {
	phone := "13800138000"
	email := "test@example.com"
	u := &model.User{
		ID: uuid.New(), Name: "张三",
		Phone: &phone, Email: &email,
	}
	masked := service.MaskUser(u)
	require.NotNil(t, masked)
	assert.Equal(t, "138****8000", *masked.Phone)
	assert.Contains(t, *masked.Email, "@example.com")
	assert.Contains(t, *masked.Email, "****")
}