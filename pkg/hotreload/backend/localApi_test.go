package backend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azrod/zr/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestBackendLocalAPI(t *testing.T) {
	t.Parallel()

	b, err := BackendLocalAPI()
	assert.Nil(t, err)

	assert.NotNil(t, b)
	assert.Equal(t, DefaultCfgAddress, b.cfg.Address)
	assert.Equal(t, DefaultCfgPort, b.cfg.Port)

	// Test with config
	c := ConfigBackendLocalAPI{
		Address: "",
		Port:    0,
	}

	b, err = BackendLocalAPI(c)
	assert.Nil(t, err)

	assert.NotNil(t, b)
	assert.Equal(t, DefaultCfgAddress, b.cfg.Address)
	assert.Equal(t, DefaultCfgPort, b.cfg.Port)

	errS := b.Shutdown()
	assert.Nil(t, errS)

}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	c := ConfigBackendLocalAPI{}
	c.defaultConfig()
	assert.Equal(t, DefaultCfgAddress, c.Address)
	assert.Equal(t, DefaultCfgPort, c.Port)

}

func TestGetSetRequests(t *testing.T) {
	t.Parallel()

	var x data.DataHotReload

	x.LogFormat = "json"
	x.LogLevel = "debug"

	// json Marshal
	y := x.JSONMarshal()

	b, err := BackendLocalAPI()
	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/get", ioutil.NopCloser(bytes.NewBuffer(y)))
	w := httptest.NewRecorder()

	b.setConfig(w, r)

	b.getConfig(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, []byte(`{"log_level":"debug","log_format":"json"}`), w.Body.Bytes())

	d := b.Reload()
	assert.Equal(t, "json", d.LogFormat)
	assert.Equal(t, "debug", d.LogLevel)
}
