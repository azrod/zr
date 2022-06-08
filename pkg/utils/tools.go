package utils

import (
	"errors"

	"github.com/rs/zerolog"
)

var (
	// Default level
	default_level = zerolog.InfoLevel
	ErrLogLevel   = errors.New("level not allowed, use info")
)

func ParseLogLevel(level string) (zerolog.Level, error) {

	zlevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return default_level, ErrLogLevel
	}

	return zlevel, nil

}
