package hub

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestServiceCRUDDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	h := &Hub{ID: "HUB-001", Name: "Bangkok Hub"}
	repo.On("FindAll").Return([]Hub{*h}, nil)
	repo.On("FindByID", h.ID).Return(h, nil)
	repo.On("Create", h).Return(nil)
	repo.On("Delete", h.ID).Return(nil)

	service := NewService(repo)
	hubs, err := service.List(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []Hub{*h}, hubs)
	found, err := service.GetByID(ctx, h.ID)
	assert.NoError(t, err)
	assert.Equal(t, h, found)
	assert.NoError(t, service.Create(ctx, h))
	assert.NoError(t, service.Delete(ctx, h.ID))
	repo.AssertExpectations(t)
}

func TestServiceUpdatePreservesURLID(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(&Hub{ID: "HUB-001"}, nil)
	repo.On("Save", mock.MatchedBy(func(h *Hub) bool {
		return h.ID == "HUB-001" && h.Name == "Updated"
	})).Return(nil)

	err := NewService(repo).Update(context.Background(), "HUB-001", &Hub{ID: "HUB-999", Name: "Updated"})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestServiceUpdatePropagatesRepositoryError(t *testing.T) {
	repo := new(mockRepo)
	want := errors.New("database unavailable")
	repo.On("FindByID", "HUB-001").Return(nil, want)

	err := NewService(repo).Update(context.Background(), "HUB-001", &Hub{})

	assert.ErrorIs(t, err, want)
	repo.AssertExpectations(t)
}
