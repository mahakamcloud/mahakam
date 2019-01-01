package tfmodule

import (
	"github.com/mahakamcloud/mahakam/pkg/tfmodule/templates"
	"github.com/mahakamcloud/mahakam/pkg/tfmodule/templates/basic"
	"github.com/mahakamcloud/mahakam/pkg/tfmodule/templates/controlplane"
	"github.com/mahakamcloud/mahakam/pkg/tfmodule/templates/dhcp"
	"github.com/mahakamcloud/mahakam/pkg/tfmodule/templates/dns"
	"github.com/mahakamcloud/mahakam/pkg/tfmodule/templates/worker"
)

// CreateNode creates VM with generic configuration through terraform
func CreateNode(name, destdir string, data map[string]string) error {
	var basicNode = &TerraformProvisioner{
		Name:    name,
		DestDir: destdir,
		Files: []TerraformFile{
			TerraformFile{"backend", templates.Backend, destdir, "backend.tf"},
			TerraformFile{"data", basic.Data, destdir, "data.tf"},
			TerraformFile{"main", templates.MainFile, destdir, "main.tf"},
			TerraformFile{"tfvars", basic.TFVars, destdir, "terraform.tfvars"},
			TerraformFile{"vars", basic.Vars, destdir, "vars.tf"},
			TerraformFile{"cloudinit", basic.CloudInit, destdir + "/templates/", "user_data.tpl"},
		},
	}

	basicNode.GenerateProvisionerFiles(data)
	err := basicNode.ExecuteProvisioner()
	if err != nil {
		return err
	}
	return nil
}

// CreateControlPlaneNode creates VM with kubernetes control plane configuration
// through terraform
func CreateControlPlaneNode(name, destdir string, data map[string]string) error {
	cpNode := &TerraformProvisioner{
		Name:    name,
		DestDir: destdir,
		Files: []TerraformFile{
			TerraformFile{"backend", templates.Backend, destdir, "backend.tf"},
			TerraformFile{"data", controlplane.Data, destdir, "data.tf"},
			TerraformFile{"main", templates.MainFile, destdir, "main.tf"},
			TerraformFile{"tfvars", controlplane.TFVars, destdir, "terraform.tfvars"},
			TerraformFile{"vars", controlplane.Vars, destdir, "vars.tf"},
			TerraformFile{"cloudinit", controlplane.CloudInit, destdir + "/templates/", "user_data.tpl"},
		},
	}

	cpNode.GenerateProvisionerFiles(data)
	err := cpNode.ExecuteProvisioner()
	if err != nil {
		return err
	}
	return nil
}

// CreateNetworkDNS creates VM with DNS server configuration through terraform
func CreateNetworkDNS(name, destdir string, data map[string]string) error {
	dnsNode := &TerraformProvisioner{
		Name:    name,
		DestDir: destdir,
		Files: []TerraformFile{
			TerraformFile{"backend", templates.Backend, destdir, "backend.tf"},
			TerraformFile{"data", dns.Data, destdir, "data.tf"},
			TerraformFile{"main", templates.MainFile, destdir, "main.tf"},
			TerraformFile{"tfvars", dns.TFVars, destdir, "terraform.tfvars"},
			TerraformFile{"vars", dns.Vars, destdir, "vars.tf"},
			TerraformFile{"cloudinit", dns.CloudInit, destdir + "/templates/", "user_data.tpl"},
		},
	}

	dnsNode.GenerateProvisionerFiles(data)
	err := dnsNode.ExecuteProvisioner()
	if err != nil {
		return err
	}
	return nil
}

// CreateNetworkDHCP creates VM with DNS server configuration through terraform
func CreateNetworkDHCP(name, destdir string, data map[string]string) error {
	dhcpNode := &TerraformProvisioner{
		Name:    name,
		DestDir: destdir,
		Files: []TerraformFile{
			TerraformFile{"backend", templates.Backend, destdir, "backend.tf"},
			TerraformFile{"data", dhcp.Data, destdir, "data.tf"},
			TerraformFile{"main", templates.MainFile, destdir, "main.tf"},
			TerraformFile{"tfvars", dhcp.TFVars, destdir, "terraform.tfvars"},
			TerraformFile{"vars", dhcp.Vars, destdir, "vars.tf"},
			TerraformFile{"cloudinit", dhcp.CloudInit, destdir + "/templates/", "user_data.tpl"},
		},
	}

	dhcpNode.GenerateProvisionerFiles(data)
	err := dhcpNode.ExecuteProvisioner()
	if err != nil {
		return err
	}
	return nil
}

// CreateWorkerNode creates VM with kubernetes worker configuration
// through terraform
func CreateWorkerNode(name, destdir string, data map[string]string) error {
	wNode := &TerraformProvisioner{
		Name:    name,
		DestDir: destdir,
		Files: []TerraformFile{
			TerraformFile{"backend", templates.Backend, destdir, "backend.tf"},
			TerraformFile{"data", worker.Data, destdir, "data.tf"},
			TerraformFile{"main", templates.MainFile, destdir, "main.tf"},
			TerraformFile{"tfvars", worker.TFVars, destdir, "terraform.tfvars"},
			TerraformFile{"vars", worker.Vars, destdir, "vars.tf"},
			TerraformFile{"cloudinit", worker.CloudInit, destdir + "/templates/", "user_data.tpl"},
		},
	}

	wNode.GenerateProvisionerFiles(data)
	err := wNode.ExecuteProvisioner()
	if err != nil {
		return err
	}
	return nil
}
