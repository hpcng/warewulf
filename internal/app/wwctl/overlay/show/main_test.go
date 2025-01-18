package show

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/warewulfd"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

var (
	overlayEmail = `
{{ if .Tags.email }}eMail: {{ .Tags.email }}{{else}} noMail{{- end }}
`
	overlayOverlay = `
overlay name {{ .Overlay }}
`
)

func Test_Overlay_List(t *testing.T) {
	env := testenv.New(t)
	defer env.RemoveAll()

	env.WriteFile("etc/warewulf/warewulf.conf", `ipaddr: 192.168.0.1/24
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
    export options: ro,sync,no_root_squash`)

	env.WriteFile("etc/warewulf/nodes.conf",
		`nodeprofiles:
  default:
    tags:
      email: admin@localhost
  empty: {}
nodes:
  node1:
    tags:
      email: admin@node1
  node2:
    profiles:
      - default
  node3:
    profiles:
      - empty
`)

	env.WriteFile("usr/share/warewulf/overlays/testoverlay/email.ww", overlayEmail)
	env.WriteFile("usr/share/warewulf/overlays/testoverlay/overlay.ww", overlayOverlay)
	env.WriteFile("usr/share/warewulf/overlays/dist/foo.ww", "foo")
	env.WriteFile("var/lib/warewulf/overlays/dist/foo.ww", "foobaar")

	warewulfd.SetNoDaemon()
	t.Run("overlay show raw", func(t *testing.T) {
		baseCmd.SetArgs([]string{"testoverlay", "email.ww"})
		baseCmd := GetCommand()
		buf := new(bytes.Buffer)
		baseCmd.SetOut(buf)
		baseCmd.SetErr(buf)
		wwlog.SetLogWriter(buf)
		err := baseCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), overlayEmail)
	})
	t.Run("overlay show rendered node tag", func(t *testing.T) {
		baseCmd.SetArgs([]string{"-r", "node1", "testoverlay", "email.ww"})
		baseCmd := GetCommand()
		buf := new(bytes.Buffer)
		baseCmd.SetOut(buf)
		baseCmd.SetErr(buf)
		wwlog.SetLogWriter(buf)
		err := baseCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "admin@node1")
	})
	t.Run("overlay show rendered profile tag", func(t *testing.T) {
		baseCmd.SetArgs([]string{"-r", "node2", "testoverlay", "email.ww"})
		baseCmd := GetCommand()
		buf := new(bytes.Buffer)
		baseCmd.SetOut(buf)
		baseCmd.SetErr(buf)
		wwlog.SetLogWriter(buf)
		err := baseCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "admin@localhost")
	})
	t.Run("overlay show no tag", func(t *testing.T) {
		baseCmd.SetArgs([]string{"-r", "node3", "testoverlay", "email.ww"})
		baseCmd := GetCommand()
		buf := new(bytes.Buffer)
		baseCmd.SetOut(buf)
		baseCmd.SetErr(buf)
		wwlog.SetLogWriter(buf)
		err := baseCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "noMail")
	})
	t.Run("overlay shows overlay", func(t *testing.T) {
		baseCmd.SetArgs([]string{"-r", "node1", "testoverlay", "overlay.ww"})
		baseCmd := GetCommand()
		buf := new(bytes.Buffer)
		baseCmd.SetOut(buf)
		baseCmd.SetErr(buf)
		wwlog.SetLogWriter(buf)
		err := baseCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "testoverlay")
	})
	t.Run("site overlays precede", func(t *testing.T) {
		baseCmd.SetArgs([]string{"-r", "node1", "dist", "foo.ww"})
		baseCmd := GetCommand()
		buf := new(bytes.Buffer)
		baseCmd.SetOut(buf)
		baseCmd.SetErr(buf)
		wwlog.SetLogWriter(buf)
		err := baseCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "foobaar")
	})
}

func TestShowServerTemplate(t *testing.T) {
	const template = `
	Id: {{.Id}}
	ClusterName: {{.ClusterName}}
	BuildHost: {{.BuildHost}}
	`

	env := testenv.New(t)
	env.WriteFile("etc/warewulf/nodes.conf",
		`nodeprofiles:
  default:
    tags:
      email: admin@localhost
  empty: {}
nodes:
  node1:
    tags:
      email: admin@node1
  node2: {}
  node3:
    profiles:
      - empty
`)

	env.WriteFile(path.Join(testenv.WWOverlaydir, "testoverlay/template.ww"), template)
	defer env.RemoveAll()
	warewulfd.SetNoDaemon()

	host, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("overlay render host template using 'host' value", func(t *testing.T) {
		baseCmd.SetArgs([]string{"-r", "host", "testoverlay", "template.ww"})
		baseCmd := GetCommand()
		buf := new(bytes.Buffer)
		baseCmd.SetOut(buf)
		baseCmd.SetErr(buf)
		wwlog.SetLogWriter(buf)
		err := baseCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Id: "+host)
		assert.Contains(t, buf.String(), "ClusterName: "+host)
		assert.Contains(t, buf.String(), "BuildHost: "+host)
	})

	t.Run("overlay render host template using host domain value", func(t *testing.T) {
		host, err := os.Hostname()
		if err != nil {
			t.Fatal(err)
		}
		baseCmd.SetArgs([]string{"-r", host, "testoverlay", "template.ww"})
		baseCmd := GetCommand()
		buf := new(bytes.Buffer)
		baseCmd.SetOut(buf)
		baseCmd.SetErr(buf)
		wwlog.SetLogWriter(buf)
		err = baseCmd.Execute()
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Id: "+host)
		assert.Contains(t, buf.String(), "ClusterName: "+host)
		assert.Contains(t, buf.String(), "BuildHost: "+host)
	})
}
