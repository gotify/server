package database

import (
	"testing"
	"time"

	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
)

func (s *DatabaseSuite) TestMessage() {
	assert.Nil(s.T(), s.db.GetMessageByID(5), "not existing message")

	user := &model.User{Name: "test", Pass: []byte{1}}
	s.db.CreateUser(user)
	assert.NotEqual(s.T(), 0, user.ID)

	backupServer := &model.Application{UserID: user.ID, ID: "A0000000000", Name: "backupserver"}
	s.db.CreateApplication(backupServer)
	assert.NotEqual(s.T(), 0, backupServer.ID)

	msgs := s.db.GetMessagesByUser(user.ID)
	assert.Empty(s.T(), msgs)
	msgs = s.db.GetMessagesByApplication(backupServer.ID)
	assert.Empty(s.T(), msgs)

	backupdone := &model.Message{ApplicationID: backupServer.ID, Message: "backup done", Title: "backup", Priority: 1, Date: time.Now()}
	s.db.CreateMessage(backupdone)
	assert.NotEqual(s.T(), 0, backupdone.ID)

	assertEquals(s.T(), s.db.GetMessageByID(backupdone.ID), backupdone)

	msgs = s.db.GetMessagesByUser(user.ID)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], backupdone)

	msgs = s.db.GetMessagesByApplication(backupServer.ID)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], backupdone)

	loginServer := &model.Application{UserID: user.ID, ID: "A0000000001", Name: "loginserver"}
	s.db.CreateApplication(loginServer)
	assert.NotEqual(s.T(), 0, loginServer.ID)

	logindone := &model.Message{ApplicationID: loginServer.ID, Message: "login done", Title: "login", Priority: 1, Date: time.Now()}
	s.db.CreateMessage(logindone)
	assert.NotEqual(s.T(), 0, logindone.ID)

	msgs = s.db.GetMessagesByUser(user.ID)
	assert.Len(s.T(), msgs, 2)
	assertEquals(s.T(), msgs[0], logindone)
	assertEquals(s.T(), msgs[1], backupdone)

	msgs = s.db.GetMessagesByApplication(backupServer.ID)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], backupdone)

	loginfailed := &model.Message{ApplicationID: loginServer.ID, Message: "login failed", Title: "login", Priority: 1, Date: time.Now()}
	s.db.CreateMessage(loginfailed)
	assert.NotEqual(s.T(), 0, loginfailed.ID)

	msgs = s.db.GetMessagesByApplication(backupServer.ID)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], backupdone)

	msgs = s.db.GetMessagesByApplication(loginServer.ID)
	assert.Len(s.T(), msgs, 2)
	assertEquals(s.T(), msgs[0], loginfailed)
	assertEquals(s.T(), msgs[1], logindone)

	msgs = s.db.GetMessagesByUser(user.ID)
	assert.Len(s.T(), msgs, 3)
	assertEquals(s.T(), msgs[0], loginfailed)
	assertEquals(s.T(), msgs[1], logindone)
	assertEquals(s.T(), msgs[2], backupdone)

	backupfailed := &model.Message{ApplicationID: backupServer.ID, Message: "backup failed", Title: "backup", Priority: 1, Date: time.Now()}
	s.db.CreateMessage(backupfailed)
	assert.NotEqual(s.T(), 0, backupfailed.ID)

	msgs = s.db.GetMessagesByUser(user.ID)
	assert.Len(s.T(), msgs, 4)
	assertEquals(s.T(), msgs[0], backupfailed)
	assertEquals(s.T(), msgs[1], loginfailed)
	assertEquals(s.T(), msgs[2], logindone)
	assertEquals(s.T(), msgs[3], backupdone)

	msgs = s.db.GetMessagesByApplication(loginServer.ID)
	assert.Len(s.T(), msgs, 2)
	assertEquals(s.T(), msgs[0], loginfailed)
	assertEquals(s.T(), msgs[1], logindone)

	s.db.DeleteMessagesByApplication(loginServer.ID)
	assert.Empty(s.T(), s.db.GetMessagesByApplication(loginServer.ID))

	msgs = s.db.GetMessagesByUser(user.ID)
	assert.Len(s.T(), msgs, 2)
	assertEquals(s.T(), msgs[0], backupfailed)
	assertEquals(s.T(), msgs[1], backupdone)

	logindone = &model.Message{ApplicationID: loginServer.ID, Message: "login done", Title: "login", Priority: 1, Date: time.Now()}
	s.db.CreateMessage(logindone)
	assert.NotEqual(s.T(), 0, logindone.ID)

	msgs = s.db.GetMessagesByUser(user.ID)
	assert.Len(s.T(), msgs, 3)
	assertEquals(s.T(), msgs[0], logindone)
	assertEquals(s.T(), msgs[1], backupfailed)
	assertEquals(s.T(), msgs[2], backupdone)

	s.db.DeleteMessagesByUser(user.ID)
	assert.Empty(s.T(), s.db.GetMessagesByUser(user.ID))

	logout := &model.Message{ApplicationID: loginServer.ID, Message: "logout success", Title: "logout", Priority: 1, Date: time.Now()}
	s.db.CreateMessage(logout)
	msgs = s.db.GetMessagesByUser(user.ID)
	assert.Len(s.T(), msgs, 1)
	assertEquals(s.T(), msgs[0], logout)

	s.db.DeleteMessageByID(logout.ID)
	assert.Empty(s.T(), s.db.GetMessagesByUser(user.ID))
}

// assertEquals compares messages and correctly check dates
func assertEquals(t *testing.T, left *model.Message, right *model.Message) {
	assert.Equal(t, left.Date.Unix(), right.Date.Unix())
	left.Date = right.Date
	assert.Equal(t, left, right)
}
