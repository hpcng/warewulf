package upgrade

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/warewulf/warewulf/internal/pkg/testenv"
)

var nodesYamlUpgradeTests = []struct {
	name            string
	addDefaults     bool
	replaceOverlays bool
	files           map[string]string
	legacyYaml      string
	upgradedYaml    string
	warewulfConf    string
}{
	{
		name:            "captured vers42 example",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
nodeprofiles:
  default:
    comment: This profile is automatically included for each node
    runtime overlay: "generic"
    discoverable: false
  leap:
    comment: openSUSE leap
    kernel version: 5.14.21
    ipmi netmask: "255.255.255.0"
    keys:
      foo: baar
    network devices:
      lan1:
        gateway: 1.1.1.1
nodes:
  node01:
    system overlay: "nodeoverlay"
    discoverable: true
    network devices:
      eth0:
        ipaddr: 1.2.3.4
        default: true
`,
		upgradedYaml: `
nodeprofiles:
  default:
    comment: This profile is automatically included for each node
    runtime overlay:
      - generic
  leap:
    comment: openSUSE leap
    kernel:
      version: 5.14.21
    ipmi:
      netmask: 255.255.255.0
    network devices:
      lan1:
        gateway: 1.1.1.1
    tags:
      foo: baar
nodes:
  node01:
    discoverable: "true"
    system overlay:
      - nodeoverlay
    network devices:
      eth0:
        ipaddr: 1.2.3.4
    primary network: eth0
`,
	},
	{
		name:            "captured vers43 example",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
WW_INTERNAL: 45
nodeprofiles:
  default:
    comment: This profile is automatically included for each node
    runtime overlay:
    - generic
    discoverable: "false"
  leap:
    comment: openSUSE leap
    kernel:
      override: 5.14.21
    ipmi:
      netmask: 255.255.255.0
    discoverable: "false"
    network devices:
      lan1:
        device: lan1
        gateway: 1.1.1.1
        default: "false"
    keys:
      foo: baar
nodes:
  node01:
    system overlay:
    - nodeoverlay
    discoverable: "true"
    network devices:
      eth0:
        device: eth0
        ipaddr: 1.2.3.4
        default: "true"
`,
		upgradedYaml: `
nodeprofiles:
  default:
    comment: This profile is automatically included for each node
    runtime overlay:
      - generic
  leap:
    comment: openSUSE leap
    ipmi:
      netmask: 255.255.255.0
    network devices:
      lan1:
        device: lan1
        gateway: 1.1.1.1
    tags:
      foo: baar
nodes:
  node01:
    discoverable: "true"
    system overlay:
      - nodeoverlay
    network devices:
      eth0:
        device: eth0
        ipaddr: 1.2.3.4
    primary network: eth0
`,
	},
	{
		name:            "remove WW_INTERNAL",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml:      `WW_INTERNAL: 45`,
		upgradedYaml: `
nodeprofiles: {}
nodes: {}
`,
	},
	{
		name:            "disabled is obsolete",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
nodes:
  n1:
    disabled: true
nodeprofiles:
  default:
    disabled: true
`,
		upgradedYaml: `
nodeprofiles:
  default: {}
nodes:
  n1: {}
`,
	},
	{
		name:            "inline IPMI settings",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
nodes:
  n1:
    ipmi escapechar: "~"
    ipmi gateway: 192.168.0.1
    ipmi interface: lanplus
    ipmi ipaddr: 192.168.0.100
    ipmi netmask: 255.255.255.0
    ipmi password: password
    ipmi port: 623
    ipmi username: admin
    ipmi write: true
nodeprofiles:
  default:
    ipmi escapechar: "~"
    ipmi gateway: 192.168.0.1
    ipmi interface: lanplus
    ipmi ipaddr: 192.168.0.100
    ipmi netmask: 255.255.255.0
    ipmi password: password
    ipmi port: 623
    ipmi username: admin
    ipmi write: true
`,
		upgradedYaml: `
nodeprofiles:
  default:
    ipmi:
      username: admin
      password: password
      ipaddr: 192.168.0.100
      gateway: 192.168.0.1
      netmask: 255.255.255.0
      port: "623"
      interface: lanplus
      escapechar: "~"
      write: "true"
nodes:
  n1:
    ipmi:
      username: admin
      password: password
      ipaddr: 192.168.0.100
      gateway: 192.168.0.1
      netmask: 255.255.255.0
      port: "623"
      interface: lanplus
      escapechar: "~"
      write: "true"
`,
	},
	{
		name:            "inline Kernel settings",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
nodeprofiles:
  default:
    kernel args: quiet
    kernel override: rockylinux-9
    kernel version: 2.6
nodes:
  n1:
    kernel args: quiet
    kernel override: rockylinux-9
    kernel version: 2.6
`,
		upgradedYaml: `
nodeprofiles:
  default:
    kernel:
      version: "2.6"
      args: quiet
nodes:
  n1:
    kernel:
      version: "2.6"
      args: quiet
`,
	},
	{
		name:            "keys and tags",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
nodeprofiles:
  default:
    keys:
      key1: val1
      key2: val2
    tags:
      key2: valB
      key3: valC
      key4: valD
    tagsdel:
      - key4
    network devices:
      default:
        tags:
          key2: valB
          key3: valC
          key4: valD
        tagsdel:
          - key4
    ipmi:
      tags:
        key2: valB
        key3: valC
        key4: valD
      tagsdel:
        - key4
nodes:
  n1:
    keys:
      key1: val1
      key2: val2
    tags:
      key2: valB
      key3: valC
      key4: valD
    tagsdel:
      - key4
    network devices:
      default:
        tags:
          key2: valB
          key3: valC
          key4: valD
        tagsdel:
          - key4
    ipmi:
      tags:
        key2: valB
        key3: valC
        key4: valD
      tagsdel:
        - key4
`,
		upgradedYaml: `
nodeprofiles:
  default:
    ipmi:
      tags:
        key2: valB
        key3: valC
    network devices:
      default:
        tags:
          key2: valB
          key3: valC
    tags:
      key1: val1
      key2: valB
      key3: valC
nodes:
  n1:
    ipmi:
      tags:
        key2: valB
        key3: valC
    network devices:
      default:
        tags:
          key2: valB
          key3: valC
    tags:
      key1: val1
      key2: valB
      key3: valC
`,
	},
	{
		name:            "primary network",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
nodes:
  n1:
    network devices:
      eth0: {}
      eth1:
        default: true
  n2:
    network devices:
      eth0:
        primary: true
      eth1: {}
  n3:
    network devices:
      eth0:
        primary: true
      eth1: {}
    primary network: eth1
nodeprofiles:
  p1:
    network devices:
      eth0: {}
      eth1:
        default: true
  p2:
    network devices:
      eth0:
        primary: true
      eth1: {}
  p3:
    network devices:
      eth0:
        primary: true
      eth1: {}
    primary network: eth1
`,
		upgradedYaml: `
nodeprofiles:
  p1:
    primary network: eth1
  p2:
    primary network: eth0
  p3:
    primary network: eth1
nodes:
  n1:
    primary network: eth1
  n2:
    primary network: eth0
  n3:
    primary network: eth1
`,
	},
	{
		name:            "overlays",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
nodes:
  n1:
    runtime overlay:
    - r1
    - r2
    system overlay:
    - s1
    - s2
  n2:
    runtime overlay: r1,r2
    system overlay: s1,s2
nodeprofiles:
  p1:
    runtime overlay:
    - r1
    - r2
    system overlay:
    - s1
    - s2
  p2:
    runtime overlay: r1,r2
    system overlay: s1,s2
`,
		upgradedYaml: `
nodeprofiles:
  p1:
    runtime overlay:
      - r1
      - r2
    system overlay:
      - s1
      - s2
  p2:
    runtime overlay:
      - r1
      - r2
    system overlay:
      - s1
      - s2
nodes:
  n1:
    runtime overlay:
      - r1
      - r2
    system overlay:
      - s1
      - s2
  n2:
    runtime overlay:
      - r1
      - r2
    system overlay:
      - s1
      - s2
`,
	},
	{
		name:            "disk example",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
nodes:
  n1:
    disks:
      /dev/vda:
        wipe_table: true
        partitions:
          scratch:
            number: "1"
            should_exist: true
          swap:
            number: "2"
            size_mib: "1024"
    filesystems:
      /dev/disk/by-partlabel/scratch:
        format: btrfs
        path: /scratch
      /dev/disk/by-partlabel/swap:
        format: swap
        path: swap
`,
		upgradedYaml: `
nodeprofiles: {}
nodes:
  n1:
    disks:
      /dev/vda:
        wipe_table: true
        partitions:
          scratch:
            number: "1"
            should_exist: true
          swap:
            number: "2"
            size_mib: "1024"
    filesystems:
      /dev/disk/by-partlabel/scratch:
        format: btrfs
        path: /scratch
      /dev/disk/by-partlabel/swap:
        format: swap
        path: swap
`,
	},
	{
		name:            "add defaults",
		addDefaults:     true,
		replaceOverlays: false,
		legacyYaml: `
nodes:
  n1:
    network devices:
      default:
        ipaddr: 192.168.0.100
`,
		upgradedYaml: `
nodeprofiles:
  default:
    ipxe template: default
    runtime overlay:
      - hosts
      - ssh.authorized_keys
      - syncuser
    system overlay:
      - wwinit
      - wwclient
      - fstab
      - hostname
      - ssh.host_keys
      - issue
      - resolv
      - udev.netname
      - systemd.netname
      - ifcfg
      - NetworkManager
      - debian.interfaces
      - wicked
      - ignition
    kernel:
      args: quiet crashkernel=no vga=791 net.naming-scheme=v238
    init: /sbin/init
    root: initramfs
nodes:
  n1:
    profiles:
      - default
    network devices:
      default:
        type: ethernet
        ipaddr: 192.168.0.100
        netmask: 255.255.255.0
`,
	},
	{
		name:            "add defaults conflicts",
		addDefaults:     true,
		replaceOverlays: false,
		legacyYaml: `
nodeprofiles:
  default:
    system overlay:
      - wwinit
      - wwclient
      - fstab
      - hostname
      - ssh.host_keys
      - issue
      - resolv
      - udev.netname
      - systemd.netname
      - NetworkManager
  custom: {}
nodes:
  n1:
    profiles:
      - custom
    network devices:
      default:
        ipaddr: 10.0.0.100
        netmask: 255.255.0.0
`,
		upgradedYaml: `
nodeprofiles:
  custom: {}
  default:
    ipxe template: default
    runtime overlay:
      - hosts
      - ssh.authorized_keys
      - syncuser
    system overlay:
      - wwinit
      - wwclient
      - fstab
      - hostname
      - ssh.host_keys
      - issue
      - resolv
      - udev.netname
      - systemd.netname
      - NetworkManager
    kernel:
      args: quiet crashkernel=no vga=791 net.naming-scheme=v238
    init: /sbin/init
    root: initramfs
nodes:
  n1:
    profiles:
      - custom
    network devices:
      default:
        type: ethernet
        ipaddr: 10.0.0.100
        netmask: 255.255.0.0
`,
	},
	{
		name:            "add defaults conflicts",
		addDefaults:     false,
		replaceOverlays: true,
		legacyYaml: `
nodeprofiles:
  default:
    runtime overlay:
      - generic
    system overlay:
      - wwinit
nodes:
  n1:
    runtime overlay:
      - generic
    system overlay:
      - wwinit
`,
		upgradedYaml: `
nodeprofiles:
  default:
    runtime overlay:
      - hosts
      - ssh.authorized_keys
      - syncuser
    system overlay:
      - wwinit
      - wwclient
      - fstab
      - hostname
      - ssh.host_keys
      - issue
      - resolv
      - udev.netname
      - systemd.netname
      - ifcfg
      - NetworkManager
      - debian.interfaces
      - wicked
      - ignition
nodes:
  n1:
    runtime overlay:
      - hosts
      - ssh.authorized_keys
      - syncuser
    system overlay:
      - wwinit
      - wwclient
      - fstab
      - hostname
      - ssh.host_keys
      - issue
      - resolv
      - udev.netname
      - systemd.netname
      - ifcfg
      - NetworkManager
      - debian.interfaces
      - wicked
      - ignition
`,
	},
	{
		name:            "Kernel.Override (legacy)",
		addDefaults:     false,
		replaceOverlays: false,
		files: map[string]string{
			"/srv/warewulf/kernel/mykernel/version":                           "1.2.3",
			"/var/lib/warewulf/chroots/mycontainer/rootfs/boot/vmlinuz-1.2.3": "",
		},
		legacyYaml: `
nodeprofiles:
  default:
    container name: mycontainer
    kernel:
      override: mykernel
nodes:
  n1:
    container name: mycontainer
    kernel:
      override: mykernel
`,
		upgradedYaml: `
nodeprofiles:
  default:
    container name: mycontainer
    kernel:
      version: /boot/vmlinuz-1.2.3
nodes:
  n1:
    container name: mycontainer
    kernel:
      version: /boot/vmlinuz-1.2.3
`,
	},
	{
		name:            "Kernel.Override (upgraded)",
		addDefaults:     false,
		replaceOverlays: false,
		files: map[string]string{
			"/srv/warewulf/kernel/mykernel/version":                           "1.2.3",
			"/var/lib/warewulf/chroots/mycontainer/rootfs/boot/vmlinuz-1.2.3": "",
		},
		legacyYaml: `
nodeprofiles:
  default:
    container name: mycontainer
    kernel:
      override: /boot/vmlinuz-1.2.3
nodes:
  n1:
    container name: mycontainer
    kernel:
      override: /boot/vmlinuz-1.2.3
`,
		upgradedYaml: `
nodeprofiles:
  default:
    container name: mycontainer
    kernel:
      version: /boot/vmlinuz-1.2.3
nodes:
  n1:
    container name: mycontainer
    kernel:
      version: /boot/vmlinuz-1.2.3
`,
	},
	{
		name:            "Nested profiles",
		addDefaults:     false,
		replaceOverlays: false,
		legacyYaml: `
nodeprofiles:
  p1:
    profiles:
      - p2
  p2: {}
`,
		upgradedYaml: `
nodeprofiles:
  p1:
    profiles:
      - p2
  p2: {}
nodes: {}
`,
	},
	{
		name: "Legacy export mounts",
		legacyYaml: `
nodeprofiles:
  default: {}
`,
		upgradedYaml: `
nodeprofiles:
  default:
    resources:
      fstab:
        - spec: warewulf:/home
          file: /home
          vfstype: nfs
        - spec: warewulf:/opt
          file: /opt
          vfstype: nfs
nodes: {}
`,
		warewulfConf: `
nfs:
  exports:
  - /home
  - /opt
`,
	},
	{
		name: "Legacy extended export mounts",
		legacyYaml: `
nodeprofiles:
  default: {}
`,
		upgradedYaml: `
nodeprofiles:
  default:
    resources:
      fstab:
        - spec: warewulf:/home
          file: /home
          vfstype: nfs
        - spec: warewulf:/opt
          file: /opt
          mntops: defaults,ro
          vfstype: nfs
nodes: {}
`,
		warewulfConf: `
nfs:
  export paths:
  - path: /home
    mount: true
  - path: /opt
    mount options: defaults,ro
    mount: true
  - path: /var
    mount: false
`,
	},
}

func Test_UpgradeNodesYaml(t *testing.T) {
	for _, tt := range nodesYamlUpgradeTests {
		t.Run(tt.name, func(t *testing.T) {
			env := testenv.New(t)
			defer env.RemoveAll()
			if tt.files != nil {
				for fileName, content := range tt.files {
					env.WriteFile(fileName, content)
				}
			}
			conf, err := ParseConfig([]byte(tt.warewulfConf))
			assert.NoError(t, err)
			legacy, err := ParseNodes([]byte(tt.legacyYaml))
			assert.NoError(t, err)
			upgraded := legacy.Upgrade(tt.addDefaults, tt.replaceOverlays, conf)
			upgradedYaml, err := upgraded.Dump()
			assert.NoError(t, err)
			assert.YAMLEq(t, tt.upgradedYaml, string(upgradedYaml))
		})
	}
}
