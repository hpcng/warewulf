package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultWarewulfYaml(t *testing.T) {
	conf := New()

	assert.Equal(t, 9873, conf.Warewulf.Port)
	assert.True(t, conf.Warewulf.Secure())
	assert.Equal(t, 60, conf.Warewulf.UpdateInterval)
	assert.True(t, conf.Warewulf.AutobuildOverlays())
	assert.True(t, conf.Warewulf.EnableHostOverlay())
	assert.False(t, conf.Warewulf.Syslog())

	assert.True(t, conf.DHCP.Enabled())
	assert.Equal(t, "default", conf.DHCP.Template)
	assert.Empty(t, conf.DHCP.RangeStart)
	assert.Empty(t, conf.DHCP.RangeEnd)
	assert.Equal(t, "dhcpd", conf.DHCP.SystemdName)

	assert.True(t, conf.TFTP.Enabled())
	assert.NotEmpty(t, conf.TFTP.TftpRoot)
	assert.Equal(t, "tftp", conf.TFTP.SystemdName)
	assert.NotEmpty(t, conf.TFTP.IpxeBinaries["00:00"])
	assert.NotEmpty(t, conf.TFTP.IpxeBinaries["00:07"])
	assert.NotEmpty(t, conf.TFTP.IpxeBinaries["00:09"])
	assert.NotEmpty(t, conf.TFTP.IpxeBinaries["00:0B"])

	assert.True(t, conf.NFS.Enabled())
	assert.Empty(t, conf.NFS.ExportsExtended)
	assert.Equal(t, "nfsd", conf.NFS.SystemdName)

	assert.Equal(t, "/etc/resolv.conf", conf.MountsContainer[0].Source)
	assert.Equal(t, "/etc/resolv.conf", conf.MountsContainer[0].Dest)
	assert.False(t, conf.MountsContainer[0].ReadOnly())
	assert.Empty(t, conf.MountsContainer[0].Options)

	assert.NotEmpty(t, conf.Paths.Bindir)
	assert.NotEmpty(t, conf.Paths.Sysconfdir)
	assert.NotEmpty(t, conf.Paths.Localstatedir)
	assert.NotEmpty(t, conf.Paths.Srvdir)
	assert.NotEmpty(t, conf.Paths.Firewallddir)
	assert.NotEmpty(t, conf.Paths.Systemddir)
	assert.NotEmpty(t, conf.Paths.WWOverlaydir)
	assert.NotEmpty(t, conf.Paths.WWChrootdir)
	assert.NotEmpty(t, conf.Paths.WWProvisiondir)
	assert.NotEmpty(t, conf.Paths.WWClientdir)
	assert.NotEmpty(t, conf.Paths.Cachedir)
}

func TestInitializedFromFile(t *testing.T) {
	example_warewulf_conf := ""
	tempWarewulfConf, warewulfConfErr := os.CreateTemp("", "warewulf.conf-")
	assert.NoError(t, warewulfConfErr)
	defer os.Remove(tempWarewulfConf.Name())
	_, warewulfConfErr = tempWarewulfConf.Write([]byte(example_warewulf_conf))
	assert.NoError(t, warewulfConfErr)
	assert.NoError(t, tempWarewulfConf.Sync())

	conf := New()
	assert.False(t, conf.InitializedFromFile())
	assert.NoError(t, conf.Read(tempWarewulfConf.Name()))
	assert.True(t, conf.InitializedFromFile())
	assert.Equal(t, conf.GetWarewulfConf(), tempWarewulfConf.Name())
}

func TestExampleWarewulfYaml(t *testing.T) {
	example_warewulf_conf := `
ipaddr: 192.168.200.1
netmask: 255.255.255.0
network: 192.168.200.0
warewulf:
  port: 9873
  secure: false
  update interval: 60
  autobuild overlays: true
  host overlay: true
  syslog: false
dhcp:
  enabled: true
  range start: 192.168.200.50
  range end: 192.168.200.99
  systemd name: dhcpd
tftp:
  enabled: true
  systemd name: tftp
nfs:
  enabled: true
  export paths:
  - path: /home
    export options: rw,sync
  - path: /opt
    export options: ro,sync,no_root_squash
  systemd name: nfs-server
container mounts:
  - source: /etc/resolv.conf
    dest: /etc/resolv.conf
    readonly: true`

	tempWarewulfConf, warewulfConfErr := os.CreateTemp("", "warewulf.conf-")
	assert.NoError(t, warewulfConfErr)
	defer os.Remove(tempWarewulfConf.Name())
	_, warewulfConfErr = tempWarewulfConf.Write([]byte(example_warewulf_conf))
	assert.NoError(t, warewulfConfErr)
	assert.NoError(t, tempWarewulfConf.Sync())

	conf := New()
	assert.NoError(t, conf.Read(tempWarewulfConf.Name()))

	assert.Equal(t, "192.168.200.1", conf.Ipaddr)
	assert.Equal(t, "255.255.255.0", conf.Netmask)
	assert.Equal(t, "192.168.200.0", conf.Network)

	assert.Equal(t, 9873, conf.Warewulf.Port)
	assert.False(t, conf.Warewulf.Secure())
	assert.Equal(t, 60, conf.Warewulf.UpdateInterval)
	assert.True(t, conf.Warewulf.AutobuildOverlays())
	assert.True(t, conf.Warewulf.EnableHostOverlay())
	assert.False(t, conf.Warewulf.Syslog())

	assert.True(t, conf.DHCP.Enabled())
	assert.Equal(t, "192.168.200.50", conf.DHCP.RangeStart)
	assert.Equal(t, "192.168.200.99", conf.DHCP.RangeEnd)
	assert.Equal(t, "dhcpd", conf.DHCP.SystemdName)

	assert.True(t, conf.TFTP.Enabled())
	assert.Equal(t, "tftp", conf.TFTP.SystemdName)

	assert.True(t, conf.NFS.Enabled())
	assert.Equal(t, "/home", conf.NFS.ExportsExtended[0].Path)
	assert.Equal(t, "rw,sync", conf.NFS.ExportsExtended[0].ExportOptions)
	assert.Equal(t, "/opt", conf.NFS.ExportsExtended[1].Path)
	assert.Equal(t, "ro,sync,no_root_squash", conf.NFS.ExportsExtended[1].ExportOptions)
	assert.Equal(t, "nfs-server", conf.NFS.SystemdName)

	assert.Equal(t, "/etc/resolv.conf", conf.MountsContainer[0].Source)
	assert.Equal(t, "/etc/resolv.conf", conf.MountsContainer[0].Dest)
	assert.True(t, conf.MountsContainer[0].ReadOnly())
}

func TestCache(t *testing.T) {
	confOrig := New()
	confCached := Get()

	assert.Equal(t, 9873, confOrig.Warewulf.Port)
	assert.Equal(t, 9873, confCached.Warewulf.Port)

	confOrig.Warewulf.Port = 9999
	assert.Equal(t, 9999, confCached.Warewulf.Port)
	assert.Equal(t, 9999, Get().Warewulf.Port)

	New()
	assert.NotEqual(t, 9999, Get().Warewulf.Port)
}
