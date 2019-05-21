package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestUser() {
	if user, err := s.db.GetUserByID(55); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), user, "not existing user")
	}
	if user, err := s.db.GetUserByName("nicories"); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), user, "not existing user")
	}
	jmattheis, err := s.db.GetUserByID(1)
	if assert.NoError(s.T(), err) {
		assert.NotNil(s.T(), jmattheis, "on bootup the first user should be automatically created")
	}
	if adminCount, err := s.db.CountUser("admin = ?", true); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), 1, adminCount, 1, "there is initially one admin")
	}

	if users, err := s.db.GetUsers(); assert.NoError(s.T(), err) {
		assert.Len(s.T(), users, 1)
		assert.Contains(s.T(), users, jmattheis)
	}

	nicories := &model.User{Name: "nicories", Pass: []byte{1, 2, 3, 4}, Admin: false}
	s.db.CreateUser(nicories)
	assert.NotEqual(s.T(), 0, nicories.ID, "on create user a new id should be assigned")
	if userCount, err := s.db.CountUser(); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), 2, userCount, "two users should exist")
	}

	if user, err := s.db.GetUserByName("nicories"); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), nicories, user)
	}

	if users, err := s.db.GetUsers(); assert.NoError(s.T(), err) {
		assert.Len(s.T(), users, 2)
		assert.Contains(s.T(), users, jmattheis)
		assert.Contains(s.T(), users, nicories)
	}

	nicories.Name = "tom"
	nicories.Pass = []byte{12}
	nicories.Admin = true
	assert.NoError(s.T(), s.db.UpdateUser(nicories))

	tom, err := s.db.GetUserByID(nicories.ID)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), &model.User{ID: nicories.ID, Name: "tom", Pass: []byte{12}, Admin: true}, tom)
	}
	if users, err := s.db.GetUsers(); assert.NoError(s.T(), err) {
		assert.Len(s.T(), users, 2)
	}
	if adminCount, err := s.db.CountUser(&model.User{Admin: true}); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), 2, adminCount, "two admins exist")
	}

	assert.NoError(s.T(), s.db.DeleteUserByID(tom.ID))
	if users, err := s.db.GetUsers(); assert.NoError(s.T(), err) {
		assert.Len(s.T(), users, 1)
		assert.Contains(s.T(), users, jmattheis)
	}

	s.db.DeleteUserByID(jmattheis.ID)
	if users, err := s.db.GetUsers(); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), users)
	}
}

func (s *DatabaseSuite) TestUserPlugins() {
	assert.NoError(s.T(), s.db.CreateUser(&model.User{Name: "geek", ID: 16}))
	if geekUser, err := s.db.GetUserByName("geek"); assert.NoError(s.T(), err) {
		s.db.CreatePluginConf(&model.PluginConf{
			UserID:     geekUser.ID,
			ModulePath: "github.com/gotify/example-plugin",
			Token:      "P1234",
			Enabled:    true,
		})
		s.db.CreatePluginConf(&model.PluginConf{
			UserID:     geekUser.ID,
			ModulePath: "github.com/gotify/example-plugin/v2",
			Token:      "P5678",
			Enabled:    true,
		})
	}

	if geekUser, err := s.db.GetUserByName("geek"); assert.NoError(s.T(), err) {
		if pluginConfs, err := s.db.GetPluginConfByUser(geekUser.ID); assert.NoError(s.T(), err) {
			assert.Len(s.T(), pluginConfs, 2)
		}
	}
	if pluginConf, err := s.db.GetPluginConfByToken("P1234"); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), "github.com/gotify/example-plugin", pluginConf.ModulePath)
	}

}

func (s *DatabaseSuite) TestDeleteUserDeletesApplicationsAndClientsAndPluginConfs() {
	assert.NoError(s.T(), s.db.CreateUser(&model.User{Name: "nicories", ID: 10}))
	assert.NoError(s.T(), s.db.CreateApplication(&model.Application{ID: 100, Token: "apptoken", UserID: 10}))
	assert.NoError(s.T(), s.db.CreateMessage(&model.Message{ID: 1000, ApplicationID: 100}))
	assert.NoError(s.T(), s.db.CreateClient(&model.Client{ID: 10000, Token: "clienttoken", UserID: 10}))
	assert.NoError(s.T(), s.db.CreatePluginConf(&model.PluginConf{ID: 1000, Token: "plugintoken", UserID: 10}))

	assert.NoError(s.T(), s.db.CreateUser(&model.User{Name: "nicories2", ID: 20}))
	assert.NoError(s.T(), s.db.CreateApplication(&model.Application{ID: 200, Token: "apptoken2", UserID: 20}))
	assert.NoError(s.T(), s.db.CreateMessage(&model.Message{ID: 2000, ApplicationID: 200}))
	assert.NoError(s.T(), s.db.CreateClient(&model.Client{ID: 20000, Token: "clienttoken2", UserID: 20}))
	assert.NoError(s.T(), s.db.CreatePluginConf(&model.PluginConf{ID: 2000, Token: "plugintoken2", UserID: 20}))

	assert.NoError(s.T(), s.db.DeleteUserByID(10))

	if app, err := s.db.GetApplicationByToken("apptoken"); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), app)
	}
	if client, err := s.db.GetClientByToken("clienttoken"); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), client)
	}
	if clients, err := s.db.GetClientsByUser(10); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), clients)
	}
	if apps, err := s.db.GetApplicationsByUser(10); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), apps)
	}
	if msgs, err := s.db.GetMessagesByApplication(100); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), msgs)
	}
	if msgs, err := s.db.GetMessagesByUser(10); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), msgs)
	}
	if pluginConfs, err := s.db.GetPluginConfByUser(10); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), pluginConfs)
	}
	if msg, err := s.db.GetMessageByID(1000); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), msg)
	}

	if apps, err := s.db.GetApplicationByToken("apptoken2"); assert.NoError(s.T(), err) {
		assert.NotNil(s.T(), apps)
	}
	if clients, err := s.db.GetClientByToken("clienttoken2"); assert.NoError(s.T(), err) {
		assert.NotNil(s.T(), clients)
	}
	if clients, err := s.db.GetClientsByUser(20); assert.NoError(s.T(), err) {
		assert.NotEmpty(s.T(), clients)
	}
	if apps, err := s.db.GetApplicationsByUser(20); assert.NoError(s.T(), err) {
		assert.NotEmpty(s.T(), apps)
	}
	if pluginConf, err := s.db.GetPluginConfByUser(20); assert.NoError(s.T(), err) {
		assert.NotEmpty(s.T(), pluginConf)
	}
	if msgs, err := s.db.GetMessagesByApplication(200); assert.NoError(s.T(), err) {
		assert.NotEmpty(s.T(), msgs)
	}
	if msgs, err := s.db.GetMessagesByUser(20); assert.NoError(s.T(), err) {
		assert.NotEmpty(s.T(), msgs)
	}
	if msg, err := s.db.GetMessageByID(2000); assert.NoError(s.T(), err) {
		assert.NotNil(s.T(), msg)
	}
}
