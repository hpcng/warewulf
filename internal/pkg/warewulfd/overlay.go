package warewulfd

import (
	"net/http"
	"os"
	"path"
	"strconv"

	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/overlay"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func OverlaySend(w http.ResponseWriter, req *http.Request) {
	conf := warewulfconf.Get()
	overlaySourceDir := overlay.OverlaySourceDir("wwroot")
	if !util.IsDir(overlaySourceDir) {
		wwlog.Error("Overlay source dir wwroot doesn't exist")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	rinfo, err := parseReqRender(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		wwlog.ErrorExc(err, "got error: %w")
		return
	}

	wwlog.Info("recv: render req node: %s, overlay: %s", rinfo.node, rinfo.overlay)
	if conf.Warewulf.Secure() {
		if rinfo.remoteport >= 1024 {
			wwlog.Denied("Non-privileged port: %s", req.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	overlayFile := path.Join(overlaySourceDir, rinfo.overlay)
	if !path.IsAbs(overlayFile) {
		wwlog.Denied("Path %s isn't absolute", overlayFile)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	nodeName := rinfo.node
	if nodeName != "" {
		nodeDB, err := node.New()
		if err != nil {
			wwlog.ErrorExc(err, "error when opening node database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		node, err := nodeDB.GetNode(nodeName)
		if err != nil {
			wwlog.Warn("Unknown node %s from: %s", nodeName, req.RemoteAddr)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if !util.IsFile(overlayFile) {
			overlayFile += ".ww"
			wwlog.Debug("trying to use .ww for file: %s", overlayFile)
		}
		if !util.IsFile(overlayFile) {
			wwlog.Denied("file doesn't exists: %s", overlayFile)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		tstruct, err := overlay.InitStruct(node)
		if err != nil {
			wwlog.ErrorExc(err, "error when initializing template data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tstruct.BuildSource = overlayFile
		buffer, _, _, err := overlay.RenderTemplateFile(overlayFile, tstruct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			wwlog.ErrorExc(err, "")
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))
		_, err = buffer.WriteTo(w)
		if err != nil {
			wwlog.ErrorExc(err, "")
		}
		wwlog.Info("%15s: %s", node.Id(), overlayFile)
	} else {
		if !util.IsFile(overlayFile) {
			wwlog.Denied("file doesn't exists: %s", overlayFile)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fileBytes, err := os.ReadFile(overlayFile)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			wwlog.ErrorExc(err, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = w.Write(fileBytes)
		if err != nil {
			wwlog.ErrorExc(err, "")
		}
		wwlog.Info("send overlay for node %s: %s", nodeName, overlayFile)

	}
}
