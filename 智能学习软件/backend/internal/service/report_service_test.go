package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/internal/service"
)

// stubReportRepo 简单 stub。
type stubReportRepo struct{}

func (s *stubReportRepo) Upsert(_ context.Context, _ *model.StudyReport) error { return nil }
func (s *stubReportRepo) GetByPeriod(_ context.Context, _ uuid.UUID, _ string, _ time.Time) (*model.StudyReport, error) {
	return nil, errors.New("not implemented")
}

// stubSubjectRepo stub 科目。
type stubSubjectRepo struct{}

func (s *stubSubjectRepo) List(_ context.Context) ([]model.Subject, error) {
	return []model.Subject{
		{ID: 1, Name: "数学"},
		{ID: 2, Name: "英语"},
	}, nil
}

// stubKnowledgeRepo stub 知识点。
type stubKnowledgeRepo struct{}

func (s *stubKnowledgeRepo) ListTree(_ context.Context, subjectID uint64) ([]repository.KnowledgeNode, error) {
	return []repository.KnowledgeNode{
		{ID: 1, SubjectID: subjectID, Name: "数与代数", Level: 1, Path: "1"},
	}, nil
}

func (s *stubKnowledgeRepo) GetByID(_ context.Context, id uint64) (*model.KnowledgePoint, error) {
	return &model.KnowledgePoint{ID: id, Name: "一级知识点", Level: 1}, nil
}

func (s *stubKnowledgeRepo) ListBySubject(_ context.Context, _ uint64) ([]model.KnowledgePoint, error) {
	return []model.KnowledgePoint{
		{ID: 1, Name: "实数", Level: 1},
		{ID: 2, Name: "一元二次方程", Level: 2},
	}, nil
}

func TestReportService_Summary(t *testing.T) {
	planRepo := newMockPlanRepo()
	exRepo := newMockExerciseRepo()
	mistakeRepo := newMockMistakeRepoForExercise()
	reportSvc := service.NewReportService(planRepo, exRepo, mistakeRepo, &stubSubjectRepo{}, &stubKnowledgeRepo{}, &stubReportRepo{})

	sum, err := reportSvc.Summary(context.Background(), uuid.New())
	require.NoError(t, err)
	assert.NotNil(t, sum)
	assert.Equal(t, 0, sum.StreakDays)
	assert.Equal(t, int64(0), sum.UnmasteredMistakes)
}

func TestReportService_Detail(t *testing.T) {
	planRepo := newMockPlanRepo()
	exRepo := newMockExerciseRepo()
	mistakeRepo := newMockMistakeRepoForExercise()
	reportSvc := service.NewReportService(planRepo, exRepo, mistakeRepo, &stubSubjectRepo{}, &stubKnowledgeRepo{}, &stubReportRepo{})

	start, _ := time.Parse("2006-01-02", "2026-07-01")
	dto, err := reportSvc.Detail(context.Background(), uuid.New(), "weekly", start)
	require.NoError(t, err)
	assert.Equal(t, "weekly", dto.Period.Type)
	assert.Equal(t, 7, len(dto.DailyStats))
}

func TestReportService_Mastery(t *testing.T) {
	planRepo := newMockPlanRepo()
	exRepo := newMockExerciseRepo()
	mistakeRepo := newMockMistakeRepoForExercise()
	reportSvc := service.NewReportService(planRepo, exRepo, mistakeRepo, &stubSubjectRepo{}, &stubKnowledgeRepo{}, &stubReportRepo{})

	dto, err := reportSvc.Mastery(context.Background(), uuid.New(), 1)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), dto.Subject.ID)
	assert.Equal(t, "数学", dto.Subject.Name)
	assert.Greater(t, len(dto.KnowledgePoints), 0)
}

func TestReportService_Trend(t *testing.T) {
	planRepo := newMockPlanRepo()
	exRepo := newMockExerciseRepo()
	mistakeRepo := newMockMistakeRepoForExercise()
	reportSvc := service.NewReportService(planRepo, exRepo, mistakeRepo, &stubSubjectRepo{}, &stubKnowledgeRepo{}, &stubReportRepo{})

	dto, err := reportSvc.Trend(context.Background(), uuid.New(), 7)
	require.NoError(t, err)
	assert.Equal(t, 7, dto.Days)
	assert.Len(t, dto.DailyPoints, 7)
}

func TestReportService_Trend_ClampDays(t *testing.T) {
	planRepo := newMockPlanRepo()
	exRepo := newMockExerciseRepo()
	mistakeRepo := newMockMistakeRepoForExercise()
	reportSvc := service.NewReportService(planRepo, exRepo, mistakeRepo, &stubSubjectRepo{}, &stubKnowledgeRepo{}, &stubReportRepo{})

	dto, err := reportSvc.Trend(context.Background(), uuid.New(), 100)
	require.NoError(t, err)
	assert.Equal(t, 60, dto.Days)
}

func TestReportService_Trend_DefaultDays(t *testing.T) {
	planRepo := newMockPlanRepo()
	exRepo := newMockExerciseRepo()
	mistakeRepo := newMockMistakeRepoForExercise()
	reportSvc := service.NewReportService(planRepo, exRepo, mistakeRepo, &stubSubjectRepo{}, &stubKnowledgeRepo{}, &stubReportRepo{})

	dto, err := reportSvc.Trend(context.Background(), uuid.New(), 0)
	require.NoError(t, err)
	assert.Equal(t, 7, dto.Days)
}

func TestReportService_Detail_Daily(t *testing.T) {
	planRepo := newMockPlanRepo()
	exRepo := newMockExerciseRepo()
	mistakeRepo := newMockMistakeRepoForExercise()
	reportSvc := service.NewReportService(planRepo, exRepo, mistakeRepo, &stubSubjectRepo{}, &stubKnowledgeRepo{}, &stubReportRepo{})

	start, _ := time.Parse("2006-01-02", "2026-07-01")
	dto, err := reportSvc.Detail(context.Background(), uuid.New(), "daily", start)
	require.NoError(t, err)
	assert.Equal(t, "daily", dto.Period.Type)
	assert.Len(t, dto.DailyStats, 1)
}

func TestReportService_Detail_Monthly(t *testing.T) {
	planRepo := newMockPlanRepo()
	exRepo := newMockExerciseRepo()
	mistakeRepo := newMockMistakeRepoForExercise()
	reportSvc := service.NewReportService(planRepo, exRepo, mistakeRepo, &stubSubjectRepo{}, &stubKnowledgeRepo{}, &stubReportRepo{})

	start, _ := time.Parse("2006-01-02", "2026-07-01")
	dto, err := reportSvc.Detail(context.Background(), uuid.New(), "monthly", start)
	require.NoError(t, err)
	assert.Equal(t, "monthly", dto.Period.Type)
}