package models

import "time"

type ShipmentEvent struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ShipmentID  uint      `gorm:"not null;index" json:"shipment_id"`
	Status      string    `gorm:"not null" json:"status"`
	Location    string    `json:"location,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
