package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // enable the mysql dialect
	_ "github.com/jinzhu/gorm/dialects/postgres" // enable the postgres dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // enable the sqlite3 dialect
	"github.com/jmattheis/memo/auth"
	"github.com/jmattheis/memo/model"
)

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, defaultUser, defaultPass string) (*GormDatabase, error) {
	db, err := gorm.Open(dialect, connection)
	if err != nil {
		return nil, err
	}
	if !db.HasTable(new(model.User)) && !db.HasTable(new(model.Message)) &&
		!db.HasTable(new(model.Client)) && !db.HasTable(new(model.Application)) {
		db.AutoMigrate(new(model.User), new(model.Application), new(model.Message), new(model.Client))
		db.Create(&model.User{Name: defaultUser, Pass: auth.CreatePassword(defaultPass), Admin: true})
	}

	return &GormDatabase{DB: db}, nil
}

// GormDatabase is a wrapper for the gorm framework.
type GormDatabase struct {
	DB *gorm.DB
}

// Close closes the gorm database connection.
func (d *GormDatabase) Close() {
	d.DB.Close()
}
