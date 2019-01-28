package utils

import (
	"net"
	"strconv"

	"github.com/mahakamcloud/mahakam/pkg/cmd_runner"
)

type IPAssigner interface {
	Assign(ip net.IP, mask net.IPMask, netif string) (string, error)
}

type IPUtil struct {
	runner cmd_runner.CmdRunner
}

func NewIPUtil() *IPUtil {
	runner := cmd_runner.New()
	return &IPUtil{
		runner: runner,
	}
}

func (i *IPUtil) Assign(ip net.IP, mask net.IPMask, netif string) (string, error) {
	ones, _ := mask.Size()
	ipaddr := ip.String() + "/" + strconv.Itoa(ones)
	args := []string{"addr", "add", ipaddr, "dev", netif}
	return i.runner.CombinedOutput("ip", args...)
}
