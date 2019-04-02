package kvstore

import (
	"fmt"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/mahakamcloud/mahakam/pkg/config"
)

// StorageBackend represents type of storage for kv store backend
type StorageBackend string

const (
	BackendConsul StorageBackend = "consul"
)

// New creates resource store backed by choice of storage backend type
func New() (*KVStore, error) {
	c := config.GetConfig().KVStoreConfig

	switch c.BackendType {
	case string(BackendConsul):
		kv, err := newConsulKVStore(c)
		if err != nil {
			return nil, fmt.Errorf("create resource store with consul error: %s", err)
		}

		return NewKVStore(kv), nil
	default:
		return nil, fmt.Errorf("create resource store error: %s not supported", c.BackendType)
	}
}

func newConsulKVStore(c config.StorageBackendConfig) (store.Store, error) {
	consul.Register()

	return libkv.NewStore(
		store.Backend(c.BackendType),
		[]string{c.Address},
		&store.Config{
			Bucket:            c.Bucket,
			ConnectionTimeout: 1 * time.Second,
			PersistConnection: true,
			Username:          c.Username,
			Password:          c.Password,
		},
	)
}
