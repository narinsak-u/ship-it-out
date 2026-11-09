package presentation

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/hub/application"
	"github.com/narinsak-u/backend/internal/hub/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type Hub = domain.Hub
type GeoPoint = domain.GeoPoint
var NewService = application.NewService
var ErrNotFound = application.ErrNotFound

type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindAll(_ context.Context) ([]Hub, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Hub), args.Error(1)
}

func (m *mockRepo) FindByID(_ context.Context, id string) (*Hub, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Hub), args.Error(1)
}

func (m *mockRepo) Create(_ context.Context, hub *Hub) error {
	args := m.Called(hub)
	return args.Error(0)
}

func (m *mockRepo) Save(_ context.Context, hub *Hub) error {
	args := m.Called(hub)
	return args.Error(0)
}

func (m *mockRepo) Delete(_ context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestHubList_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindAll").Return([]Hub{
		{ID: "HUB-001", Name: "Test Hub", Status: "active"},
	}, nil)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Get("/api/hubs", h.List)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/hubs", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
	repo.AssertExpectations(t)
}

func TestHubGetByID_Found(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(&Hub{ID: "HUB-001", Name: "Test"}, nil)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Get("/api/hubs/:id", h.GetByID)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/hubs/HUB-001", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHubGetByID_NotFound(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-999").Return(nil, ErrNotFound)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Get("/api/hubs/:id", h.GetByID)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/hubs/HUB-999", nil), 1000)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestHubCreate_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("Create", mock.MatchedBy(func(h *Hub) bool {
		return h.Name == "New Hub"
	})).Return(nil)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Post("/api/hubs", h.Create)

	body := `{"name":"New Hub","carrierId":"CAR-001","address":"Addr","coords":{"lat":10,"lng":20},"capacity":100,"currentUtilization":50,"status":"active"}`
	req := httptest.NewRequest("POST", "/api/hubs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 1000)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var response struct {
		Success bool `json:"success"`
		Data    Hub  `json:"data"`
	}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
	assert.True(t, response.Success)
	assert.Equal(t, (GeoPoint{Lat: 10, Lng: 20}), response.Data.Coords)
}

func TestHubCreate_InvalidBody(t *testing.T) {
	repo := new(mockRepo)
	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Post("/api/hubs", h.Create)

	req := httptest.NewRequest("POST", "/api/hubs", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestHubUpdate_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(&Hub{ID: "HUB-001"}, nil).Twice()
	repo.On("Save", mock.MatchedBy(func(h *Hub) bool {
		return h.ID == "HUB-001" && h.Name == "Updated"
	})).Return(nil)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Put("/api/hubs/:id", h.Update)

	body := `{"name":"Updated"}`
	req := httptest.NewRequest("PUT", "/api/hubs/HUB-001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHubUpdate_NotFound(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-999").Return(nil, ErrNotFound)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Put("/api/hubs/:id", h.Update)

	body := `not json`
	req := httptest.NewRequest("PUT", "/api/hubs/HUB-999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestHubDelete_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("Delete", "HUB-001").Return(nil)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Delete("/api/hubs/:id", h.Delete)

	resp, _ := app.Test(httptest.NewRequest("DELETE", "/api/hubs/HUB-001", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHubUpdate_InvalidBody(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(&Hub{ID: "HUB-001"}, nil)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Put("/api/hubs/:id", h.Update)

	req := httptest.NewRequest("PUT", "/api/hubs/HUB-001", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestHubDelete_Error(t *testing.T) {
	repo := new(mockRepo)
	repo.On("Delete", "HUB-999").Return(gorm.ErrRecordNotFound)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Delete("/api/hubs/:id", h.Delete)

	resp, _ := app.Test(httptest.NewRequest("DELETE", "/api/hubs/HUB-999", nil), 1000)
	assert.Equal(t, 500, resp.StatusCode)
}
