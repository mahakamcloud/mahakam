package resourcestore_test

import (
	"net"
	"testing"

	. "github.com/mahakamcloud/mahakam/pkg/resource_store"
	"github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/stretchr/testify/assert"
)

func TestApplyFilter(t *testing.T) {
	_, fakenet, _ := net.ParseCIDR("1.1.1.1/28")

	testcases := []struct {
		name           string
		filterLabels   resource.Labels
		resourceLabels resource.Labels
		expectedMatch  bool
	}{
		{
			name:           "Filter with one label should match resource with one label",
			filterLabels:   resource.Labels{resource.Label{Key: "foo", Value: "foo"}},
			resourceLabels: resource.Labels{resource.Label{Key: "foo", Value: "foo"}},
			expectedMatch:  true,
		},
		{
			name:           "Filter with one label should match resource with multiple labels",
			filterLabels:   resource.Labels{resource.Label{Key: "foo", Value: "foo"}},
			resourceLabels: resource.Labels{resource.Label{Key: "foo", Value: "foo"}, resource.Label{Key: "bar", Value: "bar"}},
			expectedMatch:  true,
		},
		{
			name:           "Filter with one label should not match resource with one label",
			filterLabels:   resource.Labels{resource.Label{Key: "foo", Value: "foo"}},
			resourceLabels: resource.Labels{resource.Label{Key: "foo", Value: "bar"}},
			expectedMatch:  false,
		},
		{
			name:           "Filter with one label should not match resource with multiple labels",
			filterLabels:   resource.Labels{resource.Label{Key: "foo", Value: "foo"}},
			resourceLabels: resource.Labels{resource.Label{Key: "foo", Value: "bar"}, resource.Label{Key: "bar", Value: "foo"}},
			expectedMatch:  false,
		},
		{
			name:           "Filter with one label should not match resource with no label",
			filterLabels:   resource.Labels{resource.Label{Key: "foo", Value: "foo"}},
			resourceLabels: resource.Labels{},
			expectedMatch:  false,
		},
		{
			name:           "Filter with multiple labels should match resource with multiple labels",
			filterLabels:   resource.Labels{resource.Label{Key: "foo", Value: "bar"}, resource.Label{Key: "bar", Value: "foo"}},
			resourceLabels: resource.Labels{resource.Label{Key: "foo", Value: "bar"}, resource.Label{Key: "bar", Value: "foo"}},
			expectedMatch:  true,
		},
		{
			name:           "Filter with multiple labels should not match resource with multiple labels",
			filterLabels:   resource.Labels{resource.Label{Key: "foo", Value: "foo"}, resource.Label{Key: "bar", Value: "bar"}},
			resourceLabels: resource.Labels{resource.Label{Key: "foo", Value: "bar"}, resource.Label{Key: "bar", Value: "foo"}},
			expectedMatch:  false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			filter := &FilterScope{}
			for _, l := range testcase.filterLabels {
				filter.Add(l)
			}
			fakeResource := resource.NewResourceIPPool(*fakenet)
			fakeResource.WithLabels(testcase.resourceLabels)

			match := ApplyFilter(filter, fakeResource)
			assert.Equal(t, testcase.expectedMatch, match)
		})
	}
}
