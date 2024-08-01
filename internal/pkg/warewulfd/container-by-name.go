package warewulfd

import (
	"net/http"
	"path"
	"strings"

	"github.com/warewulf/warewulf/internal/pkg/container"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func ContainerByNameSend(w http.ResponseWriter, req *http.Request) {
	wwlog.Debug("Requested URL: %s", req.URL.String())

	url := strings.Split(req.URL.Path, "?")[0]
	path_parts := strings.Split(url, "/")

	if len(path_parts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		wwlog.Error("invalid /container-by-name/$name URL")
		return
	}

	container_name := path_parts[2]
	wwlog.Debug("container-by-name: %s", container_name)

	stage_file := path.Join(container.ImageParentDir(), container_name)
	wwlog.Serv("stage_file '%s'", stage_file)

	if util.IsFile(stage_file) {
		err := sendFile(w, req, stage_file, "")
		if err != nil {
			wwlog.ErrorExc(err, "")
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		wwlog.Error("container-by-name: not found: %s", stage_file)
	}
}
