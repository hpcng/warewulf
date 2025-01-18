package upgrade

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configUpgradeTests = []struct {
	name         string
	legacyYaml   string
	upgradedYaml string
}{
	{
		name:         "empty",
		legacyYaml:   ``,
		upgradedYaml: `{}`,
	},
	{
		name: "v4.0.0",
		legacyYaml: `
ipaddr: 192.168.1.1
netmask: 255.255.255.0
warewulf:
  port: 9873
  secure: true
  update interval: 60
dhcp:
  enabled: true
  range start: 192.168.1.150
  range end: 192.168.1.200
  template: default
  systemd name: dhcpd
tftp:
  enabled: true
  tftproot: /var/lib/tftpboot
  systemd name: tftp
nfs:
  systemd name: nfs-server
  exports:
  - /home
  - /var/warewulf
`,
		upgradedYaml: `
ipaddr: 192.168.1.1
netmask: 255.255.255.0
warewulf:
  port: 9873
  secure: true
  update interval: 60
dhcp:
  enabled: true
  template: default
  range start: 192.168.1.150
  range end: 192.168.1.200
  systemd name: dhcpd
tftp:
  enabled: true
  tftproot: /var/lib/tftpboot
  systemd name: tftp
nfs:
  export paths:
    - path: /home
    - path: /var/warewulf
  systemd name: nfs-server
`,
	},
	{
		name: "v4.1.0",
		legacyYaml: `
ipaddr: 192.168.1.1
netmask: 255.255.255.0
warewulf:
  port: 9873
  secure: true
  update interval: 60
dhcp:
  enabled: true
  range start: 192.168.1.150
  range end: 192.168.1.200
  template: default
  systemd name: dhcpd
tftp:
  enabled: true
  tftproot: /var/lib/tftpboot
  systemd name: tftp
nfs:
  systemd name: nfs-server
  exports:
  - /home
  - /var/warewulf
`,
		upgradedYaml: `
ipaddr: 192.168.1.1
netmask: 255.255.255.0
warewulf:
  port: 9873
  secure: true
  update interval: 60
dhcp:
  enabled: true
  template: default
  range start: 192.168.1.150
  range end: 192.168.1.200
  systemd name: dhcpd
tftp:
  enabled: true
  tftproot: /var/lib/tftpboot
  systemd name: tftp
nfs:
  export paths:
    - path: /home
    - path: /var/warewulf
  systemd name: nfs-server
`,
	},
	{
		name: "v4.2.0",
		legacyYaml: `
ipaddr: 192.168.200.1
netmask: 255.255.255.0
warewulf:
  port: 9873
  secure: true
  autobuild overlays: true
  update interval: 60
  syslog: false
dhcp:
  enabled: true
  range start: 192.168.200.50
  range end: 192.168.200.99
  template: default
  systemd name: dhcpd
tftp:
  enabled: true
  tftproot: /var/lib/tftpboot
  systemd name: tftp
nfs:
  systemd name: nfs-server
  exports:
  - /home
  - /var/warewulf
`,
		upgradedYaml: `
ipaddr: 192.168.200.1
netmask: 255.255.255.0
warewulf:
  port: 9873
  secure: true
  update interval: 60
  autobuild overlays: true
  syslog: false
dhcp:
  enabled: true
  template: default
  range start: 192.168.200.50
  range end: 192.168.200.99
  systemd name: dhcpd
tftp:
  enabled: true
  tftproot: /var/lib/tftpboot
  systemd name: tftp
nfs:
  export paths:
    - path: /home
    - path: /var/warewulf
  systemd name: nfs-server
`,
	},
	{
		name: "v4.3.0",
		legacyYaml: `
WW_INTERNAL: 43
ipaddr: 192.168.200.1
netmask: 255.255.255.0
network: 192.168.200.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: false
  datastore: ""
dhcp:
  enabled: true
  template: default
  range start: 192.168.200.50
  range end: 192.168.200.99
  systemd name: dhcpd
tftp:
  enabled: true
  tftproot: ""
  systemd name: tftp
nfs:
  enabled: true
  export paths:
  - path: /home
    export options: rw,sync
    mount options: defaults
    mount: true
  - path: /opt
    export options: ro,sync,no_root_squash
    mount options: defaults
    mount: false
  systemd name: nfs-server
`,
		upgradedYaml: `
ipaddr: 192.168.200.1
netmask: 255.255.255.0
network: 192.168.200.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: false
dhcp:
  enabled: true
  template: default
  range start: 192.168.200.50
  range end: 192.168.200.99
  systemd name: dhcpd
tftp:
  enabled: true
  systemd name: tftp
nfs:
  enabled: true
  export paths:
    - path: /home
      export options: rw,sync
    - path: /opt
      export options: ro,sync,no_root_squash
  systemd name: nfs-server
`,
	},
	{
		name: "v4.4.1",
		legacyYaml: `
WW_INTERNAL: 43
ipaddr: 192.168.200.1
netmask: 255.255.255.0
network: 192.168.200.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: false
dhcp:
  enabled: true
  range start: 192.168.200.50
  range end: 192.168.200.99
  systemd name: dhcpd
tftp:
  enabled: true
  systemd name: tftp
nfs:
  enabled: true
  export paths:
  - path: /home
    export options: rw,sync
    mount options: defaults
    mount: true
  - path: /opt
    export options: ro,sync,no_root_squash
    mount options: defaults
    mount: false
  systemd name: nfs-server
`,
		upgradedYaml: `
ipaddr: 192.168.200.1
netmask: 255.255.255.0
network: 192.168.200.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: false
dhcp:
  enabled: true
  range start: 192.168.200.50
  range end: 192.168.200.99
  systemd name: dhcpd
tftp:
  enabled: true
  systemd name: tftp
nfs:
  enabled: true
  export paths:
    - path: /home
      export options: rw,sync
    - path: /opt
      export options: ro,sync,no_root_squash
  systemd name: nfs-server
`,
	},
	{
		name: "v4.5.8",
		legacyYaml: `
WW_INTERNAL: 45
ipaddr: 10.0.0.1
netmask: 255.255.252.0
network: 10.0.0.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: false
dhcp:
  enabled: true
  range start: 10.0.1.1
  range end: 10.0.1.255
  systemd name: dhcpd
tftp:
  enabled: true
  systemd name: tftp
  ipxe:
    00:09: ipxe-snponly-x86_64.efi
    00:00: undionly.kpxe
    00:0B: arm64-efi/snponly.efi
    00:07: ipxe-snponly-x86_64.efi
nfs:
  enabled: true
  export paths:
  - path: /home
    export options: rw,sync
    mount options: defaults
    mount: true
  - path: /opt
    export options: ro,sync,no_root_squash
    mount options: defaults
    mount: false
  systemd name: nfs-server
container mounts:
  - source: /etc/resolv.conf
    dest: /etc/resolv.conf
    readonly: true
ssh:
  key types:
    - rsa
    - dsa
    - ecdsa
    - ed25519
`,
		upgradedYaml: `
ipaddr: 10.0.0.1
netmask: 255.255.252.0
network: 10.0.0.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: false
dhcp:
  enabled: true
  range start: 10.0.1.1
  range end: 10.0.1.255
  systemd name: dhcpd
tftp:
  enabled: true
  systemd name: tftp
  ipxe:
    00:0B: arm64-efi/snponly.efi
    "00:00": undionly.kpxe
    "00:07": ipxe-snponly-x86_64.efi
    "00:09": ipxe-snponly-x86_64.efi
nfs:
  enabled: true
  export paths:
    - path: /home
      export options: rw,sync
    - path: /opt
      export options: ro,sync,no_root_squash
  systemd name: nfs-server
ssh:
  key types:
    - rsa
    - dsa
    - ecdsa
    - ed25519
container mounts:
  - source: /etc/resolv.conf
    dest: /etc/resolv.conf
    readonly: true
`,
	},
	{
		name: "v4.6.0",
		legacyYaml: `
ipaddr: 10.0.0.1
netmask: 255.255.252.0
network: 10.0.0.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: false
  datastore: /usr/share
dhcp:
  enabled: true
  range start: 10.0.1.1
  range end: 10.0.1.255
  systemd name: dhcpd
tftp:
  enabled: true
  systemd name: tftp
  ipxe:
    00:0B: arm64-efi/snponly.efi
    "00:00": undionly.kpxe
    "00:07": ipxe-snponly-x86_64.efi
    "00:09": ipxe-snponly-x86_64.efi
nfs:
  enabled: true
  export paths:
    - path: /home
      export options: rw,sync
      mount options: defaults
      mount: true
    - path: /opt
      export options: ro,sync,no_root_squash
      mount options: defaults
      mount: false
  systemd name: nfs-server
ssh:
  key types:
    - rsa
    - dsa
    - ecdsa
    - ed25519
container mounts:
  - source: /etc/resolv.conf
    dest: /etc/resolv.conf
    readonly: true
`,
		upgradedYaml: `
ipaddr: 10.0.0.1
netmask: 255.255.252.0
network: 10.0.0.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: false
dhcp:
  enabled: true
  range start: 10.0.1.1
  range end: 10.0.1.255
  systemd name: dhcpd
tftp:
  enabled: true
  systemd name: tftp
  ipxe:
    00:0B: arm64-efi/snponly.efi
    "00:00": undionly.kpxe
    "00:07": ipxe-snponly-x86_64.efi
    "00:09": ipxe-snponly-x86_64.efi
nfs:
  enabled: true
  export paths:
    - path: /home
      export options: rw,sync
    - path: /opt
      export options: ro,sync,no_root_squash
  systemd name: nfs-server
ssh:
  key types:
    - rsa
    - dsa
    - ecdsa
    - ed25519
container mounts:
  - source: /etc/resolv.conf
    dest: /etc/resolv.conf
    readonly: true
paths:
  datadir: /usr/share
`,
	},
}

func Test_UpgradeConfig(t *testing.T) {
	for _, tt := range configUpgradeTests {
		t.Run(tt.name, func(t *testing.T) {
			legacy, err := ParseConfig([]byte(tt.legacyYaml))
			assert.NoError(t, err)
			upgraded := legacy.Upgrade()
			upgradedYaml, err := upgraded.Dump()
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tt.upgradedYaml), strings.TrimSpace(string(upgradedYaml)))
		})
	}
}
