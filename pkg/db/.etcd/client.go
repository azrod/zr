package etcd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"time"

	"github.com/azrod/zr/pkg/db"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ConfigDBETCD struct {
	// Endpoints is a list of ETCD Servers.
	Endpoints []string
	// DialTimeout is the timeout for failing to connect to a server.
	DialTimeout time.Duration
	// TLS is a reference to the TLS Config.
	TLS *tls.Config
	// Path get the value from etcd.
	Path string
}

/*
DBETCD is a backend for hot reload.
*/
type DBETCD struct {
	cli    *clientv3.Client
	Errors []error
}

/*
NewClient is a constructor for DBETCD.
*/
func NewClient(cfg ConfigDBETCD) (*DBETCD, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
		TLS:         cfg.TLS,
	})
	if err != nil {
		return nil, err
	}
	return &DBETCD{
		cli: cli,
	}, nil
}

/*
Get is a function which gets the value from etcd.
*/
func (b *DBETCD) Get(key string, timeout int) (db.DataHotReload, error) {

	// check if key is valid
	if b.keyInvalid(key) {
		return db.DataHotReload{}, ErrKeyInvalid
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	resp, err := b.cli.Get(ctx, key)
	if err != nil {
		return db.DataHotReload{}, err
	}
	if len(resp.Kvs) == 0 {
		return db.DataHotReload{}, nil
	} else {
		x := db.DataHotReload{}
		err = json.Unmarshal(resp.Kvs[0].Value, &x)
		return x, err
	}
}

/*
Set is a function which sets the value to etcd.
*/
func (b *DBETCD) Set(key string, value db.DataHotReload, timeout int) error {

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
func (b *DBETCD) Close() error {
	return b.cli.Close()
}

/*
KeyInvalid is a function which checks the key is valid or not.
*/
func (b *DBETCD) keyInvalid(key string) bool {
	return key == ""
}

/*
ValueInvalid is a function which checks the value is valid or not.
*/
func (b *DBETCD) valueInvalid(value string) bool {
	return value == ""
}
