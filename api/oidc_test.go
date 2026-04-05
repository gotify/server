package api

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/decaymap"
	"github.com/gotify/server/v2/mode"
	"github.com/gotify/server/v2/test"
	"github.com/gotify/server/v2/test/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

var origGenClientToken = generateClientToken

func TestOIDCSuite(t *testing.T) {
	suite.Run(t, new(OIDCSuite))
}

type OIDCSuite struct {
	suite.Suite
	db       *testdb.Database
	a        *OIDCAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
	notified bool
}

func (s *OIDCSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.db = testdb.NewDB(s.T())
	s.notified = false
	notifier := new(UserChangeNotifier)
	notifier.OnUserAdded(func(uint) error {
		s.notified = true
		return nil
	})
	s.a = &OIDCAPI{
		DB:                 s.db.GormDatabase,
		UserChangeNotifier: notifier,
		UsernameClaim:      "preferred_username",
		AutoRegister:       true,
		pendingSessions:    decaymap.NewDecayMap[string, *pendingOIDCSession](time.Now(), pendingSessionMaxAge),
	}
}

func (s *OIDCSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

func (s *OIDCSuite) Test_GenerateState_Unique() {
	s1, _ := s.a.generateState()
	s2, _ := s.a.generateState()
	assert.NotEqual(s.T(), s1, s2)
}

func (s *OIDCSuite) Test_ResolveUser_ExistingUser() {
	s.db.NewUserWithName(1, "alice")

	info := &oidc.UserInfo{Claims: map[string]any{"preferred_username": "alice"}}
	user, status, err := s.a.resolveUser(info)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, status)
	assert.Equal(s.T(), "alice", user.Name)
	assert.Equal(s.T(), uint(1), user.ID)
	assert.False(s.T(), s.notified)
}

func (s *OIDCSuite) Test_ResolveUser_AutoRegister() {
	info := &oidc.UserInfo{Claims: map[string]any{"preferred_username": "newuser"}}
	user, status, err := s.a.resolveUser(info)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, status)
	assert.Equal(s.T(), "newuser", user.Name)
	assert.False(s.T(), user.Admin)
	assert.True(s.T(), s.notified)

	// verify persisted
	dbUser, err := s.db.GetUserByName("newuser")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), dbUser)
}

func (s *OIDCSuite) Test_ResolveUser_AutoRegisterDisabled() {
	s.a.AutoRegister = false
	info := &oidc.UserInfo{Claims: map[string]any{"preferred_username": "newuser"}}

	_, status, err := s.a.resolveUser(info)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), 403, status)
	s.db.AssertUsernameNotExist("newuser")
}

func (s *OIDCSuite) Test_ResolveUser_MissingClaim() {
	info := &oidc.UserInfo{Claims: map[string]any{}}

	_, status, err := s.a.resolveUser(info)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), 500, status)
	assert.Contains(s.T(), err.Error(), "preferred_username")
}

func (s *OIDCSuite) Test_ResolveUser_EmptyClaim() {
	info := &oidc.UserInfo{Claims: map[string]any{"preferred_username": ""}}

	_, status, err := s.a.resolveUser(info)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), 500, status)
}

func (s *OIDCSuite) Test_ResolveUser_NilClaim() {
	info := &oidc.UserInfo{Claims: map[string]any{"preferred_username": nil}}

	_, status, err := s.a.resolveUser(info)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), 500, status)
}

func (s *OIDCSuite) Test_ResolveUser_CustomClaim() {
	s.a.UsernameClaim = "email"
	s.db.NewUserWithName(1, "alice@example.com")

	info := &oidc.UserInfo{Claims: map[string]any{"email": "alice@example.com"}}
	user, _, err := s.a.resolveUser(info)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "alice@example.com", user.Name)
}

// --- createClient ---

func (s *OIDCSuite) Test_CreateClient() {
	generateClientToken = test.Tokens("Ctesttoken00001")
	defer func() { generateClientToken = origGenClientToken }()

	s.db.NewUser(1)
	client, err := s.a.createClient("MyPhone", 1)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "MyPhone", client.Name)
	assert.Equal(s.T(), "Ctesttoken00001", client.Token)
	assert.Equal(s.T(), uint(1), client.UserID)

	dbClient, err := s.db.GetClientByToken("Ctesttoken00001")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), dbClient)
}

// --- ExternalAuthorizeHandler ---

func (s *OIDCSuite) Test_ExternalAuthorizeHandler_MissingFields() {
	s.ctx.Request = httptest.NewRequest("POST", "/auth/oidc/external/authorize", strings.NewReader(`{}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.ExternalAuthorizeHandler(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

// --- ExternalTokenHandler ---

func (s *OIDCSuite) Test_ExternalTokenHandler_InvalidJSON() {
	s.ctx.Request = httptest.NewRequest("POST", "/auth/oidc/external/token", strings.NewReader(`{bad`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.ExternalTokenHandler(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *OIDCSuite) Test_ExternalTokenHandler_UnknownState() {
	s.ctx.Request = httptest.NewRequest("POST", "/auth/oidc/external/token", strings.NewReader(
		`{"code":"abc","state":"bogus","code_verifier":"v"}`,
	))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.a.ExternalTokenHandler(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}
