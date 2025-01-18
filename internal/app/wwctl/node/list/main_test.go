package list

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/warewulfd"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_List(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		stdout  string
		inDb    string
	}{
		{
			name:    "single node list",
			args:    []string{},
			wantErr: false,
			stdout: `
NODE NAME  PROFILES  NETWORK
---------  --------  -------
n01        default   --
`,
			inDb: `nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name:    "multiple nodes list",
			args:    []string{},
			wantErr: false,
			stdout: `
NODE NAME  PROFILES  NETWORK
---------  --------  -------
n01        default   --
n02        default   --
`,
			inDb: `nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
  n02:
   profiles:
   - default
`,
		},
		{
			name:    "node list returns multiple nodes",
			args:    []string{"n01,n02"},
			wantErr: false,
			stdout: `
NODE NAME  PROFILES  NETWORK
---------  --------  -------
n01        default   --
n02        default   --
`,
			inDb: `nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
  n02:
   profiles:
   - default
`,
		},
		{
			name:    "node list returns multiple nodes (case 2)",
			args:    []string{"n01,n03"},
			wantErr: false,
			stdout: `
NODE NAME  PROFILES  NETWORK
---------  --------  -------
n01        default   --
n03        default   --
`,
			inDb: `nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
  n02:
   profiles:
   - default
  n03:
   profiles:
   - default
  n04:
   profiles:
   - default
  n05:
   profiles:
   - default
`,
		},
		{
			name:    "node list returns one node",
			args:    []string{"n01,"},
			wantErr: false,
			stdout: `
NODE NAME  PROFILES  NETWORK
---------  --------  -------
n01        default   --
`,
			inDb: `nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
  n02:
   profiles:
   - default
`,
		},
		{
			name:    "node list profile with network",
			args:    []string{},
			wantErr: false,
			stdout: `
NODE NAME  PROFILES  NETWORK
---------  --------  -------
n01        default   default
`,
			inDb: `nodeprofiles:
  default:
    network devices:
      default:
        device: eth0
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name:    "node list profile with comment",
			args:    []string{"-a"},
			wantErr: false,
			stdout: `
NODE  FIELD     PROFILE  VALUE
----  -----     -------  -----
n01   Profiles  --       default
n01   Comment   default  profilecomment
`,
			inDb: `nodeprofiles:
  default:
    comment: profilecomment
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name:    "node list profile with comment superseded",
			args:    []string{"-a"},
			wantErr: false,
			stdout: `
NODE  FIELD     PROFILE     VALUE
----  -----     -------     -----
n01   Profiles  --          default
n01   Comment   SUPERSEDED  nodecomment
`,
			inDb: `nodeprofiles:
  default:
    comment: profilecomment
nodes:
  n01:
    comment: nodecomment
    profiles:
    - default
`,
		},
		{
			name:    "node list profile with ipmi user",
			args:    []string{"-i"},
			wantErr: false,
			stdout: `
NODE  IPMI IPADDR  IPMI PORT  IPMI USERNAME  IPMI INTERFACE
----  -----------  ---------  -------------  --------------
n01   <nil>        --         admin          --
`,
			inDb: `nodeprofiles:
  default:
    ipmi:
      username: admin
nodes:
  n01:
    profiles:
    - default
`,
		},
		{
			name:    "node list profile with ipmi user superseded",
			args:    []string{"-i"},
			wantErr: false,
			stdout: `
NODE  IPMI IPADDR  IPMI PORT  IPMI USERNAME  IPMI INTERFACE
----  -----------  ---------  -------------  --------------
n01   <nil>        --         user           --
`,
			inDb: `nodeprofiles:
  default:
    ipmi:
      username: admin
nodes:
  n01:
    ipmi:
      username: user
    profiles:
    - default
`,
		},
		{
			name:    "multiple profiles list",
			args:    []string{},
			wantErr: false,
			stdout: `
NODE NAME  PROFILES  NETWORK
---------  --------  -------
n01        p1,p2     --
`,
			inDb: `nodeprofiles:
  p1: {}
  p2: {}
nodes:
  n01:
    profiles:
    - p1
    - p2
`,
		},
		{
			name:    "multiple profiles list all",
			args:    []string{"-a"},
			wantErr: false,
			stdout: `
NODE  FIELD     PROFILE  VALUE
----  -----     -------  -----
n01   Profiles  --       p1,p2
`,
			inDb: `nodeprofiles:
  p1: {}
  p2: {}
nodes:
  n01:
    profiles:
    - p1
    - p2
`,
		},
		{
			name:    "multiple overlays list long with negation",
			args:    []string{"-l"},
			wantErr: false,
			stdout: `
NODE NAME  KERNEL VERSION  CONTAINER  OVERLAYS (S/R)
---------  --------------  ---------  --------------
n01        --              --         /rop1,rop2
`,
			inDb: `nodeprofiles:
  p1:
    runtime overlay:
    - rop1
    - rop2
nodes:
  n01:
    profiles:
    - p1
`,
		},
		{
			name:    "multiple overlays list long",
			args:    []string{"-l"},
			wantErr: false,
			stdout: `
NODE NAME  KERNEL VERSION  CONTAINER  OVERLAYS (S/R)
---------  --------------  ---------  --------------
n01        --              --         sop1/rop1,rop2,nop1,~rop1
`,
			inDb: `nodeprofiles:
  p1:
    system overlay:
    - sop1
    runtime overlay:
    - rop1
    - rop2
nodes:
  n01:
    profiles:
    - p1
    runtime overlay:
    - nop1
    - ~rop1
`,
		},
		{
			name:    "multiple overlays list all with negation",
			args:    []string{"-a"},
			wantErr: false,
			stdout: `
NODE  FIELD           PROFILE  VALUE
----  -----           -------  -----
n01   Profiles        --       p1
n01   RuntimeOverlay  p1,n01   rop1,rop2,nop1,~rop1
n01   SystemOverlay   p1       sop1
`,
			inDb: `nodeprofiles:
  p1:
    system overlay:
    - sop1
    runtime overlay:
    - rop1
    - rop2
nodes:
  n01:
    profiles:
    - p1
    runtime overlay:
    - nop1
    - ~rop1
`,
		},
		{
			name:    "multiple overlays list all",
			args:    []string{"-a"},
			wantErr: false,
			stdout: `
NODE  FIELD           PROFILE  VALUE
----  -----           -------  -----
n01   Profiles        --       p1
n01   RuntimeOverlay  p1,n01   rop1,rop2,nop1
`,
			inDb: `nodeprofiles:
  p1:
    runtime overlay:
    - rop1
    - rop2
nodes:
  n01:
    profiles:
    - p1
    runtime overlay:
    - nop1
`,
		},
		{
			name:    "network onboot",
			args:    []string{"-a"},
			wantErr: false,
			stdout: `
NODE  FIELD                    PROFILE  VALUE
----  -----                    -------  -----
n1    NetDevs[default].OnBoot  --       true
`,
			inDb: `nodes:
  n1:
    network devices:
      default:
        onboot: true
`,
		},
		{
			name:    "empty network device",
			args:    []string{"-a"},
			wantErr: false,
			stdout: `
NODE  FIELD  PROFILE  VALUE
----  -----  -------  -----
`,
			inDb: `nodes:
  wwnode1:
    network devices:
      default: {}
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
			buf := new(bytes.Buffer)
			baseCmd.SetOut(buf)
			baseCmd.SetErr(buf)
			wwlog.SetLogWriter(buf)
			wwlog.SetLogWriterErr(buf)
			wwlog.SetLogWriterInfo(buf)
			err := baseCmd.Execute()
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tt.stdout), strings.TrimSpace(buf.String()))
		})
	}
}

