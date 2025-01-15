package debug

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/app/wwctl/overlay/show"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_debugOverlay(t *testing.T) {
	variableData := regexp.MustCompile(`(?m)(BuildTime|BuildTimeUnix|BuildSource|DataStore):.*$`)
	hostname, _ := os.Hostname()

	env := testenv.New(t)
	defer env.RemoveAll()
	env.ImportFile("etc/warewulf/nodes.conf", "nodes.conf")
	env.ImportFile("var/lib/warewulf/overlays/debug/rootfs/warewulf/template-variables.md.ww", "../rootfs/warewulf/template-variables.md.ww")

	tests := []struct {
		name string
		args []string
		log  string
	}{
		{
			name: "debug",
			args: []string{"--render", "node1", "debug", "warewulf/template-variables.md.ww"},
			log:  debug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := show.GetCommand()
			cmd.SetArgs(tt.args)
			stdout := bytes.NewBufferString("")
			stderr := bytes.NewBufferString("")
			logbuf := bytes.NewBufferString("")
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			wwlog.SetLogWriter(logbuf)
			err := cmd.Execute()
			assert.NoError(t, err)
			assert.Empty(t, stdout.String())
			assert.Empty(t, stderr.String())
			assert.Equal(t, strings.Replace(tt.log, "%HOSTNAME%", hostname, -1), variableData.ReplaceAllString(logbuf.String(), "${1}: REMOVED_BY_TEST"))
		})
	}
}

