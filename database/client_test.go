package database

import (
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestClient() {
	if client, err := s.db.GetClientByID(1); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), client, "not existing client")
	}
	if client, err := s.db.GetClientByToken("asdasd"); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), client, "not existing client")
	}

	user := &model.User{Name: "test", Pass: []byte{1}}
	s.db.CreateUser(user)
	assert.NotEqual(s.T(), 0, user.ID)

	if clients, err := s.db.GetClientsByUser(user.ID); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), clients)
	}

	client := &model.Client{UserID: user.ID, Token: "C0000000000", Name: "android"}
	assert.NoError(s.T(), s.db.CreateClient(client))

	if clients, err := s.db.GetClientsByUser(user.ID); assert.NoError(s.T(), err) {
		assert.Len(s.T(), clients, 1)
		assert.Contains(s.T(), clients, client)
	}

	newClient, err := s.db.GetClientByID(client.ID)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), client, newClient)
	}

	if newClient, err := s.db.GetClientByToken(client.Token); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), client, newClient)
	}

	updateClient := &model.Client{ID: client.ID, UserID: user.ID, Token: "C0000000000", Name: "new_name"}
	s.db.UpdateClient(updateClient)
	if updatedClient, err := s.db.GetClientByID(client.ID); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), updateClient, updatedClient)
	}

	s.db.DeleteClientByID(client.ID)

	if clients, err := s.db.GetClientsByUser(user.ID); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), clients)
	}

	if client, err := s.db.GetClientByID(client.ID); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), client)
	}
}
