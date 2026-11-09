package infrastructure

import (
	"testing"

	"github.com/narinsak-u/backend/internal/shipment/domain"
	"github.com/stretchr/testify/assert"
)

func TestShipmentRecordRoundTripPreservesCoordinates(t *testing.T) {
	shipment := domain.Shipment{
		Customer: domain.Contact{Coords: domain.GeoPoint{Lat: 10, Lng: 20}},
		Receiver: domain.Contact{Coords: domain.GeoPoint{Lat: 30, Lng: 40}},
		CurrentCoords: domain.GeoPoint{Lat: 25, Lng: 35},
	}

	stored := shipmentRecordFromDomain(shipment)
	assert.Equal(t, 10.0, stored.CustomerLat)
	assert.Equal(t, 20.0, stored.CustomerLng)
	assert.Equal(t, 30.0, stored.ReceiverLat)
	assert.Equal(t, 40.0, stored.ReceiverLng)
	assert.Equal(t, 25.0, stored.CurrentLat)
	assert.Equal(t, 35.0, stored.CurrentLng)

	loaded := shipmentDomain(stored)
	assert.Equal(t, shipment.Customer.Coords, loaded.Customer.Coords)
	assert.Equal(t, shipment.Receiver.Coords, loaded.Receiver.Coords)
	assert.Equal(t, shipment.CurrentCoords, loaded.CurrentCoords)
}
