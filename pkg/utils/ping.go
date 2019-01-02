package utils

import (
	"fmt"
	"net"
	"time"

	"github.com/sparrc/go-ping"

	log "github.com/sirupsen/logrus"
)

const (
	ICMPDefaultPingCount = 1
)

// PortPing is used to check whether specific port in a network
// is opened or not. Nil means success.
func PortPing(addressWithPort string, timeout time.Duration, log log.FieldLogger) error {
	log.Infof("pinging telnet node %s", addressWithPort)

	const network = "tcp"
	conn, err := net.DialTimeout(network, addressWithPort, timeout)
	if conn != nil {
		defer conn.Close()
	}
	return err
}

// PingNWithDelay calls Ping N number of times with delay in between
func PingNWithDelay(pingFn func(string, time.Duration, log.FieldLogger) error,
	count int, delay time.Duration) bool {

	for i := 0; i < count; i++ {
		if err := pingFn; err == nil {
			return true
		}
		time.Sleep(delay)
	}
	return false
}

// ICMPPing is used to check if node is pinging
func ICMPPing(address string, timeout time.Duration, log log.FieldLogger) error {
	log.Infof("pinging icmp node %s", address)

	pinger, err := ping.NewPinger(address)
	if err != nil {
		return err
	}

	pinger.Count = ICMPDefaultPingCount
	pinger.Timeout = timeout

	// blocking run
	pinger.Run()

	if pinger.Statistics().PacketsRecv <= 0 {
		return fmt.Errorf("not able to ping %s", address)
	}

	return nil
}
