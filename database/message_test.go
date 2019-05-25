package database

import (
	"testing"
	"time"

	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DatabaseSuite) TestMessage() {
	messages, err := s.db.GetMessageByID(5)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), messages, "not existing message")

	user := &model.User{Name: "test", Pass: []byte{1}}
	s.db.CreateUser(user)
	assert.NotEqual(s.T(), 0, user.ID)

	backupServer := &model.Application{UserID: user.ID, Token: "A0000000000", Name: "backupserver"}
	s.db.CreateApplication(backupServer)
	assert.NotEqual(s.T(), 0, backupServer.ID)

	msgs, err := s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), msgs)

	msgs, err = s.db.GetMessagesByApplication(backupServer.ID)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), msgs)

	backupdone := &model.Message{ApplicationID: backupServer.ID, Message: "backup done", Title: "backup", Priority: 1, Date: time.Now()}
	require.NoError(s.T(), s.db.CreateMessage(backupdone))
	assert.NotEqual(s.T(), 0, backupdone.ID)

	messages, err = s.db.GetMessageByID(backupdone.ID)
	require.NoError(s.T(), err)
	assertEquals(s.T(), messages, backupdone)

	msgs, err = s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], backupdone)

	msgs, err = s.db.GetMessagesByApplication(backupServer.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], backupdone)

	loginServer := &model.Application{UserID: user.ID, Token: "A0000000001", Name: "loginserver"}
	require.NoError(s.T(), s.db.CreateApplication(loginServer))
	assert.NotEqual(s.T(), 0, loginServer.ID)

	logindone := &model.Message{ApplicationID: loginServer.ID, Message: "login done", Title: "login", Priority: 1, Date: time.Now()}
	require.NoError(s.T(), s.db.CreateMessage(logindone))
	assert.NotEqual(s.T(), 0, logindone.ID)

	msgs, err = s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 2)
	assertEquals(s.T(), msgs[0], logindone)
	assertEquals(s.T(), msgs[1], backupdone)

	msgs, err = s.db.GetMessagesByApplication(backupServer.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], backupdone)

	loginfailed := &model.Message{ApplicationID: loginServer.ID, Message: "login failed", Title: "login", Priority: 1, Date: time.Now()}
	require.NoError(s.T(), s.db.CreateMessage(loginfailed))
	assert.NotEqual(s.T(), 0, loginfailed.ID)

	msgs, err = s.db.GetMessagesByApplication(backupServer.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], backupdone)

	msgs, err = s.db.GetMessagesByApplication(loginServer.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 2)
	assertEquals(s.T(), msgs[0], loginfailed)
	assertEquals(s.T(), msgs[1], logindone)

	msgs, err = s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 3)
	assertEquals(s.T(), msgs[0], loginfailed)
	assertEquals(s.T(), msgs[1], logindone)
	assertEquals(s.T(), msgs[2], backupdone)

	backupfailed := &model.Message{ApplicationID: backupServer.ID, Message: "backup failed", Title: "backup", Priority: 1, Date: time.Now()}
	require.NoError(s.T(), s.db.CreateMessage(backupfailed))
	assert.NotEqual(s.T(), 0, backupfailed.ID)

	msgs, err = s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 4)
	assertEquals(s.T(), msgs[0], backupfailed)
	assertEquals(s.T(), msgs[1], loginfailed)
	assertEquals(s.T(), msgs[2], logindone)
	assertEquals(s.T(), msgs[3], backupdone)

	msgs, err = s.db.GetMessagesByApplication(loginServer.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 2)
	assertEquals(s.T(), msgs[0], loginfailed)
	assertEquals(s.T(), msgs[1], logindone)

	require.NoError(s.T(), s.db.DeleteMessagesByApplication(loginServer.ID))
	msgs, err = s.db.GetMessagesByApplication(loginServer.ID)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), msgs)

	msgs, err = s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 2)
	assertEquals(s.T(), msgs[0], backupfailed)
	assertEquals(s.T(), msgs[1], backupdone)

	logindone = &model.Message{ApplicationID: loginServer.ID, Message: "login done", Title: "login", Priority: 1, Date: time.Now()}
	require.NoError(s.T(), s.db.CreateMessage(logindone))
	assert.NotEqual(s.T(), 0, logindone.ID)

	msgs, err = s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 3)
	assertEquals(s.T(), msgs[0], logindone)
	assertEquals(s.T(), msgs[1], backupfailed)
	assertEquals(s.T(), msgs[2], backupdone)

	s.db.DeleteMessagesByUser(user.ID)
	msgs, err = s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), msgs)

	logout := &model.Message{ApplicationID: loginServer.ID, Message: "logout success", Title: "logout", Priority: 1, Date: time.Now()}
	s.db.CreateMessage(logout)
	msgs, err = s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], logout)

	require.NoError(s.T(), s.db.DeleteMessageByID(logout.ID))
	msgs, err = s.db.GetMessagesByUser(user.ID)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), msgs)

}

