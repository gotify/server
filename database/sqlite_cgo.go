//go:build cgo

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func sqliteDialector(connection string) gorm.Dialector {
	return sqlite.Open(connection)
}
