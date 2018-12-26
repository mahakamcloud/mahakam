package cmd_runner

import (
	"os/exec"
)

type CmdRunner interface {
	CombinedOutput(cmd string, args ...string) (output string, err error)
	CombinedOutputWithDir(dir, cmd string, args ...string) (output string, err error)
}

type realCmdRunner struct{}

func New() *realCmdRunner {
	return &realCmdRunner{}
}

func (runner *realCmdRunner) CombinedOutput(cmd string, args ...string) (string, error) {
	output, err := exec.Command(cmd, args...).CombinedOutput()
	return string(output), err
}

func (runner *realCmdRunner) CombinedOutputWithDir(dir, cmd string, args ...string) (string, error) {
	run := exec.Command(cmd, args...)
	run.Dir = dir
	output, err := run.CombinedOutput()
	return string(output), err
}
