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
		tfMainParser    TerraformParser
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

		var mainData = map[string]string{
			"Name":              "mahakam-spike-01",
			"LibvirtModulePath": "git::https://source.golabs.io/terraform/libvirt-vm-private.git?ref=v3.3.2",
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

		tfMainParser = TerraformParser{
			"main",
			templates.MainFile,
			"/tmp/mahakam/terraform/main.tf",
			mainData,
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

	parsedMainData := `module "mahakam-spike-01" {
    source = "git::https://source.golabs.io/terraform/libvirt-vm-private.git?ref=v3.3.2"

    instance_name = "${var.hostname}"
    libvirt_host  = "${var.host}"
    source_path   = "${var.image_source_path}"
    mac_address   = "${var.mac_address}"
    ip_address    = "${var.ip_address}"

    user_data = "${data.template_file.user_data.rendered}"
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

		Context("With main.tf data", func() {
			It("should be able to parse templates.Data and return the correct string", func() {
				result := tfMainParser.ParseTemplate()
				Expect(result).To(Equal(parsedMainData))
			})
		})
	})
})
