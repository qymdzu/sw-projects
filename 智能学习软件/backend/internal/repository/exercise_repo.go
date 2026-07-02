package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smart-learning/internal/model"
)

// QuestionFilter 是题目查询过滤条件。
type QuestionFilter struct {
	SubjectID        *uint64
	KnowledgePointID *uint64
	Type             string
	Difficulty       *int
	Page             int
	PageSize         int
}

// KPRate 是知识点维度的正确率聚合。
type KPRate struct {
	KnowledgePointID uint64  `json:"kp_id"`
	Name             string  `json:"name"`
	Total            int64   `json:"total"`
	Correct          int64   `json:"correct"`
	Rate             float64 `json:"rate"`
}

// ExerciseRepository 是练习数据访问接口。
type ExerciseRepository interface {
	// Question
	GetQuestionByID(ctx context.Context, id uint64) (*model.Question, error)
	ListQuestions(ctx context.Context, filter QuestionFilter) ([]model.Question, int64, error)
	ListQuestionsByIDs(ctx context.Context, ids []uint64) ([]model.Question, error)
	RandomQuestions(ctx context.Context, subjectID, kpID *uint64, count int) ([]model.Question, error)

	// Record
	CreateRecord(ctx context.Context, rec *model.ExerciseRecord) error
	CountByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error)
	CorrectCountByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error)
	HistoryByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]model.ExerciseRecord, int64, error)
	CorrectRateByKP(ctx context.Context, userID uuid.UUID, subjectID uint64) ([]KPRate, error)
}

type exerciseRepo struct {
	db *gorm.DB
}

// NewExerciseRepository 构造 ExerciseRepository。
func NewExerciseRepository(db *gorm.DB) ExerciseRepository {
	return &exerciseRepo{db: db}
}

func (r *exerciseRepo) GetQuestionByID(ctx context.Context, id uint64) (*model.Question, error) {
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

func (r *exerciseRepo) ListQuestions(ctx context.Context, f QuestionFilter) ([]model.Question, int64, error) {
	var items []model.Question
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Question{})
	if f.SubjectID != nil {
		q = q.Where("subject_id = ?", *f.SubjectID)
	}
	if f.KnowledgePointID != nil {
		q = q.Where("knowledge_point_id = ?", *f.KnowledgePointID)
	}
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.Difficulty != nil {
		q = q.Where("difficulty = ?", *f.Difficulty)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Preload("KnowledgePoint").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *exerciseRepo) ListQuestionsByIDs(ctx context.Context, ids []uint64) ([]model.Question, error) {
	var items []model.Question
	if len(ids) == 0 {
		return items, nil
	}
	err := r.db.WithContext(ctx).Preload("KnowledgePoint").Where("id IN ?", ids).Find(&items).Error
	return items, err
}

func (r *exerciseRepo) RandomQuestions(ctx context.Context, subjectID, kpID *uint64, count int) ([]model.Question, error) {
	var items []model.Question
	q := r.db.WithContext(ctx).Model(&model.Question{})
	if subjectID != nil {
		q = q.Where("subject_id = ?", *subjectID)
	}
	if kpID != nil {
		q = q.Where("knowledge_point_id = ?", *kpID)
	}
	// 用 ORDER BY RANDOM() 简化实现，MVP 数据量可控
	err := q.Order("random()").Limit(count).Preload("KnowledgePoint").Find(&items).Error
	return items, err
}

func (r *exerciseRepo) CreateRecord(ctx context.Context, rec *model.ExerciseRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *exerciseRepo) CountByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&model.ExerciseRecord{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, from, to).
		Count(&total).Error
	return total, err
}

func (r *exerciseRepo) CorrectCountByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&model.ExerciseRecord{}).
		Where("user_id = ? AND is_correct = ? AND created_at BETWEEN ? AND ?", userID, true, from, to).
		Count(&total).Error
	return total, err
}

func (r *exerciseRepo) HistoryByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]model.ExerciseRecord, int64, error) {
	var items []model.ExerciseRecord
	var total int64
	q := r.db.WithContext(ctx).Model(&model.ExerciseRecord{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *exerciseRepo) CorrectRateByKP(ctx context.Context, userID uuid.UUID, subjectID uint64) ([]KPRate, error) {
	type row struct {
		KnowledgePointID uint64
		Name             string
		Total            int64
		Correct          int64
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Table("exercise_records er").
		Select("q.knowledge_point_id AS knowledge_point_id, kp.name AS name, COUNT(*) AS total, SUM(CASE WHEN er.is_correct THEN 1 ELSE 0 END) AS correct").
		Joins("JOIN questions q ON q.id = er.question_id").
		Joins("JOIN knowledge_points kp ON kp.id = q.knowledge_point_id").
		Where("er.user_id = ? AND q.subject_id = ?", userID, subjectID).
		Group("q.knowledge_point_id, kp.name").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]KPRate, 0, len(rows))
	for _, r := range rows {
		rate := 0.0
		if r.Total > 0 {
			rate = float64(r.Correct) / float64(r.Total)
		}
		out = append(out, KPRate{
			KnowledgePointID: r.KnowledgePointID,
			Name:             r.Name,
			Total:            r.Total,
			Correct:          r.Correct,
			Rate:             rate,
		})
	}
	return out, nil
}