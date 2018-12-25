package templates

var Data = `data "template_file" "user_data" {
    template = "${file("${path.module}/templates/user_data.tpl")}"

    vars {
        hostname                 = "${var.hostname}"
        password                 = "${var.password}"
        ip_address               = "${var.ip_address}"
        netmask                  = "${var.netmask}"
        dns_address              = "${var.dns_address}"
    }
}`
