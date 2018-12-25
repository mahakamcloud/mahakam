package network

import (
	"crypto/rand"
	"fmt"
)

func GenerateMacAddress() string {
	buf := make([]byte, 6)
	rand.Read(buf)
	buf[0] = (buf[0] | 2) & 0xfe // Set local bit, ensure unicast address
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}
