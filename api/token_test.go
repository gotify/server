package api

import (
	"math/rand"
	"net/http/httptest"
	"testing"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
	"github.com/gotify/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	firstApplicationToken  = "APorrUa5b1IIK3y"
	secondApplicationToken = "AKo_Pp6ww_9vZal"
	firstClientToken       = "CPorrUa5b1IIK3y"
	secondClientToken      = "CKo_Pp6ww_9vZal"
)

func TestTokenSuite(t *testing.T) {
	suite.Run(t, new(TokenSuite))
}

type TokenSuite struct {
	suite.Suite
	db       *test.Database
	a        *TokenAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func (s *TokenSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	rand.Seed(50)
	s.recorder = httptest.NewRecorder()
	s.db = test.NewDB(s.T())
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.a = &TokenAPI{DB: s.db}
}

func (s *TokenSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

// test application api

func (s *TokenSuite) Test_CreateApplication_mapAllParameters() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name&description=description_text")
	s.a.CreateApplication(s.ctx)

	expected := &model.Application{ID: 1, Token: firstApplicationToken, UserID: 5, Name: "custom_name", Description: "description_text"}
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Equal(s.T(), expected, s.db.GetApplicationByID(1))
}

func (s *TokenSuite) Test_CreateApplication_expectBadRequestOnEmptyName() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=&description=description_text")
	s.a.CreateApplication(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
	assert.Empty(s.T(), s.db.GetApplicationsByUser(5))
}

func (s *TokenSuite) Test_DeleteApplication_expectNotFoundOnCurrentUserIsNotOwner() {
	s.db.User(2)
	s.db.User(5).App(5)

	test.WithUser(s.ctx, 2)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/5", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "5"}}

	s.a.DeleteApplication(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
	s.db.AssertAppExist(5)
}

func (s *TokenSuite) Test_CreateApplication_onlyRequiredParameters() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")
	s.a.CreateApplication(s.ctx)

	expected := &model.Application{ID: 1, Token: firstApplicationToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.db.GetApplicationsByUser(5), expected)
}
func (s *TokenSuite) Test_ensureApplicationHasCorrectJsonRepresentation() {
	actual := &model.Application{ID: 1, UserID: 2, Token: "Aasdasfgeeg", Name: "myapp", Description: "mydesc"}
	test.JSONEquals(s.T(), actual, `{"id":1,"token":"Aasdasfgeeg","name":"myapp","description":"mydesc"}`)
}

func (s *TokenSuite) Test_CreateApplication_returnsApplicationWithID() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateApplication(s.ctx)

	expected := &model.Application{ID: 1, Token: firstApplicationToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), expected, s.recorder)
}

func (s *TokenSuite) Test_CreateApplication_withExistingToken() {
	s.db.User(5)
	s.db.User(6).AppWithToken(1, firstApplicationToken)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateApplication(s.ctx)

	expected := &model.Application{ID: 2, Token: secondApplicationToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.db.GetApplicationsByUser(5), expected)
}

func (s *TokenSuite) Test_GetApplications() {
	userBuilder := s.db.User(5)
	first := userBuilder.NewAppWithToken(1, "perfper")
	second := userBuilder.NewAppWithToken(2, "asdasd")

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("GET", "/tokens", nil)

	s.a.GetApplications(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), []*model.Application{first, second}, s.recorder)
}

func (s *TokenSuite) Test_DeleteApplication_expectNotFound() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstApplicationToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "4"}}

	s.a.DeleteApplication(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *TokenSuite) Test_DeleteApplication() {
	s.db.User(5).App(1)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstApplicationToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	s.a.DeleteApplication(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertAppNotExist(1)
}

// test client api

func (s *TokenSuite) Test_ensureClientHasCorrectJsonRepresentation() {
	actual := &model.Client{ID: 1, UserID: 2, Token: "Casdasfgeeg", Name: "myclient"}
	test.JSONEquals(s.T(), actual, `{"id":1,"token":"Casdasfgeeg","name":"myclient"}`)
}

func (s *TokenSuite) Test_CreateClient_mapAllParameters() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name&description=description_text")

	s.a.CreateClient(s.ctx)

	expected := &model.Client{ID: 1, Token: firstClientToken, UserID: 5, Name: "custom_name"}
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.db.GetClientsByUser(5), expected)
}

func (s *TokenSuite) Test_CreateClient_expectBadRequestOnEmptyName() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=&description=description_text")

	s.a.CreateClient(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
	assert.Empty(s.T(), s.db.GetClientsByUser(5))
}

func (s *TokenSuite) Test_DeleteClient_expectNotFoundOnCurrentUserIsNotOwner() {
	s.db.User(5).Client(7)
	s.db.User(2)

	test.WithUser(s.ctx, 2)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/7", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "7"}}

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
	s.db.AssertClientExist(7)
}

func (s *TokenSuite) Test_CreateClient_returnsClientWithID() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateClient(s.ctx)

	expected := &model.Client{ID: 1, Token: firstClientToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), expected, s.recorder)
}

func (s *TokenSuite) Test_CreateClient_withExistingToken() {
	s.db.User(5).ClientWithToken(1, firstClientToken)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateClient(s.ctx)

	expected := &model.Client{ID: 2, Token: secondClientToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), expected, s.recorder)
}

func (s *TokenSuite) Test_GetClients() {
	userBuilder := s.db.User(5)
	first := userBuilder.NewClientWithToken(1, "perfper")
	second := userBuilder.NewClientWithToken(2, "asdasd")

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("GET", "/tokens", nil)

	s.a.GetClients(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), []*model.Client{first, second}, s.recorder)
}

func (s *TokenSuite) Test_DeleteClient_expectNotFound() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstClientToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "8"}}

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

//
func (s *TokenSuite) Test_DeleteClient() {
	s.db.User(5).Client(8)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstClientToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "8"}}

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertClientNotExist(8)
}

func (s *TokenSuite) withFormData(formData string) {
	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(formData))
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}
