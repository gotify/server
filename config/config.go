package config

import "github.com/jinzhu/configor"

// Configuration the application config can be set per env variables or config file (config.yml).
type Configuration struct {
	Port int `default:"8080"`
	Database struct {
		Dialect    string `default:"sqlite3"`
		Connection string `default:"memo.db"`
	}
	DefaultUser struct {
		Name string `default:"admin"`
		Pass string `default:"admin"`
	}
}

// Get returns the configuration extracted from env variables or config file.
func Get() *Configuration {
	conf := new(Configuration)
	configor.New(&configor.Config{ENVPrefix: "MEMO"}).Load(conf, "config.yml", "/etc/memo/config.yml")
	return conf
}
