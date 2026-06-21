package config

import (
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

type LetsEncrypt struct {
	Enabled      bool
	AcceptTOS    bool
	Cache        string
	DirectoryURL string
	Hosts        []string
}

type SSL struct {
	Enabled         bool
	RedirectToHTTPS bool
	ListenAddr      string
	Port            int
	CertFile        string
	CertKey         string
	LetsEncrypt     LetsEncrypt
}

type Stream struct {
	PingPeriodSeconds int
	AllowedOrigins    []string
}

type Cors struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

type Server struct {
	KeepAlivePeriodSeconds int
	ListenAddr             string
	Port                   int
	SSL                    SSL
	ResponseHeaders        map[string]string
	Stream                 Stream
	Cors                   Cors
	TrustedProxies         []string
	SecureCookie           bool
}

type Database struct {
	Dialect    string
	Connection string
}

type DefaultUser struct {
	Name string
	Pass string
}

type OIDC struct {
	Enabled        bool
	Issuer         string
	ClientID       string
	ClientSecret   string
	UsernameClaim  string
	RedirectURL    string
	AutoRegister   bool
	LinkByUsername bool
	Scopes         []string
}

type Configuration struct {
	LogLevel          LogLevel
	Server            Server
	Database          Database
	DefaultUser       DefaultUser
	PassStrength      int
	UploadedImagesDir string
	PluginsDir        string
	Registration      bool
	OIDC              OIDC
	NoColor           string
}

// Get returns the configuration extracted from env variables.
func Get() (*Configuration, []FutureLog) {
	c := &Configuration{
		LogLevel: LogLevel(zerolog.InfoLevel),
		Server: Server{
			Port: 80,
			SSL: SSL{
				RedirectToHTTPS: true,
				Port:            443,
				LetsEncrypt: LetsEncrypt{
					Cache: "data/certs",
				},
			},
			Stream: Stream{
				PingPeriodSeconds: 45,
			},
		},
		Database: Database{
			Dialect:    "sqlite3",
			Connection: "data/gotify.db",
		},
		DefaultUser: DefaultUser{
			Name: "admin",
			Pass: "admin",
		},
		PassStrength:      10,
		UploadedImagesDir: "data/images",
		PluginsDir:        "data/plugins",
		OIDC: OIDC{
			UsernameClaim: "preferred_username",
			AutoRegister:  true,
			Scopes:        []string{"openid", "profile", "email"},
		},
	}

	logs := loadFiles()

	add := func(err error) {
		if err != nil {
			logs = append(logs, futureFatal(err.Error()))
		}
	}

	add(parseLogLevel(&c.LogLevel, EnvLogLevel))

	add(parseInt(&c.Server.KeepAlivePeriodSeconds, EnvServerKeepAlivePeriodSeconds))
	add(parseString(&c.Server.ListenAddr, EnvServerListenAddr))
	add(parseInt(&c.Server.Port, EnvServerPort))

	add(parseBool(&c.Server.SSL.Enabled, EnvServerSSLEnabled))
	add(parseBool(&c.Server.SSL.RedirectToHTTPS, EnvServerSSLRedirectToHTTPS))
	add(parseString(&c.Server.SSL.ListenAddr, EnvServerSSLListenAddr))
	add(parseInt(&c.Server.SSL.Port, EnvServerSSLPort))
	add(parseString(&c.Server.SSL.CertFile, EnvServerSSLCertFile))
	add(parseString(&c.Server.SSL.CertKey, EnvServerSSLCertKey))

	add(parseBool(&c.Server.SSL.LetsEncrypt.Enabled, EnvServerSSLLetsEncryptEnabled))
	add(parseBool(&c.Server.SSL.LetsEncrypt.AcceptTOS, EnvServerSSLLetsEncryptAcceptTOS))
	add(parseString(&c.Server.SSL.LetsEncrypt.Cache, EnvServerSSLLetsEncryptCache))
	add(parseString(&c.Server.SSL.LetsEncrypt.DirectoryURL, EnvServerSSLLetsEncryptDirectoryURL))
	add(parseList(&c.Server.SSL.LetsEncrypt.Hosts, EnvServerSSLLetsEncryptHosts))

	add(parseMap(&c.Server.ResponseHeaders, EnvServerResponseHeaders))

	add(parseInt(&c.Server.Stream.PingPeriodSeconds, EnvServerStreamPingPeriodSeconds))
	add(parseList(&c.Server.Stream.AllowedOrigins, EnvServerStreamAllowedOrigins))

	add(parseList(&c.Server.Cors.AllowOrigins, EnvServerCorsAllowOrigins))
	add(parseList(&c.Server.Cors.AllowMethods, EnvServerCorsAllowMethods))
	add(parseList(&c.Server.Cors.AllowHeaders, EnvServerCorsAllowHeaders))

	add(parseList(&c.Server.TrustedProxies, EnvServerTrustedProxies))
	add(parseBool(&c.Server.SecureCookie, EnvServerSecureCookie))

	add(parseString(&c.Database.Dialect, EnvDatabaseDialect))
	add(parseString(&c.Database.Connection, EnvDatabaseConnection))

	add(parseString(&c.DefaultUser.Name, EnvDefaultUserName))
	add(parseString(&c.DefaultUser.Pass, EnvDefaultUserPass))

	add(parseInt(&c.PassStrength, EnvPassStrength))
	add(parseString(&c.UploadedImagesDir, EnvUploadedImagesDir))
	add(parseString(&c.PluginsDir, EnvPluginsDir))
	add(parseBool(&c.Registration, EnvRegistration))

	add(parseBool(&c.OIDC.Enabled, EnvOIDCEnabled))
	add(parseString(&c.OIDC.Issuer, EnvOIDCIssuer))
	add(parseString(&c.OIDC.ClientID, EnvOIDCClientID))
	add(parseString(&c.OIDC.ClientSecret, EnvOIDCClientSecret))
	add(parseString(&c.OIDC.UsernameClaim, EnvOIDCUsernameClaim))
	add(parseString(&c.OIDC.RedirectURL, EnvOIDCRedirectURL))
	add(parseBool(&c.OIDC.AutoRegister, EnvOIDCAutoRegister))
	add(parseBool(&c.OIDC.LinkByUsername, EnvOIDCLinkByUsername))
	add(parseList(&c.OIDC.Scopes, EnvOIDCScopes))

	add(parseString(&c.NoColor, EnvNoColor))

	addTrailingSlashToPaths(c)

	return c, logs
}

func addTrailingSlashToPaths(conf *Configuration) {
	if !strings.HasSuffix(conf.UploadedImagesDir, "/") && !strings.HasSuffix(conf.UploadedImagesDir, "\\") {
		conf.UploadedImagesDir += string(filepath.Separator)
	}
}
