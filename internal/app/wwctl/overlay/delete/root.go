package delete

import (
	"github.com/hpcng/warewulf/internal/pkg/overlay"
	"github.com/spf13/cobra"
)

var (
	baseCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "delete [OPTIONS] OVERLAY_NAME [FILE [FILE ...]]",
		Short:                 "Delete Warewulf Overlay or files",
		Long:                  "This command will delete FILEs within OVERLAY_NAME or the entire OVERLAY_NAME if no\nfiles are listed. Use with caution!",
		RunE:                  CobraRunE,
		Args:                  cobra.RangeArgs(1, 2),
		Aliases:               []string{"rm", "del"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				list, _ := overlay.FindOverlays()
				return list, cobra.ShellCompDirectiveNoFileComp
			} else if len(args) == 1 {
				ret, err := overlay.OverlayGetFiles(args[0])
				if err == nil {
					return ret, cobra.ShellCompDirectiveNoFileComp
				}
			}
			return []string{""}, cobra.ShellCompDirectiveNoFileComp
		},
	}
	Force   bool
	Parents bool
)

func init() {
	baseCmd.PersistentFlags().BoolVarP(&Force, "force", "f", false, "Force deletion of a non-empty overlay")
	baseCmd.PersistentFlags().BoolVarP(&Parents, "parents", "p", false, "Remove empty parent directories")
}

// GetRootCommand returns the root cobra.Command for the application.
func GetCommand() *cobra.Command {
	return baseCmd
}
