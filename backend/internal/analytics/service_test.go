package analytics

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	total      int64
	active     int64
	statuses   []StatusCount
	provinces  []ProvinceCount
	months     []MonthCount
	days       []DayCount
	err        error
}

func (r fakeRepository) Count() (int64, error) { return r.total, r.err }
func (r fakeRepository) CountActive() (int64, error) { return r.active, r.err }
func (r fakeRepository) CountByStatus() ([]StatusCount, error) { return r.statuses, r.err }
func (r fakeRepository) CountByProvince() ([]ProvinceCount, error) { return r.provinces, r.err }
func (r fakeRepository) CountByMonth() ([]MonthCount, error) { return r.months, r.err }
func (r fakeRepository) CountByDayOfWeek() ([]DayCount, error) { return r.days, r.err }

func TestServiceOverviewAggregatesDeliveredAndRegions(t *testing.T) {
	service := NewService(fakeRepository{
		total: 10, active: 7,
		statuses: []StatusCount{{Status: "delivered", Count: 2}, {Status: "pending", Count: 8}},
		provinces: []ProvinceCount{{Province: "ชลบุรี", Total: 5}, {Province: "กรุงเทพมหานคร", Total: 5}, {Province: "unknown", Total: 7}},
	})

	got, err := service.Overview()
	require.NoError(t, err)
	assert.Equal(t, int64(10), got.Total)
	assert.Equal(t, int64(7), got.Active)
	assert.Equal(t, int64(2), got.Delivered)
	assert.Equal(t, []StatusCount{{Status: "delivered", Count: 2}, {Status: "pending", Count: 8}}, got.ByStatus)
	regions := map[string]int64{}
	for _, region := range got.ByRegion {
		regions[region.Name] = region.Total
	}
	assert.Equal(t, int64(5), regions["Central"])
	assert.Equal(t, int64(5), regions["East"])
	assert.Len(t, got.ByRegion, 6)
}

func TestServiceTimeSeriesReturnsRepositoryResults(t *testing.T) {
	months := []MonthCount{{Month: "2026-01", Count: 3}}
	days := []DayCount{{Day: "Monday", Count: 2}}
	got, err := NewService(fakeRepository{months: months, days: days}).TimeSeries()
	require.NoError(t, err)
	assert.Equal(t, months, got.ByMonth)
	assert.Equal(t, days, got.ByDayOfWeek)
}

func TestServicePropagatesQueryErrors(t *testing.T) {
	want := errors.New("database unavailable")
	_, err := NewService(fakeRepository{err: want}).Overview()
	assert.ErrorIs(t, err, want)
	_, err = NewService(fakeRepository{err: want}).TimeSeries()
	assert.ErrorIs(t, err, want)
}
