package wwctl

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/app/wwctl/clean"
	"github.com/warewulf/warewulf/internal/app/wwctl/configure"
	"github.com/warewulf/warewulf/internal/app/wwctl/container"
	"github.com/warewulf/warewulf/internal/app/wwctl/genconf"
	"github.com/warewulf/warewulf/internal/app/wwctl/kernel"
	"github.com/warewulf/warewulf/internal/app/wwctl/node"
	"github.com/warewulf/warewulf/internal/app/wwctl/overlay"
	"github.com/warewulf/warewulf/internal/app/wwctl/power"
	"github.com/warewulf/warewulf/internal/app/wwctl/profile"
	"github.com/warewulf/warewulf/internal/app/wwctl/server"
	"github.com/warewulf/warewulf/internal/app/wwctl/ssh"
	"github.com/warewulf/warewulf/internal/app/wwctl/version"
	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/help"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

var (
	rootCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "wwctl COMMAND [OPTIONS]",
		Short:                 "Warewulf Control",
		Long:                  "Control interface to the Warewulf Cluster Provisioning System.",
		PersistentPreRunE:     rootPersistentPreRunE,
		SilenceUsage:          true,
		SilenceErrors:         true,
	}
	verboseArg      bool
	DebugFlag       bool
	LogLevel        int
	WarewulfConfArg string
	AllowEmptyConf  bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseArg, "verbose", "v", false, "Run with increased verbosity.")
	rootCmd.PersistentFlags().BoolVarP(&DebugFlag, "debug", "d", false, "Run with debugging messages enabled.")
	rootCmd.PersistentFlags().IntVar(&LogLevel, "loglevel", wwlog.INFO, "Set log level to given string")
	_ = rootCmd.PersistentFlags().MarkHidden("loglevel")
	rootCmd.PersistentFlags().StringVar(&WarewulfConfArg, "warewulfconf", "", "Set the warewulf configuration file")
	rootCmd.PersistentFlags().BoolVar(&AllowEmptyConf, "emptyconf", false, "Allow empty configuration")
	_ = rootCmd.PersistentFlags().MarkHidden("emptyconf")
	rootCmd.SetUsageTemplate(help.UsageTemplate)
	rootCmd.SetHelpTemplate(help.HelpTemplate)
	rootCmd.AddCommand(overlay.GetCommand())
	rootCmd.AddCommand(container.GetCommand())
	rootCmd.AddCommand(node.GetCommand())
	rootCmd.AddCommand(kernel.GetCommand())
	rootCmd.AddCommand(power.GetCommand())
	rootCmd.AddCommand(profile.GetCommand())
	rootCmd.AddCommand(configure.GetCommand())
	rootCmd.AddCommand(server.GetCommand())
	rootCmd.AddCommand(version.GetCommand())
	rootCmd.AddCommand(ssh.GetCommand())
	rootCmd.AddCommand(genconf.GetCommand())
	rootCmd.AddCommand(clean.GetCommand())
}

// GetRootCommand returns the root cobra.Command for the application.
func GetRootCommand() *cobra.Command {
	return rootCmd
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) (err error) {
	if DebugFlag {
		wwlog.SetLogLevel(wwlog.DEBUG)
	} else if verboseArg {
		wwlog.SetLogLevel(wwlog.VERBOSE)
	} else {
		wwlog.SetLogLevel(wwlog.INFO)
	}
	if LogLevel != wwlog.INFO {
		wwlog.SetLogLevel(LogLevel)
	}
	conf := warewulfconf.Get()
	if !AllowEmptyConf && !conf.InitializedFromFile() {
		if WarewulfConfArg != "" {
			err = conf.Read(WarewulfConfArg)
		} else if os.Getenv("WAREWULFCONF") != "" {
			err = conf.Read(os.Getenv("WAREWULFCONF"))
		} else {
			err = conf.Read(warewulfconf.ConfigFile)
		}
	}
	if err != nil {
		wwlog.Error("version: %s relase: %s WW_INTERNAL: %s", warewulfconf.Version, warewulfconf.Release, warewulfconf.Confversion)
		return
	}
	err = conf.SetDynamicDefaults()
	return
}
