package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smart-learning/internal/model"
)

// PlanRepository 是学习计划数据访问接口。
type PlanRepository interface {
	Create(ctx context.Context, plan *model.StudyPlan) error
	GetByID(ctx context.Context, id uint64) (*model.StudyPlan, error)
	ListByUser(ctx context.Context, userID uuid.UUID, status string, page, pageSize int) ([]model.StudyPlan, int64, error)
	Update(ctx context.Context, plan *model.StudyPlan) error
	Delete(ctx context.Context, id uint64) error

	CreateRecord(ctx context.Context, rec *model.StudyRecord) error
	ListRecordsByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.StudyRecord, error)
	SumDurationByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error)
	CountStreakDays(ctx context.Context, userID uuid.UUID) (int, error)
}

type planRepo struct {
	db *gorm.DB
}

// NewPlanRepository 构造 PlanRepository。
func NewPlanRepository(db *gorm.DB) PlanRepository {
	return &planRepo{db: db}
}

func (r *planRepo) Create(ctx context.Context, plan *model.StudyPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *planRepo) GetByID(ctx context.Context, id uint64) (*model.StudyPlan, error) {
	var p model.StudyPlan
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *planRepo) ListByUser(ctx context.Context, userID uuid.UUID, status string, page, pageSize int) ([]model.StudyPlan, int64, error) {
	var items []model.StudyPlan
	var total int64
	q := r.db.WithContext(ctx).Model(&model.StudyPlan{}).Where("user_id = ?", userID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *planRepo) Update(ctx context.Context, plan *model.StudyPlan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

func (r *planRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.StudyPlan{}, "id = ?", id).Error
}

func (r *planRepo) CreateRecord(ctx context.Context, rec *model.StudyRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *planRepo) ListRecordsByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.StudyRecord, error) {
	var items []model.StudyRecord
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, from, to).
		Order("date DESC").
		Find(&items).Error
	return items, err
}

func (r *planRepo) SumDurationByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error) {
	var total int
	err := r.db.WithContext(ctx).
		Model(&model.StudyRecord{}).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, from, to).
		Select("COALESCE(SUM(duration_min), 0)").
		Scan(&total).Error
	return total, err
}

func (r *planRepo) CountStreakDays(ctx context.Context, userID uuid.UUID) (int, error) {
	// 简化策略：统计最近连续有记录的日期天数
	var dates []time.Time
	err := r.db.WithContext(ctx).
		Model(&model.StudyRecord{}).
		Where("user_id = ?", userID).
		Distinct("date").
		Order("date DESC").
		Limit(60).
		Pluck("date", &dates).Error
	if err != nil {
		return 0, err
	}
	if len(dates) == 0 {
		return 0, nil
	}
	streak := 1
	today := time.Now().Truncate(24 * time.Hour)
	prev := today
	for i, d := range dates {
		d = d.Truncate(24 * time.Hour)
		if i == 0 {
			// 第一天必须是今天或昨天
			if d.Equal(today) || d.Equal(today.AddDate(0, 0, -1)) {
				prev = d
				continue
			}
			return 0, nil
		}
		if d.Equal(prev.AddDate(0, 0, -1)) {
			streak++
			prev = d
		} else {
			break
		}
	}
	return streak, nil
}