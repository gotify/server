package migrate

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gotify/server/v2/config"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type oldConfig struct {
	Server struct {
		KeepAlivePeriodSeconds *int
		ListenAddr             *string
		Port                   *int
		SSL                    struct {
			Enabled         *bool
			RedirectToHTTPS *bool
			ListenAddr      *string
			Port            *int
			CertFile        *string
			CertKey         *string
			LetsEncrypt     struct {
				Enabled      *bool
				AcceptTOS    *bool
				Cache        *string
				DirectoryURL *string
				Hosts        []string
			}
		}
		ResponseHeaders map[string]string
		Stream          struct {
			PingPeriodSeconds *int
			AllowedOrigins    []string
		}
		Cors struct {
			AllowOrigins []string
			AllowMethods []string
			AllowHeaders []string
		}
		TrustedProxies []string
		SecureCookie   *bool
	}
	Database struct {
		Dialect    *string
		Connection *string
	}
	DefaultUser struct {
		Name *string
		Pass *string
	}
	PassStrength      *int
	UploadedImagesDir *string
	PluginsDir        *string
	Registration      *bool
	OIDC              struct {
		Enabled       *bool
		Issuer        *string
		ClientID      *string
		ClientSecret  *string
		UsernameClaim *string
		RedirectURL   *string
		AutoRegister  *bool
		Scopes        []string
	}
}

func Config(file string) (string, error) {
	if file == "" {
		return "", errors.New("migrate-config requires one argument: the path to the old config.yml")
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("cannot read config file %s: %w", file, err)
	}

	var migrated oldConfig
	if err := yaml.Unmarshal(data, &migrated); err != nil {
		return "", fmt.Errorf("cannot parse config file %s: %w", file, err)
	}

	content, err := godotenv.Marshal(buildEnv(migrated))
	if err != nil {
		return "", fmt.Errorf("cannot render config: %w", err)
	}

	return content, nil
}

func buildEnv(c oldConfig) map[string]string {
	out := map[string]string{}
	str := func(key string, value *string) {
		if value != nil {
			out[key] = *value
		}
	}
	num := func(key string, value *int) {
		if value != nil {
			out[key] = strconv.Itoa(*value)
		}
	}
	boolean := func(key string, value *bool) {
		if value != nil {
			out[key] = strconv.FormatBool(*value)
		}
	}
	list := func(key string, value []string) {
		if value != nil {
			out[key] = marshalList(value)
		}
	}
	headers := func(key string, value map[string]string) {
		if value != nil {
			out[key] = marshalMap(value)
		}
	}

	num(config.EnvServerKeepAlivePeriodSeconds, c.Server.KeepAlivePeriodSeconds)
	str(config.EnvServerListenAddr, c.Server.ListenAddr)
	num(config.EnvServerPort, c.Server.Port)
	boolean(config.EnvServerSSLEnabled, c.Server.SSL.Enabled)
	boolean(config.EnvServerSSLRedirectToHTTPS, c.Server.SSL.RedirectToHTTPS)
	str(config.EnvServerSSLListenAddr, c.Server.SSL.ListenAddr)
	num(config.EnvServerSSLPort, c.Server.SSL.Port)
	str(config.EnvServerSSLCertFile, c.Server.SSL.CertFile)
	str(config.EnvServerSSLCertKey, c.Server.SSL.CertKey)
	boolean(config.EnvServerSSLLetsEncryptEnabled, c.Server.SSL.LetsEncrypt.Enabled)
	boolean(config.EnvServerSSLLetsEncryptAcceptTOS, c.Server.SSL.LetsEncrypt.AcceptTOS)
	str(config.EnvServerSSLLetsEncryptCache, c.Server.SSL.LetsEncrypt.Cache)
	str(config.EnvServerSSLLetsEncryptDirectoryURL, c.Server.SSL.LetsEncrypt.DirectoryURL)
	list(config.EnvServerSSLLetsEncryptHosts, c.Server.SSL.LetsEncrypt.Hosts)
	headers(config.EnvServerResponseHeaders, c.Server.ResponseHeaders)
	num(config.EnvServerStreamPingPeriodSeconds, c.Server.Stream.PingPeriodSeconds)
	list(config.EnvServerStreamAllowedOrigins, c.Server.Stream.AllowedOrigins)
	list(config.EnvServerCorsAllowOrigins, c.Server.Cors.AllowOrigins)
	list(config.EnvServerCorsAllowMethods, c.Server.Cors.AllowMethods)
	list(config.EnvServerCorsAllowHeaders, c.Server.Cors.AllowHeaders)
	list(config.EnvServerTrustedProxies, c.Server.TrustedProxies)
	boolean(config.EnvServerSecureCookie, c.Server.SecureCookie)
	str(config.EnvDatabaseDialect, c.Database.Dialect)
	str(config.EnvDatabaseConnection, c.Database.Connection)
	str(config.EnvDefaultUserName, c.DefaultUser.Name)
	str(config.EnvDefaultUserPass, c.DefaultUser.Pass)
	num(config.EnvPassStrength, c.PassStrength)
	str(config.EnvUploadedImagesDir, c.UploadedImagesDir)
	str(config.EnvPluginsDir, c.PluginsDir)
	boolean(config.EnvRegistration, c.Registration)
	boolean(config.EnvOIDCEnabled, c.OIDC.Enabled)
	str(config.EnvOIDCIssuer, c.OIDC.Issuer)
	str(config.EnvOIDCClientID, c.OIDC.ClientID)
	str(config.EnvOIDCClientSecret, c.OIDC.ClientSecret)
	str(config.EnvOIDCUsernameClaim, c.OIDC.UsernameClaim)
	str(config.EnvOIDCRedirectURL, c.OIDC.RedirectURL)
	boolean(config.EnvOIDCAutoRegister, c.OIDC.AutoRegister)
	list(config.EnvOIDCScopes, c.OIDC.Scopes)
	return out
}

func marshalMap(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}

func marshalList(values []string) string {
	var sb strings.Builder
	writer := csv.NewWriter(&sb)
	writer.UseCRLF = false
	if err := writer.Write(values); err != nil {
		return ""
	}
	writer.Flush()
	return strings.TrimRight(sb.String(), "\n")
}
