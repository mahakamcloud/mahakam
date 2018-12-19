package resourcestore

import (
	"fmt"
	"strings"

	"github.com/mahakamcloud/mahakam/pkg/config"
)

// Status represents current state of task or entity
type Status string

// StorageBackend represents type of storage for kv store backend
type StorageBackend string

// ResourceKind represents stored resource kind
type ResourceKind string

const (
	StatusPending  Status = "Pending"
	StatusCreating Status = "Creating"
	StatusReady    Status = "Ready"

	BackendConsul StorageBackend = "consul"

	KindCluster          ResourceKind = "cluster"
	KindTerraformBackend ResourceKind = "terraform backend"
	KindTask             ResourceKind = "task"
)

// Labels are filterable metadata as key pairs
type Labels map[string]string

func buildKey(rk ResourceKind, owner string, id ...string) string {
	sub := strings.Join(id, "/")
	return fmt.Sprintf("%s/%s/%s", rk, owner, sub)
}

// ResourceStore is docker libkv wrapper
type ResourceStore interface {
	Add(resource Resource) (id string, err error)
	Get(resource Resource) error
	List(owner string, resources interface{}) error
	Update(resource Resource) (revision int64, err error)
	Delete(owner string, id string, resource Resource) error
}

// New creates resource store backed by choice of storage backend type
func New(c config.StorageBackendConfig) (ResourceStore, error) {
	switch c.BackendType {
	case string(BackendConsul):
		kv, err := newConsulKVStore(c)
		if err != nil {
			return nil, fmt.Errorf("Create resource store with consul error: %s", err)
		}
		return NewKVResourceStore(kv), nil
	default:
		return nil, fmt.Errorf("Create resource store error: %s not supported", c.BackendType)
	}
}
