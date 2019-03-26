package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		name  string
		kind  string
		owner string
		role  string
	}{
		{"fake-name", "fake-kind", "fake-owner", "fake-role"},
		{"", "fake-kind", "fake-owner", "fake-role"},
	}

	for _, test := range tests {
		b := &BareMetalHostBuilder{}
		b.Build(test.name, test.kind, test.owner, test.role)

		assert.Equal(t, test.name, b.resource.Name)
		assert.Equal(t, test.kind, b.resource.Kind)
		assert.Equal(t, test.owner, b.resource.Owner)

		for _, l := range b.resource.Labels {
			assert.Equal(t, test.role, l.Value)
		}
	}
}

func TestBuildMetadata(t *testing.T) {
	tests := []struct {
		id string
	}{
		{"fake-id"},
		{""},
	}

	for _, test := range tests {
		b := &BareMetalHostBuilder{}
		b.BuildMetadata()

		assert.NotNil(t, test.id)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		kind  string
		owner string
		role  string

		expectedError bool
	}{
		{"fake-name", "fake-kind", "fake-owner", "fake-role", false},
		{"fake-name", "fake-kind", "fake-owner", "", true},
		{"", "fake-kind", "fake-owner", "fake-role", true},
	}

	for _, test := range tests {
		b := &BareMetalHostBuilder{}
		b.Build(test.name, test.kind, test.owner, test.role)

		err := b.Validate()
		if test.expectedError {
			assert.Error(t, err)
		}
		if !test.expectedError {
			assert.NoError(t, err)
		}
	}
}
