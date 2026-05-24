// Package shipment provides HTTP handlers for shipment CRUD operations.
package shipment

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

// CreateRequest is the JSON body for creating a new shipment.
// Uses the same nested ContactInfo structure as the frontend Order type.
type CreateRequest struct {
	Customer models.ContactInfo `json:"customer"`
	Receiver models.ContactInfo `json:"receiver"`
	Carrier  string             `json:"carrier"`
	Weight   string             `json:"weight"`
	Items    int                `json:"items"`
}

// generateTrackingNumber creates a unique identifier for each shipment.
// Format: "TH" + current year + 5-digit millisecond-based number
// Example: "TH202596374"
func generateTrackingNumber() string {
	return fmt.Sprintf("TH%d%05d", time.Now().Year(), time.Now().UnixMilli()%100000)
}

// composeAddress builds a display string from a ContactInfo's address fields.
func composeAddress(c models.ContactInfo) string {
	return fmt.Sprintf("%s, %s, %s", c.SubDistrict, c.District, c.Province)
}

// List returns all shipments from the database.
func List(c *fiber.Ctx) error {
	var shipments []models.Shipment
	database.DB.Find(&shipments)
	return utils.Success(c, shipments)
}

// Create adds a new shipment to the database.
func Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	shipment := models.Shipment{
		TrackingNumber:    generateTrackingNumber(),
		Customer:          req.Customer,
		Receiver:          req.Receiver,
		Origin:            composeAddress(req.Customer),
		Destination:       composeAddress(req.Receiver),
		CurrentCoords:     req.Customer.Coords,
		Status:            "pending",
		Carrier:           req.Carrier,
		Weight:            req.Weight,
		Items:             req.Items,
		EstimatedDelivery: time.Now().Add(72 * time.Hour),
		Progress:          0,
	}

	if result := database.DB.Create(&shipment); result.Error != nil {
		return utils.Error(c, 500, "failed to create shipment")
	}

	return utils.Success(c, shipment)
}

// GetByID fetches a single shipment by its primary key ID.
func GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, 400, "invalid id")
	}

	var shipment models.Shipment
	if result := database.DB.First(&shipment, id); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	return utils.Success(c, shipment)
}

// UpdateStatus changes the status of a shipment and records the change
// as a ShipmentEvent for the tracking timeline.
func UpdateStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, 400, "invalid id")
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var shipment models.Shipment
	if result := database.DB.First(&shipment, id); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	shipment.Status = body.Status
	database.DB.Save(&shipment)

	event := models.ShipmentEvent{
		ShipmentID: shipment.ID,
		Status:     body.Status,
		Location: models.Location{
			Name: shipment.Destination,
			Lat:  shipment.CurrentCoords.Lat,
			Lng:  shipment.CurrentCoords.Lng,
		},
		Description: fmt.Sprintf("Status updated to %s", body.Status),
	}
	database.DB.Create(&event)

	return utils.Success(c, shipment)
}
