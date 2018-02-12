package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestApplication() {
	assert.Nil(s.T(), s.db.GetApplicationByID("asdasdf"), "not existing app")

	user := &model.User{Name: "test", Pass: []byte{1}}
	s.db.CreateUser(user)
	assert.NotEqual(s.T(), 0, user.ID)

	apps := s.db.GetApplicationsByUser(user.ID)
	assert.Empty(s.T(), apps)

	app := &model.Application{UserID: user.ID, ID: "C0000000000", Name: "backupserver"}
	s.db.CreateApplication(app)

	apps = s.db.GetApplicationsByUser(user.ID)
	assert.Len(s.T(), apps, 1)
	assert.Contains(s.T(), apps, app)

	newApp := s.db.GetApplicationByID(app.ID)
	assert.Equal(s.T(), app, newApp)

	s.db.DeleteApplicationByID(app.ID)

	apps = s.db.GetApplicationsByUser(user.ID)
	assert.Empty(s.T(), apps)

	assert.Nil(s.T(), s.db.GetApplicationByID(app.ID))
}
