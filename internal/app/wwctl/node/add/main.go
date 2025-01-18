package add

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
	apinode "github.com/warewulf/warewulf/internal/pkg/api/node"
	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

/*
RunE needs a function of type func(*cobraCommand,[]string) err, but
in order to avoid global variables which mess up testing a function of
the required type is returned
*/
func CobraRunE(vars *variables) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// remove the UNDEF network as all network values are assigned
		// to this network
		if !node.ObjectIsEmpty(vars.nodeConf.NetDevs["UNDEF"]) {
			netDev := *vars.nodeConf.NetDevs["UNDEF"]
			vars.nodeConf.NetDevs[vars.nodeAdd.Net] = &netDev
			vars.nodeConf.NetDevs[vars.nodeAdd.Net].Tags = vars.nodeAdd.NetTagsAdd

		}
		delete(vars.nodeConf.NetDevs, "UNDEF")
		if vars.nodeAdd.FsName != "" {
			if !strings.HasPrefix(vars.nodeAdd.FsName, "/dev") {
				if vars.nodeAdd.FsName == vars.nodeAdd.PartName {
					vars.nodeAdd.FsName = "/dev/disk/by-partlabel/" + vars.nodeAdd.PartName
				} else {
					return fmt.Errorf("filesystems need to have a underlying blockdev")
				}
			}
			fs := *vars.nodeConf.FileSystems["UNDEF"]
			vars.nodeConf.FileSystems[vars.nodeAdd.FsName] = &fs
		}
		delete(vars.nodeConf.FileSystems, "UNDEF")
		if vars.nodeAdd.DiskName != "" && vars.nodeAdd.PartName != "" {
			prt := *vars.nodeConf.Disks["UNDEF"].Partitions["UNDEF"]
			vars.nodeConf.Disks["UNDEF"].Partitions[vars.nodeAdd.PartName] = &prt
			delete(vars.nodeConf.Disks["UNDEF"].Partitions, "UNDEF")
			dsk := *vars.nodeConf.Disks["UNDEF"]
			vars.nodeConf.Disks[vars.nodeAdd.DiskName] = &dsk
		}
		if (vars.nodeAdd.DiskName != "") != (vars.nodeAdd.PartName != "") {
			return fmt.Errorf("partition and disk must be specified")
		}
		delete(vars.nodeConf.Disks, "UNDEF")
		vars.nodeConf.Ipmi.Tags = vars.nodeAdd.IpmiTagsAdd
		if len(vars.nodeConf.Profiles) == 0 {
			vars.nodeConf.Profiles = []string{"default"}
		}
		buffer, err := yaml.Marshal(vars.nodeConf)
		if err != nil {
			wwlog.Error("Can't marshall nodeInfo", err)
			return err
		}
		set := wwapiv1.NodeAddParameter{
			NodeConfYaml: string(buffer[:]),
			NodeNames:    args,
			Force:        true,
		}
		return apinode.NodeAdd(&set)
	}
}
