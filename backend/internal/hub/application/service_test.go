package application

import (
	"context"
	"errors"
	"testing"

	"github.com/narinsak-u/backend/internal/hub/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindAll(_ context.Context) ([]domain.Hub, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Hub), args.Error(1)
}

func (m *mockRepo) FindByID(_ context.Context, id string) (*domain.Hub, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Hub), args.Error(1)
}

func (m *mockRepo) Create(_ context.Context, h *domain.Hub) error {
	return m.Called(h).Error(0)
}

func (m *mockRepo) Save(_ context.Context, h *domain.Hub) error {
	return m.Called(h).Error(0)
}

func (m *mockRepo) Delete(_ context.Context, id string) error {
	return m.Called(id).Error(0)
}


func TestServiceCRUDDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	h := &domain.Hub{ID: "HUB-001", Name: "Bangkok Hub"}
	repo.On("FindAll").Return([]domain.Hub{*h}, nil)
	repo.On("FindByID", h.ID).Return(h, nil)
	repo.On("Create", h).Return(nil)
	repo.On("Delete", h.ID).Return(nil)

	service := NewService(repo)
	hubs, err := service.List(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []domain.Hub{*h}, hubs)
	found, err := service.GetByID(ctx, h.ID)
	assert.NoError(t, err)
	assert.Equal(t, h, found)
	assert.NoError(t, service.Create(ctx, h))
	assert.NoError(t, service.Delete(ctx, h.ID))
	repo.AssertExpectations(t)
}

func TestServiceUpdatePreservesURLID(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(&domain.Hub{ID: "HUB-001"}, nil)
	repo.On("Save", mock.MatchedBy(func(h *domain.Hub) bool {
		return h.ID == "HUB-001" && h.Name == "Updated"
	})).Return(nil)

	err := NewService(repo).Update(context.Background(), "HUB-001", &domain.Hub{ID: "HUB-999", Name: "Updated"})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestServiceUpdatePropagatesRepositoryError(t *testing.T) {
	repo := new(mockRepo)
	want := errors.New("database unavailable")
	repo.On("FindByID", "HUB-001").Return(nil, want)

	err := NewService(repo).Update(context.Background(), "HUB-001", &domain.Hub{})

	assert.ErrorIs(t, err, want)
	repo.AssertExpectations(t)
}
