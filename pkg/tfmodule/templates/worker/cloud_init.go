package worker

var CloudInit = `password: ${password}
chpasswd: { expire: False }
ssh_pwauth: True
hostname: ${hostname}
fqdn: ${hostname}.${dns_domain_name}
ssh_authorized_keys:
  - ${ssh_public_key}

resolv_conf:
  nameservers: ['8.8.8.8']
  searchdomains:
    - ${dns_domain_name}

package_upgrade: true

write_files:
  - path: /opt/cloud-init/setup-network.sh
    permissions: 0644
    content: |
      cat <<EOF >/etc/network/interfaces
      auto lo
      iface lo inet loopback

      auto ens3
      iface ens3 inet static
        address ${ip_address}
        netmask ${netmask}
        dns-nameservers ${dns_address}
        gateway ${gateway}
      EOF

      ifdown ens3 && ifup ens3
      systemctl restart networking

  - path: /opt/cloud-init/setup-docker.sh
    permissions: 0644
    content: |
      echo "Setup Docker"

      curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
      add-apt-repository \
        "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
        $(lsb_release -cs) \
        stable"
      apt-get update
      apt-get install -y docker-ce=17.03.3~ce-0~ubuntu-xenial

  - path: /opt/cloud-init/kubeadm-join.sh
    owner: root:root
    permissions: 0644
    content: |
      apt-key adv --keyserver hkp://keyserver.ubuntu.com --recv-keys 0xF76221572C52609D 0x3746C208A7317B0F
      echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
      apt-get update && apt-get install -y --allow-unauthenticated kubelet=1.12.3-00 kubeadm=1.12.3-00 kubectl=1.12.3-00 kubernetes-cni

      systemctl daemon-reload
      systemctl enable docker
      systemctl enable kubelet
      systemctl start docker

      echo "127.0.1.1 ${hostname}" >> /etc/hosts
      sleep 120 && for i in $(seq 10); do echo "kubeadm join $i" && kubeadm join --token=${kubeadm_token} --discovery-token-unsafe-skip-ca-verification ${controlplane_ip}:6443 && break || sleep 15; done

      sysctl -w net.bridge.bridge-nf-call-iptables=1

packages:
  - curl
  - apt-transport-https
  - ca-certificates

runcmd:
  - echo "Starting worker VM" >> /var/log/start-worker.log
  - bash -ex /opt/cloud-init/setup-public-key.sh >> /var/log/start-worker.log 2>&1
  - bash -ex /opt/cloud-init/setup-network.sh >> /var/log/start-worker.log 2>&1
  - bash -ex /opt/cloud-init/setup-docker.sh >> /var/log/start-worker.log 2>&1
  - bash -ex /opt/cloud-init/kubeadm-join.sh >> /var/log/start-worker.log 2>&1

final_message: "The system is finally up, after $UPTIME seconds"
`
