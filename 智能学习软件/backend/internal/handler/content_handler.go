package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart-learning/internal/service"
	"smart-learning/pkg/response"
)

// SubjectHandler 科目 HTTP 处理器。
type SubjectHandler struct {
	svc service.SubjectService
}

// NewSubjectHandler 构造 SubjectHandler。
func NewSubjectHandler(svc service.SubjectService) *SubjectHandler {
	return &SubjectHandler{svc: svc}
}

// List GET /api/v1/subjects。
func (h *SubjectHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"items": items})
}

// KnowledgeHandler 知识点 HTTP 处理器。
type KnowledgeHandler struct {
	svc service.KnowledgeService
}

// NewKnowledgeHandler 构造 KnowledgeHandler。
func NewKnowledgeHandler(svc service.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{svc: svc}
}

// Tree GET /api/v1/knowledge-points?subject_id=1。
func (h *KnowledgeHandler) Tree(c *gin.Context) {
	sidStr := c.Query("subject_id")
	if sidStr == "" {
		response.BadRequest(c, "subject_id 不能为空", nil)
		return
	}
	sid, err := strconv.ParseUint(sidStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "subject_id 格式错误", nil)
		return
	}
	tree, err := h.svc.GetTree(c.Request.Context(), sid)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"tree": tree})
}