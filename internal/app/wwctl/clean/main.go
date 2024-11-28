package clean

import (
	"github.com/warewulf/warewulf/internal/pkg/api/clean"

	"github.com/spf13/cobra"
)

func CobraRunE(vars *variables) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		return clean.Clean()
	}
}
