package models

import (
	"time"

	"gorm.io/gorm"
)

// GeoPoint matches the frontend's GeoPoint — a simple lat/lng coordinate pair.
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// ContactInfo matches the frontend's ContactInfo — reusable for both customer
// and receiver. Coords is stored in dedicated columns on Shipment and
// populated via GORM hooks (BeforeSave / AfterFind).
type ContactInfo struct {
	Name        string   `json:"name"`
	Zipcode     string   `json:"zipcode"`
	SubDistrict string   `json:"subDistrict"`
	District    string   `json:"district"`
	Province    string   `json:"province"`
	Coords      GeoPoint `gorm:"-" json:"coords"`
}

// Shipment maps to the frontend's Order type.
type Shipment struct {
	ID                uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	TrackingNumber    string      `gorm:"unique;index;not null" json:"trackingNumber"`
	Customer          ContactInfo `gorm:"embedded;embeddedPrefix:customer_" json:"customer"`
	CustomerLat       float64     `gorm:"column:customer_lat" json:"-"`
	CustomerLng       float64     `gorm:"column:customer_lng" json:"-"`
	Receiver          ContactInfo `gorm:"embedded;embeddedPrefix:receiver_" json:"receiver"`
	ReceiverLat       float64     `gorm:"column:receiver_lat" json:"-"`
	ReceiverLng       float64     `gorm:"column:receiver_lng" json:"-"`
	CurrentCoords     GeoPoint    `gorm:"-" json:"currentCoords"`
	CurrentLat        float64     `gorm:"column:current_lat" json:"-"`
	CurrentLng        float64     `gorm:"column:current_lng" json:"-"`
	Origin            string      `gorm:"not null" json:"origin"`
	Destination       string      `gorm:"not null" json:"destination"`
	Status            string      `gorm:"not null;default:pending" json:"status"`
	Carrier           string      `json:"carrier"`
	DriverID          string      `json:"driverId,omitempty"`
	Weight            string      `json:"weight"`
	Items             int         `json:"items"`
	EstimatedDelivery time.Time   `json:"estimatedDelivery"`
	CreatedAt         time.Time   `json:"createdAt"`
	Progress          float64     `json:"progress"`
}

// BeforeSave copies nested Coords into flat DB columns before insert/update.
func (s *Shipment) BeforeSave(_ *gorm.DB) error {
	s.CustomerLat = s.Customer.Coords.Lat
	s.CustomerLng = s.Customer.Coords.Lng
	s.ReceiverLat = s.Receiver.Coords.Lat
	s.ReceiverLng = s.Receiver.Coords.Lng
	s.CurrentLat = s.CurrentCoords.Lat
	s.CurrentLng = s.CurrentCoords.Lng
	return nil
}

// AfterFind reconstructs nested Coords from flat DB columns after a query.
func (s *Shipment) AfterFind(_ *gorm.DB) error {
	s.Customer.Coords = GeoPoint{Lat: s.CustomerLat, Lng: s.CustomerLng}
	s.Receiver.Coords = GeoPoint{Lat: s.ReceiverLat, Lng: s.ReceiverLng}
	s.CurrentCoords = GeoPoint{Lat: s.CurrentLat, Lng: s.CurrentLng}
	return nil
}
