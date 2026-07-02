package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"smart-learning/internal/service"
)

func TestUserService_GetMe(t *testing.T) {
	auth, repo := newAuthSvc()
	resp, err := auth.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)

	userSvc := service.NewUserService(repo)
	u, err := userSvc.GetMe(context.Background(), resp.User.ID)
	require.NoError(t, err)
	assert.Equal(t, "张三", u.Name)
}

func TestUserService_GetMe_NotFound(t *testing.T) {
	_, repo := newAuthSvc()
	userSvc := service.NewUserService(repo)
	_, err := userSvc.GetMe(context.Background(), uuid.New())
	assert.Error(t, err)
}

func TestUserService_UpdateMe(t *testing.T) {
	auth, repo := newAuthSvc()
	resp, err := auth.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)

	userSvc := service.NewUserService(repo)
	u, err := userSvc.UpdateMe(context.Background(), resp.User.ID, service.UpdateUserRequest{
		Name: "张三（更新）",
	})
	require.NoError(t, err)
	assert.Equal(t, "张三（更新）", u.Name)
}

func TestUserService_ChangePassword(t *testing.T) {
	auth, repo := newAuthSvc()
	resp, err := auth.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)

	userSvc := service.NewUserService(repo)
	err = userSvc.ChangePassword(context.Background(), resp.User.ID, "Abc12345!", "NewPass123!")
	require.NoError(t, err)

	// 旧密码登录失败
	_, err = auth.Login(context.Background(), service.LoginRequest{
		Account: "13800138000", Password: "Abc12345!",
	})
	assert.Error(t, err)
	// 新密码登录成功
	_, err = auth.Login(context.Background(), service.LoginRequest{
		Account: "13800138000", Password: "NewPass123!",
	})
	require.NoError(t, err)
}

func TestUserService_ChangePassword_WrongOldPwd(t *testing.T) {
	auth, repo := newAuthSvc()
	resp, err := auth.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)

	userSvc := service.NewUserService(repo)
	err = userSvc.ChangePassword(context.Background(), resp.User.ID, "WrongOld1!", "NewPass123!")
	assert.ErrorIs(t, err, service.ErrPasswordInvalid)
}

func TestUserService_ChangePassword_Weak(t *testing.T) {
	auth, repo := newAuthSvc()
	resp, err := auth.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)

	userSvc := service.NewUserService(repo)
	err = userSvc.ChangePassword(context.Background(), resp.User.ID, "Abc12345!", "weak")
	assert.Error(t, err)
}

func TestUserService_UpdateAvatar(t *testing.T) {
	auth, repo := newAuthSvc()
	resp, err := auth.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)

	userSvc := service.NewUserService(repo)
	err = userSvc.UpdateAvatar(context.Background(), resp.User.ID, "https://cdn.example.com/a.jpg")
	require.NoError(t, err)

	u, err := userSvc.GetMe(context.Background(), resp.User.ID)
	require.NoError(t, err)
	require.NotNil(t, u.AvatarURL)
	assert.Equal(t, "https://cdn.example.com/a.jpg", *u.AvatarURL)
}

func TestUserService_UpdateAvatar_Empty(t *testing.T) {
	auth, repo := newAuthSvc()
	resp, err := auth.Register(context.Background(), service.RegisterRequest{
		Name: "张三", Phone: "13800138000", Password: "Abc12345!",
	})
	require.NoError(t, err)

	userSvc := service.NewUserService(repo)
	err = userSvc.UpdateAvatar(context.Background(), resp.User.ID, "")
	assert.Error(t, err)
}