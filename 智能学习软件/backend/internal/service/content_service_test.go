package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"smart-learning/internal/model"
	"smart-learning/internal/repository"
	"smart-learning/internal/service"
)

// mockSubjectRepoForContent 是 SubjectRepository 的 mock。
type mockSubjectRepoForContent struct {
	items []model.Subject
}

func (m *mockSubjectRepoForContent) List(_ context.Context) ([]model.Subject, error) {
	return m.items, nil
}

// mockKnowledgeRepoForContent 是 KnowledgeRepository 的 mock。
type mockKnowledgeRepoForContent struct {
	tree []repository.KnowledgeNode
}

func (m *mockKnowledgeRepoForContent) ListTree(_ context.Context, _ uint64) ([]repository.KnowledgeNode, error) {
	return m.tree, nil
}

func (m *mockKnowledgeRepoForContent) GetByID(_ context.Context, id uint64) (*model.KnowledgePoint, error) {
	return &model.KnowledgePoint{ID: id, Name: "测试 KP"}, nil
}

func (m *mockKnowledgeRepoForContent) ListBySubject(_ context.Context, _ uint64) ([]model.KnowledgePoint, error) {
	return nil, nil
}

func TestSubjectService_List(t *testing.T) {
	repo := &mockSubjectRepoForContent{items: []model.Subject{
		{ID: 1, Name: "数学"},
		{ID: 2, Name: "英语"},
	}}
	svc := service.NewSubjectService(repo)
	items, err := svc.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestKnowledgeService_GetTree(t *testing.T) {
	repo := &mockKnowledgeRepoForContent{
		tree: []repository.KnowledgeNode{
			{ID: 1, Name: "数与代数", Level: 1, Children: []repository.KnowledgeNode{
				{ID: 3, Name: "一元二次方程", Level: 2},
			}},
		},
	}
	svc := service.NewKnowledgeService(repo)
	tree, err := svc.GetTree(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, tree, 1)
	assert.Len(t, tree[0].Children, 1)
}

func TestKnowledgeService_GetTree_EmptySubjectID(t *testing.T) {
	repo := &mockKnowledgeRepoForContent{}
	svc := service.NewKnowledgeService(repo)
	_, err := svc.GetTree(context.Background(), 0)
	assert.Error(t, err)
}