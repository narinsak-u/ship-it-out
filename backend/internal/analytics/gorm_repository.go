package analytics

import (
	"github.com/narinsak-u/backend/internal/models"
	"gorm.io/gorm"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Shipment{}).Count(&count).Error
	return count, err
}

func (r *GormRepository) CountActive() (int64, error) {
	var count int64
	err := r.db.Model(&models.Shipment{}).Where("status NOT IN ?", []string{"delivered", "returned"}).Count(&count).Error
	return count, err
}

func (r *GormRepository) CountByStatus() ([]StatusCount, error) {
	var results []StatusCount
	err := r.db.Model(&models.Shipment{}).Select("status, count(*) as count").Group("status").Scan(&results).Error
	return results, err
}

func (r *GormRepository) CountByProvince() ([]ProvinceCount, error) {
	var results []ProvinceCount
	err := r.db.Model(&models.Shipment{}).Select("receiver_province as province, count(*) as total").Group("receiver_province").Scan(&results).Error
	return results, err
}

func (r *GormRepository) CountByMonth() ([]MonthCount, error) {
	var results []MonthCount
	err := r.db.Model(&models.Shipment{}).Select("to_char(created_at, 'YYYY-MM') as month, count(*) as count").Group("month").Order("month").Scan(&results).Error
	return results, err
}

func (r *GormRepository) CountByDayOfWeek() ([]DayCount, error) {
	var results []DayCount
	err := r.db.Model(&models.Shipment{}).Select("trim(to_char(created_at, 'Day')) as day, count(*) as count").Group("day").Order("min(extract(dow from created_at))").Scan(&results).Error
	return results, err
}
