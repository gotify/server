package api

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGeneration(t *testing.T) {
	assert.Regexp(t, regexp.MustCompile("^C(.+)$"), generateClientToken())
	assert.Regexp(t, regexp.MustCompile("^A(.+)$"), generateApplicationToken())
	assert.Regexp(t, regexp.MustCompile("^(.+)$"), generateImageName())
}
