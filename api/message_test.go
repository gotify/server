package api

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
	"github.com/gotify/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"strings"

	"github.com/bouk/monkey"
	"github.com/gotify/server/auth"
)

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageSuite))
}

type MessageSuite struct {
	suite.Suite
	db       *test.Database
	a        *MessageAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
	notified bool
}

func (s *MessageSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.db = test.NewDB(s.T())
	s.notified = false
	s.a = &MessageAPI{DB: s.db, Notifier: s}
}

func (s *MessageSuite) AfterTest(string, string) {
	s.db.Close()
}

func (s *MessageSuite) Notify(userID uint, msg *model.Message) {
	s.notified = true
}

func (s *MessageSuite) Test_ensureCorrectJsonRepresentation() {
	t, _ := time.Parse("2006/01/02", "2017/01/02")

	actual := model.Message{ID: 55, ApplicationID: 2, Message: "hi", Title: "hi", Date: t, Priority: 4}
	test.JSONEquals(s.T(), actual, `{"id":55,"appid":2,"message":"hi","title":"hi","priority":4,"date":"2017-01-02T00:00:00Z"}`)
}

func (s *MessageSuite) Test_GetMessages() {
	user := s.db.User(5)
	first := user.App(1).NewMessage(1)
	second := user.App(2).NewMessage(2)

	test.WithUser(s.ctx, 5)
	s.a.GetMessages(s.ctx)

	test.BodyEquals(s.T(), &[]model.Message{first, second}, s.recorder)
}

func (s *MessageSuite) Test_GetMessagesWithToken() {
	expected := s.db.User(4).App(2).NewMessage(1)

	test.WithUser(s.ctx, 4)
	s.ctx.Params = gin.Params{{Key: "id", Value: "2"}}
	s.a.GetMessagesWithApplication(s.ctx)

	test.BodyEquals(s.T(), []model.Message{expected}, s.recorder)
}

