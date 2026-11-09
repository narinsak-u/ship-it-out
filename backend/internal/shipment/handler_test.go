package shipment

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gofiber/fiber/v2"
	hub "github.com/narinsak-u/backend/internal/hub/domain"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Mock Repositories ----

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) List(_ context.Context, filter ShipmentFilter) ([]models.Shipment, int64, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Shipment), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) FindByOrderID(_ context.Context, orderID string) (*models.Shipment, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Shipment), args.Error(1)
}

func (m *mockRepo) FindByTrackingNumber(trackingNumber string) (*models.Shipment, error) {
	args := m.Called(trackingNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Shipment), args.Error(1)
}

func (m *mockRepo) Create(shipment *models.Shipment) error {
	args := m.Called(shipment)
	return args.Error(0)
}
func (m *mockRepo) CreateWithEvent(_ context.Context, shipment *models.Shipment, event *models.ShipmentEvent) error {
	if err := m.Create(shipment); err != nil {
		return err
	}
	event.ShipmentID = shipment.ID
	return m.CreateEvent(event)
}

func (m *mockRepo) SaveWithEvent(_ context.Context, shipment *models.Shipment, event *models.ShipmentEvent) error {
	if err := m.Save(shipment); err != nil || event == nil {
		return err
	}
	return m.CreateEvent(event)
}

func (m *mockRepo) Save(shipment *models.Shipment) error {
	args := m.Called(shipment)
	return args.Error(0)
}

func (m *mockRepo) Delete(shipment *models.Shipment) error {
	args := m.Called(shipment)
	return args.Error(0)
}

func (m *mockRepo) CreateEvent(event *models.ShipmentEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *mockRepo) FindEventsByShipmentID(_ context.Context, shipmentID uint) ([]models.ShipmentEvent, error) {
	args := m.Called(shipmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ShipmentEvent), args.Error(1)
}

func (m *mockRepo) DeleteWithEvents(_ context.Context, shipment *models.Shipment) error {
	if err := m.DeleteShipmentEvents(shipment.ID); err != nil {
		return err
	}
	return m.Delete(shipment)
}

func (m *mockRepo) DeleteShipmentEvents(shipmentID uint) error {
	args := m.Called(shipmentID)
	return args.Error(0)
}

type mockHubRepo struct {
	mock.Mock
}

func (m *mockHubRepo) FindByID(_ context.Context, id string) (*hub.Hub, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*hub.Hub), args.Error(1)
}

// ---- Tests ----

func TestShipmentList(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("List", mock.Anything).Return([]models.Shipment{
		{OrderID: "ORD-10245", TrackingNumber: "TH202600001", Status: "in_transit"},
	}, int64(1), nil)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Get("/api/shipments", h.List)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/shipments", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))
	data := result["data"].([]interface{})
	assert.Len(t, data, 1)
	repo.AssertExpectations(t)
}

