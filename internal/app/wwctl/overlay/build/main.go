package build

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/hostlist"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/overlay"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	controller := warewulfconf.Get()
	nodeDB, err := node.New()
	if err != nil {
		return errors.Wrap(err, "couldn't open node configuration")
	}

	db, err := nodeDB.FindAllNodes()
	if err != nil {
		return errors.Wrap(err, "could not get node list")
	}

	if len(args) > 0 {
		args = hostlist.Expand(args)
		db = node.FilterByName(db, args)

		if len(db) < len(args) {
			return errors.New("failed to find nodes")
		}
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
			if len(db) != 1 {
				return errors.New("nust specify one node to build overlay")
			}

			for _, node := range db {
				return overlay.BuildOverlayIndir(node, OverlayNames, OverlayDir)
			}
		} else {
			// TODO this seems different than what is set in BuildHostOverlay
			hostname, _ := os.Hostname()
			node := node.NewNode(hostname)
			wwlog.Info("building overlay for host: %s", hostname)
			return overlay.BuildOverlayIndir(node, OverlayNames, OverlayDir)

		}

	}

	if BuildHost && controller.Warewulf.EnableHostOverlay {
		err := overlay.BuildHostOverlay()
		if err != nil {
			return fmt.Errorf("host overlay could not be built: %s", err)
		}
	}

	if BuildNodes || (!BuildHost && !BuildNodes) {
		if len(OverlayNames) > 0 {
			err = overlay.BuildSpecificOverlays(db, OverlayNames)
		} else {
			err = overlay.BuildAllOverlays(db)
		}

		if err != nil {
			return fmt.Errorf("Some overlays failed to be generated: %s", err)
		}
	}
	return nil
}
