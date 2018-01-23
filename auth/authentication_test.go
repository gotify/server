package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationSuite))
}

type AuthenticationSuite struct {
	suite.Suite
	auth *Auth
}

func (s *AuthenticationSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	s.auth = &Auth{&DBMock{}}
}

func (s *AuthenticationSuite) TestQueryToken() {
	s.assertQueryRequest("token", "ergerogerg", s.auth.RequireWrite, 401)
	s.assertQueryRequest("token", "ergerogerg", s.auth.RequireAll, 401)
	s.assertQueryRequest("token", "ergerogerg", s.auth.RequireAdmin, 401)

	s.assertQueryRequest("tokenx", "all", s.auth.RequireWrite, 401)
	s.assertQueryRequest("tokenx", "all", s.auth.RequireAll, 401)
	s.assertQueryRequest("tokenx", "all", s.auth.RequireAdmin, 401)

	s.assertQueryRequest("token", "writeonly", s.auth.RequireWrite, 200)
	s.assertQueryRequest("token", "writeonly", s.auth.RequireAll, 401)
	s.assertQueryRequest("token", "writeonly", s.auth.RequireAdmin, 401)

	s.assertQueryRequest("token", "all", s.auth.RequireWrite, 200)
	s.assertQueryRequest("token", "all", s.auth.RequireAll, 200)
	s.assertQueryRequest("token", "all", s.auth.RequireAdmin, 401)

	s.assertQueryRequest("token", "admin", s.auth.RequireWrite, 200)
	s.assertQueryRequest("token", "admin", s.auth.RequireAll, 200)
	s.assertQueryRequest("token", "admin", s.auth.RequireAdmin, 200)
}

func (s *AuthenticationSuite) assertQueryRequest(key, value string, f fMiddleware, code int) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", fmt.Sprintf("/?%s=%s", key, value), nil)
	f()(ctx)
	assert.Equal(s.T(), code, recorder.Code)
}

func (s *AuthenticationSuite) TestNothingProvided() {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/", nil)
	s.auth.RequireWrite()(ctx)
	assert.Equal(s.T(), 401, recorder.Code)
}

func (s *AuthenticationSuite) TestHeaderApiKeyToken() {
	s.assertHeaderRequest("Authorization", "ergerogerg", s.auth.RequireWrite, 401)
	s.assertHeaderRequest("Authorization", "ergerogerg", s.auth.RequireAll, 401)
	s.assertHeaderRequest("Authorization", "ergerogerg", s.auth.RequireAdmin, 401)

	s.assertHeaderRequest("Authorization", "ApiKey ergerogerg", s.auth.RequireWrite, 401)
	s.assertHeaderRequest("Authorization", "ApiKey ergerogerg", s.auth.RequireAll, 401)
	s.assertHeaderRequest("Authorization", "ApiKey ergerogerg", s.auth.RequireAdmin, 401)

	s.assertHeaderRequest("Authorizationx", "ApiKey all", s.auth.RequireWrite, 401)
	s.assertHeaderRequest("Authorizationx", "ApiKey all", s.auth.RequireAll, 401)
	s.assertHeaderRequest("Authorizationx", "ApiKey all", s.auth.RequireAdmin, 401)

	s.assertHeaderRequest("Authorization", "ApiKey writeonly", s.auth.RequireWrite, 200)
	s.assertHeaderRequest("Authorization", "ApiKey writeonly", s.auth.RequireAll, 401)
	s.assertHeaderRequest("Authorization", "ApiKey writeonly", s.auth.RequireAdmin, 401)

	s.assertHeaderRequest("Authorization", "ApiKey all", s.auth.RequireWrite, 200)
	s.assertHeaderRequest("Authorization", "ApiKey all", s.auth.RequireAll, 200)
	s.assertHeaderRequest("Authorization", "ApiKey all", s.auth.RequireAdmin, 401)

	s.assertHeaderRequest("Authorization", "ApiKey admin", s.auth.RequireWrite, 200)
	s.assertHeaderRequest("Authorization", "ApiKey admin", s.auth.RequireAll, 200)
	s.assertHeaderRequest("Authorization", "ApiKey admin", s.auth.RequireAdmin, 200)
}

func (s *AuthenticationSuite) TestBasicAuth() {
	s.assertHeaderRequest("Authorization", "Basic ergerogerg", s.auth.RequireWrite, 401)
	s.assertHeaderRequest("Authorization", "Basic ergerogerg", s.auth.RequireAll, 401)
	s.assertHeaderRequest("Authorization", "Basic ergerogerg", s.auth.RequireAdmin, 401)

	// user existing:pw
	s.assertHeaderRequest("Authorization", "Basic ZXhpc3Rpbmc6cHc=", s.auth.RequireWrite, 200)
	s.assertHeaderRequest("Authorization", "Basic ZXhpc3Rpbmc6cHc=", s.auth.RequireAll, 200)
	s.assertHeaderRequest("Authorization", "Basic ZXhpc3Rpbmc6cHc=", s.auth.RequireAdmin, 401)

	// user admin:pw
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHc=", s.auth.RequireWrite, 200)
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHc=", s.auth.RequireAll, 200)
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHc=", s.auth.RequireAdmin, 200)

	// user admin:pwx
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHd4", s.auth.RequireWrite, 401)
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHd4", s.auth.RequireAll, 401)
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHd4", s.auth.RequireAdmin, 401)
}

func (s *AuthenticationSuite) assertHeaderRequest(key, value string, f fMiddleware, code int) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/", nil)
	ctx.Request.Header.Set(key, value)
	f()(ctx)
	assert.Equal(s.T(), code, recorder.Code)
}

type fMiddleware func() gin.HandlerFunc
type DBMock struct{}

func (d *DBMock) GetTokenById(id string) *model.Token {
	if id == "writeonly" {
		return &model.Token{Id: "valid", WriteOnly: true, UserId: 1}
	}
	if id == "all" {
		return &model.Token{Id: "valid", WriteOnly: false, UserId: 1}
	}
	if id == "admin" {
		return &model.Token{Id: "valid", WriteOnly: false, UserId: 2}
	}
	return nil
}

func (d *DBMock) GetUserByName(name string) *model.User {
	if name == "existing" {
		return &model.User{Name: "existing", Pass: CreatePassword("pw")}
	}
	if name == "admin" {
		return &model.User{Name: "admin", Pass: CreatePassword("pw"), Admin: true}
	}
	return nil
}
func (d *DBMock) GetUserById(id uint) *model.User {
	if id == 1 {
		return &model.User{Name: "existing", Pass: CreatePassword("pw"), Admin: false}
	}

	if id == 2 {
		return &model.User{Name: "existing", Pass: CreatePassword("pw"), Admin: true}
	}
	return nil
}
