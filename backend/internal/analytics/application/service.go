package application

import (
	"context"

	"github.com/narinsak-u/backend/internal/data"
)

type StatusCount struct {
	Status string
	Count  int64
}

type ProvinceCount struct {
	Province string
	Total    int64
}

type MonthCount struct {
	Month string
	Count int64
}

type DayCount struct {
	Day   string
	Count int64
}

type OverviewReadPort interface {
	Count(context.Context) (int64, error)
	CountActive(context.Context) (int64, error)
	CountByStatus(context.Context) ([]StatusCount, error)
	CountByProvince(context.Context) ([]ProvinceCount, error)
}

type TimeSeriesReadPort interface {
	CountByMonth(context.Context) ([]MonthCount, error)
	CountByDayOfWeek(context.Context) ([]DayCount, error)
}

type RegionCount struct {
	Name  string
	Total int64
}

type Overview struct {
	Total     int64
	Active    int64
	Delivered int64
	ByStatus  []StatusCount
	ByRegion  []RegionCount
}

type TimeSeries struct {
	ByMonth     []MonthCount
	ByDayOfWeek []DayCount
}

type Service struct {
	overview   OverviewReadPort
	timeSeries TimeSeriesReadPort
}

func NewService(overview OverviewReadPort, timeSeries TimeSeriesReadPort) *Service {
	return &Service{overview: overview, timeSeries: timeSeries}
}

func (s *Service) Overview(ctx context.Context) (Overview, error) {
	var result Overview
	var err error
	if result.Total, err = s.overview.Count(ctx); err != nil {
		return Overview{}, err
	}
	if result.Active, err = s.overview.CountActive(ctx); err != nil {
		return Overview{}, err
	}
	if result.ByStatus, err = s.overview.CountByStatus(ctx); err != nil {
		return Overview{}, err
	}
	for _, status := range result.ByStatus {
		if status.Status == "delivered" {
			result.Delivered = status.Count
			break
		}
	}
	provinces, err := s.overview.CountByProvince(ctx)
	if err != nil {
		return Overview{}, err
	}
	regions := map[string]int64{"Central": 0, "East": 0, "North": 0, "West": 0, "North-east": 0, "South": 0}
	for _, province := range provinces {
		if region := data.ThailandProvinceRegion[province.Province]; region != "" {
			regions[region] += province.Total
		}
	}
	result.ByRegion = make([]RegionCount, 0, len(regions))
	for name, total := range regions {
		result.ByRegion = append(result.ByRegion, RegionCount{Name: name, Total: total})
	}
	return result, nil
}

func (s *Service) TimeSeries(ctx context.Context) (TimeSeries, error) {
	months, err := s.timeSeries.CountByMonth(ctx)
	if err != nil {
		return TimeSeries{}, err
	}
	days, err := s.timeSeries.CountByDayOfWeek(ctx)
	if err != nil {
		return TimeSeries{}, err
	}
	return TimeSeries{ByMonth: months, ByDayOfWeek: days}, nil
}
