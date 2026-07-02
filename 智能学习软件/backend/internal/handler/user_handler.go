package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"smart-learning/internal/middleware"
	"smart-learning/internal/service"
	"smart-learning/pkg/response"
)

// UserHandler 用户相关 HTTP 处理器。
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler 构造 UserHandler。
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// currentUserID 从 gin.Context 提取当前用户 UUID。
func currentUserID(c *gin.Context) (uuid.UUID, error) {
	v, ok := c.Get(middleware.CtxUserID)
	if !ok {
		return uuid.Nil, errors.New("未登录")
	}
	s, _ := v.(string)
	return uuid.Parse(s)
}

// GetMe GET /api/v1/users/me。
func (h *UserHandler) GetMe(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	u, err := h.svc.GetMe(c.Request.Context(), uid)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.OK(c, service.MaskUser(u))
}

// updateMeReqBody 更新个人信息请求。
type updateMeReqBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateMe PUT /api/v1/users/me。
func (h *UserHandler) UpdateMe(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	var req updateMeReqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	u, err := h.svc.UpdateMe(c.Request.Context(), uid, service.UpdateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, service.MaskUser(u))
}

// changePwdReqBody 修改密码请求。
type changePwdReqBody struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangePassword PUT /api/v1/users/me/password。
func (h *UserHandler) ChangePassword(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	var req changePwdReqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	if err := h.svc.ChangePassword(c.Request.Context(), uid, req.OldPassword, req.NewPassword); err != nil {
		switch {
		case errors.Is(err, service.ErrPasswordInvalid):
			response.Unauthorized(c, "原密码错误")
		default:
			response.BadRequest(c, err.Error(), nil)
		}
		return
	}
	response.OK(c, nil)
}

// UpdateAvatar POST /api/v1/users/me/avatar。
// MVP 简化：通过 JSON body 传入 URL（实际项目可改为 multipart/form-data + OSS）。
func (h *UserHandler) UpdateAvatar(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	var body struct {
		AvatarURL string `json:"avatar_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	if err := h.svc.UpdateAvatar(c.Request.Context(), uid, body.AvatarURL); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, gin.H{"avatar_url": body.AvatarURL})
}