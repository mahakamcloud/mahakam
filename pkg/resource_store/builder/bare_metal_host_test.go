package builder

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
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
		{"fake-name", "", "", "fake-role"},
	}

	for _, test := range tests {
		b := &BareMetalHostBuilder{}
		b.Build(test.name, test.kind, test.owner, test.role)

		assert.Equal(t, test.name, swag.StringValue(b.resource.Name))
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
		b.Build("fake-name", "", "", "")
		b.BuildMetadata()

		assert.NotNil(t, test.id)
	}
}

func TestBuildWithModel(t *testing.T) {
	tests := []struct {
		name  string
		kind  string
		owner string
	}{
		{"fake-name", "fake-kind", "fake-owner"},
		{"fake-name", "", ""},
	}

	for _, test := range tests {
		b := &BareMetalHostBuilder{}
		m := &models.BareMetalHost{
			BaseResource: models.BaseResource{
				Name:  swag.String(test.name),
				Kind:  test.kind,
				Owner: test.owner,
			},
		}

		b.BuildWithModel(m)

		if test.kind == "" {
			assert.NotEmpty(t, b.resource.Kind)
		}

		if test.owner == "" {
			assert.NotEmpty(t, b.resource.Owner)
		}

		assert.Equal(t, test.name, swag.StringValue(b.resource.Name))
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
