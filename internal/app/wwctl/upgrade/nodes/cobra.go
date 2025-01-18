package nodes

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/upgrade"
	"github.com/warewulf/warewulf/internal/pkg/util"
)

var (
	addDefaults     bool
	replaceOverlays bool
	inputPath       string
	outputPath      string
	inputConfPath   string
)

func GetCommand() *cobra.Command {
	command := &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "nodes [OPTIONS]",
		Short:                 "Upgrade an existing nodes.conf",
		Long: `Upgrades nodes.conf from a previous version of Warewulf 4 to a format
supported by the current version.`,
		RunE: UpgradeNodesConf,
	}
	command.Flags().BoolVar(&addDefaults, "add-defaults", false, "Configure a default profile and set default node values")
	command.Flags().BoolVar(&replaceOverlays, "replace-overlays", false, "Replace 'wwinit' and 'generic' overlays with their split replacements")
	command.Flags().StringVarP(&inputPath, "input-path", "i", "", "Path to a legacy nodes.conf")
	command.Flags().StringVarP(&outputPath, "output-path", "o", "", "Path to write the upgraded nodes.conf to")
	command.Flags().StringVar(&inputConfPath, "with-warewulfconf", config.ConfigFile, "Path to a legacy warewulf.conf")
	if err := command.MarkFlagRequired("add-defaults"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("replace-overlays"); err != nil {
		panic(err)
	}
	return command
}

func UpgradeNodesConf(cmd *cobra.Command, args []string) error {
	inputPath := inputPath
	if inputPath == "" {
		inputPath = config.Get().Paths.NodesConf()
	}
	outputPath := outputPath
	if outputPath == "" {
		outputPath = config.Get().Paths.NodesConf()
	}

	confData, err := os.ReadFile(inputConfPath)
	if err != nil {
		return err
	}
	warewulfConf, err := upgrade.ParseConfig(confData)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	legacy, err := upgrade.ParseNodes(data)
	if err != nil {
		return err
	}
	upgraded := legacy.Upgrade(addDefaults, replaceOverlays, warewulfConf)
	if outputPath == "-" {
		upgradedYaml, err := upgraded.Dump()
		if err != nil {
			return err
		}
		fmt.Print(string(upgradedYaml))
		return nil
	} else {
		if err := util.CopyFile(outputPath, outputPath+"-old"); err != nil {
			return err
		}
		return upgraded.PersistToFile(outputPath)
	}
}
