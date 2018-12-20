package templates

var MainFile = `module "test-vm" {
	source = "{{.LibvirtModulePath}}"
	# git::https://source.golabs.io/terraform/libvirt-vm-private.git?ref=v3.3.2
	instance_name = "${var.hostname}"
	libvirt_host  = "${var.host}"
	source_path   = "${var.image_source_path}"
	mac_address   = "${var.mac_address}"
	ip_address    = "${var.ip_address}"
  
	user_data = "${data.template_file.user_data.rendered}"
  }`
