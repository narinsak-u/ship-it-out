package shipment

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

// Handler exposes shipment operations through the workflow service.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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
	shipments, total, err := h.service.List(c.UserContext(), filter)
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
	shipment, err := h.service.Create(c.UserContext(), req)
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
	shipment, err := h.service.GetByID(c.UserContext(), orderID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		return utils.Error(c, 500, "failed to find shipment")
	}
	return utils.Success(c, shipment)
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
	shipment, err := h.service.UpdateStatus(c.UserContext(), c.Params("orderId"), body.Status, body.HubID)
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
	shipment, err := h.service.Update(c.UserContext(), c.Params("orderId"), req)
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
	if err := h.service.Delete(c.UserContext(), c.Params("orderId")); err != nil {
		if errors.Is(err, ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		return utils.Error(c, 500, "failed to delete shipment")
	}
	return utils.Success(c, fiber.Map{"message": "shipment deleted"})
}
