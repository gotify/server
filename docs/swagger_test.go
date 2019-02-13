package docs

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/mode"
	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	mode.Set(mode.TestDev)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	withURL(ctx, "http", "example.com")

	ctx.Request = httptest.NewRequest("GET", "/swagger?base="+url.QueryEscape("127.0.0.1/proxy/"), nil)

	Serve(ctx)

	content := recorder.Body.String()
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "127.0.0.1/proxy/")
}

func withURL(ctx *gin.Context, scheme, host string) {
	ctx.Set("location", &url.URL{Scheme: scheme, Host: host})
}
