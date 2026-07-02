package repository

import (
	"context"

	"smart-learning/internal/model"
	"gorm.io/gorm"
)

// SubjectRepository 是科目数据访问接口。
type SubjectRepository interface {
	List(ctx context.Context) ([]model.Subject, error)
}

type subjectRepo struct {
	db *gorm.DB
}

// NewSubjectRepository 构造 SubjectRepository。
func NewSubjectRepository(db *gorm.DB) SubjectRepository {
	return &subjectRepo{db: db}
}

func (r *subjectRepo) List(ctx context.Context) ([]model.Subject, error) {
	var items []model.Subject
	err := r.db.WithContext(ctx).Order("id ASC").Find(&items).Error
	return items, err
}