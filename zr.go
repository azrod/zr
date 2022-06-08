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

/*
WithCustomLevel is an option which sets up the custom log level.
zr is called on for custom log level.
if this option is not used the log level is info.
*/
func WithCustomLevel(level string) ExtraZrOptions {
	return func(t *zr) {
		t.level, _ = utils.ParseLogLevel(level)
	}
}

/*
WithCustomFormat is an option which sets up the custom log format.
zr is called on for custom log format.
Allowed format are json and human.
if this option is not used the log level is json.
*/
func WithCustomFormat(format string) ExtraZrOptions {
	return func(t *zr) {
		t.format = format
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
