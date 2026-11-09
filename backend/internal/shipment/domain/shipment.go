package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound   = errors.New("shipment not found")
	ErrInvalidHub = errors.New("invalid hub ID")
)

type Status string

const (
	StatusPending        Status = "pending"
	StatusPickedUp       Status = "picked_up"
	StatusDeparted       Status = "departed"
	StatusInTransit      Status = "in_transit"
	StatusOutForDelivery Status = "out_for_delivery"
	StatusDelivered      Status = "delivered"
	StatusDelayed        Status = "delayed"
)

type GeoPoint struct {
	Lat float64
	Lng float64
}

type Contact struct {
	Name        string
	Zipcode     string
	SubDistrict string
	District    string
	Province    string
	Coords      GeoPoint
}

type Location struct {
	Name string
	Lat  float64
	Lng  float64
}

type ShipmentEvent struct {
	ID          uint
	ShipmentID  uint
	Status      string
	Location    Location
	Description string
	CreatedAt   time.Time
}

type Shipment struct {
	ID                uint
	OrderID           string
	TrackingNumber    string
	Customer          Contact
	Receiver          Contact
	CurrentCoords     GeoPoint
	Origin            string
	Destination       string
	Status            Status
	Carrier           string
	HubID             string
	Weight            float64
	Items             int
	EstimatedDelivery time.Time
	CreatedAt         time.Time
	Progress          float64
}

type Filter struct {
	Page          int
	Limit         int
	Search        string
	Status        string
	ExcludeStatus string
}
