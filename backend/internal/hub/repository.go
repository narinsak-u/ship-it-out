package hub

import "github.com/narinsak-u/backend/internal/models"

type Repository interface {
	FindAll() ([]models.Hub, error)
	FindByID(id string) (*models.Hub, error)
	Create(hub *models.Hub) error
	Save(hub *models.Hub) error
	Delete(id string) error
}
