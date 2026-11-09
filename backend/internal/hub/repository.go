// Package hub provides the hub service, HTTP handler, and focused persistence contract.
package hub

import (
	"context"

	"github.com/narinsak-u/backend/internal/models"
)

// Repository defines the data-access contract used by hub workflows.
type Repository interface {
	FindAll(context.Context) ([]models.Hub, error)
	FindByID(context.Context, string) (*models.Hub, error)
	Create(context.Context, *models.Hub) error
	Save(context.Context, *models.Hub) error
	Delete(context.Context, string) error
}
