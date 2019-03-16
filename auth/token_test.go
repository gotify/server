package auth

import (
	"crypto/rand"
	"fmt"
	"strings"
	"testing"

	"github.com/gotify/server/test"

	"github.com/stretchr/testify/assert"
)

func TestTokenHavePrefix(t *testing.T) {
	for i := 0; i < 50; i++ {
		assert.True(t, strings.HasPrefix(GenerateApplicationToken(), "A"))
		assert.True(t, strings.HasPrefix(GenerateClientToken(), "C"))
		assert.True(t, strings.HasPrefix(GeneratePluginToken(), "P"))
		assert.NotEmpty(t, GenerateImageName())
	}
}

func TestGenerateNotExistingToken(t *testing.T) {
	count := 5
	token := GenerateNotExistingToken(func() string {
		return fmt.Sprint(count)
	}, func(token string) bool {
		count--
		if token == "0" {
			return false
		}
		return true
	})
	assert.Equal(t, "0", token)
}

func TestBadCryptoReaderPanics(t *testing.T) {
	assert.Panics(t, func() {
		randReader = test.UnreadableReader()
		defer func() {
			randReader = rand.Reader
		}()
		randIntn(2)
	})
}
