package auth

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUtilSuite(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}

type UtilSuite struct {
	suite.Suite
}

func (s *UtilSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
}

func (s *UtilSuite) Test_getID() {
	s.expectUserIDWith(&model.User{ID: 2}, 0, 2)
	s.expectUserIDWith(nil, 5, 5)
	assert.Panics(s.T(), func() {
		s.expectUserIDWith(nil, 0, 0)
	})
}

func (s *UtilSuite) Test_getToken() {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterAuthentication(ctx, nil, 1, "asdasda")
	actualID := GetTokenID(ctx)
	assert.Equal(s.T(), "asdasda", actualID)
}

func (s *UtilSuite) expectUserIDWith(user *model.User, tokenUserID uint, expectedID uint) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterAuthentication(ctx, user, tokenUserID, "")
	actualID := GetUserID(ctx)
	assert.Equal(s.T(), expectedID, actualID)
}
