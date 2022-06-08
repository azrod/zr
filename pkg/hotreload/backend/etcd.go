package backend

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"time"

	"github.com/azrod/zr/pkg/data"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	DefaultPathEtcd    = "/zerolog/config"
	DefaultDialTimeout = 5 * time.Second

	ErrEndpointsEmpty = errors.New("endpoints is empty")
	ErrKeyInvalid     = errors.New("key is invalid")
	ErrValueInvalid   = errors.New("value is invalid")
)

type ConfigBackendETCD struct {
	// Endpoints is a list of ETCD Servers.
	Endpoints []string
	// DialTimeout is the timeout for failing to connect to a server.
	DialTimeout time.Duration
	// TLS is a reference to the TLS Config.
	TLS *tls.Config
	// Path get the value from etcd.
	Path string
}

type BETCD struct {
	cli *clientv3.Client
	cfg ConfigBackendETCD
}

/*
hotReloadFromETCD is a function which sets up the hot reload from ETCD.
*/
func BackendETCD(cfg ConfigBackendETCD) (*BETCD, error) {

	if cfg.Path == "" {
		cfg.Path = DefaultPathEtcd
	}

	if cfg.Endpoints == nil {
		return nil, ErrEndpointsEmpty
	}

	if cfg.DialTimeout == time.Duration(0*time.Second) {
		cfg.DialTimeout = DefaultDialTimeout
	}

	b, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &BETCD{
		b,
		cfg,
	}, nil

}

/*
Reload is a function which reloads the hot reload from ETCD.
*/
func (b *BETCD) Reload() data.DataHotReload {
	x, _ := b.Get(b.cfg.Path, 5)
	return x
}

/*
Shutdown is a function which shuts down the hot reload from ETCD.
*/
func (b *BETCD) Shutdown() error {
	return nil
}

/*
NewClient is a constructor for DBETCD.
*/
func NewClient(cfg ConfigBackendETCD) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
		TLS:         cfg.TLS,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

/*
Get is a function which gets the value from etcd.
*/
func (b *BETCD) Get(key string, timeout int) (data.DataHotReload, error) {

	// check if key is valid
	if b.keyInvalid(key) {
		return data.DataHotReload{}, ErrKeyInvalid
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	resp, err := b.cli.Get(ctx, key)
	if err != nil {
		return data.DataHotReload{}, err
	}
	if len(resp.Kvs) == 0 {
		return data.DataHotReload{}, nil
	} else {
		x := data.DataHotReload{}
		err = json.Unmarshal(resp.Kvs[0].Value, &x)
		return x, err
	}
}

/*
Set is a function which sets the value to etcd.
*/
func (b *BETCD) Set(key string, value data.DataHotReload, timeout int) error {

	// check if key is valid
	if b.keyInvalid(key) {
		return ErrKeyInvalid
	}

	// check if value is valid
	if b.valueInvalid(value.JSONMarshal().String()) {
		return ErrValueInvalid
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	_, err := b.cli.Put(ctx, key, value.JSONMarshal().String())
	return err

}

/*
Close is a function which closes the connection to etcd.
*/
func (b *BETCD) Close() error {
	return b.cli.Close()
}

/*
KeyInvalid is a function which checks the key is valid or not.
*/
func (b *BETCD) keyInvalid(key string) bool {
	return key == ""
}

/*
ValueInvalid is a function which checks the value is valid or not.
*/
func (b *BETCD) valueInvalid(value string) bool {
	return value == ""
}
