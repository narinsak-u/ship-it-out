package tracking

import (
	"context"
	"errors"

	"github.com/narinsak-u/backend/internal/models"
	shipmentapplication "github.com/narinsak-u/backend/internal/shipment/application"
)

var (
	ErrShipmentNotFound  = errors.New("shipment not found")
	ErrEventsUnavailable = errors.New("failed to fetch events")
)

type ShipmentQuery interface {
	FindByTrackingNumber(context.Context, string) (*models.Shipment, error)
	FindEventsByShipmentID(context.Context, uint) ([]models.ShipmentEvent, error)
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

func (s *Service) Track(ctx context.Context, trackingNumber string) (*TrackingResult, error) {
	shipmentRecord, err := s.query.FindByTrackingNumber(ctx, trackingNumber)
	if errors.Is(err, shipmentapplication.ErrNotFound) {
		return nil, ErrShipmentNotFound
	}
	if err != nil {
		return nil, err
	}
	events, err := s.query.FindEventsByShipmentID(ctx, shipmentRecord.ID)
	if err != nil {
		return nil, ErrEventsUnavailable
	}
	return &TrackingResult{Shipment: shipmentRecord, Events: events}, nil
}
