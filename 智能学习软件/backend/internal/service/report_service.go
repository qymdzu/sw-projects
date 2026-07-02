package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"smart-learning/internal/repository"
)

// SummaryDTO 学习概览 DTO。
type SummaryDTO struct {
	TodayDurationMin    int     `json:"today_duration_min"`
	TotalDurationMin    int     `json:"total_duration_min"`
	TotalExercises      int64   `json:"total_exercises"`
	OverallCorrectRate  float64 `json:"overall_correct_rate"`
	StreakDays          int     `json:"streak_days"`
	ActivePlanCount     int64   `json:"active_plan_count"`
	UnmasteredMistakes  int64   `json:"unmastered_mistakes"`
}

// DetailDTO 学习报告详情 DTO。
type DetailDTO struct {
	Period          PeriodInfo   `json:"period"`
	TotalDurationMin int         `json:"total_duration_min"`
	TotalExercises  int64        `json:"total_exercises"`
	CorrectRate     float64      `json:"correct_rate"`
	MasteryByKP     []KPStatDTO  `json:"mastery_by_kp"`
	StreakDays      int          `json:"streak_days"`
	DailyStats      []DailyStat  `json:"daily_stats"`
}

// PeriodInfo 周期信息。
type PeriodInfo struct {
	Type  string `json:"type"`
	Start string `json:"start"`
	End   string `json:"end"`
}

// KPStatDTO 知识点统计。
type KPStatDTO struct {
	KPID uint64  `json:"kp_id"`
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}

// DailyStat 每日统计。
type DailyStat struct {
	Date         string  `json:"date"`
	DurationMin  int     `json:"duration_min"`
	Exercises    int64   `json:"exercises"`
	CorrectRate  float64 `json:"correct_rate"`
}

// MasteryDTO 掌握度 DTO。
type MasteryDTO struct {
	Subject         SubjectInfo     `json:"subject"`
	KnowledgePoints []MasteryKPDTO  `json:"knowledge_points"`
}

// SubjectInfo 科目信息。
type SubjectInfo struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// MasteryKPDTO 知识点掌握度。
type MasteryKPDTO struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	Level  int     `json:"level"`
	Rate   float64 `json:"rate"`
	Status string  `json:"status"` // mastered/learning/weak
}

// TrendDTO 趋势 DTO。
type TrendDTO struct {
	Days        int         `json:"days"`
	DailyPoints []DailyStat `json:"daily_points"`
}

// ReportService 学习报告服务接口。
type ReportService interface {
	Summary(ctx context.Context, userID uuid.UUID) (*SummaryDTO, error)
	Detail(ctx context.Context, userID uuid.UUID, periodType string, start time.Time) (*DetailDTO, error)
	Mastery(ctx context.Context, userID uuid.UUID, subjectID uint64) (*MasteryDTO, error)
	Trend(ctx context.Context, userID uuid.UUID, days int) (*TrendDTO, error)
}

type reportService struct {
	plans     repository.PlanRepository
	exercises repository.ExerciseRepository
	mistakes  repository.MistakeRepository
	subjects  repository.SubjectRepository
	kpRepo    repository.KnowledgeRepository
	reports   repository.ReportRepository
}

// NewReportService 构造 ReportService。
func NewReportService(
	plans repository.PlanRepository,
	exercises repository.ExerciseRepository,
	mistakes repository.MistakeRepository,
	subjects repository.SubjectRepository,
	kpRepo repository.KnowledgeRepository,
	reports repository.ReportRepository,
) ReportService {
	return &reportService{
		plans:     plans,
		exercises: exercises,
		mistakes:  mistakes,
		subjects:  subjects,
		kpRepo:    kpRepo,
		reports:   reports,
	}
}

func dayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func dayEnd(t time.Time) time.Time {
	return dayStart(t).Add(24*time.Hour - time.Nanosecond)
}

func classifyMastery(rate float64) string {
	switch {
	case rate >= 0.85:
		return "mastered"
	case rate >= 0.6:
		return "learning"
	default:
		return "weak"
	}
}

