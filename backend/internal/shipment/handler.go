package shipment

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

// HubRepository is the interface the shipment handler uses to look up a hub by its ID
// when updating status (e.g. to set CurrentCoords to the hub's location). It's defined here
// rather than in the hub package to avoid an import cycle — the handler depends on both.
type HubRepository interface {
	FindByID(id string) (*models.Hub, error)
}

// Handler holds the shipment read repository and workflow service.
type Handler struct {
	repo    HandlerRepository
	service *Service
}

// NewHandler wires shipment reads and the workflow service.
func NewHandler(repo HandlerRepository, service *Service) *Handler {
	return &Handler{repo: repo, service: service}
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

// List handles GET /api/shipments. Reads query params (page, limit, search, status, exclude_status),
// passes them to the repository, and returns paginated results. Public endpoint (no auth required).
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

// Create handles POST /api/shipments. Validates the JSON body, builds a Shipment model with
// auto-generated OrderID/TrackingNumber, sets origin/destination from customer/receiver addresses,
// sets status to "pending" with a 72h estimated delivery, creates an initial "Label Created"
// tracking event, and saves everything to Postgres. Requires staff authorization.
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	shipment, err := h.service.Create(req)
	if err != nil {
		if err.Error() == "customer name, province, and zipcode are required" ||
			err.Error() == "receiver name, province, and zipcode are required" ||
			err.Error() == "weight must be greater than 0" ||
			err.Error() == "items must be at least 1" ||
			err.Error() == "carrier is required" {
			return utils.Error(c, 400, err.Error())
		}
		return utils.Error(c, 500, "failed to create shipment")
	}
	return utils.Success(c, shipment)
}

// GetByID handles GET /api/shipments/:orderId. Looks up a single shipment by its ORD-xxxxx
// order ID and returns it directly. Returns 404 if not found. Public endpoint (no auth required).
func (h *Handler) GetByID(c *fiber.Ctx) error {
	orderID := c.Params("orderId")
	shipment, err := h.repo.FindByOrderID(orderID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		return utils.Error(c, 500, "failed to find shipment")
	}
	return utils.Success(c, shipment)
}

// statusToEvent translates a raw status string (e.g. "in_transit", "delivered") into a
// human-readable ShipmentEvent with a label, description, and location. Location logic:
//   - pending/picked_up → customer's address
//   - delivered → receiver's address
//   - departed/in_transit/out_for_delivery/delayed → hub location if hub is provided,
//     otherwise falls back to current coords, origin, or destination
//   - any unknown status → generic "Status updated." at destination
//
// This is a pure helper function — it does not save anything to the database.
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

// UpdateStatus handles PATCH /api/shipments/:orderId/status. Changes the shipment's status,
// optionally links a hub (updating CurrentCoords to the hub's lat/lng), saves the shipment,
// and creates a tracking event via statusToEvent with context-aware location.
// Requires staff authorization. Request body: { "status": "in_transit", "hubId": "HUB-003" }.
func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	var body struct {
		Status string `json:"status"`
		HubID  string `json:"hubId,omitempty"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	shipment, err := h.service.UpdateStatus(c.Params("orderId"), body.Status, body.HubID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		if errors.Is(err, ErrInvalidHub) || err.Error() == "status is required" {
			return utils.Error(c, 400, err.Error())
		}
		return utils.Error(c, 500, "failed to update shipment")
	}
	return utils.Success(c, shipment)
}

// UpdateRequest is the JSON body for updating an existing shipment. ALL fields are optional
// (pointer types) — only non-nil fields are applied. This lets callers send partial updates
// without accidentally zeroing out unset fields.
type UpdateRequest struct {
	Customer          *models.ContactInfo `json:"customer,omitempty"`
	Receiver          *models.ContactInfo `json:"receiver,omitempty"`
	Carrier           *string             `json:"carrier,omitempty"`
	Weight            *float64            `json:"weight,omitempty"`
	Items             *int                `json:"items,omitempty"`
	EstimatedDelivery *time.Time          `json:"estimatedDelivery,omitempty"`
}

// Update handles PUT /api/shipments/:orderId. Finds the existing shipment, applies only the
// fields the caller provided in the JSON body (pointer fields allow partial updates), recalculates
// Origin/Destination/CurrentCoords when customer or receiver changes, then saves. Requires staff authorization.
func (h *Handler) Update(c *fiber.Ctx) error {
	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	shipment, err := h.service.Update(c.Params("orderId"), req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		return utils.Error(c, 500, "failed to update shipment")
	}
	return utils.Success(c, shipment)
}

// Delete handles DELETE /api/shipments/:orderId. Deletes all tracking events first (to avoid
// foreign-key constraint errors), then deletes the shipment itself. Returns a success message.
// Requires staff authorization.
func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Params("orderId")); err != nil {
		if errors.Is(err, ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		return utils.Error(c, 500, "failed to delete shipment")
	}
	return utils.Success(c, fiber.Map{"message": "shipment deleted"})
}
