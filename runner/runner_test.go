package runner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	cases := []struct {
		Request string
		TLS     int
		Expect  string
	}{
		{Request: "http://gotify.net/meow", TLS: 443, Expect: "https://gotify.net:443/meow"},
		{Request: "http://gotify.net:8080/meow", TLS: 443, Expect: "https://gotify.net:443/meow"},
		{Request: "http://gotify.net:8080/meow", TLS: 8443, Expect: "https://gotify.net:8443/meow"},
	}

	for _, testCase := range cases {
		name := fmt.Sprintf("%s -- %d -> %s", testCase.Request, testCase.TLS, testCase.Expect)
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", testCase.Request, nil)
			rec := httptest.NewRecorder()

			redirectToHTTPS(fmt.Sprint(testCase.TLS)).ServeHTTP(rec, req)

			assert.Equal(t, http.StatusFound, rec.Result().StatusCode)
			assert.Equal(t, testCase.Expect, rec.Header().Get("location"))
		})
	}
}
