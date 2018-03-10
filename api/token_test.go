package api

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	apimock "github.com/gotify/server/api/mock"
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	db       *apimock.MockTokenDatabase
	a        *TokenAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func (s *TokenSuite) BeforeTest(suiteName, testName string) {
	gin.SetMode(gin.TestMode)
	rand.Seed(50)
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.db = &apimock.MockTokenDatabase{}
	s.a = &TokenAPI{DB: s.db}
}

// test application api

func (s *TokenSuite) Test_CreateApplication_mapAllParameters() {
	expected := &model.Application{Token: firstApplicationToken, UserID: 5, Name: "custom_name", Description: "description_text"}

	s.ctx.Set("user", &model.User{ID: 5})
	s.withFormData("name=custom_name&description=description_text")

	s.db.On("GetApplicationByToken", firstApplicationToken).Return(nil)
	s.db.On("CreateApplication", expected).Return(nil)

	s.a.CreateApplication(s.ctx)

	s.db.AssertCalled(s.T(), "CreateApplication", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) Test_CreateApplication_expectBadRequestOnEmptyName() {
	s.ctx.Set("user", &model.User{ID: 5})
	s.withFormData("name=&description=description_text")

	s.a.CreateApplication(s.ctx)

	s.db.AssertNotCalled(s.T(), "CreateApplication", mock.Anything)
	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *TokenSuite) Test_DeleteApplication_expectNotFoundOnCurrentUserIsNotOwner() {
	s.ctx.Set("user", &model.User{ID: 2})
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/5", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "5"}}

	s.db.On("GetApplicationByID", uint(5)).Return(&model.Application{ID: 5, Token: firstApplicationToken, UserID: 5})

	s.a.DeleteApplication(s.ctx)

	s.db.AssertNotCalled(s.T(), "DeleteApplicationByID", mock.Anything)
	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *TokenSuite) Test_CreateApplication_onlyRequiredParameters() {
	expected := &model.Application{Token: firstApplicationToken, Name: "custom_name", UserID: 5}

	s.ctx.Set("user", &model.User{ID: 5})
	s.withFormData("name=custom_name")

	s.db.On("GetApplicationByToken", firstApplicationToken).Return(nil)
	s.db.On("CreateApplication", expected).Return(nil)

	s.a.CreateApplication(s.ctx)

	s.db.AssertCalled(s.T(), "CreateApplication", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) Test_CreateApplication_returnsApplicationWithID() {
	expected := &model.Application{Token: firstApplicationToken, Name: "custom_name", UserID: 5}

	s.ctx.Set("user", &model.User{ID: 5})
	s.withFormData("name=custom_name")

	s.db.On("GetApplicationByToken", firstApplicationToken).Return(nil)
	s.db.On("CreateApplication", expected).Return(nil)

	s.a.CreateApplication(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	bytes, _ := ioutil.ReadAll(s.recorder.Body)

	assert.Equal(s.T(), `{"id":0,"token":"APorrUa5b1IIK3y","name":"custom_name","description":""}`, string(bytes))
}

func (s *TokenSuite) Test_CreateApplication_withExistingToken() {
	expected := &model.Application{Token: secondApplicationToken, Name: "custom_name", UserID: 5}

	s.ctx.Set("user", &model.User{ID: 5})
	s.withFormData("name=custom_name")

	s.db.On("GetApplicationByToken", firstApplicationToken).Return(&model.Application{Token: firstApplicationToken})
	s.db.On("GetApplicationByToken", secondApplicationToken).Return(nil)
	s.db.On("CreateApplication", expected).Return(nil)

	s.a.CreateApplication(s.ctx)

	s.db.AssertCalled(s.T(), "CreateApplication", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) Test_GetApplications() {
	s.ctx.Set("user", &model.User{ID: 5})
	s.ctx.Request = httptest.NewRequest("GET", "/tokens", nil)

	s.db.On("GetApplicationsByUser", uint(5)).Return([]*model.Application{
		{Token: "perfper", Name: "first", Description: "desc"},
		{Token: "asdasd", Name: "second", Description: "desc2"},
	})
	s.a.GetApplications(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	bytes, _ := ioutil.ReadAll(s.recorder.Body)

	assert.Equal(s.T(), `[{"id":0,"token":"perfper","name":"first","description":"desc"},{"id":0,"token":"asdasd","name":"second","description":"desc2"}]`, string(bytes))
}

func (s *TokenSuite) Test_DeleteApplication_expectNotFound() {
	s.ctx.Set("user", &model.User{ID: 5})
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstApplicationToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "4"}}

	s.db.On("DeleteApplicationByID", uint(4)).Return(errors.New("what? that does not exist"))
	s.db.On("GetApplicationByID", uint(4)).Return(nil)

	s.a.DeleteApplication(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *TokenSuite) Test_DeleteApplication() {
	s.ctx.Set("user", &model.User{ID: 5})
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstApplicationToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	s.db.On("DeleteApplicationByID", uint(1)).Return(nil)
	s.db.On("GetApplicationByID", uint(1)).Return(&model.Application{Token: firstApplicationToken, Name: "custom_name", UserID: 5})

	s.a.DeleteApplication(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
}

// test client api

func (s *TokenSuite) Test_CreateClient_mapAllParameters() {
	expected := &model.Client{Token: firstClientToken, UserID: 5, Name: "custom_name"}

	s.ctx.Set("user", &model.User{ID: 5})
	s.withFormData("name=custom_name&description=description_text")

	s.db.On("GetClientByToken", firstClientToken).Return(nil)
	s.db.On("CreateClient", expected).Return(nil)

	s.a.CreateClient(s.ctx)

	s.db.AssertCalled(s.T(), "CreateClient", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) Test_CreateClient_expectBadRequestOnEmptyName() {
	s.ctx.Set("user", &model.User{ID: 5})
	s.withFormData("name=&description=description_text")

	s.a.CreateClient(s.ctx)

	s.db.AssertNotCalled(s.T(), "CreateClient", mock.Anything)
	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *TokenSuite) Test_DeleteClient_expectNotFoundOnCurrentUserIsNotOwner() {
	s.ctx.Set("user", &model.User{ID: 2})
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/7", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "7"}}

	s.db.On("GetClientByID", uint(7)).Return(&model.Client{Token: firstClientToken, UserID: 5})

	s.a.DeleteClient(s.ctx)

	s.db.AssertNotCalled(s.T(), "DeleteClientByID", mock.Anything)
	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *TokenSuite) Test_CreateClient_returnsClientWithID() {
	expected := &model.Client{Token: firstClientToken, Name: "custom_name", UserID: 5}

	s.ctx.Set("user", &model.User{ID: 5})
	s.withFormData("name=custom_name")

	s.db.On("GetClientByToken", firstClientToken).Return(nil)
	s.db.On("CreateClient", expected).Return(nil)

	s.a.CreateClient(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	bytes, _ := ioutil.ReadAll(s.recorder.Body)

	assert.Equal(s.T(), `{"id":0,"token":"CPorrUa5b1IIK3y","name":"custom_name"}`, string(bytes))
}

func (s *TokenSuite) Test_CreateClient_withExistingToken() {
	expected := &model.Client{Token: secondClientToken, Name: "custom_name", UserID: 5}

	s.ctx.Set("user", &model.User{ID: 5})
	s.withFormData("name=custom_name")

	s.db.On("GetClientByToken", firstClientToken).Return(&model.Client{Token: firstClientToken})
	s.db.On("GetClientByToken", secondClientToken).Return(nil)
	s.db.On("CreateClient", expected).Return(nil)

	s.a.CreateClient(s.ctx)

	s.db.AssertCalled(s.T(), "CreateClient", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) Test_GetClients() {
	s.ctx.Set("user", &model.User{ID: 5})
	s.ctx.Request = httptest.NewRequest("GET", "/tokens", nil)

	s.db.On("GetClientsByUser", uint(5)).Return([]*model.Client{
		{Token: "perfper", Name: "first"},
		{Token: "asdasd", Name: "second"},
	})
	s.a.GetClients(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	bytes, _ := ioutil.ReadAll(s.recorder.Body)

	assert.Equal(s.T(), `[{"id":0,"token":"perfper","name":"first"},{"id":0,"token":"asdasd","name":"second"}]`, string(bytes))
}

func (s *TokenSuite) Test_DeleteClient_expectNotFound() {
	s.ctx.Set("user", &model.User{ID: 5})
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstClientToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "8"}}

	s.db.On("DeleteClientByID", uint(8)).Return(errors.New("what? that does not exist"))
	s.db.On("GetClientByID", uint(8)).Return(nil)

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *TokenSuite) Test_DeleteClient() {
	s.ctx.Set("user", &model.User{ID: 5})
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstClientToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "8"}}

	s.db.On("DeleteClientByID", uint(8)).Return(nil)
	s.db.On("GetClientByID", uint(8)).Return(&model.Client{Token: firstClientToken, Name: "custom_name", UserID: 5})

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) withFormData(formData string) {
	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(formData))
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}
