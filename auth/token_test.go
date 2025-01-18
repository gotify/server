package auth

import (
	"crypto/rand"
	"fmt"
	"strings"
	"testing"

	"github.com/gotify/server/v2/test"
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

func TestHashTokenStable(t *testing.T) {
	salt1 := []byte("salt")
	salt2 := []byte("pepper")
	seen := make(map[string]bool)
	for _, plain := range []string{"", "a", "b", "c", "a\x00", "a\n"} {
		hash1, err := HashToken(plain, salt1)
		assert.NoError(t, err)
		hash1Again, err := HashToken(plain, salt1)
		assert.NoError(t, err)
		assert.Equal(t, hash1, hash1Again)
		hash2, err := HashToken(plain, salt2)
		assert.NoError(t, err)
		hash2Again, err := HashToken(plain, salt2)
		assert.NoError(t, err)
		assert.Equal(t, hash2, hash2Again)

		assert.NotEqual(t, hash1, hash2)
		assert.False(t, seen[hash1])
		assert.False(t, seen[hash2])
		seen[hash1] = true
		seen[hash2] = true
	}
}

func TestCompareToken(t *testing.T) {
	salt := []byte("salt")
	tokenPlain := GenerateApplicationToken()
	hashed, err := HashToken(tokenPlain, salt)
	assert.NoError(t, err)
	cmpPlain, upgPlain, err := CompareToken(tokenPlain, tokenPlain)
	assert.NoError(t, err)
	assert.True(t, cmpPlain)
	assert.NotEmpty(t, *upgPlain)

	cmpHashed, upgHashed, err := CompareToken(tokenPlain, hashed)
	assert.NoError(t, err)
	assert.True(t, cmpHashed)
	assert.Nil(t, upgHashed)
}

func TestGenerateNotExistingToken(t *testing.T) {
	count := 5
	token := GenerateNotExistingToken(func() string {
		return fmt.Sprint(count)
	}, func(token string) bool {
		count--
		return token != "0"
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
