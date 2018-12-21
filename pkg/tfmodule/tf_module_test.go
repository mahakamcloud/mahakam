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
		tfFile                 TerraformFile
		tfProvisioner          TerraformProvisioner
		backendData            map[string]string
		backednDestinationFile string
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
		backednDestinationFile = filepath.Join(tfFile.DestDir, tfFile.DestFile)
	})

	parsedBackendData := `terraform {
    backend "s3" {
        bucket = "tf-mahakam"
        key    = "gofinance-k8s/control-plane/terraform.tfstate"
        region = "ap-southeast-1"
    }
}`

	tfDataFile := TerraformFile{
		FileType: "data",
		Source:   templates.UserData,
		DestDir:  "/tmp/mahakam/terraform/mahakam-cluster-01/",
		DestFile: "data.tf",
	}

	Describe("Generating the Provisioner file", func() {
		Context("With backend.tf data", func() {
			It("Backend.tf file should be created at Destination with ParsedBackendData", func() {
				tfProvisioner.GenerateProvisionerFiles(backendData)
				_, err := os.Stat(backednDestinationFile)
				Expect(os.IsNotExist(err)).To(Equal(false))

				readData, err := ioutil.ReadFile(backednDestinationFile)
				Expect(os.IsNotExist(err)).To(Equal(false))
				Expect(string(readData)).To(Equal(parsedBackendData))

			})
		})
		Context("Update the terraform provisioner", func() {
			It("With data.tf file", func() {
				(&tfProvisioner).UpdateProvisionerFile("data", templates.UserData, "data.tf")
				Expect(tfProvisioner.Files).To(Equal([]TerraformFile{tfFile, tfDataFile}))
			})
		})

	})
})
