package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/pkg/hash"
	"smart-learning/pkg/validator"
)

// UpdateUserRequest 更新用户请求。
type UpdateUserRequest struct {
	Name  string
	Email string
}

// UserService 用户服务接口。
type UserService interface {
	GetMe(ctx context.Context, userID uuid.UUID) (*model.User, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*model.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPwd, newPwd string) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, url string) error
}

type userService struct {
	users repository.UserRepository
}

// NewUserService 构造 UserService。
func NewUserService(users repository.UserRepository) UserService {
	return &userService{users: users}
}

func (s *userService) GetMe(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrResourceMissing
		}
		return nil, err
	}
	return u, nil
}

func (s *userService) UpdateMe(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*model.User, error) {
	if req.Name != "" {
		if err := validator.RequireString("姓名", req.Name); err != nil {
			return nil, err
		}
	}
	if req.Email != "" {
		if err := validator.Email(req.Email); err != nil {
			return nil, err
		}
	}
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		e := req.Email
		u.Email = &e
	}
	if err := s.users.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPwd, newPwd string) error {
	if err := validator.Password(newPwd); err != nil {
		return err
	}
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := hash.Verify(u.PasswordHash, oldPwd); err != nil {
		return ErrPasswordInvalid
	}
	hashed, err := hash.Hash(newPwd)
	if err != nil {
		return err
	}
	return s.users.UpdatePassword(ctx, userID, hashed)
}

func (s *userService) UpdateAvatar(ctx context.Context, userID uuid.UUID, url string) error {
	if url == "" {
		return errors.New("头像 URL 不能为空")
	}
	return s.users.UpdateAvatar(ctx, userID, url)
}