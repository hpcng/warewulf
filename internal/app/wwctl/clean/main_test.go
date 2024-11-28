package clean

import (
	"path"
	"testing"

	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_Clean(t *testing.T) {
	wwlog.SetLogLevel(wwlog.DEBUG)
	env := testenv.New(t)
	env.WriteFile(t, "etc/warewulf/nodes.conf",
		`WW_INTERNAL: 45
nodeprofiles: {}
nodes:
  node1: {}
`)
	wwconf := warewulfconf.Get()
	env.WriteFileAbs(t, path.Join(wwconf.Paths.WWProvisiondir, "overlays/node1/__SYSTEM__.img"), "Fake System")
	env.WriteFileAbs(t, path.Join(wwconf.Paths.WWProvisiondir, "overlays/node2/__SYSTEM__.img"), "Fake System")
	env.WriteFileAbs(t, path.Join(wwconf.Paths.Cachedir, "warewulf/test"), "Nothing to see here")
	baseCmd := GetCommand()
	err := baseCmd.Execute()
	assert.NoError(t, err)
	assert.FileExists(t, path.Join(wwconf.Paths.WWProvisiondir, "overlays/node1/__SYSTEM__.img"))
	assert.NoFileExists(t, path.Join(wwconf.Paths.WWProvisiondir, "overlays/node2/__SYSTEM__.img"))
	assert.NoDirExists(t, path.Join(wwconf.Paths.Cachedir, "warewulf"))
}
