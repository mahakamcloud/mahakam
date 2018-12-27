package utils

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// Ping is used to check whether specific port in a network
// is opened or not. Nil means success.
func Ping(network, address string, timeout time.Duration) error {
	conn, err := net.DialTimeout(network, address, timeout)
	if conn != nil {
		defer conn.Close()
	}
	return err
}

// PingNWithDelay calls Ping x number of times with delay in between
func PingNWithDelay(network, address string, timeout time.Duration, count int,
	delay time.Duration, log log.FieldLogger) bool {

	for i := 0; i < count; i++ {
		log.Infof("pinging control plane %s", address)
		if err := Ping(network, address, timeout); err == nil {
			return true
		}
		time.Sleep(delay)
	}
	return false
}
