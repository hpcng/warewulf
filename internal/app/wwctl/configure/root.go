package configure

import (
	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/app/wwctl/configure/dhcp"
	"github.com/warewulf/warewulf/internal/app/wwctl/configure/hostfile"
	"github.com/warewulf/warewulf/internal/app/wwctl/configure/nfs"
	"github.com/warewulf/warewulf/internal/app/wwctl/configure/rsync"
	"github.com/warewulf/warewulf/internal/app/wwctl/configure/ssh"
	"github.com/warewulf/warewulf/internal/app/wwctl/configure/tftp"
)

var (
	baseCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "configure [OPTIONS]",
		Short:                 "Manage system services",
		Long: "This application allows you to manage and initialize Warewulf dependent system\n" +
			"services based on the configuration in the warewulf.conf file.",
		RunE: CobraRunE,
	}
	allFunctions bool
)

func init() {
	baseCmd.AddCommand(dhcp.GetCommand())
	baseCmd.AddCommand(tftp.GetCommand())
	baseCmd.AddCommand(ssh.GetCommand())
	baseCmd.AddCommand(nfs.GetCommand())
	baseCmd.AddCommand(hostfile.GetCommand())
	baseCmd.AddCommand(rsync.GetCommand())

	baseCmd.Flags().BoolVarP(&allFunctions, "all", "a", false, "Configure all services")
}

// GetRootCommand returns the root cobra.Command for the application.
func GetCommand() *cobra.Command {
	return baseCmd
}
