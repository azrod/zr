package zr

import (
	"testing"

	hr "github.com/azrod/zr/pkg/hotreload"

	"github.com/azrod/zr/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rzajac/zltest"
	"github.com/stretchr/testify/assert"
)

func Test_ServiceLogsProperly_LogFormat(t *testing.T) {
	// --- Given ---
	// Crate zerolog test helper.
	tst := zltest.New(t)

	output = tst

	err := setFormat("json")
	assert.Nil(t, err)
	assert.Equal(t, zerolog.InfoLevel, default_level)

	// Configure zerolog and pas tester as a writer.
	log := zerolog.New(tst).With().Timestamp().Logger()

	// Log something.
	log.Info().Int("key0", 123).Str("str0", "string").Msg("Hello world!")

	// Test if log messages were generated properly.
	ent := tst.LastEntry()
	ent.ExpNum("key0", 123)
	ent.ExpStr("str0", "string")
	ent.ExpMsg("Hello world!")
	ent.ExpLevel(zerolog.InfoLevel)
}

// TODO LogFormatInvalid

func Test_ServiceLogsProperly_LogLevel(t *testing.T) {
	var err error

	// --- Given ---
	// Crate zerolog test helper.
	tst := zltest.New(t)

	output = tst

	lv, err := utils.ParseLogLevel("debug")
	assert.Nil(t, err, "Parse log level should not return error")
	assert.Equal(t, zerolog.DebugLevel, lv)

	setLevel(lv)

	// Configure zerolog and pas tester as a writer.
	log := zerolog.New(tst).With().Timestamp().Logger()

	// Log something.
	log.Debug().Int("key0", 123).Str("str0", "string").Msg("Hello world!")

	// Test if log messages were generated properly.
	ent := tst.LastEntry()
	ent.ExpNum("key0", 123)
	ent.ExpStr("str0", "string")
	ent.ExpMsg("Hello world!")
	ent.ExpLevel(zerolog.DebugLevel)
}

func Test_ServiceLogsProperly_WithCustomLevel_Invalid(t *testing.T) {

	lv, err := utils.ParseLogLevel("unknown")
	assert.Equal(t, err, utils.ErrLogLevel, "Level should be allowed")
	assert.Equal(t, zerolog.InfoLevel, lv, "Level should be set")
	assert.NotNil(t, err, "Parse log level return error for unknown level")

}

func Test_Setup(t *testing.T) {
	err := Setup()
	assert.Nil(t, err, "Setup should not return error")

	assert.Equal(t, zerolog.InfoLevel, default_level)
	assert.Equal(t, "json", default_format)

	Done()
}

func Test_WithCustomInterval(t *testing.T) {
	interval := 10

	err := Setup(
		WithCustomHotReload(
			hr.WithCustomInterval(interval),
		),
	)

	assert.Equal(t, z.hotReload.Interval, interval, "Hot reload interval should be set")
	assert.Nil(t, err, "Setup should not return error")
}

func Test_WithNoHotReload(t *testing.T) {
	err := Setup(
		WithCustomHotReload(
			hr.WithNoHotReload(),
		),
	)

	assert.Equal(t, z.hotReload.Enabled, false, "Hot reload should be disabled")
	assert.Nil(t, err, "Setup should not return error")
}

func Test_WithCustomLogLevel(t *testing.T) {
	err := Setup(
		WithCustomLevel("debug"),
	)

	assert.Nil(t, err, "Setup should not return error")
	assert.Equal(t, zerolog.GlobalLevel(), zerolog.DebugLevel, "Log level should be set")
}

func Test_WithCustomLogLevel_Invalid(t *testing.T) {
	Setup(
		WithCustomLevel("unknown"),
	)

	assert.Equal(t, z.level, default_level, "Log level should be set")

}

func Test_WithCustomFormat(t *testing.T) {
	err := Setup(
		WithCustomFormat("human"),
	)

	assert.Nil(t, err, "Setup should not return error")
	assert.Equal(t, z.format, "human", "Log format should be set")
}

func Test_WithCustomFormat_Invalid(t *testing.T) {
	err := Setup(
		WithCustomFormat("unknown"),
	)
	assert.NotNil(t, err, "Setup should not return error")

}

// 	// --- Given ---
// 	// Crate zerolog test helper.
// 	tst := zltest.New(t)

// 	output = tst

// 	// Configure zerolog and pas tester as a writer.
// 	log := zerolog.New(tst).With().Timestamp().Logger()

// 	// Log something.
// 	log.Info().Int("key0", 123).Str("str0", "string").Msg("Hello world!")

// 	// Test if log messages were generated properly.
// 	ent := tst.LastEntry()
// 	ent.ExpNum("key0", 123)
// 	ent.ExpStr("str0", "string")
// 	ent.ExpMsg("Hello world!")
// 	ent.ExpLevel(zerolog.InfoLevel)

// 	// --- When ---
// 	// Setup hot reload.
// 	err := Setup(
// 		WithCustomHotReload(
// 			hr.WithCustomInterval(1),
// 		),
// 	)
// 	assert.Nil(t, err, "Setup should not return error")

// 	// Start hot reload.
// 	go z.loop()

// 	// --- Then ---
// 	// Log something.
// 	log.Info().Int("key1", 123).Str("str1", "string").Msg("Hello world!")

// 	// Test if log messages were generated properly.
// 	ent = tst.LastEntry()
// 	ent.ExpNum("key1", 123)
// 	ent.ExpStr("str1", "string")
// 	ent.ExpMsg("Hello world!")
// 	ent.ExpLevel(zerolog.InfoLevel)
// }
