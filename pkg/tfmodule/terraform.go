package tfmodule

import (
	"fmt"

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

func (t *Terraform) Init(tfDir string) (string, error) {
	args := []string{"init", tfDir}
	return t.runner.CombinedOutput("terraform", args...)
}

func (t *Terraform) Plan(options string, tfDir string) (string, error) {
	args := []string{"plan", options, tfDir}
	fmt.Println(args)
	return t.runner.CombinedOutput("terraform", args...)
}

func (t *Terraform) Apply(options string, tfDir string) (string, error) {
	args := []string{"apply", options, tfDir}
	fmt.Println(args)
	return t.runner.CombinedOutput("terraform", args...)
}

func (t *Terraform) Destroy(tfDir string) (string, error) {
	args := []string{"destroy", tfDir}
	fmt.Println(args)
	return t.runner.CombinedOutput("terraform", args...)
}
