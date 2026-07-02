package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
)

// SubmitRequest 提交答案请求。
type SubmitRequest struct {
	QuestionID  uint64
	Answer      string
	DurationSec int
}

// SubmitResponse 提交答案响应。
type SubmitResponse struct {
	IsCorrect        bool
	CorrectAnswer    string
	Analysis         string
	MistakeRecorded  bool
	RecordID         uint64
}

// ExerciseService 练习服务接口。
type ExerciseService interface {
	List(ctx context.Context, filter repository.QuestionFilter) ([]model.Question, int64, error)
	Random(ctx context.Context, subjectID, kpID *uint64, count int) ([]model.Question, error)
	Submit(ctx context.Context, userID uuid.UUID, req SubmitRequest) (*SubmitResponse, error)
	Recommend(ctx context.Context, userID uuid.UUID, count int) ([]model.Question, error)
	ByKnowledgePoint(ctx context.Context, kpID uint64, page, pageSize int) ([]model.Question, int64, error)
	History(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]model.ExerciseRecord, int64, error)
}

type exerciseService struct {
	exercises repository.ExerciseRepository
	mistakes  repository.MistakeRepository
}

// NewExerciseService 构造 ExerciseService。
func NewExerciseService(ex repository.ExerciseRepository, m repository.MistakeRepository) ExerciseService {
	return &exerciseService{exercises: ex, mistakes: m}
}

// judgeAnswer 根据题目类型批改答案。
func judgeAnswer(q *model.Question, userAnswer string) bool {
	userAns := strings.TrimSpace(userAnswer)
	correct := strings.TrimSpace(q.Answer)
	switch q.Type {
	case "choice", "judge":
		return strings.EqualFold(userAns, correct)
	case "fill":
		// 填空题支持多个答案以 | 分隔
		for _, alt := range strings.Split(correct, "|") {
			if strings.EqualFold(strings.TrimSpace(userAns), strings.TrimSpace(alt)) {
				return true
			}
		}
		return false
	case "subjective":
		// MVP 主观题默认不自动判分，由调用方根据 answer 关键词自行决定
		// 这里采用宽松匹配：答案包含 correct 中的关键词即视为正确
		return strings.Contains(strings.ToLower(userAns), strings.ToLower(correct))
	default:
		return false
	}
}

// recordMistakeHelper 错误时收录错题（幂等）。
func (s *exerciseService) recordMistakeHelper(ctx context.Context, userID uuid.UUID, q *model.Question, wrongAnswer string) (bool, error) {
	existing, err := s.mistakes.GetByUserQuestion(ctx, userID, q.ID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return false, err
	}
	if existing != nil {
		if err := s.mistakes.IncrementCount(ctx, existing.ID, time.Now()); err != nil {
			return false, err
		}
		return false, nil
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
	if err := s.mistakes.Create(ctx, m); err != nil {
		return false, err
	}
	return true, nil
}

func (s *exerciseService) List(ctx context.Context, filter repository.QuestionFilter) ([]model.Question, int64, error) {
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	return s.exercises.ListQuestions(ctx, filter)
}

func (s *exerciseService) Random(ctx context.Context, subjectID, kpID *uint64, count int) ([]model.Question, error) {
	if count <= 0 {
		count = 10
	}
	return s.exercises.RandomQuestions(ctx, subjectID, kpID, count)
}

func (s *exerciseService) Submit(ctx context.Context, userID uuid.UUID, req SubmitRequest) (*SubmitResponse, error) {
	if req.QuestionID == 0 {
		return nil, errors.New("题目 ID 不能为空")
	}
	if strings.TrimSpace(req.Answer) == "" {
		return nil, errors.New("答案不能为空")
	}
	q, err := s.exercises.GetQuestionByID(ctx, req.QuestionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrResourceMissing
		}
		return nil, err
	}

	isCorrect := judgeAnswer(q, req.Answer)
	durSec := req.DurationSec
	var durPtr *int
	if durSec > 0 {
		durPtr = &durSec
	}
	rec := &model.ExerciseRecord{
		UserID:      userID,
		QuestionID:  q.ID,
		Answer:      req.Answer,
		IsCorrect:   isCorrect,
		DurationSec: durPtr,
	}
	if err := s.exercises.CreateRecord(ctx, rec); err != nil {
		return nil, err
	}

	resp := &SubmitResponse{
		IsCorrect:     isCorrect,
		CorrectAnswer: q.Answer,
		RecordID:      rec.ID,
	}
	if q.Analysis != nil {
		resp.Analysis = *q.Analysis
	}

	if !isCorrect {
		recorded, err := s.recordMistakeHelper(ctx, userID, q, req.Answer)
		if err != nil {
			return nil, err
		}
		resp.MistakeRecorded = recorded
	}
	return resp, nil
}

// Recommend 智能推荐：优先返回错题本中未掌握知识点对应的题目，再补充随机题。
func (s *exerciseService) Recommend(ctx context.Context, userID uuid.UUID, count int) ([]model.Question, error) {
	if count <= 0 {
		count = 10
	}
	// 1. 取错题本中未掌握题目
	mistakeQuestions, err := s.mistakes.ListUnmasteredQuestions(ctx, userID, nil, count)
	if err != nil {
		return nil, err
	}
	if len(mistakeQuestions) >= count {
		return mistakeQuestions[:count], nil
	}
	// 2. 不足时用随机题补充
	need := count - len(mistakeQuestions)
	randQ, err := s.exercises.RandomQuestions(ctx, nil, nil, need)
	if err != nil {
		return nil, err
	}
	out := append(mistakeQuestions, randQ...)
	return out, nil
}

func (s *exerciseService) ByKnowledgePoint(ctx context.Context, kpID uint64, page, pageSize int) ([]model.Question, int64, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	filter := repository.QuestionFilter{
		KnowledgePointID: &kpID,
		Page:             page,
		PageSize:         pageSize,
	}
	return s.exercises.ListQuestions(ctx, filter)
}

func (s *exerciseService) History(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]model.ExerciseRecord, int64, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	return s.exercises.HistoryByUser(ctx, userID, page, pageSize)
}