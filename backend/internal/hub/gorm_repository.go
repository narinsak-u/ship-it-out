package hub

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) FindAll(ctx context.Context) ([]Hub, error) {
	var records []Record
	if err := r.db.WithContext(ctx).Find(&records).Error; err != nil { return nil, err }
	hubs := make([]Hub, len(records))
	for i := range records { hubs[i] = records[i].hub() }
	return hubs, nil
}

func (r *GormRepository) FindByID(ctx context.Context, id string) (*Hub, error) {
	var record Record
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrNotFound }
		return nil, err
	}
	hub := record.hub()
	return &hub, nil
}

func (r *GormRepository) Create(ctx context.Context, hub *Hub) error {
	if hub.ID == "" { hub.ID = r.generateHubID(ctx) }
	return r.db.WithContext(ctx).Create(recordFromHub(*hub)).Error
}

func (r *GormRepository) Save(ctx context.Context, hub *Hub) error {
	return r.db.WithContext(ctx).Save(recordFromHub(*hub)).Error
}

func (r *GormRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Record{}, "id = ?", id).Error
}

func (r *GormRepository) generateHubID(ctx context.Context) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var records []Record
	r.db.WithContext(ctx).Select("id").Find(&records)
	maxNum := 0
	for _, record := range records {
		parts := strings.SplitN(record.ID, "-", 2)
		if len(parts) == 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil && n > maxNum { maxNum = n }
		}
	}
	return fmt.Sprintf("HUB-%03d", maxNum+1)
}
