package writers_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mahakamcloud/mahakam/pkg/terraform/parsers"
	"github.com/mahakamcloud/mahakam/pkg/terraform/templates"

	. "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	. "github.com/mahakamcloud/mahakam/pkg/terraform/writers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TerraformWriter", func() {
	var (
		tfWriter TerraformWriter
		tfFile   string
	)

	BeforeEach(func() {
		var data = map[string]string{
			"Bucket": "tf-mahakam",
			"Key":    "gofinance-k8s/control-plane/terraform.tfstate",
			"Region": "ap-southeast-1",
		}
		terraformResource := NewResourceTerraform("backend.tf", data)

		backendParser := parsers.TerraformParser{
			"backend",
			templates.Backend,
			terraformResource.GetName(),
			terraformResource.GetAttributes(),
		}
		backendTf := backendParser.ParseTemplate()
		//		fmt.Println(backendTf)
		tfWriter = TerraformWriter{
			backendTf,
			"/tmp/mahakam/terraform/",
			"backend.tf",
		}

		tfFile = filepath.Join(tfWriter.DestDirectory, tfWriter.DestFile)
	})

	Describe("Writing Terraform files", func() {
		Context("with directory structure that does not exist yet", func() {
			It("creates the required directory structure", func() {
				tfWriter.WriteFile()
				_, err := os.Stat(tfWriter.DestDirectory)
				Expect(os.IsNotExist(err)).To(Equal(false))
			})

			It("creates the destination file", func() {
				tfWriter.WriteFile()
				_, err := os.Stat(tfFile)
				Expect(os.IsNotExist(err)).To(Equal(false))
			})

			It("writes the correct data to the file", func() {
				tfWriter.WriteFile()
				readData, err := ioutil.ReadFile(tfFile)
				parsedBackendData := `terraform {
    backend "s3" {
        bucket = "tf-mahakam"
        key    = "gofinance-k8s/control-plane/terraform.tfstate"
        region = "ap-southeast-1"
    }
}`

				Expect(os.IsNotExist(err)).To(Equal(false))
				Expect(string(readData)).To(Equal(parsedBackendData))
			})
		})
	})
})
