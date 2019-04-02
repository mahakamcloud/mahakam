package service

import (
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/repository"
)

type BareMetalHostService struct {
	repo *repository.BareMetalHostRepository
}

func NewBareMetalHostService() (*BareMetalHostService, error) {
	r, err := repository.NewBareMetalHostRepository()
	if err != nil {
		return nil, err
	}
	return &BareMetalHostService{r}, nil
}

func (s *BareMetalHostService) RegisterBareMetalHost(b *models.BareMetalHost) error {
	return s.repo.Put(b)
}

func (s *BareMetalHostService) GetAll() ([]*models.BareMetalHost, error) {
	return s.repo.List()
}
