package resourcestore

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Status represents current state of task or entity
type Status string

const (
	StatusPending  Status = "Pending"
	StatusCreating Status = "Creating"
	StatusReady    Status = "Ready"
)

// ResourceKind represents stored resource kind
type ResourceKind string

// Labels are filterable metadata as key pairs
type Labels map[string]string

// Resource is base interface for all stored resources or objects
type Resource interface{}

// BaseResource is the base struct for all stored resources or objects
type BaseResource struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Owner        string    `json:"owner"`
	CreatedTime  time.Time `json:"createdTime,omitempty"`
	ModifiedTime time.Time `json:"modifiedTimeomitempty"`
	Revision     uint64    `json:"revision"`
	Status       Status    `json:"status"`
	Labels       Labels    `json:"labels"`
}

func buildKey(rk ResourceKind, owner string, id ...string) string {
	sub := strings.Join(id, "/")
	return fmt.Sprintf("%s/%s/%s", rk, owner, sub)
}

// ResourceStore is docker libkv wrapper
type ResourceStore interface {
	Add(ctx context.Context, resource Resource) (id string, err error)
	Get(ctx context.Context, owner string, key string, resource Resource) error
	List(ctx context.Context, owner string, resources interface{}) error
	Update(ctx context.Context, resource Resource) (revision int64, err error)
	Delete(ctx context.Context, owner string, id string, resource Resource) error
}

// StorageBackendConfig stores metadata for storage backend that we use
type StorageBackendConfig struct {
	BackendType string
	Address     string
	Username    string
	Password    string
	Bucket      string
}

// New creates resource store backed by choice of storage backend type
func New(config StorageBackendConfig) (ResourceStore, error) {
	switch config.BackendType {
	case "postgres":
		p, err := NewPostgresResourceStore(config)
		if err != nil {
			return nil, fmt.Errorf("Create resource store with postgres error: %s", err)
		}
		return p, nil
	default:
		return nil, fmt.Errorf("Create resource store error: %s not supported", config.BackendType)
	}
}