const debug = `backupFile: true
writeFile: true
Filename: warewulf/template-variables.md
# Warewulf template variables

This Warewulf template serves as a complete example of the variables
available to Warewulf templates. It may also be rendered against a
node to debug its apparent configuration.

    wwctl overlay show --render $nodename debug /warewulf/template-variables.md.ww

The template data structure is defined in
internal/pkg/overlay/datastructure.go, though it also references
data from other structures.


## Node

- Id: node1
- Hostname: node1
- Comment: 
- ClusterName: 
- ContainerName: rockylinux-9
- Ipxe: default
- RuntimeOverlay: hosts,ssh.authorized_keys,syncuser
- SystemOverlay: wwinit,wwclient,fstab,hostname,ssh.host_keys,issue,resolv,udev.netname,systemd.netname,ifcfg,NetworkManager,debian.interfaces,wicked,ignition
- Init: /sbin/init
- Root: initramfs
- AssetKey: 
- Discoverable: 
- Profiles: default
- Tags: 

- Kernel:
  - Version: 
  - Args: quiet crashkernel=no vga=791 net.naming-scheme=v238

- Ipmi:
  - UserName: user
  - Password: password
  - Ipaddr: 192.168.4.21
  - Netmask: 255.255.255.0
  - Port: 
  - Gateway: 192.168.4.1
  - Interface: 
  - Write: true
  - Tags: 
- NetDevs[default]:
  - Type: 
  - OnBoot: true
  - Device: wwnet0
  - Hwaddr: e6:92:39:49:7b:03
  - Ipaddr: 192.168.3.21
  - Ipaddr6: <nil>
  - Prefix: <nil>
  - Netmask: 255.255.255.0
  - Gateway: 192.168.3.1
  - MTU: 
  - Primary: false
  - Tags: 
- NetDevs[secondary]:
  - Type: 
  - OnBoot: true
  - Device: wwnet1
  - Hwaddr: 9a:77:29:73:14:f1
  - Ipaddr: 192.168.3.22
  - Ipaddr6: <nil>
  - Prefix: <nil>
  - Netmask: 255.255.255.0
  - Gateway: 192.168.3.1
  - MTU: 
  - Primary: false
  - Tags: DNS1=8.8.8.8 DNS2=8.8.4.4 


## Build variables

- BuildHost: %HOSTNAME%
- BuildTime: REMOVED_BY_TEST
- BuildTimeUnix: REMOVED_BY_TEST
- BuildSource: REMOVED_BY_TEST
- Overlay: debug


## Network

- Ipaddr: 
- Ipaddr6: 
- Netmask: 
- Network: 
- NetworkCIDR: 
- Ipv6: false


## Services

### DHCP

- Dhcp.Enabled: true
- Dhcp.Template: default
- Dhcp.RangeStart: 
- Dhcp.RangeEnd: 
- Dhcp.SystemdName: dhcpd


### NFS

- Enabled: true
- SystemdName: nfsd

### SSH
- Key types:
  - rsa
  - dsa
  - ecdsa
  - ed25519
- First key type: rsa

### Warewulf

- Port: 9873
- Secure: true
- UpdateInterval: 60
- AutobuildOverlays: true
- EnableHostOverlay: true
- Syslog: false


### Other nodes

- AllNodes[0]:
  - Id: node1
  - Comment: 
  - ClusterName: 
  - ContainerName: rockylinux-9
  - Ipxe: default
  - RuntimeOverlay: [hosts ssh.authorized_keys syncuser]
  - SystemOverlay: [wwinit wwclient fstab hostname ssh.host_keys issue resolv udev.netname systemd.netname ifcfg NetworkManager debian.interfaces wicked ignition]
  - Root: initramfs
  - Discoverable: 
  - Init: /sbin/init
  - AssetKey: 
  - Profiles: [default]
  - Tags: 

  - Kernel
    - Version: 
    - Args: quiet crashkernel=no vga=791 net.naming-scheme=v238

  - Ipmi:
    - Ipaddr: 192.168.4.21
    - Netmask: 255.255.255.0
    - Port: 
    - Gateway: 192.168.4.1
    - UserName: user
    - Password: password
    - Interface: 
    - Write: true
    - Tags: 
  - NetDevs[default]:
    - Type: 
    - OnBoot: true
    - Device: wwnet0
    - Hwaddr: e6:92:39:49:7b:03
    - Ipaddr: 192.168.3.21
    - IpCIDR: 192.168.3.21/24
    - Ipaddr6: <nil>
    - Prefix: <nil>
    - Netmask: 255.255.255.0
    - Gateway: 192.168.3.1
    - MTU: 
    - Primary: true
    - Tags: 
  - NetDevs[secondary]:
    - Type: 
    - OnBoot: true
    - Device: wwnet1
    - Hwaddr: 9a:77:29:73:14:f1
    - Ipaddr: 192.168.3.22
    - IpCIDR: 192.168.3.22/24
    - Ipaddr6: <nil>
    - Prefix: <nil>
    - Netmask: 255.255.255.0
    - Gateway: 192.168.3.1
    - MTU: 
    - Primary: false
    - Tags: DNS1=8.8.8.8 DNS2=8.8.4.4 
- AllNodes[1]:
  - Id: node2
  - Comment: 
  - ClusterName: 
  - ContainerName: 
  - Ipxe: default
  - RuntimeOverlay: [hosts ssh.authorized_keys syncuser]
  - SystemOverlay: [wwinit wwclient fstab hostname ssh.host_keys issue resolv udev.netname systemd.netname ifcfg NetworkManager debian.interfaces wicked ignition]
  - Root: initramfs
  - Discoverable: 
  - Init: /sbin/init
  - AssetKey: 
  - Profiles: [default]
  - Tags: 

  - Kernel
    - Version: 
    - Args: quiet crashkernel=no vga=791 net.naming-scheme=v238

  - NetDevs[default]:
    - Type: 
    - OnBoot: true
    - Device: wwnet0
    - Hwaddr: e6:92:39:49:7b:04
    - Ipaddr: 192.168.3.23
    - IpCIDR: 192.168.3.23/24
    - Ipaddr6: <nil>
    - Prefix: <nil>
    - Netmask: 255.255.255.0
    - Gateway: 192.168.3.1
    - MTU: 
    - Primary: true
    - Tags: 

`
