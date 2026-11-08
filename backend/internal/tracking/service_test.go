package tracking

import (
	"context"
	"errors"
	"testing"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeShipmentQuery struct {
	shipment *models.Shipment
	events   []models.ShipmentEvent
	findErr  error
	eventErr error
}

func (f fakeShipmentQuery) FindByTrackingNumber(context.Context, string) (*models.Shipment, error) {
	return f.shipment, f.findErr
}

func (f fakeShipmentQuery) FindEventsByShipmentID(context.Context, uint) ([]models.ShipmentEvent, error) {
	return f.events, f.eventErr
}

func TestServiceTrack(t *testing.T) {
	record := &models.Shipment{ID: 7, TrackingNumber: "TH-7"}
	events := []models.ShipmentEvent{{Status: "Label Created"}}
	result, err := NewService(fakeShipmentQuery{shipment: record, events: events}).Track(context.Background(), "TH-7")
	require.NoError(t, err)
	assert.Same(t, record, result.Shipment)
	assert.Equal(t, events, result.Events)
}

func TestServiceTrack_NotFound(t *testing.T) {
	_, err := NewService(fakeShipmentQuery{findErr: shipment.ErrNotFound}).Track(context.Background(), "TH-7")
	assert.ErrorIs(t, err, ErrShipmentNotFound)
}

func TestServiceTrack_LookupFailure(t *testing.T) {
	want := errors.New("database unavailable")
	_, err := NewService(fakeShipmentQuery{findErr: want}).Track(context.Background(), "TH-7")
	assert.ErrorIs(t, err, want)
}

func TestServiceTrack_EventQueryFailure(t *testing.T) {
	_, err := NewService(fakeShipmentQuery{
		shipment: &models.Shipment{ID: 7},
		eventErr: errors.New("database unavailable"),
	}).Track(context.Background(), "TH-7")
	assert.ErrorIs(t, err, ErrEventsUnavailable)
}