func (s *MessageSuite) Test_GetMessagesWithToken_withWrongUser_expectNotFound() {
	s.db.User(4)
	s.db.User(5).App(2).Message(66)

	test.WithUser(s.ctx, 4)
	s.ctx.Params = gin.Params{{Key: "id", Value: "2"}}
	s.a.GetMessagesWithApplication(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessage_invalidID() {
	s.ctx.Params = gin.Params{{Key: "id", Value: "string"}}

	s.a.DeleteMessage(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessage_notExistingID() {
	s.db.User(1).App(5).Message(55)

	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	s.a.DeleteMessage(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessage_existingIDButNotOwner() {
	s.db.User(1).App(10).Message(100)
	s.db.User(2)

	test.WithUser(s.ctx, 2)
	s.ctx.Params = gin.Params{{Key: "id", Value: "100"}}
	s.a.DeleteMessage(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessage() {
	s.db.User(6).App(1).Message(50)

	test.WithUser(s.ctx, 6)
	s.ctx.Params = gin.Params{{Key: "id", Value: "50"}}
	s.a.DeleteMessage(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertMessageNotExist(50)
}

func (s *MessageSuite) Test_DeleteMessageWithID() {
	s.db.User(2).AppWithToken(5, "mytoken").Message(55)

	test.WithUser(s.ctx, 2)
	s.ctx.Params = gin.Params{{Key: "id", Value: "5"}}
	s.a.DeleteMessageWithApplication(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertMessageNotExist(55)
}

func (s *MessageSuite) Test_DeleteMessageWithToken_notExistingID() {
	s.db.User(2).AppWithToken(1, "wrong").Message(1)

	test.WithUser(s.ctx, 2)
	s.ctx.Params = gin.Params{{Key: "id", Value: "55"}}
	s.a.DeleteMessageWithApplication(s.ctx)

	s.db.AssertMessageExist(1)
	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessageWithToken_notOwner() {
	s.db.User(4)
	s.db.User(2).App(55).Message(5)

	test.WithUser(s.ctx, 4)
	s.ctx.Params = gin.Params{{Key: "id", Value: "55"}}
	s.a.DeleteMessageWithApplication(s.ctx)

	s.db.AssertMessageExist(5)
	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessages() {
	userBuilder := s.db.User(4)
	userBuilder.App(5).Message(5).Message(6)
	userBuilder.App(2).Message(7).Message(8)
	s.db.User(5).App(7).Message(22)

	test.WithUser(s.ctx, 4)
	s.a.DeleteMessages(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertMessageExist(22)
	s.db.AssertMessageNotExist(5, 6, 7, 8)
}

func (s *MessageSuite) Test_CreateMessage_onJson_allParams() {
	t, _ := time.Parse("2006/01/02", "2017/01/02")
	patch := monkey.Patch(time.Now, func() time.Time { return t })
	defer patch.Unpatch()

	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")
	s.db.User(4).AppWithToken(7, "app-token")
	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"title": "mytitle", "message": "mymessage", "priority": 1}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateMessage(s.ctx)

	msgs := s.db.GetMessagesByApplication(7)
	expected := &model.Message{ID: 1, ApplicationID: 7, Title: "mytitle", Message: "mymessage", Priority: 1, Date: t}
	assert.Len(s.T(), msgs, 1)
	assert.Contains(s.T(), msgs, expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.True(s.T(), s.notified)
}

func (s *MessageSuite) Test_CreateMessage_onlyRequired() {
	t, _ := time.Parse("2006/01/02", "2017/01/02")
	patch := monkey.Patch(time.Now, func() time.Time { return t })
	defer patch.Unpatch()

	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")
	s.db.User(4).AppWithToken(5, "app-token")
	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"title": "mytitle", "message": "mymessage"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateMessage(s.ctx)

	msgs := s.db.GetMessagesByApplication(5)
	expected := &model.Message{ID: 1, ApplicationID: 5, Title: "mytitle", Message: "mymessage", Date: t}
	assert.Len(s.T(), msgs, 1)
	assert.Contains(s.T(), msgs, expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.True(s.T(), s.notified)
}

func (s *MessageSuite) Test_CreateMessage_failWhenNoMessage() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")
	s.db.User(4).AppWithToken(1, "app-token")

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"title": "mytitle"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateMessage(s.ctx)

	assert.Empty(s.T(), s.db.GetMessagesByApplication(1))
	assert.Equal(s.T(), 400, s.recorder.Code)
	assert.False(s.T(), s.notified)
}

func (s *MessageSuite) Test_CreateMessage_failWhenNoTitle() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")
	s.db.User(4).AppWithToken(8, "app-token")

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"message": "mymessage"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateMessage(s.ctx)

	assert.Empty(s.T(), s.db.GetMessagesByApplication(8))
	assert.Equal(s.T(), 400, s.recorder.Code)
	assert.False(s.T(), s.notified)
}

func (s *MessageSuite) Test_CreateMessage_failWhenPriorityNotNumber() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")
	s.db.User(4).AppWithToken(8, "app-token")

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"title": "mytitle", "message": "mymessage", "priority": "asd"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateMessage(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
	assert.False(s.T(), s.notified)
	assert.Empty(s.T(), s.db.GetMessagesByApplication(1))
}

func (s *MessageSuite) Test_CreateMessage_onQueryData() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")
	s.db.User(4).AppWithToken(2, "app-token")

	t, _ := time.Parse("2006/01/02", "2017/01/02")
	patch := monkey.Patch(time.Now, func() time.Time { return t })
	defer patch.Unpatch()

	s.ctx.Request = httptest.NewRequest("POST", "/token?title=mytitle&message=mymessage&priority=1", nil)
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.a.CreateMessage(s.ctx)

	expected := &model.Message{ID: 1, ApplicationID: 2, Title: "mytitle", Message: "mymessage", Priority: 1, Date: t}

	msgs := s.db.GetMessagesByApplication(2)
	assert.Len(s.T(), msgs, 1)
	assert.Contains(s.T(), msgs, expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.True(s.T(), s.notified)
}

func (s *MessageSuite) Test_CreateMessage_onFormData() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")
	s.db.User(4).AppWithToken(99, "app-token")

	t, _ := time.Parse("2006/01/02", "2017/01/02")
	patch := monkey.Patch(time.Now, func() time.Time { return t })
	defer patch.Unpatch()

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader("title=mytitle&message=mymessage&priority=1"))
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.a.CreateMessage(s.ctx)

	expected := &model.Message{ID: 1, ApplicationID: 99, Title: "mytitle", Message: "mymessage", Priority: 1, Date: t}
	msgs := s.db.GetMessagesByApplication(99)
	assert.Len(s.T(), msgs, 1)
	assert.Contains(s.T(), msgs, expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.True(s.T(), s.notified)
}
