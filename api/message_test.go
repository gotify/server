package api

import (
	"github.com/bouk/monkey"
	"github.com/gin-gonic/gin"
	apimock "github.com/jmattheis/memo/api/mock"
	"github.com/jmattheis/memo/auth"
	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageSuite))
}

type MessageSuite struct {
	suite.Suite
	db       *apimock.MockMessageDatabase
	a        *MessageAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func (s *MessageSuite) BeforeTest(suiteName, testName string) {
	gin.SetMode(gin.TestMode)
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.db = &apimock.MockMessageDatabase{}
	s.a = &MessageAPI{DB: s.db}
}

func (s *MessageSuite) Test_GetMessages() {
	auth.RegisterAuthentication(s.ctx, nil, 5, "")
	t, _ := time.Parse("2006/01/02", "2017/01/02")
	s.db.On("GetMessagesByUser", uint(5)).Return([]*model.Message{{ID: 1, ApplicationID: "asd", Message: "OH HELLO THERE", Date: t, Title: "wup", Priority: 2}, {ID: 2, ApplicationID: "cloud", Message: "hi", Title: "hi", Date: t, Priority: 4}})

	s.a.GetMessages(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	bytes, _ := ioutil.ReadAll(s.recorder.Body)

	assert.JSONEq(s.T(), `[{"id":1,"appid":"asd","message":"OH HELLO THERE","title":"wup","priority":2,"date":"2017-01-02T00:00:00Z"},{"id":2,"appid":"cloud","message":"hi","title":"hi","priority":4,"date":"2017-01-02T00:00:00Z"}]`, string(bytes))
}

func (s *MessageSuite) Test_GetMessagesWithToken() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "")
	t, _ := time.Parse("2006/01/02", "2021/01/02")
	s.db.On("GetMessagesByUserAndApplication", uint(4), "mytoken").Return([]*model.Message{{ID: 2, ApplicationID: "mytoken", Message: "hi", Title: "hi", Date: t, Priority: 4}})
	s.ctx.Params = gin.Params{{Key: "appid", Value: "mytoken"}}

	s.a.GetMessagesWithApplication(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	bytes, _ := ioutil.ReadAll(s.recorder.Body)
	assert.JSONEq(s.T(), `[{"id":2,"appid":"mytoken","message":"hi","title":"hi","priority":4,"date":"2021-01-02T00:00:00Z"}]`, string(bytes))
}

func (s *MessageSuite) Test_DeleteMessage_invalidID() {
	s.ctx.Params = gin.Params{{Key: "id", Value: "string"}}

	s.a.DeleteMessage(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessage_notExistingID() {
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	s.db.On("GetMessageByID", uint(1)).Return(nil)

	s.a.DeleteMessage(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessage_existingIDButNotOwner() {
	auth.RegisterAuthentication(s.ctx, nil, 6, "")
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	s.db.On("GetMessageByID", uint(1)).Return(&model.Message{ID: 1, ApplicationID: "token"})
	s.db.On("GetApplicationByID", "token").Return(&model.Application{ID: "token", UserID: 2})

	s.a.DeleteMessage(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessage() {
	auth.RegisterAuthentication(s.ctx, nil, 2, "")
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	s.db.On("GetMessageByID", uint(1)).Return(&model.Message{ID: 1, ApplicationID: "token"})
	s.db.On("GetApplicationByID", "token").Return(&model.Application{ID: "token", UserID: 2})
	s.db.On("DeleteMessageByID", uint(1)).Return(nil)

	s.a.DeleteMessage(s.ctx)

	s.db.AssertCalled(s.T(), "DeleteMessageByID", uint(1))
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessageWithToken() {
	auth.RegisterAuthentication(s.ctx, nil, 2, "")
	s.ctx.Params = gin.Params{{Key: "appid", Value: "mytoken"}}
	s.db.On("GetApplicationByID", "mytoken").Return(&model.Application{ID: "mytoken", UserID: 2})
	s.db.On("DeleteMessagesByApplication", "mytoken").Return(nil)

	s.a.DeleteMessageWithApplication(s.ctx)

	s.db.AssertCalled(s.T(), "DeleteMessagesByApplication", "mytoken")
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessageWithToken_notExistingToken() {
	auth.RegisterAuthentication(s.ctx, nil, 2, "")
	s.ctx.Params = gin.Params{{Key: "appid", Value: "asdasdasd"}}
	s.db.On("GetApplicationByID", "asdasdasd").Return(nil)
	s.db.On("DeleteMessagesByApplication", "asdasdasd").Return(nil)

	s.a.DeleteMessageWithApplication(s.ctx)

	s.db.AssertNotCalled(s.T(), "DeleteMessagesByApplication", "mytoken")
	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessageWithToken_notOwner() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "")
	s.ctx.Params = gin.Params{{Key: "appid", Value: "mytoken"}}
	s.db.On("GetApplicationByID", "mytoken").Return(&model.Application{ID: "mytoken", UserID: 2})
	s.db.On("DeleteMessagesByApplication", "mytoken").Return(nil)

	s.a.DeleteMessageWithApplication(s.ctx)

	s.db.AssertNotCalled(s.T(), "DeleteMessagesByApplication", "mytoken")
	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *MessageSuite) Test_DeleteMessages() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "")
	s.db.On("DeleteMessagesByUser", uint(4)).Return(nil)

	s.a.DeleteMessages(s.ctx)

	s.db.AssertCalled(s.T(), "DeleteMessagesByUser", uint(4))
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *MessageSuite) Test_CreateMessage_onJson_allParams() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")

	t, _ := time.Parse("2006/01/02", "2017/01/02")
	patch := monkey.Patch(time.Now, func() time.Time { return t })
	defer patch.Unpatch()
	expected := &model.Message{ID: 0, ApplicationID: "app-token", Title: "mytitle", Message: "mymessage", Priority: 1, Date: t}

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"title": "mytitle", "message": "mymessage", "priority": 1}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.db.On("CreateMessage", expected).Return(nil)

	s.a.CreateMessage(s.ctx)

	s.db.AssertCalled(s.T(), "CreateMessage", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *MessageSuite) Test_CreateMessage_onlyRequired() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")

	t, _ := time.Parse("2006/01/02", "2017/01/02")
	patch := monkey.Patch(time.Now, func() time.Time { return t })
	defer patch.Unpatch()
	expected := &model.Message{ID: 0, ApplicationID: "app-token", Title: "mytitle", Message: "mymessage", Date: t}

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"title": "mytitle", "message": "mymessage"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.db.On("CreateMessage", expected).Return(nil)

	s.a.CreateMessage(s.ctx)

	s.db.AssertCalled(s.T(), "CreateMessage", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *MessageSuite) Test_CreateMessage_failWhenNoMessage() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"title": "mytitle"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateMessage(s.ctx)

	s.db.AssertNotCalled(s.T(), "CreateMessage", mock.Anything)
	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *MessageSuite) Test_CreateMessage_failWhenNoTitle() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"message": "mymessage"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateMessage(s.ctx)

	s.db.AssertNotCalled(s.T(), "CreateMessage", mock.Anything)
	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *MessageSuite) Test_CreateMessage_failWhenPriorityNotNumber() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(`{"title": "mytitle", "message": "mymessage", "priority": "asd"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateMessage(s.ctx)

	s.db.AssertNotCalled(s.T(), "CreateMessage", mock.Anything)
	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *MessageSuite) Test_CreateMessage_onQueryData() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")

	t, _ := time.Parse("2006/01/02", "2017/01/02")
	patch := monkey.Patch(time.Now, func() time.Time { return t })
	defer patch.Unpatch()
	expected := &model.Message{ID: 0, ApplicationID: "app-token", Title: "mytitle", Message: "mymessage", Priority: 1, Date: t}

	s.ctx.Request = httptest.NewRequest("POST", "/token?title=mytitle&message=mymessage&priority=1", nil)
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.db.On("CreateMessage", expected).Return(nil)

	s.a.CreateMessage(s.ctx)

	s.db.AssertCalled(s.T(), "CreateMessage", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *MessageSuite) Test_CreateMessage_onFormData() {
	auth.RegisterAuthentication(s.ctx, nil, 4, "app-token")

	t, _ := time.Parse("2006/01/02", "2017/01/02")
	patch := monkey.Patch(time.Now, func() time.Time { return t })
	defer patch.Unpatch()
	expected := &model.Message{ID: 0, ApplicationID: "app-token", Title: "mytitle", Message: "mymessage", Priority: 1, Date: t}

	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader("title=mytitle&message=mymessage&priority=1"))
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.db.On("CreateMessage", expected).Return(nil)

	s.a.CreateMessage(s.ctx)

	s.db.AssertCalled(s.T(), "CreateMessage", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}
