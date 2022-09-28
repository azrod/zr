package format

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Format struct {
	output io.Writer
	format LogFormat
}

type Options func(*Format) error
type LogFormat string

const (
	LogFormatJson  LogFormat = "json"
	LogFormatHuman LogFormat = "human"
)

// ParseLogFormat parses the log format
// Returning error if the format not a known
func (f *Format) ParseLogFormat(v string) (format LogFormat, err error) {
	switch v {
	case "json":
		return LogFormatJson, nil
	case "human":
		return LogFormatHuman, nil
	default:
		return format, ErrLogFormat
	}
}

func (l LogFormat) String() string { return string(l) }

func (f *Format) SetFormat(lf LogFormat) error {
	switch lf {
	case LogFormatHuman:
		log.Logger = zerolog.New(f.output).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: f.output})
		f.format = LogFormatHuman
	case LogFormatJson:
		log.Logger = zerolog.New(f.output).With().Timestamp().Logger()
		f.format = LogFormatJson
	default:
		return ErrLogFormat
	}

	return nil
}

func Setup() *Format {

	f := &Format{
		output: os.Stdout,
		format: LogFormatJson,
	}

	return f
}

func (f *Format) SetOptions(opts ...Options) error {

	for _, opt := range opts {
		opt(f)
	}

	return nil
}

func CustomOutput(output io.Writer) Options {
	return func(f *Format) error {
		f.output = output
		return nil
	}
}

func (f *Format) GetFormat() LogFormat {
	return f.format
}
