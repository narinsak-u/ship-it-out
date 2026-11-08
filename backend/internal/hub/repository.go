// Package hub provides the hub service, HTTP handler, and focused persistence contract.
package hub

import "github.com/narinsak-u/backend/internal/models"

// Repository defines the data-access contract used by hub workflows.
type Repository interface {
	FindAll() ([]models.Hub, error)
	FindByID(id string) (*models.Hub, error)
	Create(hub *models.Hub) error
	Save(hub *models.Hub) error
	Delete(id string) error
}
