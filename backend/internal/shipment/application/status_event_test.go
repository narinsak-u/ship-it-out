package application

import (
	"testing"

	hubdomain "github.com/narinsak-u/backend/internal/hub/domain"
	"github.com/narinsak-u/backend/internal/shipment/domain"
)

func TestStatusEventMappings(t *testing.T) {
	shipment := domain.Shipment{
		Customer: domain.Contact{
			SubDistrict: "S", District: "D", Province: "P",
			Coords: domain.GeoPoint{Lat: 10, Lng: 20},
		},
		Receiver: domain.Contact{
			SubDistrict: "R", District: "R2", Province: "P2",
			Coords: domain.GeoPoint{Lat: 30, Lng: 40},
		},
		CurrentCoords: domain.GeoPoint{Lat: 15, Lng: 25},
		Origin: "S, D, P", Destination: "R, R2, P2",
	}
	hub := &hubdomain.Hub{
		Name: "Bangkok Hub", Address: "Bangkok",
		Coords: hubdomain.GeoPoint{Lat: 13.7563, Lng: 100.5018},
	}
	tests := []struct {
		name, wantStatus, wantDescription, wantLocation string
		status domain.Status
		withHub bool
	}{
		{name: "pending", status: domain.StatusPending, wantStatus: "Label Created", wantDescription: "Awaiting pickup.", wantLocation: "S, D, P"},
		{name: "picked up", status: domain.StatusPickedUp, wantStatus: "Picked Up", wantDescription: "Parcel collected from sender.", wantLocation: "S, D, P"},
		{name: "departed hub", status: domain.StatusDeparted, wantStatus: "Departed", wantDescription: "In transit to hub.", wantLocation: "Bangkok Hub, Bangkok", withHub: true},
		{name: "departed fallback", status: domain.StatusDeparted, wantStatus: "Departed", wantDescription: "In transit to hub.", wantLocation: "S, D, P"},
		{name: "transit hub", status: domain.StatusInTransit, wantStatus: "In Transit", wantDescription: "Transit to next hub.", wantLocation: "Bangkok Hub, Bangkok", withHub: true},
		{name: "transit fallback", status: domain.StatusInTransit, wantStatus: "In Transit", wantDescription: "Transit to next hub.", wantLocation: "R, R2, P2"},
		{name: "out for delivery hub", status: domain.StatusOutForDelivery, wantStatus: "Out for Delivery", wantDescription: "Out for delivery.", wantLocation: "Bangkok Hub, Bangkok", withHub: true},
		{name: "out for delivery fallback", status: domain.StatusOutForDelivery, wantStatus: "Out for Delivery", wantDescription: "Out for delivery.", wantLocation: "R, R2, P2"},
		{name: "delivered", status: domain.StatusDelivered, wantStatus: "Delivered", wantDescription: "Delivered to recipient.", wantLocation: "R, R2, P2"},
		{name: "delayed hub", status: domain.StatusDelayed, wantStatus: "Delayed", wantDescription: "Unexpected issue encountered.", wantLocation: "Bangkok Hub, Bangkok", withHub: true},
		{name: "delayed fallback", status: domain.StatusDelayed, wantStatus: "Delayed", wantDescription: "Unexpected issue encountered.", wantLocation: "R, R2, P2"},
		{name: "custom", status: "custom_status", wantStatus: "custom_status", wantDescription: "Status updated.", wantLocation: "R, R2, P2"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var selectedHub *hubdomain.Hub
			if test.withHub {
				selectedHub = hub
			}
			event := statusEvent(shipment, selectedHub, test.status)
			if event.Status != test.wantStatus || event.Description != test.wantDescription || event.Location.Name != test.wantLocation {
				t.Fatalf("got %+v", event)
			}
		})
	}

	pending := statusEvent(shipment, nil, domain.StatusPending)
	if pending.Location.Lat != 10 || pending.Location.Lng != 20 {
		t.Fatalf("sender coordinates lost: %+v", pending.Location)
	}
	delivered := statusEvent(shipment, nil, domain.StatusDelivered)
	if delivered.Location.Lat != 30 || delivered.Location.Lng != 40 {
		t.Fatalf("receiver coordinates lost: %+v", delivered.Location)
	}
	transit := statusEvent(shipment, hub, domain.StatusInTransit)
	if transit.Location.Lat != hub.Coords.Lat || transit.Location.Lng != hub.Coords.Lng {
		t.Fatalf("hub coordinates lost: %+v", transit.Location)
	}
}
