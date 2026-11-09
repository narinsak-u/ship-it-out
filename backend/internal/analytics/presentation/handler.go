package presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/analytics/application"
	"github.com/narinsak-u/backend/pkg/utils"
)

type statusCountResponse struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type regionCountResponse struct {
	Name  string `json:"name"`
	Total int64  `json:"total"`
}

type overviewResponse struct {
	Total     int64                 `json:"total"`
	Active    int64                 `json:"active"`
	Delivered int64                 `json:"delivered"`
	ByStatus  []statusCountResponse `json:"by_status"`
	ByRegion  []regionCountResponse `json:"by_region"`
}

type monthCountResponse struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type dayCountResponse struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}

type timeSeriesResponse struct {
	ByMonth     []monthCountResponse `json:"by_month"`
	ByDayOfWeek []dayCountResponse   `json:"by_day_of_week"`
}

func presentOverview(result application.Overview) overviewResponse {
	response := overviewResponse{
		Total: result.Total, Active: result.Active, Delivered: result.Delivered,
	}
	if result.ByStatus != nil {
		response.ByStatus = make([]statusCountResponse, len(result.ByStatus))
		for i, status := range result.ByStatus {
			response.ByStatus[i] = statusCountResponse{Status: status.Status, Count: status.Count}
		}
	}
	if result.ByRegion != nil {
		response.ByRegion = make([]regionCountResponse, len(result.ByRegion))
		for i, region := range result.ByRegion {
			response.ByRegion[i] = regionCountResponse{Name: region.Name, Total: region.Total}
		}
	}
	return response
}

func presentTimeSeries(result application.TimeSeries) timeSeriesResponse {
	response := timeSeriesResponse{}
	if result.ByMonth != nil {
		response.ByMonth = make([]monthCountResponse, len(result.ByMonth))
		for i, month := range result.ByMonth {
			response.ByMonth[i] = monthCountResponse{Month: month.Month, Count: month.Count}
		}
	}
	if result.ByDayOfWeek != nil {
		response.ByDayOfWeek = make([]dayCountResponse, len(result.ByDayOfWeek))
		for i, day := range result.ByDayOfWeek {
			response.ByDayOfWeek[i] = dayCountResponse{Day: day.Day, Count: day.Count}
		}
	}
	return response
}

type Handler struct{ service *application.Service }

func NewHandler(service *application.Service) *Handler { return &Handler{service: service} }

// Overview handles GET /api/analytics/overview. Public endpoint (no auth required).
func (h *Handler) Overview(c *fiber.Ctx) error {
	result, err := h.service.Overview(c.UserContext())
	if err != nil {
		return utils.Error(c, 500, "failed to fetch analytics")
	}
	return utils.Success(c, presentOverview(result))
}

// TimeSeries handles GET /api/analytics/timeseries. Public endpoint (no auth required).
func (h *Handler) TimeSeries(c *fiber.Ctx) error {
	result, err := h.service.TimeSeries(c.UserContext())
	if err != nil {
		return utils.Error(c, 500, "failed to fetch time series")
	}
	return utils.Success(c, presentTimeSeries(result))
}
