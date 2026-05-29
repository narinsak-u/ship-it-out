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

func (r *GormRepository) FindByOrderID(orderID string) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := r.db.Where("order_id = ?", orderID).First(&shipment).Error; err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *GormRepository) FindByTrackingNumber(trackingNumber string) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := r.db.Where("tracking_number = ?", trackingNumber).First(&shipment).Error; err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *GormRepository) Create(shipment *models.Shipment) error {
	shipment.OrderID = r.generateOrderID()
	shipment.TrackingNumber = r.generateTrackingNumber()
	return r.db.Create(shipment).Error
}

func (r *GormRepository) Save(shipment *models.Shipment) error {
	return r.db.Save(shipment).Error
}

func (r *GormRepository) Delete(shipment *models.Shipment) error {
	return r.db.Delete(shipment).Error
}

func (r *GormRepository) CreateEvent(event *models.ShipmentEvent) error {
	return r.db.Create(event).Error
}

func (r *GormRepository) FindEventsByShipmentID(shipmentID uint) ([]models.ShipmentEvent, error) {
	var events []models.ShipmentEvent
	if err := r.db.Where("shipment_id = ?", shipmentID).Order("created_at asc").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *GormRepository) DeleteShipmentEvents(shipmentID uint) error {
	return r.db.Where("shipment_id = ?", shipmentID).Delete(&models.ShipmentEvent{}).Error
}

func (r *GormRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Shipment{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormRepository) CountActive() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Shipment{}).Where("status NOT IN ?", []string{"delivered", "returned"}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormRepository) CountByStatus() ([]StatusCountResult, error) {
	var results []StatusCountResult
	if err := r.db.Model(&models.Shipment{}).Select("status, count(*) as count").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *GormRepository) CountByMonth() ([]MonthCountResult, error) {
	var results []MonthCountResult
	if err := r.db.Model(&models.Shipment{}).
		Select("to_char(created_at, 'YYYY-MM') as month, count(*) as count").
		Group("month").Order("month").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *GormRepository) CountByDayOfWeek() ([]DayCountResult, error) {
	var results []DayCountResult
	if err := r.db.Model(&models.Shipment{}).
		Select("trim(to_char(created_at, 'Day')) as day, count(*) as count").
		Group("day").Order("min(extract(dow from created_at))").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *GormRepository) generateTrackingNumber() string {
	return fmt.Sprintf("TH%d%05d", time.Now().Year(), time.Now().UnixMilli()%100000)
}

func (r *GormRepository) generateOrderID() string {
	r.orderIDLock.Lock()
	defer r.orderIDLock.Unlock()
	var shipments []models.Shipment
	r.db.Select("order_id").Find(&shipments)
	maxNum := 10245
	for _, s := range shipments {
		parts := strings.SplitN(s.OrderID, "-", 2)
		if len(parts) == 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil && n > maxNum {
				maxNum = n
			}
		}
	}
	return fmt.Sprintf("ORD-%d", maxNum+1)
}
