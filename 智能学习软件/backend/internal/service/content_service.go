package service

import (
	"context"
	"errors"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
)

// SubjectService 科目服务接口。
type SubjectService interface {
	List(ctx context.Context) ([]model.Subject, error)
}

type subjectService struct {
	subjects repository.SubjectRepository
}

// NewSubjectService 构造 SubjectService。
func NewSubjectService(s repository.SubjectRepository) SubjectService {
	return &subjectService{subjects: s}
}

func (s *subjectService) List(ctx context.Context) ([]model.Subject, error) {
	return s.subjects.List(ctx)
}

// KnowledgeService 知识点服务接口。
type KnowledgeService interface {
	GetTree(ctx context.Context, subjectID uint64) ([]repository.KnowledgeNode, error)
}

type knowledgeService struct {
	kp repository.KnowledgeRepository
}

// NewKnowledgeService 构造 KnowledgeService。
func NewKnowledgeService(kp repository.KnowledgeRepository) KnowledgeService {
	return &knowledgeService{kp: kp}
}

func (s *knowledgeService) GetTree(ctx context.Context, subjectID uint64) ([]repository.KnowledgeNode, error) {
	if subjectID == 0 {
		return nil, errors.New("subject_id 不能为空")
	}
	return s.kp.ListTree(ctx, subjectID)
}