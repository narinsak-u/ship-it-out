// Package shipment provides HTTP handlers for shipment CRUD operations.
package shipment

import (
	"fmt"
	"strconv"
	"strings"
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

// generateOrderID creates a human-readable order ID like "ORD-10251".
// Scans existing order IDs to find the highest numeric suffix and increments it.
func generateOrderID() string {
	var shipments []models.Shipment
	database.DB.Select("order_id").Find(&shipments)
	maxNum := 10245
	for _, s := range shipments {
		parts := strings.SplitN(s.OrderID, "-", 2)
		if len(parts) == 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil && n > maxNum {
				maxNum = n
			}
		}
	}
	return fmt.Sprintf("ORD-%d", maxNum+1)
}

// Create adds a new shipment to the database.
func Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	shipment := models.Shipment{
		OrderID:           generateOrderID(),
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

	// Create the initial tracking event for the new shipment
	event := models.ShipmentEvent{
		ShipmentID: shipment.ID,
		Status:     "Label created",
		Location: models.Location{
			Name: composeAddress(req.Customer),
			Lat:  req.Customer.Coords.Lat,
			Lng:  req.Customer.Coords.Lng,
		},
		Description: "Awaiting pickup.",
	}
	database.DB.Create(&event)

	return utils.Success(c, shipment)
}

// GetByID fetches a single shipment by its order ID.
func GetByID(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	var shipment models.Shipment
	if result := database.DB.Where("order_id = ?", orderID).First(&shipment); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	return utils.Success(c, shipment)
}

// UpdateStatus changes the status of a shipment and records the change
// as a ShipmentEvent for the tracking timeline.
func UpdateStatus(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var shipment models.Shipment
	if result := database.DB.Where("order_id = ?", orderID).First(&shipment); result.Error != nil {
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

// UpdateRequest is the JSON body for updating an existing shipment.
type UpdateRequest struct {
	Customer          *models.ContactInfo `json:"customer,omitempty"`
	Receiver          *models.ContactInfo `json:"receiver,omitempty"`
	Carrier           *string             `json:"carrier,omitempty"`
	Weight            *string             `json:"weight,omitempty"`
	Items             *int                `json:"items,omitempty"`
	EstimatedDelivery *time.Time          `json:"estimatedDelivery,omitempty"`
}

// Update modifies an existing shipment's fields.
func Update(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var shipment models.Shipment
	if result := database.DB.Where("order_id = ?", orderID).First(&shipment); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	if req.Customer != nil {
		shipment.Customer = *req.Customer
		shipment.CurrentCoords = req.Customer.Coords
		shipment.Origin = composeAddress(*req.Customer)
	}
	if req.Receiver != nil {
		shipment.Receiver = *req.Receiver
		shipment.Destination = composeAddress(*req.Receiver)
	}
	if req.Carrier != nil {
		shipment.Carrier = *req.Carrier
	}
	if req.Weight != nil {
		shipment.Weight = *req.Weight
	}
	if req.Items != nil {
		shipment.Items = *req.Items
	}
	if req.EstimatedDelivery != nil {
		shipment.EstimatedDelivery = *req.EstimatedDelivery
	}

	database.DB.Save(&shipment)
	return utils.Success(c, shipment)
}

// Delete removes a shipment and its events from the database.
func Delete(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	var shipment models.Shipment
	if result := database.DB.Where("order_id = ?", orderID).First(&shipment); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	database.DB.Where("shipment_id = ?", shipment.ID).Delete(&models.ShipmentEvent{})
	database.DB.Delete(&shipment)

	return utils.Success(c, fiber.Map{"message": "shipment deleted"})
}
