// Package jwt 提供基于 golang-jwt/jwt v5 的 Token 生成与校验能力。
//
// 支持 access_token (短期) + refresh_token (长期) 双 Token 机制，
// 详见 docs/design/系统架构设计.md 第 4.1 节。
package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Token 类型常量。
const (
	TypeAccess  = "access"
	TypeRefresh = "refresh"
)

// 业务错误码（与 API 设计 1.3 节保持一致）。
var (
	ErrTokenInvalid = errors.New("token 无效")
	ErrTokenExpired = errors.New("token 已过期")
	ErrTokenWrongType = errors.New("token 类型不匹配")
)

// Claims 是自定义的 JWT Claims。
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"`
	jwtv5.RegisteredClaims
}

// TokenPair 是 access + refresh 双 Token 对。
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // access_token 剩余秒数
}

// Manager 是 JWT 管理器。
type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

// NewManager 构造一个 Manager。
func NewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     "smart-learning",
	}
}

// GenerateTokenPair 为指定用户生成 access + refresh 双 Token。
func (m *Manager) GenerateTokenPair(userID, role string) (*TokenPair, error) {
	now := time.Now()
	
	// 生成唯一 token ID，确保每次生成的 token 不同
	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	tokenID := hex.EncodeToString(idBytes)
	
	accessClaims := Claims{
		UserID: userID,
		Role:   role,
		Type:   TypeAccess,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ID:        tokenID,
			Issuer:    m.issuer,
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(m.accessTTL)),
		},
	}
	refreshClaims := Claims{
		UserID: userID,
		Role:   role,
		Type:   TypeRefresh,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ID:        tokenID,
			Issuer:    m.issuer,
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(m.refreshTTL)),
		},
	}

	accessToken := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, accessClaims)
	refreshToken := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, refreshClaims)

	accessStr, err := accessToken.SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("签名 access_token 失败: %w", err)
	}
	refreshStr, err := refreshToken.SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("签名 refresh_token 失败: %w", err)
	}
	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		ExpiresIn:    int64(m.accessTTL.Seconds()),
	}, nil
}

// ParseToken 解析并校验 token 签名与有效期。
func (m *Manager) ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwtv5.ParseWithClaims(tokenStr, claims, func(t *jwtv5.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwtv5.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	if !tok.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

// RefreshToken 用 refresh_token 换新 token 对。
func (m *Manager) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := m.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}
	if claims.Type != TypeRefresh {
		return nil, ErrTokenWrongType
	}
	return m.GenerateTokenPair(claims.UserID, claims.Role)
}