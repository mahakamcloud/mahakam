package util

import (
	"net"
	"strings"
)

func LastOctet(ip net.IP) string {
	octets := strings.Split(ip.String(), ".")
	return octets[len(octets)-1]
}
