package shipment

import (
	"testing"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestStatusToEvent_AllStatuses(t *testing.T) {
	shipment := models.Shipment{
		Customer:    models.ContactInfo{SubDistrict: "S", District: "D", Province: "P"},
		Receiver:    models.ContactInfo{SubDistrict: "R", District: "R2", Province: "P2"},
		CustomerLat: 10.0, CustomerLng: 20.0,
		ReceiverLat: 30.0, ReceiverLng: 40.0,
		CurrentLat: 15.0, CurrentLng: 25.0,
		Origin:      "S, D, P",
		Destination: "R, R2, P2",
	}

	hub := &models.Hub{
		ID: "HUB-001", Name: "Bangkok Hub", Address: "Bangkok",
		Coords: models.GeoPoint{Lat: 13.7563, Lng: 100.5018},
	}

	tests := []struct {
		name     string
		status   string
		hasHub   bool
		wantDesc string
	}{
		{"pending", "pending", false, "Awaiting pickup."},
		{"picked_up", "picked_up", false, "Parcel collected from sender."},
		{"departed with hub", "departed", true, "In transit to hub."},
		{"departed without hub", "departed", false, "In transit to hub."},
		{"in_transit with hub", "in_transit", true, "Transit to next hub."},
		{"in_transit without hub", "in_transit", false, "Transit to next hub."},
		{"out_for_delivery with hub", "out_for_delivery", true, "Out for delivery."},
		{"out_for_delivery without hub", "out_for_delivery", false, "Out for delivery."},
		{"delivered", "delivered", false, "Delivered to recipient."},
		{"delayed with hub", "delayed", true, "Unexpected issue encountered."},
		{"delayed without hub", "delayed", false, "Unexpected issue encountered."},
		{"default status", "custom_status", false, "Status updated."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h *models.Hub
			if tt.hasHub {
				h = hub
			}
			event := statusToEvent(shipment, h, tt.status)
			assert.Equal(t, tt.wantDesc, event.Description)
			assert.NotEmpty(t, event.Status)
			assert.NotEmpty(t, event.Location.Name)
		})
	}
}

func TestStatusToEvent_LocationSources(t *testing.T) {
	shipment := models.Shipment{
		Customer:    models.ContactInfo{SubDistrict: "S", District: "D", Province: "P"},
		Receiver:    models.ContactInfo{SubDistrict: "R", District: "R2", Province: "P2"},
		CustomerLat: 10.0, CustomerLng: 20.0,
		ReceiverLat: 30.0, ReceiverLng: 40.0,
		CurrentLat: 15.0, CurrentLng: 25.0,
		Origin:      "S, D, P",
		Destination: "R, R2, P2",
	}

	e1 := statusToEvent(shipment, nil, "pending")
	assert.Equal(t, 10.0, e1.Location.Lat)
	assert.Equal(t, 20.0, e1.Location.Lng)

	e2 := statusToEvent(shipment, nil, "delivered")
	assert.Equal(t, 30.0, e2.Location.Lat)
	assert.Equal(t, 40.0, e2.Location.Lng)

	hub := &models.Hub{Coords: models.GeoPoint{Lat: 50.0, Lng: 60.0}}
	e3 := statusToEvent(shipment, hub, "in_transit")
	assert.Equal(t, 50.0, e3.Location.Lat)
	assert.Equal(t, 60.0, e3.Location.Lng)
}
