package database

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

type DatabaseSuite struct {
	suite.Suite
	db *GormDatabase
}

func (s *DatabaseSuite) BeforeTest(suiteName, testName string) {
	db, err := New("sqlite3", "testdb.db", "defaultUser", "defaultPass", 5, true)
	assert.Nil(s.T(), err)
	s.db = db
}

func (s *DatabaseSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
	assert.Nil(s.T(), os.Remove("testdb.db"))
}

func TestInvalidDialect(t *testing.T) {
	_, err := New("asdf", "testdb.db", "defaultUser", "defaultPass", 5, true)
	assert.NotNil(t, err)
}

func TestCreateSqliteFolder(t *testing.T) {
	// ensure path not exists
	os.RemoveAll("somepath")

	db, err := New("sqlite3", "somepath/testdb.db", "defaultUser", "defaultPass", 5, true)
	assert.Nil(t, err)
	assert.DirExists(t, "somepath")
	db.Close()

	assert.Nil(t, os.RemoveAll("somepath"))
}

func TestWithAlreadyExistingSqliteFolder(t *testing.T) {
	// ensure path not exists
	os.RemoveAll("somepath")
	os.MkdirAll("somepath", 0777)

	db, err := New("sqlite3", "somepath/testdb.db", "defaultUser", "defaultPass", 5, true)
	assert.Nil(t, err)
	assert.DirExists(t, "somepath")
	db.Close()

	assert.Nil(t, os.RemoveAll("somepath"))
}

func TestPanicsOnMkdirError(t *testing.T) {
	os.RemoveAll("somepath")
	mkdirAll = func(path string, perm os.FileMode) error {
		return errors.New("ERROR")
	}
	assert.Panics(t, func() {
		New("sqlite3", "somepath/test.db", "defaultUser", "defaultPass", 5, true)
	})
}
