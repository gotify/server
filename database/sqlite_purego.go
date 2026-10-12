//go:build !cgo

package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func sqliteDialector(connection string) gorm.Dialector {
	return sqlite.Open(connection)
}
