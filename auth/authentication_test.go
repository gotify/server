// +build !race

package auth

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	authmock "github.com/gotify/server/auth/mock"
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationSuite))
}

type AuthenticationSuite struct {
	suite.Suite
	auth *Auth
	DB   *authmock.MockDatabase
}

func (s *AuthenticationSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	s.DB = &authmock.MockDatabase{}
	s.auth = &Auth{s.DB}
	s.DB.On("GetClientByID", "clienttoken").Return(&model.Client{ID: "clienttoken", UserID: 1, Name: "android phone"})
	s.DB.On("GetClientByID", "clienttoken_admin").Return(&model.Client{ID: "clienttoken", UserID: 2, Name: "android phone2"})
	s.DB.On("GetClientByID", mock.Anything).Return(nil)
	s.DB.On("GetApplicationByID", "apptoken").Return(&model.Application{ID: "apptoken", UserID: 1, Name: "backup server", Description: "irrelevant"})
	s.DB.On("GetApplicationByID", "apptoken_admin").Return(&model.Application{ID: "apptoken", UserID: 2, Name: "backup server", Description: "irrelevant"})
	s.DB.On("GetApplicationByID", mock.Anything).Return(nil)

	s.DB.On("GetUserByID", uint(1)).Return(&model.User{ID: 1, Name: "irrelevant", Admin: false})
	s.DB.On("GetUserByID", uint(2)).Return(&model.User{ID: 2, Name: "irrelevant", Admin: true})

	s.DB.On("GetUserByName", "existing").Return(&model.User{Name: "existing", Pass: CreatePassword("pw")})
	s.DB.On("GetUserByName", "admin").Return(&model.User{Name: "admin", Pass: CreatePassword("pw"), Admin: true})
	s.DB.On("GetUserByName", mock.Anything).Return(nil)
}

func (s *AuthenticationSuite) TestQueryToken() {
	// not existing token
	s.assertQueryRequest("token", "ergerogerg", s.auth.RequireApplicationToken, 401)
	s.assertQueryRequest("token", "ergerogerg", s.auth.RequireClient, 401)
	s.assertQueryRequest("token", "ergerogerg", s.auth.RequireAdmin, 401)

	// not existing key
	s.assertQueryRequest("tokenx", "clienttoken", s.auth.RequireApplicationToken, 401)
	s.assertQueryRequest("tokenx", "clienttoken", s.auth.RequireClient, 401)
	s.assertQueryRequest("tokenx", "clienttoken", s.auth.RequireAdmin, 401)

	// apptoken
	s.assertQueryRequest("token", "apptoken", s.auth.RequireApplicationToken, 200)
	s.assertQueryRequest("token", "apptoken", s.auth.RequireClient, 401)
	s.assertQueryRequest("token", "apptoken", s.auth.RequireAdmin, 401)
	s.assertQueryRequest("token", "apptoken_admin", s.auth.RequireApplicationToken, 200)
	s.assertQueryRequest("token", "apptoken_admin", s.auth.RequireClient, 401)
	s.assertQueryRequest("token", "apptoken_admin", s.auth.RequireAdmin, 401)

	// clienttoken
	s.assertQueryRequest("token", "clienttoken", s.auth.RequireApplicationToken, 401)
	s.assertQueryRequest("token", "clienttoken", s.auth.RequireClient, 200)
	s.assertQueryRequest("token", "clienttoken", s.auth.RequireAdmin, 403)
	s.assertQueryRequest("token", "clienttoken_admin", s.auth.RequireApplicationToken, 401)
	s.assertQueryRequest("token", "clienttoken_admin", s.auth.RequireClient, 200)
	s.assertQueryRequest("token", "clienttoken_admin", s.auth.RequireAdmin, 200)
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
	s.auth.RequireApplicationToken()(ctx)
	assert.Equal(s.T(), 401, recorder.Code)
}

func (s *AuthenticationSuite) TestHeaderApiKeyToken() {
	// not existing token
	s.assertHeaderRequest("Authorization", "ApiKey ergerogerg", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "ApiKey ergerogerg", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "ApiKey ergerogerg", s.auth.RequireAdmin, 401)

	// no authentication schema
	s.assertHeaderRequest("Authorization", "ergerogerg", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "ergerogerg", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "ergerogerg", s.auth.RequireAdmin, 401)

	// wrong authentication schema
	s.assertHeaderRequest("Authorization", "ApiKeyx clienttoken", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "ApiKeyx clienttoken", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "ApiKeyx clienttoken", s.auth.RequireAdmin, 401)

	// not existing key
	s.assertHeaderRequest("Authorizationx", "ApiKey clienttoken", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorizationx", "ApiKey clienttoken", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorizationx", "ApiKey clienttoken", s.auth.RequireAdmin, 401)

	// apptoken
	s.assertHeaderRequest("Authorization", "ApiKey apptoken", s.auth.RequireApplicationToken, 200)
	s.assertHeaderRequest("Authorization", "ApiKey apptoken", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "ApiKey apptoken", s.auth.RequireAdmin, 401)
	s.assertHeaderRequest("Authorization", "ApiKey apptoken_admin", s.auth.RequireApplicationToken, 200)
	s.assertHeaderRequest("Authorization", "ApiKey apptoken_admin", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "ApiKey apptoken_admin", s.auth.RequireAdmin, 401)

	// clienttoken
	s.assertHeaderRequest("Authorization", "ApiKey clienttoken", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "ApiKey clienttoken", s.auth.RequireClient, 200)
	s.assertHeaderRequest("Authorization", "ApiKey clienttoken", s.auth.RequireAdmin, 403)
	s.assertHeaderRequest("Authorization", "ApiKey clienttoken_admin", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "ApiKey clienttoken_admin", s.auth.RequireClient, 200)
	s.assertHeaderRequest("Authorization", "ApiKey clienttoken_admin", s.auth.RequireAdmin, 200)
}

func (s *AuthenticationSuite) TestBasicAuth() {
	s.assertHeaderRequest("Authorization", "Basic ergerogerg", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "Basic ergerogerg", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "Basic ergerogerg", s.auth.RequireAdmin, 401)

	// user existing:pw
	s.assertHeaderRequest("Authorization", "Basic ZXhpc3Rpbmc6cHc=", s.auth.RequireApplicationToken, 403)
	s.assertHeaderRequest("Authorization", "Basic ZXhpc3Rpbmc6cHc=", s.auth.RequireClient, 200)
	s.assertHeaderRequest("Authorization", "Basic ZXhpc3Rpbmc6cHc=", s.auth.RequireAdmin, 403)

	// user admin:pw
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHc=", s.auth.RequireApplicationToken, 403)
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHc=", s.auth.RequireClient, 200)
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHc=", s.auth.RequireAdmin, 200)

	// user admin:pwx
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHd4", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHd4", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "Basic YWRtaW46cHd4", s.auth.RequireAdmin, 401)

	// user notexisting:pw
	s.assertHeaderRequest("Authorization", "Basic bm90ZXhpc3Rpbmc6cHc=", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "Basic bm90ZXhpc3Rpbmc6cHc=", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "Basic bm90ZXhpc3Rpbmc6cHc=", s.auth.RequireAdmin, 401)
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
