package builder

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/go-openapi/strfmt"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
)

// NodeBuilder is builder object for Node model
type NodeBuilder struct {
	Resource models.Node
}

// GetResourceID returns Resource ID
func (n *NodeBuilder) GetResourceID() string {
	return n.Resource.ID.String()
}

// Validate returns error in case NodeBuilder is invalid
func (n *NodeBuilder) Validate() error {
	if n.Resource.Owner == "" {
		return fmt.Errorf("resource owner attribute cannot be empty")
	}

	var validName = regexp.MustCompile(`[\w.-]`)
	if !validName.MatchString(n.Resource.Name) {
		return fmt.Errorf("resource name %s is invalid", n.Resource.Name)
	}
	return nil
}

// Build returns a resource
func (n *NodeBuilder) Build(name, kind, owner, role string) ResourceBuilder {
	n.Resource = models.Node{
		BaseResource: models.BaseResource{
			Name:  name,
			Owner: owner,
			Kind:  string(n.Resource.Kind),
		},
	}
	return n
}

// BuildMetadata returns a resource
func (n *NodeBuilder) BuildMetadata() ResourceBuilder {
	n.Resource.ID = strfmt.UUID(uuid.NewV4().String())

	now := time.Now()
	n.Resource.CreatedAt = strfmt.DateTime(now)
	n.Resource.ModifiedAt = strfmt.DateTime(now)

	return n
}

// Update updates a existing resource
func (n *NodeBuilder) Update() ResourceBuilder {
	now := time.Now()
	n.Resource.ModifiedAt = strfmt.DateTime(now)

	return n
}

// Update updates a existing resource
func (n *NodeBuilder) UpdateRevision(i uint64) ResourceBuilder {
	n.Resource.Revision = i
	return n
}

// BuildKey returns stringified key
func (n *NodeBuilder) BuildKey(optKeys ...string) (string, error) {
	keys := strings.Join(optKeys, "/")

	if n.Resource.Owner == "" {
		return "", fmt.Errorf("owner parameter is required for getting resource")
	}
	if n.Resource.Name == "" {
		return "", fmt.Errorf("name parameter is required for getting resource")
	}
	if n.Resource.Kind == "" {
		return "", fmt.Errorf("kind parameter is required for getting resource")
	}

	return fmt.Sprintf("%s/%s/%s/%s", n.Resource.Kind, n.Resource.Owner, n.Resource.Name, keys), nil
}

// BuildChildKey returns stringified key from parent key
func (n *NodeBuilder) BuildChildKey(parentKey string, key string) string {
	return fmt.Sprintf("%s/%s/", parentKey, key)
}

// GetLabels is getter for Labels attribute
func (n *NodeBuilder) GetLabels() []*models.BaseResourceLabelsItems0 {
	return n.Resource.Labels
}

// NewResourceBuilderFromKind returns empty resource of that type
func NewResourceBuilderFromKind(resKind Kind) (ResourceBuilder, error) {
	switch resKind {
	case KindNode:
		return &NodeBuilder{}, nil
	default:
		return nil, fmt.Errorf("invalid resourceKind : %s", resKind)
	}
}
