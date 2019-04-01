package resourcestore

import (
	"fmt"
	"strings"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/model"
	"github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
)

// StorageBackend represents type of storage for kv store backend
type StorageBackend string

const (
	BackendConsul StorageBackend = "consul"
)

// Labels are filterable metadata as key pairs
type Labels map[string]string

func buildKey(rk resource.ResourceKind, owner string, id ...string) string {
	sub := strings.Join(id, "/")
	return fmt.Sprintf("%s/%s/%s/", rk, owner, sub)
}

// ResourceStore is docker libkv wrapper
type ResourceStore interface {
	AddV1(builder model.ResourceBuilder) (id string, err error)
	ListV1(owner string, kind model.ResourceKind, list model.ResourceBuilderList) error

	// TODO(vjdhama): deprecate after moving to builder pattern with
	// swagger generated model completely
	Add(resource resource.Resource) (id string, err error)
	Get(resource resource.Resource) error
	List(owner string, kind resource.ResourceKind, list resource.ResourceList) error
	Update(resource resource.Resource) (revision int64, err error)
	Delete(owner string, id string, resource resource.Resource) error

	// Quick helper hack for interacting with kvstore by given path
	AddFromPath(path string, resource resource.Resource) (id string, err error)
	GetFromPath(path string, resource resource.Resource) error
	ListFromPath(path string, filter Filter, resources resource.ResourceList) error
	ListKeysFromPath(path string) (keys []string, err error)
	UpdateFromPath(path string, resource resource.Resource) (revision int64, err error)
	DeleteFromPath(path string) error

	KeyExists(path string) bool
}

// New creates resource store backed by choice of storage backend type
func New(c config.StorageBackendConfig) (ResourceStore, error) {
	switch c.BackendType {
	case string(BackendConsul):
		kv, err := newConsulKVStore(c)
		if err != nil {
			return nil, fmt.Errorf("create resource store with consul error: %s", err)
		}
		return NewKVResourceStore(kv), nil
	default:
		return nil, fmt.Errorf("create resource store error: %s not supported", c.BackendType)
	}
}
