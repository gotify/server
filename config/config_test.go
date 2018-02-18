package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEnv(t *testing.T) {
	os.Setenv("GOTIFY_DEFAULTUSER_NAME", "jmattheis")
	os.Setenv("GOTIFY_SERVER_SSL_LETSENCRYPT_HOSTS", "- push.example.tld\n- push.other.tld")
	conf := Get()
	assert.Equal(t, 80, conf.Server.Port, "should use defaults")
	assert.Equal(t, "jmattheis", conf.DefaultUser.Name, "should not use default but env var")
	assert.Equal(t, []string{"push.example.tld", "push.other.tld"}, conf.Server.SSL.LetsEncrypt.Hosts)
	os.Unsetenv("GOTIFY_DEFAULTUSER_NAME")
	os.Unsetenv("GOTIFY_SERVER_SSL_LETSENCRYPT_HOSTS")
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
database:
  dialect: mysql
  connection: user name
defaultuser:
  name: nicories
  pass: 12345
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

	assert.Nil(t, os.Remove("config.yml"))
}
