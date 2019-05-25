package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DatabaseSuite) TestUser() {
	user, err := s.db.GetUserByID(55)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), user, "not existing user")

	user, err = s.db.GetUserByName("nicories")
	require.NoError(s.T(), err)
	assert.Nil(s.T(), user, "not existing user")

	jmattheis, err := s.db.GetUserByID(1)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), jmattheis, "on bootup the first user should be automatically created")

	adminCount, err := s.db.CountUser("admin = ?", true)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, adminCount, 1, "there is initially one admin")

	users, err := s.db.GetUsers()
	require.NoError(s.T(), err)
	assert.Len(s.T(), users, 1)
	assert.Contains(s.T(), users, jmattheis)

	nicories := &model.User{Name: "nicories", Pass: []byte{1, 2, 3, 4}, Admin: false}
	s.db.CreateUser(nicories)
	assert.NotEqual(s.T(), 0, nicories.ID, "on create user a new id should be assigned")
	userCount, err := s.db.CountUser()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, userCount, "two users should exist")

	user, err = s.db.GetUserByName("nicories")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), nicories, user)

	users, err = s.db.GetUsers()
	require.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
	assert.Contains(s.T(), users, jmattheis)
	assert.Contains(s.T(), users, nicories)

	nicories.Name = "tom"
	nicories.Pass = []byte{12}
	nicories.Admin = true
	require.NoError(s.T(), s.db.UpdateUser(nicories))

	tom, err := s.db.GetUserByID(nicories.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), &model.User{ID: nicories.ID, Name: "tom", Pass: []byte{12}, Admin: true}, tom)

	users, err = s.db.GetUsers()
	require.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)

	adminCount, err = s.db.CountUser(&model.User{Admin: true})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, adminCount, "two admins exist")

	require.NoError(s.T(), s.db.DeleteUserByID(tom.ID))
	users, err = s.db.GetUsers()
	require.NoError(s.T(), err)
	assert.Len(s.T(), users, 1)
	assert.Contains(s.T(), users, jmattheis)

	s.db.DeleteUserByID(jmattheis.ID)
	users, err = s.db.GetUsers()
	require.NoError(s.T(), err)
	assert.Empty(s.T(), users)

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
	require.NoError(s.T(), s.db.CreateUser(&model.User{Name: "nicories", ID: 10}))
	require.NoError(s.T(), s.db.CreateApplication(&model.Application{ID: 100, Token: "apptoken", UserID: 10}))
	require.NoError(s.T(), s.db.CreateMessage(&model.Message{ID: 1000, ApplicationID: 100}))
	require.NoError(s.T(), s.db.CreateClient(&model.Client{ID: 10000, Token: "clienttoken", UserID: 10}))
	require.NoError(s.T(), s.db.CreatePluginConf(&model.PluginConf{ID: 1000, Token: "plugintoken", UserID: 10}))

	require.NoError(s.T(), s.db.CreateUser(&model.User{Name: "nicories2", ID: 20}))
	require.NoError(s.T(), s.db.CreateApplication(&model.Application{ID: 200, Token: "apptoken2", UserID: 20}))
	require.NoError(s.T(), s.db.CreateMessage(&model.Message{ID: 2000, ApplicationID: 200}))
	require.NoError(s.T(), s.db.CreateClient(&model.Client{ID: 20000, Token: "clienttoken2", UserID: 20}))
	require.NoError(s.T(), s.db.CreatePluginConf(&model.PluginConf{ID: 2000, Token: "plugintoken2", UserID: 20}))

	require.NoError(s.T(), s.db.DeleteUserByID(10))

	app, err := s.db.GetApplicationByToken("apptoken")
	require.NoError(s.T(), err)
	assert.Nil(s.T(), app)

	client, err := s.db.GetClientByToken("clienttoken")
	require.NoError(s.T(), err)
	assert.Nil(s.T(), client)

	clients, err := s.db.GetClientsByUser(10)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), clients)

	apps, err := s.db.GetApplicationsByUser(10)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), apps)

	msgs, err := s.db.GetMessagesByApplication(100)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), msgs)

	msgs, err = s.db.GetMessagesByUser(10)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), msgs)

	pluginConfs, err := s.db.GetPluginConfByUser(10)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), pluginConfs)

	msg, err := s.db.GetMessageByID(1000)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), msg)

	app, err = s.db.GetApplicationByToken("apptoken2")
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), app)

	client, err = s.db.GetClientByToken("clienttoken2")
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), client)

	clients, err = s.db.GetClientsByUser(20)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), clients)

	apps, err = s.db.GetApplicationsByUser(20)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), apps)

	pluginConf, err := s.db.GetPluginConfByUser(20)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), pluginConf)

	msgs, err = s.db.GetMessagesByApplication(200)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), msgs)

	msgs, err = s.db.GetMessagesByUser(20)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), msgs)

	msg, err = s.db.GetMessageByID(2000)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), msg)

}
