package warewulfd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

var provisionSendTests = []struct {
	description string
	url         string
	body        string
	status      int
	ip          string
}{
	{"system overlay", "/overlay-system/00:00:00:ff:ff:ff", "system overlay", 200, "10.10.10.10:9873"},
	{"runtime overlay", "/overlay-runtime/00:00:00:ff:ff:ff", "runtime overlay", 200, "10.10.10.10:9873"},
	{"fake overlay", "/overlay-system/00:00:00:ff:ff:ff?overlay=fake", "", 400, "10.10.10.10:9873:9873"},
	{"specific overlay", "/overlay-system/00:00:00:ff:ff:ff?overlay=o1", "specific overlay", 200, "10.10.10.10:9873"},
	{"find shim", "/efiboot/shim.efi", "", 200, "10.10.10.10:9873"},
	{"find shim", "/efiboot/shim.efi", "", 404, "10.10.10.11:9873"},
	{"find grub", "/efiboot/grub.efi", "", 200, "10.10.10.10:9873"},
	{"find grub", "/efiboot/grub.efi", "", 404, "10.10.10.11:9873"},
	{"find initramfs", "/provision/00:00:00:ff:ff:ff?stage=initramfs", "", 200, "10.10.10.10:9873"},
	{"ipxe test with NetDevs and KernelOverrides", "/provision/00:00:00:00:00:ff?stage=ipxe", "1.1.1 ifname=net:00:00:00:00:00:ff ", 200, "10.10.10.12:9873"},
	{"find grub.cfg", "/efiboot/grub.cfg", "dracut", 200, "10.10.10.11:9873"},
}

func Test_ProvisionSend(t *testing.T) {
	env := testenv.New(t)

	env.WriteFile(t, "etc/warewulf/nodes.conf", `nodeprofiles:
  default:
    container name: suse
nodes:
  n1:
    network devices:
      default:
        hwaddr: 00:00:00:ff:ff:ff
    profiles:
    - default
  n2:
    network devices:
      default:
        hwaddr: 00:00:00:00:ff:ff
    container name: none
    tags:
      GrubMenuEntry: dracut
  n3:
    network devices:
      default:
        hwaddr: 00:00:00:00:00:ff
        device: net
    ipxe template: test
    kernel:
      override: 1.1.1`)
	// create a  arp file as for grub we look up the ip address through the arp cache

	env.WriteFile(t, "arpcache", `IP address       HW type     Flags       HW address            Mask     Device
10.10.10.10    0x1         0x2         00:00:00:ff:ff:ff     *        dummy
10.10.10.11    0x1         0x2         00:00:00:00:ff:ff     *        dummy
10.10.10.12    0x1         0x2         00:00:00:00:00:ff     *        dummy`)
	SetArpFile(path.Join(env.BaseDir, "arpcache"))
	conf := warewulfconf.Get()
	containerDir, imageDirErr := os.MkdirTemp(os.TempDir(), "ww-test-container-*")
	assert.NoError(t, imageDirErr)
	defer os.RemoveAll(containerDir)
	conf.Paths.WWChrootdir = containerDir

	sysConfDir, sysConfDirErr := os.MkdirTemp(os.TempDir(), "ww-test-sysconf-*")
	assert.NoError(t, sysConfDirErr)
	defer os.RemoveAll(sysConfDir)
	conf.Paths.Sysconfdir = sysConfDir

	assert.NoError(t, os.MkdirAll(path.Join(containerDir, "suse/rootfs/usr/lib64/efi"), 0700))
	{
		_, err := os.Create(path.Join(containerDir, "suse/rootfs/usr/lib64/efi", "shim.efi"))
		assert.NoError(t, err)
	}
	assert.NoError(t, os.MkdirAll(path.Join(containerDir, "suse/rootfs/usr/share/efi/x86_64/"), 0700))
	{
		_, err := os.Create(path.Join(containerDir, "suse/rootfs/usr/share/efi/x86_64/", "grub.efi"))
		assert.NoError(t, err)
	}
	assert.NoError(t, os.MkdirAll(path.Join(containerDir, "suse/rootfs/boot"), 0700))
	{
		_, err := os.Create(path.Join(containerDir, "suse/rootfs/boot", "initramfs-.img"))
		assert.NoError(t, err)
	}
	assert.NoError(t, os.MkdirAll(path.Join(conf.Paths.Sysconfdir, "warewulf/ipxe"), 0700))
	{
		assert.NoError(t, os.WriteFile(path.Join(conf.Paths.Sysconfdir, "warewulf/ipxe", "test.ipxe"), []byte("{{.KernelOverride}}{{range $devname, $netdev := .NetDevs}}{{if and $netdev.Hwaddr $netdev.Device}} ifname={{$netdev.Device}}:{{$netdev.Hwaddr}} {{end}}{{end}}"), 0600))
	}
	assert.NoError(t, os.MkdirAll(path.Join(conf.Paths.Sysconfdir, "warewulf/grub"), 0700))
	{
		assert.NoError(t, os.WriteFile(path.Join(conf.Paths.Sysconfdir, "warewulf/grub", "grub.cfg.ww"), []byte("{{ .Tags.GrubMenuEntry }}"), 0600))
	}

	dbErr := LoadNodeDB()
	assert.NoError(t, dbErr)

	secureFalse := false
	conf.Warewulf.SecureP = &secureFalse
	assert.NoError(t, os.MkdirAll(path.Join(conf.Paths.OverlayProvisiondir(), "n1"), 0700))
	assert.NoError(t, os.WriteFile(path.Join(conf.Paths.OverlayProvisiondir(), "n1", "__SYSTEM__.img"), []byte("system overlay"), 0600))
	assert.NoError(t, os.WriteFile(path.Join(conf.Paths.OverlayProvisiondir(), "n1", "__RUNTIME__.img"), []byte("runtime overlay"), 0600))
	assert.NoError(t, os.WriteFile(path.Join(conf.Paths.OverlayProvisiondir(), "n1", "o1.img"), []byte("specific overlay"), 0600))

	wwlog.SetLogLevel(wwlog.DEBUG)
	for _, tt := range provisionSendTests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			req.RemoteAddr = tt.ip
			w := httptest.NewRecorder()
			ProvisionSend(w, req)
			res := w.Result()
			defer res.Body.Close()

			data, readErr := io.ReadAll(res.Body)
			assert.NoError(t, readErr)
			if tt.body != "" {
				assert.Equal(t, tt.body, string(data))
			}
			assert.Equal(t, tt.status, res.StatusCode)
		})
	}
}
