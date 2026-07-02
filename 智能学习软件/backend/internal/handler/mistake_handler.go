package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"smart-learning/internal/service"
	"smart-learning/pkg/pagination"
	"smart-learning/pkg/response"
)

// MistakeHandler 错题 HTTP 处理器。
type MistakeHandler struct {
	svc service.MistakeService
}

// NewMistakeHandler 构造 MistakeHandler。
func NewMistakeHandler(svc service.MistakeService) *MistakeHandler {
	return &MistakeHandler{svc: svc}
}

// List GET /api/v1/mistakes。
func (h *MistakeHandler) List(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	pg := pagination.Parse(c)
	var kpID *uint64
	if v := c.Query("knowledge_point_id"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			kpID = &n
		}
	}
	var mastered *bool
	if v := c.Query("mastered"); v != "" {
		b := v == "true" || v == "1"
		mastered = &b
	}
	items, total, err := h.svc.List(c.Request.Context(), uid, kpID, mastered, pg.Page, pg.PageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, items, total, pg.Page, pg.PageSize)
}

// Groups GET /api/v1/mistakes/groups。
func (h *MistakeHandler) Groups(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	groups, err := h.svc.GroupByKP(c.Request.Context(), uid)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"groups": groups})
}

// MarkMasteredReq 标记掌握请求。
type MarkMasteredReq struct {
	Mastered bool `json:"mastered"`
}

// MarkMastered PUT /api/v1/mistakes/:id/master。
func (h *MistakeHandler) MarkMastered(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id 格式错误", nil)
		return
	}
	var req MarkMasteredReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	if err := h.svc.MarkMastered(c.Request.Context(), uid, id, req.Mastered); err != nil {
		switch {
		case errors.Is(err, service.ErrResourceMissing):
			response.NotFound(c, "错题不存在")
		case errors.Is(err, service.ErrForbidden):
			response.Forbidden(c, "无权操作该错题")
		default:
			response.ServerError(c, err.Error())
		}
		return
	}
	response.OK(c, nil)
}

// ReviewReq 错题重练请求。
type ReviewReq struct {
	KnowledgePointIDs []uint64 `json:"knowledge_point_ids"`
	Count             int      `json:"count"`
}

// Review POST /api/v1/mistakes/review。
func (h *MistakeHandler) Review(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	var req ReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	resp, err := h.svc.Review(c.Request.Context(), uid, service.ReviewRequest{
		KnowledgePointIDs: req.KnowledgePointIDs,
		Count:             req.Count,
	})
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// Delete DELETE /api/v1/mistakes/:id。
func (h *MistakeHandler) Delete(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "id 格式错误", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uid, id); err != nil {
		switch {
		case errors.Is(err, service.ErrResourceMissing):
			response.NotFound(c, "错题不存在")
		case errors.Is(err, service.ErrForbidden):
			response.Forbidden(c, "无权操作该错题")
		default:
			response.ServerError(c, err.Error())
		}
		return
	}
	response.OK(c, nil)
}