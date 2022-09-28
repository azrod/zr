package level

import "github.com/rs/zerolog"

type Level struct {
	LogLevel
}

type Options func(*Level) error
type LogLevel zerolog.Level

func ParseLogLevel(v string) (level LogLevel, err error) {
	x, err := zerolog.ParseLevel(v)
	if err != nil {
		return level, ErrLogLevel
	}
	return LogLevel(x), nil
}

func (l LogLevel) String() string {
	return zerolog.Level(l).String()
}

func (l *Level) SetLevel(level LogLevel) error {
	l.LogLevel = level
	zerolog.SetGlobalLevel(zerolog.Level(level))

	return nil
}

func Setup() *Level {

	l := &Level{
		LogLevel: LogLevel(zerolog.InfoLevel),
	}

	return l

}

func (l *Level) GetLevel() LogLevel {
	return l.LogLevel
}
