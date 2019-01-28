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
type Labels map[string]string

const (
	StatusPending  Status = "Pending"
	StatusCreating Status = "Creating"
	StatusReady    Status = "Ready"

	KindCluster          ResourceKind = "cluster"
	KindTerraformBackend ResourceKind = "terraform backend"
	KindTask             ResourceKind = "task"
	KindNetwork          ResourceKind = "network"
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
}

// ResourceList represents list of resources
type ResourceList interface {
	Resource() Resource
	SetItems(items []Resource)
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
		return fmt.Errorf("BaseResource owner attribute cannot be empty")
	}

	var validName = regexp.MustCompile(`[\w.-]`)
	if !validName.MatchString(br.Name) {
		return fmt.Errorf("BaseResource name %s is invalid", br.Name)
	}
	return nil
}

func (br *BaseResource) BuildKey(optKeys ...string) string {
	keys := strings.Join(optKeys, "/")
	return fmt.Sprintf("%s/%s/%s/%s", br.Kind, br.Owner, br.Name, keys)
}

func (br *BaseResource) BuildChildKey(parentKey string, key string) string {
	return fmt.Sprintf("%s/%s/", parentKey, key)
}
