package analytics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/pkg/utils"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

// Overview handles GET /api/analytics/overview. Public endpoint (no auth required).
func (h *Handler) Overview(c *fiber.Ctx) error {
	result, err := h.service.Overview(c.UserContext())
	if err != nil { return utils.Error(c, 500, "failed to fetch analytics") }
	return utils.Success(c, result)
}

// TimeSeries handles GET /api/analytics/timeseries. Public endpoint (no auth required).
func (h *Handler) TimeSeries(c *fiber.Ctx) error {
	result, err := h.service.TimeSeries(c.UserContext())
	if err != nil { return utils.Error(c, 500, "failed to fetch time series") }
	return utils.Success(c, result)
}
