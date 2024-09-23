package add

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/pkg/container"
	"github.com/warewulf/warewulf/internal/pkg/kernel"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/overlay"
)

// Holds the variables which are needed in CobraRunE
type variables struct {
	netName  string
	fsName   string
	partName string
	diskName string
	nodeConf node.NodeConf
}

// Returns the newly created command
func GetCommand() *cobra.Command {
	vars := variables{}
	vars.nodeConf = node.NewNode("")
	baseCmd := &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "add [OPTIONS] NODENAME",
		Short:                 "Add new node to Warewulf",
		Long:                  "This command will add a new node named NODENAME to Warewulf.",
		Aliases:               []string{"new", "create"},
		RunE:                  CobraRunE(&vars),
		Args:                  cobra.MinimumNArgs(1),
	}
	vars.nodeConf.CreateFlags(baseCmd)
	baseCmd.PersistentFlags().StringVar(&vars.netName, "netname", "default", "Set network name for network options")
	baseCmd.PersistentFlags().StringVar(&vars.fsName, "fsname", "", "set the file system name which must match a partition name")
	baseCmd.PersistentFlags().StringVar(&vars.partName, "partname", "", "set the partition name so it can be used by a file system")
	baseCmd.PersistentFlags().StringVar(&vars.diskName, "diskname", "", "set disk device name for the partition")
	// register the command line completions
	if err := baseCmd.RegisterFlagCompletionFunc("container", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, _ := container.ListSources()
		return list, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Println(err)
	}
	if err := baseCmd.RegisterFlagCompletionFunc("kerneloverride", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, _ := kernel.ListKernels()
		return list, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
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
	if err := baseCmd.RegisterFlagCompletionFunc("profile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var list []string
		nodeDB, _ := node.New()
		profiles, _ := nodeDB.FindAllProfiles()
		for _, profile := range profiles {
			list = append(list, profile.Id())
		}
		return list, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Println(err)
	}

	// GetRootCommand returns the root cobra.Command for the application.
	return baseCmd
}
