package application

import (
	"context"
	"errors"

	shipmentapp "github.com/narinsak-u/backend/internal/shipment/application"
	"github.com/narinsak-u/backend/internal/shipment/domain"
)

var (
	ErrShipmentNotFound = errors.New("shipment not found")
	ErrEventsUnavailable = errors.New("failed to fetch events")
)

type TrackingResult struct {
	Shipment *domain.Shipment
	Events []domain.ShipmentEvent
}

type Service struct { reader shipmentapp.TrackingReader }

func NewService(reader shipmentapp.TrackingReader) *Service { return &Service{reader: reader} }

func (s *Service) Track(ctx context.Context, trackingNumber string) (*TrackingResult, error) {
	shipment, err := s.reader.FindByTrackingNumber(ctx, trackingNumber)
	if errors.Is(err, domain.ErrNotFound) || errors.Is(err, shipmentapp.ErrNotFound) {
		return nil, ErrShipmentNotFound
	}
	if err != nil { return nil, err }
	events, err := s.reader.FindEventsByShipmentID(ctx, shipment.ID)
	if err != nil { return nil, errors.Join(ErrEventsUnavailable, err) }
	return &TrackingResult{Shipment: shipment, Events: events}, nil
}
