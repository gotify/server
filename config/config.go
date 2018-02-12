package config

import "github.com/jinzhu/configor"

// Configuration is stuff that can be configured externally per env variables or config file (config.yml).
type Configuration struct {
	Port     int `default:"8080"`
	Database struct {
		Dialect    string `default:"sqlite3"`
		Connection string `default:"gotify.db"`
	}
	DefaultUser struct {
		Name string `default:"admin"`
		Pass string `default:"admin"`
	}
}

// Get returns the configuration extracted from env variables or config file.
func Get() *Configuration {
	conf := new(Configuration)
	configor.New(&configor.Config{ENVPrefix: "GOTIFY"}).Load(conf, "config.yml", "/etc/gotify/config.yml")
	return conf
}
