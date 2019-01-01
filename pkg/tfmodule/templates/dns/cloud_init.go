package dns

var CloudInit = `password: ${password}
chpasswd: { expire: False }
ssh_pwauth: True
hostname: ${hostname}
fqdn: ${hostname}.${dns_zone_name}
ssh_authorized_keys:
  - ${ssh_public_key}

resolv_conf:
  nameservers: ['8.8.8.8']
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
        dns-nameservers 8.8.8.8
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

packages:
  - curl
  - apt-transport-https
  - ca-certificates
  - bind9
  - bind9utils
  - bind9-doc

  - path: /opt/cloud-init/setup-bind.sh
    permissions: 0644
    content: |
      echo "Setup Bind"

      cat <<EOF >/etc/bind/named.conf.options
      options {
        directory "/var/cache/bind";
      
        recursion yes;
        allow-recursion { any; };
        listen-on { ${ip_address} };
        allow-transfer { none; };
        zone-statistics yes;
      
        forwarders {
          8.8.8.8;
          1.1.1.1;
        };
      };
      
      statistics-channels {
        inet {{ ansible_host }} port 8080;
      };      
      EOF

      cat <<EOF >/etc/bind/named.conf.local
      zone "${dns_zone_name}" {
        type master;
        file "/var/lib/bind/zones/db.${dns_zone_name}";
        allow-transfer { 10.30.1.13; };
        allow-update { ${ip_address} };
      };
      
      zone "30.10.in-addr.arpa" {
        type master;
        file "/var/lib/bind/zones/db.10.30";
        allow-transfer { 10.30.1.13; };
        allow-update { ${ip_address} };
      };      
      EOF

      mkdir -p /var/lib/bind/zones
      chmod g+s to /var/lib/bind/zones

      cat <<EOF >/var/lib/bind/zones/db.${dns_zone_name}
      $TTL    604800
      @       IN      SOA     ns1.${dns_zone_name}. admin.${dns_zone_name}. (
                        3     ; Serial
                  604800     ; Refresh
                    86400     ; Retry
                  2419200     ; Expire
                  604800 )   ; Negative Cache TTL
      ;
      ; name servers - NS records
          IN      NS      ns1.${dns_zone_name}.
          IN      NS      ns2.${dns_zone_name}.

      ; name servers - A records
      $TTL 60
      ns1.${dns_zone_name}.          IN      A       10.30.1.3
      ns2.${dns_zone_name}.          IN      A       10.30.1.13
      EOF

      cat <<EOF >/var/lib/bind/zones/db.${dns_zone_name}.tpl
      $TTL    604800
      @       IN      SOA     ns1.${dns_zone_name}. admin.${dns_zone_name}. (
                  {{ keyOrDefault "zones/${dns_zone_name}/serial"  "0" }}
                  604800     ; Refresh
                    86400     ; Retry
                  2419200     ; Expire
                  604800 )   ; Negative Cache TTL
      ;
      ; name servers - NS records
          IN      NS      ns1.${dns_zone_name}.
          IN      NS      ns2.${dns_zone_name}.

      ; name servers - A records
      $TTL 60
      ns1.${dns_zone_name}.          IN      A       10.30.1.3
      ns2.${dns_zone_name}.          IN      A       10.30.1.13

      {{ range ls zones/${dns_zone_name}/hosts/ }}
      {{ .Value }}
      {{ end }}
      EOF

      cat <<EOF >/var/lib/bind/zones/db.10.30
      $TTL    604800
      @       IN      SOA     ${dns_zone_name}. admin.${dns_zone_name}. (
                                    4         ; Serial
                              604800         ; Refresh
                                86400         ; Retry
                              2419200         ; Expire
                              604800 )       ; Negative Cache TTL
      ; name servers
            IN      NS      ns1.${dns_zone_name}.
            IN      NS      ns2.${dns_zone_name}.

      ; PTR Records
      3.1    IN      PTR     ns1.${dns_zone_name}.
      13.1   IN      PTR     ns2.${dns_zone_name}.
      3.1    IN      PTR     i-dctv-dns-01.${dns_zone_name}.
      13.1   IN      PTR     i-dctv-dns-02.${dns_zone_name}.
      EOF

      systemctl restart bind9.service

runcmd:
  - echo "Configuring DNS VM" >> /var/log/start-dns.log
  - bash -ex /opt/cloud-init/setup-network.sh >> /var/log/start-dns.log 2>&1
  - bash -ex /opt/cloud-init/setup-docker.sh >> /var/log/start-dns.log 2>&1
  - bash -ex /opt/cloud-init/setup-bind.sh >> /var/log/start-dns.log 2>&1

final_message: "The system is finally up, after $UPTIME seconds"
`
