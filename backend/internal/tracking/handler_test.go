package tracking

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
	shipmentapplication "github.com/narinsak-u/backend/internal/shipment/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) FindByTrackingNumber(_ context.Context, tn string) (*models.Shipment, error) {
	args := m.Called(tn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Shipment), args.Error(1)
}

func (m *mockRepo) FindEventsByShipmentID(_ context.Context, id uint) ([]models.ShipmentEvent, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ShipmentEvent), args.Error(1)
}

func TestTrack_Found(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByTrackingNumber", "TH202600001").Return(&models.Shipment{
		ID: 1, OrderID: "ORD-10245", TrackingNumber: "TH202600001", Status: "in_transit",
	}, nil)
	repo.On("FindEventsByShipmentID", uint(1)).Return([]models.ShipmentEvent{
		{Status: "Label Created", Description: "Awaiting pickup."},
	}, nil)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Get("/api/track/:trackingNumber", h.Track)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/track/TH202600001", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))
	data := result["data"].(map[string]interface{})
	assert.NotNil(t, data["shipment"])
	assert.NotNil(t, data["events"])
	repo.AssertExpectations(t)
}

func TestTrack_NotFound(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByTrackingNumber", "INVALID").Return(nil, shipmentapplication.ErrNotFound)

	app := fiber.New()
	h := NewHandler(NewService(repo))

	app.Get("/api/track/:trackingNumber", h.Track)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/track/INVALID", nil), 1000)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestTrack_LookupFailure(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByTrackingNumber", "TH202600001").Return(nil, assert.AnError)

	app := fiber.New()
	app.Get("/api/track/:trackingNumber", NewHandler(NewService(repo)).Track)
	resp, _ := app.Test(httptest.NewRequest("GET", "/api/track/TH202600001", nil), 1000)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestTrack_EventQueryFailure(t *testing.T) {
	repo := new(mockRepo)
	repo.On("FindByTrackingNumber", "TH202600001").Return(&models.Shipment{ID: 1}, nil)
	repo.On("FindEventsByShipmentID", uint(1)).Return(nil, assert.AnError)

	app := fiber.New()
	h := NewHandler(NewService(repo))
	app.Get("/api/track/:trackingNumber", h.Track)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/track/TH202600001", nil), 1000)
	assert.Equal(t, 500, resp.StatusCode)
}
