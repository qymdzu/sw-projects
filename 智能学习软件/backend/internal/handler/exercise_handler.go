package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/internal/service"
	"smart-learning/pkg/pagination"
	"smart-learning/pkg/response"
)

// ExerciseHandler 练习相关 HTTP 处理器。
type ExerciseHandler struct {
	svc service.ExerciseService
}

// NewExerciseHandler 构造 ExerciseHandler。
func NewExerciseHandler(svc service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{svc: svc}
}

// SafeQuestion 是题目安全视图（不返回答案和解析）。
type SafeQuestion struct {
	ID               uint64                 `json:"id"`
	SubjectID        uint64                 `json:"subject_id"`
	KnowledgePointID uint64                 `json:"knowledge_point_id"`
	Type             string                 `json:"type"`
	Difficulty       int                    `json:"difficulty"`
	Content          map[string]interface{} `json:"content"`
	Options          []map[string]string    `json:"options,omitempty"`
	KnowledgePoint   *kpInfo                `json:"knowledge_point,omitempty"`
}

type kpInfo struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

func toSafeQuestions(items []model.Question) []SafeQuestion {
	out := make([]SafeQuestion, 0, len(items))
	for _, q := range items {
		sq := SafeQuestion{
			ID:               q.ID,
			SubjectID:        q.SubjectID,
			KnowledgePointID: q.KnowledgePointID,
			Type:             q.Type,
			Difficulty:       q.Difficulty,
			Content:          map[string]interface{}{},
		}
		// 解析 Content/Options JSON
		if len(q.Content) > 0 {
			var c map[string]interface{}
			if err := jsonUnmarshal(q.Content, &c); err == nil {
				sq.Content = c
			}
		}
		if len(q.Options) > 0 {
			var opts []map[string]string
			if err := jsonUnmarshal(q.Options, &opts); err == nil {
				sq.Options = opts
			}
		}
		if q.KnowledgePoint != nil {
			sq.KnowledgePoint = &kpInfo{
				ID:    q.KnowledgePoint.ID,
				Name:  q.KnowledgePoint.Name,
				Level: q.KnowledgePoint.Level,
			}
		}
		out = append(out, sq)
	}
	return out
}

// List GET /api/v1/exercises。
func (h *ExerciseHandler) List(c *gin.Context) {
	pg := pagination.Parse(c)
	filter := repository.QuestionFilter{
		Type:     c.Query("type"),
		Page:     pg.Page,
		PageSize: pg.PageSize,
	}
	if sid := c.Query("subject_id"); sid != "" {
		if v, err := strconv.ParseUint(sid, 10, 64); err == nil {
			filter.SubjectID = &v
		}
	}
	if kp := c.Query("knowledge_point_id"); kp != "" {
		if v, err := strconv.ParseUint(kp, 10, 64); err == nil {
			filter.KnowledgePointID = &v
		}
	}
	if d := c.Query("difficulty"); d != "" {
		if v, err := strconv.Atoi(d); err == nil {
			filter.Difficulty = &v
		}
	}
	items, total, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, toSafeQuestions(items), total, pg.Page, pg.PageSize)
}

// Random GET /api/v1/exercises/random。
func (h *ExerciseHandler) Random(c *gin.Context) {
	count := 10
	if v, err := strconv.Atoi(c.Query("count")); err == nil && v > 0 {
		count = v
	}
	var subjectID, kpID *uint64
	if sid := c.Query("subject_id"); sid != "" {
		if v, err := strconv.ParseUint(sid, 10, 64); err == nil {
			subjectID = &v
		}
	}
	if kp := c.Query("knowledge_point_id"); kp != "" {
		if v, err := strconv.ParseUint(kp, 10, 64); err == nil {
			kpID = &v
		}
	}
	items, err := h.svc.Random(c.Request.Context(), subjectID, kpID, count)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"items": toSafeQuestions(items)})
}

// submitReq 提交答案请求。
type submitReq struct {
	QuestionID  uint64 `json:"question_id" binding:"required"`
	Answer      string `json:"answer" binding:"required"`
	DurationSec int    `json:"duration_sec"`
}

// Submit POST /api/v1/exercises/submit。
func (h *ExerciseHandler) Submit(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	var req submitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
		return
	}
	resp, err := h.svc.Submit(c.Request.Context(), uid, service.SubmitRequest{
		QuestionID:  req.QuestionID,
		Answer:      req.Answer,
		DurationSec: req.DurationSec,
	})
	if err != nil {
		if errors.Is(err, service.ErrResourceMissing) {
			response.NotFound(c, "题目不存在")
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, resp)
}

// Recommend GET /api/v1/exercises/recommend。
func (h *ExerciseHandler) Recommend(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	count := 10
	if v, err := strconv.Atoi(c.Query("count")); err == nil && v > 0 {
		count = v
	}
	items, err := h.svc.Recommend(c.Request.Context(), uid, count)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{
		"items":              toSafeQuestions(items),
		"recommend_strategy": "weak_point",
		"ai_used":            false,
	})
}

// ByKnowledgePoint GET /api/v1/exercises/kp/:kp_id。
func (h *ExerciseHandler) ByKnowledgePoint(c *gin.Context) {
	kpID, err := strconv.ParseUint(c.Param("kp_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "kp_id 格式错误", nil)
		return
	}
	pg := pagination.Parse(c)
	items, total, err := h.svc.ByKnowledgePoint(c.Request.Context(), kpID, pg.Page, pg.PageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, toSafeQuestions(items), total, pg.Page, pg.PageSize)
}

// History GET /api/v1/exercises/history。
func (h *ExerciseHandler) History(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	pg := pagination.Parse(c)
	items, total, err := h.svc.History(c.Request.Context(), uid, pg.Page, pg.PageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, items, total, pg.Page, pg.PageSize)
}