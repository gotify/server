package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockRandSource(t *testing.T) {
	UseTokenSource(&MockRandSource{
		Tokens: []string{"test1", "test2", "test3"},
	})
	defer UseCryptoRand()

	assert.Equal(t, "Ctest1", GenerateClientToken())
	assert.Equal(t, "Atest2", GenerateApplicationToken())
}
