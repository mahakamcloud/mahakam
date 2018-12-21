package tfmodule_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTfmodule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tfmodule Suite")
}
