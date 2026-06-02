package api

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGeneration(t *testing.T) {
	clientPub, clientPriv := generateClientToken()
	assert.Regexp(t, regexp.MustCompile(`^gtfy_client\.(.+)$`), clientPub)
	assert.Regexp(t, regexp.MustCompile(`^gtfy_client\.(.+)$`), clientPriv)
	applicationPub, applicationPriv := generateApplicationToken()
	assert.Regexp(t, regexp.MustCompile(`^gtfy_app\.(.+)$`), applicationPub)
	assert.Regexp(t, regexp.MustCompile(`^gtfy_app\.(.+)$`), applicationPriv)
	imageName := generateImageName()
	assert.Regexp(t, regexp.MustCompile(`^(.+)$`), imageName)
}
