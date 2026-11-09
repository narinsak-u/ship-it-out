package application

import (
	"context"

	"github.com/narinsak-u/backend/internal/data"
)

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type ProvinceCount struct {
	Province string
	Total    int64
}

type MonthCount struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type DayCount struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
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
	Name  string `json:"name"`
	Total int64  `json:"total"`
}

type Overview struct {
	Total     int64         `json:"total"`
	Active    int64         `json:"active"`
	Delivered int64         `json:"delivered"`
	ByStatus  []StatusCount `json:"by_status"`
	ByRegion  []RegionCount `json:"by_region"`
}

type TimeSeries struct {
	ByMonth     []MonthCount `json:"by_month"`
	ByDayOfWeek []DayCount   `json:"by_day_of_week"`
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
