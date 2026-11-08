package shipment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/narinsak-u/backend/internal/models"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// List returns shipments matching the filter, supports pagination, status/exclude_status filtering,
// and text search across order_id, tracking_number, customer_name, and destination (case-insensitive ILIKE).
// Returns the shipments slice, total matching count, and any error.
func (r *GormRepository) List(ctx context.Context, filter ShipmentFilter) ([]models.Shipment, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Shipment{})
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
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
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
func (r *GormRepository) FindByOrderID(ctx context.Context, orderID string) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&shipment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &shipment, nil
}

// FindByTrackingNumber looks up a single shipment by its TH2026xxxxx tracking number.
// This is the public-facing lookup used on the tracking page (no auth required in the handler).
func (r *GormRepository) FindByTrackingNumber(ctx context.Context, trackingNumber string) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := r.db.WithContext(ctx).Where("tracking_number = ?", trackingNumber).First(&shipment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &shipment, nil
}


func (r *GormRepository) CreateWithEvent(ctx context.Context, shipment *models.Shipment, event *models.ShipmentEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.allocateIDs(tx, shipment); err != nil {
			return err
		}
		if err := tx.Create(shipment).Error; err != nil {
			return err
		}
		event.ShipmentID = shipment.ID
		return tx.Create(event).Error
	})
}

func (r *GormRepository) allocateIDs(tx *gorm.DB, shipment *models.Shipment) error {
	if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", int64(30030)).Error; err != nil {
		return err
	}
	var orderID int64
	if err := tx.Model(&models.Shipment{}).Select("COALESCE(MAX(CAST(SUBSTRING(order_id FROM 5) AS BIGINT)), 10245)").Where("order_id ~ ?", `^ORD-[0-9]+$`).Scan(&orderID).Error; err != nil {
		return err
	}
	var trackingID int64
	year := time.Now().Year()
	if err := tx.Model(&models.Shipment{}).Select("COALESCE(MAX(CAST(SUBSTRING(tracking_number FROM 7) AS INTEGER)), -1)").Where("tracking_number ~ ? AND SUBSTRING(tracking_number FROM 3 FOR 4) = ?", `^TH[0-9]{9}$`, fmt.Sprint(year)).Scan(&trackingID).Error; err != nil {
		return err
	}
	if trackingID >= 99999 {
		return errors.New("tracking number space exhausted for current year")
	}
	shipment.OrderID = fmt.Sprintf("ORD-%d", orderID+1)
	shipment.TrackingNumber = fmt.Sprintf("TH%d%05d", year, trackingID+1)
	return nil
}

func (r *GormRepository) SaveWithEvent(ctx context.Context, shipment *models.Shipment, event *models.ShipmentEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(shipment).Error; err != nil {
			return err
		}
		if event != nil {
			return tx.Create(event).Error
		}
		return nil
	})
}

func (r *GormRepository) DeleteWithEvents(ctx context.Context, shipment *models.Shipment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("shipment_id = ?", shipment.ID).Delete(&models.ShipmentEvent{}).Error; err != nil {
			return err
		}
		return tx.Delete(shipment).Error
	})
}

// FindEventsByShipmentID returns all tracking events for a given shipment, ordered oldest-first. Returns empty slice (not nil) if none exist.
func (r *GormRepository) FindEventsByShipmentID(ctx context.Context, shipmentID uint) ([]models.ShipmentEvent, error) {
	var events []models.ShipmentEvent
	if err := r.db.WithContext(ctx).Where("shipment_id = ?", shipmentID).Order("created_at asc").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}


// IDs are allocated while holding a PostgreSQL transaction-scoped advisory lock.
