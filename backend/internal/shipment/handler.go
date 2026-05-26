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

// List returns shipments from the database with optional pagination and filtering.
// Query params: page (default 1), limit (default 10, use -1 for all), search, status, exclude_status.
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	status := c.Query("status", "")
	excludeStatus := c.Query("exclude_status", "")

	if page < 1 {
		page = 1
	}

	var total int64
	query := database.DB.Model(&models.Shipment{})

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if excludeStatus != "" {
		query = query.Where("status != ?", excludeStatus)
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"order_id ILIKE ? OR tracking_number ILIKE ? OR customer_name ILIKE ? OR destination ILIKE ?",
			like, like, like, like,
		)
	}

	query.Count(&total)

	var shipments []models.Shipment
	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}
	query.Order("created_at DESC").Find(&shipments)

	return utils.SuccessWithPagination(c, shipments, page, limit, int(total))
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
		ShipmentID:  shipment.ID,
		Status:      "Label Created",
		Description: "Awaiting pickup.",
		Location: models.Location{
			Name: composeAddress(req.Customer),
			Lat:  req.Customer.Coords.Lat,
			Lng:  req.Customer.Coords.Lng,
		},
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

// statusToEvent builds the tracking event details (status label, description, location)
// for a given shipment status. Some statuses require a hub reference for location.
func statusToEvent(shipment models.Shipment, hub *models.Hub, targetStatus string) models.ShipmentEvent {
	var eventStatus, description string
	var location models.Location

	switch targetStatus {
	case "pending":
		eventStatus = "Label Created"
		description = "Awaiting pickup."
		location = models.Location{
			Name: composeAddress(shipment.Customer),
			Lat:  shipment.CustomerLat,
			Lng:  shipment.CustomerLng,
		}
	case "picked_up":
		eventStatus = "Picked Up"
		description = "Parcel collected from sender."
		location = models.Location{
			Name: composeAddress(shipment.Customer),
			Lat:  shipment.CustomerLat,
			Lng:  shipment.CustomerLng,
		}
	case "departed":
		eventStatus = "Departed"
		description = "In transit to hub."
		if hub != nil {
			location = models.Location{Name: hub.Name + ", " + hub.Address, Lat: hub.Coords.Lat, Lng: hub.Coords.Lng}
		} else {
			location = models.Location{Name: shipment.Origin, Lat: shipment.CurrentLat, Lng: shipment.CurrentLng}
		}
	case "in_transit":
		eventStatus = "In Transit"
		description = "Transit to next hub."
		if hub != nil {
			location = models.Location{Name: hub.Name + ", " + hub.Address, Lat: hub.Coords.Lat, Lng: hub.Coords.Lng}
		} else {
			location = models.Location{Name: shipment.Destination, Lat: shipment.CurrentLat, Lng: shipment.CurrentLng}
		}
	case "out_for_delivery":
		eventStatus = "Out for Delivery"
		description = "Out for delivery."
		if hub != nil {
			location = models.Location{Name: hub.Name + ", " + hub.Address, Lat: hub.Coords.Lat, Lng: hub.Coords.Lng}
		} else {
			location = models.Location{Name: shipment.Destination, Lat: shipment.ReceiverLat, Lng: shipment.ReceiverLng}
		}
	case "delivered":
		eventStatus = "Delivered"
		description = "Delivered to recipient."
		location = models.Location{
			Name: composeAddress(shipment.Receiver),
			Lat:  shipment.ReceiverLat,
			Lng:  shipment.ReceiverLng,
		}
	case "delayed":
		eventStatus = "Delayed"
		description = "Unexpected issue encountered."
		if hub != nil {
			location = models.Location{Name: hub.Name + ", " + hub.Address, Lat: hub.Coords.Lat, Lng: hub.Coords.Lng}
		} else {
			location = models.Location{Name: shipment.Destination, Lat: shipment.CurrentLat, Lng: shipment.CurrentLng}
		}
	default:
		eventStatus = targetStatus
		description = "Status updated."
		location = models.Location{Name: shipment.Destination, Lat: shipment.CurrentLat, Lng: shipment.CurrentLng}
	}

	return models.ShipmentEvent{
		Status:      eventStatus,
		Description: description,
		Location:    location,
	}
}

// UpdateStatus changes the status of a shipment and records a tracking event
// with context-aware status label, description, and location.
func UpdateStatus(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	var body struct {
		Status string `json:"status"`
		HubID  string `json:"hubId,omitempty"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var shipment models.Shipment
	if result := database.DB.Where("order_id = ?", orderID).First(&shipment); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	shipment.Status = body.Status
	if body.HubID != "" {
		shipment.HubID = body.HubID
	}

	// Look up hub if provided (for statuses where location = hub address)
	var hub *models.Hub
	if body.HubID != "" {
		var h models.Hub
		if err := database.DB.Where("id = ?", body.HubID).First(&h); err.Error == nil {
			hub = &h
			shipment.CurrentCoords.Lat = h.Lat
			shipment.CurrentCoords.Lng = h.Lng
		}
	}

	database.DB.Save(&shipment)

	event := statusToEvent(shipment, hub, body.Status)
	event.ShipmentID = shipment.ID
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
