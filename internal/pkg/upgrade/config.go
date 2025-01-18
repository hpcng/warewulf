package upgrade

import (
	"gopkg.in/yaml.v3"

	"github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func ParseConfig(data []byte) (warewulfYaml *WarewulfYaml, err error) {
	warewulfYaml = new(WarewulfYaml)
	if err = yaml.Unmarshal(data, warewulfYaml); err != nil {
		return warewulfYaml, err
	}
	return warewulfYaml, nil
}

type WarewulfYaml struct {
	WWInternal      string        `yaml:"WW_INTERNAL"`
	Comment         string        `yaml:"comment"`
	Ipaddr          string        `yaml:"ipaddr"`
	Ipaddr6         string        `yaml:"ipaddr6"`
	Netmask         string        `yaml:"netmask"`
	Network         string        `yaml:"network"`
	Ipv6net         string        `yaml:"ipv6net"`
	Fqdn            string        `yaml:"fqdn"`
	Warewulf        *WarewulfConf `yaml:"warewulf"`
	DHCP            *DHCPConf     `yaml:"dhcp"`
	TFTP            *TFTPConf     `yaml:"tftp"`
	NFS             *NFSConf      `yaml:"nfs"`
	SSH             *SSHConf      `yaml:"ssh"`
	MountsContainer []*MountEntry `yaml:"container mounts"`
	Paths           *BuildConfig  `yaml:"paths"`
	WWClient        *WWClientConf `yaml:"wwclient"`
}

func (this *WarewulfYaml) Upgrade() (upgraded *config.WarewulfYaml) {
	upgraded = new(config.WarewulfYaml)
	if this.WWInternal != "" {
		logIgnore("WW_INTERNAL", this.WWInternal, "obsolete")
	}
	upgraded.Comment = this.Comment
	upgraded.Ipaddr = this.Ipaddr
	upgraded.Ipaddr6 = this.Ipaddr6
	upgraded.Netmask = this.Netmask
	upgraded.Network = this.Network
	upgraded.Ipv6net = this.Ipv6net
	upgraded.Fqdn = this.Fqdn
	if this.Warewulf != nil {
		upgraded.Warewulf = this.Warewulf.Upgrade()
	}
	if this.DHCP != nil {
		upgraded.DHCP = this.DHCP.Upgrade()
	}
	if this.TFTP != nil {
		upgraded.TFTP = this.TFTP.Upgrade()
	}
	if this.NFS != nil {
		upgraded.NFS = this.NFS.Upgrade()
	}
	if this.SSH != nil {
		upgraded.SSH = this.SSH.Upgrade()
	}
	upgraded.MountsContainer = make([]*config.MountEntry, 0)
	for _, mount := range this.MountsContainer {
		upgraded.MountsContainer = append(upgraded.MountsContainer, mount.Upgrade())
	}
	if this.Paths != nil {
		upgraded.Paths = this.Paths.Upgrade()
	}
	if this.WWClient != nil {
		upgraded.WWClient = this.WWClient.Upgrade()
	}
	if this.Warewulf != nil && this.Warewulf.DataStore != "" {
		if upgraded.Paths == nil {
			upgraded.Paths = new(config.BuildConfig)
		}
		if upgraded.Paths.Datadir == "" {
			upgraded.Paths.Datadir = this.Warewulf.DataStore
		}
	}
	return upgraded
}

type WarewulfConf struct {
	Port              int    `yaml:"port"`
	Secure            *bool  `yaml:"secure"`
	UpdateInterval    int    `yaml:"update interval"`
	AutobuildOverlays *bool  `yaml:"autobuild overlays"`
	EnableHostOverlay *bool  `yaml:"host overlay"`
	Syslog            *bool  `yaml:"syslog"`
	DataStore         string `yaml:"datastore"`
	GrubBoot          *bool  `yaml:"grubboot"`
}

func (this *WarewulfConf) Upgrade() (upgraded *config.WarewulfConf) {
	upgraded = new(config.WarewulfConf)
	upgraded.Port = this.Port
	upgraded.SecureP = this.Secure
	upgraded.UpdateInterval = this.UpdateInterval
	upgraded.AutobuildOverlaysP = this.AutobuildOverlays
	upgraded.EnableHostOverlayP = this.EnableHostOverlay
	upgraded.SyslogP = this.Syslog
	upgraded.GrubBootP = this.GrubBoot
	return upgraded
}

type DHCPConf struct {
	Enabled     *bool  `yaml:"enabled"`
	Template    string `yaml:"template"`
	RangeStart  string `yaml:"range start"`
	RangeEnd    string `yaml:"range end"`
	SystemdName string `yaml:"systemd name"`
}

func (this *DHCPConf) Upgrade() (upgraded *config.DHCPConf) {
	upgraded = new(config.DHCPConf)
	upgraded.EnabledP = this.Enabled
	upgraded.Template = this.Template
	upgraded.RangeStart = this.RangeStart
	upgraded.RangeEnd = this.RangeEnd
	upgraded.SystemdName = this.SystemdName
	return upgraded
}

