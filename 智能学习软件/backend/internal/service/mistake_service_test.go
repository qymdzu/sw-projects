package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"smart-learning/internal/service"
)

func TestMistakeService_RecordMistake_New(t *testing.T) {
	qRepo := newMockExerciseRepo()
	q := qRepo.seed("choice", 1, "A")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewMistakeService(mRepo, qRepo)

	uid := uuid.New()
	require.NoError(t, svc.RecordMistake(context.Background(), uid, q.ID, "B"))

	m, err := mRepo.GetByUserQuestion(context.Background(), uid, q.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, m.MistakeCount)
}

func TestMistakeService_RecordMistake_DuplicateIncrement(t *testing.T) {
	qRepo := newMockExerciseRepo()
	q := qRepo.seed("choice", 1, "A")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewMistakeService(mRepo, qRepo)

	uid := uuid.New()
	require.NoError(t, svc.RecordMistake(context.Background(), uid, q.ID, "B"))
	require.NoError(t, svc.RecordMistake(context.Background(), uid, q.ID, "C"))

	m, err := mRepo.GetByUserQuestion(context.Background(), uid, q.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, m.MistakeCount)
}

func TestMistakeService_MarkMastered_Forbidden(t *testing.T) {
	qRepo := newMockExerciseRepo()
	q := qRepo.seed("choice", 1, "A")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewMistakeService(mRepo, qRepo)

	owner := uuid.New()
	require.NoError(t, svc.RecordMistake(context.Background(), owner, q.ID, "B"))
	m, _ := mRepo.GetByUserQuestion(context.Background(), owner, q.ID)

	err := svc.MarkMastered(context.Background(), uuid.New(), m.ID, true)
	assert.ErrorIs(t, err, service.ErrForbidden)
}

func TestMistakeService_MarkMastered_Success(t *testing.T) {
	qRepo := newMockExerciseRepo()
	q := qRepo.seed("choice", 1, "A")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewMistakeService(mRepo, qRepo)

	owner := uuid.New()
	require.NoError(t, svc.RecordMistake(context.Background(), owner, q.ID, "B"))
	m, _ := mRepo.GetByUserQuestion(context.Background(), owner, q.ID)

	require.NoError(t, svc.MarkMastered(context.Background(), owner, m.ID, true))
	updated, _ := mRepo.GetByID(context.Background(), m.ID)
	assert.True(t, updated.Mastered)
}

func TestMistakeService_Delete_Forbidden(t *testing.T) {
	qRepo := newMockExerciseRepo()
	q := qRepo.seed("choice", 1, "A")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewMistakeService(mRepo, qRepo)

	owner := uuid.New()
	require.NoError(t, svc.RecordMistake(context.Background(), owner, q.ID, "B"))
	m, _ := mRepo.GetByUserQuestion(context.Background(), owner, q.ID)

	err := svc.Delete(context.Background(), uuid.New(), m.ID)
	assert.ErrorIs(t, err, service.ErrForbidden)
}

func TestMistakeService_Delete_NotFound(t *testing.T) {
	qRepo := newMockExerciseRepo()
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewMistakeService(mRepo, qRepo)

	err := svc.Delete(context.Background(), uuid.New(), 9999)
	assert.Error(t, err)
}

func TestMistakeService_Review(t *testing.T) {
	qRepo := newMockExerciseRepo()
	qRepo.seed("choice", 1, "A")
	mRepo := newMockMistakeRepoForExercise()
	svc := service.NewMistakeService(mRepo, qRepo)

	resp, err := svc.Review(context.Background(), uuid.New(), service.ReviewRequest{Count: 5})
	require.NoError(t, err)
	assert.Equal(t, "unmastered_only", resp.Strategy)
}