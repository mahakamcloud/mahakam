package dns

// CloudInit stores cloud-init template
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

network:
  version: 2
  ethernets:
    ens3:
      match:
        macaddress: ${mac_address}
      wakeonlan: true
      dhcp4: true
      gateway4: ${gateway}
      nameservers:
        search: [${dns_zone_name}]
        addresses: [${dns_address}]
  # static routes
  routes:
   - to: 0.0.0.0/0
     via: ${gateway}
     metric: 3

package_upgrade: true

packages:
  - curl
  - apt-transport-https
  - ca-certificates
  - unzip
  - bind9
  - bind9utils
  - bind9-doc

write_files:
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
      mkdir /consul/data
      docker run -d --name=dev-consul -e CONSUL_BIND_INTERFACE=ens3 -v /consul/data:/consul/data -p 8500:8500 -p 8600:8600 consul

  - path: /opt/cloud-init/setup-bind-zone-dir.sh
    permissions: 0644
    content: |
      echo "Setup Bind"

      mkdir -p /var/lib/bind/zones
      chmod g+s /var/lib/bind/zones

  - path: /etc/bind/named.conf.options
    permissions: 0644
    content: |
      # run resolvconf?
      RESOLVCONF=no

      # startup options for the server
      OPTIONS="-u bind -4"

  - path: /etc/bind/named.conf.options
    permissions: 0644
    content: |
      options {
        directory "/var/cache/bind";

        recursion yes;
        allow-recursion { any; };
        listen-on { ${ip_address}; };
        allow-update { ${ip_address}; };
        zone-statistics yes;

        forwarders {
          ${dns_address};
        };
      };

      statistics-channels {
        inet ${ip_address} port 8080;
      };

  - path: /etc/bind/named.conf.local
    permissions: 0644
    content: |
      zone "${dns_zone_name}" {
        type master;
        file "/var/lib/bind/zones/db.${dns_zone_name}";
        allow-update { ${ip_address}; };
      };

  - path: /var/lib/bind/zones/db.${dns_zone_name}
    permissions: 0644
    content: |
      $ORIGIN .
      $TTL 604800 ; 1 week
      ${dns_zone_name}.      IN      SOA     ns1.${dns_zone_name}. admin.${dns_zone_name}. (
                            3         ; Serial
                        604800        ; Refresh
                          86400       ; Retry
                        2419200       ; Expire
                        604800        ; Negative Cache TTL
          )
      ;
      ; name servers - NS records
                         NS      ns1.${dns_zone_name}.

      ; name servers - A records
      $TTL 60
      ns1.${dns_zone_name}.          IN      A       ${ip_address}

      $ORIGIN ${dns_zone_name}.
      $TTL 60 ; 1 minute

  - path: /var/lib/bind/zones/db.${dns_zone_name}.tpl
    permissions: 0644
    content: |
      $ORIGIN .
      $TTL 604800 ; 1 week
      ${dns_zone_name}.      IN      SOA     ns1.${dns_zone_name}. admin.${dns_zone_name}. (
                            [[ keyOrDefault "zones/${dns_zone_name}/serial"  "0" ]]         ; Serial
                        604800        ; Refresh
                          86400       ; Retry
                        2419200       ; Expire
                        604800        ; Negative Cache TTL
          )
      ;
      ; name servers - NS records
                         NS      ns1.${dns_zone_name}.

      ; name servers - A records
      $TTL 60
      ns1.${dns_zone_name}.          IN      A       ${ip_address}

      $ORIGIN ${dns_zone_name}.
      $TTL 60 ; 1 minute

      [[ range ls "zones/${dns_zone_name}/hosts/" ]]
      [[ .Key ]]        A        [[ .Value ]][[ end ]]

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
        source = "/var/lib/bind/zones/db.mgmt.gocloud.io.tpl"
        destination = "/var/lib/bind/zones/db.mgmt.gocloud.io"
        create_dest_dirs = true
        command = "bash -c 'chown root:bind /var/lib/bind/zones/db.mgmt.gocloud.io && systemctl reload bind9'"
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

bootcmd:
  - echo "127.0.1.1 ${hostname}" >> /etc/hosts

runcmd:
  - echo "Configuring DNS VM"
  - bash -ex /opt/cloud-init/setup-docker.sh
  - bash -ex /opt/cloud-init/setup-consul-docker.sh
  - bash -ex /opt/cloud-init/setup-bind-zone-dir.sh
  - [ systemctl, daemon-reload ]
  - [ systemctl, enable, bind9.service ]
  - [ systemctl, start, bind9.service ]
  - bash -ex /opt/cloud-init/setup-consul-template.sh
  - [ systemctl, daemon-reload ]
  - [ systemctl, enable, consul-template.service ]
  - [ systemctl, start, consul-template.service ]

final_message: "The system is finally up, after $UPTIME seconds"
`
