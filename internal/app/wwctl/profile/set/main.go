package set

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	apiprofile "github.com/warewulf/warewulf/internal/pkg/api/profile"
	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/api/util"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
	"gopkg.in/yaml.v3"
)

func CobraRunE(vars *variables) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) error {

		// remove the default network as the all network values are assigned
		// to this network
		if !node.ObjectIsEmpty(vars.profileConf.NetDevs["UNDEF"]) || len(vars.profileAdd.NetTagsAdd) > 0 {
			netDev := *vars.profileConf.NetDevs["UNDEF"]
			vars.profileConf.NetDevs[vars.profileAdd.Net] = &netDev
			vars.profileConf.NetDevs[vars.profileAdd.Net].Tags = vars.profileAdd.NetTagsAdd
		}
		delete(vars.profileConf.NetDevs, "UNDEF")
		if vars.fsName != "" {
			if !strings.HasPrefix(vars.fsName, "/dev") {
				if vars.fsName == vars.partName {
					vars.fsName = "/dev/disk/by-partlabel/" + vars.partName
				} else {
					return fmt.Errorf("filesystems need to have a underlying blockdev")
				}
			}
			fs := *vars.profileConf.FileSystems["UNDEF"]
			vars.profileConf.FileSystems[vars.fsName] = &fs
		}
		delete(vars.profileConf.FileSystems, "UNDEF")
		if vars.diskName != "" && vars.partName != "" {
			prt := *vars.profileConf.Disks["UNDEF"].Partitions["UNDEF"]
			vars.profileConf.Disks["UNDEF"].Partitions[vars.partName] = &prt
			delete(vars.profileConf.Disks["UNDEF"].Partitions, "UNDEF")
			dsk := *vars.profileConf.Disks["UNDEF"]
			vars.profileConf.Disks[vars.diskName] = &dsk
		}
		if (vars.diskName != "") != (vars.partName != "") {
			return fmt.Errorf("partition and disk must be specified")
		}
		delete(vars.profileConf.Disks, "UNDEF")
		vars.profileConf.Ipmi.Tags = vars.profileAdd.IpmiTagsAdd
		buffer, err := yaml.Marshal(vars.profileConf)
		if err != nil {
			return fmt.Errorf("can not marshall nodeInfo: %s", err)
		}
		wwlog.Debug("sending following values: %s", string(buffer))
		set := wwapiv1.ConfSetParameter{
			NodeConfYaml:     string(buffer[:]),
			NetdevDelete:     vars.setNetDevDel,
			PartitionDelete:  vars.setPartDel,
			DiskDelete:       vars.setDiskDel,
			FilesystemDelete: vars.setFsDel,
			TagAdd:           vars.profileAdd.TagsAdd,
			TagDel:           vars.profileDel.TagsDel,
			NetTagAdd:        vars.profileAdd.NetTagsAdd,
			NetTagDel:        vars.profileDel.NetTagsDel,
			IpmiTagAdd:       vars.profileAdd.IpmiTagsAdd,
			IpmiTagDel:       vars.profileDel.IpmiTagsDel,

			AllConfs: vars.setNodeAll,
			Force:    vars.setForce,
			ConfList: args,
		}

		if !vars.setYes {
			var profileCount uint
			// The checks run twice in the prompt case.
			// Avoiding putting in a blocking prompt in an API.
			_, profileCount, err = apiprofile.ProfileSetParameterCheck(&set)
			if err != nil {
				return err
			}
			yes := util.ConfirmationPrompt(fmt.Sprintf("Are you sure you want to modify %d profile(s)", profileCount))
			if !yes {
				return err
			}
		}
		return apiprofile.ProfileSet(&set)
	}
}
