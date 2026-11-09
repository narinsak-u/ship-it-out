// Package tracking provides the public-facing shipment tracking endpoint.
// It is the only handler accessible without authentication, enabling customers
// to look up their shipment status and event history via a tracking number.
package tracking

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Track(c *fiber.Ctx) error {
	result, err := h.service.Track(c.UserContext(), c.Params("trackingNumber"))
	if errors.Is(err, ErrShipmentNotFound) {
		return utils.Error(c, 404, "shipment not found")
	}
	if err != nil {
		return utils.Error(c, 500, "failed to fetch events")
	}

	return utils.Success(c, fiber.Map{
		"shipment": result.Shipment,
		"events":   result.Events,
	})
}
