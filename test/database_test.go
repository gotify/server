package test_test

import (
	"testing"

	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
	"github.com/gotify/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func Test_WithDefault(t *testing.T) {
	db := test.NewDBWithDefaultUser(t)
	assert.NotNil(t, db.GetUserByName("admin"))
	db.Close()
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

type DatabaseSuite struct {
	suite.Suite
	db *test.Database
}

func (s *DatabaseSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	s.db = test.NewDB(s.T())
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

	assert.Equal(s.T(), users, s.db.GetUsers())
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
	assert.Equal(s.T(), userOneExpected, s.db.GetClientsByUser(1))
	userTwoExpected := []*model.Client{{ID: 5, Token: "client5", UserID: 2}}
	assert.Equal(s.T(), userTwoExpected, s.db.GetClientsByUser(2))

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

	s.db.User(2).App(5)

	newAppExpected := &model.Application{ID: 2, Token: "asdf", UserID: 1}

	assert.Equal(s.T(), newAppExpected, newAppActual)

	userOneExpected := []*model.Application{{ID: 1, Token: "app1", UserID: 1}, {ID: 2, Token: "asdf", UserID: 1}}
	assert.Equal(s.T(), userOneExpected, s.db.GetApplicationsByUser(1))
	userTwoExpected := []*model.Application{{ID: 5, Token: "app5", UserID: 2}}
	assert.Equal(s.T(), userTwoExpected, s.db.GetApplicationsByUser(2))

	s.db.AssertAppExist(1)
	s.db.AssertAppExist(2)
	s.db.AssertAppNotExist(3)
	s.db.AssertAppNotExist(4)
	s.db.AssertAppExist(5)
	s.db.AssertAppNotExist(6)

	s.db.DeleteApplicationByID(2)

	s.db.AssertAppNotExist(2)
}

func (s *DatabaseSuite) Test_Messages() {
	s.db.User(1).App(1).Message(1).Message(2)
	s.db.User(2).App(2).Message(4).Message(5)

	userOneExpected := []*model.Message{{ID: 1, ApplicationID: 1}, {ID: 2, ApplicationID: 1}}
	assert.Equal(s.T(), userOneExpected, s.db.GetMessagesByUser(1))
	userTwoExpected := []*model.Message{{ID: 4, ApplicationID: 2}, {ID: 5, ApplicationID: 2}}
	assert.Equal(s.T(), userTwoExpected, s.db.GetMessagesByUser(2))

	s.db.AssertMessageExist(1)
	s.db.AssertMessageExist(2)
	s.db.AssertMessageExist(4)
	s.db.AssertMessageExist(5)

	s.db.AssertMessageNotExist(3, 6, 7, 8)

	s.db.DeleteMessageByID(2)

	s.db.AssertMessageNotExist(2)
}
