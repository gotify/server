package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGeneration(t *testing.T) {
	mockTokenFunc := Tokens("a", "b", "c")

	for _, expected := range []string{"a", "b", "c", "a", "b", "c"} {
		assert.Equal(t, expected, mockTokenFunc())
	}
}
