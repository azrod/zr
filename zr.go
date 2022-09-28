package zr

import (
	"strconv"

	"github.com/azrod/zr/pkg/format"
	hr "github.com/azrod/zr/pkg/hotreload"
	"github.com/azrod/zr/pkg/level"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// Default level
	default_level = level.LogLevel(zerolog.InfoLevel)
	// Default log format
	default_format = format.LogFormatJson
)

var (
	hrd *hr.HotReload
	z   *zr
)

type zr struct {
	format    *format.Format
	level     *level.Level
	hotReload *hr.HotReload
	// logger   *zerolog.Logger
}

type ExtraZrOptions func(*zr)

/*
WithCustomLevel is an option which sets up the custom log level.
zr is called on for custom log level.
if this option is not used the log level is info.
*/
func Level(logLevel level.LogLevel) ExtraZrOptions {
	return func(t *zr) {
		t.level.SetLevel(logLevel)
	}
}

/*
Format is an option which sets up the custom log format.
zr is called on for custom log format.
Allowed format are json and human.
if this option is not used the log format is json.
*/
func Format(format format.LogFormat) ExtraZrOptions {
	return func(t *zr) {
		t.format.SetFormat(format)
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
	CustomFormatOptions is an option which sets up the custom log format.

zr is called on for custom log format.
*/
func CustomFormatOptions(format format.Options) ExtraZrOptions {
	return func(t *zr) {
		t.format.SetOptions(format)
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
		format:    format.Setup(),
		level:     level.Setup(),
		hotReload: hrd,
	}

	for _, opt := range opts {
		opt(z)
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
					if x.LogFormat.String() != z.format.GetFormat().String() && x.LogFormat.String() != "" {
						err := z.format.SetFormat(x.LogFormat)
						if err != nil {
							log.Error().Msgf("error setting log format: %s", err)
						} else {
							log.Info().Msgf("log format changed from %s to: %s", z.format.GetFormat().String(), x.LogFormat.String())
						}
					}

					if x.LogLevel != "" {
						// Parse log level
						lv, err := level.ParseLogLevel(x.LogLevel)
						if err != nil {
							log.Error().Msgf("error parsing log level: %s", err)
						} else {
							if lv.String() != z.level.String() {
								err := z.level.SetLevel(lv)
								if err != nil {
									log.Error().Msgf("error setting log level: %s", err)
								} else {
									log.Info().Msgf("log level changed from %s to: %s", z.level.String(), lv.String())
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
