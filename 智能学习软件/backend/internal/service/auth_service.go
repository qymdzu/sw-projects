package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/pkg/hash"
	"smart-learning/pkg/jwt"
	"smart-learning/pkg/validator"
)

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Name     string
	Phone    string
	Email    string
	Password string
	Role     string // 可选，默认 student
}

// LoginRequest 登录请求。
type LoginRequest struct {
	Account  string // 支持手机号或邮箱
	Password string
}

// AuthResponse 认证响应（与 API设计 2.1 一致）。
type AuthResponse struct {
	User         *model.User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// AuthService 认证服务接口。
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error)
}

type authService struct {
	users repository.UserRepository
	jwt   *jwt.Manager
}

// NewAuthService 构造 AuthService。
func NewAuthService(users repository.UserRepository, mgr *jwt.Manager) AuthService {
	return &authService{users: users, jwt: mgr}
}

// maskPhone 脱敏手机号 138****8000。
func maskPhone(p string) string {
	if len(p) < 7 {
		return p
	}
	return p[:3] + "****" + p[len(p)-4:]
}

// maskEmail 脱敏邮箱 z****n@example.com。
func maskEmail(e string) string {
	at := strings.Index(e, "@")
	if at <= 1 {
		return e
	}
	return e[:1] + "****" + e[at-1:] // z****n@example.com
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	if err := validator.RequireString("姓名", req.Name); err != nil {
		return nil, err
	}
	if req.Phone == "" && req.Email == "" {
		return nil, errors.New("手机号和邮箱至少填写一个")
	}
	if req.Phone != "" {
		if err := validator.Phone(req.Phone); err != nil {
			return nil, err
		}
	}
	if req.Email != "" {
		if err := validator.Email(req.Email); err != nil {
			return nil, err
		}
	}
	if err := validator.Password(req.Password); err != nil {
		return nil, err
	}
	role := req.Role
	if role == "" {
		role = "student"
	}
	if err := validator.Role(role); err != nil {
		return nil, err
	}

	// 唯一性校验
	if req.Phone != "" {
		if existing, _ := s.users.GetByPhone(ctx, req.Phone); existing != nil {
			return nil, ErrAccountConflict
		}
	}
	if req.Email != "" {
		if existing, _ := s.users.GetByEmail(ctx, req.Email); existing != nil {
			return nil, ErrAccountConflict
		}
	}

	hashed, err := hash.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Name:         req.Name,
		PasswordHash: hashed,
		Role:         role,
	}
	if req.Phone != "" {
		p := req.Phone
		user.Phone = &p
	}
	if req.Email != "" {
		e := req.Email
		user.Email = &e
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	tokens, err := s.jwt.GenerateTokenPair(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	if err := validator.RequireString("账号", req.Account); err != nil {
		return nil, err
	}
	if err := validator.RequireString("密码", req.Password); err != nil {
		return nil, err
	}

	var user *model.User
	var err error
	if validator.Phone(req.Account) == nil {
		user, err = s.users.GetByPhone(ctx, req.Account)
	} else if validator.Email(req.Account) == nil {
		user, err = s.users.GetByEmail(ctx, req.Account)
	} else {
		return nil, errors.New("账号必须是手机号或邮箱")
	}
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	if err := hash.Verify(user.PasswordHash, req.Password); err != nil {
		return nil, ErrPasswordInvalid
	}

	tokens, err := s.jwt.GenerateTokenPair(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	tokens, err := s.jwt.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	claims, err := s.jwt.ParseToken(tokens.AccessToken)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}
	user, err := s.users.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

// MaskUser 对用户敏感信息脱敏，供 handler 序列化使用。
func MaskUser(u *model.User) *model.User {
	if u == nil {
		return nil
	}
	cp := *u
	if u.Phone != nil {
		p := maskPhone(*u.Phone)
		cp.Phone = &p
	}
	if u.Email != nil {
		e := maskEmail(*u.Email)
		cp.Email = &e
	}
	return &cp
}