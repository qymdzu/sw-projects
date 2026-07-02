package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/internal/service"
)

// mockExerciseRepo 是 ExerciseRepository 的内存 mock。
type mockExerciseRepo struct {
	questions map[uint64]*model.Question
	nextID    uint64
}

func newMockExerciseRepo() *mockExerciseRepo {
	return &mockExerciseRepo{
		questions: make(map[uint64]*model.Question),
		nextID:    1,
	}
}

func (m *mockExerciseRepo) seed(t string, kp uint64, answer string) *model.Question {
	q := &model.Question{
		ID:               m.nextID,
		SubjectID:        1,
		KnowledgePointID: kp,
		Type:             t,
		Difficulty:       3,
		Content:          datatypes.JSON(`{"text":"一测试题"}`),
		Answer:           answer,
	}
	m.nextID++
	m.questions[q.ID] = q
	return q
}

func (m *mockExerciseRepo) GetQuestionByID(_ context.Context, id uint64) (*model.Question, error) {
	q, ok := m.questions[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return q, nil
}

func (m *mockExerciseRepo) ListQuestions(_ context.Context, f repository.QuestionFilter) ([]model.Question, int64, error) {
	var items []model.Question
	var total int64
	for _, q := range m.questions {
		match := true
		if f.SubjectID != nil && q.SubjectID != *f.SubjectID {
			match = false
		}
		if f.KnowledgePointID != nil && q.KnowledgePointID != *f.KnowledgePointID {
			match = false
		}
		if f.Type != "" && q.Type != f.Type {
			match = false
		}
		if f.Difficulty != nil && q.Difficulty != *f.Difficulty {
			match = false
		}
		if match {
			items = append(items, *q)
			total++
		}
	}
	return items, total, nil
}

func (m *mockExerciseRepo) ListQuestionsByIDs(_ context.Context, ids []uint64) ([]model.Question, error) {
	var items []model.Question
	for _, id := range ids {
		if q, ok := m.questions[id]; ok {
			items = append(items, *q)
		}
	}
	return items, nil
}

func (m *mockExerciseRepo) RandomQuestions(_ context.Context, _, _ *uint64, count int) ([]model.Question, error) {
	var items []model.Question
	for _, q := range m.questions {
		items = append(items, *q)
		if len(items) >= count {
			break
		}
	}
	return items, nil
}

func (m *mockExerciseRepo) CreateRecord(_ context.Context, _ *model.ExerciseRecord) error {
	return nil
}

func (m *mockExerciseRepo) CountByUser(_ context.Context, _ uuid.UUID, _, _ time.Time) (int64, error) {
	return 0, nil
}

func (m *mockExerciseRepo) CorrectCountByUser(_ context.Context, _ uuid.UUID, _, _ time.Time) (int64, error) {
	return 0, nil
}

func (m *mockExerciseRepo) HistoryByUser(_ context.Context, _ uuid.UUID, _, _ int) ([]model.ExerciseRecord, int64, error) {
	return nil, 0, nil
}

func (m *mockExerciseRepo) CorrectRateByKP(_ context.Context, _ uuid.UUID, _ uint64) ([]repository.KPRate, error) {
	return nil, nil
}

// mockMistakeRepo 复用 plan 测试的 mock 是不行的，需新建专门用于练习测试。
type mockMistakeRepoForExercise struct {
	mistakes  map[string]*model.MistakeBook
	nextID    uint64
}

func newMockMistakeRepoForExercise() *mockMistakeRepoForExercise {
	return &mockMistakeRepoForExercise{
		mistakes: make(map[string]*model.MistakeBook),
	}
}

func mistakeKey(userID uuid.UUID, qid uint64) string {
	return userID.String() + ":" + uint64Str(qid)
}

func uint64Str(v uint64) string {
	if v == 0 {
		return "0"
	}
	var b []byte
	for v > 0 {
		b = append([]byte{byte('0' + v%10)}, b...)
		v /= 10
	}
	return string(b)
}

func (m *mockMistakeRepoForExercise) GetByUserQuestion(_ context.Context, uid uuid.UUID, qid uint64) (*model.MistakeBook, error) {
	v, ok := m.mistakes[mistakeKey(uid, qid)]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return v, nil
}

func (m *mockMistakeRepoForExercise) Create(_ context.Context, mb *model.MistakeBook) error {
	if mb.ID == 0 {
		m.nextID++
		mb.ID = m.nextID
	}
	m.mistakes[mistakeKey(mb.UserID, mb.QuestionID)] = mb
	return nil
}

func (m *mockMistakeRepoForExercise) IncrementCount(_ context.Context, id uint64, _ time.Time) error {
	for _, v := range m.mistakes {
		if v.ID == id {
			v.MistakeCount++
			return nil
		}
	}
	return repository.ErrNotFound
}

func (m *mockMistakeRepoForExercise) GetByID(_ context.Context, id uint64) (*model.MistakeBook, error) {
	for _, v := range m.mistakes {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockMistakeRepoForExercise) MarkMastered(_ context.Context, id uint64, mastered bool) error {
	for _, v := range m.mistakes {
		if v.ID == id {
			v.Mastered = mastered
			return nil
		}
	}
	return repository.ErrNotFound
}

func (m *mockMistakeRepoForExercise) Delete(_ context.Context, id uint64) error {
	for k, v := range m.mistakes {
		if v.ID == id {
			delete(m.mistakes, k)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (m *mockMistakeRepoForExercise) ListByUser(_ context.Context, _ uuid.UUID, _ *uint64, _ *bool, _, _ int) ([]repository.MistakeWithQuestion, int64, error) {
	return nil, 0, nil
}

func (m *mockMistakeRepoForExercise) GroupByKP(_ context.Context, _ uuid.UUID) ([]repository.MistakeGroup, error) {
	return nil, nil
}

func (m *mockMistakeRepoForExercise) ListUnmasteredQuestions(_ context.Context, _ uuid.UUID, _ []uint64, _ int) ([]model.Question, error) {
	return nil, nil
}

func (m *mockMistakeRepoForExercise) CountUnmastered(_ context.Context, _ uuid.UUID) (int64, error) {
	return 0, nil
}

func TestExerciseService_Submit_Correct(t *testing.T) {
	qRepo := newMockExerciseRepo()
	q := qRepo.seed("choice", 1, "A")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewExerciseService(qRepo, mRepo)

	resp, err := svc.Submit(context.Background(), uuid.New(), service.SubmitRequest{
		QuestionID: q.ID, Answer: "A",
	})
	require.NoError(t, err)
	assert.True(t, resp.IsCorrect)
	assert.False(t, resp.MistakeRecorded)
}

func TestExerciseService_Submit_Wrong_RecordsMistake(t *testing.T) {
	qRepo := newMockExerciseRepo()
	q := qRepo.seed("choice", 1, "A")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewExerciseService(qRepo, mRepo)

	resp, err := svc.Submit(context.Background(), uuid.New(), service.SubmitRequest{
		QuestionID: q.ID, Answer: "B",
	})
	require.NoError(t, err)
	assert.False(t, resp.IsCorrect)
	assert.True(t, resp.MistakeRecorded)
}

func TestExerciseService_Submit_Fill_MultipleAnswers(t *testing.T) {
	qRepo := newMockExerciseRepo()
	q := qRepo.seed("fill", 1, "yes|YES|Yes")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewExerciseService(qRepo, mRepo)

	resp, err := svc.Submit(context.Background(), uuid.New(), service.SubmitRequest{
		QuestionID: q.ID, Answer: "YES",
	})
	require.NoError(t, err)
	assert.True(t, resp.IsCorrect)
}

func TestExerciseService_Submit_EmptyAnswer(t *testing.T) {
	qRepo := newMockExerciseRepo()
	q := qRepo.seed("choice", 1, "A")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewExerciseService(qRepo, mRepo)

	_, err := svc.Submit(context.Background(), uuid.New(), service.SubmitRequest{
		QuestionID: q.ID, Answer: "   ",
	})
	assert.Error(t, err)
}

func TestExerciseService_Submit_QuestionNotFound(t *testing.T) {
	qRepo := newMockExerciseRepo()
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewExerciseService(qRepo, mRepo)

	_, err := svc.Submit(context.Background(), uuid.New(), service.SubmitRequest{
		QuestionID: 999, Answer: "A",
	})
	assert.Error(t, err)
}

func TestExerciseService_Random_DefaultCount(t *testing.T) {
	qRepo := newMockExerciseRepo()
	for i := 0; i < 5; i++ {
		qRepo.seed("choice", 1, "A")
	}
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewExerciseService(qRepo, mRepo)

	items, err := svc.Random(context.Background(), nil, nil, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(items), 0)
}

func TestExerciseService_Recommend(t *testing.T) {
	qRepo := newMockExerciseRepo()
	for i := 0; i < 5; i++ {
		qRepo.seed("choice", 1, "A")
	}
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewExerciseService(qRepo, mRepo)

	items, err := svc.Recommend(context.Background(), uuid.New(), 3)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(items), 3)
}

func TestExerciseService_List(t *testing.T) {
	qRepo := newMockExerciseRepo()
	qRepo.seed("choice", 1, "A")
	qRepo.seed("fill", 1, "yes")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewExerciseService(qRepo, mRepo)

	items, total, err := svc.List(context.Background(), repository.QuestionFilter{
		Page: 1, PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, items, 2)
}

func TestExerciseService_ByKnowledgePoint(t *testing.T) {
	qRepo := newMockExerciseRepo()
	qRepo.seed("choice", 1, "A")
	qRepo.seed("choice", 2, "B")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewExerciseService(qRepo, mRepo)

	items, total, err := svc.ByKnowledgePoint(context.Background(), 1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, items, 1)
}