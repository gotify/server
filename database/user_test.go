package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestUser() {
	assert.Nil(s.T(), s.db.GetUserByID(55), "not existing user")
	assert.Nil(s.T(), s.db.GetUserByName("nicories"), "not existing user")

	jmattheis := s.db.GetUserByID(1)
	assert.NotNil(s.T(), jmattheis, "on bootup the first user should be automatically created")
	assert.Equal(s.T(), 1, s.db.CountUser("admin = ?", true), 1, "there is initially one admin")

	users := s.db.GetUsers()
	assert.Len(s.T(), users, 1)
	assert.Contains(s.T(), users, jmattheis)

	nicories := &model.User{Name: "nicories", Pass: []byte{1, 2, 3, 4}, Admin: false}
	s.db.CreateUser(nicories)
	assert.NotEqual(s.T(), 0, nicories.ID, "on create user a new id should be assigned")
	assert.Equal(s.T(), 2, s.db.CountUser(), "two users should exist")

	assert.Equal(s.T(), nicories, s.db.GetUserByName("nicories"))

	users = s.db.GetUsers()
	assert.Len(s.T(), users, 2)
	assert.Contains(s.T(), users, jmattheis)
	assert.Contains(s.T(), users, nicories)

	nicories.Name = "tom"
	nicories.Pass = []byte{12}
	nicories.Admin = true
	s.db.UpdateUser(nicories)
	tom := s.db.GetUserByID(nicories.ID)
	assert.Equal(s.T(), &model.User{ID: nicories.ID, Name: "tom", Pass: []byte{12}, Admin: true}, tom)
	users = s.db.GetUsers()
	assert.Len(s.T(), users, 2)
	assert.Equal(s.T(), 2, s.db.CountUser(&model.User{Admin: true}), "two admins exist")

	s.db.DeleteUserByID(tom.ID)
	users = s.db.GetUsers()
	assert.Len(s.T(), users, 1)
	assert.Contains(s.T(), users, jmattheis)

	s.db.DeleteUserByID(jmattheis.ID)
	users = s.db.GetUsers()
	assert.Empty(s.T(), users)
}

func (s *DatabaseSuite) TestUserPlugins() {
	s.db.CreateUser(&model.User{Name: "geek", ID: 16})
	s.db.CreatePluginConf(&model.PluginConf{
		UserID:     s.db.GetUserByName("geek").ID,
		ModulePath: "github.com/gotify/example-plugin",
		Token:      "P1234",
		Enabled:    true,
	})
	s.db.CreatePluginConf(&model.PluginConf{
		UserID:     s.db.GetUserByName("geek").ID,
		ModulePath: "github.com/gotify/example-plugin/v2",
		Token:      "P5678",
		Enabled:    true,
	})

	assert.Len(s.T(), s.db.GetPluginConfByUser(s.db.GetUserByName("geek").ID), 2)
	assert.Equal(s.T(), "github.com/gotify/example-plugin", s.db.GetPluginConfByToken("P1234").ModulePath)

}

func (s *DatabaseSuite) TestDeleteUserDeletesApplicationsAndClientsAndPluginConfs() {
	s.db.CreateUser(&model.User{Name: "nicories", ID: 10})
	s.db.CreateApplication(&model.Application{ID: 100, Token: "apptoken", UserID: 10})
	s.db.CreateMessage(&model.Message{ID: 1000, ApplicationID: 100})
	s.db.CreateClient(&model.Client{ID: 10000, Token: "clienttoken", UserID: 10})
	s.db.CreatePluginConf(&model.PluginConf{ID: 1000, Token: "plugintoken", UserID: 10})

	s.db.CreateUser(&model.User{Name: "nicories2", ID: 20})
	s.db.CreateApplication(&model.Application{ID: 200, Token: "apptoken2", UserID: 20})
	s.db.CreateMessage(&model.Message{ID: 2000, ApplicationID: 200})
	s.db.CreateClient(&model.Client{ID: 20000, Token: "clienttoken2", UserID: 20})
	s.db.CreatePluginConf(&model.PluginConf{ID: 2000, Token: "plugintoken2", UserID: 20})

	s.db.DeleteUserByID(10)

	assert.Nil(s.T(), s.db.GetApplicationByToken("apptoken"))
	assert.Nil(s.T(), s.db.GetClientByToken("clienttoken"))
	assert.Empty(s.T(), s.db.GetClientsByUser(10))
	assert.Empty(s.T(), s.db.GetApplicationsByUser(10))
	assert.Empty(s.T(), s.db.GetMessagesByApplication(100))
	assert.Empty(s.T(), s.db.GetMessagesByUser(10))
	assert.Empty(s.T(), s.db.GetPluginConfByUser(10))
	assert.Nil(s.T(), s.db.GetMessageByID(1000))

	assert.NotNil(s.T(), s.db.GetApplicationByToken("apptoken2"))
	assert.NotNil(s.T(), s.db.GetClientByToken("clienttoken2"))
	assert.NotEmpty(s.T(), s.db.GetClientsByUser(20))
	assert.NotEmpty(s.T(), s.db.GetApplicationsByUser(20))
	assert.NotEmpty(s.T(), s.db.GetPluginConfByUser(20))
	assert.NotEmpty(s.T(), s.db.GetMessagesByApplication(200))
	assert.NotEmpty(s.T(), s.db.GetMessagesByUser(20))
	assert.NotNil(s.T(), s.db.GetMessageByID(2000))
}
