package docs

import (
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v3/mode"
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

func TestUIExternalResourcesHaveIntegrity(t *testing.T) {
	tags := regexp.MustCompile(`<(?:script|link)[^>]*https://[^>]*>`).FindAllString(ui, -1)
	assert.Len(t, tags, 3)
	for _, tag := range tags {
		assert.Regexp(t, `integrity="sha512-[A-Za-z0-9+/=]+"`, tag)
		assert.Contains(t, tag, `crossorigin="anonymous"`)
	}
}
