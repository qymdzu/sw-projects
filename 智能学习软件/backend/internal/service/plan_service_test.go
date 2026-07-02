package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/internal/service"
)

// mockPlanRepo 是 PlanRepository 的内存 mock。
type mockPlanRepo struct {
	plans       map[uint64]*model.StudyPlan
	records     []model.StudyRecord
	nextPlanID  uint64
	nextRecID   uint64
}

func newMockPlanRepo() *mockPlanRepo {
	return &mockPlanRepo{
		plans:      make(map[uint64]*model.StudyPlan),
		nextPlanID: 1,
		nextRecID:  1,
	}
}

func (m *mockPlanRepo) Create(_ context.Context, p *model.StudyPlan) error {
	if p.ID == 0 {
		p.ID = m.nextPlanID
		m.nextPlanID++
	}
	m.plans[p.ID] = p
	return nil
}

func (m *mockPlanRepo) GetByID(_ context.Context, id uint64) (*model.StudyPlan, error) {
	p, ok := m.plans[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return p, nil
}

func (m *mockPlanRepo) ListByUser(_ context.Context, userID uuid.UUID, status string, page, pageSize int) ([]model.StudyPlan, int64, error) {
	var items []model.StudyPlan
	var total int64
	for _, p := range m.plans {
		if p.UserID == userID {
			if status == "" || p.Status == status {
				items = append(items, *p)
				total++
			}
		}
	}
	return items, total, nil
}

func (m *mockPlanRepo) Update(_ context.Context, p *model.StudyPlan) error {
	m.plans[p.ID] = p
	return nil
}

func (m *mockPlanRepo) Delete(_ context.Context, id uint64) error {
	delete(m.plans, id)
	return nil
}

func (m *mockPlanRepo) CreateRecord(_ context.Context, r *model.StudyRecord) error {
	if r.ID == 0 {
		r.ID = m.nextRecID
		m.nextRecID++
	}
	m.records = append(m.records, *r)
	return nil
}

func (m *mockPlanRepo) ListRecordsByUser(_ context.Context, userID uuid.UUID, from, to time.Time) ([]model.StudyRecord, error) {
	var out []model.StudyRecord
	for _, r := range m.records {
		if r.UserID == userID && !r.Date.Before(from) && !r.Date.After(to) {
			out = append(out, r)
		}
	}
	return out, nil
}

func (m *mockPlanRepo) SumDurationByUser(_ context.Context, userID uuid.UUID, from, to time.Time) (int, error) {
	var sum int
	for _, r := range m.records {
		if r.UserID == userID && !r.Date.Before(from) && !r.Date.After(to) {
			sum += r.DurationMin
		}
	}
	return sum, nil
}

func (m *mockPlanRepo) CountStreakDays(_ context.Context, userID uuid.UUID) (int, error) {
	dates := map[time.Time]bool{}
	for _, r := range m.records {
		if r.UserID == userID {
			d := time.Date(r.Date.Year(), r.Date.Month(), r.Date.Day(), 0, 0, 0, 0, time.UTC)
			dates[d] = true
		}
	}
	today := time.Now().Truncate(24 * time.Hour)
	if dates[today] || dates[today.AddDate(0, 0, -1)] {
		streak := 1
		prev := today
		if !dates[today] {
			prev = today.AddDate(0, 0, -1)
		}
		for i := 1; i < 60; i++ {
			d := prev.AddDate(0, 0, -i)
			if dates[d] {
				streak++
			} else {
				break
			}
		}
		return streak, nil
	}
	return 0, nil
}

func TestPlanService_Create(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	uid := uuid.New()
	start := time.Now()
	end := start.AddDate(0, 0, 7)

	p, err := svc.Create(context.Background(), uid, service.CreatePlanRequest{
		Goal:      "复习一元二次方程",
		StartDate: start,
		EndDate:   end,
		Items: []service.PlanItem{
			{Date: "2026-07-01", KnowledgePointIDs: []uint64{1, 3}, DurationMin: 60},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, uid, p.UserID)
	assert.Equal(t, "active", p.Status)
	assert.False(t, p.AIGenerated)
}

func TestPlanService_Create_InvalidDates(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	start := time.Now()
	end := start.AddDate(0, 0, -1)
	_, err := svc.Create(context.Background(), uuid.New(), service.CreatePlanRequest{
		Goal: "x", StartDate: start, EndDate: end,
	})
	assert.Error(t, err)
}

func TestPlanService_Create_EmptyGoal(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	_, err := svc.Create(context.Background(), uuid.New(), service.CreatePlanRequest{
		StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7),
	})
	assert.Error(t, err)
}

func TestPlanService_GetByID_NotFound(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	_, err := svc.GetByID(context.Background(), uuid.New(), 999)
	assert.Error(t, err)
}

func TestPlanService_GetByID_Forbidden(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	owner := uuid.New()
	p, err := svc.Create(context.Background(), owner, service.CreatePlanRequest{
		Goal: "x", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7),
	})
	require.NoError(t, err)
	_, err = svc.GetByID(context.Background(), uuid.New(), p.ID)
	assert.ErrorIs(t, err, service.ErrForbidden)
}

func TestPlanService_Update(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	uid := uuid.New()
	p, err := svc.Create(context.Background(), uid, service.CreatePlanRequest{
		Goal: "old", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7),
	})
	require.NoError(t, err)

	updated, err := svc.Update(context.Background(), uid, p.ID, service.CreatePlanRequest{
		Goal: "new", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 14),
		Items: []service.PlanItem{
			{Date: "2026-07-01", DurationMin: 60},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "new", updated.Goal)
}

func TestPlanService_Delete(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	uid := uuid.New()
	p, err := svc.Create(context.Background(), uid, service.CreatePlanRequest{
		Goal: "x", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7),
	})
	require.NoError(t, err)
	require.NoError(t, svc.Delete(context.Background(), uid, p.ID))
}

