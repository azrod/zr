# Getting Started

## Installation 

```sh 
    go get -u github.com/azrod/zr
```

## Basic usage

```go linenums="1"

    package main
    
    import (
        "github.com/azrod/zr"
        "github.com/rs/zerolog/log"
    )
    
    func main() {
        zr.Setup()

        log.Info().Msg("Hello world")
    }

```

**Output**
```json
{"level":"info","time":"2022-06-04T14:36:24+02:00","caller":"basic.go:15","message":"hello world"}
```

## Examples

A lot of examples are available in the [examples](https://github.com/azrod/zr/tree/main/examples) directory.

