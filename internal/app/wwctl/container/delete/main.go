package delete

import (
	"fmt"

	"github.com/warewulf/warewulf/internal/pkg/api/container"
	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/util"

	"github.com/spf13/cobra"
)

func CobraRunE(cmd *cobra.Command, args []string) (err error) {
	cdp := &wwapiv1.ContainerDeleteParameter{
		ContainerNames: args,
	}

	if !SetYes {
		yes := util.Confirm(fmt.Sprintf("Are you sure you want to delete container %s", args))
		if !yes {
			return
		}

	}
	return container.ContainerDelete(cdp)
}
