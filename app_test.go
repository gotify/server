package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantCode int
		stdout   string // substring expected on stdout
		stderr   string // substring expected on stderr
	}{
		{"version", []string{"version"}, 0, "Version: ", ""},
		{"unknown command", []string{"bogus"}, 2, "", "unknown command"},
		{"unknown flag", []string{"--nope"}, 2, "", "not defined"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(c.args, &stdout, &stderr)
			assert.Equal(t, c.wantCode, code)
			assert.Contains(t, stdout.String(), c.stdout)
			assert.Contains(t, stderr.String(), c.stderr)
		})
	}
}
