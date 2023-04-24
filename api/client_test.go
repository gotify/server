package api

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/mode"
	"github.com/gotify/server/v2/model"
	"github.com/gotify/server/v2/test"
	"github.com/gotify/server/v2/test/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	firstClientToken  = "Caaaaaaaaaaaaaa"
	secondClientToken = "Cbbbbbbbbbbbbbb"
)

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

type ClientSuite struct {
	suite.Suite
	db       *testdb.Database
	a        *ClientAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
	notified bool
}

var originalGenerateClientToken func() string

func (s *ClientSuite) BeforeTest(suiteName, testName string) {
	originalGenerateClientToken = generateClientToken
	generateClientToken = test.Tokens(firstClientToken, secondClientToken)
	mode.Set(mode.TestDev)
	s.recorder = httptest.NewRecorder()
	s.db = testdb.NewDB(s.T())
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	withURL(s.ctx, "http", "example.com")
	s.notified = false
	s.a = &ClientAPI{DB: s.db, NotifyDeleted: s.notify}
}

func (s *ClientSuite) notify(uint, string) {
	s.notified = true
}

func (s *ClientSuite) AfterTest(suiteName, testName string) {
	generateClientToken = originalGenerateClientToken
	s.db.Close()
}

func (s *ClientSuite) Test_ensureClientHasCorrectJsonRepresentation() {
	actual := &model.Client{ID: 1, UserID: 2, Token: "Casdasfgeeg", Name: "myclient"}
	test.JSONEquals(s.T(), actual, `{"id":1,"token":"Casdasfgeeg","name":"myclient"}`)
}

func (s *ClientSuite) Test_CreateClient_mapAllParameters() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name&description=description_text")

	s.a.CreateClient(s.ctx)

	expected := &model.Client{ID: 1, Token: firstClientToken, UserID: 5, Name: "custom_name"}
	assert.Equal(s.T(), 200, s.recorder.Code)
	if clients, err := s.db.GetClientsByUser(5); assert.NoError(s.T(), err) {
		assert.Contains(s.T(), clients, expected)
	}
}

func (s *ClientSuite) Test_CreateClient_ignoresReadOnlyPropertiesInParams() {
	s.db.User(5)
	test.WithUser(s.ctx, 5)

	s.withFormData("name=myclient&ID=45&Token=12341234&UserID=333")

	s.a.CreateClient(s.ctx)
	expected := &model.Client{ID: 1, UserID: 5, Token: firstClientToken, Name: "myclient"}

	assert.Equal(s.T(), 200, s.recorder.Code)
	if clients, err := s.db.GetClientsByUser(5); assert.NoError(s.T(), err) {
		assert.Contains(s.T(), clients, expected)
	}
}

func (s *ClientSuite) Test_CreateClient_expectBadRequestOnEmptyName() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=&description=description_text")

	s.a.CreateClient(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
	if clients, err := s.db.GetClientsByUser(5); assert.NoError(s.T(), err) {
		assert.Empty(s.T(), clients)
	}
}

func (s *ClientSuite) Test_DeleteClient_expectNotFoundOnCurrentUserIsNotOwner() {
	s.db.User(5).Client(7)
	s.db.User(2)

	test.WithUser(s.ctx, 2)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/7", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "7"}}

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
	s.db.AssertClientExist(7)
}

func (s *ClientSuite) Test_CreateClient_returnsClientWithID() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateClient(s.ctx)

	expected := &model.Client{ID: 1, Token: firstClientToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), expected, s.recorder)
}

func (s *ClientSuite) Test_CreateClient_withExistingToken() {
	s.db.User(5).ClientWithToken(1, firstClientToken)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateClient(s.ctx)

	expected := &model.Client{ID: 2, Token: secondClientToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), expected, s.recorder)
}

func (s *ClientSuite) Test_GetClients() {
	userBuilder := s.db.User(5)
	first := userBuilder.NewClientWithToken(1, "perfper")
	second := userBuilder.NewClientWithToken(2, "asdasd")

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("GET", "/tokens", nil)

	s.a.GetClients(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), []*model.Client{first, second}, s.recorder)
}

func (s *ClientSuite) Test_DeleteClient_expectNotFound() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstClientToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "8"}}

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *ClientSuite) Test_DeleteClient() {
	s.db.User(5).Client(8)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstClientToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "8"}}

	assert.False(s.T(), s.notified)

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertClientNotExist(8)
	assert.True(s.T(), s.notified)
}

func (s *ClientSuite) Test_UpdateClient_expectSuccess() {
	s.db.User(5).NewClientWithToken(1, firstClientToken)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=firefox")
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	s.a.UpdateClient(s.ctx)

	expected := &model.Client{
		ID:     1,
		Token:  firstClientToken,
		UserID: 5,
		Name:   "firefox",
	}

	assert.Equal(s.T(), 200, s.recorder.Code)
	if client, err := s.db.GetClientByID(1); assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expected, client)
	}
}

func (s *ClientSuite) Test_UpdateClient_expectNotFound() {
	test.WithUser(s.ctx, 5)
	s.ctx.Params = gin.Params{{Key: "id", Value: "2"}}
	s.a.UpdateClient(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *ClientSuite) Test_UpdateClient_WithMissingAttributes_expectBadRequest() {
	test.WithUser(s.ctx, 5)
	s.a.UpdateClient(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *ClientSuite) withFormData(formData string) {
	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(formData))
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

func withURL(ctx *gin.Context, scheme, host string) {
	ctx.Set("location", &url.URL{Scheme: scheme, Host: host})
}
