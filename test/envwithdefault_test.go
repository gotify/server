package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	assert.Equal(t, GetEnv("THIS_ENV_SHOULD_NOT_EXIST", "default"), "default")
	assert.NoError(t, os.Setenv("THIS_ENV_IS_TEMPORARY", "test"))
	defer os.Unsetenv("THIS_ENV_IS_TEMPORARY")
	assert.Equal(t, GetEnv("THIS_ENV_IS_TEMPORARY", ""), "test")
}
