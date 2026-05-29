package models_test

import (
	"testing"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHub_BeforeSave(t *testing.T) {
	h := &models.Hub{Coords: models.GeoPoint{Lat: 13.0, Lng: 100.0}}
	err := h.BeforeSave(nil)
	assert.NoError(t, err)
	assert.Equal(t, 13.0, h.Lat)
	assert.Equal(t, 100.0, h.Lng)
}

func TestHub_AfterFind(t *testing.T) {
	h := &models.Hub{Lat: 13.0, Lng: 100.0}
	err := h.AfterFind(nil)
	assert.NoError(t, err)
	assert.Equal(t, 13.0, h.Coords.Lat)
	assert.Equal(t, 100.0, h.Coords.Lng)
}

func TestHub_RoundTrip(t *testing.T) {
	h := &models.Hub{Coords: models.GeoPoint{Lat: 10.0, Lng: 20.0}}
	h.BeforeSave(nil)
	h.AfterFind(nil)
	assert.Equal(t, 10.0, h.Coords.Lat)
	assert.Equal(t, 20.0, h.Coords.Lng)
}
