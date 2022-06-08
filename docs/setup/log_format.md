# Log format

A log format is a string that can be used to format log messages. The default log format is `json`.

Two formats are available :

* json
* human

## JSON format

**Output**

```json
{"level": "info","time": "2022-01-01T14:00:00+00:00","caller": "demo.go:17","message": "hello world"}
```

**Setup**
    
``` go linenums="1" hl_lines="2"

    zr.Setup(
        zr.WithCustomFormat("json"), // (1)
    )

```

1.  :octicons-info-24: This is a default value



## Human format

**Output**

```sh
2022-01-01T14:00:00+00:00 INF demo.go:17 > hello world
```

**Setup**

```go linenums="1" hl_lines="2"

    zr.Setup(
        zr.WithCustomFormat("human"), 
    )

```
