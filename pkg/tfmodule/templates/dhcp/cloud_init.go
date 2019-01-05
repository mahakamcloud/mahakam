package dhcp

var CloudInit = `#cloud-config
password: ${password}
chpasswd: { expire: False }
ssh_pwauth: True
hostname: ${hostname}
fqdn: ${hostname}.${dns_zone_name}

ssh_authorized_keys:
  - ${ssh_public_key}

resolv_conf:
  nameservers: ['${dns_address}']
  searchdomains:
    - ${dns_zone_name}

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

  - path: /opt/cloud-init/setup-dhcp.sh
    permissions: 0644
    content: |
      apt install -y isc-dhcp-server
      systemctl stop ufw.service

      cat <<EOF >/etc/dhcp/dhcpd.conf
      ddns-update-style none;

      authoritative;

      option domain-name "${dns_zone_name}";
      option domain-name-servers ${dns_address};

      default-lease-time 7200;
      max-lease-time 14400;

      log-facility local7;

      subnet ${subnet_address} netmask ${subnet_mask} {
        range ${ip_address} ${ip_address};
        option domain-name-servers ${dns_address};
        option domain-name "${dns_zone_name}";
        option subnet-mask ${subnet_mask};
        option routers ${gateway};
        option broadcast-address ${broadcast_address};
        default-lease-time 7200;
        max-lease-time 14400;
      }
      EOF

      cat <<EOF >/etc/dhcp/dhcpd.conf.tpl
      ddns-update-style none;

      authoritative;

      option domain-name "${dns_zone_name}";
      option domain-name-servers ${dns_address};

      default-lease-time 7200;
      max-lease-time 14400;

      log-facility local7;

      subnet ${subnet_address} netmask ${subnet_mask} {
        range ${ip_address} ${ip_address};
        option domain-name-servers ${dns_address};
        option domain-name "${dns_zone_name}";
        option subnet-mask ${subnet_mask};
        option routers ${gateway};
        option broadcast-address ${broadcast_address};
        default-lease-time 7200;
        max-lease-time 14400;
      }

      [[ range ls "hosts" ]][[ .Value ]][[ end ]]
      EOF

      cat <<EOF >/etc/default/isc-dhcp-server
      INTERFACES="ens3"
      EOF

      systemctl start isc-dhcp-server

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

  - path: /opt/cloud-init/setup-consul-docker.sh
    permissions: 0644
    content: |
      docker run -d --name=dev-consul -e CONSUL_BIND_INTERFACE=eth0 -p 8500:8500 -p 8600:8600 consul

  - path: /opt/cloud-init/setup-consul-template.sh
    permissions: 0644
    content: |
      echo "Downloading consul-template 0.19.5"
      curl --silent --output /tmp/consul-template_0.19.5_linux_amd64.zip https://releases.hashicorp.com/consul-template/0.19.5/consul-template_0.19.5_linux_amd64.zip

      echo "Setup consul-template user and group"
      groupadd -r consul-template
      useradd -r -g consul-template -d /var/lib/consul-template -s /sbin/nologin -c "consul-template user" consul-template
      mkdir -p /etc/consul-template.d

      echo "Installing consul-template"
      apt install -y unzip
      sudo unzip -o /tmp/consul-template_0.19.5_linux_amd64.zip -d /usr/local/bin/
      sudo chmod 0755 /usr/local/bin/consul-template
      sudo chown consul-template:consul-template /usr/local/bin/consul-template

      echo "/usr/local/bin/consul-template --version: $(/usr/local/bin/consul-template --version)"

      echo "Adding consul-template config file"
      cat <<EOF >/etc/consul-template.d/bind.hcl
      consul {
        address = "127.0.0.1:8500"
        retry {
          enabled = true
          attempts = 12
          backoff = "250ms"
          max_backoff = "1m"
        }
      }

      reload_signal = "SIGHUP"
      kill_signal = "SIGINT"
      max_stale = "10m"
      log_level = "warn"

      wait {
        min = "5s"
        max = "10s"
      }

      template {
        source = "/etc/dhcp/dhcpd.conf.tpl"
        destination = "/etc/dhcp/dhcpd.conf"
        create_dest_dirs = true
        command = "systemctl restart isc-dhcp-server"
        command_timeout = "60s"
        error_on_missing_key = false
        perms = 0644
        backup = true
        left_delimiter  = "[["
        right_delimiter = "]]"
        wait {
          min = "2s"
          max = "10s"
        }
      }
      EOF

      echo "Configuring consul-template"
      sudo mkdir -pm 0755 /etc/consul-template.d /opt/consul-template/data
      sudo chown -R consul-template:consul-template /etc/consul-template.d /opt/consul-template/data
      sudo chmod -R 0644 /etc/consul-template.d/*

      echo "Installing consul template systemd service and config"

      cat <<EOF >/etc/systemd/system/consul-template.service
      [Unit]
      Description=consul-template
      Requires=network-online.target
      After=network-online.target consul.service

      [Service]
      ExecStart=/usr/local/bin/consul-template -config=/etc/consul-template.d
      KillSignal=SIGINT
      ExecReload=/bin/kill -HUP $MAINPID
      Restart=always
      RestartSec=5

      [Install]
      WantedBy=multi-user.target
      EOF

      systemctl daemon-reload
      systemctl enable consul-template
      systemctl start consul-template

      echo "Completed consul template setup"

runcmd:
  - echo "Configuring DNS VM" 
  - bash -ex /opt/cloud-init/setup-network.sh 
  - bash -ex /opt/cloud-init/setup-dhcp.sh 
  - bash -ex /opt/cloud-init/setup-docker.sh 
  - bash -ex /opt/cloud-init/setup-consul-docker.sh 
  - bash -ex /opt/cloud-init/setup-consul-template.sh 

final_message: "The system is finally up, after $UPTIME seconds"
`
