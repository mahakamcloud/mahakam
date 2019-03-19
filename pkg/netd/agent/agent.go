package agent

type Agent interface {
	Run()
	Execute() error
}
