package infrastructure

import (
	"context"

	"github.com/narinsak-u/backend/internal/analytics/application"
	"gorm.io/gorm"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("shipments").Count(&count).Error
	return count, err
}

func (r *GormRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("shipments").Where("status NOT IN ?", []string{"delivered", "returned"}).Count(&count).Error
	return count, err
}

func (r *GormRepository) CountByStatus(ctx context.Context) ([]application.StatusCount, error) {
	var results []application.StatusCount
	err := r.db.WithContext(ctx).Table("shipments").Select("status, count(*) as count").Group("status").Scan(&results).Error
	return results, err
}

func (r *GormRepository) CountByProvince(ctx context.Context) ([]application.ProvinceCount, error) {
	var results []application.ProvinceCount
	err := r.db.WithContext(ctx).Table("shipments").Select("receiver_province as province, count(*) as total").Group("receiver_province").Scan(&results).Error
	return results, err
}

func (r *GormRepository) CountByMonth(ctx context.Context) ([]application.MonthCount, error) {
	var results []application.MonthCount
	err := r.db.WithContext(ctx).Table("shipments").Select("to_char(created_at, 'YYYY-MM') as month, count(*) as count").Group("month").Order("month").Scan(&results).Error
	return results, err
}

func (r *GormRepository) CountByDayOfWeek(ctx context.Context) ([]application.DayCount, error) {
	var results []application.DayCount
	err := r.db.WithContext(ctx).Table("shipments").Select("trim(to_char(created_at, 'Day')) as day, count(*) as count").Group("day").Order("min(extract(dow from created_at))").Scan(&results).Error
	return results, err
}
