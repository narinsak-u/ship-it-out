package hub

import (
	"errors"

	"github.com/narinsak-u/backend/internal/models"
)

var ErrNotFound = errors.New("hub not found")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]models.Hub, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByID(id string) (*models.Hub, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(hub *models.Hub) error {
	return s.repo.Create(hub)
}

func (s *Service) Update(id string, hub *models.Hub) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}
	hub.ID = id
	return s.repo.Save(hub)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
