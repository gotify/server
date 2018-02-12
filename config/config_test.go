package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEnv(t *testing.T) {
	os.Setenv("GOTIFY_DEFAULTUSER_NAME", "jmattheis")
	conf := Get()
	assert.Equal(t, 8080, conf.Port, "should use defaults")
	assert.Equal(t, "jmattheis", conf.DefaultUser.Name, "should not use default but env var")
	os.Unsetenv("GOTIFY_DEFAULTUSER_NAME")
}

func TestConfigFile(t *testing.T) {
	file, err := os.Create("config.yml")
	defer func() {
		file.Close()
	}()
	assert.Nil(t, err)
	_, err = file.WriteString(`
port: 1234
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
	assert.Equal(t, 1234, conf.Port)
	assert.Equal(t, "nicories", conf.DefaultUser.Name)
	assert.Equal(t, "12345", conf.DefaultUser.Pass)
	assert.Equal(t, "mysql", conf.Database.Dialect)
	assert.Equal(t, "user name", conf.Database.Connection)

	assert.Nil(t, os.Remove("config.yml"))
}
