package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/narinsak-u/backend/internal/analytics/application"
)

type fakeOverviewReader struct{}

func (fakeOverviewReader) Count(_ context.Context) (int64, error) { return 10, nil }
func (fakeOverviewReader) CountActive(_ context.Context) (int64, error) { return 7, nil }
func (fakeOverviewReader) CountByStatus(_ context.Context) ([]application.StatusCount, error) {
	return []application.StatusCount{{Status: "delivered", Count: 2}}, nil
}
func (fakeOverviewReader) CountByProvince(_ context.Context) ([]application.ProvinceCount, error) {
	return []application.ProvinceCount{{Province: "ชลบุรี", Total: 5}}, nil
}

type fakeTimeSeriesReader struct{}

func (fakeTimeSeriesReader) CountByMonth(_ context.Context) ([]application.MonthCount, error) {
	return []application.MonthCount{{Month: "2026-01", Count: 3}}, nil
}
func (fakeTimeSeriesReader) CountByDayOfWeek(_ context.Context) ([]application.DayCount, error) {
	return []application.DayCount{{Day: "Monday", Count: 2}}, nil
}

type failingOverviewReader struct{ fakeOverviewReader }

func (failingOverviewReader) Count(context.Context) (int64, error) { return 0, errors.New("database unavailable") }

func TestAnalyticsRoutesRemainPublicAndKeepResponseShapes(t *testing.T) {
	service := application.NewService(fakeOverviewReader{}, fakeTimeSeriesReader{})
	handler := NewHandler(service)
	app := fiber.New()
	app.Get("/api/analytics/overview", handler.Overview)
	app.Get("/api/analytics/timeseries", handler.TimeSeries)

	for _, test := range []struct {
		path string
		check func(*testing.T, map[string]any)
	}{
		{path: "/api/analytics/overview", check: func(t *testing.T, data map[string]any) {
			assert.Equal(t, float64(10), data["total"])
			assert.Equal(t, float64(7), data["active"])
			assert.Equal(t, float64(2), data["delivered"])
			assert.NotNil(t, data["by_status"])
			assert.NotNil(t, data["by_region"])
		}},
		{path: "/api/analytics/timeseries", check: func(t *testing.T, data map[string]any) {
			assert.Equal(t, "2026-01", data["by_month"].([]any)[0].(map[string]any)["month"])
			assert.Equal(t, "Monday", data["by_day_of_week"].([]any)[0].(map[string]any)["day"])
		}},
	} {
		resp, err := app.Test(httptest.NewRequest("GET", test.path, nil), 1000)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		var result map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.True(t, result["success"].(bool))
		test.check(t, result["data"].(map[string]any))
	}
}

func TestAnalyticsHandlerMapsQueryErrors(t *testing.T) {
	service := application.NewService(failingOverviewReader{}, fakeTimeSeriesReader{})
	app := fiber.New()
	app.Get("/overview", NewHandler(service).Overview)
	resp, err := app.Test(httptest.NewRequest("GET", "/overview", nil), 1000)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}
