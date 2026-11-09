package presentation

import (
	"time"

	"github.com/narinsak-u/backend/internal/hub/domain"
)

type HubRequest struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	CarrierID          string         `json:"carrierId"`
	Address            string         `json:"address"`
	Coords             domain.GeoPoint `json:"coords"`
	Capacity           int            `json:"capacity"`
	CurrentUtilization int            `json:"currentUtilization"`
	Status             string         `json:"status"`
	CreatedAt          time.Time      `json:"createdAt"`
}

func (r HubRequest) hub() *domain.Hub {
	return &domain.Hub{ID: r.ID, Name: r.Name, CarrierID: r.CarrierID, Address: r.Address, Coords: r.Coords, Capacity: r.Capacity, CurrentUtilization: r.CurrentUtilization, Status: r.Status, CreatedAt: r.CreatedAt}
}
