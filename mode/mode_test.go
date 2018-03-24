package mode

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDevMode(t *testing.T) {
	Set(Dev)
	assert.Equal(t, Get(), Dev)
	assert.True(t, IsDev())
	assert.Equal(t, gin.Mode(), gin.DebugMode)
}

func TestTestDevMode(t *testing.T) {
	Set(TestDev)
	assert.Equal(t, Get(), TestDev)
	assert.True(t, IsDev())
	assert.Equal(t, gin.Mode(), gin.TestMode)
}

func TestProdMode(t *testing.T) {
	Set(Prod)
	assert.Equal(t, Get(), Prod)
	assert.False(t, IsDev())
	assert.Equal(t, gin.Mode(), gin.ReleaseMode)
}

func TestInvalidMode(t *testing.T) {
	assert.Panics(t, func() {
		Set("asdasda")
	})
}
