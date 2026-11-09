package shipment

import (
	"context"
	"errors"
	"time"

	hubapplication "github.com/narinsak-u/backend/internal/hub/application"
	hub "github.com/narinsak-u/backend/internal/hub/domain"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

var (
	ErrInvalidHub = errors.New("invalid hub ID")
	ErrNotFound   = errors.New("shipment not found")
)

const (
	StatusPending        = "pending"
	StatusPickedUp       = "picked_up"
	StatusDeparted       = "departed"
	StatusInTransit      = "in_transit"
	StatusOutForDelivery = "out_for_delivery"
	StatusDelivered      = "delivered"
	StatusDelayed        = "delayed"
)

type HubRepository interface {
	FindByID(context.Context, string) (*hub.Hub, error)
}

type WorkflowRepository interface {
	FindByOrderID(context.Context, string) (*models.Shipment, error)
	List(context.Context, ShipmentFilter) ([]models.Shipment, int64, error)
	CreateWithEvent(context.Context, *models.Shipment, *models.ShipmentEvent) error
	SaveWithEvent(context.Context, *models.Shipment, *models.ShipmentEvent) error
	DeleteWithEvents(context.Context, *models.Shipment) error
}

type Service struct {
	repo WorkflowRepository
	hubs HubRepository
}

func NewService(repo WorkflowRepository, hubs HubRepository) *Service {
	return &Service{repo: repo, hubs: hubs}
}

func (s *Service) List(ctx context.Context, filter ShipmentFilter) ([]models.Shipment, int64, error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, orderID string) (*models.Shipment, error) {
	return s.repo.FindByOrderID(ctx, orderID)
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*models.Shipment, error) {
	if req.Customer.Name == "" || req.Customer.Province == "" || req.Customer.Zipcode == "" {
		return nil, errors.New("customer name, province, and zipcode are required")
	}
	if req.Receiver.Name == "" || req.Receiver.Province == "" || req.Receiver.Zipcode == "" {
		return nil, errors.New("receiver name, province, and zipcode are required")
	}
	if req.Weight <= 0 {
		return nil, errors.New("weight must be greater than 0")
	}
	if req.Items < 1 {
		return nil, errors.New("items must be at least 1")
	}
	if req.Carrier == "" {
		return nil, errors.New("carrier is required")
	}
	shipment := &models.Shipment{
		Customer: req.Customer, Receiver: req.Receiver,
		Origin: utils.ComposeAddress(req.Customer), Destination: utils.ComposeAddress(req.Receiver),
		CurrentCoords: req.Customer.Coords, Status: StatusPending, Carrier: req.Carrier,
		Weight: req.Weight, Items: req.Items, EstimatedDelivery: time.Now().Add(72 * time.Hour), Progress: 0,
	}
	event := initialEvent(req.Customer)
	if err := s.repo.CreateWithEvent(ctx, shipment, &event); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *Service) UpdateStatus(ctx context.Context, orderID, status, hubID string) (*models.Shipment, error) {
	if status == "" {
		return nil, errors.New("status is required")
	}
	shipment, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	shipment.Status = status
	var hubRecord *hub.Hub
	if hubID != "" {
		hubRecord, err = s.hubs.FindByID(ctx, hubID)
		if err != nil {
			if errors.Is(err, hubapplication.ErrNotFound) {
				return nil, ErrInvalidHub
			}
			return nil, err
		}
		shipment.HubID = hubID
		shipment.CurrentCoords.Lat = hubRecord.Coords.Lat
		shipment.CurrentCoords.Lng = hubRecord.Coords.Lng
	}
	event := statusToEvent(*shipment, hubRecord, status)
	event.ShipmentID = shipment.ID
	if err := s.repo.SaveWithEvent(ctx, shipment, &event); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *Service) Update(ctx context.Context, orderID string, req UpdateRequest) (*models.Shipment, error) {
	shipment, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
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
	if err := s.repo.SaveWithEvent(ctx, shipment, nil); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *Service) Delete(ctx context.Context, orderID string) error {
	shipment, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return err
	}
	return s.repo.DeleteWithEvents(ctx, shipment)
}

func initialEvent(customer models.ContactInfo) models.ShipmentEvent {
	return models.ShipmentEvent{Status: "Label Created", Description: "Awaiting pickup.", Location: models.Location{
		Name: utils.ComposeAddress(customer), Lat: customer.Coords.Lat, Lng: customer.Coords.Lng,
	}}
}

func statusToEvent(shipment models.Shipment, hubRecord *hub.Hub, targetStatus string) models.ShipmentEvent {
	var eventStatus, description string
	var location models.Location
	hubLocation := func(fallback string, lat, lng float64) models.Location {
		if hubRecord != nil {
			return models.Location{Name: hubRecord.Name + ", " + hubRecord.Address, Lat: hubRecord.Coords.Lat, Lng: hubRecord.Coords.Lng}
		}
		return models.Location{Name: fallback, Lat: lat, Lng: lng}
	}
	switch targetStatus {
	case StatusPending:
		eventStatus, description = "Label Created", "Awaiting pickup."
		location = models.Location{Name: utils.ComposeAddress(shipment.Customer), Lat: shipment.CustomerLat, Lng: shipment.CustomerLng}
	case StatusPickedUp:
		eventStatus, description = "Picked Up", "Parcel collected from sender."
		location = models.Location{Name: utils.ComposeAddress(shipment.Customer), Lat: shipment.CustomerLat, Lng: shipment.CustomerLng}
	case StatusDeparted:
		eventStatus, description = "Departed", "In transit to hub."
		location = hubLocation(shipment.Origin, shipment.CurrentLat, shipment.CurrentLng)
	case StatusInTransit:
		eventStatus, description = "In Transit", "Transit to next hub."
		location = hubLocation(shipment.Destination, shipment.CurrentLat, shipment.CurrentLng)
	case StatusOutForDelivery:
		eventStatus, description = "Out for Delivery", "Out for delivery."
		location = hubLocation(shipment.Destination, shipment.ReceiverLat, shipment.ReceiverLng)
	case StatusDelivered:
		eventStatus, description = "Delivered", "Delivered to recipient."
		location = models.Location{Name: utils.ComposeAddress(shipment.Receiver), Lat: shipment.ReceiverLat, Lng: shipment.ReceiverLng}
	case StatusDelayed:
		eventStatus, description = "Delayed", "Unexpected issue encountered."
		location = hubLocation(shipment.Destination, shipment.CurrentLat, shipment.CurrentLng)
	default:
		eventStatus, description = targetStatus, "Status updated."
		location = models.Location{Name: shipment.Destination, Lat: shipment.CurrentLat, Lng: shipment.CurrentLng}
	}
	return models.ShipmentEvent{Status: eventStatus, Description: description, Location: location}
}
