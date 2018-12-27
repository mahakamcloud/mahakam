package utils

import (
	"fmt"

	"github.com/mahakamcloud/mahakam/pkg/cmd_runner"
)

type SCPConfig struct {
	Username        string
	RemoteIPAddress string
	PrivateKeyPath  string
	RemoteFilePath  string
	LocalFilePath   string
}

type SCPClient struct {
	runner cmd_runner.CmdRunner
}

func NewSCPClient() *SCPClient {
	runner := cmd_runner.New()
	return &SCPClient{
		runner: runner,
	}
}

func (s *SCPClient) CopyRemoteFile(config SCPConfig) (string, error) {
	conn := fmt.Sprintf("%s@%s:%s", config.Username, config.RemoteIPAddress, config.RemoteFilePath)
	args := []string{conn, config.LocalFilePath}
	return s.runner.CombinedOutput("scp", args...)
}
