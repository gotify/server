package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v3/mode"
	"github.com/gotify/server/v3/test/testdb"
	"github.com/stretchr/testify/assert"
)

func TestMCP_HandleWithoutApplication(t *testing.T) {
	mode.Set(mode.TestDev)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("POST", "/mcp", nil)

	NewMCP(&MessageAPI{}, "1.0.0").Handle(ctx)

	assert.Equal(t, 403, recorder.Code)
	assert.Len(t, ctx.Errors, 1)
}

func TestMCP_ServerForRequestWithoutApplication(t *testing.T) {
	a := NewMCP(&MessageAPI{}, "1.0.0")
	assert.Nil(t, a.serverForRequest(httptest.NewRequest("POST", "/mcp", nil)))
}

func TestMCP_SendMessageDatabaseError(t *testing.T) {
	mode.Set(mode.TestDev)
	db := testdb.NewDBWithDefaultUser(t)
	app := db.User(5).NewAppWithToken(3, "apptoken")
	db.Close()

	a := NewMCP(&MessageAPI{DB: db}, "1.0.0")
	res, _, err := a.sendMessage(app, SendMessageInput{Message: "hi"})

	assert.Nil(t, res)
	assert.EqualError(t, err, "failed to send message")
}
