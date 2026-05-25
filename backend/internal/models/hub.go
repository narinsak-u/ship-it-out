package models

import (
	"time"

	"gorm.io/gorm"
)

// Hub maps to the frontend's Hub type — a warehouse or中转中心 where
// shipments are sorted and transferred between carriers.
type Hub struct {
	ID                 string    `gorm:"primaryKey" json:"id"`
	Name               string    `gorm:"not null" json:"name"`
	CarrierID          string    `gorm:"not null;index" json:"carrierId"`
	Address            string    `gorm:"not null" json:"address"`
	Lat                float64   `gorm:"column:lat" json:"-"`
	Lng                float64   `gorm:"column:lng" json:"-"`
	Coords             GeoPoint  `gorm:"-" json:"coords"`
	Capacity           int       `gorm:"not null" json:"capacity"`
	CurrentUtilization int       `gorm:"not null" json:"currentUtilization"`
	Status             string    `gorm:"not null;default:active" json:"status"`
	CreatedAt          time.Time `json:"createdAt"`
}

func (h *Hub) BeforeSave(_ *gorm.DB) error {
	h.Lat = h.Coords.Lat
	h.Lng = h.Coords.Lng
	return nil
}

func (h *Hub) AfterFind(_ *gorm.DB) error {
	h.Coords = GeoPoint{Lat: h.Lat, Lng: h.Lng}
	return nil
}
