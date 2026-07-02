package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
)

// ReviewRequest 错题重练请求。
type ReviewRequest struct {
	KnowledgePointIDs []uint64
	Count             int
}

// ReviewResponse 错题重练响应。
type ReviewResponse struct {
	Questions         []model.Question `json:"questions"`
	TotalUnmastered   int64            `json:"total_unmastered"`
	Strategy          string           `json:"strategy"`
}

// MistakeService 错题服务接口。
type MistakeService interface {
	RecordMistake(ctx context.Context, userID uuid.UUID, questionID uint64, wrongAnswer string) error
	List(ctx context.Context, userID uuid.UUID, kpID *uint64, mastered *bool, page, pageSize int) ([]repository.MistakeWithQuestion, int64, error)
	GroupByKP(ctx context.Context, userID uuid.UUID) ([]repository.MistakeGroup, error)
	MarkMastered(ctx context.Context, userID uuid.UUID, mistakeID uint64, mastered bool) error
	Review(ctx context.Context, userID uuid.UUID, req ReviewRequest) (*ReviewResponse, error)
	Delete(ctx context.Context, userID uuid.UUID, mistakeID uint64) error
}

type mistakeService struct {
	mistakes  repository.MistakeRepository
	exercises repository.ExerciseRepository
}

// NewMistakeService 构造 MistakeService。
func NewMistakeService(m repository.MistakeRepository, ex repository.ExerciseRepository) MistakeService {
	return &mistakeService{mistakes: m, exercises: ex}
}

func (s *mistakeService) RecordMistake(ctx context.Context, userID uuid.UUID, questionID uint64, wrongAnswer string) error {
	q, err := s.exercises.GetQuestionByID(ctx, questionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrResourceMissing
		}
		return err
	}
	existing, err := s.mistakes.GetByUserQuestion(ctx, userID, q.ID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}
	if existing != nil {
		return s.mistakes.IncrementCount(ctx, existing.ID, time.Now())
	}
	m := &model.MistakeBook{
		UserID:           userID,
		QuestionID:       q.ID,
		KnowledgePointID: q.KnowledgePointID,
		WrongAnswer:      wrongAnswer,
		MistakeCount:     1,
		Mastered:         false,
		LastReviewedAt:   time.Now(),
	}
	return s.mistakes.Create(ctx, m)
}

func (s *mistakeService) List(ctx context.Context, userID uuid.UUID, kpID *uint64, mastered *bool, page, pageSize int) ([]repository.MistakeWithQuestion, int64, error) {
	return s.mistakes.ListByUser(ctx, userID, kpID, mastered, page, pageSize)
}

func (s *mistakeService) GroupByKP(ctx context.Context, userID uuid.UUID) ([]repository.MistakeGroup, error) {
	return s.mistakes.GroupByKP(ctx, userID)
}

func (s *mistakeService) MarkMastered(ctx context.Context, userID uuid.UUID, mistakeID uint64, mastered bool) error {
	m, err := s.mistakes.GetByID(ctx, mistakeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrResourceMissing
		}
		return err
	}
	if m.UserID != userID {
		return ErrForbidden
	}
	return s.mistakes.MarkMastered(ctx, mistakeID, mastered)
}

func (s *mistakeService) Review(ctx context.Context, userID uuid.UUID, req ReviewRequest) (*ReviewResponse, error) {
	if req.Count <= 0 {
		req.Count = 10
	}
	questions, err := s.mistakes.ListUnmasteredQuestions(ctx, userID, req.KnowledgePointIDs, req.Count)
	if err != nil {
		return nil, err
	}
	total, err := s.mistakes.CountUnmastered(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &ReviewResponse{
		Questions:       questions,
		TotalUnmastered: total,
		Strategy:        "unmastered_only",
	}, nil
}

func (s *mistakeService) Delete(ctx context.Context, userID uuid.UUID, mistakeID uint64) error {
	m, err := s.mistakes.GetByID(ctx, mistakeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrResourceMissing
		}
		return err
	}
	if m.UserID != userID {
		return ErrForbidden
	}
	return s.mistakes.Delete(ctx, mistakeID)
}