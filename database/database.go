package database

import (
	"database/sql"
	"errors"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/gotify/server/v2/auth/password"
	"github.com/gotify/server/v2/fracdex"
	"github.com/gotify/server/v2/model"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// gormLogWriter routes gorm logger output through zerolog.
type gormLogWriter struct{}

func (gormLogWriter) Printf(format string, args ...interface{}) {
	log.Warn().Str("component", "gorm").Msgf(format, args...)
}

var mkdirAll = os.MkdirAll

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, defaultUser, defaultPass string, strength int, createDefaultUserIfNotExist bool, now func() time.Time) (*GormDatabase, error) {
	createDirectoryIfSqlite(dialect, connection)

	dbLogger := logger.New(gormLogWriter{}, logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  isatty.IsTerminal(os.Stderr.Fd()),
	})
	gormConfig := &gorm.Config{
		Logger:                                   dbLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		TranslateError:                           true,
		NowFunc:                                  now,
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

	if err := db.Transaction(fillMissingSortKeys, &sql.TxOptions{Isolation: sql.LevelSerializable}); err != nil {
		return nil, err
	}

	if err := db.Transaction(func(tx *gorm.DB) error { return fillMissingCreatedAt(tx, now()) }, &sql.TxOptions{Isolation: sql.LevelSerializable}); err != nil {
		return nil, err
	}

	return &GormDatabase{DB: db}, nil
}

func fillMissingCreatedAt(db *gorm.DB, now time.Time) error {
	models := []any{
		new(model.User),
		new(model.Application),
		new(model.Client),
		new(model.PluginConf),
	}
	for _, m := range models {
		if err := db.Model(m).Where("created_at IS NULL").UpdateColumn("created_at", now).Error; err != nil {
			return err
		}
	}
	return nil
}

func fillMissingSortKeys(db *gorm.DB) error {
	missingSort := int64(0)
	if err := db.Model(new(model.Application)).Where("sort_key IS NULL OR sort_key = ''").Count(&missingSort).Error; err != nil {
		return err
	}

	if missingSort == 0 {
		return nil
	}

	var apps []*model.Application
	if err := db.Order("user_id, sort_key, id ASC").Find(&apps).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	log.Info().Int("count", len(apps)).Msg("Migrating application sort keys")

	sortKey := ""
	currentUser := uint(math.MaxUint)
	var err error
	for _, app := range apps {
		if currentUser != app.UserID {
			sortKey = ""
			currentUser = app.UserID
		}
		sortKey, err = fracdex.KeyBetween(sortKey, "")
		if err != nil {
			return err
		}
		app.SortKey = sortKey
	}
	return db.Save(apps).Error
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
