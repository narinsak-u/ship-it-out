package hub

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHubGetByID_HidesRepositoryError(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(nil, errors.New("database unavailable"))

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Get("/api/hubs/:id", h.GetByID)

	resp, err := app.Test(httptest.NewRequest("GET", "/api/hubs/HUB-001", nil), 1000)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
	repo.AssertExpectations(t)
}