func TestPlanService_AIGenerate(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	uid := uuid.New()
	start := time.Now()
	end := start.AddDate(0, 0, 6)

	p, aiUsed, err := svc.AIGenerate(context.Background(), uid, service.AIGenerateRequest{
		Goal:             "复习计划",
		StartDate:        start,
		EndDate:          end,
		DailyDurationMin: 60,
	})
	require.NoError(t, err)
	assert.True(t, p.AIGenerated)
	assert.False(t, aiUsed) // MVP 走规则降级

	// 解析 items
	var items []service.PlanItem
	require.NoError(t, jsonUnmarshalBytes(p.Items, &items))
	assert.Len(t, items, 7)
}

func TestPlanService_Checkin(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	uid := uuid.New()
	p, err := svc.Create(context.Background(), uid, service.CreatePlanRequest{
		Goal: "x", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7),
	})
	require.NoError(t, err)

	rec, err := svc.Checkin(context.Background(), uid, p.ID, service.CheckinRequest{
		Date:        time.Now(),
		DurationMin: 60,
		Status:      "done",
	})
	require.NoError(t, err)
	assert.Equal(t, "done", rec.Status)
}

func TestPlanService_Checkin_InvalidStatus(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	uid := uuid.New()
	p, err := svc.Create(context.Background(), uid, service.CreatePlanRequest{
		Goal: "x", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7),
	})
	require.NoError(t, err)

	_, err = svc.Checkin(context.Background(), uid, p.ID, service.CheckinRequest{
		Date: time.Now(), DurationMin: 60, Status: "bogus",
	})
	assert.Error(t, err)
}

func TestPlanService_Checkin_InvalidDuration(t *testing.T) {
	repo := newMockPlanRepo()
	svc := service.NewPlanService(repo)
	uid := uuid.New()
	p, err := svc.Create(context.Background(), uid, service.CreatePlanRequest{
		Goal: "x", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7),
	})
	require.NoError(t, err)

	_, err = svc.Checkin(context.Background(), uid, p.ID, service.CheckinRequest{
		Date: time.Now(), DurationMin: 0,
	})
	assert.Error(t, err)
}

func jsonUnmarshalBytes(b datatypes.JSON, v interface{}) error {
	return jsonStdUnmarshal(b, v)
}