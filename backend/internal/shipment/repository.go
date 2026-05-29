package shipment

import "github.com/narinsak-u/backend/internal/models"

type ShipmentFilter struct {
	Page          int
	Limit         int
	Search        string
	Status        string
	ExcludeStatus string
}

type StatusCountResult struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type MonthCountResult struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type DayCountResult struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}

type Repository interface {
	List(filter ShipmentFilter) ([]models.Shipment, int64, error)
	FindByOrderID(orderID string) (*models.Shipment, error)
	FindByTrackingNumber(trackingNumber string) (*models.Shipment, error)
	Create(shipment *models.Shipment) error
	Save(shipment *models.Shipment) error
	Delete(shipment *models.Shipment) error
	CreateEvent(event *models.ShipmentEvent) error
	FindEventsByShipmentID(shipmentID uint) ([]models.ShipmentEvent, error)
	DeleteShipmentEvents(shipmentID uint) error
	Count() (int64, error)
	CountActive() (int64, error)
	CountByStatus() ([]StatusCountResult, error)
	CountByMonth() ([]MonthCountResult, error)
	CountByDayOfWeek() ([]DayCountResult, error)
}
