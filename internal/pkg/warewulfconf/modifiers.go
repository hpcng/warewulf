package warewulfconf

import (
	"os"

	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"gopkg.in/yaml.v2"
)

func (controller *ControllerConf) Persist() error {

	out, err := yaml.Marshal(controller)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(ConfigFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "%s\n", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(string(out) + "\n")
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Unable to write to warewulf.conf\n")
		return err
	}

	return nil
}
