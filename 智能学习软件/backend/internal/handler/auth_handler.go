// Package handler 提供 HTTP 处理器层。
//
// 仅负责参数绑定、调用 service、统一响应包装；
// 不包含业务逻辑，符合 docs/design/目录结构.md 第 4.2 节。
package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"smart-learning/internal/service"
	"smart-learning/pkg/response"
)

// AuthHandler 是认证相关 HTTP 处理器。
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler 构造 AuthHandler。
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// registerReqBody 注册请求 body。
type registerReqBody struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

// Register POST /api/v1/auth/register。
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	resp, err := h.svc.Register(c.Request.Context(), service.RegisterRequest{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		if errors.Is(err, service.ErrAccountConflict) {
			response.Conflict(c, "账号已存在")
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, gin.H{
		"user":          service.MaskUser(resp.User),
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
	})
}

// loginReqBody 登录请求 body。
type loginReqBody struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login POST /api/v1/auth/login。
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	resp, err := h.svc.Login(c.Request.Context(), service.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccountNotFound), errors.Is(err, service.ErrPasswordInvalid):
			response.Unauthorized(c, "账号或密码错误")
		default:
			response.BadRequest(c, err.Error(), nil)
		}
		return
	}
	response.OK(c, gin.H{
		"user":          service.MaskUser(resp.User),
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
	})
}

// refreshReqBody 刷新 token 请求 body。
type refreshReqBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh POST /api/v1/auth/refresh。
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshReqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	resp, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrResourceMissing):
			response.NotFound(c, "用户不存在")
		default:
			response.TokenInvalid(c)
		}
		return
	}
	response.OK(c, gin.H{
		"user":          service.MaskUser(resp.User),
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
	})
}