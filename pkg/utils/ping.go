package utils

import (
	"fmt"
	"net"
	"time"

	ping "github.com/sparrc/go-ping"

	log "github.com/sirupsen/logrus"
)

const (
	ICMPDefaultPingCount = 1
)

type PingChecker interface {
	PortPingNWithDelay(addressWithPort string, timeout time.Duration, log log.FieldLogger,
		count int, delay time.Duration) bool
	ICMPPingNWithDelay(address string, timeout time.Duration, log log.FieldLogger,
		count int, delay time.Duration) bool
}

type PingCheck struct{}

func NewPingCheck() PingChecker {
	return &PingCheck{}
}

// portPing is used to check whether specific port in a network
// is opened or not. Nil means success.
func (p PingCheck) portPing(addressWithPort string, timeout time.Duration, log log.FieldLogger) error {
	log.Infof("pinging telnet node %s", addressWithPort)

	const network = "tcp"
	conn, err := net.DialTimeout(network, addressWithPort, timeout)
	if conn != nil {
		defer conn.Close()
	}
	return err
}

// PortPingNWithDelay calls PortPing N number of times with delay in between
func (p PingCheck) PortPingNWithDelay(addressWithPort string, timeout time.Duration, log log.FieldLogger,
	count int, delay time.Duration) bool {

	for i := 0; i < count; i++ {
		if err := p.portPing(addressWithPort, timeout, log); err == nil {
			log.Infof("pinging telnet node %s successful", addressWithPort)
			return true
		}
		time.Sleep(delay)
	}
	log.Errorf("pinging telnet node %s timeout", addressWithPort)
	return false
}

// icmpPing is used to check if node is pinging
func (p PingCheck) icmpPing(address string, timeout time.Duration, log log.FieldLogger) error {
	log.Infof("pinging icmp node %s", address)

	pinger, err := ping.NewPinger(address)
	if err != nil {
		return err
	}

	pinger.Count = ICMPDefaultPingCount
	pinger.Timeout = timeout
	pinger.SetPrivileged(true)

	// blocking run
	pinger.Run()

	if pinger.Statistics().PacketsRecv <= 0 {
		return fmt.Errorf("not able to ping %s", address)
	}

	return nil
}

// ICMPPingNWithDelay calls ICMPPing N number of times with delay in between
func (p PingCheck) ICMPPingNWithDelay(address string, timeout time.Duration, log log.FieldLogger,
	count int, delay time.Duration) bool {

	for i := 0; i < count; i++ {
		if err := p.icmpPing(address, timeout, log); err == nil {
			log.Infof("pinging icmp node %s successful", address)
			return true
		}
		time.Sleep(delay)
	}
	log.Errorf("pinging icmp node %s timeout", address)
	return false
}
