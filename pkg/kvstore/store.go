package kvstore

import (
	"fmt"

	"github.com/docker/libkv/store"
)

type KVStore struct {
	store store.Store
}

// NewKVResourceStore creates key value store for resources
func NewKVStore(s store.Store) *KVStore {
	return &KVStore{
		store: s,
	}
}

func (k *KVStore) Put(key string, value []byte) error {
	opts := &store.WriteOptions{
		IsDir: false,
	}
	_, _, err := k.store.AtomicPut(key, value, nil, opts)
	if err != nil {
		return fmt.Errorf("add kv resource atomic put error: %s", err)
	}

	return nil
}

func (k *KVStore) List(key string) ([][]byte, error) {
	kvpairs, err := k.store.List(key)
	if err != nil {
		return nil, err
	}

	values := make([][]byte, 0)
	for _, kv := range kvpairs {
		values = append(values, kv.Value)
	}
	return values, nil
}
