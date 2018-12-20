package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// ReadFile reads content bytes of file
func ReadFile(filePath string) ([]byte, error) {
	bytes := []byte{}

	path, err := filepath.Abs(filePath)
	if err != nil {
		return bytes, fmt.Errorf("Error getting file path '%s': %s", filePath, err)
	}

	file, err := os.Open(path)
	if err != nil {
		return bytes, fmt.Errorf("Error opening file '%s': %s", path, err)
	}
	defer file.Close()

	bytes, err = ioutil.ReadAll(file)
	if err != nil {
		return bytes, fmt.Errorf("Error reading file '%s': %s", path, err)
	}

	return bytes, nil
}

// CidrToKeyString converts CIDR format to key string
// i.e. from 1.2.3.4/16 to 1.2.3.4-16
func CidrToKeyString(cidr net.IPNet) string {
	return strings.Replace(cidr.String(), "/", "-", -1)
}
