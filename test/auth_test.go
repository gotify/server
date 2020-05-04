package test_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/auth"
	"github.com/gotify/server/v2/mode"
	"github.com/gotify/server/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestFakeAuth(t *testing.T) {
	mode.Set(mode.TestDev)

	ctx, _ := gin.CreateTestContext(nil)
	test.WithUser(ctx, 5)
	assert.Equal(t, uint(5), auth.GetUserID(ctx))
}
