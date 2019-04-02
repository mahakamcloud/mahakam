package resourcestore

import (
	"fmt"
	"testing"

	"github.com/docker/libkv/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/mahakamcloud/mahakam/pkg/model"
)

type fakeResource struct{}

func NewFakeResource() model.ResourceBuilder {
	return &fakeResource{}
}

func (f *fakeResource) Build(name, owner, kind, role string) model.ResourceBuilder {
	return f
}

func (f *fakeResource) BuildWithMetadata(name, owner, kind, role string) model.ResourceBuilder {
	return f
}

func (f *fakeResource) Validate() error {
	return nil
}

func (f *fakeResource) BuildKey(opts ...string) string {
	return "fake-resource-key"
}

func (f *fakeResource) AddMetadata() model.ResourceBuilder {
	return f
}

func (f *fakeResource) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (f *fakeResource) UnmarshalJSON(in []byte) error {
	return nil
}

func (f *fakeResource) ID() string {
	return "fake-resource-id"
}

func TestAddV1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := NewMockStore(ctrl)
	kvr := NewKVResourceStore(s)

	b := NewFakeResource()
	s.EXPECT().AtomicPut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, &store.KVPair{}, nil)

	id, err := kvr.AddV1(b)
	assert.Equal(t, "fake-resource-id", id)
	assert.NoError(t, err)

	s.EXPECT().AtomicPut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, &store.KVPair{}, fmt.Errorf("fake-kvstore-error"))
	id, err = kvr.AddV1(b)
	assert.Equal(t, "", id)
	assert.Error(t, err)
}
