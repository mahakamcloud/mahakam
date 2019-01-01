package gateway

// Data specfies template of terraform data for dns
var Data = `data "template_file" "user_data" {
    template = "${file("${path.module}/templates/user_data.tpl")}"

    vars {
        hostname                 = "${var.hostname}"
        password                 = "${var.password}"
        ssh_public_key           = "${var.ssh_public_key}"
        ip_address               = "${var.ip_address}"
        netmask                  = "${var.netmask}"
        gateway                  = "${var.gateway}"
        dns_address              = "${var.dns_address}"
        dns_zone_name            = "${var.dns_zone_name}"
        network_cidr             = "${var.network_cidr}"
    }
}`
