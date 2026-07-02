package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"smart-learning/internal/service"
	"smart-learning/pkg/response"
)

// ReportHandler 学习看板 HTTP 处理器。
type ReportHandler struct {
	svc service.ReportService
}

// NewReportHandler 构造 ReportHandler。
func NewReportHandler(svc service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// Summary GET /api/v1/reports/summary。
func (h *ReportHandler) Summary(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	sum, err := h.svc.Summary(c.Request.Context(), uid)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, sum)
}

// Detail GET /api/v1/reports/detail。
func (h *ReportHandler) Detail(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	periodType := c.DefaultQuery("period_type", "weekly")
	startStr := c.Query("period_start")
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		response.BadRequest(c, "period_start 格式错误", nil)
		return
	}
	dto, err := h.svc.Detail(c.Request.Context(), uid, periodType, start)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, dto)
}

// Mastery GET /api/v1/reports/mastery。
func (h *ReportHandler) Mastery(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
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
	dto, err := h.svc.Mastery(c.Request.Context(), uid, sid)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, dto)
}

// Trend GET /api/v1/reports/trend。
func (h *ReportHandler) Trend(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	days := 7
	if v := c.Query("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			days = n
		}
	}
	dto, err := h.svc.Trend(c.Request.Context(), uid, days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, dto)
}