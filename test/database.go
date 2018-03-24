package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gotify/server/database"
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
)

// Database is the wrapper for the gorm database with sleek helper methods.
type Database struct {
	*database.GormDatabase
	t *testing.T
}

// AppClientBuilder has helper methods to create applications and clients.
type AppClientBuilder struct {
	userID uint
	db     *Database
}

// MessageBuilder has helper methods to create messages.
type MessageBuilder struct {
	appID uint
	db    *Database
}

// NewDBWithDefaultUser creates a new test db instance with the default user.
func NewDBWithDefaultUser(t *testing.T) *Database {
	db, err := database.New("sqlite3", fmt.Sprintf("file:%s?mode=memory&cache=shared", fmt.Sprint(time.Now().Unix())), "admin", "pw", 5, true)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	return &Database{GormDatabase: db, t: t}
}

// NewDB creates a new test db instance.
func NewDB(t *testing.T) *Database {
	db, err := database.New("sqlite3", fmt.Sprintf("file:%s?mode=memory&cache=shared", fmt.Sprint(time.Now().Unix())), "admin", "pw", 5, false)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	return &Database{GormDatabase: db, t: t}
}

// User creates a user and returns a builder for applications and clients.
func (d *Database) User(id uint) *AppClientBuilder {
	d.NewUser(id)
	return &AppClientBuilder{db: d, userID: id}
}

// NewUser creates a user and returns the user.
func (d *Database) NewUser(id uint) *model.User {
	return d.NewUserWithName(id, "user"+fmt.Sprint(id))
}

// NewUserWithName creates a user with a name and returns the user.
func (d *Database) NewUserWithName(id uint, name string) *model.User {
	user := &model.User{ID: id, Name: name}
	d.CreateUser(user)
	return user
}

// App creates an application and returns a message builder.
func (ab *AppClientBuilder) App(id uint) *MessageBuilder {
	return ab.AppWithToken(id, "app"+fmt.Sprint(id))
}

// AppWithToken creates an application with a token and returns a message builder.
func (ab *AppClientBuilder) AppWithToken(id uint, token string) *MessageBuilder {
	ab.NewAppWithToken(id, token)
	return &MessageBuilder{db: ab.db, appID: id}
}

// NewAppWithToken creates an application with a token and returns the app.
func (ab *AppClientBuilder) NewAppWithToken(id uint, token string) *model.Application {
	application := &model.Application{ID: id, UserID: ab.userID, Token: token}
	ab.db.CreateApplication(application)
	return application
}

// Client creates a client and returns itself.
func (ab *AppClientBuilder) Client(id uint) *AppClientBuilder {
	return ab.ClientWithToken(id, "client"+fmt.Sprint(id))
}

// ClientWithToken creates a client with a token and returns itself.
func (ab *AppClientBuilder) ClientWithToken(id uint, token string) *AppClientBuilder {
	ab.NewClientWithToken(id, token)
	return ab
}

// NewClientWithToken creates a client with a token and returns the client.
func (ab *AppClientBuilder) NewClientWithToken(id uint, token string) *model.Client {
	client := &model.Client{ID: id, Token: token, UserID: ab.userID}
	ab.db.CreateClient(client)
	return client
}

// Message creates a message and returns itself
func (mb *MessageBuilder) Message(id uint) *MessageBuilder {
	mb.NewMessage(id)
	return mb
}

// NewMessage creates a message and returns the message.
func (mb *MessageBuilder) NewMessage(id uint) model.Message {
	message := model.Message{ID: id, ApplicationID: mb.appID}
	mb.db.CreateMessage(&message)
	return message
}

// AssertAppNotExist asserts that the app does not exist.
func (d *Database) AssertAppNotExist(id uint) {
	assert.True(d.t, d.GetApplicationByID(id) == nil, "app %d must not exist", id)
}

// AssertUserNotExist asserts that the user does not exist.
func (d *Database) AssertUserNotExist(id uint) {
	assert.True(d.t, d.GetUserByID(id) == nil, "user %d must not exist", id)
}

// AssertClientNotExist asserts that the client does not exist.
func (d *Database) AssertClientNotExist(id uint) {
	assert.True(d.t, d.GetClientByID(id) == nil, "client %d must not exist", id)
}

// AssertMessageNotExist asserts that the messages does not exist.
func (d *Database) AssertMessageNotExist(ids ...uint) {
	for _, id := range ids {
		assert.True(d.t, d.GetMessageByID(id) == nil, "message %d must not exist", id)
	}
}

// AssertAppExist asserts that the app does exist.
func (d *Database) AssertAppExist(id uint) {
	assert.False(d.t, d.GetApplicationByID(id) == nil, "app %d must exist", id)
}

// AssertUserExist asserts that the user does exist.
func (d *Database) AssertUserExist(id uint) {
	assert.False(d.t, d.GetUserByID(id) == nil, "user %d must exist", id)
}

// AssertClientExist asserts that the client does exist.
func (d *Database) AssertClientExist(id uint) {
	assert.False(d.t, d.GetClientByID(id) == nil, "client %d must exist", id)
}

// AssertMessageExist asserts that the message does exist.
func (d *Database) AssertMessageExist(id uint) {
	assert.False(d.t, d.GetMessageByID(id) == nil, "message %d must exist", id)
}
