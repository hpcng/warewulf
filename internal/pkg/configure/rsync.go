package configure

import (
	"fmt"
	"os"
	"path"

	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/overlay"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func RSYNC() (err error) {
	controller := warewulfconf.Get()
	rsyncTemplate := path.Join(overlay.OverlaySourceDir("host"), controller.RSYNC.Conf+".ww")
	if !(util.IsFile(rsyncTemplate)) {
		return fmt.Errorf("'the overlay template '%s' does not exists in 'host' overlay", controller.RSYNC.Conf+".ww")
	}
	hostname, _ := os.Hostname()
	tstruct, err := overlay.InitStruct(node.NewNode(hostname))
	if err != nil {
		return err
	}
	buffer, backupFile, writeFile, err := overlay.RenderTemplateFile(
		rsyncTemplate,
		tstruct)
	if err != nil {
		return
	}
	info, err := os.Stat(rsyncTemplate)
	if err != nil {
		return
	}
	if writeFile {
		err = overlay.CarefulWriteBuffer(controller.RSYNC.Conf, buffer, backupFile, info.Mode())
		if err != nil {
			return fmt.Errorf("could not write file from template: %w", err)
		}
	}
	wwlog.Info("Enabling and restarting the rsyncd services")
	if controller.RSYNC.SystemdName == "" {
		return fmt.Errorf("no name for rsyncd service defined in warewulf.conf")
	} else {
		err := util.SystemdStart(controller.RSYNC.SystemdName)
		if err != nil {
			return fmt.Errorf("failed to start: %w", err)
		}
	}
	return

}
