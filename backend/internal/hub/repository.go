// Package hub provides HTTP handlers and a GORM-backed repository for logistics hub management.
package hub

import "github.com/narinsak-u/backend/internal/models"

// Repository defines the data-access contract for hubs. The handler depends on this interface,
// making it straightforward to swap in a mock for unit tests.
type Repository interface {
	// FindAll returns every hub row. Used by GET /api/hubs (public list endpoint).
	FindAll() ([]models.Hub, error)
	// FindByID looks up a hub by its string PK (e.g. "HUB-001"). Returns error if not found.
	FindByID(id string) (*models.Hub, error)
	// Create inserts a new hub. Auto-generates a "HUB-NNN" ID if hub.ID is empty.
	Create(hub *models.Hub) error
	// Save writes all fields of an existing hub back to Postgres (full update).
	Save(hub *models.Hub) error
	// Delete removes a hub by its string ID.
	Delete(id string) error
}
