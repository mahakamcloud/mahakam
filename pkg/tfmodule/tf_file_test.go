package tfmodule_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mahakamcloud/mahakam/pkg/tfmodule"
	"github.com/mahakamcloud/mahakam/pkg/tfmodule/templates"
)

var _ = Describe("TerraformFile", func() {
	var (
		tfFile          TerraformFile
		backendData     map[string]string
		destinationFile string
	)

	BeforeEach(func() {

		tfFile = TerraformFile{
			FileType: "backend",
			Source:   templates.Backend,
			DestDir:  "/tmp/mahakam/terraform",
			DestFile: "backend.tf",
		}
		backendData = map[string]string{
			"Bucket": "tf-mahakam",
			"Key":    "gofinance-k8s/control-plane/terraform.tfstate",
			"Region": "ap-southeast-1",
		}
		destinationFile = filepath.Join(tfFile.DestDir, tfFile.DestFile)
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

	Describe("Writing the parsed file", func() {
		Context("With backend.tf data", func() {
			It("creates the required directory structure", func() {
				tfFile.WriteTerraformFile(parsedBackendData)
				_, err := os.Stat(tfFile.DestDir)
				Expect(os.IsNotExist(err)).To(Equal(false))
			})

			It("creates the destination file", func() {
				tfFile.WriteTerraformFile(parsedBackendData)
				_, err := os.Stat(destinationFile)
				Expect(os.IsNotExist(err)).To(Equal(false))
			})

			It("writes the correct data to the file", func() {
				tfFile.WriteTerraformFile(parsedBackendData)
				readData, err := ioutil.ReadFile(destinationFile)
				Expect(os.IsNotExist(err)).To(Equal(false))
				Expect(string(readData)).To(Equal(parsedBackendData))

			})

			It("overwrite file data if file is generated again", func() {
				tfFile.WriteTerraformFile(parsedBackendData)
				readData, err := ioutil.ReadFile(destinationFile)
				Expect(os.IsNotExist(err)).To(Equal(false))
				Expect(string(readData)).To(Equal(parsedBackendData))

			})
		})
	})
})
