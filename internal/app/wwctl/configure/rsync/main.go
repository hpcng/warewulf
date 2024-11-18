package rsync

import (
	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/pkg/configure"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	return configure.RSYNC()
}
