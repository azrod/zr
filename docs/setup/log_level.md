# Log level

A log level is a string that can be used to filter log messages. The default log level is `info`.
Log levels are defined in the `zerolog` package. 

The following levels are available :

* debug
* info
* warn
* error
* fatal
* panic

```go hl_lines="6" linenums="1"

(...)

func main() {

	zr.Setup(
        zr.WithCustomLevel("debug"),
    )

	log.Debug().Msg("hello world")

}

```