package data

import (
	"testing"

	"github.com/azrod/zr/pkg/format"
	"github.com/stretchr/testify/assert"
)

func TestDataHotReload(t *testing.T) {
	d := DataHotReload{
		LogLevel:  "debug",
		LogFormat: format.LogFormatJson,
	}
	d.JSONMarshal()

	assert.Equal(t, "debug", d.LogLevel)
	assert.Equal(t, format.LogFormatJson, d.LogFormat)

	assert.JSONEq(t, `{"log_level":"debug","log_format":"json"}`, d.JSONMarshal().String())

}
