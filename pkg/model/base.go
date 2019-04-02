package model

// ResourceKind represents stored resource kind
type ResourceKind string

const (
	KindBareMetalHost ResourceKind = "bare-metal-host"
)

// ResourceBuilder interface
type ResourceBuilder interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	BuildKey(optKeys ...string) string
	AddMetadata() ResourceBuilder
	ID() string
}

// ResourceBuilderList
type ResourceBuilderList interface {
	ResourceBuilder() ResourceBuilder
	WithItems(items []ResourceBuilder)
}
