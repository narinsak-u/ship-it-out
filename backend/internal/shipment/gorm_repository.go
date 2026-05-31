package shipment

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/narinsak-u/backend/internal/models"
	"gorm.io/gorm"
)

type GormRepository struct {
	db          *gorm.DB
	orderIDLock sync.Mutex
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// List returns shipments matching the filter, supports pagination, status/exclude_status filtering,
// and text search across order_id, tracking_number, customer_name, and destination (case-insensitive ILIKE).
// Returns the shipments slice, total matching count, and any error.
func (r *GormRepository) List(filter ShipmentFilter) ([]models.Shipment, int64, error) {
	query := r.db.Model(&models.Shipment{})
	if filter.Status != "" && filter.Status != "all" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ExcludeStatus != "" {
		query = query.Where("status != ?", filter.ExcludeStatus)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where(
			"order_id ILIKE ? OR tracking_number ILIKE ? OR customer_name ILIKE ? OR destination ILIKE ?",
			like, like, like, like,
		)
	}
	var total int64
	query.Count(&total)
	if filter.Page < 1 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit
	var shipments []models.Shipment
	if err := query.Offset(offset).Limit(filter.Limit).Order("created_at DESC").Find(&shipments).Error; err != nil {
		return nil, 0, err
	}
	return shipments, total, nil
}

// FindByOrderID looks up a single shipment by its ORD-xxxxx order ID. Returns nil and an error
// (which the caller should map to 404) if no shipment matches.
func (r *GormRepository) FindByOrderID(orderID string) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := r.db.Where("order_id = ?", orderID).First(&shipment).Error; err != nil {
		return nil, err
	}
	return &shipment, nil
}

// FindByTrackingNumber looks up a single shipment by its TH2026xxxxx tracking number.
// This is the public-facing lookup used on the tracking page (no auth required in the handler).
func (r *GormRepository) FindByTrackingNumber(trackingNumber string) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := r.db.Where("tracking_number = ?", trackingNumber).First(&shipment).Error; err != nil {
		return nil, err
	}
	return &shipment, nil
}

// Create auto-generates an OrderID (ORD-xxxxx) and TrackingNumber (TH2026xxxxx), then inserts
// the shipment row into Postgres. The caller is responsible for populating all other fields.
func (r *GormRepository) Create(shipment *models.Shipment) error {
	shipment.OrderID = r.generateOrderID()
	shipment.TrackingNumber = r.generateTrackingNumber()
	return r.db.Create(shipment).Error
}

// Save writes all fields of the shipment back to Postgres (full update, not partial).
// GORM hooks like BeforeSave will fire, which is how Coords get synced to the flat lat/lng columns.
func (r *GormRepository) Save(shipment *models.Shipment) error {
	return r.db.Save(shipment).Error
}

// Delete removes the shipment row from Postgres. Does NOT cascade-delete events —
// call DeleteShipmentEvents separately if you need to clean up the event log.
func (r *GormRepository) Delete(shipment *models.Shipment) error {
	return r.db.Delete(shipment).Error
}

// CreateEvent inserts a tracking event (e.g. "Picked Up", "In Transit") linked to a shipment.
// The caller must set ShipmentID, Status, Description, and Location before calling this.
func (r *GormRepository) CreateEvent(event *models.ShipmentEvent) error {
	return r.db.Create(event).Error
}

// FindEventsByShipmentID returns all tracking events for a given shipment, ordered oldest-first. Returns empty slice (not nil) if none exist.
func (r *GormRepository) FindEventsByShipmentID(shipmentID uint) ([]models.ShipmentEvent, error) {
	var events []models.ShipmentEvent
	if err := r.db.Where("shipment_id = ?", shipmentID).Order("created_at asc").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// DeleteShipmentEvents deletes ALL tracking events belonging to a shipment. Used when deleting
// a shipment itself (events must be removed first to avoid FK constraint violations).
func (r *GormRepository) DeleteShipmentEvents(shipmentID uint) error {
	return r.db.Where("shipment_id = ?", shipmentID).Delete(&models.ShipmentEvent{}).Error
}

// Count returns the total number of shipments in the database. Used by the analytics overview endpoint.
func (r *GormRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Shipment{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountActive returns the number of shipments whose status is NOT "delivered" or "returned"
// (i.e. shipments still in progress). Used by the analytics overview endpoint.
func (r *GormRepository) CountActive() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Shipment{}).Where("status NOT IN ?", []string{"delivered", "returned"}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByStatus groups all shipments by their status field and returns a count per status.
// Example result: [{status: "pending", count: 2}, {status: "in_transit", count: 4}].
func (r *GormRepository) CountByStatus() ([]StatusCountResult, error) {
	var results []StatusCountResult
	if err := r.db.Model(&models.Shipment{}).Select("status, count(*) as count").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// CountByMonth groups shipments by their creation month (YYYY-MM format) and returns a count
// per month, sorted chronologically. Used by the analytics timeseries endpoint.
func (r *GormRepository) CountByMonth() ([]MonthCountResult, error) {
	var results []MonthCountResult
	if err := r.db.Model(&models.Shipment{}).
		Select("to_char(created_at, 'YYYY-MM') as month, count(*) as count").
		Group("month").Order("month").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// CountByDayOfWeek groups shipments by the day-of-week they were created (Monday, Tuesday, etc.)
// and returns a count per day, ordered Monday-first. Used by the analytics timeseries endpoint.
func (r *GormRepository) CountByDayOfWeek() ([]DayCountResult, error) {
	var results []DayCountResult
	if err := r.db.Model(&models.Shipment{}).
		Select("trim(to_char(created_at, 'Day')) as day, count(*) as count").
		Group("day").Order("min(extract(dow from created_at))").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// CountByProvince groups shipments by the receiver's province and returns a count per province.
// Used by the analytics overview endpoint to derive region-level breakdowns.
func (r *GormRepository) CountByProvince() ([]ProvinceCountResult, error) {
	var results []ProvinceCountResult
	if err := r.db.Model(&models.Shipment{}).
		Select("receiver_province as province, count(*) as total").
		Group("receiver_province").
		Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// generateTrackingNumber creates a unique tracking ID in the format TH<year><5-digit-microsecond-hash>.
// Example: "TH202612345". Uniqueness is probabilistic (based on microsecond timestamp).
func (r *GormRepository) generateTrackingNumber() string {
	return fmt.Sprintf("TH%d%05d", time.Now().Year(), time.Now().UnixMicro()%100000)
}

// generateOrderID creates the next OrderID by finding the highest existing ORD-xxxxx suffix and
// incrementing by one (e.g. ORD-10245 → ORD-10246). Uses a mutex to prevent duplicates under
// concurrent requests. Falls back to ORD-10246 if the DB is empty or the format doesn't match.
func (r *GormRepository) generateOrderID() string {
	r.orderIDLock.Lock()
	defer r.orderIDLock.Unlock()

	var max models.Shipment
	if err := r.db.Select("order_id").Order("order_id DESC").First(&max).Error; err != nil {
		return "ORD-10246"
	}

	parts := strings.SplitN(max.OrderID, "-", 2)
	if len(parts) == 2 {
		if n, err := strconv.Atoi(parts[1]); err == nil {
			return fmt.Sprintf("ORD-%d", n+1)
		}
	}
	return "ORD-10246"
}
