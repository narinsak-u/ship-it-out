package presentation

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/shipment/application"
	"github.com/narinsak-u/backend/internal/shipment/domain"
	utils "github.com/narinsak-u/backend/internal/platform/http"
)

type Handler struct{ service *application.Service }

func NewHandler(service *application.Service) *Handler { return &Handler{service: service} }

type ContactDTO struct {
	Name        string   `json:"name"`
	Zipcode     string   `json:"zipcode"`
	SubDistrict string   `json:"subDistrict"`
	District    string   `json:"district"`
	Province    string   `json:"province"`
	Coords      PointDTO `json:"coords"`
}

type PointDTO struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type ShipmentDTO struct {
	ID                string     `json:"id"`
	TrackingNumber    string     `json:"trackingNumber"`
	Customer          ContactDTO `json:"customer"`
	Receiver          ContactDTO `json:"receiver"`
	CurrentCoords     PointDTO   `json:"currentCoords"`
	Origin            string     `json:"origin"`
	Destination       string     `json:"destination"`
	Status            string     `json:"status"`
	Carrier           string     `json:"carrier"`
	HubID             string     `json:"hubId,omitempty"`
	Weight            float64    `json:"weight"`
	Items             int        `json:"items"`
	EstimatedDelivery time.Time  `json:"estimatedDelivery"`
	CreatedAt         time.Time  `json:"createdAt"`
	Progress          float64    `json:"progress"`
}

type CreateRequest struct {
	Customer ContactDTO `json:"customer"`
	Receiver ContactDTO `json:"receiver"`
	Carrier  string     `json:"carrier"`
	Weight   float64    `json:"weight"`
	Items    int        `json:"items"`
}

type UpdateRequest struct {
	Customer          *ContactDTO `json:"customer,omitempty"`
	Receiver          *ContactDTO `json:"receiver,omitempty"`
	Carrier           *string     `json:"carrier,omitempty"`
	Weight            *float64    `json:"weight,omitempty"`
	Items             *int        `json:"items,omitempty"`
	EstimatedDelivery *time.Time  `json:"estimatedDelivery,omitempty"`
}

func (h *Handler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	filter := domain.Filter{
		Page: page, Limit: limit,
		Search: c.Query("search", ""), Status: c.Query("status", ""),
		ExcludeStatus: c.Query("exclude_status", ""),
	}
	shipments, total, err := h.service.List(c.UserContext(), filter)
	if err != nil {
		return utils.Error(c, 500, "failed to list shipments")
	}
	return utils.SuccessWithPagination(c, shipmentDTOs(shipments), filter.Page, filter.Limit, int(total))
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	shipment, err := h.service.Create(c.UserContext(), application.CreateInput{
		Customer: req.Customer.contact(), Receiver: req.Receiver.contact(),
		Carrier: req.Carrier, Weight: req.Weight, Items: req.Items,
	})
	if err != nil {
		if errors.Is(err, application.ErrInvalidInput) {
			message := strings.TrimPrefix(err.Error(), application.ErrInvalidInput.Error()+": ")
			return utils.Error(c, 400, message)
		}
		return utils.Error(c, 500, "failed to create shipment")
	}
	return utils.Success(c, shipmentDTO(*shipment))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	shipment, err := h.service.GetByID(c.UserContext(), c.Params("orderId"))
	if err != nil {
		if errors.Is(err, application.ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		return utils.Error(c, 500, "failed to find shipment")
	}
	return utils.Success(c, shipmentDTO(*shipment))
}

func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	var req struct {
		Status string `json:"status"`
		HubID  string `json:"hubId,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	shipment, err := h.service.UpdateStatus(c.UserContext(), c.Params("orderId"), domain.Status(req.Status), req.HubID)
	if err != nil {
		if errors.Is(err, application.ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		if errors.Is(err, application.ErrInvalidHub) || errors.Is(err, application.ErrInvalidInput) {
			message := strings.TrimPrefix(err.Error(), application.ErrInvalidInput.Error()+": ")
			return utils.Error(c, 400, message)
		}
		return utils.Error(c, 500, "failed to update shipment")
	}
	return utils.Success(c, shipmentDTO(*shipment))
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	shipment, err := h.service.Update(c.UserContext(), c.Params("orderId"), application.UpdateInput{
		Customer: contactPtr(req.Customer), Receiver: contactPtr(req.Receiver),
		Carrier: req.Carrier, Weight: req.Weight, Items: req.Items,
		EstimatedDelivery: req.EstimatedDelivery,
	})
	if err != nil {
		if errors.Is(err, application.ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		return utils.Error(c, 500, "failed to update shipment")
	}
	return utils.Success(c, shipmentDTO(*shipment))
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.UserContext(), c.Params("orderId")); err != nil {
		if errors.Is(err, application.ErrNotFound) {
			return utils.Error(c, 404, "shipment not found")
		}
		return utils.Error(c, 500, "failed to delete shipment")
	}
	return utils.Success(c, fiber.Map{"message": "shipment deleted"})
}

func (c ContactDTO) contact() domain.Contact {
	return domain.Contact{
		Name: c.Name, Zipcode: c.Zipcode, SubDistrict: c.SubDistrict,
		District: c.District, Province: c.Province,
		Coords: domain.GeoPoint{Lat: c.Coords.Lat, Lng: c.Coords.Lng},
	}
}

func contactPtr(contact *ContactDTO) *domain.Contact {
	if contact == nil {
		return nil
	}
	value := contact.contact()
	return &value
}

func shipmentDTO(shipment domain.Shipment) ShipmentDTO {
	return ShipmentDTO{
		ID: shipment.OrderID, TrackingNumber: shipment.TrackingNumber,
		Customer: contactDTO(shipment.Customer), Receiver: contactDTO(shipment.Receiver),
		CurrentCoords: PointDTO{Lat: shipment.CurrentCoords.Lat, Lng: shipment.CurrentCoords.Lng},
		Origin: shipment.Origin, Destination: shipment.Destination, Status: string(shipment.Status),
		Carrier: shipment.Carrier, HubID: shipment.HubID, Weight: shipment.Weight, Items: shipment.Items,
		EstimatedDelivery: shipment.EstimatedDelivery, CreatedAt: shipment.CreatedAt, Progress: shipment.Progress,
	}
}

func shipmentDTOs(shipments []domain.Shipment) []ShipmentDTO {
	responses := make([]ShipmentDTO, len(shipments))
	for i, shipment := range shipments {
		responses[i] = shipmentDTO(shipment)
	}
	return responses
}

func contactDTO(contact domain.Contact) ContactDTO {
	return ContactDTO{
		Name: contact.Name, Zipcode: contact.Zipcode, SubDistrict: contact.SubDistrict,
		District: contact.District, Province: contact.Province,
		Coords: PointDTO{Lat: contact.Coords.Lat, Lng: contact.Coords.Lng},
	}
}
