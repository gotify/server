package auth

import (
	"crypto/rand"
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestNewComplexToken(t *testing.T) {
	token := NewEnhancedToken("a")
	canonicalizedExpected := token.PublicForm()
	tokenParsed, err := ParseEnhancedToken(token.String())
	assert.NoError(t, err)
	tokenSigned, err := tokenParsed.Sign(12345)
	assert.NoError(t, err)
	tokenSignedStr := tokenSigned.String()
	assert.True(t, tokenSigned.ValidateTimestamp(12345+1))
	assert.False(t, tokenSigned.ValidateTimestamp(12345+maxTimestampDiffSeconds+1))
	canonicalizedActual := tokenSigned.PublicForm()
	assert.Equal(t, canonicalizedExpected, canonicalizedActual)

	tokenParsed, err = ParseEnhancedToken(tokenSignedStr)
	assert.NoError(t, err)
	canonicalizedActual = tokenParsed.PublicForm()
	assert.Equal(t, canonicalizedExpected, canonicalizedActual)
	_, err = tokenSigned.Sign(12345)
	assert.ErrorIs(t, err, errNoPrivateKey)
	tokenParsed, err = ParseEnhancedToken(tokenSignedStr)
	assert.NoError(t, err)
	canonicalizedActual = tokenParsed.PublicForm()
	assert.Equal(t, canonicalizedExpected, canonicalizedActual)
	var tokenSignedStrMutated string
	lastDotIdx := strings.LastIndex(tokenSignedStr, ".")
	if tokenSignedStr[lastDotIdx+1] != 'A' {
		tokenSignedStrMutated = tokenSignedStr[:lastDotIdx+1] + "A" + tokenSignedStr[lastDotIdx+2:]
	} else {
		tokenSignedStrMutated = tokenSignedStr[:lastDotIdx+1] + "B" + tokenSignedStr[lastDotIdx+2:]
	}
	_, err = ParseEnhancedToken(tokenSignedStrMutated)
	assert.ErrorIs(t, err, errInvalidToken)
}

func TestBadCryptoReaderPanics(t *testing.T) {
	assert.Panics(t, func() {
		randReader = iotest.ErrReader(errors.New("this reader cannot be read"))
		defer func() {
			randReader = rand.Reader
		}()
		randIntn(2)
	})
}
