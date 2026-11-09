package application

import (
	"context"
	"errors"
	"testing"

	"github.com/narinsak-u/backend/internal/shipment/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type reader struct {
	shipment *domain.Shipment
	events []domain.ShipmentEvent
	shipmentErr error
	eventsErr error
}

func (r reader) FindByTrackingNumber(context.Context, string) (*domain.Shipment, error) { return r.shipment, r.shipmentErr }
func (r reader) FindEventsByShipmentID(context.Context, uint) ([]domain.ShipmentEvent, error) { return r.events, r.eventsErr }

func TestTrack(t *testing.T) {
	shipment := &domain.Shipment{ID: 7, OrderID: "ORD-7", TrackingNumber: "TH-7"}
	events := []domain.ShipmentEvent{{Status: "in_transit"}}
	result, err := NewService(reader{shipment: shipment, events: events}).Track(context.Background(), "TH-7")
	require.NoError(t, err)
	assert.Same(t, shipment, result.Shipment)
	assert.Equal(t, events, result.Events)
}

func TestTrackNotFound(t *testing.T) {
	_, err := NewService(reader{shipmentErr: domain.ErrNotFound}).Track(context.Background(), "TH-7")
	assert.ErrorIs(t, err, ErrShipmentNotFound)
}

func TestTrackRepositoryFailures(t *testing.T) {
	dbErr := errors.New("database unavailable")
	_, err := NewService(reader{shipmentErr: dbErr}).Track(context.Background(), "TH-7")
	assert.ErrorIs(t, err, dbErr)

	_, err = NewService(reader{shipment: &domain.Shipment{ID: 7}, eventsErr: dbErr}).Track(context.Background(), "TH-7")
	assert.ErrorIs(t, err, ErrEventsUnavailable)
	assert.ErrorIs(t, err, dbErr)
}
