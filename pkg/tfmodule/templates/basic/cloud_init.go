package basic

var CloudInit = `password: ${password}
chpasswd: { expire: False }
ssh_pwauth: False
hostname: ${hostname}
fqdn: ${hostname}.gocloud.io
`
