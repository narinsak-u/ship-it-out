// Package hub provides HTTP handlers for managing logistics hubs.
package hub

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

// Handler routes HTTP requests to the hub repository. Public routes (List, GetByID) need no auth;
// write routes (Create, Update, Delete) are protected by AuthRequired middleware in main.go.
// Create via NewHandler — never instantiate directly.
type Handler struct {
	repo Repository
}

// NewHandler creates a hub Handler backed by the given Repository (typically a GormRepository).
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// List handles GET /api/hubs. Returns every hub in the database. Public endpoint (no auth).
func (h *Handler) List(c *fiber.Ctx) error {
	hubs, err := h.repo.FindAll()
	if err != nil {
		return utils.Error(c, 500, "failed to fetch hubs")
	}
	return utils.Success(c, hubs)
}

// GetByID handles GET /api/hubs/:id. Looks up a single hub by its string ID (e.g. "HUB-001").
// Returns 404 if not found. Public endpoint (no auth).
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	hub, err := h.repo.FindByID(id)
	if err != nil {
		return utils.Error(c, 404, "hub not found")
	}
	return utils.Success(c, hub)
}

// Create handles POST /api/hubs. Parses the JSON body into a Hub model and inserts it.
// If the JSON body does not include an "id" field, the repository auto-generates "HUB-NNN".
// Requires JWT auth. Request body: { name, carrierId, address, coords, capacity, ... }.
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

// Update handles PUT /api/hubs/:id. Finds the existing hub, overwrites its fields with the
// request body, re-sets the ID from the URL param (BodyParser clears the PK field), then saves.
// Requires JWT auth.
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

// Delete handles DELETE /api/hubs/:id. Removes the hub row from Postgres by its string ID.
// Returns a confirmation message. Requires JWT auth.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.Delete(id); err != nil {
		return utils.Error(c, 500, "failed to delete hub")
	}
	return utils.Success(c, fiber.Map{"message": "hub deleted"})
}
