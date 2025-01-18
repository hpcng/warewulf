package overlay

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/container"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

/*
Reads a file file from the host fs. If the file has nor '/' prefix
the path is relative to Paths.SysconfdirTemplates in the file are no evaluated.
*/
func templateFileInclude(inc string) string {
	conf := warewulfconf.Get()
	if !strings.HasPrefix(inc, "/") {
		inc = path.Join(conf.Paths.Sysconfdir, "warewulf", inc)
	}
	wwlog.Debug("Including file into template: %s", inc)
	content, err := os.ReadFile(inc)
	if err != nil {
		wwlog.Verbose("Could not include file into template: %s", err)
	}
	return strings.TrimSuffix(string(content), "\n")
}

/*
Reads a file into template the abort string is found in a line. First argument
is the file to read, the second the abort string
Templates in the file are no evaluated.
*/
func templateFileBlock(inc string, abortStr string) (string, error) {
	conf := warewulfconf.Get()
	if !strings.HasPrefix(inc, "/") {
		inc = path.Join(conf.Paths.Sysconfdir, "warewulf", inc)
	}
	wwlog.Debug("Including file block into template: %s", inc)
	readFile, err := os.Open(inc)
	if err != nil {
		wwlog.Info("couldn't read block %s: %s", inc, err)
		return abortStr, nil
	}
	defer readFile.Close()
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

	// NOTE: the text originally contains N-1 newlines for N lines, but the above
	// loop will always add one at the end
	// Avoids adding a blank line that was not present in the original file
	// by adding 'abort' string to the end of the included block (without a newline)
	// instead of manually in the template
	cont += abortStr

	return cont, nil

}

/*
Reads a file relative to given container.
Templates in the file are no evaluated.
*/
func templateContainerFileInclude(containername string, filepath string) string {
	wwlog.Verbose("Including file from Container into template: %s:%s", containername, filepath)

	if containername == "" {
		wwlog.Warn("Container is not defined for node: %s", filepath)
		return ""
	}

	if !container.ValidSource(containername) {
		wwlog.Warn("Template requires file(s) from non-existant container: %s:%s", containername, filepath)
		return ""
	}

	containerDir := container.RootFsDir(containername)

	wwlog.Debug("Including file from container: %s:%s", containerDir, filepath)

	if !util.IsFile(path.Join(containerDir, filepath)) {
		wwlog.Warn("Requested file from container does not exist: %s:%s", containername, filepath)
		return ""
	}

	content, err := os.ReadFile(path.Join(containerDir, filepath))

	if err != nil {
		wwlog.Error("Template include failed: %s", err)
	}
	return strings.TrimSuffix(string(content), "\n")
}

// don't return an error as we use this function for template evaluation,
// so error will turn up there as the return string
func createIgnitionJson(node *node.Node) string {
	conf, rep, err := node.GetConfig()
	if len(conf.Storage.Disks) == 0 && len(conf.Storage.Filesystems) == 0 {
		wwlog.Debug("no disks or filesystems present, don't create a json object")
		return ""
	}
	if err != nil {
		wwlog.Error("disk, filesystem configuration has following error: ", fmt.Sprint(err))
		return fmt.Sprint(err)
	}
	if rep != "" {
		wwlog.Warn("%s storage configuration has following non fatal problems: %s", node.Id, rep)
	}
	tmpYaml, _ := json.Marshal(&conf)
	return string(tmpYaml)
}

func importSoftlink(lnk string) string {
	target, err := filepath.EvalSymlinks(lnk)
	if err != nil {
		return "abort"
	}
	wwlog.Debug("importing softlink pointing to: %s", target)
	return softlink(target)
}

func softlink(target string) string {
	return fmt.Sprintf("{{ /* softlink \"%s\" */ }}", target)
}
