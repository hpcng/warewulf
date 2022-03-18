package create

import (
	"github.com/spf13/cobra"
)

var (
	baseCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "create [OPTIONS] OVERLAY_NAME",
		Short:                 "Initialize a new Overlay",
		Long:                  "This command creates a new empty overlay with the given OVERLAY_NAME.",
		RunE:                  CobraRunE,
		Args:                  cobra.ExactArgs(1),
	}
)

func init() {
}

// GetRootCommand returns the root cobra.Command for the application.
func GetCommand() *cobra.Command {
	return baseCmd
}
