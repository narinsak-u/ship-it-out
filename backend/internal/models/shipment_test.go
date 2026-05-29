package models_test

import (
	"testing"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestShipment_BeforeSave(t *testing.T) {
	s := &models.Shipment{
		Customer:      models.ContactInfo{Coords: models.GeoPoint{Lat: 10.0, Lng: 20.0}},
		Receiver:      models.ContactInfo{Coords: models.GeoPoint{Lat: 30.0, Lng: 40.0}},
		CurrentCoords: models.GeoPoint{Lat: 25.0, Lng: 35.0},
	}
	err := s.BeforeSave(nil)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, s.CustomerLat)
	assert.Equal(t, 20.0, s.CustomerLng)
	assert.Equal(t, 30.0, s.ReceiverLat)
	assert.Equal(t, 40.0, s.ReceiverLng)
	assert.Equal(t, 25.0, s.CurrentLat)
	assert.Equal(t, 35.0, s.CurrentLng)
}

func TestShipment_AfterFind(t *testing.T) {
	s := &models.Shipment{
		CustomerLat: 10.0, CustomerLng: 20.0,
		ReceiverLat: 30.0, ReceiverLng: 40.0,
		CurrentLat: 25.0, CurrentLng: 35.0,
	}
	err := s.AfterFind(nil)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, s.Customer.Coords.Lat)
	assert.Equal(t, 20.0, s.Customer.Coords.Lng)
	assert.Equal(t, 30.0, s.Receiver.Coords.Lat)
	assert.Equal(t, 40.0, s.Receiver.Coords.Lng)
	assert.Equal(t, 25.0, s.CurrentCoords.Lat)
	assert.Equal(t, 35.0, s.CurrentCoords.Lng)
}

func TestShipment_RoundTrip(t *testing.T) {
	s := &models.Shipment{
		Customer:      models.ContactInfo{Coords: models.GeoPoint{Lat: 10.0, Lng: 20.0}},
		Receiver:      models.ContactInfo{Coords: models.GeoPoint{Lat: 30.0, Lng: 40.0}},
		CurrentCoords: models.GeoPoint{Lat: 25.0, Lng: 35.0},
	}
	s.BeforeSave(nil)
	s.AfterFind(nil)
	assert.Equal(t, 10.0, s.Customer.Coords.Lat)
	assert.Equal(t, 20.0, s.Customer.Coords.Lng)
	assert.Equal(t, 30.0, s.Receiver.Coords.Lat)
	assert.Equal(t, 40.0, s.Receiver.Coords.Lng)
	assert.Equal(t, 25.0, s.CurrentCoords.Lat)
	assert.Equal(t, 35.0, s.CurrentCoords.Lng)
}
