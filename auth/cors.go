package auth

import (
	"github.com/gin-contrib/cors"
	"github.com/gotify/server/config"
	"github.com/gotify/server/mode"
	"regexp"
	"strings"
	"time"
)

// CorsConfig generates a config to use in gin cors middleware based on server configuration
func CorsConfig(conf *config.Configuration) cors.Config {
	corsConf := cors.Config{
		MaxAge: 12 * time.Hour,
	}
	if mode.IsDev() {
		corsConf.AllowAllOrigins = true
		corsConf.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"}
		corsConf.AllowHeaders = []string{"X-Gotify-Key", "Authorization", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host"}
	} else {
		compiledOrigins := compileAllowedCORSOrigins(conf.Server.Cors.AllowOrigins)
		corsConf.AllowMethods = conf.Server.Cors.AllowMethods
		corsConf.AllowHeaders = conf.Server.Cors.AllowHeaders
		corsConf.AllowOriginFunc = func(origin string) bool {
			for _, compiledOrigin := range compiledOrigins {
				if compiledOrigin.Match([]byte(strings.ToLower(origin))) {
					return true
				}
			}
			return false
		}
	}

	return corsConf
}

func compileAllowedCORSOrigins(allowedOrigins []string) []*regexp.Regexp {
	var compiledAllowedOrigins []*regexp.Regexp
	for _, origin := range allowedOrigins {
		compiledAllowedOrigins = append(compiledAllowedOrigins, regexp.MustCompile(origin))
	}

	return compiledAllowedOrigins
}
