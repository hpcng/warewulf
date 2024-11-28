package clean

import (
	"os"
	"path"

	_ "golang.org/x/exp/slices"

	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

/*
Cleans up the OCI cache and remains of deleted nodes
*/
func Clean() (err error) {
	warewulfconf := warewulfconf.Get()
	wwlog.Verbose("removing oci cache dir: %s", path.Join(warewulfconf.Paths.Cachedir+"/warewulf"))
	err = os.RemoveAll(path.Join(warewulfconf.Paths.Cachedir + "/warewulf"))
	if err != nil {
		return err
	}
	nodeDB, err := node.New()
	if err != nil {
		return err
	}
	nodes := nodeDB.NodeList()
	dirList, err := os.ReadDir(path.Join(warewulfconf.Paths.WWProvisiondir, "overlays/"))
	if err != nil {
		return err
	}
	for _, item := range dirList {
		if !item.IsDir() {
			continue
		}
		if !util.InSlice(nodes, item.Name()) {
			wwlog.Verbose("removing overlays of delete node: %s", item.Name())
			err = os.RemoveAll(path.Join(warewulfconf.Paths.WWProvisiondir, "overlays/", item.Name()))
			if err != nil {
				return err
			}
		}
	}
	return
}
