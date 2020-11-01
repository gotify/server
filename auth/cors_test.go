package auth

import (
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gotify/server/v2/config"
	"github.com/gotify/server/v2/mode"
	"github.com/stretchr/testify/assert"
)

func TestCorsConfig(t *testing.T) {
	mode.Set(mode.Prod)
	serverConf := config.Configuration{}
	serverConf.Server.Cors.AllowOrigins = []string{"http://test.com"}
	serverConf.Server.Cors.AllowHeaders = []string{"content-type"}
	serverConf.Server.Cors.AllowMethods = []string{"GET"}

	actual := CorsConfig(&serverConf)
	allowF := actual.AllowOriginFunc
	actual.AllowOriginFunc = nil // func cannot be checked with equal

	assert.Equal(t, cors.Config{
		AllowAllOrigins:        false,
		AllowHeaders:           []string{"content-type"},
		AllowMethods:           []string{"GET"},
		MaxAge:                 12 * time.Hour,
		AllowBrowserExtensions: true,
	}, actual)
	assert.NotNil(t, allowF)
	assert.True(t, allowF("http://test.com"))
	assert.False(t, allowF("https://test.com"))
	assert.False(t, allowF("https://other.com"))
}

func TestEmptyCorsConfigWithResponseHeaders(t *testing.T) {
	mode.Set(mode.Prod)
	serverConf := config.Configuration{}
	serverConf.Server.ResponseHeaders = map[string]string{"Access-control-allow-origin": "https://example.com"}

	actual := CorsConfig(&serverConf)
	assert.NotNil(t, actual.AllowOriginFunc)
	actual.AllowOriginFunc = nil // func cannot be checked with equal

	assert.Equal(t, cors.Config{
		AllowAllOrigins:        false,
		AllowOrigins:           []string{"https://example.com"},
		MaxAge:                 12 * time.Hour,
		AllowBrowserExtensions: true,
	}, actual)
}

func TestDevCorsConfig(t *testing.T) {
	mode.Set(mode.Dev)
	serverConf := config.Configuration{}
	serverConf.Server.Cors.AllowOrigins = []string{"http://test.com"}
	serverConf.Server.Cors.AllowHeaders = []string{"content-type"}
	serverConf.Server.Cors.AllowMethods = []string{"GET"}

	actual := CorsConfig(&serverConf)

	assert.Equal(t, cors.Config{
		AllowHeaders: []string{
			"X-Gotify-Key", "Authorization", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host",
		},
		AllowMethods:           []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		MaxAge:                 12 * time.Hour,
		AllowAllOrigins:        true,
		AllowBrowserExtensions: true,
	}, actual)
}
