// Package presentation provides HTTP handlers for managing logistics hubs.
package presentation

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/hub/application"
	"github.com/narinsak-u/backend/pkg/utils"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *fiber.Ctx) error {
	hubs, err := h.service.List(c.UserContext())
	if err != nil {
		return utils.Error(c, 500, "failed to fetch hubs")
	}
	return utils.Success(c, hubs)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	hub, err := h.service.GetByID(c.UserContext(), c.Params("id"))
	if err != nil {
		if errors.Is(err, application.ErrNotFound) {
			return utils.Error(c, 404, "hub not found")
		}
		return utils.Error(c, 500, "failed to fetch hub")
	}
	return utils.Success(c, hub)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var request HubRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	hub := request.hub()
	if err := h.service.Create(c.UserContext(), hub); err != nil {
		return utils.Error(c, 500, "failed to create hub")
	}
	return utils.Success(c, hub)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := h.service.GetByID(c.UserContext(), id); err != nil {
		if errors.Is(err, application.ErrNotFound) {
			return utils.Error(c, 404, "hub not found")
		}
		return utils.Error(c, 500, "failed to update hub")
	}
	var request HubRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	hub := request.hub()
	if err := h.service.Update(c.UserContext(), id, hub); err != nil {
		if errors.Is(err, application.ErrNotFound) {
			return utils.Error(c, 404, "hub not found")
		}
		return utils.Error(c, 500, "failed to update hub")
	}
	return utils.Success(c, hub)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.UserContext(), c.Params("id")); err != nil {
		return utils.Error(c, 500, "failed to delete hub")
	}
	return utils.Success(c, fiber.Map{"message": "hub deleted"})
}
