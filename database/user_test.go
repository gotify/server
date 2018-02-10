package database

import (
	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestUser() {
	assert.Nil(s.T(), s.db.GetUserByID(55), "not existing user")
	assert.Nil(s.T(), s.db.GetUserByName("nicories"), "not existing user")

	jmattheis := s.db.GetUserByID(1)
	assert.NotNil(s.T(), jmattheis, "on bootup the first user should be automatically created")

	users := s.db.GetUsers()
	assert.Len(s.T(), users, 1)
	assert.Contains(s.T(), users, jmattheis)

	nicories := &model.User{Name: "nicories", Pass: []byte{1, 2, 3, 4}, Admin: false}
	s.db.CreateUser(nicories)
	assert.NotEqual(s.T(), 0, nicories.ID, "on create user a new id should be assigned")

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

	s.db.DeleteUserByID(tom.ID)
	users = s.db.GetUsers()
	assert.Len(s.T(), users, 1)
	assert.Contains(s.T(), users, jmattheis)

	s.db.DeleteUserByID(jmattheis.ID)
	users = s.db.GetUsers()
	assert.Empty(s.T(), users)
}
