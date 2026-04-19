package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/auth"
	"github.com/gotify/server/v2/auth/password"
	"github.com/gotify/server/v2/mode"
	"github.com/gotify/server/v2/model"
	"github.com/gotify/server/v2/test"
	"github.com/gotify/server/v2/test/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSessionSuite(t *testing.T) {
	suite.Run(t, new(SessionSuite))
}

type SessionSuite struct {
	suite.Suite
	db       *testdb.Database
	a        *SessionAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
	notified bool
}

func (s *SessionSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	s.recorder = httptest.NewRecorder()
	s.db = testdb.NewDB(s.T())
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	withURL(s.ctx, "http", "example.com")
	s.notified = false
	s.a = &SessionAPI{DB: s.db, NotifyDeleted: s.notify}

	s.db.CreateUser(&model.User{
		Name: "testuser",
		Pass: password.CreatePassword("testpass", 5),
	})
}

func (s *SessionSuite) notify(uint, string) {
	s.notified = true
}

func (s *SessionSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

func (s *SessionSuite) Test_Login_Success() {
	originalGenerateClientToken := generateClientToken
	defer func() { generateClientToken = originalGenerateClientToken }()
	generateClientToken = test.Tokens("Ctesttoken12345")

	s.ctx.Request = httptest.NewRequest("POST", "/auth/local/login", strings.NewReader("name=test-browser"))
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.ctx.Request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testuser:testpass")))

	s.a.Login(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)

	// Verify HttpOnly cookie is set
	cookies := s.recorder.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == auth.CookieName {
			sessionCookie = c
			break
		}
	}
	assert.NotNil(s.T(), sessionCookie)
	assert.Equal(s.T(), "Ctesttoken12345", sessionCookie.Value)
	assert.True(s.T(), sessionCookie.HttpOnly)
	assert.Equal(s.T(), "/", sessionCookie.Path)
	assert.Equal(s.T(), http.SameSiteStrictMode, sessionCookie.SameSite)

	body := s.recorder.Body.String()
	assert.Contains(s.T(), body, "testuser")
	assert.NotContains(s.T(), body, "Ctesttoken12345")

	clients, err := s.db.GetClientsByUser(1)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), clients, 1)
	assert.Equal(s.T(), "test-browser", clients[0].Name)
}

func (s *SessionSuite) Test_Login_WrongPassword() {
	s.ctx.Request = httptest.NewRequest("POST", "/auth/local/login", strings.NewReader("name=test-browser"))
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.ctx.Request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testuser:wrongpass")))

	s.a.Login(s.ctx)

	assert.Equal(s.T(), 401, s.recorder.Code)

	// No cookie should be set
	cookies := s.recorder.Result().Cookies()
	for _, c := range cookies {
		assert.NotEqual(s.T(), auth.CookieName, c.Name)
	}
}

func (s *SessionSuite) Test_Logout_Success() {
	builder := s.db.User(5)
	builder.ClientWithToken(1, "Ctesttoken12345")

	s.ctx.Request = httptest.NewRequest("POST", "/auth/logout", nil)
	test.WithUser(s.ctx, 5)
	s.ctx.Set("tokenid", "Ctesttoken12345")

	s.a.Logout(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.True(s.T(), s.notified)

	cookies := s.recorder.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == auth.CookieName {
			sessionCookie = c
			break
		}
	}
	assert.NotNil(s.T(), sessionCookie)
	assert.Equal(s.T(), "", sessionCookie.Value)
	assert.True(s.T(), sessionCookie.MaxAge < 0)

	s.db.AssertClientNotExist(1)
}
