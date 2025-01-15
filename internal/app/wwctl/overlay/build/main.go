package build

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/pkg/hostlist"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/overlay"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	nodeDB, err := node.New()
	if err != nil {
		return fmt.Errorf("could not open node configuration: %s", err)
	}

	allNodes, err := nodeDB.FindAllNodes()
	if err != nil {
		return fmt.Errorf("could not get node list: %s", err)
	}

	var filteredNodes []node.Node
	if len(args) > 0 {
		args = hostlist.Expand(args)
		filteredNodes = node.FilterNodeListByName(allNodes, args)

		if len(filteredNodes) < len(args) {
			return errors.New("failed to find nodes")
		}
	} else {
		filteredNodes = allNodes
	}

	// NOTE: this is to keep backward compatible
	// passing -O a,b,c versus -O a -O b -O c, but will also accept -O a,b -O c
	overlayNames := []string{}
	for _, name := range OverlayNames {
		names := strings.Split(name, ",")
		overlayNames = append(overlayNames, names...)
	}
	OverlayNames = overlayNames

	if OverlayDir != "" {
		if len(OverlayNames) == 0 {
			// TODO: should this behave the same as OverlayDir == "", and build default
			// set to overlays?
			return errors.New("must specify overlay(s) to build")
		}

		if len(args) > 0 {
			if len(filteredNodes) != 1 {
				return errors.New("must specify one node to build overlay")
			}

			for _, node := range filteredNodes {
				return overlay.BuildOverlayIndir(node, allNodes, OverlayNames, OverlayDir)
			}
		} else {
			return errors.New("must specify a node to build overlay")
		}
	}

	oldMask := syscall.Umask(000)
	defer syscall.Umask(oldMask)

	if len(OverlayNames) > 0 {
		err = overlay.BuildSpecificOverlays(filteredNodes, allNodes, OverlayNames, Workers)
	} else {
		err = overlay.BuildAllOverlays(filteredNodes, allNodes, Workers)
	}

	if err != nil {
		return fmt.Errorf("some overlays failed to be generated: %s", err)
	}
	return nil
}
