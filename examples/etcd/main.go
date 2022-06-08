package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/azrod/zr"
	hr "github.com/azrod/zr/pkg/hotreload"
	"github.com/azrod/zr/pkg/hotreload/backend"
	"github.com/rs/zerolog/log"
)

func main() {

	zr.Setup(
		zr.WithCustomHotReload(
			hr.WithBackendETCD(backend.ConfigBackendETCD{
				Endpoints: []string{"http://localhost:2379"},
				Path:      "/zr/basic/config",
			}),
		),
	)

	log.Info().Msg("This is a test for level info")

	go func() {
		for {
			log.Info().Msg("This is a test for level info")
			log.Debug().Msg("This is a test for level debug")
			log.Warn().Msg("This is a test for level warn")
			log.Error().Msg("This is a test for level error")
			time.Sleep(10 * time.Second)
		}
	}()

	/*
		etcdctl --endpoints=localhost:2379 put /zr/basic/config '{"log_level":"info","log_format":"json"}'
	*/

	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	zr.Done()

}
