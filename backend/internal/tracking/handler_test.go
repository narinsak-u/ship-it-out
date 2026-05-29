package tracking

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) List(f shipment.ShipmentFilter) ([]models.Shipment, int64, error) {
	args := m.Called(f)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Shipment), args.Get(1).(int64), args.Error(2)
}

func (m *mockRepo) FindByOrderID(id string) (*models.Shipment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Shipment), args.Error(1)
}

func (m *mockRepo) FindByTrackingNumber(tn string) (*models.Shipment, error) {
	args := m.Called(tn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Shipment), args.Error(1)
}

func (m *mockRepo) Create(*models.Shipment) error           { return m.Called().Error(0) }
func (m *mockRepo) Save(*models.Shipment) error             { return m.Called().Error(0) }
func (m *mockRepo) Delete(*models.Shipment) error           { return m.Called().Error(0) }
func (m *mockRepo) CreateEvent(*models.ShipmentEvent) error { return m.Called().Error(0) }

func (m *mockRepo) FindEventsByShipmentID(id uint) ([]models.ShipmentEvent, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ShipmentEvent), args.Error(1)
}

func (m *mockRepo) DeleteShipmentEvents(uint) error { return m.Called().Error(0) }
func (m *mockRepo) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockRepo) CountActive() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) CountByStatus() ([]shipment.StatusCountResult, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]shipment.StatusCountResult), args.Error(1)
}

func (m *mockRepo) CountByMonth() ([]shipment.MonthCountResult, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]shipment.MonthCountResult), args.Error(1)
}

func (m *mockRepo) CountByDayOfWeek() ([]shipment.DayCountResult, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]shipment.DayCountResult), args.Error(1)
}

func (m *mockRepo) CountByProvince() ([]shipment.ProvinceCountResult, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]shipment.ProvinceCountResult), args.Error(1)
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
	h := NewHandler(repo)
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
	repo.On("FindByTrackingNumber", "INVALID").Return(nil, assert.AnError)

	app := fiber.New()
	h := NewHandler(repo)
	app.Get("/api/track/:trackingNumber", h.Track)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/track/INVALID", nil), 1000)
	assert.Equal(t, 404, resp.StatusCode)
}
