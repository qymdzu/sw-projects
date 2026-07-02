package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"smart-learning/internal/model"
)

// KnowledgeNode 是知识点树节点（API响应）。
type KnowledgeNode struct {
	ID       uint64          `json:"id"`
	SubjectID uint64         `json:"subject_id"`
	ParentID *uint64         `json:"parent_id,omitempty"`
	Name     string          `json:"name"`
	Level    int             `json:"level"`
	Path     string          `json:"path"`
	Children []KnowledgeNode `json:"children"`
}

// KnowledgeRepository 是知识点数据访问接口。
type KnowledgeRepository interface {
	ListTree(ctx context.Context, subjectID uint64) ([]KnowledgeNode, error)
	GetByID(ctx context.Context, id uint64) (*model.KnowledgePoint, error)
	ListBySubject(ctx context.Context, subjectID uint64) ([]model.KnowledgePoint, error)
}

type knowledgeRepo struct {
	db *gorm.DB
}

// NewKnowledgeRepository 构造 KnowledgeRepository。
func NewKnowledgeRepository(db *gorm.DB) KnowledgeRepository {
	return &knowledgeRepo{db: db}
}

func (r *knowledgeRepo) ListTree(ctx context.Context, subjectID uint64) ([]KnowledgeNode, error) {
	rows, err := r.ListBySubject(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	return buildTree(rows), nil
}

func (r *knowledgeRepo) GetByID(ctx context.Context, id uint64) (*model.KnowledgePoint, error) {
	var kp model.KnowledgePoint
	err := r.db.WithContext(ctx).First(&kp, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &kp, nil
}

func (r *knowledgeRepo) ListBySubject(ctx context.Context, subjectID uint64) ([]model.KnowledgePoint, error) {
	var items []model.KnowledgePoint
	err := r.db.WithContext(ctx).Where("subject_id = ?", subjectID).Order("level ASC, id ASC").Find(&items).Error
	return items, err
}

// buildTree 根据 parent_id 自引用构建树。
func buildTree(rows []model.KnowledgePoint) []KnowledgeNode {
	nodes := make(map[uint64]*KnowledgeNode, len(rows))
	roots := make([]*KnowledgeNode, 0)
	for i := range rows {
		r := &rows[i]
		n := &KnowledgeNode{
			ID: r.ID, SubjectID: r.SubjectID, ParentID: r.ParentID,
			Name: r.Name, Level: r.Level, Path: r.Path,
			Children: []KnowledgeNode{},
		}
		nodes[r.ID] = n
		if r.ParentID == nil {
			roots = append(roots, n)
		}
	}
	for _, r := range rows {
		if r.ParentID != nil {
			parent, ok := nodes[*r.ParentID]
			if ok {
				parent.Children = append(parent.Children, *nodes[r.ID])
			}
		}
	}
	out := make([]KnowledgeNode, 0, len(roots))
	for _, r := range roots {
		out = append(out, *r)
	}
	return out
}