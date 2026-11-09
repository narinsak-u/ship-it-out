package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/narinsak-u/backend/internal/hub/application"
	"github.com/narinsak-u/backend/internal/hub/domain"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) FindAll(ctx context.Context) ([]domain.Hub, error) {
	var records []Record
	if err := r.db.WithContext(ctx).Find(&records).Error; err != nil {
		return nil, err
	}
	hubs := make([]domain.Hub, len(records))
	for i := range records {
		hubs[i] = records[i].hub()
	}
	return hubs, nil
}

func (r *GormRepository) FindByID(ctx context.Context, id string) (*domain.Hub, error) {
	var record Record
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, application.ErrNotFound }
		return nil, err
	}
	hub := record.hub()
	return &hub, nil
}

func (r *GormRepository) Create(ctx context.Context, hub *domain.Hub) error {
	if hub.ID == "" {
		id, err := r.generateHubID(ctx)
		if err != nil {
			return err
		}
		hub.ID = id
	}
	record := recordFromHub(*hub)
	return r.db.WithContext(ctx).Create(&record).Error
}

func (r *GormRepository) Save(ctx context.Context, hub *domain.Hub) error {
	record := recordFromHub(*hub)
	return r.db.WithContext(ctx).Save(&record).Error
}

func (r *GormRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Record{}, "id = ?", id).Error
}

func (r *GormRepository) generateHubID(ctx context.Context) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var records []Record
	if err := r.db.WithContext(ctx).Select("id").Find(&records).Error; err != nil {
		return "", err
	}
	maxNum := 0
	for _, record := range records {
		parts := strings.SplitN(record.ID, "-", 2)
		if len(parts) == 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil && n > maxNum {
				maxNum = n
			}
		}
	}
	return fmt.Sprintf("HUB-%03d", maxNum+1), nil
}
