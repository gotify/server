package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEnv(t *testing.T) {
	os.Setenv("GOTIFY_DEFAULTUSER_NAME", "jmattheis")
	os.Setenv("GOTIFY_SERVER_SSL_LETSENCRYPT_HOSTS", "- push.example.tld\n- push.other.tld")
	os.Setenv("GOTIFY_SERVER_RESPONSEHEADERS",
		"Access-Control-Allow-Origin: \"*\"\nAccess-Control-Allow-Methods: \"GET,POST\"",
	)
	os.Setenv("GOTIFY_SERVER_STREAM_ALLOWEDORIGINS", "- \".+.example.com\"\n- \"otherdomain.com\"")

	conf := Get()
	assert.Equal(t, 80, conf.Server.Port, "should use defaults")
	assert.Equal(t, "jmattheis", conf.DefaultUser.Name, "should not use default but env var")
	assert.Equal(t, []string{"push.example.tld", "push.other.tld"}, conf.Server.SSL.LetsEncrypt.Hosts)
	assert.Equal(t, "*", conf.Server.ResponseHeaders["Access-Control-Allow-Origin"])
	assert.Equal(t, "GET,POST", conf.Server.ResponseHeaders["Access-Control-Allow-Methods"])
	assert.Equal(t, []string{".+.example.com", "otherdomain.com"}, conf.Server.Stream.AllowedOrigins)

	os.Unsetenv("GOTIFY_DEFAULTUSER_NAME")
	os.Unsetenv("GOTIFY_SERVER_SSL_LETSENCRYPT_HOSTS")
	os.Unsetenv("GOTIFY_SERVER_RESPONSEHEADERS")
	os.Unsetenv("GOTIFY_SERVER_STREAM_ALLOWEDORIGINS")
}

func TestAddSlash(t *testing.T) {
	os.Setenv("GOTIFY_UPLOADEDIMAGESDIR", "../data/images")
	conf := Get()
	assert.Equal(t, "../data/images"+string(filepath.Separator), conf.UploadedImagesDir)
	os.Unsetenv("GOTIFY_UPLOADEDIMAGESDIR")
}

func TestNotAddSlash(t *testing.T) {
	os.Setenv("GOTIFY_UPLOADEDIMAGESDIR", "../data/")
	conf := Get()
	assert.Equal(t, "../data/", conf.UploadedImagesDir)
	os.Unsetenv("GOTIFY_UPLOADEDIMAGESDIR")
}

func TestFileWithSyntaxErrors(t *testing.T) {
	file, err := os.Create("config.yml")
	defer func() {
		file.Close()
	}()
	assert.Nil(t, err)
	_, err = file.WriteString(`
sdgsgsdfgsdfg
`)
	file.Close()
	assert.Nil(t, err)
	assert.Panics(t, func() {
		Get()
	})

	assert.Nil(t, os.Remove("config.yml"))
}

func TestConfigFile(t *testing.T) {
	file, err := os.Create("config.yml")
	defer func() {
		file.Close()
	}()
	assert.Nil(t, err)
	_, err = file.WriteString(`
server:
  port: 1234
  ssl:
    port: 3333
    letsencrypt:
      hosts:
        - push.example.tld
  responseheaders:
    Access-Control-Allow-Origin: "*"
    Access-Control-Allow-Methods: "GET,POST"
  stream:
    allowedorigins:
      - ".+.example.com"
      - "otherdomain.com"
database:
  dialect: mysql
  connection: user name
defaultuser:
  name: nicories
  pass: 12345
pluginsdir: data/plugins
`)
	file.Close()
	assert.Nil(t, err)
	conf := Get()
	assert.Equal(t, 1234, conf.Server.Port)
	assert.Equal(t, 3333, conf.Server.SSL.Port)
	assert.Equal(t, []string{"push.example.tld"}, conf.Server.SSL.LetsEncrypt.Hosts)
	assert.Equal(t, "nicories", conf.DefaultUser.Name)
	assert.Equal(t, "12345", conf.DefaultUser.Pass)
	assert.Equal(t, "mysql", conf.Database.Dialect)
	assert.Equal(t, "user name", conf.Database.Connection)
	assert.Equal(t, "*", conf.Server.ResponseHeaders["Access-Control-Allow-Origin"])
	assert.Equal(t, "GET,POST", conf.Server.ResponseHeaders["Access-Control-Allow-Methods"])
	assert.Equal(t, []string{".+.example.com", "otherdomain.com"}, conf.Server.Stream.AllowedOrigins)
	assert.Equal(t, "data/plugins", conf.PluginsDir)

	assert.Nil(t, os.Remove("config.yml"))
}
