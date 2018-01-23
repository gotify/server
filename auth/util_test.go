package auth

import (
	"github.com/stretchr/testify/suite"
	"github.com/gin-gonic/gin"
	"testing"
	"net/http/httptest"
	"github.com/jmattheis/memo/model"
	"github.com/stretchr/testify/assert"
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
	s.expectUserIdWith(&model.User{ID: 2}, nil, 2)
	s.expectUserIdWith(nil, &model.Token{UserID: 5}, 5)
	assert.Panics(s.T(), func() {
		s.expectUserIdWith(nil, nil, 0)
	})
}

func (s *UtilSuite) expectUserIdWith(user *model.User, token *model.Token, id uint) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterAuthentication(ctx, user, token)
	actualId := GetUserID(ctx)
	assert.Equal(s.T(), id, actualId)
}
