package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gotify/server/config"
	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
	"github.com/gotify/server/test/testdb"
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
	db       *testdb.Database
	server   *httptest.Server
	closable func()
}

func (s *IntegrationSuite) BeforeTest(string, string) {
	mode.Set(mode.TestDev)
	var err error
	s.db = testdb.NewDBWithDefaultUser(s.T())
	assert.Nil(s.T(), err)
	g, closable := Create(s.db.GormDatabase,
		&model.VersionInfo{Version: "1.0.0", BuildDate: "2018-02-20-17:30:47", Commit: "asdasds"},
		&config.Configuration{PassStrength: 5},
	)
	s.closable = closable
	s.server = httptest.NewServer(g)
}

func (s *IntegrationSuite) AfterTest(string, string) {
	s.closable()
	s.db.Close()
	s.server.Close()
}

func (s *IntegrationSuite) TestVersionInfo() {
	req := s.newRequest("GET", "version", "")

	doRequestAndExpect(s.T(), req, 200, `{"version":"1.0.0", "commit":"asdasds", "buildDate":"2018-02-20-17:30:47"}`)
}

func (s *IntegrationSuite) TestHeaderInDev() {
	mode.Set(mode.TestDev)
	req := s.newRequest("GET", "version", "")

	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), res.Header.Get("Access-Control-Allow-Origin"))
}

func (s *IntegrationSuite) TestHeaderInProd() {
	mode.Set(mode.Prod)
	req := s.newRequest("GET", "version", "")

	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Empty(s.T(), res.Header.Get("Access-Control-Allow-Origin"))
}

func TestHeadersFromConfiguration(t *testing.T) {
	mode.Set(mode.Prod)
	db := testdb.NewDBWithDefaultUser(t)
	defer db.Close()

	config := config.Configuration{PassStrength: 5}
	config.Server.ResponseHeaders = map[string]string{
		"New-Cool-Header":             "Nice",
		"Access-Control-Allow-Origin": "---",
	}

	g, closable := Create(db.GormDatabase,
		&model.VersionInfo{Version: "1.0.0", BuildDate: "2018-02-20-17:30:47", Commit: "asdasds"},
		&config,
	)
	server := httptest.NewServer(g)

	defer func() {
		closable()
		server.Close()
	}()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", server.URL, "version"), nil)
	req.Header.Add("Content-Type", "application/json")
	assert.Nil(t, err)

	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, "---", res.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Nice", res.Header.Get("New-Cool-Header"))
}

func (s *IntegrationSuite) TestOptionsRequest() {
	req := s.newRequest("OPTIONS", "version", "")

	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), res.StatusCode, 200)
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
	req.Header.Add("X-Gotify-Key", token.Token)
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 200, res.StatusCode)

	req = s.newRequest("GET", "message", "")
	req.SetBasicAuth("admin", "pw")
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 200, res.StatusCode)
	msgs := &model.PagedMessages{}
	json.NewDecoder(res.Body).Decode(&msgs)
	assert.Len(s.T(), msgs.Messages, 1)

	msg := msgs.Messages[0]
	assert.Equal(s.T(), "backup done", msg.Message)
	assert.Equal(s.T(), "backup done", msg.Title)
	assert.Equal(s.T(), uint(1), msg.ID)
	assert.Equal(s.T(), token.ID, msg.ApplicationID)
}

func (s *IntegrationSuite) TestPluginLoadFail_expectPanic() {
	db := testdb.NewDBWithDefaultUser(s.T())
	defer db.Close()

	assert.Panics(s.T(), func() {
		Create(db.GormDatabase, new(model.VersionInfo), &config.Configuration{
			PluginsDir: "<THIS_PATH_IS_MALFORMED>",
		})
	})
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
