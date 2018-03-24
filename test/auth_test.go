package test_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/auth"
	"github.com/gotify/server/mode"
	"github.com/gotify/server/test"
	"github.com/stretchr/testify/assert"
)

func TestFakeAuth(t *testing.T) {
	mode.Set(mode.TestDev)

	ctx, _ := gin.CreateTestContext(nil)
	test.WithUser(ctx, 5)
	assert.Equal(t, uint(5), auth.GetUserID(ctx))
}
