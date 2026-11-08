package shipment

import (
	"errors"
	"time"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

var (
	ErrInvalidHub = errors.New("invalid hub ID")
	ErrNotFound   = errors.New("shipment not found")
)

type WorkflowRepository interface {
	FindByOrderID(string) (*models.Shipment, error)
	CreateWithEvent(*models.Shipment, *models.ShipmentEvent) error
	SaveWithEvent(*models.Shipment, *models.ShipmentEvent) error
	DeleteWithEvents(*models.Shipment) error
}

type HandlerRepository interface {
	WorkflowRepository
	List(ShipmentFilter) ([]models.Shipment, int64, error)
}

type Service struct {
	repo WorkflowRepository
	hubs HubRepository
}

func NewService(repo WorkflowRepository, hubs HubRepository) *Service {
	return &Service{repo: repo, hubs: hubs}
}

func (s *Service) Create(req CreateRequest) (*models.Shipment, error) {
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
		CurrentCoords: req.Customer.Coords, Status: "pending", Carrier: req.Carrier,
		Weight: req.Weight, Items: req.Items, EstimatedDelivery: time.Now().Add(72 * time.Hour), Progress: 0,
	}
	event := &models.ShipmentEvent{Status: "Label Created", Description: "Awaiting pickup.", Location: models.Location{
		Name: utils.ComposeAddress(req.Customer), Lat: req.Customer.Coords.Lat, Lng: req.Customer.Coords.Lng,
	}}
	if err := s.repo.CreateWithEvent(shipment, event); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *Service) UpdateStatus(orderID, status, hubID string) (*models.Shipment, error) {
	if status == "" {
		return nil, errors.New("status is required")
	}
	shipment, err := s.repo.FindByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	shipment.Status = status
	var hub *models.Hub
	if hubID != "" {
		hubRecord, err := s.hubs.FindByID(hubID)
		if err != nil {
			if err.Error() == "record not found" {
				return nil, ErrInvalidHub
			}
			return nil, err
		}
		hub = hubRecord
		shipment.HubID = hubID
		shipment.CurrentCoords.Lat = hub.Lat
		shipment.CurrentCoords.Lng = hub.Lng
	}
	event := statusToEvent(*shipment, hub, status)
	event.ShipmentID = shipment.ID
	if err := s.repo.SaveWithEvent(shipment, &event); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *Service) Update(orderID string, req UpdateRequest) (*models.Shipment, error) {
	shipment, err := s.repo.FindByOrderID(orderID)
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
	if err := s.repo.SaveWithEvent(shipment, nil); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *Service) Delete(orderID string) error {
	shipment, err := s.repo.FindByOrderID(orderID)
	if err != nil {
		return err
	}
	return s.repo.DeleteWithEvents(shipment)
}
