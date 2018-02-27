package docs

import (
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"io/ioutil"
)

func TestServe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "/swagger", nil)
	ctx.Request.URL.Host = "localhost"

	Serve(ctx)

	actualFileContent := getActualSpecFileContent(t)
	packrFileContent := recorder.Body.String()
	assert.JSONEq(t, packrFileContent, actualFileContent, "packr and spec file are out of sync")
}

func getActualSpecFileContent(t *testing.T) string {
	bytes, err := ioutil.ReadFile("spec.json")
	assert.Nil(t, err)
	return string(bytes)
}