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

func (s *UtilSuite) Test_getID() {
	s.expectUserIDWith(&model.User{ID: 2}, 0, 2)
	s.expectUserIDWith(nil, 5, 5)
	assert.Panics(s.T(), func() {
		s.expectUserIDWith(nil, 0, 0)
	})
}

func (s *UtilSuite) expectUserIDWith(user *model.User, tokenID uint, expectedID uint) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterAuthentication(ctx, user, tokenID)
	actualID := GetUserID(ctx)
	assert.Equal(s.T(), expectedID, actualID)
}
