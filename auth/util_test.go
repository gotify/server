package auth

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/mode"
	"github.com/gotify/server/v2/model"
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
	s.expectTryUserIDWith(nil, 0, nil)
}

func (s *UtilSuite) Test_getToken() {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterAuthentication(ctx, nil, 1, 0, "asdasda", 0)
	actualID := GetTokenID(ctx)
	assert.Equal(s.T(), "asdasda", actualID)
}

func (s *UtilSuite) expectUserIDWith(user *model.User, tokenUserID, expectedID uint) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterAuthentication(ctx, user, tokenUserID, 0, "", 0)
	actualID := GetUserID(ctx)
	assert.Equal(s.T(), expectedID, actualID)
}

func (s *UtilSuite) expectTryUserIDWith(user *model.User, tokenUserID uint, expectedID *uint) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterAuthentication(ctx, user, tokenUserID, 0, "", 0)
	actualID := TryGetUserID(ctx)
	assert.Equal(s.T(), expectedID, actualID)
}
