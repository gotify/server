package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestApplication() {

	if app, err := s.db.GetApplicationByToken("asdasdf"); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), app, "not existing app")
	}

	if app, err := s.db.GetApplicationByID(uint(1)); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), app, "not existing app")
	}

	user := &model.User{Name: "test", Pass: []byte{1}}
	s.db.CreateUser(user)
	assert.NotEqual(s.T(), 0, user.ID)

	if apps, err := s.db.GetApplicationsByUser(user.ID); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), apps)
	}

	app := &model.Application{UserID: user.ID, Token: "C0000000000", Name: "backupserver"}
	s.db.CreateApplication(app)

	if apps, err := s.db.GetApplicationsByUser(user.ID); assert.NoError(s.T(), err) {
		assert.Len(s.T(), apps, 1)
		assert.Contains(s.T(), apps, app)
	}

	newApp, err := s.db.GetApplicationByToken(app.Token)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), app, newApp)
	}

	newApp, err = s.db.GetApplicationByID(app.ID)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), app, newApp)
	}

	newApp.Image = "asdasd"
	assert.NoError(s.T(), s.db.UpdateApplication(newApp))

	newApp, err = s.db.GetApplicationByID(app.ID)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), "asdasd", newApp.Image)
	}

	assert.NoError(s.T(), s.db.DeleteApplicationByID(app.ID))

	if apps, err := s.db.GetApplicationsByUser(user.ID); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), apps)
	}

	if app, err := s.db.GetApplicationByID(app.ID); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), app)
	}
}

func (s *DatabaseSuite) TestDeleteAppDeletesMessages() {
	assert.NoError(s.T(), s.db.CreateApplication(&model.Application{ID: 55, Token: "token"}))
	assert.NoError(s.T(), s.db.CreateApplication(&model.Application{ID: 66, Token: "token2"}))
	assert.NoError(s.T(), s.db.CreateMessage(&model.Message{ID: 12, ApplicationID: 55}))
	assert.NoError(s.T(), s.db.CreateMessage(&model.Message{ID: 13, ApplicationID: 66}))
	assert.NoError(s.T(), s.db.CreateMessage(&model.Message{ID: 14, ApplicationID: 55}))
	assert.NoError(s.T(), s.db.CreateMessage(&model.Message{ID: 15, ApplicationID: 55}))

	assert.NoError(s.T(), s.db.DeleteApplicationByID(55))

	if msg, err := s.db.GetMessageByID(12); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), msg)
	}
	if msg, err := s.db.GetMessageByID(13); assert.NoError(s.T(), err) {
		assert.NotNil(s.T(), msg)
	}
	if msg, err := s.db.GetMessageByID(14); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), msg)
	}
	if msg, err := s.db.GetMessageByID(15); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), msg)
	}

	if msgs, err := s.db.GetMessagesByApplication(55); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), msgs)
	}
	if msgs, err := s.db.GetMessagesByApplication(66); assert.NoError(s.T(), err) {
		assert.NotEmpty(s.T(), msgs)
	}
}
