package database

import (
	"testing"

	"github.com/gotify/server/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

type DatabaseSuite struct {
	suite.Suite
	db     *GormDatabase
	tmpDir test.TmpDir
}

func (s *DatabaseSuite) BeforeTest(suiteName, testName string) {
	s.tmpDir = test.NewTmpDir("gotify_databasesuite")
	db, err := New(
		test.GetEnv("TEST_DB_DIALECT", "sqlite3"),
		test.GetEnv("TEST_DB_CONNECTION", s.tmpDir.Path("testdb.db")),
		"defaultUser",
		"defaultPass",
		5,
		true)
	assert.Nil(s.T(), err)
	s.db = db
}

func (s *DatabaseSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
	assert.Nil(s.T(), s.tmpDir.Clean())
}

func TestInvalidDialect(t *testing.T) {
	tmpDir := test.NewTmpDir("gotify_testinvaliddialect")
	defer tmpDir.Clean()
	_, err := New("asdf", tmpDir.Path("testdb.db"), "defaultUser", "defaultPass", 5, true)
	assert.Error(t, err)
}
