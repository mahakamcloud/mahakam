package tfmodule_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mahakamcloud/mahakam/pkg/terraform/templates"
	. "github.com/mahakamcloud/mahakam/pkg/tfmodule"
)

var _ = Describe("TerraformFile", func() {
	var (
		tfFile      TerraformFile
		backendData map[string]string
	)

	BeforeEach(func() {

		tfFile = TerraformFile{
			FileType:    "backend",
			Source:      templates.Backend,
			Destination: "/tmp/mahakam/terraform/backend.tf",
		}
		backendData = map[string]string{
			"Bucket": "tf-mahakam",
			"Key":    "gofinance-k8s/control-plane/terraform.tfstate",
			"Region": "ap-southeast-1",
		}
	})

	parsedBackendData := `terraform {
        backend "s3" {
        bucket = "tf-mahakam"
        key    = "gofinance-k8s/control-plane/terraform.tfstate"
        region = "ap-southeast-1"
    }
}`

	Describe("Generating the parsed file", func() {
		Context("With backend.tf data", func() {
			It("should be able to parse templates.Backend and return string", func() {
				result := tfFile.ParseTerraformFile(backendData)
				Expect(result).To(Equal(parsedBackendData))
			})
		})
	})
})
