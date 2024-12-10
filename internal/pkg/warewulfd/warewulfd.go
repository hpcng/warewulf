package warewulfd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

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

func defaultHandler() *slashFix {
	var wwHandler http.ServeMux
	wwHandler.HandleFunc("/provision/", ProvisionSend)
	wwHandler.HandleFunc("/ipxe/", ProvisionSend)
	wwHandler.HandleFunc("/efiboot/", ProvisionSend)
	wwHandler.HandleFunc("/kernel/", ProvisionSend)
	wwHandler.HandleFunc("/container/", ProvisionSend)
	wwHandler.HandleFunc("/overlay-system/", ProvisionSend)
	wwHandler.HandleFunc("/overlay-runtime/", ProvisionSend)
	wwHandler.HandleFunc("/status", StatusSend)
	return &slashFix{&wwHandler}
}

func RunServer() error {
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

	err := LoadNodeDB()
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

	conf := warewulfconf.Get()
	daemonPort := conf.Warewulf.Port

	apiHandler := apiHandler()
	defaultHandler := defaultHandler()
	dispatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			apiHandler.ServeHTTP(w, r)
		} else {
			defaultHandler.ServeHTTP(w, r)
		}
	})
	err = http.ListenAndServe(":"+strconv.Itoa(daemonPort), dispatchHandler)

	if err != nil {
		return fmt.Errorf("could not start listening service: %w", err)
	}

	return nil
}
