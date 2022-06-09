package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/azrod/zr"
	hr "github.com/azrod/zr/pkg/hotreload"
	"github.com/azrod/zr/pkg/hotreload/backend"
	"github.com/rs/zerolog/log"
)

var (
	RootCAs *x509.CertPool
)

func init() {
	RootCAs = x509.NewCertPool()
}

func main() {

	tlsConfig, err := TlsConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create tls config")
	}

	zr.Setup(
		zr.WithCustomHotReload(
			hr.WithBackendETCD(backend.ConfigBackendETCD{
				Endpoints:   []string{"http://localhost:2379"},
				Path:        "/zr/basic/config",
				DialTimeout: time.Second * 5,
				TLS:         tlsConfig,
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

func TlsConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair("tls.crt", "tls.key")
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile("tls.ca")
	if err != nil {
		return nil, err
	}
	RootCAs.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            RootCAs,
		InsecureSkipVerify: true,
	}, nil

}
