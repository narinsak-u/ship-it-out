package analytics

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

type mockRepo struct{ mock.Mock }

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

func TestOverview(t *testing.T) {
	repo := new(mockRepo)
	repo.On("Count").Return(int64(10), nil)
	repo.On("CountActive").Return(int64(7), nil)
	repo.On("CountByStatus").Return([]shipment.StatusCountResult{
		{Status: "in_transit", Count: 4},
		{Status: "pending", Count: 3},
		{Status: "delivered", Count: 2},
		{Status: "delayed", Count: 1},
	}, nil)
	repo.On("CountByProvince").Return([]shipment.ProvinceCountResult{
		{Province: "ชลบุรี", Total: 5},
		{Province: "กรุงเทพมหานคร", Total: 5},
	}, nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Get("/api/analytics/overview", h.Overview)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/analytics/overview", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))
	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(10), data["total"])
	assert.Equal(t, float64(7), data["active"])
	assert.Equal(t, float64(2), data["delivered"])
	assert.NotNil(t, data["by_status"])
	assert.NotNil(t, data["by_region"])
	repo.AssertExpectations(t)
}

func TestTimeSeries(t *testing.T) {
	repo := new(mockRepo)
	repo.On("CountByMonth").Return([]shipment.MonthCountResult{
		{Month: "2026-01", Count: 3},
		{Month: "2026-02", Count: 5},
	}, nil)
	repo.On("CountByDayOfWeek").Return([]shipment.DayCountResult{
		{Day: "Monday", Count: 2},
		{Day: "Wednesday", Count: 4},
	}, nil)

	app := fiber.New()
	h := NewHandler(repo)
	app.Get("/api/analytics/timeseries", h.TimeSeries)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/analytics/timeseries", nil), 1000)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))
	repo.AssertExpectations(t)
}
