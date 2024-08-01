package warewulfd

import (
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	warewulfconf "github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

// TODO: https://github.com/danderson/netboot/blob/master/pixiecore/dhcp.go
// TODO: https://github.com/pin/tftp
/*
wrapper type for the server mux as shim requests http://efiboot//grub.efi
which is filtered out by http to `301 Moved Permanently` what
shim.efi can't handle. So filter out `//` before they hit go/http.
Makes go/http more to behave like apache
*/
type slashFix struct {
	mux http.Handler
}

/*
Filter out the '//'
*/
func (h *slashFix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.Replace(r.URL.Path, "//", "/", -1)
	h.mux.ServeHTTP(w, r)
}

func RunServer() error {
	err := DaemonInitLogging()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize logging")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for range c {
			wwlog.Warn("Received SIGHUP, reloading...")
			err := LoadNodeDB()
			if err != nil {
				wwlog.Error("Could not load node DB: %s", err)
			}

			err = LoadNodeStatus()
			if err != nil {
				wwlog.Error("Could not prepopulate node status DB: %s", err)
			}
		}
	}()

	err = LoadNodeDB()
	if err != nil {
		wwlog.Error("Could not load database: %s", err)
	}

	err = LoadNodeStatus()
	if err != nil {
		wwlog.Error("Could not prepopulate node status DB: %s", err)
	}

	if err != nil {
		wwlog.Warn("couldn't copy default shim: %s", err)
	}
	var wwHandler http.ServeMux
	wwHandler.HandleFunc("/provision/", ProvisionSend)
	wwHandler.HandleFunc("/ipxe/", ProvisionSend)
	wwHandler.HandleFunc("/efiboot/", ProvisionSend)
	wwHandler.HandleFunc("/kernel/", ProvisionSend)
	wwHandler.HandleFunc("/kmods/", ProvisionSend)
	wwHandler.HandleFunc("/container/", ProvisionSend)
	wwHandler.HandleFunc("/overlay-system/", ProvisionSend)
	wwHandler.HandleFunc("/overlay-runtime/", ProvisionSend)
	wwHandler.HandleFunc("/status", StatusSend)
	wwHandler.HandleFunc("/container-by-name/", ContainerByNameSend)

	conf := warewulfconf.Get()

	daemonPort := conf.Warewulf.Port
	/*
		wwlog.Serv("Starting HTTPD REST service on port %d", daemonPort)
		s := &http.Server{
			Addr:         ":" + strconv.Itoa(daemonPort),
			Handler:      &slashFix{&wwHandler},
			ReadTimeout:  10 * time.Second,
			IdleTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		err = s.ListenAndServe()
	*/
	err = http.ListenAndServe(":"+strconv.Itoa(daemonPort), &slashFix{&wwHandler})

	if err != nil {
		return errors.Wrap(err, "Could not start listening service")
	}

	return nil
}
