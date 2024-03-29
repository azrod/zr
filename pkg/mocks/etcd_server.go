/*
 * Copyright 2019-present Open Networking Foundation
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package mocks

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"go.etcd.io/etcd/server/v3/embed"
)

const (
	serverStartUpTimeout          = 10 * time.Second // Maximum time allowed to wait for the Etcd server to be ready
	defaultLocalPersistentStorage = "zr.test.embed.etcd"
)

// EtcdServer represents an emibedded Etcd server.  It is used for testing only.
type EtcdServer struct {
	server *embed.Etcd
}

func islogLevelValid(logLevel string) bool {
	valid := []string{"debug", "info", "warn", "error", "panic", "fatal"}
	for _, l := range valid {
		if l == logLevel {
			return true
		}
	}
	return false
}

/*
 * MKConfig creates an embedded Etcd config
 * :param configName: A name for this config
 * :param clientPort:  The port the etcd client will connect to (do not use 2379 for unit test)
 * :param peerPort:  The port the etcd server will listen for its peers (do not use 2380 for unit test)
 * :param localPersistentStorageDir: The name of a local directory which will hold the Etcd server data
 * :param logLevel: One of debug, info, warn, error, panic, or fatal. Default 'info'.
 */
func MKConfig(configName string, clientPort, peerPort int, localPersistentStorageDir string, logLevel string) *embed.Config {
	cfg := embed.NewConfig()
	cfg.Name = configName
	cfg.Dir = localPersistentStorageDir
	cfg.Logger = "zap"
	if !islogLevelValid(logLevel) {
		log.Fatalf("Invalid log level -%s", logLevel)
	}
	cfg.LogLevel = logLevel
	acurl, err := url.Parse(fmt.Sprintf("http://localhost:%d", clientPort))
	if err != nil {
		log.Fatalf("Invalid client port -%d", clientPort)
	}
	cfg.ACUrls = []url.URL{*acurl}
	cfg.LCUrls = []url.URL{*acurl}
	apurl, err := url.Parse(fmt.Sprintf("http://localhost:%d", peerPort))
	if err != nil {
		log.Fatalf("Invalid peer port -%d", peerPort)
	}
	cfg.LPUrls = []url.URL{*apurl}
	cfg.APUrls = []url.URL{*apurl}
	cfg.ClusterState = embed.ClusterStateFlagNew
	cfg.InitialCluster = cfg.Name + "=" + apurl.String()
	return cfg
}

// getDefaultCfg specifies the default config
func getDefaultCfg() *embed.Config {
	cfg := embed.NewConfig()
	cfg.Dir = defaultLocalPersistentStorage
	cfg.Logger = "zap"
	cfg.LogLevel = "error"
	return cfg
}

// StartEtcdServer creates and starts an embedded Etcd server.  A local directory to store data is created for the
// embedded server lifetime (for the duration of a unit test.  The server runs at localhost:2379.
func StartEtcdServer(cfg *embed.Config) *EtcdServer {
	// If the server is already running, just return
	if cfg == nil {
		cfg = getDefaultCfg()
	}
	// Remove the local directory as
	// a safeguard for the case where a prior test failed
	if err := os.RemoveAll(cfg.Dir); err != nil {
		log.Fatalf("Failure removing local directory %s", cfg.Dir)
	}
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	select {
	case <-e.Server.ReadyNotify():
		log.Print("Embedded Etcd server is ready!")
	case <-time.After(serverStartUpTimeout):
		e.Server.HardStop() // trigger a shutdown
		e.Close()
		log.Fatal("Embedded Etcd server took too long to start!")
	case err := <-e.Err():
		e.Server.HardStop() // trigger a shutdown
		e.Close()
		log.Fatalf("Embedded Etcd server errored out - %s", err)
	}
	return &EtcdServer{server: e}
}

// Stop closes the embedded Etcd server and removes the local data directory as well
func (es *EtcdServer) Stop() {
	if es != nil {
		storage := es.server.Config().Dir
		es.server.Server.HardStop()
		es.server.Close()
		if err := os.RemoveAll(storage); err != nil {
			log.Fatalf("Failure removing local directory %s", es.server.Config().Dir)
		}
	}
}
