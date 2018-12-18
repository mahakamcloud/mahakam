package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
