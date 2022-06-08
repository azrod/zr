package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataHotReload(t *testing.T) {
	d := DataHotReload{
		LogLevel:  "debug",
		LogFormat: "json",
	}
	d.JSONMarshal()

	assert.Equal(t, "debug", d.LogLevel)
	assert.Equal(t, "json", d.LogFormat)

	assert.JSONEq(t, `{"log_level":"debug","log_format":"json"}`, d.JSONMarshal().String())

}
