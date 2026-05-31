package analytics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/data"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/pkg/utils"
)

// regionCount is an internal helper for aggregating shipment counts by Thai geographic region.
// Not exported — results are serialized directly in the Overview response.
type regionCount struct {
	Name  string `json:"name"`
	Total int64  `json:"total"`
}

// Handler exposes dashboard analytics as HTTP handlers. It reuses the shipment package's
// Repository interface (no separate analytics repository needed). Create via NewHandler.
type Handler struct {
	repo shipment.Repository
}

// NewHandler creates an analytics Handler backed by the shipment Repository. The same
// shipment repository serves both CRUD and analytics queries.
func NewHandler(repo shipment.Repository) *Handler {
	return &Handler{repo: repo}
}

// Overview handles GET /api/analytics/overview. Returns aggregate stats for the dashboard:
// total shipments, active (not delivered/returned), delivered count, breakdown by status,
// and breakdown by Thai geographic region (derived from receiver province).
// Public endpoint (no auth required in the current configuration).
func (h *Handler) Overview(c *fiber.Ctx) error {
	total, err := h.repo.Count()
	if err != nil {
		return utils.Error(c, 500, "failed to fetch analytics")
	}

	active, err := h.repo.CountActive()
	if err != nil {
		return utils.Error(c, 500, "failed to fetch analytics")
	}

	byStatus, err := h.repo.CountByStatus()
	if err != nil {
		return utils.Error(c, 500, "failed to fetch analytics")
	}

	var delivered int64
	for _, s := range byStatus {
		if s.Status == "delivered" {
			delivered = s.Count
			break
		}
	}

	byProvince, err := h.repo.CountByProvince()
	if err != nil {
		return utils.Error(c, 500, "failed to fetch analytics")
	}

	regionMap := map[string]*regionCount{
		"Central":    {Name: "Central"},
		"East":       {Name: "East"},
		"North":      {Name: "North"},
		"West":       {Name: "West"},
		"North-east": {Name: "North-east"},
		"South":      {Name: "South"},
	}
	for _, p := range byProvince {
		region := data.ThailandProvinceRegion[p.Province]
		if region == "" {
			continue
		}
		r := regionMap[region]
		r.Total += p.Total
	}

	byRegion := make([]regionCount, 0, len(regionMap))
	for _, r := range regionMap {
		byRegion = append(byRegion, *r)
	}

	return utils.Success(c, fiber.Map{
		"total":     total,
		"active":    active,
		"delivered": delivered,
		"by_status": byStatus,
		"by_region": byRegion,
	})
}

// TimeSeries handles GET /api/analytics/timeseries. Returns shipment creation trends
// grouped by month (YYYY-MM) and by day-of-week (Monday..Sunday), useful for line/bar
// charts on the dashboard. Public endpoint (no auth required).
func (h *Handler) TimeSeries(c *fiber.Ctx) error {
	byMonth, err := h.repo.CountByMonth()
	if err != nil {
		return utils.Error(c, 500, "failed to fetch time series")
	}

	byDay, err := h.repo.CountByDayOfWeek()
	if err != nil {
		return utils.Error(c, 500, "failed to fetch time series")
	}

	return utils.Success(c, fiber.Map{
		"by_month":       byMonth,
		"by_day_of_week": byDay,
	})
}
