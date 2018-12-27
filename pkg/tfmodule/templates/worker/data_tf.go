package worker

var Data = `data "template_file" "user_data" {
    template = "${file("${path.module}/templates/user_data.tpl")}"

    vars {
        hostname                 = "${var.hostname}"
        ssh_public_key           = "${var.ssh_public_key}"
        password                 = "${var.password}"
        ip_address               = "${var.ip_address}"
        netmask                  = "${var.netmask}"
        gateway                  = "${var.gateway}"
        dns_address              = "${var.dns_address}"
        dns_domain_name          = "${var.dns_domain_name}"
        controlplane_ip          = "${var.controlplane_ip}"
        kubeadm_token            = "${var.kubeadm_token}"
    }
}`
