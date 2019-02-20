package network

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestGenerateMacAddress(t *testing.T) {
	macAddress := GenerateMacAddress()
	firstOctet := strings.Split(macAddress, ":")[0]
	hexOctet, err := hex.DecodeString(firstOctet)
	if err != nil {
		t.Error("The MAC Address cannot be converted to hexadecimal")
	}

	if (hexOctet[0] & 1) == 1 {
		t.Errorf("expected the address %s to be unicast", macAddress)
	}

	if hexOctet[0] == 0xfe {
		t.Errorf("expected the address %s first octet to not start with 0xFE", macAddress)
	}

}
