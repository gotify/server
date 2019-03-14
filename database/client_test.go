package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestClient() {
	assert.Nil(s.T(), s.db.GetClientByID(1), "not existing client")
	assert.Nil(s.T(), s.db.GetClientByToken("asdasd"), "not existing client")

	user := &model.User{Name: "test", Pass: []byte{1}}
	s.db.CreateUser(user)
	assert.NotEqual(s.T(), 0, user.ID)

	clients := s.db.GetClientsByUser(user.ID)
	assert.Empty(s.T(), clients)

	client := &model.Client{UserID: user.ID, Token: "C0000000000", Name: "android"}
	s.db.CreateClient(client)

	clients = s.db.GetClientsByUser(user.ID)
	assert.Len(s.T(), clients, 1)
	assert.Contains(s.T(), clients, client)

	newClient := s.db.GetClientByID(client.ID)
	assert.Equal(s.T(), client, newClient)

	newClient = s.db.GetClientByToken(client.Token)
	assert.Equal(s.T(), client, newClient)

	updateClient := &model.Client{ID: client.ID, UserID: user.ID, Token: "C0000000000", Name: "new_name"}
	s.db.UpdateClient(updateClient)
	updatedClient := s.db.GetClientByID(client.ID)
	assert.Equal(s.T(), updateClient, updatedClient)

	s.db.DeleteClientByID(client.ID)

	clients = s.db.GetClientsByUser(user.ID)
	assert.Empty(s.T(), clients)

	assert.Nil(s.T(), s.db.GetClientByID(client.ID))
}
