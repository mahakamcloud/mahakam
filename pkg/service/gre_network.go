package service

import (
	"fmt"

	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/netdclient"
	"github.com/mahakamcloud/mahakam/pkg/repository"
)

// DefaultCIDRRange defines the default network CIDR for new networks
const DefaultCIDRRange = "10.40.0.0/16"

// GreNetworkService is service for GreNetworkRepository
type GreNetworkService struct {
	repo *repository.GreNetworkRepository
}

// NewGreNetworkService creates a new GreNetworkService
func NewGreNetworkService() (*GreNetworkService, error) {
	r, err := repository.NewGreNetworkRepository()
	if err != nil {
		return nil, err
	}
	return &GreNetworkService{r}, nil
}

// CreateGreNetwork interacts with netd to create networks
func (s *GreNetworkService) CreateGreNetwork(g *models.GreNetwork) error {
	greKey, err := s.generateGreKey()
	if err != nil {
		return fmt.Errorf("error generating GREKey : %s", err)
	}
	g.GREKey = greKey
	// TODO(vjdhama) : Use network allocator
	g.CIDR = DefaultCIDRRange

	bmHostService, err := NewBareMetalHostService()
	if err != nil {
		return err
	}

	bmhosts, err := bmHostService.GetAll()
	if err != nil {
		return err
	}

	netdc := netdclient.Client{}

	// TODO(vjdhama) : Save responses on KV store
	netdc.CreateNetwork(bmhosts, g)

	return s.repo.Put(g)
}

// GetAll lists all GreNetworks
func (s *GreNetworkService) GetAll() ([]*models.GreNetwork, error) {
	return s.repo.List()
}

func (s *GreNetworkService) generateGreKey() (int64, error) {
	networks, err := s.repo.List()
	if err != nil {
		return 0, err
	}

	return int64(len(networks) + 1), nil
}
