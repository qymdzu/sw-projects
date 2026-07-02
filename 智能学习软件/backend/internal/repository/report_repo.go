package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smart-learning/internal/model"
)

// ReportRepository 是学习报告数据访问接口。
type ReportRepository interface {
	Upsert(ctx context.Context, r *model.StudyReport) error
	GetByPeriod(ctx context.Context, userID uuid.UUID, periodType string, start time.Time) (*model.StudyReport, error)
}

type reportRepo struct {
	db *gorm.DB
}

// NewReportRepository 构造 ReportRepository。
func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepo{db: db}
}

func (r *reportRepo) Upsert(ctx context.Context, rep *model.StudyReport) error {
	// 先查询是否存在
	var existing model.StudyReport
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND period_type = ? AND period_start = ?", rep.UserID, rep.PeriodType, rep.PeriodStart).
		First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.WithContext(ctx).Create(rep).Error
	}
	if err != nil {
		return err
	}
	existing.PeriodEnd = rep.PeriodEnd
	existing.Stats = rep.Stats
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *reportRepo) GetByPeriod(ctx context.Context, userID uuid.UUID, periodType string, start time.Time) (*model.StudyReport, error) {
	var rep model.StudyReport
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND period_type = ? AND period_start = ?", userID, periodType, start).
		First(&rep).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &rep, nil
}