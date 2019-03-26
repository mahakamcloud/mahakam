package builder

// ResourceBuilder interface
type ResourceBuilder interface {
	Build(name, kind, owner, role string) ResourceBuilder
	BuildKey(optKeys ...string) (string, error)
	BuildMetadata() ResourceBuilder
	BuildWithMetadata(name, kind, owner, role string) ResourceBuilder
}
