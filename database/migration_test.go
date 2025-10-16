package database

import (
	"testing"

	"github.com/gotify/server/v2/model"
	"github.com/gotify/server/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMigration(t *testing.T) {
	suite.Run(t, &MigrationSuite{})
}

type MigrationSuite struct {
	suite.Suite
	tmpDir test.TmpDir
}

func (s *MigrationSuite) BeforeTest(suiteName, testName string) {
	s.tmpDir = test.NewTmpDir("gotify_migrationsuite")
	db, err := gorm.Open(sqlite.Open(s.tmpDir.Path("test_obsolete.db")), &gorm.Config{})
	assert.Nil(s.T(), err)
	defer db.DB()

	assert.Nil(s.T(), db.Migrator().CreateTable(new(model.User)))
	assert.Nil(s.T(), db.Create(&model.User{
		Name:  "test_user",
		Admin: true,
	}).Error)

	// we should not be able to create applications by now
	assert.False(s.T(), db.Migrator().HasTable(new(model.Application)))
}

func (s *MigrationSuite) AfterTest(suiteName, testName string) {
	assert.Nil(s.T(), s.tmpDir.Clean())
}

func (s *MigrationSuite) TestMigration() {
	db, err := New("sqlite3", s.tmpDir.Path("test_obsolete.db"), "admin", "admin", 6, true)
	assert.Nil(s.T(), err)
	defer db.Close()

	assert.True(s.T(), db.DB.Migrator().HasTable(new(model.Application)))

	// a user already exist, not adding a new user
	if user, err := db.GetUserByName("admin"); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), user)
	}

	// the old user should persist
	if user, err := db.GetUserByName("test_user"); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), true, user.Admin)
	}

	// we should be able to create applications
	if user, err := db.GetUserByName("test_user"); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), db.CreateApplication(&model.Application{
			Token:       "A1234",
			UserID:      user.ID,
			Description: "this is a test application",
			Name:        "test application",
		}))
	}
	if app, err := db.GetApplicationByToken("A1234"); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), "test application", app.Name)
	}
}
