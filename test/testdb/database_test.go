package testdb_test

import (
	"testing"

	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
	"github.com/gotify/server/test/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func Test_WithDefault(t *testing.T) {
	db := testdb.NewDBWithDefaultUser(t)
	if user, err := db.GetUserByName("admin"); assert.NoError(t, err) {
		assert.NotNil(t, user)
	}
	db.Close()
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

type DatabaseSuite struct {
	suite.Suite
	db *testdb.Database
}

func (s *DatabaseSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	s.db = testdb.NewDB(s.T())
}

func (s *DatabaseSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

func (s *DatabaseSuite) Test_Users() {
	s.db.User(1)
	newUserActual := s.db.NewUser(2)
	s.db.NewUserWithName(3, "tom")

	newUserExpected := &model.User{ID: 2, Name: "user2"}

	assert.Equal(s.T(), newUserExpected, newUserActual)

	users := []*model.User{{ID: 1, Name: "user1"}, {ID: 2, Name: "user2"}, {ID: 3, Name: "tom"}}

	if usersActual, err := s.db.GetUsers(); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), users, usersActual)
	}
	s.db.AssertUserExist(1)
	s.db.AssertUserExist(2)
	s.db.AssertUserExist(3)
	s.db.AssertUserNotExist(4)

	s.db.DeleteUserByID(2)

	s.db.AssertUserNotExist(2)
}

func (s *DatabaseSuite) Test_Clients() {
	userBuilder := s.db.User(1)
	userBuilder.Client(1)
	newClientActual := userBuilder.NewClientWithToken(2, "asdf")

	s.db.User(2).Client(5)

	newClientExpected := &model.Client{ID: 2, Token: "asdf", UserID: 1}

	assert.Equal(s.T(), newClientExpected, newClientActual)

	userOneExpected := []*model.Client{{ID: 1, Token: "client1", UserID: 1}, {ID: 2, Token: "asdf", UserID: 1}}
	if clients, err := s.db.GetClientsByUser(1); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), userOneExpected, clients)
	}
	userTwoExpected := []*model.Client{{ID: 5, Token: "client5", UserID: 2}}
	if clients, err := s.db.GetClientsByUser(2); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), userTwoExpected, clients)
	}

	s.db.AssertClientExist(1)
	s.db.AssertClientExist(2)
	s.db.AssertClientNotExist(3)
	s.db.AssertClientNotExist(4)
	s.db.AssertClientExist(5)
	s.db.AssertClientNotExist(6)

	s.db.DeleteClientByID(2)

	s.db.AssertClientNotExist(2)
}

func (s *DatabaseSuite) Test_Apps() {
	userBuilder := s.db.User(1)
	userBuilder.App(1)
	newAppActual := userBuilder.NewAppWithToken(2, "asdf")
	newInternalAppActual := userBuilder.NewInternalAppWithToken(3, "qwer")

	s.db.User(2).InternalApp(5)

	newAppExpected := &model.Application{ID: 2, Token: "asdf", UserID: 1}
	newInternalAppExpected := &model.Application{ID: 3, Token: "qwer", UserID: 1, Internal: true}

	assert.Equal(s.T(), newAppExpected, newAppActual)
	assert.Equal(s.T(), newInternalAppExpected, newInternalAppActual)

	userOneExpected := []*model.Application{{ID: 1, Token: "app1", UserID: 1}, {ID: 2, Token: "asdf", UserID: 1}, {ID: 3, Token: "qwer", UserID: 1, Internal: true}}
	if app, err := s.db.GetApplicationsByUser(1); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), userOneExpected, app)
	}
	userTwoExpected := []*model.Application{{ID: 5, Token: "app5", UserID: 2, Internal: true}}
	if app, err := s.db.GetApplicationsByUser(2); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), userTwoExpected, app)
	}

	newAppWithName := userBuilder.NewAppWithTokenAndName(7, "test-token", "app name")
	newAppWithNameExpected := &model.Application{ID: 7, Token: "test-token", UserID: 1, Name: "app name"}
	assert.Equal(s.T(), newAppWithNameExpected, newAppWithName)

	newInternalAppWithName := userBuilder.NewInternalAppWithTokenAndName(8, "test-tokeni", "app name")
	newInternalAppWithNameExpected := &model.Application{ID: 8, Token: "test-tokeni", UserID: 1, Name: "app name", Internal: true}
	assert.Equal(s.T(), newInternalAppWithNameExpected, newInternalAppWithName)

	userBuilder.AppWithTokenAndName(9, "test-token-2", "app name")
	userBuilder.InternalAppWithTokenAndName(10, "test-tokeni-2", "app name")
	userBuilder.AppWithToken(11, "test-token-3")
	userBuilder.InternalAppWithToken(12, "test-tokeni-3")

	s.db.AssertAppExist(1)
	s.db.AssertAppExist(2)
	s.db.AssertAppExist(3)
	s.db.AssertAppNotExist(4)
	s.db.AssertAppExist(5)
	s.db.AssertAppNotExist(6)
	s.db.AssertAppExist(7)
	s.db.AssertAppExist(8)
	s.db.AssertAppExist(9)
	s.db.AssertAppExist(10)
	s.db.AssertAppExist(11)
	s.db.AssertAppExist(12)

	s.db.DeleteApplicationByID(2)

	s.db.AssertAppNotExist(2)
}

func (s *DatabaseSuite) Test_Messages() {
	s.db.User(1).App(1).Message(1).Message(2)
	s.db.User(2).App(2).Message(4).Message(5)

	userOneExpected := []*model.Message{{ID: 2, ApplicationID: 1}, {ID: 1, ApplicationID: 1}}
	if msgs, err := s.db.GetMessagesByUser(1); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), userOneExpected, msgs)
	}
	userTwoExpected := []*model.Message{{ID: 5, ApplicationID: 2}, {ID: 4, ApplicationID: 2}}
	if msgs, err := s.db.GetMessagesByUser(2); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), userTwoExpected, msgs)
	}

	s.db.AssertMessageExist(1)
	s.db.AssertMessageExist(2)
	s.db.AssertMessageExist(4)
	s.db.AssertMessageExist(5)

	s.db.AssertMessageNotExist(3, 6, 7, 8)

	s.db.DeleteMessageByID(2)

	s.db.AssertMessageNotExist(2)
}
