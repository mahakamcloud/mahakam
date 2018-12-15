package resourcestore

import (
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
)

func newConsulKVStore(config StorageBackendConfig) (store.Store, error) {
	consul.Register()

	return libkv.NewStore(
		store.Backend(config.BackendType),
		[]string{config.Address},
		&store.Config{
			Bucket:            config.Bucket,
			ConnectionTimeout: 1 * time.Second,
			PersistConnection: true,
			Username:          config.Username,
			Password:          config.Password,
		},
	)
}