type TFTPConf struct {
	Enabled      *bool             `yaml:"enabled"`
	TftpRoot     string            `yaml:"tftproot"`
	SystemdName  string            `yaml:"systemd name"`
	IpxeBinaries map[string]string `yaml:"ipxe"`
}

func (this *TFTPConf) Upgrade() (upgraded *config.TFTPConf) {
	upgraded = new(config.TFTPConf)
	upgraded.EnabledP = this.Enabled
	upgraded.TftpRoot = this.TftpRoot
	upgraded.SystemdName = this.SystemdName
	upgraded.IpxeBinaries = make(map[string]string)
	for name, binary := range this.IpxeBinaries {
		upgraded.IpxeBinaries[name] = binary
	}
	return upgraded
}

type NFSConf struct {
	Enabled         *bool            `yaml:"enabled"`
	Exports         []string         `yaml:"exports"`
	ExportsExtended []*NFSExportConf `yaml:"export paths"`
	SystemdName     string           `yaml:"systemd name"`
}

func (this *NFSConf) Upgrade() (upgraded *config.NFSConf) {
	upgraded = new(config.NFSConf)
	upgraded.EnabledP = this.Enabled
	upgraded.ExportsExtended = make([]*config.NFSExportConf, 0)
	for _, export := range this.Exports {
		extendedExport := new(config.NFSExportConf)
		extendedExport.Path = export
		upgraded.ExportsExtended = append(upgraded.ExportsExtended, extendedExport)
	}
	for _, export := range this.ExportsExtended {
		upgraded.ExportsExtended = append(upgraded.ExportsExtended, export.Upgrade())
	}
	upgraded.SystemdName = this.SystemdName
	return upgraded
}

type NFSExportConf struct {
	Path          string `yaml:"path"`
	ExportOptions string `yaml:"export options"`
	MountOptions  string `yaml:"mount options"`
	Mount         *bool  `yaml:"mount"`
}

func (this *NFSExportConf) Upgrade() (upgraded *config.NFSExportConf) {
	upgraded = new(config.NFSExportConf)
	upgraded.Path = this.Path
	upgraded.ExportOptions = this.ExportOptions
	if *(this.Mount) {
		wwlog.Warn("Legacy mount configured for NFS export %s: use `wwctl upgrade nodes --with-warewulfconf=<original file>` to port to nodes.conf", this.Path)
	}
	return upgraded
}

type SSHConf struct {
	KeyTypes []string `yaml:"key types"`
}

func (this *SSHConf) Upgrade() (upgraded *config.SSHConf) {
	upgraded = new(config.SSHConf)
	upgraded.KeyTypes = append([]string{}, this.KeyTypes...)
	return upgraded
}

type MountEntry struct {
	Source   string `yaml:"source"`
	Dest     string `yaml:"dest"`
	ReadOnly *bool  `yaml:"readonly"`
	Options  string `yaml:"options"`
	Copy     *bool  `yaml:"copy"`
}

func (this *MountEntry) Upgrade() (upgraded *config.MountEntry) {
	upgraded = new(config.MountEntry)
	upgraded.Source = this.Source
	upgraded.Dest = this.Dest
	upgraded.ReadOnlyP = this.ReadOnly
	upgraded.Options = this.Options
	upgraded.CopyP = this.Copy
	return upgraded
}

type BuildConfig struct {
	Bindir         string
	Sysconfdir     string
	Localstatedir  string
	Cachedir       string
	Ipxesource     string
	Srvdir         string
	Firewallddir   string
	Systemddir     string
	Datadir        string
	WWOverlaydir   string
	WWChrootdir    string
	WWProvisiondir string
	WWClientdir    string
}

func (this *BuildConfig) Upgrade() (upgraded *config.BuildConfig) {
	upgraded = new(config.BuildConfig)
	upgraded.Bindir = this.Bindir
	upgraded.Sysconfdir = this.Sysconfdir
	upgraded.Localstatedir = this.Localstatedir
	upgraded.Cachedir = this.Cachedir
	upgraded.Ipxesource = this.Ipxesource
	upgraded.Srvdir = this.Srvdir
	upgraded.Firewallddir = this.Firewallddir
	upgraded.Datadir = this.Datadir
	upgraded.WWOverlaydir = this.WWOverlaydir
	upgraded.WWChrootdir = this.WWChrootdir
	upgraded.WWProvisiondir = this.WWProvisiondir
	upgraded.WWClientdir = this.WWClientdir
	return upgraded
}

type WWClientConf struct {
	Port uint16 `yaml:"port"`
}

func (this *WWClientConf) Upgrade() (upgraded *config.WWClientConf) {
	upgraded = new(config.WWClientConf)
	upgraded.Port = this.Port
	return upgraded
}
