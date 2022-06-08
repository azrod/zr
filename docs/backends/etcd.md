# Backend ETCD

The etcd backend is used to store the configuration in etcd.

```golang

zr.Setup(
	zr.WithCustomLevel("debug"),
	zr.WithCustomFormat("human"),
	zr.WithCustomHotReload(
		hr.WithBackendETCD(etcd.ConfigBackendETCD{
			Endpoints: []string{"http://etcd.databases:2379"},
			Path:      "/zr/basic/config",
		}),
	),
)

```
