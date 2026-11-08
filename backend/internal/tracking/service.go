package tracking

import (
	"errors"

	"github.com/narinsak-u/backend/internal/models"
)

var (
	ErrShipmentNotFound  = errors.New("shipment not found")
	ErrEventsUnavailable = errors.New("failed to fetch events")
)

type ShipmentQuery interface {
	FindByTrackingNumber(string) (*models.Shipment, error)
	FindEventsByShipmentID(uint) ([]models.ShipmentEvent, error)
}

type Service struct {
	query ShipmentQuery
}

type TrackingResult struct {
	Shipment *models.Shipment
	Events   []models.ShipmentEvent
}

func NewService(query ShipmentQuery) *Service {
	return &Service{query: query}
}

func (s *Service) Track(trackingNumber string) (*TrackingResult, error) {
	shipment, err := s.query.FindByTrackingNumber(trackingNumber)
	if err != nil {
		return nil, ErrShipmentNotFound
	}
	events, err := s.query.FindEventsByShipmentID(shipment.ID)
	if err != nil {
		return nil, ErrEventsUnavailable
	}
	return &TrackingResult{Shipment: shipment, Events: events}, nil
}
