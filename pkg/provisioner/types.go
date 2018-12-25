package provisioner

type Task interface {
	Run() error
}
