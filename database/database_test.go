package database

import (
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
	db, err := New("sqlite3", "testdb.db", "defaultUser", "defaultPass")
	assert.Nil(s.T(), err)
	s.db = db
}

func (s *DatabaseSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
	assert.Nil(s.T(), os.Remove("testdb.db"))
}

func TestInvalidDialect(t *testing.T) {
	_, err := New("asdf", "testdb.db", "defaultUser", "defaultPass")
	assert.NotNil(t, err)
}
