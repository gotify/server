package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/mode"
	"github.com/gotify/server/v2/model"
	"github.com/gotify/server/v2/test"
	"github.com/gotify/server/v2/test/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func testConsistentHead(t *testing.T, head, get http.Header) {
	// The server SHOULD send the same header fields in response to a HEAD request as it would have sent if the request method had been GET.
	// However, a server MAY omit header fields for which a value is determined only while generating the content.
	assert.Empty(t, head.Get("Content-Length"), "Content-Length should be empty")
	assert.Equal(t, get.Get("Content-Type"), head.Get("Content-Type"), "Content-Type should be the same")
	assert.Equal(t, get.Get("Transfer-Encoding"), head.Get("Transfer-Encoding"), "Transfer-Encoding should be the same")
	assert.Equal(t, get.Get("Connection"), head.Get("Connection"), "Connection should be the same")
}

func TestHealthSuite(t *testing.T) {
	suite.Run(t, new(HealthSuite))
}

type HealthSuite struct {
	suite.Suite
	db       *testdb.Database
	a        *HealthAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func (s *HealthSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	s.recorder = httptest.NewRecorder()
	s.db = testdb.NewDB(s.T())
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	withURL(s.ctx, "http", "example.com")
	s.a = &HealthAPI{DB: s.db}
}

func (s *HealthSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

func (s *HealthSuite) TestHealthSuccess() {
	head, err := http.NewRequest("HEAD", "/health", nil)
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx.Request = head
	s.a.Health(s.ctx)
	headHeaders := s.recorder.Header().Clone()

	request, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx.Request = request
	s.a.Health(s.ctx)

	testConsistentHead(s.T(), headHeaders, s.recorder.Header())
	test.BodyEquals(s.T(), model.Health{Health: model.StatusGreen, Database: model.StatusGreen}, s.recorder)
}

func (s *HealthSuite) TestDatabaseFailure() {
	s.db.Close()

	head, err := http.NewRequest("HEAD", "/health", nil)
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx.Request = head
	s.a.Health(s.ctx)
	headHeaders := s.recorder.Header().Clone()

	request, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx.Request = request
	s.a.Health(s.ctx)

	testConsistentHead(s.T(), headHeaders, s.recorder.Header())
	test.BodyEquals(s.T(), model.Health{Health: model.StatusOrange, Database: model.StatusRed}, s.recorder)
}
