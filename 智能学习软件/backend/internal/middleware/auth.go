// Package middleware 提供 HTTP 中间件。
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"smart-learning/pkg/jwt"
	"smart-learning/pkg/response"
)

const (
	// CtxUserID 用户 ID 注入键。
	CtxUserID = "user_id"
	// CtxUserRole 角色注入键。
	CtxUserRole = "role"
)

// JWTAuth 是 JWT 认证中间件（详见 docs/design/系统架构设计.md 4.1）。
func JWTAuth(mgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供 Authorization Header")
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "Authorization Header 格式错误")
			return
		}
		claims, err := mgr.ParseToken(parts[1])
		if err != nil {
			switch err {
			case jwt.ErrTokenExpired:
				response.TokenExpired(c)
			case jwt.ErrTokenInvalid, jwt.ErrTokenWrongType:
				response.TokenInvalid(c)
			default:
				response.Unauthorized(c, err.Error())
			}
			return
		}
		if claims.Type != jwt.TypeAccess {
			response.TokenInvalid(c)
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUserRole, claims.Role)
		c.Next()
	}
}

// RequireRole 角色守卫中间件（详见 API设计 9.3）。
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, ok := c.Get(CtxUserRole)
		if !ok {
			response.Unauthorized(c, "未登录")
			return
		}
		if _, ok := allowed[role.(string)]; !ok {
			response.Forbidden(c, "权限不足")
			return
		}
		c.Next()
	}
}