package compat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const examplePluginPath = "github.com/gotify/server/plugin/example/echo"

func TestPluginInfoStringer(t *testing.T) {
	info := Info{
		ModulePath: examplePluginPath,
	}
	assert.Equal(t, examplePluginPath, info.String())
	info.Name = "test name"
	assert.Equal(t, "test name", info.String())
}
