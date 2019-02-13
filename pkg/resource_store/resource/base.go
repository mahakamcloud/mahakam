package resource

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Status represents current state of task or entity
type Status string

// ResourceKind represents stored resource kind
type ResourceKind string

// Labels represents filterable keypairs metadata
type Labels []Label

type Label struct {
	Key   string
	Value string
}

const (
	StatusPending  Status = "Pending"
	StatusCreating Status = "Creating"
	StatusReady    Status = "Ready"

	KindCluster          ResourceKind = "cluster"
	KindTerraform        ResourceKind = "terraform"
	KindTerraformBackend ResourceKind = "terraform backend"
	KindTask             ResourceKind = "task"
	KindNetwork          ResourceKind = "network"
	KindNode             ResourceKind = "node"
	KindIPPool           ResourceKind = "ippool"
)

// Resource is interface for all stored resources or objects
type Resource interface {
	GetResource() *BaseResource
	BuildResource() Resource
	UpdateResource() Resource
	BuildKey(optKeys ...string) string
	BuildChildKey(parentKey, key string) string
	PreCheck() error
	GetLabels() Labels
}

// ResourceList represents list of resources
type ResourceList interface {
	Resource() Resource
	WithItems(items []Resource)
}

// BaseResource is the base struct for all stored resources or objects
type BaseResource struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Kind         string    `json:"kind"`
	Owner        string    `json:"owner"`
	CreatedTime  time.Time `json:"createdTime,omitempty"`
	ModifiedTime time.Time `json:"modifiedTime,omitempty"`
	Revision     uint64    `json:"revision"`
	Status       Status    `json:"status"`
	Labels       Labels    `json:"labels"`
}

func (br *BaseResource) GetResource() *BaseResource {
	return br
}

func (br *BaseResource) BuildResource() Resource {
	br.ID = uuid.NewV4().String()

	now := time.Now()
	br.CreatedTime = now
	br.ModifiedTime = now

	return br
}

func (br *BaseResource) UpdateResource() Resource {
	now := time.Now()
	br.ModifiedTime = now

	return br
}

func (br *BaseResource) PreCheck() error {
	if br.Owner == "" {
		return fmt.Errorf("resource owner attribute cannot be empty")
	}

	var validName = regexp.MustCompile(`[\w.-]`)
	if !validName.MatchString(br.Name) {
		return fmt.Errorf("resource name %s is invalid", br.Name)
	}
	return nil
}

func (br *BaseResource) BuildKey(optKeys ...string) string {
	keys := strings.Join(optKeys, "/")
	return fmt.Sprintf("%s/%s/%s/%s", br.Kind, br.Owner, br.Name, keys)
}

// BuildChildKey returns stringified key from parent key
func (br *BaseResource) BuildChildKey(parentKey string, key string) string {
	return fmt.Sprintf("%s/%s/", parentKey, key)
}

// GetLabels is getter for Labels attribute
func (br *BaseResource) GetLabels() Labels {
	return br.Labels
}

// NewResourceFromKind returns empty resource of that type
func NewResourceFromKind(resKind ResourceKind) (Resource, error) {
	switch resKind {
	case KindCluster:
		return &ResourceCluster{}, nil
	case KindIPPool:
		return &ResourceIPPool{}, nil
	case KindNetwork:
		return &ResourceNetwork{}, nil
	case KindNode:
		return &ResourceNode{}, nil
	case KindTerraform:
		return &ResourceTerraform{}, nil
	default:
		return nil, fmt.Errorf("invalid resourceKind : %s", resKind)
	}
}

// NewResourceListFromKind returns empty resource of that type
func NewResourceListFromKind(resKind ResourceKind) (ResourceList, error) {
	switch resKind {
	case KindCluster:
		return &ResourceClusterList{}, nil
	case KindIPPool:
		return &ResourceIPPoolList{}, nil
	default:
		return nil, fmt.Errorf("invalid resource kind: %s", resKind)
	}
}
