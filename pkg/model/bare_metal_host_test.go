package model

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/stretchr/testify/assert"
)

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
		b := &BareMetalHostWrapper{}
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
