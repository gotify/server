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

func (s *UtilSuite) Test_getUserID() {
	s.expectUserID(func(ctx *gin.Context) { RegisterUser(ctx, &model.User{ID: 2}) }, 2)
	s.expectUserID(func(ctx *gin.Context) { RegisterClient(ctx, &model.Client{UserID: 5}) }, 5)
	s.expectUserID(func(ctx *gin.Context) { RegisterApplication(ctx, &model.Application{UserID: 7}) }, 7)

	assert.Panics(s.T(), func() {
		s.expectUserID(func(ctx *gin.Context) {}, 0)
	})

	s.expectTryUserID(func(ctx *gin.Context) {}, nil)
}

func (s *UtilSuite) Test_GetApplication() {
	app := &model.Application{Token: "atoken"}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterApplication(ctx, app)
	assert.Same(s.T(), app, GetApplication(ctx))

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	RegisterUser(ctx, &model.User{ID: 1})
	assert.Nil(s.T(), GetApplication(ctx))

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	assert.Nil(s.T(), GetApplication(ctx))
}

func (s *UtilSuite) Test_GetClient() {
	client := &model.Client{Token: "ctoken"}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	RegisterClient(ctx, client)
	assert.Same(s.T(), client, GetClient(ctx))

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	RegisterUser(ctx, &model.User{ID: 1})
	assert.Nil(s.T(), GetClient(ctx))

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	assert.Nil(s.T(), GetClient(ctx))
}

func (s *UtilSuite) expectUserID(register func(*gin.Context), expectedID uint) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	register(ctx)
	assert.Equal(s.T(), expectedID, GetUserID(ctx))
}

func (s *UtilSuite) expectTryUserID(register func(*gin.Context), expectedID *uint) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	register(ctx)
	assert.Equal(s.T(), expectedID, TryGetUserID(ctx))
}
