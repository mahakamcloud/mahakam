package tfmodule_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mahakamcloud/mahakam/pkg/terraform/templates"
	. "github.com/mahakamcloud/mahakam/pkg/tfmodule"
)

var _ = Describe("TerraformFile", func() {
	var (
		tfFile          TerraformFile
		tfProvisioner   TerraformProvisioner
		backendData     map[string]string
		destinationFile string
	)

	BeforeEach(func() {

		tfFile = TerraformFile{
			FileType: "backend",
			Source:   templates.Backend,
			DestDir:  "/tmp/mahakam/terraform/mahakam-cluster-01/",
			DestFile: "backend.tf",
		}
		tfProvisioner = TerraformProvisioner{
			Name:    "mahakam-cluster-01",
			DestDir: "/tmp/mahakam/terraform/",
			Files: []TerraformFile{
				tfFile,
			},
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

	Describe("Generating the Provisioner file", func() {
		Context("With backend.tf data", func() {
			It("Backend.tf file should be created at Destination with ParsedBackendData", func() {
				tfProvisioner.GenerateProvisionerFiles(backendData)
				_, err := os.Stat(destinationFile)
				Expect(os.IsNotExist(err)).To(Equal(false))

				readData, err := ioutil.ReadFile(destinationFile)
				Expect(os.IsNotExist(err)).To(Equal(false))
				Expect(string(readData)).To(Equal(parsedBackendData))

			})
		})
	})
})
