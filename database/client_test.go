package database

import (
	"time"

	"github.com/gotify/server/v2/model"
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
	assert.Nil(s.T(), client.ExpiresAt, "client without inactivity seconds has no expires_at")

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

	lastUsed := time.Now().Add(-time.Hour)
	s.db.UpdateClientTokensLastUsedAndExpiresAt([]string{client.Token}, &lastUsed)
	newClient, err = s.db.GetClientByID(client.ID)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), lastUsed.Unix(), newClient.LastUsed.Unix())
	}

	s.db.DeleteClientByID(client.ID)

	if clients, err := s.db.GetClientsByUser(user.ID); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), clients)
	}

	if client, err := s.db.GetClientByID(client.ID); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), client)
	}
}

func (s *DatabaseSuite) TestCleanupExpiredClients() {
	user := &model.User{Name: "expiry", Pass: []byte{1}}
	s.db.CreateUser(user)

	now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
	staleDate := now.Add(-2 * time.Hour)
	freshDate := now.Add(-5 * time.Second)

	noExpiry := &model.Client{UserID: user.ID, Token: "C0", Name: "never", ExpiresAfterInactivitySeconds: 0, LastUsed: &staleDate}
	assert.NoError(s.T(), s.db.CreateClient(noExpiry))

	lastUsedStale := &model.Client{UserID: user.ID, Token: "C1", Name: "stale", LastUsed: &staleDate, ExpiresAfterInactivitySeconds: 60}
	assert.NoError(s.T(), s.db.CreateClient(lastUsedStale))

	lastUsedFresh := &model.Client{UserID: user.ID, Token: "C2", Name: "fresh", LastUsed: &freshDate, ExpiresAfterInactivitySeconds: 60}
	assert.NoError(s.T(), s.db.CreateClient(lastUsedFresh))

	createdAtFresh := &model.Client{UserID: user.ID, Token: "C3", Name: "unused", ExpiresAfterInactivitySeconds: 60, CreatedAt: freshDate}
	assert.NoError(s.T(), s.db.CreateClient(createdAtFresh))

	createdAtStale := &model.Client{UserID: user.ID, Token: "C4", Name: "unused", ExpiresAfterInactivitySeconds: 60, CreatedAt: staleDate}
	assert.NoError(s.T(), s.db.CreateClient(createdAtStale))

	expired, err := s.db.CleanupExpiredClients(now)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), expired, 2)

	expiredTokens := []string{}
	for _, c := range expired {
		expiredTokens = append(expiredTokens, c.Token)
	}
	assert.Contains(s.T(), expiredTokens, lastUsedStale.Token)
	assert.Contains(s.T(), expiredTokens, createdAtStale.Token)

	for _, id := range []uint{noExpiry.ID, lastUsedFresh.ID, createdAtFresh.ID} {
		if c, err := s.db.GetClientByID(id); assert.NoError(s.T(), err) {
			assert.NotNil(s.T(), c)
		}
	}
	for _, id := range []uint{lastUsedStale.ID, createdAtStale.ID} {
		if c, err := s.db.GetClientByID(id); assert.NoError(s.T(), err) {
			assert.Nil(s.T(), c)
		}
	}
}

func (s *DatabaseSuite) TestClientExpiresAtMemoization() {
	user := &model.User{Name: "memo", Pass: []byte{1}}
	s.db.CreateUser(user)

	// CreateClient populates expires_at from CreatedAt + seconds.
	client := &model.Client{
		UserID:                        user.ID,
		Token:                         "TOKEN",
		Name:                          "client",
		CreatedAt:                     now,
		ExpiresAfterInactivitySeconds: 60,
	}
	assert.NoError(s.T(), s.db.CreateClient(client))
	if assert.NotNil(s.T(), client.ExpiresAt) {
		assert.Equal(s.T(), now.Add(60*time.Second), *client.ExpiresAt)
	}

	fetched, err := s.db.GetClientByID(client.ID)
	if assert.NoError(s.T(), err) && assert.NotNil(s.T(), fetched) {
		assert.Equal(s.T(), client.ExpiresAt.Unix(), fetched.ExpiresAt.Unix())
	}

	used := now.Add(30 * time.Second)
	assert.NoError(s.T(), s.db.UpdateClientTokensLastUsedAndExpiresAt([]string{client.Token}, &used))
	fetched, err = s.db.GetClientByID(client.ID)
	if assert.NoError(s.T(), err) && assert.NotNil(s.T(), fetched) {
		assert.Equal(s.T(), used.Unix(), fetched.LastUsed.Unix())
		assert.Equal(s.T(), used.Add(60*time.Second).Unix(), fetched.ExpiresAt.Unix())
	}

	fetched.ExpiresAfterInactivitySeconds = 120
	assert.NoError(s.T(), s.db.UpdateClient(fetched))
	if assert.NotNil(s.T(), fetched.ExpiresAt) {
		assert.Equal(s.T(), used.Add(120*time.Second).Unix(), fetched.ExpiresAt.Unix())
	}

	fetched.ExpiresAfterInactivitySeconds = 0
	assert.NoError(s.T(), s.db.UpdateClient(fetched))
	assert.Nil(s.T(), fetched.ExpiresAt)
}

func (s *DatabaseSuite) TestClientSelectHidesExpired() {
	user := &model.User{Name: "predicate", Pass: []byte{1}}
	s.db.CreateUser(user)

	pastExpiry := now.Add(-time.Second)
	futureExpiry := now.Add(time.Hour)

	expired := &model.Client{
		UserID:                        user.ID,
		Token:                         "EXPIRED",
		Name:                          "expired",
		ExpiresAfterInactivitySeconds: 60,
		ExpiresAt:                     &pastExpiry,
	}
	assert.NoError(s.T(), s.db.DB.Create(expired).Error)

	valid := &model.Client{
		UserID:                        user.ID,
		Token:                         "VALID",
		Name:                          "valid",
		ExpiresAfterInactivitySeconds: 60,
		ExpiresAt:                     &futureExpiry,
	}
	assert.NoError(s.T(), s.db.DB.Create(valid).Error)

	if c, err := s.db.GetClientByID(expired.ID); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), c)
	}
	if c, err := s.db.GetClientByToken(expired.Token); assert.NoError(s.T(), err) {
		assert.Nil(s.T(), c)
	}
	if c, err := s.db.GetClientByID(valid.ID); assert.NoError(s.T(), err) {
		assert.NotNil(s.T(), c)
	}
	if c, err := s.db.GetClientByToken(valid.Token); assert.NoError(s.T(), err) {
		assert.NotNil(s.T(), c)
	}
	if cs, err := s.db.GetClientsByUser(user.ID); assert.NoError(s.T(), err) {
		assert.Len(s.T(), cs, 1)
		assert.Equal(s.T(), valid.Token, cs[0].Token)
	}
}
