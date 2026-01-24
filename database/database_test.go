package database

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gotify/server/v2/model"
	"github.com/gotify/server/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
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

func TestMigrateSortKey(t *testing.T) {
	db, err := New("sqlite3", fmt.Sprintf("file:%s?mode=memory&cache=shared", fmt.Sprint(time.Now().UnixNano())), "admin", "pw", 5, true)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.CreateApplication(&model.Application{Name: "one", Token: "one", UserID: 1})
	assert.NoError(t, err)
	err = db.CreateApplication(&model.Application{Name: "two", Token: "two", UserID: 1})
	assert.NoError(t, err)
	err = db.CreateApplication(&model.Application{Name: "three", Token: "three", UserID: 1})
	assert.NoError(t, err)
	err = db.CreateApplication(&model.Application{Name: "one-other", Token: "one-other", UserID: 2})
	assert.NoError(t, err)

	err = db.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(new(model.Application)).UpdateColumn("sort_key", nil).Error
	assert.NoError(t, err)

	err = fillMissingSortKeys(db.DB)
	assert.NoError(t, err)

	apps, err := db.GetApplicationsByUser(1)
	assert.NoError(t, err)

	assert.Len(t, apps, 3)
	assert.Equal(t, apps[0].Name, "one")
	assert.Equal(t, apps[0].SortKey, "a0")
	assert.Equal(t, apps[1].Name, "two")
	assert.Equal(t, apps[1].SortKey, "a1")
	assert.Equal(t, apps[2].Name, "three")
	assert.Equal(t, apps[2].SortKey, "a2")

	apps, err = db.GetApplicationsByUser(2)
	assert.NoError(t, err)

	assert.Len(t, apps, 1)
	assert.Equal(t, apps[0].Name, "one-other")
	assert.Equal(t, apps[0].SortKey, "a0")
}
