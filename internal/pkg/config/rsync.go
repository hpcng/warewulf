package config

// RSYNCConf represents the configuration for the rsync service that
// Warewulf will configure.
type RSYNCConf struct {
	SystemdName string `yaml:"systemd name,omitempty" default:"rsyncd"`
	Conf        string `yaml:"conf,omitempty" default:"/etc/rsyncd.conf"`
}
