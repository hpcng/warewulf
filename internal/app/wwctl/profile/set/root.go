package set

import (
	"log"

	"github.com/hpcng/warewulf/internal/pkg/container"
	"github.com/hpcng/warewulf/internal/pkg/kernel"
	"github.com/hpcng/warewulf/internal/pkg/node"
	"github.com/hpcng/warewulf/internal/pkg/overlay"
	"github.com/spf13/cobra"
)

var (
	baseCmd = &cobra.Command{
		Use:   "set [OPTIONS] [PROFILE ...]",
		Short: "Configure node profile properties",
		Long: "This command sets configuration properties for the node PROFILE(s).\n\n" +
			"Note: use the string 'UNSET' to remove a configuration",
		Args: cobra.MinimumNArgs(1),
		RunE: CobraRunE,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			nodeDB, _ := node.New()
			profiles, _ := nodeDB.FindAllProfiles()
			var p_names []string
			for _, profile := range profiles {
				p_names = append(p_names, profile.Id.Get())
			}
			return p_names, cobra.ShellCompDirectiveNoFileComp
		},
	}
	SetNetDevDel string
	SetNodeAll   bool
	SetYes       bool
	SetForce     bool
	NetName      string
	ProfileConf  node.NodeConf
	Converters   []func() error
)

func init() {
	ProfileConf = node.NewConf()
	Converters = ProfileConf.CreateFlags(baseCmd,
		[]string{"ipaddr", "ipaddr6", "ipmiaddr", "profile"})
	baseCmd.PersistentFlags().StringVar(&NetName, "netname", "default", "Set network name for network options")
	baseCmd.PersistentFlags().StringVarP(&SetNetDevDel, "netdel", "D", "", "Delete the node's network device")
	baseCmd.PersistentFlags().BoolVarP(&SetNodeAll, "all", "a", false, "Set all nodes")
	baseCmd.PersistentFlags().BoolVarP(&SetYes, "yes", "y", false, "Set 'yes' to all questions asked")
	baseCmd.PersistentFlags().BoolVarP(&SetForce, "force", "f", false, "Force configuration (even on error)")
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

}

// GetRootCommand returns the root cobra.Command for the application.
func GetCommand() *cobra.Command {
	return baseCmd
}
