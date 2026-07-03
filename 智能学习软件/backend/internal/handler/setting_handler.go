// Package handler 提供 HTTP 处理器层。

package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"smart-learning/internal/middleware"
	"smart-learning/internal/service"
	"smart-learning/pkg/response"
)

// SettingHandler 是模型配置相关 HTTP 处理器。
type SettingHandler struct {
	svc service.SettingService
}

// NewSettingHandler 构造 SettingHandler。
func NewSettingHandler(svc service.SettingService) *SettingHandler {
	return &SettingHandler{svc: svc}
}

// saveReqBody POST /api/v1/settings/model 请求体。
type saveReqBody struct {
	Provider    string         `json:"provider" binding:"required"`
	APIEndpoint string         `json:"api_endpoint" binding:"required"`
	APIKey      string         `json:"api_key" binding:"required"`
	Model       string         `json:"model" binding:"required"`
	ExtraConfig map[string]any `json:"extra_config"`
}

// activateReqBody PUT /api/v1/settings/model 请求体。
type activateReqBody struct {
	Provider string `json:"provider" binding:"required"`
}

// CreateOrUpdate POST /api/v1/settings/model。
//
// 业务语义：保存（或覆盖）一条 provider 配置，并把它设为当前 default。
// 出参 api_key 用 "***" 掩码，避免在响应中泄露刚保存的明文（保存后前端会回显自己填的）。
func (h *SettingHandler) CreateOrUpdate(c *gin.Context) {
	var req saveReqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	uidVal, ok := c.Get(middleware.CtxUserID)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	uidStr, ok := uidVal.(string)
	if !ok || uidStr == "" {
		response.Unauthorized(c, "无效的用户身份")
		return
	}
	userID, err := uuid.Parse(uidStr)
	if err != nil {
		response.Unauthorized(c, "无效的用户 ID")
		return
	}

	dto, err := h.svc.Save(c.Request.Context(), userID, service.SaveSettingRequest{
		Provider:    req.Provider,
		APIEndpoint: req.APIEndpoint,
		APIKey:      req.APIKey,
		Model:       req.Model,
		ExtraConfig: req.ExtraConfig,
	})
	if err != nil {
		writeSettingError(c, err)
		return
	}
	response.OK(c, gin.H{
		"provider":     dto.Provider,
		"api_endpoint": dto.APIEndpoint,
		"api_key":      "***", // 创建/更新场景不回显明文
		"model":        dto.Model,
		"extra_config": dto.ExtraConfig,
		"is_default":   dto.IsDefault,
		"updated_at":   dto.UpdatedAt,
	})
}

// GetActive GET /api/v1/settings/model —— 取当前 default 配置。
//
// 未配置时返回 200 + data:null（前端据此判断是否需要展示配置引导）。
func (h *SettingHandler) GetActive(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	dto, err := h.svc.GetActive(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrSettingNotFound) {
			response.OK(c, gin.H{"setting": nil, "message": "尚未配置模型"})
			return
		}
		writeSettingError(c, err)
		return
	}
	response.OK(c, gin.H{
		"provider":     dto.Provider,
		"api_endpoint": dto.APIEndpoint,
		"api_key":      dto.APIKey, // 主动查询场景，回明文
		"model":        dto.Model,
		"extra_config": dto.ExtraConfig,
		"is_default":   dto.IsDefault,
		"updated_at":   dto.UpdatedAt,
	})
}

// Update PUT /api/v1/settings/model —— 切换 default 到已存在的 provider。
func (h *SettingHandler) Update(c *gin.Context) {
	var req activateReqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	if err := h.svc.Activate(c.Request.Context(), userID, req.Provider); err != nil {
		if errors.Is(err, service.ErrSettingNotFound) {
			response.NotFound(c, "配置不存在，请先创建")
			return
		}
		writeSettingError(c, err)
		return
	}
	response.OK(c, gin.H{
		"provider":   req.Provider,
		"is_default": true,
	})
}

// Delete DELETE /api/v1/settings/model?provider=xxx —— 删除某 provider 配置。
func (h *SettingHandler) Delete(c *gin.Context) {
	provider := c.Query("provider")
	if provider == "" {
		response.BadRequest(c, "provider 参数必填", nil)
		return
	}
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), userID, provider); err != nil {
		if errors.Is(err, service.ErrSettingNotFound) {
			response.NotFound(c, "配置不存在")
			return
		}
		writeSettingError(c, err)
		return
	}
	c.JSON(204, nil)
}

// getUserID 从 gin.Context 取出用户 UUID（重复模式抽出来便于单测）。
func getUserID(c *gin.Context) (uuid.UUID, bool) {
	uidVal, ok := c.Get(middleware.CtxUserID)
	if !ok {
		response.Unauthorized(c, "未登录")
		return uuid.Nil, false
	}
	uidStr, ok := uidVal.(string)
	if !ok || uidStr == "" {
		response.Unauthorized(c, "无效的用户身份")
		return uuid.Nil, false
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		response.Unauthorized(c, "无效的用户 ID")
		return uuid.Nil, false
	}
	return uid, true
}

// writeSettingError 把 service 错误统一映射到 HTTP 响应。
//
// 注意：绝不把 ErrInternalCrypto 等底层错误细节回传给客户端（P1-03 一致策略）。
func writeSettingError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnsupportedProvider),
		errors.Is(err, service.ErrSettingParamInvalid),
		errors.Is(err, service.ErrEndpointInvalid),
		errors.Is(err, service.ErrEndpointUnsafe):
		response.BadRequest(c, err.Error(), nil)
	default:
		// 其他未知错误（含 ErrInternalCrypto）一律返回 500 通用消息
		response.ServerError(c, "服务器内部错误")
	}
}