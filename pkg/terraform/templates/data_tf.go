package templates

var UserData = `data "template_file" "user_data" {
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
