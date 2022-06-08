package hr

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/azrod/zr/pkg/hotreload/backend"
	mk "github.com/azrod/zr/pkg/mocks"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
)

func TestGetHostname(t *testing.T) {

	defer func() { osHostname = os.Hostname }()
	osHostname = func() (string, error) { return "", errors.New("fail") }
	got, err := fromHostname()
	if err == nil {
		t.Errorf("getHostname() = (%v, nil), want error", got)
	}
	assert.NotNil(t, err)

	defer func() { osHostname = os.Hostname }()
	osHostname = func() (string, error) { return "mock.server", nil }
	hn, err := fromHostname()
	if err != nil {
		t.Errorf("getHostname() = (%v, nil), want error : %s", hn, err)
	}
	assert.Equal(t, "mock.server", hn)
	assert.Nil(t, err)
}

func TestWithCustomInterval(t *testing.T) {
	x := HotReload{}
	WithCustomInterval(1)(&x)

	assert.Equal(t, 1, x.Interval)

}

func TestWithNoHotReload(t *testing.T) {
	x := HotReload{}
	WithNoHotReload()(&x)

	assert.Equal(t, false, x.Enabled)
}

func TestWithBackendETCD(t *testing.T) {

	var etcdServer *mk.EtcdServer
	var client *backend.BETCD

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

	x, err := Setup()
	if err != nil {
		etcdServer.Stop()
		assert.Nil(t, err)
	}

	WithBackendETCD(backend.ConfigBackendETCD{
		Endpoints:   []string{clientAddr},
		DialTimeout: 5 * time.Second,
		TLS:         nil,
		Path:        "/zr/mock/test",
	})(x)

	go func() {
		x.LoopHotReload()
	}()
	runtime.Gosched()

	x.DoneLoop <- true

	if client != nil {
		client.Close()
	}
	if etcdServer != nil {
		etcdServer.Stop()
	}

}

func TestWithBackendLocalAPI(t *testing.T) {

	x, err := Setup()
	assert.Nil(t, err)

	WithBackendLocalAPI(backend.ConfigBackendLocalAPI{
		Address: "127.0.0.1",
		Port:    6789,
	})(x)

	go func() {
		x.LoopHotReload()
	}()
	runtime.Gosched()

	x.DoneLoop <- true

}
