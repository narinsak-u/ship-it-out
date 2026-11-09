package domain

import "time"

type GeoPoint struct {
	Lat float64
	Lng float64
}

type Hub struct {
	ID                 string
	Name               string
	CarrierID          string
	Address            string
	Coords             GeoPoint
	Capacity           int
	CurrentUtilization int
	Status             string
	CreatedAt          time.Time
}
