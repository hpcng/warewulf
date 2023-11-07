package list

import (
	"bytes"
	"os"
	"strings"
	"testing"

	warewulfconf "github.com/hpcng/warewulf/internal/pkg/config"
	"github.com/hpcng/warewulf/internal/pkg/node"
	"github.com/hpcng/warewulf/internal/pkg/warewulfd"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/stretchr/testify/assert"
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
			stdout: `  NODE NAME  PROFILES  NETWORK  
  n01        default            
`,
			inDb: `WW_INTERNAL: 43
nodeprofiles:
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
			stdout: `  NODE NAME  PROFILES  NETWORK
  n01        default            
  n02        default
`,
			inDb: `WW_INTERNAL: 43
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
			name:    "node list returns multiple nodes",
			args:    []string{"n01,n02"},
			wantErr: false,
			stdout: `  NODE NAME  PROFILES  NETWORK
  n01        default            
  n02        default
`,
			inDb: `WW_INTERNAL: 43
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
			name:    "node list returns multiple nodes (case 2)",
			args:    []string{"n01,n03"},
			wantErr: false,
			stdout: `  NODE NAME  PROFILES  NETWORK
  n01        default            
  n03        default
`,
			inDb: `WW_INTERNAL: 43
nodeprofiles:
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
			stdout: `  NODE NAME  PROFILES  NETWORK
  n01        default            
`,
			inDb: `WW_INTERNAL: 43
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
	}
	conf_yml := `WW_INTERNAL: 0`
	tempWarewulfConf, warewulfConfErr := os.CreateTemp("", "warewulf.conf-")
	assert.NoError(t, warewulfConfErr)
	defer os.Remove(tempWarewulfConf.Name())
	_, warewulfConfErr = tempWarewulfConf.Write([]byte(conf_yml))
	assert.NoError(t, warewulfConfErr)
	assert.NoError(t, tempWarewulfConf.Sync())
	assert.NoError(t, warewulfconf.New().Read(tempWarewulfConf.Name()))

	tempNodeConf, nodesConfErr := os.CreateTemp("", "nodes.conf-")
	assert.NoError(t, nodesConfErr)
	defer os.Remove(tempNodeConf.Name())
	node.ConfigFile = tempNodeConf.Name()
	warewulfd.SetNoDaemon()
	for _, tt := range tests {
		var err error
		_, err = tempNodeConf.Seek(0, 0)
		assert.NoError(t, err)
		assert.NoError(t, tempNodeConf.Truncate(0))
		_, err = tempNodeConf.Write([]byte(tt.inDb))
		assert.NoError(t, err)
		assert.NoError(t, tempNodeConf.Sync())
		t.Logf("Running test: %s\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			baseCmd := GetCommand()
			baseCmd.SetArgs(tt.args)
			buf := new(bytes.Buffer)
			baseCmd.SetOut(buf)
			baseCmd.SetErr(buf)
			wwlog.SetLogWriter(buf)
			err = baseCmd.Execute()
			assert.NoError(t, err)
			assert.Equal(t, strings.ReplaceAll(buf.String(), " ", ""), strings.ReplaceAll(tt.stdout, " ", ""), "Got wrong output")
		})
	}
}
