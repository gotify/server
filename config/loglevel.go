package config

import (
	"errors"

	"github.com/rs/zerolog"
)

// LogLevel type that provides helper methods for decoding.
type LogLevel zerolog.Level

// Decode decodes a string to a log level.
func (ll *LogLevel) Decode(value string) error {
	if level, err := zerolog.ParseLevel(value); err == nil {
		*ll = LogLevel(level)
		return nil
	}
	*ll = LogLevel(zerolog.InfoLevel)
	return errors.New("unknown log level")
}

// AsZeroLogLevel converts the LogLevel to a zerolog.Level.
func (ll LogLevel) AsZeroLogLevel() zerolog.Level {
	return zerolog.Level(ll)
}
