package docs

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/swagger", nil)
	ctx.Request.URL.Host = "localhost"

	Serve(ctx)

	content := recorder.Body.String()
	assert.NotEmpty(t, content)
}
