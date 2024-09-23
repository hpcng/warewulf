package console

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/pkg/hostlist"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/power"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	var returnErr error = nil

	nodeDB, err := node.New()
	if err != nil {
		wwlog.Error("Could not open node configuration: %s", err)
		os.Exit(1)
	}

	nodes, err := nodeDB.FindAllNodes()
	if err != nil {
		wwlog.Error("Could not get node list: %s", err)
		os.Exit(1)
	}

	args = hostlist.Expand(args)

	if len(args) > 0 {
		nodes = node.FilterByName(nodes, args)
	} else {
		//nolint:errcheck
		cmd.Usage()
		os.Exit(1)
	}

	if len(nodes) == 0 {
		fmt.Printf("No nodes found\n")
		os.Exit(1)
	}

	for _, node := range nodes {

		if node.Ipmi.Ipaddr.IsUnspecified() {
			wwlog.Error("%s: No IPMI IP address", node.Id())
			continue
		}

		ipmiCmd := power.IPMI{
			NodeName:   node.Id(),
			HostName:   node.Ipmi.Ipaddr.String(),
			Port:       node.Ipmi.Port,
			User:       node.Ipmi.UserName,
			Password:   node.Ipmi.Password,
			AuthType:   "MD5",
			Interface:  node.Ipmi.Interface,
			EscapeChar: node.Ipmi.EscapeChar,
		}

		err := ipmiCmd.Console()

		if err != nil {
			wwlog.Error("%s: Console problem", node.Id())
			returnErr = err
			continue
		}

	}

	return returnErr
}