func TestListMultipleFormats(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdout  string
		inDb    string
		wantErr bool
	}{
		{
			name: "single node list yaml output",
			args: []string{"-y"},
			stdout: `
- profiles:
    - default
  kernel: {}
  ipmi: {}
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
			name: "single node list json output",
			args: []string{"-j"},
			stdout: `
[
  {
    "Discoverable": "",
    "AssetKey": "",
    "Profiles": [
      "default"
    ],
    "Comment": "",
    "ClusterName": "",
    "ContainerName": "",
    "Ipxe": "",
    "RuntimeOverlay": null,
    "SystemOverlay": null,
    "Kernel": {},
    "Ipmi": {
      "UserName": "",
      "Password": "",
      "Ipaddr": "",
      "Gateway": "",
      "Netmask": "",
      "Port": "",
      "Interface": "",
      "EscapeChar": "",
      "Write": "",
      "Template": "",
      "Tags": {}
    },
    "Init": "",
    "Root": "",
    "NetDevs": {},
    "Tags": {},
    "PrimaryNetDev": "",
    "Disks": null,
    "FileSystems": null,
    "Resources": null
  }
]
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
			name: "multiple nodes list json output",
			args: []string{"-j"},
			stdout: `
[
  {
    "Discoverable": "",
    "AssetKey": "",
    "Profiles": [
      "default"
    ],
    "Comment": "",
    "ClusterName": "",
    "ContainerName": "",
    "Ipxe": "",
    "RuntimeOverlay": null,
    "SystemOverlay": null,
    "Kernel": {},
    "Ipmi": {
      "UserName": "",
      "Password": "",
      "Ipaddr": "",
      "Gateway": "",
      "Netmask": "",
      "Port": "",
      "Interface": "",
      "EscapeChar": "",
      "Write": "",
      "Template": "",
      "Tags": {}
    },
    "Init": "",
    "Root": "",
    "NetDevs": {},
    "Tags": {},
    "PrimaryNetDev": "",
    "Disks": null,
    "FileSystems": null,
    "Resources": null
  },
  {
    "Discoverable": "",
    "AssetKey": "",
    "Profiles": [
      "default"
    ],
    "Comment": "",
    "ClusterName": "",
    "ContainerName": "",
    "Ipxe": "",
    "RuntimeOverlay": null,
    "SystemOverlay": null,
    "Kernel": {},
    "Ipmi": {
      "UserName": "",
      "Password": "",
      "Ipaddr": "",
      "Gateway": "",
      "Netmask": "",
      "Port": "",
      "Interface": "",
      "EscapeChar": "",
      "Write": "",
      "Template": "",
      "Tags": {}
    },
    "Init": "",
    "Root": "",
    "NetDevs": {},
    "Tags": {},
    "PrimaryNetDev": "",
    "Disks": null,
    "FileSystems": null,
    "Resources": null
  }
]
`,
			inDb: `
nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
  n02:
    profiles:
    - default
`,
		},
		{
			name: "multiple nodes list yaml output",
			args: []string{"-y"},
			stdout: `
- profiles:
    - default
  kernel: {}
  ipmi: {}
- profiles:
    - default
  kernel: {}
  ipmi: {}
`,
			inDb: `
nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
  n02:
    profiles:
    - default
`,
		},
		{
			name:    "single node list with network",
			args:    []string{"-a"},
			wantErr: false,
			stdout: `
NODE  FIELD                    PROFILE  VALUE
----  -----                    -------  -----
n01   Profiles                 --       default
n01   NetDevs[default].Hwaddr  --       aa:bb:cc:dd:ee:ff
n01   NetDevs[default].Ipaddr  --       1.1.1.1
`,
			inDb: `nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default
    network devices:
      default:
        hwaddr: aa:bb:cc:dd:ee:ff
        ipaddr: 1.1.1.1
`,
		},
	}

	warewulfd.SetNoDaemon()
	for _, tt := range tests {
		env := testenv.New(t)
		env.WriteFile("etc/warewulf/nodes.conf", tt.inDb)
		var err error
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			wwlog.SetLogWriter(buf)
			baseCmd := GetCommand()
			baseCmd.SetArgs(tt.args)
			baseCmd.SetOut(buf)
			baseCmd.SetErr(buf)
			err = baseCmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, strings.TrimSpace(tt.stdout), strings.TrimSpace(buf.String()))
		})
	}
}
