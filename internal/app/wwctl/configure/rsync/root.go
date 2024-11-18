package rsync

import (
	"github.com/spf13/cobra"
)

var (
	baseCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "rsync [OPTIONS]",
		Short:                 "Manage and initialize DHCP",
		Long: "rsync is needed for persistent installation for Warewulf. This command will configure rsync as defined\n" +
			"in the warewulf.conf file.",
		RunE: CobraRunE,
	}
)

func init() {
}

// GetRootCommand returns the root cobra.Command for the application.
func GetCommand() *cobra.Command {
	return baseCmd
}
