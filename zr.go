package zr

import (
	"os"
	"strconv"

	hr "github.com/azrod/zr/pkg/hotreload"
	"github.com/azrod/zr/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// Default level
	default_level = zerolog.InfoLevel
	// Default log format
	default_format = "json"
)

var (
	hrd *hr.HotReload
	z   *zr
)

type zr struct {
	level     zerolog.Level
	format    string
	hotReload *hr.HotReload
	// logger   *zerolog.Logger
}

type ExtraZrOptions func(*zr)
type FormatsOptions func(*zr)
type LevelsOptions func(*zr)

/*
WithCustomLevel is an option which sets up the custom log level.
zr is called on for custom log level.
if this option is not used the log level is info.
*/
func WithCustomLevel(level LevelsOptions) ExtraZrOptions {
	return func(t *zr) {
		level(t)
	}
}

/*
LevelPanic is a helper function which sets the log level to panic.
*/
func LevelPanic() LevelsOptions {
	return func(t *zr) {
		t.level = zerolog.PanicLevel
	}
}

/*
LevelFatal is a helper function which sets the log level to fatal.
*/
func LevelFatal() LevelsOptions {
	return func(t *zr) {
		t.level = zerolog.FatalLevel
	}
}

/*
LevelWarning is a helper function which sets the log level to warning.
*/
func LevelWarn() LevelsOptions {
	return func(t *zr) {
		t.level = zerolog.WarnLevel
	}
}

/*
LevelTrace is a helper function which sets the log level to trace.
*/
func LevelTrace() LevelsOptions {
	return func(t *zr) {
		t.level = zerolog.TraceLevel
	}
}

/*
LevelDebug is a helper function which sets the log level to debug.
*/
func LevelDebug() LevelsOptions {
	return func(t *zr) {
		t.level = zerolog.DebugLevel
	}
}

/*
LevelError is a helper function which sets the log level to error.
*/
func LevelError() LevelsOptions {
	return func(t *zr) {
		t.level = zerolog.ErrorLevel
	}
}

/*
LevelInfo is a helper function which sets the log level to info.
*/
func LevelInfo() LevelsOptions {
	return func(t *zr) {
		t.level = zerolog.InfoLevel
	}
}

/*
WithCustomFormat is an option which sets up the custom log format.
zr is called on for custom log format.
Allowed format are json and human.
if this option is not used the log level is json.
*/
func WithCustomFormat(format FormatsOptions) ExtraZrOptions {
	return func(t *zr) {
		format(t)
	}
}

/*
WithCustomHotReload is an option which sets up the custom hot reload.
zr is called on for custom hot reload.
*/
func WithCustomHotReload(opts ...hr.ExtraHotReloadOptions) ExtraZrOptions {
	return func(t *zr) {
		for _, opt := range opts {
			opt(t.hotReload)
		}
	}
}

/*
FormatHuman is an option which sets up the human log format.
*/
func FormatHuman() FormatsOptions {
	return func(t *zr) {
		t.format = "human"
	}
}

/*
FormatJSON is an option which sets up the json log format.
*/
func FormatJSON() FormatsOptions {
	return func(t *zr) {
		t.format = "json"
	}
}

/*
Setup is a constructor for zr.
It takes a list of options which can be used to customize the zr.
Available options are:

	WithCustomLevel(level string)
	WithCustomFormat

Default values are:

	level: info
	format: json
*/
func Setup(opts ...ExtraZrOptions) error {

	var err error

	hrd, err = hr.Setup()
	if err != nil {
		return err
	}

	z = &zr{
		level:     default_level,
		format:    default_format,
		hotReload: hrd,
	}

	for _, opt := range opts {
		opt(z)
	}

	setLevel(z.level)
	// Err return always nil

	output = os.Stderr

	err = setFormat(z.format)
	if err != nil {
		return err
	}

	// SRC https://github.com/rs/zerolog#add-file-and-line-number-to-log
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	log.Logger = log.With().Caller().Logger()

	// Setup hotReload
	if z.hotReload.Enabled {
		// Here we can setup hot reload

		z.hotReload.LoopHotReload()

		// for wait chan
		go func() {
			for {
				select {
				case x := <-z.hotReload.Chan:
					if x.LogFormat != z.format && x.LogFormat != "" {
						err := setFormat(x.LogFormat)
						if err != nil {
							log.Error().Msgf("error setting log format: %s", err)
						} else {
							log.Info().Msgf("log format changed from %s to: %s", z.format, x.LogFormat)
							z.format = x.LogFormat
						}
					}

					if x.LogLevel != "" {
						// Parse log level
						level, err := utils.ParseLogLevel(x.LogLevel)
						if err != nil {
							log.Error().Msgf("error parsing log level: %s", err)
						} else {
							if level != z.level {
								err := setLevel(level)
								if err != nil {
									log.Error().Msgf("error setting log level: %s", err)
								} else {
									log.Info().Msgf("log level changed from %s to: %s", z.level, level)
									z.level = level
								}
							}
						}
					}
				}
			}
		}()
	}

	return nil
}

func Done() {
	z.hotReload.DoneLoop <- true
}
