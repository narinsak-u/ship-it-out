package hub

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/narinsak-u/backend/internal/models"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindAll() ([]models.Hub, error) {
	var hubs []models.Hub
	if err := r.db.Find(&hubs).Error; err != nil {
		return nil, err
	}
	return hubs, nil
}

func (r *GormRepository) FindByID(id string) (*models.Hub, error) {
	var hub models.Hub
	if err := r.db.First(&hub, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &hub, nil
}

func (r *GormRepository) Create(hub *models.Hub) error {
	if hub.ID == "" {
		hub.ID = r.generateHubID()
	}
	return r.db.Create(hub).Error
}

func (r *GormRepository) Save(hub *models.Hub) error {
	return r.db.Save(hub).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&models.Hub{}, "id = ?", id).Error
}

func (r *GormRepository) generateHubID() string {
	var hubs []models.Hub
	r.db.Select("id").Find(&hubs)
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
