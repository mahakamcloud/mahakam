package model

// ResourceKind represents stored resource kind
type ResourceKind string

const (
	KindBareMetalHost ResourceKind = "bare-metal-host"
)

// ResourceBuilder interface
type ResourceBuilder interface {
	ToJSON() ([]byte, error)
	FromJSON([]byte) error
	BuildKey(optKeys ...string) string
	AddMetadata() ResourceBuilder
	GetID() string
	Validate() error
}

// ResourceBuilderList
type ResourceBuilderList interface {
	ResourceBuilder() ResourceBuilder
	WithItems(items []ResourceBuilder)
}
