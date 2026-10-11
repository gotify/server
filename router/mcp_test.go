package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gotify/server/v3/config"
	"github.com/gotify/server/v3/mode"
	"github.com/gotify/server/v3/model"
	"github.com/gotify/server/v3/test/testdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type headerTransport struct {
	header string
	value  string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set(t.header, t.value)
	return http.DefaultTransport.RoundTrip(req)
}

func startMCPServer(t *testing.T, enabled bool) (*testdb.Database, *httptest.Server) {
	mode.Set(mode.TestDev)
	db := testdb.NewDBWithDefaultUser(t)
	g, closable := Create(
		db.GormDatabase,
		&model.VersionInfo{Version: "1.0.0"},
		&config.Configuration{PassStrength: 5, LocalAuthEnabled: true, MCP: config.MCP{Enabled: enabled}},
	)
	server := httptest.NewServer(g)
	t.Cleanup(func() {
		server.Close()
		closable()
		db.Close()
	})
	return db, server
}

func connectMCP(t *testing.T, url, header, value string) (*mcp.ClientSession, error) {
	c := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	return c.Connect(context.Background(), &mcp.StreamableClientTransport{
		Endpoint:             url,
		HTTPClient:           &http.Client{Transport: &headerTransport{header: header, value: value}},
		DisableStandaloneSSE: true,
	}, nil)
}

func TestMCP_SendMessage(t *testing.T) {
	db, server := startMCPServer(t, true)
	app := db.User(5).NewAppWithTokenAndDefaultPriority(3, "apptoken", 4)

	session, err := connectMCP(t, server.URL+"/mcp", "Authorization", "Bearer apptoken")
	require.NoError(t, err)
	defer session.Close()

	tools, err := session.ListTools(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, tools.Tools, 1)
	assert.Equal(t, "send_message", tools.Tools[0].Name)

	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "send_message",
		Arguments: map[string]any{
			"message":       "**backup** done",
			"markdown":      true,
			"click_url":     "https://example.com",
			"big_image_url": "https://example.com/image.png",
		},
	})
	require.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Equal(t, "Message sent (id 1).", res.Content[0].(*mcp.TextContent).Text)

	msgs, err := db.GetMessagesByApplication(app.ID)
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	msg := msgs[0]
	assert.Equal(t, "**backup** done", msg.Message)
	assert.Equal(t, app.Name, msg.Title)
	assert.Equal(t, 4, msg.Priority)
	assert.JSONEq(t, `{
		"client::display": {"contentType": "text/markdown"},
		"client::notification": {"click": {"url": "https://example.com"}, "bigImageUrl": "https://example.com/image.png"}
	}`, string(msg.Extras))

	res, err = session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "send_message",
		Arguments: map[string]any{"message": "plain", "title": "custom", "priority": 8},
	})
	require.NoError(t, err)
	assert.False(t, res.IsError)

	msgs, err = db.GetMessagesByApplication(app.ID)
	require.NoError(t, err)
	require.Len(t, msgs, 2)
	msg = msgs[0]
	assert.Equal(t, "plain", msg.Message)
	assert.Equal(t, "custom", msg.Title)
	assert.Equal(t, 8, msg.Priority)
	assert.Empty(t, msg.Extras)
}

func TestMCP_EmptyMessage(t *testing.T) {
	db, server := startMCPServer(t, true)
	db.User(5).AppWithToken(3, "apptoken")

	session, err := connectMCP(t, server.URL+"/mcp?token=apptoken", "X-Ignored", "")
	require.NoError(t, err)
	defer session.Close()

	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "send_message",
		Arguments: map[string]any{"message": ""},
	})
	require.NoError(t, err)
	assert.True(t, res.IsError)
	db.AssertMessageNotExist(1)
}

func TestMCP_RequiresApplicationToken(t *testing.T) {
	db, server := startMCPServer(t, true)
	db.User(5).ClientWithToken(1, "clienttoken")

	for _, token := range []string{"clienttoken", "unknown"} {
		_, err := connectMCP(t, server.URL+"/mcp", "X-Gotify-Key", token)
		assert.Error(t, err, token)
	}

	for token, code := range map[string]int{"clienttoken": 401, "": 401} {
		req, err := http.NewRequest("POST", server.URL+"/mcp", strings.NewReader(`{}`))
		require.NoError(t, err)
		req.Header.Set("X-Gotify-Key", token)
		res, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, code, res.StatusCode, token)
	}

	req, err := http.NewRequest("POST", server.URL+"/mcp", strings.NewReader(`{}`))
	require.NoError(t, err)
	req.SetBasicAuth("admin", "pw")
	res, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 403, res.StatusCode)
}

func TestMCP_Disabled(t *testing.T) {
	db, server := startMCPServer(t, false)
	db.User(5).AppWithToken(3, "apptoken")

	req, err := http.NewRequest("POST", server.URL+"/mcp", strings.NewReader(`{}`))
	require.NoError(t, err)
	req.Header.Set("X-Gotify-Key", "apptoken")
	res, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
}
