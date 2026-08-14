package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gotify/server/v2/mode"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestConfigEnv(t *testing.T) {
	mode.Set(mode.TestDev)
	t.Setenv("GOTIFY_DEFAULTUSER_NAME", "jmattheis")
	t.Setenv("GOTIFY_SERVER_SSL_LETSENCRYPT_HOSTS", "push.example.tld,push.other.tld")
	t.Setenv(
		"GOTIFY_SERVER_RESPONSEHEADERS",
		`{"Access-Control-Allow-Origin":"*","Access-Control-Allow-Methods":"GET,POST"}`,
	)
	t.Setenv("GOTIFY_SERVER_CORS_ALLOWORIGINS", ".+.example.com,otherdomain.com")
	t.Setenv("GOTIFY_SERVER_CORS_ALLOWMETHODS", "GET,POST")
	t.Setenv("GOTIFY_SERVER_CORS_ALLOWHEADERS", "Authorization,content-type")
	t.Setenv("GOTIFY_SERVER_STREAM_ALLOWEDORIGINS", ".+.example.com,otherdomain.com")

	conf, _ := Get()
	assert.Equal(t, 80, conf.Server.Port, "should use defaults")
	assert.Equal(t, "jmattheis", conf.DefaultUser.Name, "should not use default but env var")
	assert.Equal(t, []string{"push.example.tld", "push.other.tld"}, conf.Server.SSL.LetsEncrypt.Hosts)
	assert.Equal(t, "*", conf.Server.ResponseHeaders["Access-Control-Allow-Origin"])
	assert.Equal(t, "GET,POST", conf.Server.ResponseHeaders["Access-Control-Allow-Methods"])
	assert.Equal(t, []string{".+.example.com", "otherdomain.com"}, conf.Server.Cors.AllowOrigins)
	assert.Equal(t, []string{"GET", "POST"}, conf.Server.Cors.AllowMethods)
	assert.Equal(t, []string{"Authorization", "content-type"}, conf.Server.Cors.AllowHeaders)
	assert.Equal(t, []string{".+.example.com", "otherdomain.com"}, conf.Server.Stream.AllowedOrigins)
}

func TestLocalAuthDisabled(t *testing.T) {
	tests := []struct {
		name   string
		env    map[string]string
		fatals []FutureLog
	}{
		{
			name: "with oidc",
			env:  map[string]string{EnvLocalAuthEnabled: "false", EnvOIDCEnabled: "true"},
		},
		{
			name:   "without oidc",
			env:    map[string]string{EnvLocalAuthEnabled: "false"},
			fatals: []FutureLog{futureFatal("either local authentication or OIDC must be enabled")},
		},
		{
			name: "with registration",
			env: map[string]string{
				EnvLocalAuthEnabled: "false",
				EnvOIDCEnabled:      "true",
				EnvRegistration:     "true",
			},
			fatals: []FutureLog{futureFatal("registration requires local authentication to be enabled")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mode.Set(mode.TestDev)
			for key, value := range tc.env {
				t.Setenv(key, value)
			}

			conf, logs := Get()
			assert.False(t, conf.LocalAuthEnabled)

			var fatals []FutureLog
			for _, entry := range logs {
				if entry.Level == zerolog.FatalLevel {
					fatals = append(fatals, entry)
				}
			}
			assert.Equal(t, tc.fatals, fatals)
		})
	}
}

func TestFile(t *testing.T) {
	mode.Set(mode.TestDev)
	dir := t.TempDir()
	passPath := filepath.Join(dir, "pass")
	hostsPath := filepath.Join(dir, "hosts")
	assert.Nil(t, os.WriteFile(passPath, []byte("filesecret\n"), 0o600))
	assert.Nil(t, os.WriteFile(hostsPath, []byte("a.example.com,b.example.com"), 0o600))

	t.Setenv("GOTIFY_DEFAULTUSER_PASS_FILE", passPath)
	t.Setenv("GOTIFY_SERVER_SSL_LETSENCRYPT_HOSTS_FILE", hostsPath)

	conf, _ := Get()
	assert.Equal(t, "filesecret", conf.DefaultUser.Pass)
	assert.Equal(t, []string{"a.example.com", "b.example.com"}, conf.Server.SSL.LetsEncrypt.Hosts)
}

func TestGotifyConfigFile(t *testing.T) {
	mode.Set(mode.TestDev)
	dir := t.TempDir()
	configPath := filepath.Join(dir, "custom.env")
	assert.Nil(t, os.WriteFile(configPath, []byte("GOTIFY_DEFAULTUSER_NAME=fromfile\n"), 0o600))

	t.Setenv("GOTIFY_CONFIG_FILE", configPath)

	conf, _ := Get()
	assert.Equal(t, "fromfile", conf.DefaultUser.Name)
}

func TestAddSlash(t *testing.T) {
	mode.Set(mode.TestDev)
	t.Setenv("GOTIFY_UPLOADEDIMAGESDIR", "../data/images")
	conf, _ := Get()
	assert.Equal(t, "../data/images"+string(filepath.Separator), conf.UploadedImagesDir)
}

func TestNotAddSlash(t *testing.T) {
	mode.Set(mode.TestDev)
	t.Setenv("GOTIFY_UPLOADEDIMAGESDIR", "../data/")
	conf, _ := Get()
	assert.Equal(t, "../data/", conf.UploadedImagesDir)
}

func TestParseList(t *testing.T) {
	const env = "GOTIFY_TEST_PARSELIST"

	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{name: "escaped quotes", raw: `"a,b","c""d",e`, want: []string{`a,b`, `c"d`, `e`}},
		{name: "lazy bare quote", raw: `a"b,c`, want: []string{`a"b`, `c`}},
		{name: "lazy quote in quoted field", raw: `"ab"cd",test`, want: []string{`ab"cd`, `test`}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(env, tc.raw)

			var got []string
			assert.Nil(t, parseList(&got, env))
			assert.Equal(t, tc.want, got)
		})
	}
}