func (s *reportService) Summary(ctx context.Context, userID uuid.UUID) (*SummaryDTO, error) {
	now := time.Now()
	todayStart := dayStart(now)
	totalStart := now.AddDate(0, 0, -365)

	totalDuration, err := s.plans.SumDurationByUser(ctx, userID, totalStart, dayEnd(now))
	if err != nil {
		return nil, err
	}
	todayDuration, err := s.plans.SumDurationByUser(ctx, userID, todayStart, dayEnd(now))
	if err != nil {
		return nil, err
	}
	totalEx, err := s.exercises.CountByUser(ctx, userID, totalStart, dayEnd(now))
	if err != nil {
		return nil, err
	}
	correct, err := s.exercises.CorrectCountByUser(ctx, userID, totalStart, dayEnd(now))
	if err != nil {
		return nil, err
	}
	rate := 0.0
	if totalEx > 0 {
		rate = float64(correct) / float64(totalEx)
	}
	streak, err := s.plans.CountStreakDays(ctx, userID)
	if err != nil {
		return nil, err
	}
	activePlans, _, err := s.plans.ListByUser(ctx, userID, "active", 1, 1)
	if err != nil {
		return nil, err
	}
	unmastered, err := s.mistakes.CountUnmastered(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &SummaryDTO{
		TodayDurationMin:   todayDuration,
		TotalDurationMin:   totalDuration,
		TotalExercises:     totalEx,
		OverallCorrectRate: rate,
		StreakDays:         streak,
		ActivePlanCount:    int64(len(activePlans)),
		UnmasteredMistakes: unmastered,
	}, nil
}

func (s *reportService) Detail(ctx context.Context, userID uuid.UUID, periodType string, start time.Time) (*DetailDTO, error) {
	var end time.Time
	switch periodType {
	case "daily":
		end = dayEnd(start)
	case "weekly":
		end = dayEnd(start.AddDate(0, 0, 6))
	case "monthly":
		end = dayEnd(start.AddDate(0, 1, -1))
	default:
		end = dayEnd(start)
		periodType = "daily"
	}

	duration, err := s.plans.SumDurationByUser(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}
	totalEx, err := s.exercises.CountByUser(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}
	correct, err := s.exercises.CorrectCountByUser(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}
	rate := 0.0
	if totalEx > 0 {
		rate = float64(correct) / float64(totalEx)
	}
	streak, err := s.plans.CountStreakDays(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 知识点维度（取 subject_id=0 表示全科目聚合过于昂贵，MVP 默认取默认科目）
	// MVP 简化：仅展示空数组，由 mastery 接口单独提供
	mastery := []KPStatDTO{}

	// 按日聚合
	daily := []DailyStat{}
	days := int(end.Sub(start).Hours()/24) + 1
	if days > 60 {
		days = 60
	}
	for i := 0; i < days; i++ {
		ds := start.AddDate(0, 0, i)
		de := dayEnd(ds)
		if de.After(end) {
			de = end
		}
		d, _ := s.plans.SumDurationByUser(ctx, userID, ds, de)
		exN, _ := s.exercises.CountByUser(ctx, userID, ds, de)
		crN, _ := s.exercises.CorrectCountByUser(ctx, userID, ds, de)
		r := 0.0
		if exN > 0 {
			r = float64(crN) / float64(exN)
		}
		daily = append(daily, DailyStat{
			Date:        ds.Format("2006-01-02"),
			DurationMin: d,
			Exercises:   exN,
			CorrectRate: r,
		})
	}

	dto := &DetailDTO{
		Period: PeriodInfo{
			Type:  periodType,
			Start: start.Format("2006-01-02"),
			End:   end.Format("2006-01-02"),
		},
		TotalDurationMin: duration,
		TotalExercises:  totalEx,
		CorrectRate:     rate,
		MasteryByKP:     mastery,
		StreakDays:      streak,
		DailyStats:      daily,
	}

	// MVP 简化：暂不持久化 study_reports 表，由看板查询时实时聚合
	// s.reports 保留字段用于 V1.5 接入异步报告生成
	_ = s.reports

	return dto, nil
}

func (s *reportService) Mastery(ctx context.Context, userID uuid.UUID, subjectID uint64) (*MasteryDTO, error) {
	subjects, err := s.subjects.List(ctx)
	if err != nil {
		return nil, err
	}
	subjInfo := SubjectInfo{ID: subjectID}
	for _, sb := range subjects {
		if sb.ID == subjectID {
			subjInfo.Name = sb.Name
			break
		}
	}

	kps, err := s.kpRepo.ListBySubject(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	rates, err := s.exercises.CorrectRateByKP(ctx, userID, subjectID)
	if err != nil {
		return nil, err
	}
	rateMap := make(map[uint64]float64, len(rates))
	totalMap := make(map[uint64]int64, len(rates))
	for _, r := range rates {
		rateMap[r.KnowledgePointID] = r.Rate
		totalMap[r.KnowledgePointID] = r.Total
	}
	out := []MasteryKPDTO{}
	for _, kp := range kps {
		rate := rateMap[kp.ID]
		total := totalMap[kp.ID]
		status := "learning"
		if total == 0 {
			status = "learning"
			rate = 0
		} else {
			status = classifyMastery(rate)
		}
		out = append(out, MasteryKPDTO{
			ID:     kp.ID,
			Name:   kp.Name,
			Level:  kp.Level,
			Rate:   rate,
			Status: status,
		})
	}
	return &MasteryDTO{
		Subject:         subjInfo,
		KnowledgePoints: out,
	}, nil
}

func (s *reportService) Trend(ctx context.Context, userID uuid.UUID, days int) (*TrendDTO, error) {
	if days <= 0 {
		days = 7
	}
	if days > 60 {
		days = 60
	}
	now := time.Now()
	start := dayStart(now.AddDate(0, 0, -days+1))
	out := []DailyStat{}
	for i := 0; i < days; i++ {
		ds := start.AddDate(0, 0, i)
		de := dayEnd(ds)
		d, _ := s.plans.SumDurationByUser(ctx, userID, ds, de)
		exN, _ := s.exercises.CountByUser(ctx, userID, ds, de)
		crN, _ := s.exercises.CorrectCountByUser(ctx, userID, ds, de)
		r := 0.0
		if exN > 0 {
			r = float64(crN) / float64(exN)
		}
		out = append(out, DailyStat{
			Date:        ds.Format("2006-01-02"),
			DurationMin: d,
			Exercises:   exN,
			CorrectRate: r,
		})
	}
	return &TrendDTO{Days: days, DailyPoints: out}, nil
}