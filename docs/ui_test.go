package docs

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/mode"
	"github.com/stretchr/testify/assert"
)

func TestUI(t *testing.T) {
	mode.Set(mode.TestDev)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	withURL(ctx, "http", "example.com")

	ctx.Request = httptest.NewRequest("GET", "/swagger", nil)

	UI(ctx)

	content := recorder.Body.String()
	assert.NotEmpty(t, content)
}
