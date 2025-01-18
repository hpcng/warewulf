package set

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	apinode "github.com/warewulf/warewulf/internal/pkg/api/node"
	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/api/util"
	"github.com/warewulf/warewulf/internal/pkg/hostlist"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
	"gopkg.in/yaml.v3"
)

func CobraRunE(vars *variables) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) error {
		// remove the default network as the all network values are assigned
		// to this network
		if !node.ObjectIsEmpty(vars.nodeConf.NetDevs["UNDEF"]) || len(vars.nodeAdd.NetTagsAdd) > 0 {
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
		buffer, err := yaml.Marshal(vars.nodeConf)
		if err != nil {
			return fmt.Errorf("can not marshall nodeInfo: %s", err)
		}
		wwlog.Debug("sending following values: %s", string(buffer))
		args = hostlist.Expand(args)
		set := wwapiv1.ConfSetParameter{
			NodeConfYaml: string(buffer),

			NetdevDelete:     vars.nodeDel.NetDel,
			PartitionDelete:  vars.nodeDel.PartDel,
			DiskDelete:       vars.nodeDel.DiskDel,
			FilesystemDelete: vars.nodeDel.FsDel,
			TagAdd:           vars.nodeAdd.TagsAdd,
			TagDel:           vars.nodeDel.TagsDel,
			NetTagAdd:        vars.nodeAdd.NetTagsAdd,
			NetTagDel:        vars.nodeDel.NetTagsDel,
			IpmiTagAdd:       vars.nodeAdd.IpmiTagsAdd,
			IpmiTagDel:       vars.nodeDel.IpmiTagsDel,
			Netdev:           vars.nodeAdd.Net,
			AllConfs:         vars.setNodeAll,
			Force:            vars.setForce,
			ConfList:         args,
		}

		if !vars.setYes {
			var nodeCount uint
			// The checks run twice in the prompt case.
			// Avoiding putting in a blocking prompt in an API.
			_, nodeCount, err = apinode.NodeSetParameterCheck(&set)
			if err != nil {
				return nil
			}
			yes := util.ConfirmationPrompt(fmt.Sprintf("Are you sure you want to modify %d nodes(s)", nodeCount))
			if !yes {
				return nil
			}
		}
		return apinode.NodeSet(&set)
	}
}
