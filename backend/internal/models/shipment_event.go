package models

import "time"

// Location matches the frontend's Location type — every tracking event records
// which hub or place it happened, along with its geographic coordinates.
type Location struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

// ShipmentEvent tracks a single status change or checkpoint for a shipment.
// It maps to the frontend's TrackingEvent type.
type ShipmentEvent struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ShipmentID  uint      `gorm:"not null;index" json:"shipmentId"`
	Status      string    `gorm:"not null" json:"status"`
	Location    Location  `gorm:"embedded;embeddedPrefix:location_" json:"location"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"timestamp"`
}
