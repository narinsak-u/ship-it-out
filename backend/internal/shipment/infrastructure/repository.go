package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/internal/shipment/application"
	"github.com/narinsak-u/backend/internal/shipment/domain"
	"gorm.io/gorm"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) List(ctx context.Context, filter domain.Filter) ([]domain.Shipment, int64, error) {
	query := r.db.WithContext(ctx).Model(&shipmentRecord{})
	if filter.Status != "" && filter.Status != "all" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ExcludeStatus != "" {
		query = query.Where("status != ?", filter.ExcludeStatus)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("order_id ILIKE ? OR tracking_number ILIKE ? OR customer_name ILIKE ? OR destination ILIKE ?", like, like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit
	var rows []shipmentRecord
	if err := query.Offset(offset).Limit(filter.Limit).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	shipments := make([]domain.Shipment, len(rows))
	for i, row := range rows {
		shipments[i] = shipmentDomain(row)
	}
	return shipments, total, nil
}

func (r *GormRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Shipment, error) {
	var row shipmentRecord
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&row).Error; err != nil {
		return nil, recordError(err)
	}
	shipment := shipmentDomain(row)
	return &shipment, nil
}

func (r *GormRepository) CreateWithEvent(ctx context.Context, shipment *domain.Shipment, event *domain.ShipmentEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		record := shipmentRecordFromDomain(*shipment)
		if err := allocateIDs(tx, &record); err != nil {
			return err
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		*shipment = shipmentDomain(record)
		event.ShipmentID = record.ID
		eventRecord := eventRecordFromDomain(*event)
		if err := tx.Create(&eventRecord).Error; err != nil {
			return err
		}
		event.ID, event.CreatedAt = eventRecord.ID, eventRecord.CreatedAt
		return nil
	})
}

func (r *GormRepository) SaveWithEvent(ctx context.Context, shipment *domain.Shipment, event *domain.ShipmentEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		record := shipmentRecordFromDomain(*shipment)
		if err := tx.Save(&record).Error; err != nil {
			return err
		}
		if event == nil {
			return nil
		}
		eventRecord := eventRecordFromDomain(*event)
		if err := tx.Create(&eventRecord).Error; err != nil {
			return err
		}
		event.ID, event.CreatedAt = eventRecord.ID, eventRecord.CreatedAt
		return nil
	})
}

func (r *GormRepository) DeleteWithEvents(ctx context.Context, shipment *domain.Shipment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("shipment_id = ?", shipment.ID).Delete(&eventRecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&shipmentRecord{}, shipment.ID).Error
	})
}

func allocateIDs(tx *gorm.DB, shipment *shipmentRecord) error {
	if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", int64(30030)).Error; err != nil {
		return err
	}
	var orderID int64
	if err := tx.Model(&shipmentRecord{}).Select("COALESCE(MAX(CAST(SUBSTRING(order_id FROM 5) AS BIGINT)), 10245)").Where("order_id ~ ?", `^ORD-[0-9]+$`).Scan(&orderID).Error; err != nil {
		return err
	}
	var trackingID int64
	year := time.Now().Year()
	if err := tx.Model(&shipmentRecord{}).Select("COALESCE(MAX(CAST(SUBSTRING(tracking_number FROM 7) AS INTEGER)), -1)").Where("tracking_number ~ ? AND SUBSTRING(tracking_number FROM 3 FOR 4) = ?", `^TH[0-9]{9}$`, fmt.Sprint(year)).Scan(&trackingID).Error; err != nil {
		return err
	}
	if trackingID >= 99999 {
		return errors.New("tracking number space exhausted for current year")
	}
	shipment.OrderID = fmt.Sprintf("ORD-%d", orderID+1)
	shipment.TrackingNumber = fmt.Sprintf("TH%d%05d", year, trackingID+1)
	return nil
}

func (r *GormRepository) FindByTrackingNumber(ctx context.Context, trackingNumber string) (*models.Shipment, error) {
	var row shipmentRecord
	if err := r.db.WithContext(ctx).Where("tracking_number = ?", trackingNumber).First(&row).Error; err != nil {
		return nil, recordError(err)
	}
	shipment := modelsShipment(row)
	return &shipment, nil
}

func (r *GormRepository) FindEventsByShipmentID(ctx context.Context, shipmentID uint) ([]models.ShipmentEvent, error) {
	var rows []eventRecord
	if err := r.db.WithContext(ctx).Where("shipment_id = ?", shipmentID).Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	events := make([]models.ShipmentEvent, len(rows))
	for i, row := range rows {
		events[i] = models.ShipmentEvent{
			ID: row.ID, ShipmentID: row.ShipmentID, Status: row.Status,
			Location: models.Location{Name: row.LocationName, Lat: row.LocationLat, Lng: row.LocationLng},
			Description: row.Description, CreatedAt: row.CreatedAt,
		}
	}
	return events, nil
}

func recordError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return application.ErrNotFound
	}
	return err
}

