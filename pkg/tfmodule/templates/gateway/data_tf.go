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
        public_ip_address        = "${var.public_ip_address}"
        public_netmask           = "${var.public_netmask}"
        public_gateway           = "${var.gateway}"
        dns_address              = "${var.dns_address}"
        dns_zone_name            = "${var.dns_domain_name}"
        network_cidr             = "${var.network_cidr}"
    }
}`
