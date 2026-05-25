// Package hub provides HTTP handlers for managing logistics hubs.
package hub

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

// generateHubID creates the next hub ID like "HUB-007".
func generateHubID() string {
	var hubs []models.Hub
	database.DB.Select("id").Find(&hubs)
	maxNum := 0
	for _, h := range hubs {
		parts := strings.SplitN(h.ID, "-", 2)
		if len(parts) == 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil && n > maxNum {
				maxNum = n
			}
		}
	}
	return fmt.Sprintf("HUB-%03d", maxNum+1)
}

// List returns all hubs from the database.
func List(c *fiber.Ctx) error {
	var hubs []models.Hub
	database.DB.Find(&hubs)
	return utils.Success(c, hubs)
}

// GetByID fetches a single hub by its primary key ID.
func GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var hub models.Hub
	if result := database.DB.First(&hub, "id = ?", id); result.Error != nil {
		return utils.Error(c, 404, "hub not found")
	}
	return utils.Success(c, hub)
}

// Create adds a new hub to the database. Auto-generates a "HUB-xxx" ID if not provided.
func Create(c *fiber.Ctx) error {
	var hub models.Hub
	if err := c.BodyParser(&hub); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	if hub.ID == "" {
		hub.ID = generateHubID()
	}
	if result := database.DB.Create(&hub); result.Error != nil {
		return utils.Error(c, 500, "failed to create hub")
	}
	return utils.Success(c, hub)
}

// Update modifies an existing hub's fields.
func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var hub models.Hub
	if result := database.DB.First(&hub, "id = ?", id); result.Error != nil {
		return utils.Error(c, 404, "hub not found")
	}
	if err := c.BodyParser(&hub); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	hub.ID = id
	database.DB.Save(&hub)
	return utils.Success(c, hub)
}

// Delete removes a hub from the database.
func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if result := database.DB.Delete(&models.Hub{}, "id = ?", id); result.Error != nil {
		return utils.Error(c, 500, "failed to delete hub")
	}
	return utils.Success(c, fiber.Map{"message": "hub deleted"})
}
