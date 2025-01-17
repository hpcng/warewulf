package config

// RSYNCConf represents the configuration for the rsync service that
// Warewulf will configure.
type RSYNCConf struct {
	EnabledP    *bool  `yaml:"enabled,omitempty" default:"true"`
	SystemdName string `yaml:"systemd name,omitempty" default:"rsyncd"`
	Conf        string `yaml:"conf,omitempty" default:"/etc/rsyncd.conf"`
}

func (conf *RSYNCConf) Enabled() bool {
	return *conf.EnabledP
}
