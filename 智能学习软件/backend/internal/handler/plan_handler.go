package handler

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"smart-learning/internal/service"
	"smart-learning/pkg/pagination"
	"smart-learning/pkg/response"
)

// PlanHandler 学习计划 HTTP 处理器。
type PlanHandler struct {
	svc service.PlanService
}

// NewPlanHandler 构造 PlanHandler。
func NewPlanHandler(svc service.PlanService) *PlanHandler {
	return &PlanHandler{svc: svc}
}

// planItemReq 是计划项请求结构。
type planItemReq struct {
	Date             string   `json:"date"`
	KnowledgePointIDs []uint64 `json:"knowledge_point_ids"`
	DurationMin      int      `json:"duration_min"`
	ExercisesCount   int      `json:"exercises_count,omitempty"`
	Goal             string   `json:"goal,omitempty"`
}

// createPlanReq 创建计划请求。
type createPlanReq struct {
	Goal      string        `json:"goal" binding:"required"`
	StartDate string        `json:"start_date" binding:"required"`
	EndDate   string        `json:"end_date" binding:"required"`
	Items     []planItemReq `json:"items"`
}

func parsePlanItems(items []planItemReq) ([]service.PlanItem, error) {
	out := make([]service.PlanItem, 0, len(items))
	for _, it := range items {
		pi := service.PlanItem{
			KnowledgePointIDs: it.KnowledgePointIDs,
			DurationMin:       it.DurationMin,
			ExercisesCount:    it.ExercisesCount,
			Goal:              it.Goal,
		}
		if it.Date != "" {
			_, err := time.Parse("2006-01-02", it.Date)
			if err != nil {
				return nil, errors.New("items.date 格式必须为 YYYY-MM-DD")
			}
			pi.Date = it.Date
		}
		out = append(out, pi)
	}
	return out, nil
}

// Create POST /api/v1/plans。
func (h *PlanHandler) Create(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	var req createPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.BadRequest(c, "start_date 格式错误", nil)
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.BadRequest(c, "end_date 格式错误", nil)
		return
	}
	items, err := parsePlanItems(req.Items)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	p, err := h.svc.Create(c.Request.Context(), uid, service.CreatePlanRequest{
		Goal:      req.Goal,
		StartDate: startDate,
		EndDate:   endDate,
		Items:     items,
	})
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, p)
}

// List GET /api/v1/plans。
func (h *PlanHandler) List(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	pg := pagination.Parse(c)
	status := c.Query("status")
	items, total, err := h.svc.List(c.Request.Context(), uid, status, pg.Page, pg.PageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, items, total, pg.Page, pg.PageSize)
}

// GetByID GET /api/v1/plans/:id。
func (h *PlanHandler) GetByID(c *gin.Context) {
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
	p, err := h.svc.GetByID(c.Request.Context(), uid, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrResourceMissing):
			response.NotFound(c, "计划不存在")
		case errors.Is(err, service.ErrForbidden):
			response.Forbidden(c, "无权访问该计划")
		default:
			response.ServerError(c, err.Error())
		}
		return
	}
	response.OK(c, p)
}

// Update PUT /api/v1/plans/:id。
func (h *PlanHandler) Update(c *gin.Context) {
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
	var req createPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.BadRequest(c, "start_date 格式错误", nil)
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.BadRequest(c, "end_date 格式错误", nil)
		return
	}
	items, err := parsePlanItems(req.Items)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	p, err := h.svc.Update(c.Request.Context(), uid, id, service.CreatePlanRequest{
		Goal:      req.Goal,
		StartDate: startDate,
		EndDate:   endDate,
		Items:     items,
	})
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, p)
}

// aiGenerateReq AI 生成计划请求。
type aiGenerateReq struct {
	Goal             string `json:"goal" binding:"required"`
	StartDate        string `json:"start_date" binding:"required"`
	EndDate          string `json:"end_date" binding:"required"`
	DailyDurationMin int    `json:"daily_duration_min"`
}

// AIGenerate POST /api/v1/plans/ai-generate。
func (h *PlanHandler) AIGenerate(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	var req aiGenerateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.BadRequest(c, "start_date 格式错误", nil)
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.BadRequest(c, "end_date 格式错误", nil)
		return
	}
	p, aiUsed, err := h.svc.AIGenerate(c.Request.Context(), uid, service.AIGenerateRequest{
		Goal:             req.Goal,
		StartDate:        startDate,
		EndDate:          endDate,
		DailyDurationMin: req.DailyDurationMin,
	})
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, gin.H{"plan": p, "ai_used": aiUsed})
}

// checkinReq 打卡请求。
type checkinReq struct {
	Date        string `json:"date" binding:"required"`
	DurationMin int    `json:"duration_min" binding:"required"`
	Status      string `json:"status"`
	Memo        string `json:"memo"`
}

// Checkin POST /api/v1/plans/:id/checkin。
func (h *PlanHandler) Checkin(c *gin.Context) {
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
	var req checkinReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		response.BadRequest(c, "date 格式错误", nil)
		return
	}
	rec, err := h.svc.Checkin(c.Request.Context(), uid, id, service.CheckinRequest{
		Date:        date,
		DurationMin: req.DurationMin,
		Status:      req.Status,
		Memo:        req.Memo,
	})
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, rec)
}

// Delete DELETE /api/v1/plans/:id。
func (h *PlanHandler) Delete(c *gin.Context) {
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
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, nil)
}