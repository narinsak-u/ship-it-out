package tracking

import (
	"errors"
	"testing"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeShipmentQuery struct {
	shipment *models.Shipment
	events   []models.ShipmentEvent
	findErr  error
	eventErr error
}

func (f fakeShipmentQuery) FindByTrackingNumber(string) (*models.Shipment, error) {
	return f.shipment, f.findErr
}

func (f fakeShipmentQuery) FindEventsByShipmentID(uint) ([]models.ShipmentEvent, error) {
	return f.events, f.eventErr
}

func TestServiceTrack(t *testing.T) {
	shipment := &models.Shipment{ID: 7, TrackingNumber: "TH-7"}
	events := []models.ShipmentEvent{{Status: "Label Created"}}
	result, err := NewService(fakeShipmentQuery{shipment: shipment, events: events}).Track("TH-7")
	require.NoError(t, err)
	assert.Same(t, shipment, result.Shipment)
	assert.Equal(t, events, result.Events)
}

func TestServiceTrack_LookupFailure(t *testing.T) {
	_, err := NewService(fakeShipmentQuery{findErr: errors.New("database unavailable")}).Track("TH-7")
	assert.ErrorIs(t, err, ErrShipmentNotFound)
}

func TestServiceTrack_EventQueryFailure(t *testing.T) {
	_, err := NewService(fakeShipmentQuery{
		shipment: &models.Shipment{ID: 7},
		eventErr: errors.New("database unavailable"),
	}).Track("TH-7")
	assert.ErrorIs(t, err, ErrEventsUnavailable)
}
