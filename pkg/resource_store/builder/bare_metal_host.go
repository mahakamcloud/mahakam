package builder

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	uuid "github.com/satori/go.uuid"

	"github.com/go-openapi/strfmt"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/config"
)

// RoleLabelKey represents key for Label role
const RoleLabelKey = "Role"

// RoleBareMetalHostLabelValue represents role value of bare metal host
const RoleBareMetalHostLabelValue = "bare-metal-host"

// BareMetalHostBuilder is wrapper of BareMetalHost model
type BareMetalHostBuilder struct {
	resource *models.BareMetalHost
}

// Build BareMetalHost resource
func (b *BareMetalHostBuilder) Build(name, kind, owner, role string) ResourceBuilder {
	label := &models.Label{
		Key:   RoleLabelKey,
		Value: role,
	}

	b.resource = &models.BareMetalHost{
		BaseResource: models.BaseResource{
			Name:   swag.String(name),
			Kind:   kind,
			Owner:  owner,
			Labels: []*models.Label{label},
		},
	}

	return b
}

// BuildWithModel generates BareMetalHost resource from swagger generated model
func (b *BareMetalHostBuilder) BuildWithModel(bmhost *models.BareMetalHost) *models.BareMetalHost {
	b.resource = bmhost

	if b.resource.Kind == "" {
		b.resource.Kind = string(KindBareMetalHost)
	}

	if b.resource.Owner == "" {
		b.resource.Owner = config.ResourceOwnerMahakam
	}

	if len(b.resource.Labels) == 0 {
		b.resource.Labels = []*models.Label{
			&models.Label{
				Key:   RoleLabelKey,
				Value: RoleBareMetalHostLabelValue,
			},
		}
	}

	return b.resource
}

// BuildWithMetadata generates BareMetalHost resource with metadata to persist
func (b *BareMetalHostBuilder) BuildWithMetadata(name, kind, owner, role string) ResourceBuilder {
	return b.Build(name, kind, owner, role).BuildMetadata()
}

// BuildKey generates key for a resource
func (b *BareMetalHostBuilder) BuildKey(optKeys ...string) string {
	keys := strings.Join(optKeys, "/")
	return fmt.Sprintf("%s/%s/%s/%s", b.resource.Kind, b.resource.Owner, swag.StringValue(b.resource.Name), keys)
}

// BuildMetadata returns a resource
func (b *BareMetalHostBuilder) BuildMetadata() ResourceBuilder {
	if b.resource.ID == "" {
		b.resource.ID = strfmt.UUID(uuid.NewV4().String())
	}

	now := time.Now()
	b.resource.CreatedAt = strfmt.DateTime(now)
	b.resource.ModifiedAt = strfmt.DateTime(now)

	return b
}

// ToJSON returns a resource
func (b *BareMetalHostBuilder) ToJSON() ([]byte, error) {
	return json.Marshal(b.resource)
}

// GetID return ID of BareMetalHostBuilder resource
func (b *BareMetalHostBuilder) GetID() string {
	return string(b.resource.ID)
}

// Validate BareMetalHost resource
func (b *BareMetalHostBuilder) Validate() error {
	if b.resource.Owner == "" {
		return fmt.Errorf("resource owner not found")
	}
	if swag.StringValue(b.resource.Name) == "" {
		return fmt.Errorf("resource name not found")
	}
	if b.resource.Kind == "" {
		return fmt.Errorf("resource kind not found")
	}

	for _, l := range b.resource.Labels {
		if l.Key == RoleLabelKey && l.Value == "" {
			return fmt.Errorf("resource labels are empty")
		}
	}

	return nil
}
