package api

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGeneration(t *testing.T) {
	clientPub, clientPriv := generateClientToken()
	assert.Regexp(t, regexp.MustCompile(`^gtfyc\.(.+)$`), clientPub)
	assert.Regexp(t, regexp.MustCompile(`^gtfyc\.(.+)$`), clientPriv)
	applicationPub, applicationPriv := generateApplicationToken()
	assert.Regexp(t, regexp.MustCompile(`^gtfya\.(.+)$`), applicationPub)
	assert.Regexp(t, regexp.MustCompile(`^gtfya\.(.+)$`), applicationPriv)
	imageName := generateImageName()
	assert.Regexp(t, regexp.MustCompile(`^(.+)$`), imageName)
}
