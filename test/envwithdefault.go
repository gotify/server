package test

import "os"

// GetEnv returns the value of an environment variable named by the key or the default value, if empty.
func GetEnv(name string, defaultVal string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal
	}
	return val
}
