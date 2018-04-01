package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestApplication() {
	assert.Nil(s.T(), s.db.GetApplicationByToken("asdasdf"), "not existing app")
	assert.Nil(s.T(), s.db.GetApplicationByID(uint(1)), "not existing app")

	user := &model.User{Name: "test", Pass: []byte{1}}
	s.db.CreateUser(user)
	assert.NotEqual(s.T(), 0, user.ID)

	apps := s.db.GetApplicationsByUser(user.ID)
	assert.Empty(s.T(), apps)

	app := &model.Application{UserID: user.ID, Token: "C0000000000", Name: "backupserver"}
	s.db.CreateApplication(app)

	apps = s.db.GetApplicationsByUser(user.ID)
	assert.Len(s.T(), apps, 1)
	assert.Contains(s.T(), apps, app)

	newApp := s.db.GetApplicationByToken(app.Token)
	assert.Equal(s.T(), app, newApp)

	newApp = s.db.GetApplicationByID(app.ID)
	assert.Equal(s.T(), app, newApp)

	newApp.Image = "asdasd"
	s.db.UpdateApplication(newApp)

	newApp = s.db.GetApplicationByID(app.ID)
	assert.Equal(s.T(), "asdasd", newApp.Image)

	s.db.DeleteApplicationByID(app.ID)

	apps = s.db.GetApplicationsByUser(user.ID)
	assert.Empty(s.T(), apps)

	assert.Nil(s.T(), s.db.GetApplicationByID(app.ID))
}

func (s *DatabaseSuite) TestDeleteAppDeletesMessages() {
	s.db.CreateApplication(&model.Application{ID: 55, Token: "token"})
	s.db.CreateApplication(&model.Application{ID: 66, Token: "token2"})
	s.db.CreateMessage(&model.Message{ID: 12, ApplicationID: 55})
	s.db.CreateMessage(&model.Message{ID: 13, ApplicationID: 66})
	s.db.CreateMessage(&model.Message{ID: 14, ApplicationID: 55})
	s.db.CreateMessage(&model.Message{ID: 15, ApplicationID: 55})

	s.db.DeleteApplicationByID(55)

	assert.Nil(s.T(), s.db.GetMessageByID(12))
	assert.NotNil(s.T(), s.db.GetMessageByID(13))
	assert.Nil(s.T(), s.db.GetMessageByID(14))
	assert.Nil(s.T(), s.db.GetMessageByID(15))
	assert.Empty(s.T(), s.db.GetMessagesByApplication(55))
	assert.NotEmpty(s.T(), s.db.GetMessagesByApplication(66))
}
