package builder

// ResourceKind represents stored resource kind
type ResourceKind string

const (
	KindBareMetalHost ResourceKind = "bare-metal-host"
)

// ResourceBuilder interface
type ResourceBuilder interface {
	Build(name, kind, owner, role string) ResourceBuilder
	ToJSON() ([]byte, error)
	FromJSON([]byte) error
	BuildKey(optKeys ...string) string
	BuildMetadata() ResourceBuilder
	BuildWithMetadata(name, kind, owner, role string) ResourceBuilder
	GetID() string
	Validate() error
}

// ResourceBuilderList
type ResourceBuilderList interface {
	ResourceBuilder() ResourceBuilder
	WithItems(items []ResourceBuilder)
}
