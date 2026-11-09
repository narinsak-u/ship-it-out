package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	hubdomain "github.com/narinsak-u/backend/internal/hub/domain"
	"github.com/narinsak-u/backend/internal/shipment/domain"
)

var (
	ErrNotFound    = domain.ErrNotFound
	ErrInvalidHub  = domain.ErrInvalidHub
	ErrInvalidInput = errors.New("invalid shipment input")
)

type Repository interface {
	List(context.Context, domain.Filter) ([]domain.Shipment, int64, error)
	FindByOrderID(context.Context, string) (*domain.Shipment, error)
	CreateWithEvent(context.Context, *domain.Shipment, *domain.ShipmentEvent) error
	SaveWithEvent(context.Context, *domain.Shipment, *domain.ShipmentEvent) error
	DeleteWithEvents(context.Context, *domain.Shipment) error
}
type TrackingReader interface {
	FindByTrackingNumber(context.Context, string) (*domain.Shipment, error)
	FindEventsByShipmentID(context.Context, uint) ([]domain.ShipmentEvent, error)
}

type HubReader interface {
	FindByID(context.Context, string) (*hubdomain.Hub, error)
}

type Service struct {
	repo Repository
	hubs HubReader
	now  func() time.Time
}

func NewService(repo Repository, hubs HubReader) *Service {
	return &Service{repo: repo, hubs: hubs, now: time.Now}
}

type CreateInput struct {
	Customer domain.Contact
	Receiver domain.Contact
	Carrier  string
	Weight   float64
	Items    int
}

type UpdateInput struct {
	Customer          *domain.Contact
	Receiver          *domain.Contact
	Carrier           *string
	Weight            *float64
	Items             *int
	EstimatedDelivery *time.Time
}

func (s *Service) List(ctx context.Context, filter domain.Filter) ([]domain.Shipment, int64, error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, orderID string) (*domain.Shipment, error) {
	return s.repo.FindByOrderID(ctx, orderID)
}

func (s *Service) Create(ctx context.Context, req CreateInput) (*domain.Shipment, error) {
	if req.Customer.Name == "" || req.Customer.Province == "" || req.Customer.Zipcode == "" {
		return nil, fmt.Errorf("%w: customer name, province, and zipcode are required", ErrInvalidInput)
	}
	if req.Receiver.Name == "" || req.Receiver.Province == "" || req.Receiver.Zipcode == "" {
		return nil, fmt.Errorf("%w: receiver name, province, and zipcode are required", ErrInvalidInput)
	}
	if req.Weight <= 0 {
		return nil, fmt.Errorf("%w: weight must be greater than 0", ErrInvalidInput)
	}
	if req.Items < 1 {
		return nil, fmt.Errorf("%w: items must be at least 1", ErrInvalidInput)
	}
	if req.Carrier == "" {
		return nil, fmt.Errorf("%w: carrier is required", ErrInvalidInput)
	}

	shipment := &domain.Shipment{
		Customer: req.Customer, Receiver: req.Receiver,
		Origin: composeAddress(req.Customer), Destination: composeAddress(req.Receiver),
		CurrentCoords: req.Customer.Coords, Status: domain.StatusPending, Carrier: req.Carrier,
		Weight: req.Weight, Items: req.Items,
		EstimatedDelivery: s.now().Add(72 * time.Hour), Progress: 0,
	}
	event := domain.ShipmentEvent{
		Status: "Label Created", Description: "Awaiting pickup.",
		Location: domain.Location{Name: composeAddress(req.Customer), Lat: req.Customer.Coords.Lat, Lng: req.Customer.Coords.Lng},
	}
	if err := s.repo.CreateWithEvent(ctx, shipment, &event); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *Service) UpdateStatus(ctx context.Context, orderID string, status domain.Status, hubID string) (*domain.Shipment, error) {
	if status == "" {
		return nil, fmt.Errorf("%w: status is required", ErrInvalidInput)
	}
	shipment, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	shipment.Status = status

	var hub *hubdomain.Hub
	if hubID != "" {
		hub, err = s.hubs.FindByID(ctx, hubID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, ErrInvalidHub
			}
			return nil, err
		}
		shipment.HubID = hubID
		shipment.CurrentCoords = domain.GeoPoint{Lat: hub.Coords.Lat, Lng: hub.Coords.Lng}
	}

	event := statusEvent(*shipment, hub, status)
	event.ShipmentID = shipment.ID
	if err := s.repo.SaveWithEvent(ctx, shipment, &event); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *Service) Update(ctx context.Context, orderID string, req UpdateInput) (*domain.Shipment, error) {
	shipment, err := s.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
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

func composeAddress(c domain.Contact) string {
	return fmt.Sprintf("%s, %s, %s", c.SubDistrict, c.District, c.Province)
}

func statusEvent(shipment domain.Shipment, hub *hubdomain.Hub, targetStatus domain.Status) domain.ShipmentEvent {
	var eventStatus, description string
	var location domain.Location
	hubLocation := func(fallback string, lat, lng float64) domain.Location {
		if hub != nil {
			return domain.Location{Name: hub.Name + ", " + hub.Address, Lat: hub.Coords.Lat, Lng: hub.Coords.Lng}
		}
		return domain.Location{Name: fallback, Lat: lat, Lng: lng}
	}

	switch targetStatus {
	case domain.StatusPending:
		eventStatus, description = "Label Created", "Awaiting pickup."
		location = domain.Location{Name: composeAddress(shipment.Customer), Lat: shipment.Customer.Coords.Lat, Lng: shipment.Customer.Coords.Lng}
	case domain.StatusPickedUp:
		eventStatus, description = "Picked Up", "Parcel collected from sender."
		location = domain.Location{Name: composeAddress(shipment.Customer), Lat: shipment.Customer.Coords.Lat, Lng: shipment.Customer.Coords.Lng}
	case domain.StatusDeparted:
		eventStatus, description = "Departed", "In transit to hub."
		location = hubLocation(shipment.Origin, shipment.CurrentCoords.Lat, shipment.CurrentCoords.Lng)
	case domain.StatusInTransit:
		eventStatus, description = "In Transit", "Transit to next hub."
		location = hubLocation(shipment.Destination, shipment.CurrentCoords.Lat, shipment.CurrentCoords.Lng)
	case domain.StatusOutForDelivery:
		eventStatus, description = "Out for Delivery", "Out for delivery."
		location = hubLocation(shipment.Destination, shipment.Receiver.Coords.Lat, shipment.Receiver.Coords.Lng)
	case domain.StatusDelivered:
		eventStatus, description = "Delivered", "Delivered to recipient."
		location = domain.Location{Name: composeAddress(shipment.Receiver), Lat: shipment.Receiver.Coords.Lat, Lng: shipment.Receiver.Coords.Lng}
	case domain.StatusDelayed:
		eventStatus, description = "Delayed", "Unexpected issue encountered."
		location = hubLocation(shipment.Destination, shipment.CurrentCoords.Lat, shipment.CurrentCoords.Lng)
	default:
		eventStatus, description = string(targetStatus), "Status updated."
		location = domain.Location{Name: shipment.Destination, Lat: shipment.CurrentCoords.Lat, Lng: shipment.CurrentCoords.Lng}
	}
	return domain.ShipmentEvent{Status: eventStatus, Description: description, Location: location}
}
