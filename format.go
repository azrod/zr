package zr

import (
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/azrod/common-go"
)

var (
	output io.Writer
)

func setFormat(format string) error {
	allowedFormat := []string{"json", "human"}
	_, ok := common.Find(allowedFormat, format)
	if !ok {
		return error_log_format
	}

	if format == "human" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: output, TimeFormat: time.RFC3339})
	} else {
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}

	return nil
}

func setLevel(level zerolog.Level) error {

	zerolog.SetGlobalLevel(level)

	return nil
}
