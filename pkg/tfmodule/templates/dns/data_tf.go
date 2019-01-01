package dns

// Data specfies template of terraform data for dns
var Data = `data "template_file" "user_data" {
    template = "${file("${path.module}/templates/user_data.tpl")}"

    vars {
        hostname                 = "${var.hostname}"
        password                 = "${var.password}"
        ssh_public_key           = "${var.ssh_public_key}"
        ip_address               = "${var.ip_address}"
        mac_address              = "${var.mac_address}"
        netmask                  = "${var.netmask}"
        gateway_ip               = "${var.gateway_ip}"
        upstream_dns            = "${var.upstream_dns}"
        dns_zone_name            = "${var.dns_zone_name}"
        network_cidr             = "${var.network_cidr}"
    }
}`
