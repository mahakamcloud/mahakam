package terraform

import (
	"github.com/mahakamcloud/mahakam/pkg/cmd_runner"
)

type Terraform struct {
	runner cmd_runner.CmdRunner
}

func NewWithCmdRunner(runner cmd_runner.CmdRunner) *Terraform {
	return &Terraform{
		runner: runner,
	}
}

func New() *Terraform {
	return NewWithCmdRunner(cmd_runner.New())
}

func (t *Terraform) Init() (string, error) {
	args := []string{"init"}
	return t.runner.CombinedOutput("terraform", args...)
}

func (t *Terraform) Plan() (string, error) {
	args := []string{"plan"}
	return t.runner.CombinedOutput("terraform", args...)
}

func (t *Terraform) Apply() (string, error) {
	args := []string{"apply", "-autoapprove"}
	return t.runner.CombinedOutput("terraform", args...)
}

func (t *Terraform) Destroy() (string, error) {
	args := []string{"destroy"}
	return t.runner.CombinedOutput("terraform", args...)
}
