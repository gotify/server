package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompileAllowedOrigins(t *testing.T) {
	assert.Equal(t, 0, len(CompileAllowedOrigins([]string{})))
	assert.Equal(t, 3, len(CompileAllowedOrigins([]string{"^.*$", "", "abc"})))
}

func TestMatchesFully(t *testing.T) {
	compiledOrigins := CompileAllowedOrigins([]string{"gotify\\.net|push\\.gotify\\.net", "other\\.gotify\\.net"})

	assert.True(t, MatchesFully(compiledOrigins, "gotify.net"))
	assert.True(t, MatchesFully(compiledOrigins, "push.gotify.net"))
	assert.True(t, MatchesFully(compiledOrigins, "other.gotify.net"))
	assert.False(t, MatchesFully(compiledOrigins, "gotify.net.evil.net"))
	assert.False(t, MatchesFully(compiledOrigins, "evil-gotify.net"))
}
