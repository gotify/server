package runner

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type errRoundTripper struct{ err error }

func (e errRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, e.err
}

func TestLoggingRoundTripperReturnsTransportError(t *testing.T) {
	want := errors.New("dial tcp: lookup acme-v02.api.letsencrypt.org: no such host")
	rt := &LoggingRoundTripper{Name: "Let's Encrypt", RoundTripper: errRoundTripper{err: want}}
	req, err := http.NewRequest(http.MethodGet, "https://acme-v02.api.letsencrypt.org/directory", nil)
	require.NoError(t, err)
	resp, err := rt.RoundTrip(req)
	require.Nil(t, resp)
	require.ErrorIs(t, err, want)
}
