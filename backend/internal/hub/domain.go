package hub

import "time"

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Hub struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	CarrierID          string    `json:"carrierId"`
	Address            string    `json:"address"`
	Lat                float64   `json:"-"`
	Lng                float64   `json:"-"`
	Coords             GeoPoint  `json:"coords"`
	Capacity           int       `json:"capacity"`
	CurrentUtilization int       `json:"currentUtilization"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"createdAt"`
}
