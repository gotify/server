package config

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel_Decode_success(t *testing.T) {
	ll := new(LogLevel)
	err := ll.Decode("fatal")
	assert.Nil(t, err)
	assert.Equal(t, ll.AsZeroLogLevel(), zerolog.FatalLevel)
}

func TestLogLevel_Decode_fail(t *testing.T) {
	ll := new(LogLevel)
	err := ll.Decode("asdasdasdasdasdasd")
	assert.EqualError(t, err, "unknown log level")
	assert.Equal(t, ll.AsZeroLogLevel(), zerolog.InfoLevel)
}
