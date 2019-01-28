package filter

type Verb string

const (
	// VerbEqual evaluates equality
	VerbEqual Verb = "equal"
)

type Filter interface {
	Add(...FilterStatement) Filter
}

type FilterStatement struct {
	Verb   Verb
	Object interface{}
}
