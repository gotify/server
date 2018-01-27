package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

func TestUtilSuite(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}

type UtilSuite struct {
	suite.Suite
}

func (s *UtilSuite) BeforeTest(suiteName, testName string) {
	gin.SetMode(gin.TestMode)
}

func (s *UtilSuite) Test_getId() {
	s.expectUserIdWith(&model.User{Id: 2}, 0, 2)
	s.expectUserIdWith(nil, 5, 5)
	assert.Panics(s.T(), func() {
		s.expectUserIdWith(nil, 0, 0)
	})
}

func (s *UtilSuite) expectUserIdWith(user *model.User, tokenId uint, expectedId uint) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterAuthentication(ctx, user, tokenId)
	actualId := GetUserId(ctx)
	assert.Equal(s.T(), expectedId, actualId)
}
