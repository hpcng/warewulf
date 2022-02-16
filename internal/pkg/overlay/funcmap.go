package overlay

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/hpcng/warewulf/internal/pkg/buildconfig"
	"github.com/hpcng/warewulf/internal/pkg/container"
	"github.com/hpcng/warewulf/internal/pkg/util"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
)

/*
Reads a file file from the host fs. If the file has nor '/' prefix
the path is relative to SYSCONFDIR.
Templates in the file are no evaluated.
*/
func templateFileInclude(inc string) string {
	if !strings.HasPrefix(inc, "/") {
		inc = path.Join(buildconfig.SYSCONFDIR(), "warewulf", inc)
	}
	wwlog.Printf(wwlog.DEBUG, "Including file into template: %s\n", inc)
	content, err := ioutil.ReadFile(inc)
	if err != nil {
		wwlog.Printf(wwlog.VERBOSE, "Could not include file into template: %s\n", err)
	}
	return strings.TrimSuffix(string(content), "\n")
}

/*
Reads a file into template the abort string is found in a line. First argument
is the file to read, the second the abort string
Templates in the file are no evaluated.
*/
func templateFileBlock(inc string, abortStr string) (string, error) {
	if !strings.HasPrefix(inc, "/") {
		inc = path.Join(buildconfig.SYSCONFDIR(), "warewulf", inc)
	}
	wwlog.Printf(wwlog.DEBUG, "Including file block into template: %s\n", inc)
	readFile, err := os.Open(inc)
	defer readFile.Close()
	if err != nil {
		return "", err
	}
	var cont string
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if strings.Contains(line, abortStr) {
			break
		}
		cont += line + "\n"
	}
	return cont, nil

}

/*
Reads a file relative to given container.
Templates in the file are no evaluated.
*/
func templateContainerFileInclude(containername string, filepath string) string {
	wwlog.Printf(wwlog.VERBOSE, "Including file from Container into template: %s:%s\n", containername, filepath)

	if containername == "" {
		wwlog.Printf(wwlog.WARN, "Container is not defined for node: %s\n", filepath)
		return ""
	}

	if !container.ValidSource(containername) {
		wwlog.Printf(wwlog.WARN, "Template requires file(s) from non-existant container: %s:%s\n", containername, filepath)
		return ""
	}

	containerDir := container.RootFsDir(containername)

	wwlog.Printf(wwlog.DEBUG, "Including file from container: %s:%s\n", containerDir, filepath)

	if !util.IsFile(path.Join(containerDir, filepath)) {
		wwlog.Printf(wwlog.WARN, "Requested file from container does not exist: %s:%s\n", containername, filepath)
		return ""
	}

	content, err := ioutil.ReadFile(path.Join(containerDir, filepath))

	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Template include failed: %s\n", err)
	}
	return strings.TrimSuffix(string(content), "\n")
}
