package api

import (
	"github.com/stretchr/testify/suite"
	"github.com/gin-gonic/gin"
	"testing"
	apimock "github.com/jmattheis/memo/api/mock"
	"math/rand"
	"net/http/httptest"
	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
	"strings"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"errors"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(TokenSuite))
}

type TokenSuite struct {
	suite.Suite
	db       *apimock.MockTokenDatabase
	a        *TokenApi
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func (s *TokenSuite) BeforeTest(suiteName, testName string) {
	gin.SetMode(gin.TestMode)
	rand.Seed(50)
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.db = &apimock.MockTokenDatabase{}
	s.a = &TokenApi{DB: s.db}
}

func (s *TokenSuite) Test_mapAllParameters() {
	expected := &model.Token{Id: "PorrUa5b1IIK3yK", Name: "custom_name", UserId: 5, WriteOnly: true, Description: "description_text"}

	s.ctx.Set("user", &model.User{Id: 5})
	s.withFormData("name=custom_name&writeOnly=true&description=description_text")

	s.db.On("GetTokenById", "PorrUa5b1IIK3yK").Return(nil)
	s.db.On("CreateToken", expected).Return(nil)

	s.a.CreateToken(s.ctx)

	s.db.AssertCalled(s.T(), "CreateToken", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) Test_badRequest_emptyName() {
	s.ctx.Set("user", &model.User{Id: 5})
	s.withFormData("name=&writeOnly=true&description=description_text")

	s.a.CreateToken(s.ctx)

	s.db.AssertNotCalled(s.T(), "CreateToken", mock.Anything)
	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *TokenSuite) Test_success_withOnlyRequiredProperties() {
	expected := &model.Token{Id: "PorrUa5b1IIK3yK", Name: "custom_name", UserId: 5}

	s.ctx.Set("user", &model.User{Id: 5})
	s.withFormData("name=custom_name")

	s.db.On("GetTokenById", "PorrUa5b1IIK3yK").Return(nil)
	s.db.On("CreateToken", expected).Return(nil)

	s.a.CreateToken(s.ctx)

	s.db.AssertCalled(s.T(), "CreateToken", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) Test_success_withExistingToken() {
	expected := &model.Token{Id: "o_Pp6ww_9vZal6-", Name: "custom_name", UserId: 5}

	s.ctx.Set("user", &model.User{Id: 5})
	s.withFormData("name=custom_name")

	s.db.On("GetTokenById", "PorrUa5b1IIK3yK").Return(&model.Token{Id: "PorrUa5b1IIK3yK"})
	s.db.On("GetTokenById", "o_Pp6ww_9vZal6-").Return(nil)
	s.db.On("CreateToken", expected).Return(nil)

	s.a.CreateToken(s.ctx)

	s.db.AssertCalled(s.T(), "CreateToken", expected)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) Test_getToken() {
	s.ctx.Set("user", &model.User{Id: 5})
	s.ctx.Request = httptest.NewRequest("GET", "/tokens", nil)

	s.db.On("GetTokensByUser", uint(5)).Return([]*model.Token{
		{Id: "perfper", Name: "first", Description: "desc"},
		{Id: "asdasd", Name: "second", Description: "desc2"},
	});
	s.a.GetTokens(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	bytes, _ := ioutil.ReadAll(s.recorder.Body)

	assert.Equal(s.T(), `[{"Id":"perfper","name":"first","description":"desc","writeOnly":false},{"Id":"asdasd","name":"second","description":"desc2","writeOnly":false}]`, string(bytes))
}

func (s *TokenSuite) Test_deleteToken_fail() {
	s.ctx.Set("user", &model.User{Id: 5})
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/PorrUa5b1IIK3yK", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "PorrUa5b1IIK3yK"}}

	s.db.On("DeleteToken", "PorrUa5b1IIK3yK").Return(errors.New("what? that does not exist"))
	s.db.On("GetTokenById", "PorrUa5b1IIK3yK").Return(nil)

	s.a.DeleteToken(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *TokenSuite) Test_deleteToken_success() {
	s.ctx.Set("user", &model.User{Id: 5})
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/PorrUa5b1IIK3yK", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "PorrUa5b1IIK3yK"}}

	s.db.On("DeleteToken", "PorrUa5b1IIK3yK").Return(nil)
	s.db.On("GetTokenById", "PorrUa5b1IIK3yK").Return(&model.Token{Id: "PorrUa5b1IIK3yK", Name: "custom_name", UserId: 5})

	s.a.DeleteToken(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *TokenSuite) withFormData(formData string) {
	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(formData))
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}
