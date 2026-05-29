package shipment

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

type HubRepository interface {
	FindByID(id string) (*models.Hub, error)
}

type Handler struct {
	repo    Repository
	hubRepo HubRepository
}

func NewHandler(repo Repository, hubRepo HubRepository) *Handler {
	return &Handler{repo: repo, hubRepo: hubRepo}
}

// CreateRequest is the JSON body for creating a new shipment.
// Uses the same nested ContactInfo structure as the frontend Order type.
type CreateRequest struct {
	Customer models.ContactInfo `json:"customer"`
	Receiver models.ContactInfo `json:"receiver"`
	Carrier  string             `json:"carrier"`
	Weight   float64            `json:"weight"`
	Items    int                `json:"items"`
}

// List returns shipments from the database with optional pagination and filtering.
// Query params: page (default 1), limit (default 10), search, status, exclude_status.
func (h *Handler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	filter := ShipmentFilter{
		Page:          page,
		Limit:         limit,
		Search:        c.Query("search", ""),
		Status:        c.Query("status", ""),
		ExcludeStatus: c.Query("exclude_status", ""),
	}
	shipments, total, err := h.repo.List(filter)
	if err != nil {
		return utils.Error(c, 500, "failed to list shipments")
	}
	return utils.SuccessWithPagination(c, shipments, filter.Page, filter.Limit, int(total))
}

// Create adds a new shipment to the database.
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	if req.Customer.Name == "" || req.Customer.Province == "" || req.Customer.Zipcode == "" {
		return utils.Error(c, 400, "customer name, province, and zipcode are required")
	}
	if req.Receiver.Name == "" || req.Receiver.Province == "" || req.Receiver.Zipcode == "" {
		return utils.Error(c, 400, "receiver name, province, and zipcode are required")
	}
	if req.Weight <= 0 {
		return utils.Error(c, 400, "weight must be greater than 0")
	}
	if req.Items < 1 {
		return utils.Error(c, 400, "items must be at least 1")
	}
	if req.Carrier == "" {
		return utils.Error(c, 400, "carrier is required")
	}

	shipment := models.Shipment{
		Customer:          req.Customer,
		Receiver:          req.Receiver,
		Origin:            utils.ComposeAddress(req.Customer),
		Destination:       utils.ComposeAddress(req.Receiver),
		CurrentCoords:     req.Customer.Coords,
		Status:            "pending",
		Carrier:           req.Carrier,
		Weight:            req.Weight,
		Items:             req.Items,
		EstimatedDelivery: time.Now().Add(72 * time.Hour),
		Progress:          0,
	}

	if err := h.repo.Create(&shipment); err != nil {
		return utils.Error(c, 500, "failed to create shipment")
	}

	// Create the initial tracking event for the new shipment
	event := models.ShipmentEvent{
		ShipmentID:  shipment.ID,
		Status:      "Label Created",
		Description: "Awaiting pickup.",
		Location: models.Location{
			Name: utils.ComposeAddress(req.Customer),
			Lat:  req.Customer.Coords.Lat,
			Lng:  req.Customer.Coords.Lng,
		},
	}
	h.repo.CreateEvent(&event)

	return utils.Success(c, shipment)
}

// GetByID fetches a single shipment by its order ID.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	orderID := c.Params("orderId")
	shipment, err := h.repo.FindByOrderID(orderID)
	if err != nil {
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
			Name: utils.ComposeAddress(shipment.Customer),
			Lat:  shipment.CustomerLat,
			Lng:  shipment.CustomerLng,
		}
	case "picked_up":
		eventStatus = "Picked Up"
		description = "Parcel collected from sender."
		location = models.Location{
			Name: utils.ComposeAddress(shipment.Customer),
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
			Name: utils.ComposeAddress(shipment.Receiver),
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
func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	var body struct {
		Status string `json:"status"`
		HubID  string `json:"hubId,omitempty"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	if body.Status == "" {
		return utils.Error(c, 400, "status is required")
	}

	shipment, err := h.repo.FindByOrderID(orderID)
	if err != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	shipment.Status = body.Status

	var hub *models.Hub
	if body.HubID != "" {
		hubRecord, err := h.hubRepo.FindByID(body.HubID)
		if err != nil {
			return utils.Error(c, 400, "invalid hub ID")
		}
		hub = hubRecord
		shipment.HubID = body.HubID
		shipment.CurrentCoords.Lat = hubRecord.Lat
		shipment.CurrentCoords.Lng = hubRecord.Lng
	}

	if err := h.repo.Save(shipment); err != nil {
		return utils.Error(c, 500, "failed to update shipment")
	}

	event := statusToEvent(*shipment, hub, body.Status)
	event.ShipmentID = shipment.ID
	if err := h.repo.CreateEvent(&event); err != nil {
		return utils.Error(c, 500, "failed to create tracking event")
	}

	return utils.Success(c, shipment)
}

// UpdateRequest is the JSON body for updating an existing shipment.
type UpdateRequest struct {
	Customer          *models.ContactInfo `json:"customer,omitempty"`
	Receiver          *models.ContactInfo `json:"receiver,omitempty"`
	Carrier           *string             `json:"carrier,omitempty"`
	Weight            *float64            `json:"weight,omitempty"`
	Items             *int                `json:"items,omitempty"`
	EstimatedDelivery *time.Time          `json:"estimatedDelivery,omitempty"`
}

// Update modifies an existing shipment's fields.
func (h *Handler) Update(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	shipment, err := h.repo.FindByOrderID(orderID)
	if err != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	if req.Customer != nil {
		shipment.Customer = *req.Customer
		shipment.CurrentCoords = req.Customer.Coords
		shipment.Origin = utils.ComposeAddress(*req.Customer)
	}
	if req.Receiver != nil {
		shipment.Receiver = *req.Receiver
		shipment.Destination = utils.ComposeAddress(*req.Receiver)
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

	h.repo.Save(shipment)
	return utils.Success(c, shipment)
}

// Delete removes a shipment and its events from the database.
func (h *Handler) Delete(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	shipment, err := h.repo.FindByOrderID(orderID)
	if err != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	if err := h.repo.DeleteShipmentEvents(shipment.ID); err != nil {
		return utils.Error(c, 500, "failed to delete shipment events")
	}
	if err := h.repo.Delete(shipment); err != nil {
		return utils.Error(c, 500, "failed to delete shipment")
	}

	return utils.Success(c, fiber.Map{"message": "shipment deleted"})
}
