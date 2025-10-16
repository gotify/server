package database

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gotify/server/v2/auth/password"
	"github.com/gotify/server/v2/mode"
	"github.com/gotify/server/v2/model"
	"github.com/mattn/go-isatty"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var mkdirAll = os.MkdirAll

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, defaultUser, defaultPass string, strength int, createDefaultUserIfNotExist bool) (*GormDatabase, error) {
	createDirectoryIfSqlite(dialect, connection)

	logLevel := logger.Info
	if mode.Get() == mode.Prod {
		logLevel = logger.Warn
	}

	dbLogger := logger.New(log.New(os.Stderr, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: true,
		Colorful:                  isatty.IsTerminal(os.Stderr.Fd()),
	})
	gormConfig := &gorm.Config{
		Logger: dbLogger,
	}

	var db *gorm.DB
	err := errors.New("unsupported dialect: " + dialect)

	switch dialect {
	case "mysql":
		db, err = gorm.Open(mysql.Open(connection), gormConfig)
	case "postgres":
		db, err = gorm.Open(postgres.Open(connection), gormConfig)
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(connection), gormConfig)
	}

	if err != nil {
		return nil, err
	}

	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}

	// We normally don't need that much connections, so we limit them. F.ex. mysql complains about
	// "too many connections", while load testing Gotify.
	sqldb.SetMaxOpenConns(10)

	if dialect == "sqlite3" {
		// We use the database connection inside the handlers from the http
		// framework, therefore concurrent access occurs. Sqlite cannot handle
		// concurrent writes, so we limit sqlite to one connection.
		// see https://github.com/mattn/go-sqlite3/issues/274
		sqldb.SetMaxOpenConns(1)
	}

	if dialect == "mysql" {
		// Mysql has a setting called wait_timeout, which defines the duration
		// after which a connection may not be used anymore.
		// The default for this setting on mariadb is 10 minutes.
		// See https://github.com/docker-library/mariadb/issues/113
		sqldb.SetConnMaxLifetime(9 * time.Minute)
	}

	if err := db.AutoMigrate(new(model.User), new(model.Application), new(model.Message), new(model.Client), new(model.PluginConf)); err != nil {
		return nil, err
	}

	userCount := int64(0)
	db.Find(new(model.User)).Count(&userCount)
	if createDefaultUserIfNotExist && userCount == 0 {
		db.Create(&model.User{Name: defaultUser, Pass: password.CreatePassword(defaultPass, strength), Admin: true})
	}

	return &GormDatabase{DB: db}, nil
}

func createDirectoryIfSqlite(dialect, connection string) {
	if dialect == "sqlite3" {
		if _, err := os.Stat(filepath.Dir(connection)); os.IsNotExist(err) {
			if err := mkdirAll(filepath.Dir(connection), 0o777); err != nil {
				panic(err)
			}
		}
	}
}

// GormDatabase is a wrapper for the gorm framework.
type GormDatabase struct {
	DB *gorm.DB
}

// Close closes the gorm database connection.
func (d *GormDatabase) Close() {
	sqldb, err := d.DB.DB()
	if err != nil {
		return
	}
	sqldb.Close()
}
