package analytics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/data"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/pkg/utils"
)

type regionCount struct {
	Name  string `json:"name"`
	Total int64  `json:"total"`
}

type Handler struct {
	repo shipment.Repository
}

func NewHandler(repo shipment.Repository) *Handler {
	return &Handler{repo: repo}
}

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
