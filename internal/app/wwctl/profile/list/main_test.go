package list

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/warewulfd"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_List(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		stdout string
		inDb   string
	}{
		{
			name: "profile list test",
			args: []string{},
			stdout: `
PROFILE NAME  COMMENT/DESCRIPTION
------------  -------------------
default       --
`,
			inDb: `
nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name: "profile list returns multiple profiles",
			args: []string{"default,test"},
			stdout: `
PROFILE NAME  COMMENT/DESCRIPTION
------------  -------------------
default       --
test          --
`,
			inDb: `
nodeprofiles:
  default: {}
  test: {}
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name: "profile list returns one profile",
			args: []string{"test,"},
			stdout: `
PROFILE NAME  COMMENT/DESCRIPTION
------------  -------------------
test          --
`,
			inDb: `
nodeprofiles:
  default: {}
  test: {}
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name: "profile list returns all profiles",
			args: []string{","},
			stdout: `
PROFILE NAME  COMMENT/DESCRIPTION
------------  -------------------
default       --
test          --
`,
			inDb: `
nodeprofiles:
  default: {}
  test: {}
nodes:
  n01:
    profiles:
    - default
`,
		},
	}

	warewulfd.SetNoDaemon()
	env := testenv.New(t)
	defer env.RemoveAll()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env.WriteFile("etc/warewulf/nodes.conf", tt.inDb)
			baseCmd := GetCommand()
			baseCmd.SetArgs(tt.args)
			stdoutR, stdoutW, _ := os.Pipe()
			oriout := os.Stdout
			os.Stdout = stdoutW
			wwlog.SetLogWriter(os.Stdout)
			baseCmd.SetOut(os.Stdout)
			baseCmd.SetErr(os.Stdout)
			err := baseCmd.Execute()
			assert.NoError(t, err)

			stdoutC := make(chan string)
			go func() {
				var buf bytes.Buffer
				_, _ = io.Copy(&buf, stdoutR)
				stdoutC <- buf.String()
			}()
			stdoutW.Close()
			os.Stdout = oriout

			stdout := <-stdoutC
			assert.Equal(t, strings.TrimSpace(tt.stdout), strings.TrimSpace(stdout))
		})
	}
}

func TestListMultipleFormats(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		output string
		inDb   string
	}{
		{
			name:   "single profile list yaml output",
			args:   []string{"-y"},
			output: `default: {}`,
			inDb: `
nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name: "single profile list json output",
			args: []string{"-j"},
			output: `
{
  "default": {
    "Profiles": null,
    "Comment": "",
    "ClusterName": "",
    "ContainerName": "",
    "Ipxe": "",
    "RuntimeOverlay": null,
    "SystemOverlay": null,
    "Kernel": null,
    "Ipmi": null,
    "Init": "",
    "Root": "",
    "NetDevs": null,
    "Tags": null,
    "PrimaryNetDev": "",
    "Disks": null,
    "FileSystems": null,
    "Resources": null
  }
}
`,
			inDb: `
nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name: "multiple profiles list yaml output",
			args: []string{"-y"},
			output: `
default: {}
test: {}
`,
			inDb: `
nodeprofiles:
  default: {}
  test: {}
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name: "multiple profiles list json output",
			args: []string{"-j"},
			output: `
{
  "default": {
    "Profiles": null,
    "Comment": "",
    "ClusterName": "",
    "ContainerName": "",
    "Ipxe": "",
    "RuntimeOverlay": null,
    "SystemOverlay": null,
    "Kernel": null,
    "Ipmi": null,
    "Init": "",
    "Root": "",
    "NetDevs": null,
    "Tags": null,
    "PrimaryNetDev": "",
    "Disks": null,
    "FileSystems": null,
    "Resources": null
  },
  "test": {
    "Profiles": null,
    "Comment": "",
    "ClusterName": "",
    "ContainerName": "",
    "Ipxe": "",
    "RuntimeOverlay": null,
    "SystemOverlay": null,
    "Kernel": null,
    "Ipmi": null,
    "Init": "",
    "Root": "",
    "NetDevs": null,
    "Tags": null,
    "PrimaryNetDev": "",
    "Disks": null,
    "FileSystems": null,
    "Resources": null
  }
}
`,
			inDb: `
nodeprofiles:
  default: {}
  test: {}
nodes:
  n01:
    profiles:
    - default
`,
		},
	}

	warewulfd.SetNoDaemon()
	env := testenv.New(t)
	defer env.RemoveAll()

	for _, tt := range tests {
		env.WriteFile("etc/warewulf/nodes.conf", tt.inDb)

		t.Run(tt.name, func(t *testing.T) {
			baseCmd := GetCommand()
			baseCmd.SetArgs(tt.args)

			buf := new(bytes.Buffer)
			baseCmd.SetOut(buf)
			baseCmd.SetErr(buf)
			wwlog.SetLogWriter(buf)
			err := baseCmd.Execute()
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tt.output), strings.TrimSpace(buf.String()))
		})
	}
}
