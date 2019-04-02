package service

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/repository"
	uuid "github.com/satori/go.uuid"
)

type ResourceKind string

const (
	// RoleLabelKey represents key for Label role
	RoleLabelKey = "Role"

	// RoleBareMetalHostLabelValue represents role value of bare metal host
	RoleBareMetalHostLabelValue = "bare-metal-host"

	KindBareMetalHost ResourceKind = "bare-metal-host"
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
	s.addMandatoryFieldsToBareMetalHost(b)
	return s.repo.Put(b)
}

func (s *BareMetalHostService) GetAll() ([]*models.BareMetalHost, error) {
	return s.repo.List()
}

func (s *BareMetalHostService) addMandatoryFieldsToBareMetalHost(b *models.BareMetalHost) {
	if b.Kind == "" {
		b.Kind = string(KindBareMetalHost)
	}

	if b.Owner == "" {
		b.Owner = config.ResourceOwnerMahakam
	}

	if len(b.Labels) == 0 {
		b.Labels = []*models.Label{
			&models.Label{
				Key:   RoleLabelKey,
				Value: RoleBareMetalHostLabelValue,
			},
		}
	}

	if b.ID == "" {
		b.ID = strfmt.UUID(uuid.NewV4().String())
	}

	now := time.Now()
	b.CreatedAt = strfmt.DateTime(now)
	b.ModifiedAt = strfmt.DateTime(now)
}
