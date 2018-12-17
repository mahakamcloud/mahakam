package resourcestore

import (
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/mahakamcloud/mahakam/pkg/config"
)

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
