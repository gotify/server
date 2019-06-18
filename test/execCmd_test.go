package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecShell(t *testing.T) {
	assert.NoError(t, ExecShell("exit 0"))
}
