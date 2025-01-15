package edit

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/warewulfd"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
	"gopkg.in/yaml.v3"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	if !node.CanWriteConfig() {
		return fmt.Errorf("can not write to config: exiting")
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "/bin/vi"
	}
	wwlog.Debug("using editor: %s", editor)

	registry, regErr := node.New()
	if regErr != nil {
		return regErr
	}

	if len(args) == 0 {
		for profileID := range registry.NodeProfiles {
			args = append(args, profileID)
		}
	}
	wwlog.Debug("profile list: %v", args)

	tempFile, tempErr := os.CreateTemp(os.TempDir(), "ww4ProfileEdit*.yaml")
	if tempErr != nil {
		return fmt.Errorf("could not create temp file: %s", tempErr)
	}
	defer os.Remove(tempFile.Name())

	if !NoHeader {
		yamlTemplate := node.UnmarshalConf(node.Profile{}, nil)
		if _, err := tempFile.WriteString("#profilename:\n#  " + strings.Join(yamlTemplate, "\n#  ") + "\n"); err != nil {
			return err
		}
	}

	origProfiles := make(map[string]*node.Profile)
	for _, profileID := range args {
		if n, ok := registry.NodeProfiles[profileID]; ok {
			origProfiles[profileID] = n
		}
	}

	if origYaml, err := util.EncodeYaml(origProfiles); err != nil {
		return err
	} else if _, err := tempFile.Write(origYaml); err != nil {
		return err
	}

	sum1, sumErr := util.HashFile(tempFile)
	if sumErr != nil {
		return sumErr
	}
	wwlog.Debug("original hash: %s", sum1)

	for {
		if err := util.ExecInteractive(editor, tempFile.Name()); err != nil {
			return fmt.Errorf("editor process exited with non-zero code: %w", err)
		}

		sum2, sumErr := util.HashFile(tempFile)
		if sumErr != nil {
			return sumErr
		}
		wwlog.Debug("edited hash: %s", sum2)

		if sum1 != sum2 {
			wwlog.Debug("modified")

			editYaml, err := io.ReadAll(tempFile)
			if err != nil {
				return err
			}
			editProfiles := make(map[string]*node.Profile)
			if err := yaml.Unmarshal(editYaml, &editProfiles); err != nil {
				wwlog.Error("%v\n", err)
				if util.Confirm("Parse error: retry") {
					continue
				} else {
					break
				}
			}

			var added, deleted, updated int
			for profileID := range origProfiles {
				if editProfile, ok := editProfiles[profileID]; !ok || editProfile == nil {
					wwlog.Verbose("delete profile: %s", profileID)
					delete(registry.NodeProfiles, profileID)
					deleted += 1
				}
			}
			for profileID := range editProfiles {
				if _, ok := origProfiles[profileID]; !ok {
					wwlog.Verbose("add profile: %s", profileID)
					added += 1
					registry.NodeProfiles[profileID] = editProfiles[profileID]
				} else if equalYaml, err := util.EqualYaml(origProfiles[profileID], editProfiles[profileID]); err != nil {
					return err
				} else if !equalYaml {
					wwlog.Verbose("update profile: %s", profileID)
					updated += 1
					registry.NodeProfiles[profileID] = editProfiles[profileID]
				}
			}

			if util.Confirm(fmt.Sprintf("Are you sure you want to add %d, delete %d, and update %d profiles", added, deleted, updated)) {
				if err := registry.Persist(); err != nil {
					return err
				}

				if err := warewulfd.DaemonReload(); err != nil {
					return fmt.Errorf("failed to reload warewulf daemon: %w", err)
				}
			}
			break
		} else {
			wwlog.Verbose("No changes")
			break
		}
	}

	return nil
}
