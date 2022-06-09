# Backend Local API

Local API backend is used to provide a basic HTTP API to read and write the configuration.
Configuration is stored in memory and is not persisted.

Local API backend is a **default** backend. It is used if no other backend is specified.
Actually authentification and TLS are not supported.



A complet example is available in the [examples/local_api](https://github.com/azrod/zr/tree/main/examples/local_api).

```go linenums="1"

	zr.Setup(
		zr.WithCustomHotReload(
			hr.WithBackendLocalAPI(backend.ConfigBackendLocalAPI{
				Address: "127.0.0.1",
				Port:    6583,
			}),
		),
	)

```

By default HTTP server is configured to listen on `127.0.0.1:8080`.