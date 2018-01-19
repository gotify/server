package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordSuccess(t *testing.T) {
	password := CreatePassword("secret")
	assert.Equal(t, true, ComparePassword(password, []byte("secret")))
}

func TestPasswordFailure(t *testing.T) {
	password := CreatePassword("secret")
	assert.Equal(t, false, ComparePassword(password, []byte("secretx")))
}

func TestBCryptFailure(t *testing.T) {
	strength = 12312 // invalid value
	assert.Panics(t, func() { CreatePassword("secret") })
}