func (s *DatabaseSuite) TestGetMessagesSince() {
	user := &model.User{Name: "test", Pass: []byte{1}}
	require.NoError(s.T(), s.db.CreateUser(user))

	app := &model.Application{UserID: user.ID, Token: "A0000000000"}
	app2 := &model.Application{UserID: user.ID, Token: "A0000000001"}
	require.NoError(s.T(), s.db.CreateApplication(app))
	require.NoError(s.T(), s.db.CreateApplication(app2))

	curDate := time.Now()
	for i := 1; i <= 500; i++ {
		s.db.CreateMessage(&model.Message{ApplicationID: app.ID, Message: "abc", Date: curDate.Add(time.Duration(i) * time.Second)})
		s.db.CreateMessage(&model.Message{ApplicationID: app2.ID, Message: "abc", Date: curDate.Add(time.Duration(i) * time.Second)})
	}

	actual, err := s.db.GetMessagesByUserSince(user.ID, 50, 0)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 50)
	hasIDInclusiveBetween(s.T(), actual, 1000, 951, 1)

	actual, err = s.db.GetMessagesByUserSince(user.ID, 50, 951)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 50)
	hasIDInclusiveBetween(s.T(), actual, 950, 901, 1)

	actual, err = s.db.GetMessagesByUserSince(user.ID, 100, 951)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 100)
	hasIDInclusiveBetween(s.T(), actual, 950, 851, 1)

	actual, err = s.db.GetMessagesByUserSince(user.ID, 100, 51)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 50)
	hasIDInclusiveBetween(s.T(), actual, 50, 1, 1)

	actual, err = s.db.GetMessagesByApplicationSince(app.ID, 50, 0)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 50)
	hasIDInclusiveBetween(s.T(), actual, 999, 901, 2)

	actual, err = s.db.GetMessagesByApplicationSince(app.ID, 50, 901)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 50)
	hasIDInclusiveBetween(s.T(), actual, 899, 801, 2)

	actual, err = s.db.GetMessagesByApplicationSince(app.ID, 100, 666)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 100)
	hasIDInclusiveBetween(s.T(), actual, 665, 467, 2)

	actual, err = s.db.GetMessagesByApplicationSince(app.ID, 100, 101)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 50)
	hasIDInclusiveBetween(s.T(), actual, 99, 1, 2)

	actual, err = s.db.GetMessagesByApplicationSince(app2.ID, 50, 0)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 50)
	hasIDInclusiveBetween(s.T(), actual, 1000, 902, 2)

	actual, err = s.db.GetMessagesByApplicationSince(app2.ID, 50, 902)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 50)
	hasIDInclusiveBetween(s.T(), actual, 900, 802, 2)

	actual, err = s.db.GetMessagesByApplicationSince(app2.ID, 100, 667)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 100)
	hasIDInclusiveBetween(s.T(), actual, 666, 468, 2)

	actual, err = s.db.GetMessagesByApplicationSince(app2.ID, 100, 102)
	require.NoError(s.T(), err)
	assert.Len(s.T(), actual, 50)
	hasIDInclusiveBetween(s.T(), actual, 100, 2, 2)

}

func hasIDInclusiveBetween(t *testing.T, msgs []*model.Message, from, to, decrement int) {
	index := 0
	for expectedID := from; expectedID >= to; expectedID -= decrement {
		if !assert.Equal(t, uint(expectedID), msgs[index].ID) {
			break
		}
		index++
	}
	assert.Equal(t, index, len(msgs), "not all entries inside msgs were checked")
}

// assertEquals compares messages and correctly check dates
func assertEquals(t *testing.T, left *model.Message, right *model.Message) {
	assert.Equal(t, left.Date.Unix(), right.Date.Unix())
	left.Date = right.Date
	assert.Equal(t, left, right)
}