func TestShipmentCreate_Success(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("Create", mock.Anything).Return(nil)
	repo.On("CreateEvent", mock.Anything).Return(nil)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Post("/api/shipments", h.Create)

	body := `{"customer":{"name":"John","zipcode":"10110","province":"Bangkok"},"receiver":{"name":"Jane","zipcode":"20110","province":"Chonburi"},"carrier":"Express","weight":5.5,"items":2}`
	req := httptest.NewRequest("POST", "/api/shipments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)

	assert.Equal(t, 200, resp.StatusCode)
	repo.AssertExpectations(t)
}

func TestShipmentCreate_ValidationErrors(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Post("/api/shipments", h.Create)

	tests := []struct {
		name string
		body string
	}{
		{"missing customer", `{"receiver":{"name":"J","zipcode":"1","province":"P"},"carrier":"C","weight":1,"items":1}`},
		{"missing receiver", `{"customer":{"name":"J","zipcode":"1","province":"P"},"carrier":"C","weight":1,"items":1}`},
		{"zero weight", `{"customer":{"name":"J","zipcode":"1","province":"P"},"receiver":{"name":"J","zipcode":"1","province":"P"},"carrier":"C","weight":0,"items":1}`},
		{"zero items", `{"customer":{"name":"J","zipcode":"1","province":"P"},"receiver":{"name":"J","zipcode":"1","province":"P"},"carrier":"C","weight":1,"items":0}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/shipments", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, 1000)
			assert.Equal(t, 400, resp.StatusCode)
		})
	}
}

func TestShipmentGetByID_Found(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("FindByOrderID", "ORD-10245").Return(&models.Shipment{
		OrderID: "ORD-10245", TrackingNumber: "TH202600001", Status: "in_transit",
	}, nil)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Get("/api/shipments/:orderId", h.GetByID)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/shipments/ORD-10245", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestShipmentGetByID_NotFound(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("FindByOrderID", "ORD-99999").Return(nil, ErrNotFound)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Get("/api/shipments/:orderId", h.GetByID)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/shipments/ORD-99999", nil), 1000)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestShipmentGetByID_RepositoryError(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByOrderID", "ORD-99999").Return(nil, assert.AnError)
	handler := NewHandler(NewService(repo, nil))
	app := fiber.New()
	app.Get("/api/shipments/:orderId", handler.GetByID)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/shipments/ORD-99999", nil), 1000)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestShipmentUpdate_Success(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("FindByOrderID", "ORD-10245").Return(&models.Shipment{
		OrderID: "ORD-10245", Status: "pending",
	}, nil)
	repo.On("Save", mock.Anything).Return(nil)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Put("/api/shipments/:orderId", h.Update)

	body := `{"carrier":"New Carrier"}`
	req := httptest.NewRequest("PUT", "/api/shipments/ORD-10245", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestShipmentUpdate_NotFound(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("FindByOrderID", "ORD-99999").Return(nil, ErrNotFound)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Put("/api/shipments/:orderId", h.Update)

	req := httptest.NewRequest("PUT", "/api/shipments/ORD-99999", strings.NewReader(`{"carrier":"X"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestShipmentUpdateStatus_WithHub(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("FindByOrderID", "ORD-10245").Return(&models.Shipment{
		OrderID: "ORD-10245", Status: "pending",
		Customer: models.ContactInfo{SubDistrict: "S", District: "D", Province: "P"},
		Receiver: models.ContactInfo{SubDistrict: "S2", District: "D2", Province: "P2"},
	}, nil)
	hubRepo.On("FindByID", "HUB-001").Return(&hub.Hub{
		ID: "HUB-001", Name: "Bangkok Hub", Address: "Bangkok",
		Coords: hub.GeoPoint{Lat: 13.7563, Lng: 100.5018},
	}, nil)
	repo.On("Save", mock.Anything).Return(nil)
	repo.On("CreateEvent", mock.Anything).Return(nil)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Patch("/api/shipments/:orderId/status", h.UpdateStatus)

	body := `{"status":"in_transit","hubId":"HUB-001"}`
	req := httptest.NewRequest("PATCH", "/api/shipments/ORD-10245/status", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 200, resp.StatusCode)
	repo.AssertExpectations(t)
	hubRepo.AssertExpectations(t)
}

func TestShipmentUpdateStatus_WithoutHub(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("FindByOrderID", "ORD-10245").Return(&models.Shipment{
		OrderID: "ORD-10245", Status: "pending",
		Customer: models.ContactInfo{SubDistrict: "S", District: "D", Province: "P"},
		Receiver: models.ContactInfo{SubDistrict: "S2", District: "D2", Province: "P2"},
	}, nil)
	repo.On("Save", mock.Anything).Return(nil)
	repo.On("CreateEvent", mock.Anything).Return(nil)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Patch("/api/shipments/:orderId/status", h.UpdateStatus)

	body := `{"status":"picked_up"}`
	req := httptest.NewRequest("PATCH", "/api/shipments/ORD-10245/status", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestShipmentDelete_Success(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("FindByOrderID", "ORD-10245").Return(&models.Shipment{
		OrderID: "ORD-10245", ID: 1,
	}, nil)
	repo.On("DeleteShipmentEvents", uint(1)).Return(nil)
	repo.On("Delete", mock.Anything).Return(nil)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Delete("/api/shipments/:orderId", h.Delete)

	resp, _ := app.Test(httptest.NewRequest("DELETE", "/api/shipments/ORD-10245", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestShipmentDelete_NotFound(t *testing.T) {
	repo := new(mockRepo)
	hubRepo := new(mockHubRepo)
	repo.On("FindByOrderID", "ORD-99999").Return(nil, ErrNotFound)

	app := fiber.New()
	h := NewHandler(NewService(repo, hubRepo))
	app.Delete("/api/shipments/:orderId", h.Delete)

	resp, _ := app.Test(httptest.NewRequest("DELETE", "/api/shipments/ORD-99999", nil), 1000)
	assert.Equal(t, 404, resp.StatusCode)
}
