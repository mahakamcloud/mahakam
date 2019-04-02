package resourcestore

import (
	"github.com/mahakamcloud/mahakam/pkg/repository"
)

type fakeResource struct{}

func NewFakeResource() repository.ResourceBuilder {
	return &fakeResource{}
}

func (f *fakeResource) Build(name, owner, kind, role string) repository.ResourceBuilder {
	return f
}

func (f *fakeResource) BuildWithMetadata(name, owner, kind, role string) repository.ResourceBuilder {
	return f
}

func (f *fakeResource) Validate() error {
	return nil
}

func (f *fakeResource) BuildKey(opts ...string) string {
	return "fake-resource-key"
}

func (f *fakeResource) AddMetadata() repository.ResourceBuilder {
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
