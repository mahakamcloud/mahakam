package network

import (
	"crypto/rand"
	"fmt"
)

func GenerateMacAddress() string {
	buf := make([]byte, 6)
	rand.Read(buf)
	buf[0] = (buf[0] | 2) & 0xfe // Set local bit, ensure unicast address
	if buf[0] == 0xfe {
		buf[0] = 0xfa // 0xFE is a reserved first octet in Libvirt, therefore we set the first octet to 0xFA which follows local and unicast address convention
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}
