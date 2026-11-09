// Package application contains Hub use cases and persistence ports.
package application

import (
	"context"

	"github.com/narinsak-u/backend/internal/hub/domain"
)

type Repository interface {
	FindAll(context.Context) ([]domain.Hub, error)
	FindByID(context.Context, string) (*domain.Hub, error)
	Create(context.Context, *domain.Hub) error
	Save(context.Context, *domain.Hub) error
	Delete(context.Context, string) error
}
