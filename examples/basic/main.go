package main

import (
	"os"
	"os/signal"

	"github.com/azrod/zr"
	"github.com/azrod/zr/pkg/format"
	"github.com/azrod/zr/pkg/level"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	zr.Setup(
		zr.Format(format.LogFormatHuman),            // This is the default, but we show it here for clarity.
		zr.Level(level.LogLevel(zerolog.InfoLevel)), // This is the default, but we show it here for clarity.
	)

	log.Info().Msg("This is a test for level info")
	log.Debug().Msg("This is a test for level debug")
	log.Warn().Msg("This is a test for level warn")
	log.Error().Msg("This is a test for level error")

	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	zr.Done()

}
