package database

import (
	"errors"
	"os"
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
	db, err := New("sqlite3", s.tmpDir.Path("testdb.db"), "defaultUser", "defaultPass", 5, true)
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

func TestCreateSqliteFolder(t *testing.T) {
	tmpDir := test.NewTmpDir("gotify_testcreatesqlitefolder")
	defer tmpDir.Clean()

	db, err := New("sqlite3", tmpDir.Path("somepath/testdb.db"), "defaultUser", "defaultPass", 5, true)
	assert.Nil(t, err)
	assert.DirExists(t, tmpDir.Path("somepath"))
	db.Close()
}

func TestWithAlreadyExistingSqliteFolder(t *testing.T) {
	tmpDir := test.NewTmpDir("gotify_testwithexistingfolder")
	defer tmpDir.Clean()

	db, err := New("sqlite3", tmpDir.Path("somepath/testdb.db"), "defaultUser", "defaultPass", 5, true)
	assert.Nil(t, err)
	assert.DirExists(t, tmpDir.Path("somepath"))
	db.Close()
}

func TestPanicsOnMkdirError(t *testing.T) {
	tmpDir := test.NewTmpDir("gotify_testpanicsonmkdirerror")
	defer tmpDir.Clean()
	mkdirAll = func(path string, perm os.FileMode) error {
		return errors.New("ERROR")
	}
	assert.Panics(t, func() {
		New("sqlite3", tmpDir.Path("somepath/test.db"), "defaultUser", "defaultPass", 5, true)
	})
}
