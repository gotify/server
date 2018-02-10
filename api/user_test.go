package api

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	"github.com/gin-gonic/gin"
	apimock "github.com/jmattheis/memo/api/mock"
	"github.com/jmattheis/memo/auth"
	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	adminUser      = &model.User{ID: 1, Name: "jmattheis", Pass: []byte{1, 2}, Admin: true}
	adminUserJSON  = `{"id":1,"name":"jmattheis","admin":true}`
	normalUser     = &model.User{ID: 2, Name: "nicories", Pass: []byte{2, 3}, Admin: false}
	normalUserJSON = `{"id":2,"name":"nicories","admin":false}`
)

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

type UserSuite struct {
	suite.Suite
	db       *apimock.MockUserDatabase
	a        *UserAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func (s *UserSuite) BeforeTest(suiteName, testName string) {
	gin.SetMode(gin.TestMode)
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.db = &apimock.MockUserDatabase{}
	s.a = &UserAPI{DB: s.db}
}

func (s *UserSuite) Test_GetUsers() {
	s.db.On("GetUsers").Return([]*model.User{adminUser, normalUser})

	s.a.GetUsers(s.ctx)

	s.expectJSON(fmt.Sprintf("[%s, %s]", adminUserJSON, normalUserJSON))
}

func (s *UserSuite) Test_GetCurrentUser() {
	patch := monkey.Patch(auth.GetUserID, func(*gin.Context) uint { return 1 })
	defer patch.Unpatch()
	s.db.On("GetUserByID", uint(1)).Return(adminUser)

	s.a.GetCurrentUser(s.ctx)

	s.expectJSON(adminUserJSON)
}

func (s *UserSuite) Test_GetUserByID() {
	s.db.On("GetUserByID", uint(2)).Return(normalUser)
	s.ctx.Params = gin.Params{{Key: "id", Value: "2"}}

	s.a.GetUserByID(s.ctx)

	s.expectJSON(normalUserJSON)
}

func (s *UserSuite) Test_GetUserByID_InvalidID() {
	s.db.On("GetUserByID", uint(2)).Return(normalUser)
	s.ctx.Params = gin.Params{{Key: "id", Value: "abc"}}

	s.a.GetUserByID(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_GetUserByID_UnknownUser() {
	s.db.On("GetUserByID", mock.Anything).Return(nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "3"}}

	s.a.GetUserByID(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *UserSuite) Test_DeleteUserByID_InvalidID() {
	s.ctx.Params = gin.Params{{Key: "id", Value: "abc"}}

	s.a.DeleteUserByID(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_DeleteUserByID_UnknownUser() {
	s.db.On("GetUserByID", mock.Anything).Return(nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "3"}}

	s.a.DeleteUserByID(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *UserSuite) Test_DeleteUserByID() {
	s.db.On("GetUserByID", uint(2)).Return(normalUser)
	s.db.On("DeleteUserByID", uint(2)).Return(nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "2"}}

	s.a.DeleteUserByID(s.ctx)

	s.db.AssertCalled(s.T(), "DeleteUserByID", uint(2))
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *UserSuite) Test_CreateUser() {
	pwByte := []byte{1, 2, 3}
	patch := monkey.Patch(auth.CreatePassword, func(pw string) []byte {
		if pw == "mylittlepony" {
			return pwByte
		}
		return []byte{5, 67}
	})
	defer patch.Unpatch()

	s.db.On("GetUserByName", "tom").Return(nil)
	s.db.On("CreateUser", mock.Anything).Return(nil)

	s.ctx.Request = httptest.NewRequest("POST", "/user", strings.NewReader(`{"name": "tom", "pass": "mylittlepony", "admin": true}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateUser(s.ctx)

	s.db.AssertCalled(s.T(), "CreateUser", &model.User{Name: "tom", Pass: pwByte, Admin: true})
}

func (s *UserSuite) Test_CreateUser_NoPassword() {
	s.ctx.Request = httptest.NewRequest("POST", "/user", strings.NewReader(`{"name": "tom", "pass": "", "admin": true}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateUser(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_CreateUser_NoName() {
	s.ctx.Request = httptest.NewRequest("POST", "/user", strings.NewReader(`{"name": "", "pass": "asd", "admin": true}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateUser(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_CreateUser_NameAlreadyExists() {
	pwByte := []byte{1, 2, 3}
	monkey.Patch(auth.CreatePassword, func(pw string) []byte { return pwByte })

	s.db.On("GetUserByName", "tom").Return(&model.User{ID: 3, Name: "tom"})

	s.ctx.Request = httptest.NewRequest("POST", "/user", strings.NewReader(`{"name": "tom", "pass": "mylittlepony", "admin": true}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.CreateUser(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_UpdateUserByID_InvalidID() {
	s.ctx.Params = gin.Params{{Key: "id", Value: "abc"}}

	s.ctx.Request = httptest.NewRequest("POST", "/user/abc", strings.NewReader(`{"name": "tom", "pass": "", "admin": false}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.UpdateUserByID(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_UpdateUserByID_UnknownUser() {
	s.db.On("GetUserByID", uint(2)).Return(nil)

	s.ctx.Params = gin.Params{{Key: "id", Value: "2"}}

	s.ctx.Request = httptest.NewRequest("POST", "/user/2", strings.NewReader(`{"name": "tom", "pass": "", "admin": false}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.UpdateUserByID(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *UserSuite) Test_UpdateUserByID_UpdateNotPassword() {
	s.db.On("GetUserByID", uint(2)).Return(&model.User{Name: "nico", Pass: []byte{5}, Admin: false})
	expected := &model.User{ID: 2, Name: "tom", Pass: []byte{5}, Admin: true}

	s.db.On("UpdateUser", expected).Return(nil)

	s.ctx.Params = gin.Params{{Key: "id", Value: "2"}}

	s.ctx.Request = httptest.NewRequest("POST", "/user/2", strings.NewReader(`{"name": "tom", "pass": "", "admin": true}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.UpdateUserByID(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertCalled(s.T(), "UpdateUser", expected)
}

func (s *UserSuite) Test_UpdateUserByID_UpdatePassword() {
	pwByte := []byte{1, 2, 3}
	patch := monkey.Patch(auth.CreatePassword, func(pw string) []byte { return pwByte })
	defer patch.Unpatch()

	s.db.On("GetUserByID", uint(2)).Return(normalUser)
	expected := &model.User{ID: 2, Name: "tom", Pass: pwByte, Admin: true}

	s.db.On("UpdateUser", expected).Return(nil)

	s.ctx.Params = gin.Params{{Key: "id", Value: "2"}}

	s.ctx.Request = httptest.NewRequest("POST", "/user/2", strings.NewReader(`{"name": "tom", "pass": "secret", "admin": true}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.UpdateUserByID(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertCalled(s.T(), "UpdateUser", expected)
}

func (s *UserSuite) Test_UpdatePassword() {
	pwByte := []byte{1, 2, 3}
	createPasswordPatch := monkey.Patch(auth.CreatePassword, func(pw string) []byte { return pwByte })
	defer createPasswordPatch.Unpatch()
	patchUser := monkey.Patch(auth.GetUserID, func(*gin.Context) uint { return 1 })
	defer patchUser.Unpatch()
	s.ctx.Request = httptest.NewRequest("POST", "/user/current/password", strings.NewReader(`{"pass": "secret"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.db.On("GetUserByID", uint(1)).Return(&model.User{ID: 1, Name: "jmattheis", Pass: []byte{1}})
	s.db.On("UpdateUser", mock.Anything).Return(nil)

	s.a.ChangePassword(s.ctx)

	s.db.AssertCalled(s.T(), "UpdateUser", &model.User{ID: 1, Name: "jmattheis", Pass: pwByte})
}

func (s *UserSuite) Test_UpdatePassword_EmptyPassword() {
	patchUser := monkey.Patch(auth.GetUserID, func(*gin.Context) uint { return 1 })
	defer patchUser.Unpatch()
	s.ctx.Request = httptest.NewRequest("POST", "/user/current/password", strings.NewReader(`{"pass":""}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.db.On("UpdateUser", mock.Anything).Return(nil)

	s.db.On("GetUserByID", uint(1)).Return(&model.User{ID: 1, Name: "jmattheis", Pass: []byte{1}})

	s.a.ChangePassword(s.ctx)

	s.db.AssertNotCalled(s.T(), "UpdateUser", mock.Anything)
	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) expectJSON(json string) {
	assert.Equal(s.T(), 200, s.recorder.Code)
	bytes, _ := ioutil.ReadAll(s.recorder.Body)

	assert.JSONEq(s.T(), json, string(bytes))
}
