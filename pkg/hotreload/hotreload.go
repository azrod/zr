package hr

import (
	"fmt"
	"os"
	"time"

	"github.com/azrod/zr/pkg/data"
	"github.com/azrod/zr/pkg/hotreload/backend"
)

var (
	osHostname = os.Hostname

	// Default hot reload
	default_hotReload = true
	// Default interval
	default_interval = 60
)

type Reloader interface {
	Reload() data.DataHotReload
	Shutdown() error
}

type HotReload struct {
	Enabled     bool                     // hot reload is enabled
	Ticker      func() *time.Ticker      // ticker for hot reload
	Interval    int                      // interval for requesting value
	DoneLoop    chan bool                // done loop for hot reload
	HotReloader func() (Reloader, error) // hot Reloader
	Chan        chan data.DataHotReload  // chan for hot reload
}

type ExtraHotReloadOptions func(*HotReload)

/*
WithBackendETCD is an option which sets up the backend ETCD for hot reload.
*/
func WithBackendETCD(cfg backend.ConfigBackendETCD) ExtraHotReloadOptions {
	return func(t *HotReload) {
		t.Enabled = true
		t.HotReloader = func() (Reloader, error) {
			return backend.BackendETCD(cfg)
		}
	}
}

/*
WithBackendLocalAPI is an option which sets up the backend local API for hot reload.
*/
func WithBackendLocalAPI(cfg backend.ConfigBackendLocalAPI) ExtraHotReloadOptions {
	return func(t *HotReload) {
		t.Enabled = true
		t.HotReloader = func() (Reloader, error) {
			return backend.BackendLocalAPI(cfg)
		}
	}
}

/*
WithNoHotReload is an option which sets up the hot reload is disabled.
if this option is not used the default hot reload is disabled.
*/
func WithNoHotReload() ExtraHotReloadOptions {
	return func(t *HotReload) {
		t.Enabled = false
	}
}

/*
WithCustomInterval is an option which sets up the custom interval for hot reload.
if this option is not used the default interval is used.
*/
func WithCustomInterval(interval int) ExtraHotReloadOptions {
	return func(t *HotReload) {
		t.Interval = interval
	}
}

/*
loopHotReload is a goroutine which is used to hot reload the zr.
*/
func (h *HotReload) LoopHotReload() {

	h.Ticker = func() *time.Ticker {
		return time.NewTicker(time.Duration(h.Interval) * time.Second)
	}

	// Start hot reload loop

	x, err := h.HotReloader() // PB HERE
	if err != nil {
		fmt.Printf("hot reload error: %v", err)
	}

	go func() {
		for {
			select {
			case <-h.Ticker().C:
				h.Chan <- x.Reload()
			case <-h.DoneLoop:
				h.Ticker().Stop()
				x.Shutdown()
				return
			}
		}
	}()
}

/*
Setup is a function which sets up the hot reload.
*/
func Setup() (*HotReload, error) {

	x := HotReload{
		Enabled:  default_hotReload,
		Interval: default_interval,
		DoneLoop: make(chan bool),
		Chan:     make(chan data.DataHotReload),
	}

	x.HotReloader = func() (Reloader, error) {
		return backend.BackendLocalAPI()
	}

	return &x, nil
}

/*
getPodName is a function which is used to get the pod name.
*/
func fromHostname() (string, error) {
	// get hostname
	hn, err := osHostname()
	if err != nil {
		return "unknow", err
	}

	return hn, nil
}