func shipmentDomain(row shipmentRecord) domain.Shipment {
	return domain.Shipment{
		ID: row.ID, OrderID: row.OrderID, TrackingNumber: row.TrackingNumber,
		Customer: domain.Contact{
			Name: row.CustomerName, Zipcode: row.CustomerZipcode,
			SubDistrict: row.CustomerSubDistrict, District: row.CustomerDistrict, Province: row.CustomerProvince,
			Coords: domain.GeoPoint{Lat: row.CustomerLat, Lng: row.CustomerLng},
		},
		Receiver: domain.Contact{
			Name: row.ReceiverName, Zipcode: row.ReceiverZipcode,
			SubDistrict: row.ReceiverSubDistrict, District: row.ReceiverDistrict, Province: row.ReceiverProvince,
			Coords: domain.GeoPoint{Lat: row.ReceiverLat, Lng: row.ReceiverLng},
		},
		CurrentCoords: domain.GeoPoint{Lat: row.CurrentLat, Lng: row.CurrentLng},
		Origin: row.Origin, Destination: row.Destination, Status: domain.Status(row.Status),
		Carrier: row.Carrier, HubID: row.HubID, Weight: row.Weight, Items: row.Items,
		EstimatedDelivery: row.EstimatedDelivery, CreatedAt: row.CreatedAt, Progress: row.Progress,
	}
}

func shipmentRecordFromDomain(shipment domain.Shipment) shipmentRecord {
	return shipmentRecord{
		ID: shipment.ID, OrderID: shipment.OrderID, TrackingNumber: shipment.TrackingNumber,
		CustomerName: shipment.Customer.Name, CustomerZipcode: shipment.Customer.Zipcode,
		CustomerSubDistrict: shipment.Customer.SubDistrict, CustomerDistrict: shipment.Customer.District,
		CustomerProvince: shipment.Customer.Province, CustomerLat: shipment.Customer.Coords.Lat, CustomerLng: shipment.Customer.Coords.Lng,
		ReceiverName: shipment.Receiver.Name, ReceiverZipcode: shipment.Receiver.Zipcode,
		ReceiverSubDistrict: shipment.Receiver.SubDistrict, ReceiverDistrict: shipment.Receiver.District,
		ReceiverProvince: shipment.Receiver.Province, ReceiverLat: shipment.Receiver.Coords.Lat, ReceiverLng: shipment.Receiver.Coords.Lng,
		CurrentLat: shipment.CurrentCoords.Lat, CurrentLng: shipment.CurrentCoords.Lng,
		Origin: shipment.Origin, Destination: shipment.Destination, Status: string(shipment.Status),
		Carrier: shipment.Carrier, HubID: shipment.HubID, Weight: shipment.Weight, Items: shipment.Items,
		EstimatedDelivery: shipment.EstimatedDelivery, CreatedAt: shipment.CreatedAt, Progress: shipment.Progress,
	}
}

func eventRecordFromDomain(event domain.ShipmentEvent) eventRecord {
	return eventRecord{
		ID: event.ID, ShipmentID: event.ShipmentID, Status: event.Status,
		LocationName: event.Location.Name, LocationLat: event.Location.Lat, LocationLng: event.Location.Lng,
		Description: event.Description, CreatedAt: event.CreatedAt,
	}
}

func modelsShipment(row shipmentRecord) models.Shipment {
	return models.Shipment{
		ID: row.ID, OrderID: row.OrderID, TrackingNumber: row.TrackingNumber,
		Customer: models.ContactInfo{
			Name: row.CustomerName, Zipcode: row.CustomerZipcode,
			SubDistrict: row.CustomerSubDistrict, District: row.CustomerDistrict, Province: row.CustomerProvince,
			Coords: models.GeoPoint{Lat: row.CustomerLat, Lng: row.CustomerLng},
		},
		CustomerLat: row.CustomerLat, CustomerLng: row.CustomerLng,
		Receiver: models.ContactInfo{
			Name: row.ReceiverName, Zipcode: row.ReceiverZipcode,
			SubDistrict: row.ReceiverSubDistrict, District: row.ReceiverDistrict, Province: row.ReceiverProvince,
			Coords: models.GeoPoint{Lat: row.ReceiverLat, Lng: row.ReceiverLng},
		},
		ReceiverLat: row.ReceiverLat, ReceiverLng: row.ReceiverLng,
		CurrentCoords: models.GeoPoint{Lat: row.CurrentLat, Lng: row.CurrentLng},
		CurrentLat: row.CurrentLat, CurrentLng: row.CurrentLng,
		Origin: row.Origin, Destination: row.Destination, Status: row.Status,
		Carrier: row.Carrier, HubID: row.HubID, Weight: row.Weight, Items: row.Items,
		EstimatedDelivery: row.EstimatedDelivery, CreatedAt: row.CreatedAt, Progress: row.Progress,
	}
}

var _ application.Repository = (*GormRepository)(nil)
