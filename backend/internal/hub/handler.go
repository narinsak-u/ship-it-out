// Package hub provides HTTP handlers for managing logistics hubs.
package hub

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// List returns all hubs from the database.
func (h *Handler) List(c *fiber.Ctx) error {
	hubs, err := h.repo.FindAll()
	if err != nil {
		return utils.Error(c, 500, "failed to fetch hubs")
	}
	return utils.Success(c, hubs)
}

// GetByID fetches a single hub by its primary key ID.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	hub, err := h.repo.FindByID(id)
	if err != nil {
		return utils.Error(c, 404, "hub not found")
	}
	return utils.Success(c, hub)
}

// Create adds a new hub to the database. Auto-generates a "HUB-xxx" ID if not provided.
func (h *Handler) Create(c *fiber.Ctx) error {
	var hub models.Hub
	if err := c.BodyParser(&hub); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	if err := h.repo.Create(&hub); err != nil {
		return utils.Error(c, 500, "failed to create hub")
	}
	return utils.Success(c, hub)
}

// Update modifies an existing hub's fields.
func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := h.repo.FindByID(id); err != nil {
		return utils.Error(c, 404, "hub not found")
	}
	var hub models.Hub
	if err := c.BodyParser(&hub); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	hub.ID = id
	if err := h.repo.Save(&hub); err != nil {
		return utils.Error(c, 500, "failed to update hub")
	}
	return utils.Success(c, hub)
}

// Delete removes a hub from the database.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.Delete(id); err != nil {
		return utils.Error(c, 500, "failed to delete hub")
	}
	return utils.Success(c, fiber.Map{"message": "hub deleted"})
}
