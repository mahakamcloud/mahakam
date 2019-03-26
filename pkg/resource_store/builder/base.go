package builder

// ResourceBuilder interface
type ResourceBuilder interface {
	Build(name, kind, owner, role string) ResourceBuilder
	BuildKey(optKeys ...string) string
	BuildMetadata() ResourceBuilder
	BuildWithMetadata(name, kind, owner, role string) ResourceBuilder
	GetID() string
	Validate() error
}
