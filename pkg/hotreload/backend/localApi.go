package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/azrod/zr/pkg/data"
	"github.com/gorilla/mux"
)

const (
	DefaultCfgAddress = "127.0.0.1"
	DefaultCfgPort    = 8080
)

type ConfigBackendLocalAPI struct {
	Address string
	Port    int
}

type BLocalAPI struct {
	cfg  *ConfigBackendLocalAPI
	data data.DataHotReload
	srv  *http.Server
}

func BackendLocalAPI(cfg ...ConfigBackendLocalAPI) (*BLocalAPI, error) {

	// Get if cfg is defined
	var c ConfigBackendLocalAPI
	if len(cfg) > 0 {
		c = cfg[0]

		// Set default values
		if c.Address == "" {
			c.Address = DefaultCfgAddress
		}

		if c.Port == 0 {
			c.Port = DefaultCfgPort
		}

	} else {
		c.defaultConfig()
	}

	x := &BLocalAPI{
		cfg:  &c,
		data: data.DataHotReload{},
	}

	err := x.HttpServer()

	return x, err
}

func (c *ConfigBackendLocalAPI) defaultConfig() {
	c.Address = DefaultCfgAddress
	c.Port = DefaultCfgPort
}

func (b *BLocalAPI) srvHttpServer() {

	b.srv = &http.Server{
		Addr: b.cfg.Address + ":" + strconv.Itoa(b.cfg.Port),
	}

}

func (b *BLocalAPI) HttpServer() error {
	// Setup HTTP Server for local API
	r := mux.NewRouter()
	r.HandleFunc("/get", b.getConfig).Methods("GET")
	r.HandleFunc("/set", b.setConfig).Methods("POST")

	// Start HTTP Server
	b.srvHttpServer()
	b.srv.Handler = r

	log.Printf("Starting HTTP Server (%s)", b.srv.Addr)
	go func() {
		if err := b.srv.ListenAndServe(); err != http.ErrServerClosed {
			// return err
			fmt.Printf("HTTP Server ListenAndServe: %v", err)
		}
	}()

	return nil
}

func (b *BLocalAPI) Shutdown() error {
	return b.srv.Shutdown(context.Background())
}

func (b *BLocalAPI) Reload() data.DataHotReload {
	return b.data
}

func (b *BLocalAPI) setConfig(w http.ResponseWriter, r *http.Request) {

	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse body
	var x data.DataHotReload
	err = json.Unmarshal(body, &x)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set data
	b.data = x

	w.WriteHeader(http.StatusOK)
}

func (b *BLocalAPI) getConfig(w http.ResponseWriter, r *http.Request) {

	// Get data
	x := b.data

	// Marshal data
	body, err := json.Marshal(x)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write data
	w.WriteHeader(http.StatusOK)
	w.Write(body)

}
