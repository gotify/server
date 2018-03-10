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
	s.DB.On("GetClientByToken", "clienttoken").Return(&model.Client{ID: 1, Token: "clienttoken", UserID: 1, Name: "android phone"})
	s.DB.On("GetClientByToken", "clienttoken_admin").Return(&model.Client{ID: 2, Token: "clienttoken_admin", UserID: 2, Name: "android phone2"})
	s.DB.On("GetClientByToken", mock.Anything).Return(nil)
	s.DB.On("GetApplicationByToken", "apptoken").Return(&model.Application{ID: 3, Token: "apptoken", UserID: 1, Name: "backup server", Description: "irrelevant"})
	s.DB.On("GetApplicationByToken", "apptoken_admin").Return(&model.Application{ID: 4, Token: "apptoken_admin", UserID: 2, Name: "backup server", Description: "irrelevant"})
	s.DB.On("GetApplicationByToken", mock.Anything).Return(nil)

	s.DB.On("GetUserByID", uint(1)).Return(&model.User{ID: 1, Name: "irrelevant", Admin: false})
	s.DB.On("GetUserByID", uint(2)).Return(&model.User{ID: 2, Name: "irrelevant", Admin: true})

	s.DB.On("GetUserByName", "existing").Return(&model.User{Name: "existing", Pass: CreatePassword("pw", 5)})
	s.DB.On("GetUserByName", "admin").Return(&model.User{Name: "admin", Pass: CreatePassword("pw", 5), Admin: true})
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
	s.assertHeaderRequest("X-Gotify-Key", "ergerogerg", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("X-Gotify-Key", "ergerogerg", s.auth.RequireClient, 401)
	s.assertHeaderRequest("X-Gotify-Key", "ergerogerg", s.auth.RequireAdmin, 401)

	// no authentication schema
	s.assertHeaderRequest("Authorization", "ergerogerg", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "ergerogerg", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "ergerogerg", s.auth.RequireAdmin, 401)

	// wrong authentication schema
	s.assertHeaderRequest("Authorization", "ApiKeyx clienttoken", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("Authorization", "ApiKeyx clienttoken", s.auth.RequireClient, 401)
	s.assertHeaderRequest("Authorization", "ApiKeyx clienttoken", s.auth.RequireAdmin, 401)

	// not existing key
	s.assertHeaderRequest("X-Gotify-Keyx", "clienttoken", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("X-Gotify-Keyx", "clienttoken", s.auth.RequireClient, 401)
	s.assertHeaderRequest("X-Gotify-Keyx", "clienttoken", s.auth.RequireAdmin, 401)

	// apptoken
	s.assertHeaderRequest("X-Gotify-Key", "apptoken", s.auth.RequireApplicationToken, 200)
	s.assertHeaderRequest("X-Gotify-Key", "apptoken", s.auth.RequireClient, 401)
	s.assertHeaderRequest("X-Gotify-Key", "apptoken", s.auth.RequireAdmin, 401)
	s.assertHeaderRequest("X-Gotify-Key", "apptoken_admin", s.auth.RequireApplicationToken, 200)
	s.assertHeaderRequest("X-Gotify-Key", "apptoken_admin", s.auth.RequireClient, 401)
	s.assertHeaderRequest("X-Gotify-Key", "apptoken_admin", s.auth.RequireAdmin, 401)

	// clienttoken
	s.assertHeaderRequest("X-Gotify-Key", "clienttoken", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("X-Gotify-Key", "clienttoken", s.auth.RequireClient, 200)
	s.assertHeaderRequest("X-Gotify-Key", "clienttoken", s.auth.RequireAdmin, 403)
	s.assertHeaderRequest("X-Gotify-Key", "clienttoken_admin", s.auth.RequireApplicationToken, 401)
	s.assertHeaderRequest("X-Gotify-Key", "clienttoken_admin", s.auth.RequireClient, 200)
	s.assertHeaderRequest("X-Gotify-Key", "clienttoken_admin", s.auth.RequireAdmin, 200)
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
