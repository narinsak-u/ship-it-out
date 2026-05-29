// Package tracking provides the public-facing shipment tracking endpoint.
// It is the only handler accessible without authentication, enabling customers
// to look up their shipment status and event history via a tracking number.
package tracking

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/pkg/utils"
)

type Handler struct {
	repo shipment.Repository
}

func NewHandler(repo shipment.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Track(c *fiber.Ctx) error {
	trackingNumber := c.Params("trackingNumber")

	shipment, err := h.repo.FindByTrackingNumber(trackingNumber)
	if err != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	events, err := h.repo.FindEventsByShipmentID(shipment.ID)
	if err != nil {
		return utils.Error(c, 500, "failed to fetch events")
	}

	return utils.Success(c, fiber.Map{
		"shipment": shipment,
		"events":   events,
	})
}
