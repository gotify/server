package config

import "regexp"

// CompileAllowedOrigins compiles the patterns as fully matching regexes.
func CompileAllowedOrigins(allowedOrigins []string) []*regexp.Regexp {
	var compiledAllowedOrigins []*regexp.Regexp
	for _, origin := range allowedOrigins {
		compiledAllowedOrigins = append(compiledAllowedOrigins, regexp.MustCompile("^(?:"+origin+")$"))
	}

	return compiledAllowedOrigins
}

// MatchesFully checks if any of the regexes matches the origin.
func MatchesFully(compiledOrigins []*regexp.Regexp, origin string) bool {
	for _, compiledOrigin := range compiledOrigins {
		if compiledOrigin.MatchString(origin) {
			return true
		}
	}
	return false
}
