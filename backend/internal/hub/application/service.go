package application

import (
	"context"
	"errors"

	"github.com/narinsak-u/backend/internal/hub/domain"
)

var ErrNotFound = errors.New("hub not found")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) List(ctx context.Context) ([]domain.Hub, error) { return s.repo.FindAll(ctx) }
func (s *Service) GetByID(ctx context.Context, id string) (*domain.Hub, error) { return s.repo.FindByID(ctx, id) }
func (s *Service) Create(ctx context.Context, hub *domain.Hub) error { return s.repo.Create(ctx, hub) }
func (s *Service) Update(ctx context.Context, id string, hub *domain.Hub) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil { return err }
	hub.ID = id
	return s.repo.Save(ctx, hub)
}
func (s *Service) Delete(ctx context.Context, id string) error { return s.repo.Delete(ctx, id) }
