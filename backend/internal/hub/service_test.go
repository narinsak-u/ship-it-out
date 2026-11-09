package hub

import (
	"context"
	"errors"
	"testing"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceUpdatePreservesURLID(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(&models.Hub{ID: "HUB-001"}, nil)
	repo.On("Save", mock.MatchedBy(func(h *models.Hub) bool {
		return h.ID == "HUB-001" && h.Name == "Updated"
	})).Return(nil)

	err := NewService(repo).Update(context.Background(), "HUB-001", &models.Hub{ID: "HUB-999", Name: "Updated"})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestServiceUpdatePropagatesRepositoryError(t *testing.T) {
	repo := new(mockRepo)
	want := errors.New("database unavailable")
	repo.On("FindByID", "HUB-001").Return(nil, want)

	err := NewService(repo).Update(context.Background(), "HUB-001", &models.Hub{})

	assert.ErrorIs(t, err, want)
	repo.AssertExpectations(t)
}
