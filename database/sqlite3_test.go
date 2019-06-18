package database

import (
	"errors"
	"os"
	"testing"

	"github.com/gotify/server/test"
	"github.com/stretchr/testify/assert"
)

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
