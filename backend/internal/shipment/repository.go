package shipment

import "github.com/narinsak-u/backend/internal/models"

// ShipmentFilter holds all optional query parameters for listing shipments.
// Page and Limit control pagination (defaults: page=1, limit=10). Use limit=-1 for all rows.
// Search runs a case-insensitive ILIKE across order_id, tracking_number, customer_name, and destination.
// Status and ExcludeStatus filter/suppress by exact status match.
type ShipmentFilter struct {
	Page          int
	Limit         int
	Search        string
	Status        string
	ExcludeStatus string
}

// StatusCountResult represents a single row from GROUP BY status — one per unique status value.
type StatusCountResult struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// MonthCountResult represents a single row from GROUP BY month (YYYY-MM format).
type MonthCountResult struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

// DayCountResult represents a single row from GROUP BY day-of-week (Monday, Tuesday, etc.).
type DayCountResult struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}

// ProvinceCountResult represents a single row from GROUP BY receiver_province.
type ProvinceCountResult struct {
	Province string `json:"province"`
	Total    int64  `json:"total"`
}

// Repository defines the data-access contract for shipments. The handler depends on this
// interface (not on GormRepository directly), which makes it easy to swap implementations
// for testing — just pass a mock that satisfies this interface.
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
	CountByProvince() ([]ProvinceCountResult, error)
}
