package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
)

// PlanItem 是学习计划 items 中的单日条目。
type PlanItem struct {
	Date             string `json:"date"` // YYYY-MM-DD
	KnowledgePointIDs []uint64 `json:"knowledge_point_ids"`
	DurationMin      int    `json:"duration_min"`
	ExercisesCount   int    `json:"exercises_count,omitempty"`
	Goal             string `json:"goal,omitempty"`
}

// CreatePlanRequest 创建计划请求。
type CreatePlanRequest struct {
	Goal      string
	StartDate time.Time
	EndDate   time.Time
	Items     []PlanItem
}

// AIGenerateRequest AI 生成计划请求。
type AIGenerateRequest struct {
	Goal             string
	StartDate        time.Time
	EndDate          time.Time
	DailyDurationMin int
}

// CheckinRequest 打卡请求。
type CheckinRequest struct {
	Date        time.Time
	DurationMin int
	Status      string // done/skip/partial
	Memo        string
}

// PlanService 学习计划服务接口。
type PlanService interface {
	Create(ctx context.Context, userID uuid.UUID, req CreatePlanRequest) (*model.StudyPlan, error)
	GetByID(ctx context.Context, userID uuid.UUID, id uint64) (*model.StudyPlan, error)
	List(ctx context.Context, userID uuid.UUID, status string, page, pageSize int) ([]model.StudyPlan, int64, error)
	Update(ctx context.Context, userID uuid.UUID, id uint64, req CreatePlanRequest) (*model.StudyPlan, error)
	Delete(ctx context.Context, userID uuid.UUID, id uint64) error
	AIGenerate(ctx context.Context, userID uuid.UUID, req AIGenerateRequest) (*model.StudyPlan, bool, error)
	Checkin(ctx context.Context, userID uuid.UUID, planID uint64, req CheckinRequest) (*model.StudyRecord, error)
}

type planService struct {
	plans repository.PlanRepository
}

// NewPlanService 构造 PlanService。
func NewPlanService(plans repository.PlanRepository) PlanService {
	return &planService{plans: plans}
}

func (s *planService) validateDates(start, end time.Time) error {
	if end.Before(start) {
		return errors.New("结束日期不能早于开始日期")
	}
	return nil
}

func marshalItems(items []PlanItem) (datatypes.JSON, error) {
	b, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(b), nil
}

func (s *planService) Create(ctx context.Context, userID uuid.UUID, req CreatePlanRequest) (*model.StudyPlan, error) {
	if req.Goal == "" {
		return nil, errors.New("学习目标不能为空")
	}
	if err := s.validateDates(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}
	itemsJSON, err := marshalItems(req.Items)
	if err != nil {
		return nil, err
	}
	p := &model.StudyPlan{
		UserID:    userID,
		Goal:      req.Goal,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Items:     itemsJSON,
		Status:    "active",
	}
	if err := s.plans.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *planService) GetByID(ctx context.Context, userID uuid.UUID, id uint64) (*model.StudyPlan, error) {
	p, err := s.plans.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrResourceMissing
		}
		return nil, err
	}
	if p.UserID != userID {
		return nil, ErrForbidden
	}
	return p, nil
}

func (s *planService) List(ctx context.Context, userID uuid.UUID, status string, page, pageSize int) ([]model.StudyPlan, int64, error) {
	return s.plans.ListByUser(ctx, userID, status, page, pageSize)
}

func (s *planService) Update(ctx context.Context, userID uuid.UUID, id uint64, req CreatePlanRequest) (*model.StudyPlan, error) {
	p, err := s.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if err := s.validateDates(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}
	p.Goal = req.Goal
	p.StartDate = req.StartDate
	p.EndDate = req.EndDate
	itemsJSON, err := marshalItems(req.Items)
	if err != nil {
		return nil, err
	}
	p.Items = itemsJSON
	if err := s.plans.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *planService) Delete(ctx context.Context, userID uuid.UUID, id uint64) error {
	if _, err := s.GetByID(ctx, userID, id); err != nil {
		return err
	}
	return s.plans.Delete(ctx, id)
}

// AIGenerate 使用规则降级引擎生成学习计划。
// MVP 策略：
//   - 按 start_date~end_date 每天生成 item
//   - duration_min = daily_duration_min
//   - knowledge_point_ids 默认空（让用户后续填充或 AI 补充）
func (s *planService) AIGenerate(ctx context.Context, userID uuid.UUID, req AIGenerateRequest) (*model.StudyPlan, bool, error) {
	if req.Goal == "" {
		return nil, false, errors.New("学习目标不能为空")
	}
	if err := s.validateDates(req.StartDate, req.EndDate); err != nil {
		return nil, false, err
	}
	if req.DailyDurationMin <= 0 {
		req.DailyDurationMin = 60
	}

	items := []PlanItem{}
	days := int(req.EndDate.Sub(req.StartDate).Hours()/24) + 1
	if days > 365 {
		days = 365
	}
	for i := 0; i < days; i++ {
		d := req.StartDate.AddDate(0, 0, i)
		items = append(items, PlanItem{
			Date:             d.Format("2006-01-02"),
			KnowledgePointIDs: []uint64{},
			DurationMin:      req.DailyDurationMin,
			ExercisesCount:   10,
			Goal:             req.Goal,
		})
	}
	itemsJSON, err := marshalItems(items)
	if err != nil {
		return nil, false, err
	}
	p := &model.StudyPlan{
		UserID:      userID,
		Goal:        req.Goal,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Items:       itemsJSON,
		Status:      "active",
		AIGenerated: true,
	}
	if err := s.plans.Create(ctx, p); err != nil {
		return nil, false, err
	}
	// MVP 未对接真实 LLM，直接使用规则降级，返回 ai_used=false
	return p, false, nil
}

func (s *planService) Checkin(ctx context.Context, userID uuid.UUID, planID uint64, req CheckinRequest) (*model.StudyRecord, error) {
	if _, err := s.GetByID(ctx, userID, planID); err != nil {
		return nil, err
	}
	if req.DurationMin <= 0 {
		return nil, errors.New("学习时长必须大于 0")
	}
	status := req.Status
	if status == "" {
		status = "done"
	}
	switch status {
	case "done", "skip", "partial":
	default:
		return nil, errors.New("状态取值必须为 done/skip/partial")
	}
	var memoPtr *string
	if req.Memo != "" {
		m := req.Memo
		memoPtr = &m
	}
	rec := &model.StudyRecord{
		UserID:      userID,
		PlanID:      planID,
		Date:        req.Date,
		DurationMin: req.DurationMin,
		Status:      status,
		Memo:        memoPtr,
	}
	if err := s.plans.CreateRecord(ctx, rec); err != nil {
		return nil, err
	}
	return rec, nil
}