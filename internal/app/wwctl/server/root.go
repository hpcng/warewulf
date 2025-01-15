package server

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/pkg/warewulfd"
)

var (
	baseCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "server [OPTIONS]",
		Short:                 "Start Warewulf server",
		RunE:                  CobraRunE,
	}
)

func GetCommand() *cobra.Command {
	return baseCmd
}

func CobraRunE(cmd *cobra.Command, args []string) error {
	oldMask := syscall.Umask(000)
	defer syscall.Umask(oldMask)

	if err := warewulfd.DaemonInitLogging(); err != nil {
		return fmt.Errorf("failed to configure logging: %w", err)
	}
	if err := warewulfd.RunServer(); err != nil {
		return fmt.Errorf("failed to start Warewulf server: %w", err)
	}
	return nil
}
