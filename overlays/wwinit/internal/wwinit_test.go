package wwinit

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/app/wwctl/overlay/show"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_wwinitOverlay(t *testing.T) {
	env := testenv.New(t)
	defer env.RemoveAll()
	env.ImportFile("etc/warewulf/warewulf.conf", "warewulf.conf")
	env.ImportFile("etc/warewulf/nodes.conf", "nodes.conf")
	env.ImportFile("var/lib/warewulf/overlays/wwinit/rootfs/etc/warewulf/warewulf.conf.ww", "../rootfs/etc/warewulf/warewulf.conf.ww")
	env.ImportFile("var/lib/warewulf/overlays/wwinit/rootfs/warewulf/config.ww", "../rootfs/warewulf/config.ww")

	tests := []struct {
		name string
		args []string
		log  string
	}{
		{
			name: "wwinit:warewulf.conf.ww",
			args: []string{"--render", "node1", "wwinit", "etc/warewulf/warewulf.conf.ww"},
			log:  wwinit_warewulf_conf,
		},
		{
			name: "wwinit:config.ww",
			args: []string{"--render", "node1", "wwinit", "warewulf/config.ww"},
			log:  wwinit_config,
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
			assert.Equal(t, tt.log, logbuf.String())
		})
	}
}

const wwinit_warewulf_conf string = `backupFile: true
writeFile: true
Filename: etc/warewulf/warewulf.conf
ipaddr: 192.168.0.1/24
netmask: 255.255.255.0
network: 192.168.0.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: true
dhcp:
  enabled: true
  range start: 192.168.0.100
  range end: 192.168.0.199
tftp:
  enabled: false
nfs:
  enabled: true
  export paths:
  - path: /home
    export options: rw,sync
  - path: /opt
    export options: ro,sync,no_root_squash
`

const wwinit_config string = `backupFile: true
writeFile: true
Filename: warewulf/config
WWCONTAINER=rockylinux-9
WWHOSTNAME=node1
WWROOT=initramfs
WWINIT=/sbin/init
WWIPMI_IPADDR="192.168.4.21"
WWIPMI_NETMASK="255.255.255.0"
WWIPMI_GATEWAY="192.168.4.1"
WWIPMI_USER="user"
WWIPMI_PASSWORD="password"
WWIPMI_WRITE="true"
`
