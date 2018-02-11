package router

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"github.com/jmattheis/memo/database"
	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	client        = &http.Client{}
	forbiddenJSON = `{"error":"Forbidden", "errorCode":403, "errorDescription":"you are not allowed to access this api"}`
)

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

type IntegrationSuite struct {
	suite.Suite
	db       *database.GormDatabase
	server   *httptest.Server
	closable func()
}

func (s *IntegrationSuite) BeforeTest(string, string) {
	gin.SetMode(gin.TestMode)
	var err error
	s.db, err = database.New("sqlite3", "itest.db", "admin", "pw")
	assert.Nil(s.T(), err)
	g, closable := Create(s.db)
	s.closable = closable
	s.server = httptest.NewServer(g)
}

func (s *IntegrationSuite) AfterTest(string, string) {
	s.closable()
	s.db.Close()
	assert.Nil(s.T(), os.Remove("itest.db"))
	s.server.Close()
}

func (s *IntegrationSuite) TestSendMessage() {
	req := s.newRequest("POST", "application", `{"name": "backup-server"}`)
	req.SetBasicAuth("admin", "pw")
	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 200, res.StatusCode)
	token := &model.Application{}
	json.NewDecoder(res.Body).Decode(token)
	assert.Equal(s.T(), "backup-server", token.Name)

	req = s.newRequest("POST", "message", `{"message": "backup done", "title": "backup done"}`)
	req.Header.Add("Authorization", fmt.Sprintf("ApiKey %s", token.ID))
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 200, res.StatusCode)

	req = s.newRequest("GET", "message", "")
	req.SetBasicAuth("admin", "pw")
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 200, res.StatusCode)
	msgs := &[]*model.Message{}
	json.NewDecoder(res.Body).Decode(msgs)
	assert.Len(s.T(), *msgs, 1)
	msg := (*msgs)[0]
	assert.Equal(s.T(), "backup done", msg.Message)
	assert.Equal(s.T(), "backup done", msg.Title)
	assert.Equal(s.T(), uint(1), msg.ID)
	assert.Equal(s.T(), token.ID, msg.ApplicationID)
}

func (s *IntegrationSuite) TestAuthentication() {
	req := s.newRequest("GET", "current/user", "")
	req.SetBasicAuth("admin", "pw")
	doRequestAndExpect(s.T(), req, 200, `{"id": 1, "name": "admin", "admin": true}`)

	req = s.newRequest("GET", "current/user", "")
	req.SetBasicAuth("jmattheis", "pw")
	doRequestAndExpect(s.T(), req, 401, `{"error":"Unauthorized", "errorCode":401, "errorDescription":"you need to provide a valid access token or user credentials to access this api"}`)

	req = s.newRequest("POST", "user", `{"name": "normal", "pass": "secret"}`)
	req.SetBasicAuth("admin", "pw")
	doRequestAndExpect(s.T(), req, 200, `{"id": 2, "name": "normal", "admin": false}`)

	req = s.newRequest("POST", "user", `{"name": "normal2", "pass": "secret"}`)
	req.SetBasicAuth("normal", "secret")
	doRequestAndExpect(s.T(), req, 403, forbiddenJSON)

	req = s.newRequest("POST", "message", `{"message": "backup done", "title": "backup"}`)
	req.SetBasicAuth("normal", "secret")
	doRequestAndExpect(s.T(), req, 403, forbiddenJSON)

	req = s.newRequest("GET", "current/user", "")
	req.SetBasicAuth("normal", "secret")
	doRequestAndExpect(s.T(), req, 200, `{"id": 2, "name": "normal", "admin": false}`)

	req = s.newRequest("POST", "client", `{"name": "android-client"}`)
	req.SetBasicAuth("normal", "secret")
	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 200, res.StatusCode)
	token := &model.Application{}
	json.NewDecoder(res.Body).Decode(token)
	assert.Equal(s.T(), "android-client", token.Name)
}

func (s *IntegrationSuite) newRequest(method, url string, body string) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", s.server.URL, url), strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	assert.Nil(s.T(), err)
	return req
}

func doRequestAndExpect(t *testing.T, req *http.Request, code int, json string) {
	res, err := client.Do(req)
	assert.Nil(t, err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	assert.Equal(t, code, res.StatusCode)
	assert.JSONEq(t, json, buf.String())
}
