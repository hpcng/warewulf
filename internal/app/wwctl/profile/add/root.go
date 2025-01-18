package add

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/app/wwctl/completions"
	"github.com/warewulf/warewulf/internal/pkg/container"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/overlay"
)

type variables struct {
	profileConf node.Profile
	profileAdd  node.NodeConfAdd
}

// GetRootCommand returns the root cobra.Command for the application.
func GetCommand() *cobra.Command {
	vars := variables{}
	vars.profileConf = node.NewProfile("")
	baseCmd := &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "add PROFILE",
		Short:                 "Add a new node profile",
		Long:                  "This command adds a new named PROFILE.",
		Aliases:               []string{"new", "create"},
		RunE:                  CobraRunE(&vars),
		Args:                  cobra.ExactArgs(1),
	}
	vars.profileConf.CreateFlags(baseCmd)
	vars.profileAdd.CreateAddFlags(baseCmd)
	// register the command line completions
	if err := baseCmd.RegisterFlagCompletionFunc("container", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, _ := container.ListSources()
		return list, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Println(err)
	}
	if err := baseCmd.RegisterFlagCompletionFunc("kernelversion", completions.ProfileKernelVersion); err != nil {
		log.Println(err)
	}
	if err := baseCmd.RegisterFlagCompletionFunc("runtime", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, _ := overlay.FindOverlays()
		return list, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Println(err)
	}
	if err := baseCmd.RegisterFlagCompletionFunc("wwinit", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, _ := overlay.FindOverlays()
		return list, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Println(err)
	}
	return baseCmd
}
