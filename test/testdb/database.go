package testdb

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
	return ab.app(id, false)
}

// InternalApp creates an internal application and returns a message builder.
func (ab *AppClientBuilder) InternalApp(id uint) *MessageBuilder {
	return ab.app(id, true)
}

func (ab *AppClientBuilder) app(id uint, internal bool) *MessageBuilder {
	return ab.appWithToken(id, "app"+fmt.Sprint(id), internal)
}

// AppWithToken creates an application with a token and returns a message builder.
func (ab *AppClientBuilder) AppWithToken(id uint, token string) *MessageBuilder {
	return ab.appWithToken(id, token, false)
}

// InternalAppWithToken creates an internal application with a token and returns a message builder.
func (ab *AppClientBuilder) InternalAppWithToken(id uint, token string) *MessageBuilder {
	return ab.appWithToken(id, token, true)
}

func (ab *AppClientBuilder) appWithToken(id uint, token string, internal bool) *MessageBuilder {
	ab.newAppWithToken(id, token, internal)
	return &MessageBuilder{db: ab.db, appID: id}
}

// NewAppWithToken creates an application with a token and returns the app.
func (ab *AppClientBuilder) NewAppWithToken(id uint, token string) *model.Application {
	return ab.newAppWithToken(id, token, false)
}

// NewInternalAppWithToken creates an internal application with a token and returns the app.
func (ab *AppClientBuilder) NewInternalAppWithToken(id uint, token string) *model.Application {
	return ab.newAppWithToken(id, token, true)
}

func (ab *AppClientBuilder) newAppWithToken(id uint, token string, internal bool) *model.Application {
	application := &model.Application{ID: id, UserID: ab.userID, Token: token, Internal: internal}
	ab.db.CreateApplication(application)
	return application
}

// AppWithTokenAndName creates an application with a token and name and returns a message builder.
func (ab *AppClientBuilder) AppWithTokenAndName(id uint, token, name string) *MessageBuilder {
	return ab.appWithTokenAndName(id, token, name, false)
}

// InternalAppWithTokenAndName creates an internal application with a token and name and returns a message builder.
func (ab *AppClientBuilder) InternalAppWithTokenAndName(id uint, token, name string) *MessageBuilder {
	return ab.appWithTokenAndName(id, token, name, true)
}

func (ab *AppClientBuilder) appWithTokenAndName(id uint, token, name string, internal bool) *MessageBuilder {
	ab.newAppWithTokenAndName(id, token, name, internal)
	return &MessageBuilder{db: ab.db, appID: id}
}

// NewAppWithTokenAndName creates an application with a token and name and returns the app.
func (ab *AppClientBuilder) NewAppWithTokenAndName(id uint, token, name string) *model.Application {
	return ab.newAppWithTokenAndName(id, token, name, false)
}

// NewInternalAppWithTokenAndName creates an internal application with a token and name and returns the app.
func (ab *AppClientBuilder) NewInternalAppWithTokenAndName(id uint, token, name string) *model.Application {
	return ab.newAppWithTokenAndName(id, token, name, true)
}

func (ab *AppClientBuilder) newAppWithTokenAndName(id uint, token, name string, internal bool) *model.Application {
	application := &model.Application{ID: id, UserID: ab.userID, Token: token, Name: name, Internal: internal}
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
	if app, err := d.GetApplicationByID(id); assert.NoError(d.t, err) {
		assert.True(d.t, app == nil, "app %d must not exist", id)
	}
}

// AssertUserNotExist asserts that the user does not exist.
func (d *Database) AssertUserNotExist(id uint) {
	if user, err := d.GetUserByID(id); assert.NoError(d.t, err) {
		assert.True(d.t, user == nil, "user %d must not exist", id)
	}
}

// AssertClientNotExist asserts that the client does not exist.
func (d *Database) AssertClientNotExist(id uint) {
	if client, err := d.GetClientByID(id); assert.NoError(d.t, err) {
		assert.True(d.t, client == nil, "client %d must not exist", id)
	}
}

// AssertMessageNotExist asserts that the messages does not exist.
func (d *Database) AssertMessageNotExist(ids ...uint) {
	for _, id := range ids {
		if msg, err := d.GetMessageByID(id); assert.NoError(d.t, err) {
			assert.True(d.t, msg == nil, "message %d must not exist", id)
		}
	}
}

// AssertAppExist asserts that the app does exist.
func (d *Database) AssertAppExist(id uint) {
	if app, err := d.GetApplicationByID(id); assert.NoError(d.t, err) {
		assert.False(d.t, app == nil, "app %d must exist", id)
	}
}

// AssertUserExist asserts that the user does exist.
func (d *Database) AssertUserExist(id uint) {
	if user, err := d.GetUserByID(id); assert.NoError(d.t, err) {
		assert.False(d.t, user == nil, "user %d must exist", id)
	}
}

// AssertClientExist asserts that the client does exist.
func (d *Database) AssertClientExist(id uint) {
	if client, err := d.GetClientByID(id); assert.NoError(d.t, err) {
		assert.False(d.t, client == nil, "client %d must exist", id)
	}
}

// AssertMessageExist asserts that the message does exist.
func (d *Database) AssertMessageExist(id uint) {
	if msg, err := d.GetMessageByID(id); assert.NoError(d.t, err) {
		assert.False(d.t, msg == nil, "message %d must exist", id)
	}
}
