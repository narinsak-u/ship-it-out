package tracking

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

func Track(c *fiber.Ctx) error {
	trackingNumber := c.Params("trackingNumber")

	var shipment models.Shipment
	if result := database.DB.Where("tracking_number = ?", trackingNumber).First(&shipment); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	var events []models.ShipmentEvent
	database.DB.Where("shipment_id = ?", shipment.ID).Order("created_at asc").Find(&events)

	return utils.Success(c, fiber.Map{
		"shipment": shipment,
		"events":   events,
	})
}
