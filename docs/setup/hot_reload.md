# Hot Reload 

Hot reloading is a feature that allows you to change your zerolog configuration without restarting the application.

Hot reload is enabled by default.

By default backend [local API](../backends/local_api.md) is used. No need to configure anything.

## Backends 

Multiple backends are available : 

* [local API](../backends/local_api.md)
* [ETCD](../backends/etcd.md)

## Disable hot reload

``` go linenums="1" hl_lines="9 10 11"

import (
	"github.com/azrod/zr"
	hr "github.com/azrod/zr/pkg/hotreload"
)

func main() {

	zr.Setup(
		zr.WithCustomHotReload(
			hr.WithNoHotReload(),
		),
	)

```