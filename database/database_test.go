package database

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gotify/server/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

type DatabaseSuite struct {
	suite.Suite
	db            *GormDatabase
	tmpDir        test.TmpDir
	teardownShell string
}

func (s *DatabaseSuite) BeforeTest(suiteName, testName string) {
	time := strconv.FormatInt(time.Now().UnixNano(), 10)
	if setupShell := strings.ReplaceAll(
		os.Getenv("TEST_DB_SETUP"),
		"<time>",
		time,
	); setupShell != "" {
		require.NoError(s.T(), test.ExecShell(setupShell))
	}
	s.tmpDir = test.NewTmpDir("gotify_databasesuite")
	db, err := New(
		test.GetEnv("TEST_DB_DIALECT", "sqlite3"),
		strings.ReplaceAll(
			test.GetEnv("TEST_DB_CONNECTION", "file:<time>?mode=memory&cache=shared"),
			"<time>",
			time,
		),
		"defaultUser",
		"defaultPass",
		5,
		true)
	assert.Nil(s.T(), err)
	s.db = db
	s.teardownShell = strings.ReplaceAll(
		os.Getenv("TEST_DB_TEARDOWN"),
		"<time>",
		time,
	)
}

func (s *DatabaseSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
	if s.teardownShell != "" {
		require.NoError(s.T(), test.ExecShell(s.teardownShell))
	}
	assert.Nil(s.T(), s.tmpDir.Clean())
}

func TestInvalidDialect(t *testing.T) {
	tmpDir := test.NewTmpDir("gotify_testinvaliddialect")
	defer tmpDir.Clean()
	_, err := New("asdf", tmpDir.Path("testdb.db"), "defaultUser", "defaultPass", 5, true)
	assert.Error(t, err)
}
