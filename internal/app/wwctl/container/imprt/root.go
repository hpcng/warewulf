package imprt

import (
	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

type variables struct {
	SetForce      bool
	SetUpdate     bool
	SetBuild      bool
	SetDefault    bool
	SyncUser      bool
	OciNoHttps    bool
	OciUsername   string
	OciPassword   string
	recordChanges bool
}

func GetCommand() *cobra.Command {
	vars := variables{}
	baseCmd := &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "import [OPTIONS] SOURCE [NAME]",
		Short:                 "Import a container into Warewulf",
		Aliases:               []string{"pull"},
		Long: `This command will pull and import a container into Warewulf from SOURCE,
optionally renaming it to NAME. The SOURCE must be in a supported URI format. Formats
are:
 * docker://registry.example.org/example:latest
 * docker-daemon://example:latest
 * file://path/to/archive/tar/ball
 * /path/to/archive/tar/ball
 * /path/to/chroot/
Imported containers are used to create bootable VNFS images.`,
		Example: "wwctl container import docker://ghcr.io/warewulf/warewulf-rockylinux:8 rockylinux-8",
		RunE:    CobraRunE(&vars),
		Args:    cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if vars.SetForce && vars.SetUpdate {
				wwlog.Warn("Both --force and --update flags are set, will ignore --update flag")
			}
		},
	}

	baseCmd.PersistentFlags().BoolVarP(&vars.SetForce, "force", "f", false, "Force overwrite of an existing container")
	baseCmd.PersistentFlags().BoolVarP(&vars.SetUpdate, "update", "u", false, "Update and overwrite an existing container")
	baseCmd.PersistentFlags().BoolVarP(&vars.SetBuild, "build", "b", false, "Build container when after pulling")
	baseCmd.PersistentFlags().BoolVar(&vars.SetDefault, "setdefault", false, "Set this container for the default profile")
	baseCmd.PersistentFlags().BoolVar(&vars.SyncUser, "syncuser", false, "Synchronize UIDs/GIDs from host to container")
	baseCmd.PersistentFlags().BoolVar(&vars.OciNoHttps, "ocinohttps", false, "Ignore wrong TLS certificates, superseedes env WAREWULF_OCI_NOHTTPS")
	baseCmd.PersistentFlags().StringVar(&vars.OciUsername, "ociusername", "", "Set username for the access to the registry, superseedes env WAREWULF_OCI_USERNAME")
	baseCmd.PersistentFlags().StringVar(&vars.OciPassword, "ocipasswd", "", "Set password for the access to the registry, superseedes env WAREWULF_OCI_PASSWORD")
	baseCmd.PersistentFlags().BoolVarP(&vars.recordChanges, "record", "r", false, "Record the changes in separate layers")

	return baseCmd
}
