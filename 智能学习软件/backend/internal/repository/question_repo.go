package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"smart-learning/internal/model"
)

// QuestionRepository 是题目管理数据访问接口（MVP 简化为查询，写操作通过 service 组合 GORM）。
type QuestionRepository interface {
	GetByID(ctx context.Context, id uint64) (*model.Question, error)
	Create(ctx context.Context, q *model.Question) error
	Update(ctx context.Context, q *model.Question) error
	Delete(ctx context.Context, id uint64) error
	BatchCreate(ctx context.Context, items []model.Question) error
	ListBySubject(ctx context.Context, subjectID uint64) ([]model.Question, error)
}

type questionRepo struct {
	db *gorm.DB
}

// NewQuestionRepository 构造 QuestionRepository。
func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepo{db: db}
}

func (r *questionRepo) GetByID(ctx context.Context, id uint64) (*model.Question, error) {
	var q model.Question
	err := r.db.WithContext(ctx).Preload("KnowledgePoint").First(&q, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *questionRepo) Create(ctx context.Context, q *model.Question) error {
	return r.db.WithContext(ctx).Create(q).Error
}

func (r *questionRepo) Update(ctx context.Context, q *model.Question) error {
	return r.db.WithContext(ctx).Save(q).Error
}

func (r *questionRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Question{}, "id = ?", id).Error
}

func (r *questionRepo) BatchCreate(ctx context.Context, items []model.Question) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(items, 100).Error
}

func (r *questionRepo) ListBySubject(ctx context.Context, subjectID uint64) ([]model.Question, error) {
	var items []model.Question
	err := r.db.WithContext(ctx).Preload("KnowledgePoint").Where("subject_id = ?", subjectID).Find(&items).Error
	return items, err
}