package api

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorHandling(t *testing.T) {
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)
	successOrAbort(ctx, 500, errors.New("err"))

	if rec.Code != 500 {
		t.Fail()
	}
}
