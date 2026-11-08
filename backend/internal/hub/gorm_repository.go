package hub

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/narinsak-u/backend/internal/models"
	"gorm.io/gorm"
)

// GormRepository implements Repository using GORM backed by Postgres.
// Create via NewGormRepository — the zero value is not usable.
type GormRepository struct {
	db *gorm.DB
	mu sync.Mutex // guards ID generation under concurrent Create calls
}

// NewGormRepository opens a GORM-backed hub repository with the given *gorm.DB connection.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// FindAll returns every hub in the database with no filtering or pagination.
func (r *GormRepository) FindAll(ctx context.Context) ([]models.Hub, error) {
	var hubs []models.Hub
	if err := r.db.WithContext(ctx).Find(&hubs).Error; err != nil {
		return nil, err
	}
	return hubs, nil
}

// FindByID looks up a single hub by its string primary key (e.g. "HUB-001").
// Returns ErrNotFound if no hub matches.
func (r *GormRepository) FindByID(ctx context.Context, id string) (*models.Hub, error) {
	var hub models.Hub
	if err := r.db.WithContext(ctx).First(&hub, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &hub, nil
}

// Create inserts a new hub row. If the hub's ID is empty, it auto-generates one by finding
// the highest existing "HUB-NNN" suffix and incrementing by one (e.g. HUB-005 → HUB-006).
func (r *GormRepository) Create(ctx context.Context, hub *models.Hub) error {
	if hub.ID == "" {
		hub.ID = r.generateHubID(ctx)
	}
	return r.db.WithContext(ctx).Create(hub).Error
}

// Save writes all fields of an existing hub back to Postgres (full update, not partial).
// GORM's BeforeSave hook fires here, syncing Coords to the flat Lat/Lng columns.
func (r *GormRepository) Save(ctx context.Context, hub *models.Hub) error {
	return r.db.WithContext(ctx).Save(hub).Error
}

// Delete removes a hub row by its string ID. Does not return an error if the ID doesn't exist
// (GORM's Delete is idempotent — it just affects 0 rows).
func (r *GormRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Hub{}, "id = ?", id).Error
}

// generateHubID scans all existing hub IDs, finds the highest numeric suffix, and returns
// "HUB-NNN" with the next number (zero-padded to 3 digits). Uses a mutex to prevent
// duplicate IDs under concurrent requests.
func (r *GormRepository) generateHubID(ctx context.Context) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	var hubs []models.Hub
	r.db.WithContext(ctx).Select("id").Find(&hubs)
	maxNum := 0
	for _, h := range hubs {
		parts := strings.SplitN(h.ID, "-", 2)
		if len(parts) == 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil && n > maxNum {
				maxNum = n
			}
		}
	}
	return fmt.Sprintf("HUB-%03d", maxNum+1)
}
