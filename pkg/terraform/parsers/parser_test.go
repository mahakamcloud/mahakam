package parsers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mahakamcloud/mahakam/pkg/terraform/parsers"
	"github.com/mahakamcloud/mahakam/pkg/terraform/templates"
)

var _ = Describe("TerraformParser", func() {
	var (
		tfBackendParser TerraformParser
		tfDataParser    TerraformParser
	)

	BeforeEach(func() {
		var backendTfData = map[string]string{
			"Bucket": "tf-mahakam",
			"Key":    "gofinance-k8s/control-plane/terraform.tfstate",
			"Region": "ap-southeast-1",
		}

		var userData = map[string]string{
			"Bucket": "tf-mahakam",
			"Key":    "gofinance-k8s/control-plane/terraform.tfstate",
			"Region": "ap-southeast-1",
		}

		tfBackendParser = TerraformParser{
			"backend",
			templates.Backend,
			"/tmp/mahakam/terraform/backend.tf",
			backendTfData,
		}

		tfDataParser = TerraformParser{
			"data",
			templates.UserData,
			"/tmp/mahakam/terraform/data.tf",
			userData,
		}
	})

	parsedBackendData := `terraform {
    backend "s3" {
        bucket = "tf-mahakam"
        key    = "gofinance-k8s/control-plane/terraform.tfstate"
        region = "ap-southeast-1"
    }
}`

	parsedUserData := `data "template_file" "user_data" {
    template = "${file("${path.module}/templates/user_data.tpl")}"

    vars {
        hostname                 = "${var.hostname}"
        dns_dhcp_server_username = "${var.dns_dhcp_server_username}"
        password                 = "${var.password}"
        ip_address               = "${var.ip_address}"
        netmask                  = "${var.netmask}"
        gate_nss_api_key         = "${var.gate_nss_api_key}"
        dns_address              = "${var.dns_address}"
        #CUSTOM_PARAMETERS
    }
}`

	Describe("Parsing Templates", func() {
		Context("With backend.tf data", func() {
			It("should be able to parse templates.Backend and return the correct string", func() {
				result := tfBackendParser.ParseTemplate()
				Expect(result).To(Equal(parsedBackendData))
			})
		})

		Context("With data.tf data", func() {
			It("should be able to parse templates.Data and return the correct string", func() {
				result := tfDataParser.ParseTemplate()
				Expect(result).To(Equal(parsedUserData))
			})
		})
	})
})
