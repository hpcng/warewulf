package warewulfd

import (
	"net/http"
	"os/exec"
	"strings"
)

func InfoSend(w http.ResponseWriter, req *http.Request) {

	var cmd *exec.Cmd
	url := strings.Split(req.URL.Path, "/")

	if url[2] == "" {
		daemonLogf("ERROR: Info request from %s missing argument\n", req.RemoteAddr)
		w.WriteHeader(400)
		return
	}

	switch url[2] {
	case "nodes":
		cmd = exec.Command("/usr/bin/wwctl", "node", "list")
	case "ready":
		cmd = exec.Command("/usr/bin/wwctl", "node", "ready")
	default:
		daemonLogf("ERROR: Unrecognized info request from %s\n", req.RemoteAddr)
		w.WriteHeader(400)
		return
	}

	stdout, err := cmd.CombinedOutput()

	if err != nil {
		daemonLogf("ERROR: wwctl exec error: " + err.Error() + "\n")
		stdout = []byte(err.Error() + "\n")
	}

	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write(stdout)
	if err != nil {
		daemonLogf("ERROR: %s\n", err)
	}
}
