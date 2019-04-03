package service

import (
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/repository"
)

type GreNetworkService struct {
	repo *repository.GreNetworkRepository
}

func NewGreNetworkService() (*GreNetworkService, error) {
	r, err := repository.NewGreNetworkRepository()
	if err != nil {
		return nil, err
	}
	return &GreNetworkService{r}, nil
}

func (s *GreNetworkService) CreateGreNetwork(b *models.GreNetwork) error {
	return s.repo.Put(b)
}

func (s *GreNetworkService) GetAll() ([]*models.GreNetwork, error) {
	return s.repo.List()
}
