package presentation

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/shipment/domain"
	"github.com/narinsak-u/backend/internal/tracking/application"
	"github.com/narinsak-u/backend/pkg/utils"
)

type Handler struct { service *application.Service }
func NewHandler(service *application.Service) *Handler { return &Handler{service: service} }

type PointDTO struct { Lat float64 `json:"lat"`; Lng float64 `json:"lng"` }
type ContactDTO struct {
	Name string `json:"name"`
	Zipcode string `json:"zipcode"`
	SubDistrict string `json:"subDistrict"`
	District string `json:"district"`
	Province string `json:"province"`
	Coords PointDTO `json:"coords"`
}
type ShipmentDTO struct {
	ID string `json:"id"`
	TrackingNumber string `json:"trackingNumber"`
	Customer ContactDTO `json:"customer"`
	Receiver ContactDTO `json:"receiver"`
	CurrentCoords PointDTO `json:"currentCoords"`
	Origin string `json:"origin"`
	Destination string `json:"destination"`
	Status string `json:"status"`
	Carrier string `json:"carrier"`
	HubID string `json:"hubId,omitempty"`
	Weight float64 `json:"weight"`
	Items int `json:"items"`
	EstimatedDelivery time.Time `json:"estimatedDelivery"`
	CreatedAt time.Time `json:"createdAt"`
	Progress float64 `json:"progress"`
}
type EventDTO struct {
	ID uint `json:"id"`
	ShipmentID uint `json:"shipmentId"`
	Status string `json:"status"`
	Location LocationDTO `json:"location"`
	Description string `json:"description,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
type LocationDTO struct { Name string `json:"name"`; Lat float64 `json:"lat"`; Lng float64 `json:"lng"` }

func (h *Handler) Track(c *fiber.Ctx) error {
	result, err := h.service.Track(c.UserContext(), c.Params("trackingNumber"))
	if errors.Is(err, application.ErrShipmentNotFound) { return utils.Error(c, 404, "shipment not found") }
	if err != nil { return utils.Error(c, 500, "failed to fetch events") }
	return utils.Success(c, fiber.Map{"shipment": shipmentDTO(*result.Shipment), "events": eventDTOs(result.Events)})
}

func shipmentDTO(s domain.Shipment) ShipmentDTO {
	return ShipmentDTO{ID:s.OrderID, TrackingNumber:s.TrackingNumber, Customer:contactDTO(s.Customer), Receiver:contactDTO(s.Receiver), CurrentCoords:PointDTO{s.CurrentCoords.Lat,s.CurrentCoords.Lng}, Origin:s.Origin, Destination:s.Destination, Status:string(s.Status), Carrier:s.Carrier, HubID:s.HubID, Weight:s.Weight, Items:s.Items, EstimatedDelivery:s.EstimatedDelivery, CreatedAt:s.CreatedAt, Progress:s.Progress}
}
func contactDTO(c domain.Contact) ContactDTO {
	return ContactDTO{Name:c.Name, Zipcode:c.Zipcode, SubDistrict:c.SubDistrict, District:c.District, Province:c.Province, Coords:PointDTO{c.Coords.Lat,c.Coords.Lng}}
}
func eventDTOs(events []domain.ShipmentEvent) []EventDTO {
	result := make([]EventDTO, len(events))
	for i, e := range events { result[i] = EventDTO{ID:e.ID, ShipmentID:e.ShipmentID, Status:e.Status, Location:LocationDTO{e.Location.Name,e.Location.Lat,e.Location.Lng}, Description:e.Description, Timestamp:e.CreatedAt} }
	return result
}
