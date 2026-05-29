package hub

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindAll() ([]models.Hub, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Hub), args.Error(1)
}

func (m *mockRepo) FindByID(id string) (*models.Hub, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Hub), args.Error(1)
}

func (m *mockRepo) Create(hub *models.Hub) error {
	args := m.Called(hub)
	return args.Error(0)
}

func (m *mockRepo) Save(hub *models.Hub) error {
	args := m.Called(hub)
	return args.Error(0)
}

func (m *mockRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestHubList_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindAll").Return([]models.Hub{
		{ID: "HUB-001", Name: "Test Hub", Status: "active"},
	}, nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Get("/api/hubs", h.List)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/hubs", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
	repo.AssertExpectations(t)
}

func TestHubGetByID_Found(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(&models.Hub{ID: "HUB-001", Name: "Test"}, nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Get("/api/hubs/:id", h.GetByID)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/hubs/HUB-001", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHubGetByID_NotFound(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-999").Return(nil, gorm.ErrRecordNotFound)

	app := fiber.New()
	h := NewHandler(repo)
	app.Get("/api/hubs/:id", h.GetByID)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/hubs/HUB-999", nil), 1000)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestHubCreate_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("Create", mock.MatchedBy(func(h *models.Hub) bool {
		return h.Name == "New Hub"
	})).Return(nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Post("/api/hubs", h.Create)

	body := `{"name":"New Hub","carrierId":"CAR-001","address":"Addr","coords":{"lat":10,"lng":20},"capacity":100,"currentUtilization":50,"status":"active"}`
	req := httptest.NewRequest("POST", "/api/hubs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHubCreate_InvalidBody(t *testing.T) {
	repo := new(mockRepo)
	app := fiber.New()
	h := NewHandler(repo)
	app.Post("/api/hubs", h.Create)

	req := httptest.NewRequest("POST", "/api/hubs", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestHubUpdate_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(&models.Hub{ID: "HUB-001"}, nil)
	repo.On("Save", mock.MatchedBy(func(h *models.Hub) bool {
		return h.ID == "HUB-001" && h.Name == "Updated"
	})).Return(nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Put("/api/hubs/:id", h.Update)

	body := `{"name":"Updated"}`
	req := httptest.NewRequest("PUT", "/api/hubs/HUB-001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHubUpdate_NotFound(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-999").Return(nil, gorm.ErrRecordNotFound)

	app := fiber.New()
	h := NewHandler(repo)
	app.Put("/api/hubs/:id", h.Update)

	body := `{"name":"Updated"}`
	req := httptest.NewRequest("PUT", "/api/hubs/HUB-999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestHubDelete_Success(t *testing.T) {
	repo := new(mockRepo)
	repo.On("Delete", "HUB-001").Return(nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Delete("/api/hubs/:id", h.Delete)

	resp, _ := app.Test(httptest.NewRequest("DELETE", "/api/hubs/HUB-001", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHubUpdate_InvalidBody(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByID", "HUB-001").Return(&models.Hub{ID: "HUB-001"}, nil)

	app := fiber.New()
	h := NewHandler(repo)
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
	h := NewHandler(repo)
	app.Delete("/api/hubs/:id", h.Delete)

	resp, _ := app.Test(httptest.NewRequest("DELETE", "/api/hubs/HUB-999", nil), 1000)
	assert.Equal(t, 500, resp.StatusCode)
}
