package database

import (
	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestClient() {
	assert.Nil(s.T(), s.db.GetClientByID("asdasdf"), "not existing client")

	user := &model.User{Name: "test", Pass: []byte{1}}
	s.db.CreateUser(user)
	assert.NotEqual(s.T(), 0, user.ID)

	clients := s.db.GetClientsByUser(user.ID)
	assert.Empty(s.T(), clients)

	client := &model.Client{UserID: user.ID, ID: "C0000000000", Name: "android"}
	s.db.CreateClient(client)

	clients = s.db.GetClientsByUser(user.ID)
	assert.Len(s.T(), clients, 1)
	assert.Contains(s.T(), clients, client)

	newClient := s.db.GetClientByID(client.ID)
	assert.Equal(s.T(), client, newClient)

	s.db.DeleteClientByID(client.ID)

	clients = s.db.GetClientsByUser(user.ID)
	assert.Empty(s.T(), clients)

	assert.Nil(s.T(), s.db.GetClientByID(client.ID))
}
