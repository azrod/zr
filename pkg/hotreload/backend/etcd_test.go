package backend

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/azrod/zr/pkg/data"
	mk "github.com/azrod/zr/pkg/mocks"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
)

var etcdServer *mk.EtcdServer
var client *BETCD

func setup() {
	clientPort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	peerPort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	etcdServer = mk.StartEtcdServer(mk.MKConfig("zr.mock.test", clientPort, peerPort, "zr.lib.mocks.etcd", "error"))
	if etcdServer == nil {
		log.Fatal("Embedded server failed to start")
	}

	clientAddr := fmt.Sprintf("localhost:%d", clientPort)
	client, err = BackendETCD(ConfigBackendETCD{
		Endpoints:   []string{clientAddr},
		DialTimeout: 5 * time.Second,
		Path:        "/zr/mock/config",
		TLS:         nil,
	})
	if err != nil || client == nil {
		etcdServer.Stop()
		log.Fatal("Failed to create an Etcd client")
	}

	err = client.Set("/zr/mock/config", data.DataHotReload{
		LogLevel:  "info",
		LogFormat: "json",
	}, 10)
	if err != nil {
		etcdServer.Stop()
		log.Fatal("Failed to set a key")
	}
}

func TestNewBackendETCD(t *testing.T) {

	type args struct {
		cfg ConfigBackendETCD
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				cfg: ConfigBackendETCD{
					Endpoints:   []string{"127.0.0.1:2379"},
					DialTimeout: 5 * time.Second,
					TLS:         nil,
				},
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				cfg: ConfigBackendETCD{
					Endpoints: []string{"127.0.0.1:2379"},
					TLS:       nil,
				},
			},
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				cfg: ConfigBackendETCD{
					Endpoints: nil,
					TLS:       nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := BackendETCD(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("newBackendETCD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.args.cfg.Endpoints, c.cfg.Endpoints)
				if tt.args.cfg.DialTimeout == 0 {
					assert.Equal(t, DefaultDialTimeout, c.cfg.DialTimeout)
				}
			}

		})
	}
}

func TestGetBackendETCD(t *testing.T) {

	kv := client.Reload()

	assert.NotNil(t, kv)
	assert.Equal(t, "info", kv.LogLevel)
	assert.Equal(t, "json", kv.LogFormat)

	err := client.Shutdown()
	assert.Nil(t, err)
}

func TestNewDBETCD(t *testing.T) {

	type args struct {
		cfg ConfigBackendETCD
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				cfg: ConfigBackendETCD{
					Endpoints:   []string{"127.0.0.1:2379"},
					DialTimeout: 5 * time.Second,
					TLS:         nil,
				},
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				cfg: ConfigBackendETCD{
					Endpoints:   nil,
					DialTimeout: 5 * time.Second,
					TLS:         nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("newDBETCD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestEtcdServerRWOLD(t *testing.T) {

	key := "myKey-1"
	value := data.DataHotReload{
		LogLevel:  "debug",
		LogFormat: "json",
	}

	err := client.Set(key, value, 10)
	assert.Nil(t, err)

	kv, err := client.Get(key, 10)
	assert.Nil(t, err)
	assert.NotNil(t, kv)
	assert.Equal(t, value, kv)
}

func shutdown() {
	if client != nil {
		client.Close()
	}
	if etcdServer != nil {
		etcdServer.Stop()
	}
}
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
