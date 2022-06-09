# Backend ETCD

The [ETCD](https://etcd.io/) backend is used to store the configuration.
zr read database every time (Default interval is 60 seconds) and apply the configuration.

A complet example is available in the [examples/etcd](https://github.com/azrod/zr/tree/main/examples/etcd).

## Basic usage

```go linenums="1"

	zr.Setup(
		zr.WithCustomHotReload(
			hr.WithBackendETCD(backend.ConfigBackendETCD{
				Endpoints: []string{"http://localhost:2379"},
				Path:      "/zr/basic/config",
			}),
		),
	)

```

## Options for the backend

```go

backend.ConfigBackendETCD{
	Endpoints: []string{"http://localhost:2379"},
	Path:      "/zr/basic/config",
	DialTimeout: time.Second * 5,
	TLS:         tlsConfig,
}

```