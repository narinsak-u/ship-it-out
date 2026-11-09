package infrastructure

import (
	"time"

	"github.com/narinsak-u/backend/internal/hub/domain"
)

type Record struct {
	ID                 string    `gorm:"primaryKey"`
	Name               string    `gorm:"not null"`
	CarrierID          string    `gorm:"not null;index"`
	Address            string    `gorm:"not null"`
	Lat                float64   `gorm:"column:lat"`
	Lng                float64   `gorm:"column:lng"`
	Capacity           int       `gorm:"not null"`
	CurrentUtilization int       `gorm:"not null"`
	Status             string    `gorm:"index;not null;default:active"`
	CreatedAt          time.Time
}

func (Record) TableName() string { return "hubs" }

func recordFromHub(h domain.Hub) Record {
	return Record{ID: h.ID, Name: h.Name, CarrierID: h.CarrierID, Address: h.Address, Lat: h.Coords.Lat, Lng: h.Coords.Lng, Capacity: h.Capacity, CurrentUtilization: h.CurrentUtilization, Status: h.Status, CreatedAt: h.CreatedAt}
}

func (r Record) hub() domain.Hub {
	return domain.Hub{ID: r.ID, Name: r.Name, CarrierID: r.CarrierID, Address: r.Address, Lat: r.Lat, Lng: r.Lng, Coords: domain.GeoPoint{Lat: r.Lat, Lng: r.Lng}, Capacity: r.Capacity, CurrentUtilization: r.CurrentUtilization, Status: r.Status, CreatedAt: r.CreatedAt}
}

