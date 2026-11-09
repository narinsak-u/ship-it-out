package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeOverviewReader struct {
	total     int64
	active    int64
	statuses  []StatusCount
	provinces []ProvinceCount
	err       error
}

func (r fakeOverviewReader) Count(context.Context) (int64, error) { return r.total, r.err }
func (r fakeOverviewReader) CountActive(context.Context) (int64, error) { return r.active, r.err }
func (r fakeOverviewReader) CountByStatus(context.Context) ([]StatusCount, error) { return r.statuses, r.err }
func (r fakeOverviewReader) CountByProvince(context.Context) ([]ProvinceCount, error) { return r.provinces, r.err }

type fakeTimeSeriesReader struct {
	months []MonthCount
	days   []DayCount
	err    error
}

func (r fakeTimeSeriesReader) CountByMonth(context.Context) ([]MonthCount, error) { return r.months, r.err }
func (r fakeTimeSeriesReader) CountByDayOfWeek(context.Context) ([]DayCount, error) { return r.days, r.err }

func TestServiceOverviewAggregatesDeliveredAndProvinceRegions(t *testing.T) {
	service := NewService(fakeOverviewReader{
		total: 10, active: 7,
		statuses: []StatusCount{{Status: "delivered", Count: 2}, {Status: "pending", Count: 8}},
		provinces: []ProvinceCount{{Province: "ชลบุรี", Total: 5}, {Province: "กรุงเทพมหานคร", Total: 5}, {Province: "unknown", Total: 7}},
	}, fakeTimeSeriesReader{})

	got, err := service.Overview(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(10), got.Total)
	assert.Equal(t, int64(7), got.Active)
	assert.Equal(t, int64(2), got.Delivered)
	assert.Equal(t, []StatusCount{{Status: "delivered", Count: 2}, {Status: "pending", Count: 8}}, got.ByStatus)
	regions := make(map[string]int64, len(got.ByRegion))
	for _, region := range got.ByRegion {
		regions[region.Name] = region.Total
	}
	assert.Equal(t, map[string]int64{
		"Central": 5, "East": 5, "North": 0, "West": 0, "North-east": 0, "South": 0,
	}, regions)
}

func TestServiceTimeSeriesReturnsReadResults(t *testing.T) {
	months := []MonthCount{{Month: "2026-01", Count: 3}}
	days := []DayCount{{Day: "Monday", Count: 2}}
	got, err := NewService(fakeOverviewReader{}, fakeTimeSeriesReader{months: months, days: days}).TimeSeries(context.Background())
	require.NoError(t, err)
	assert.Equal(t, months, got.ByMonth)
	assert.Equal(t, days, got.ByDayOfWeek)
}

func TestServicePropagatesReadErrors(t *testing.T) {
	want := errors.New("database unavailable")
	service := NewService(fakeOverviewReader{err: want}, fakeTimeSeriesReader{err: want})
	_, err := service.Overview(context.Background())
	assert.ErrorIs(t, err, want)
	_, err = service.TimeSeries(context.Background())
	assert.ErrorIs(t, err, want)
}
