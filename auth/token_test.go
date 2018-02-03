package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenHavePrefix(t *testing.T) {
	for i := 0; i < 50; i++ {
		assert.True(t, strings.HasPrefix(GenerateApplicationToken(), "A"))
		assert.True(t, strings.HasPrefix(GenerateClientToken(), "C"))
	}
}
