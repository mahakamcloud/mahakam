package parsers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mahakamcloud/mahakam/pkg/terraform/parsers"
	"github.com/mahakamcloud/mahakam/pkg/terraform/templates"
)

var _ = Describe("TerraformParser", func() {
	var (
		tfParser TerraformParser
	)

	BeforeEach(func() {
		var backendData = map[string]string{
			"Bucket": "tf-mahakam",
			"Key":    "gofinance-k8s/control-plane/terraform.tfstate",
			"Region": "ap-southeast-1",
		}

		tfParser = TerraformParser{
			"backend",
			templates.Backend,
			"/tmp/mahakam/terraform/backend.tf",
			backendData,
		}
	})

	parsedBackendData := `terraform {
    backend "s3" {
        bucket = "tf-mahakam"
        key    = "gofinance-k8s/control-plane/terraform.tfstate"
        region = "ap-southeast-1"
    }
}`

	Describe("Parsing Templates", func() {
		Context("With backend.tf data", func() {
			It("should be able to parse templates.Backend and return string", func() {
				result := tfParser.ParseTemplate()
				Expect(result).To(Equal(parsedBackendData))
			})
		})
	})
})
