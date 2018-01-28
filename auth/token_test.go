package auth

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTokenHavePrefix(t *testing.T) {
	for i := 0; i < 50; i++ {
		assert.True(t, strings.HasPrefix(GenerateApplicationToken(), "A"))
		assert.True(t, strings.HasPrefix(GenerateClientToken(), "C"))
	}
}
