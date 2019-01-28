package resourcestore

import (
	"github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
)

type Filter interface {
	Add(...resource.Label) Filter
	Labels() resource.Labels
}

type FilterScope struct {
	labels resource.Labels
}

func (f *FilterScope) Add(labels ...resource.Label) Filter {
	for _, l := range labels {
		f.labels = append(f.labels, l)
	}
	return f
}

func (f *FilterScope) Labels() resource.Labels {
	return f.labels
}

func ApplyFilter(filter Filter, resource resource.Resource) bool {
	nextFilter := false
	for i, s := range filter.Labels() {
		for _, l := range resource.GetLabels() {
			if matchFilter(s, l) && i == len(filter.Labels())-1 {
				return true
			}
			if matchFilter(s, l) {
				nextFilter = true
				break
			}
		}
		if nextFilter {
			continue
		}
		return false
	}
	return false
}

func matchFilter(s, l resource.Label) bool {
	return s.Key == l.Key && s.Value == l.Value
}
