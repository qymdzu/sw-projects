package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smart-learning/internal/model"
)

// MistakeWithQuestion 是错题本记录连同题目信息。
type MistakeWithQuestion struct {
	model.MistakeBook
	Question        model.Question         `gorm:"embedded" json:"question"`
	KnowledgePoint  model.KnowledgePoint   `gorm:"embedded" json:"knowledge_point"`
}

// MistakeGroup 是错题按知识点分组的统计。
type MistakeGroup struct {
	KnowledgePointID uint64 `json:"kp_id"`
	Name             string `json:"name"`
	Total            int64  `json:"total"`
	Mastered         int64  `json:"mastered"`
	Unmastered       int64  `json:"unmastered"`
}

// MistakeRepository 是错题数据访问接口。
type MistakeRepository interface {
	GetByUserQuestion(ctx context.Context, userID uuid.UUID, questionID uint64) (*model.MistakeBook, error)
	Create(ctx context.Context, m *model.MistakeBook) error
	IncrementCount(ctx context.Context, id uint64, lastReviewed time.Time) error
	GetByID(ctx context.Context, id uint64) (*model.MistakeBook, error)
	MarkMastered(ctx context.Context, id uint64, mastered bool) error
	Delete(ctx context.Context, id uint64) error

	ListByUser(ctx context.Context, userID uuid.UUID, kpID *uint64, mastered *bool, page, pageSize int) ([]MistakeWithQuestion, int64, error)
	GroupByKP(ctx context.Context, userID uuid.UUID) ([]MistakeGroup, error)
	ListUnmasteredQuestions(ctx context.Context, userID uuid.UUID, kpIDs []uint64, limit int) ([]model.Question, error)
	CountUnmastered(ctx context.Context, userID uuid.UUID) (int64, error)
}

type mistakeRepo struct {
	db *gorm.DB
}

// NewMistakeRepository 构造 MistakeRepository。
func NewMistakeRepository(db *gorm.DB) MistakeRepository {
	return &mistakeRepo{db: db}
}

func (r *mistakeRepo) GetByUserQuestion(ctx context.Context, userID uuid.UUID, questionID uint64) (*model.MistakeBook, error) {
	var m model.MistakeBook
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND question_id = ?", userID, questionID).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *mistakeRepo) Create(ctx context.Context, m *model.MistakeBook) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *mistakeRepo) IncrementCount(ctx context.Context, id uint64, lastReviewed time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.MistakeBook{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"mistake_count":    gorm.Expr("mistake_count + 1"),
			"last_reviewed_at": lastReviewed,
		}).Error
}

func (r *mistakeRepo) GetByID(ctx context.Context, id uint64) (*model.MistakeBook, error) {
	var m model.MistakeBook
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *mistakeRepo) MarkMastered(ctx context.Context, id uint64, mastered bool) error {
	updates := map[string]interface{}{
		"mastered":   mastered,
	}
	if mastered {
		now := time.Now()
		updates["mastered_at"] = &now
	} else {
		updates["mastered_at"] = nil
	}
	return r.db.WithContext(ctx).
		Model(&model.MistakeBook{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *mistakeRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.MistakeBook{}, "id = ?", id).Error
}

func (r *mistakeRepo) ListByUser(ctx context.Context, userID uuid.UUID, kpID *uint64, mastered *bool, page, pageSize int) ([]MistakeWithQuestion, int64, error) {
	var items []MistakeWithQuestion
	var total int64
	q := r.db.WithContext(ctx).Model(&model.MistakeBook{}).Where("user_id = ?", userID)
	if kpID != nil {
		q = q.Where("knowledge_point_id = ?", *kpID)
	}
	if mastered != nil {
		q = q.Where("mastered = ?", *mastered)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("last_reviewed_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	// 加载关联题目与知识点
	if len(items) > 0 {
		qids := make([]uint64, 0, len(items))
		for _, it := range items {
			qids = append(qids, it.QuestionID)
		}
		var questions []model.Question
		if err := r.db.WithContext(ctx).Preload("KnowledgePoint").Where("id IN ?", qids).Find(&questions).Error; err != nil {
			return nil, 0, err
		}
		qmap := make(map[uint64]model.Question, len(questions))
		for _, q := range questions {
			qmap[q.ID] = q
		}
		for i := range items {
			if q, ok := qmap[items[i].QuestionID]; ok {
				items[i].Question = q
				if q.KnowledgePoint != nil {
					items[i].KnowledgePoint = *q.KnowledgePoint
				}
			}
		}
	}
	return items, total, nil
}

func (r *mistakeRepo) GroupByKP(ctx context.Context, userID uuid.UUID) ([]MistakeGroup, error) {
	type row struct {
		KnowledgePointID uint64
		Name             string
		Total            int64
		Mastered         int64
		Unmastered       int64
	}
	var rows []row
	err := r.db.WithContext(ctx).
		Table("mistake_books mb").
		Select("mb.knowledge_point_id AS knowledge_point_id, kp.name AS name, COUNT(*) AS total, SUM(CASE WHEN mb.mastered THEN 1 ELSE 0 END) AS mastered, SUM(CASE WHEN mb.mastered THEN 0 ELSE 1 END) AS unmastered").
		Joins("JOIN knowledge_points kp ON kp.id = mb.knowledge_point_id").
		Where("mb.user_id = ?", userID).
		Group("mb.knowledge_point_id, kp.name").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]MistakeGroup, 0, len(rows))
	for _, r := range rows {
		out = append(out, MistakeGroup{
			KnowledgePointID: r.KnowledgePointID,
			Name:             r.Name,
			Total:            r.Total,
			Mastered:         r.Mastered,
			Unmastered:       r.Unmastered,
		})
	}
	return out, nil
}

func (r *mistakeRepo) ListUnmasteredQuestions(ctx context.Context, userID uuid.UUID, kpIDs []uint64, limit int) ([]model.Question, error) {
	var items []model.Question
	q := r.db.WithContext(ctx).
		Table("questions q").
		Joins("JOIN mistake_books mb ON mb.question_id = q.id").
		Where("mb.user_id = ? AND mb.mastered = ?", userID, false)
	if len(kpIDs) > 0 {
		q = q.Where("q.knowledge_point_id IN ?", kpIDs)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Preload("KnowledgePoint").Find(&items).Error
	return items, err
}

func (r *mistakeRepo) CountUnmastered(ctx context.Context, userID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&model.MistakeBook{}).
		Where("user_id = ? AND mastered = ?", userID, false).
		Count(&total).Error
	return total, err
}