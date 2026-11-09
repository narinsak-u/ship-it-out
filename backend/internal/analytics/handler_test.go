package analytics

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverviewRouteContract(t *testing.T) {
	repo := fakeRepository{
		total: 10, active: 7,
		statuses: []StatusCount{{Status: "in_transit", Count: 4}, {Status: "delivered", Count: 2}},
		provinces: []ProvinceCount{{Province: "ชลบุรี", Total: 5}, {Province: "กรุงเทพมหานคร", Total: 5}},
	}
	app := fiber.New()
	app.Get("/api/analytics/overview", NewHandler(NewService(repo)).Overview)
	resp, err := app.Test(httptest.NewRequest("GET", "/api/analytics/overview", nil), 1000)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.True(t, result["success"].(bool))
	data := result["data"].(map[string]any)
	assert.Equal(t, float64(10), data["total"])
	assert.Equal(t, float64(7), data["active"])
	assert.Equal(t, float64(2), data["delivered"])
	assert.NotNil(t, data["by_status"])
	assert.NotNil(t, data["by_region"])
}

func TestTimeSeriesRouteContract(t *testing.T) {
	repo := fakeRepository{
		months: []MonthCount{{Month: "2026-01", Count: 3}},
		days: []DayCount{{Day: "Monday", Count: 2}},
	}
	app := fiber.New()
	app.Get("/api/analytics/timeseries", NewHandler(NewService(repo)).TimeSeries)
	resp, err := app.Test(httptest.NewRequest("GET", "/api/analytics/timeseries", nil), 1000)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.True(t, result["success"].(bool))
	data := result["data"].(map[string]any)
	assert.Equal(t, "2026-01", data["by_month"].([]any)[0].(map[string]any)["month"])
	assert.Equal(t, "Monday", data["by_day_of_week"].([]any)[0].(map[string]any)["day"])
}

func TestAnalyticsRoutesMapQueryFailures(t *testing.T) {
	repo := fakeRepository{err: errors.New("database unavailable")}
	app := fiber.New()
	handler := NewHandler(NewService(repo))
	app.Get("/overview", handler.Overview)
	app.Get("/timeseries", handler.TimeSeries)
	for _, path := range []string{"/overview", "/timeseries"} {
		resp, err := app.Test(httptest.NewRequest("GET", path, nil), 1000)
		require.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
	}
}
