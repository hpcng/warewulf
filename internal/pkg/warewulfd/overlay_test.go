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

var overlaySendTests = []struct {
	description string
	url         string
	body        string
	status      int
	ip          string
}{
	{"testfile should exist", "/test.template", "test", 200, "10.10.10.10:9873"},
	{"testfile shouldn't exist", "/fake.template", "", 404, "10.10.10.10:9873"},
	{"testfile.ww should exist and be rendered", "/test.template.ww?node=n1", "n1", 200, "10.10.10.10:9873"},
	{"testfile.ww should exist node not found", "/test.template.ww?node=n2", "", 404, "10.10.10.10:9873"},
	{"testfile.ww should exist even it isn't rendered", "/test.template?node=n1", "test", 200, "10.10.10.10:9873"},
	{"testfile.ww should not exist", "/test2.template.ww?node=n1", "", 404, "10.10.10.10:9873"},
	{"testfile.ww should exist in subdir and be rendered", "/dir1/test2.template.ww?node=n1", "n1", 200, "10.10.10.10:9873"},
}

func Test_OverlaySend(t *testing.T) {
	env := testenv.New(t)

	env.WriteFile(t, "etc/warewulf/nodes.conf", `nodeprofiles:
  default: {}
nodes:
  n1: {}
`)
	conf := warewulfconf.Get()
	assert.NoError(t, os.MkdirAll(path.Join(conf.Paths.WWOverlaydir, "wwroot/dir1"), 0700))
	assert.NoError(t, os.WriteFile(path.Join(conf.Paths.WWOverlaydir, "/wwroot", "test.template"), []byte("test"), 0600))
	assert.NoError(t, os.WriteFile(path.Join(conf.Paths.WWOverlaydir, "/wwroot", "test.template.ww"), []byte("{{.Id}}"), 0600))
	assert.NoError(t, os.WriteFile(path.Join(conf.Paths.WWOverlaydir, "/wwroot/dir1", "test2.template.ww"), []byte("{{.Id}}"), 0600))

	dbErr := LoadNodeDB()
	assert.NoError(t, dbErr)

	secureFalse := false
	conf.Warewulf.SecureP = &secureFalse
	wwlog.SetLogLevel(wwlog.DEBUG)
	for _, tt := range overlaySendTests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			req.RemoteAddr = tt.ip
			w := httptest.NewRecorder()
			OverlaySend(w, req)
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
