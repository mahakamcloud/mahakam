package tfmodule_test

import (
	"fmt"

	. "github.com/mahakamcloud/mahakam/pkg/tfmodule"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("TerraformFile", func() {

	data := map[string]string{
		"Bucket":                "tf-mahakam",
		"Key":                   "gofinance-k8s/control-plane/terraform.tfstate",
		"Region":                "ap-southeast-1",
		"IPAddress":             "10.30.30.1",
		"DNSDhcpServerUsername": "himani.agrawal",
		"GateNssAPIKEY":         "key",
		"Host":                  "10.30.0.1",
		"Name":                  "mahakam-test-01",
		"MacAddress":            "C4:AC:69:E4:D0:24",
		"NetMask":               "255.255.255.0",
		"DNSAddress":            "10.30.1.3",
	}

	Describe("Generating the parsed file", func() {
		Context("With backend.tf data", func() {
			It("should be able to parse templates.Backend and return string", func() {
				CreateNode(data)
				fmt.Println("done")
			})
		})
	})
})
