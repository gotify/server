package password

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordSuccess(t *testing.T) {
	password, err := CreatePassword("secret", 5)
	require.NoError(t, err)
	assert.Equal(t, true, ComparePassword(password, []byte("secret")))
}

func TestPasswordFailure(t *testing.T) {
	password, err := CreatePassword("secret", 5)
	require.NoError(t, err)
	assert.Equal(t, false, ComparePassword(password, []byte("secretx")))
}

func TestBCryptoTooLongErrorIsReturned(t *testing.T) {
	_, err := CreatePassword(strings.Repeat("a", 100), 5)
	assert.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
}
