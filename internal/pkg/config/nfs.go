package config

import (
	"github.com/creasty/defaults"
)

// NFSConf represents the NFS configuration that will be used by
// Warewulf to generate exports on the server and mounts on compute
// nodes.
type NFSConf struct {
	EnabledP        *bool            `yaml:"enabled,omitempty" default:"true"`
	ExportsExtended []*NFSExportConf `yaml:"export paths,omitempty" default:"[]"`
	SystemdName     string           `yaml:"systemd name,omitempty" default:"nfsd"`
}

func (this NFSConf) Enabled() bool {
	return BoolP(this.EnabledP)
}

// An NFSExportConf reprents a single NFS export / mount.
type NFSExportConf struct {
	Path          string `yaml:"path" default:"/dev/null"`
	ExportOptions string `yaml:"export options,omitempty" default:"rw,sync,no_subtree_check"`
}

// Implements the Unmarshal interface for NFSConf to set default
// values.
func (conf *NFSConf) Unmarshal(unmarshal func(interface{}) error) error {
	if err := defaults.Set(conf); err != nil {
		return err
	}
	return nil
}
