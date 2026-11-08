// Package hub provides HTTP handlers for managing logistics hubs.
package hub

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *fiber.Ctx) error {
	hubs, err := h.service.List()
	if err != nil {
		return utils.Error(c, 500, "failed to fetch hubs")
	}
	return utils.Success(c, hubs)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	hub, err := h.service.GetByID(c.Params("id"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return utils.Error(c, 404, "hub not found")
		}
		return utils.Error(c, 500, "failed to fetch hub")
	}
	return utils.Success(c, hub)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var hub models.Hub
	if err := c.BodyParser(&hub); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	if err := h.service.Create(&hub); err != nil {
		return utils.Error(c, 500, "failed to create hub")
	}
	return utils.Success(c, hub)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var hub models.Hub
	if err := c.BodyParser(&hub); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	if err := h.service.Update(id, &hub); err != nil {
		if errors.Is(err, ErrNotFound) {
			return utils.Error(c, 404, "hub not found")
		}
		return utils.Error(c, 500, "failed to update hub")
	}
	return utils.Success(c, &hub)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Params("id")); err != nil {
		return utils.Error(c, 500, "failed to delete hub")
	}
	return utils.Success(c, fiber.Map{"message": "hub deleted"})
}
