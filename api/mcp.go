package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v3/auth"
	"github.com/gotify/server/v3/model"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog/log"
)

type mcpAppKey struct{}

// MCPAPI provides a Model Context Protocol endpoint which lets AI agents send messages
// with an application token.
type MCPAPI struct {
	messages    *MessageAPI
	version     string
	schemaCache *mcp.SchemaCache
	handler     http.Handler
}

// SendMessageInput holds the arguments of the send_message tool.
type SendMessageInput struct {
	Message     string `json:"message" jsonschema:"The message content."`
	Title       string `json:"title,omitempty" jsonschema:"The message title. Defaults to the application name."`
	Priority    *int   `json:"priority,omitempty" jsonschema:"The message priority. Higher values are more urgent: 0 is silent, 1-3 low, 4-7 normal, 8-10 high. Defaults to the application default priority."`
	Markdown    bool   `json:"markdown,omitempty" jsonschema:"Whether the message content should be rendered as Markdown."`
	ClickURL    string `json:"click_url,omitempty" jsonschema:"A URL to open when the notification is clicked."`
	BigImageURL string `json:"big_image_url,omitempty" jsonschema:"A URL of an image to show in the notification."`
}

// NewMCP creates a new MCPAPI.
func NewMCP(messages *MessageAPI, version string) *MCPAPI {
	a := &MCPAPI{messages: messages, version: version, schemaCache: mcp.NewSchemaCache()}
	a.handler = mcp.NewStreamableHTTPHandler(a.serverForRequest, &mcp.StreamableHTTPOptions{
		Stateless:    true,
		JSONResponse: true,
		// Authentication is done via application tokens, and gotify is commonly run behind a reverse proxy on localhost.
		DisableLocalhostProtection: true,
	})
	return a
}

// Handle serves the MCP endpoint, the request must be authenticated with an application token.
func (a *MCPAPI) Handle(ctx *gin.Context) {
	app := auth.GetApplication(ctx)
	if app == nil {
		ctx.AbortWithError(403, errors.New("the mcp endpoint requires an application token"))
		return
	}
	req := ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), mcpAppKey{}, app))
	a.handler.ServeHTTP(ctx.Writer, req)
}

func (a *MCPAPI) serverForRequest(req *http.Request) *mcp.Server {
	app, ok := req.Context().Value(mcpAppKey{}).(*model.Application)
	if !ok {
		return nil
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "gotify", Version: a.version}, &mcp.ServerOptions{SchemaCache: a.schemaCache})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_message",
		Description: fmt.Sprintf("Send a push notification via Gotify as the application %q.", app.Name),
		Annotations: &mcp.ToolAnnotations{Title: "Send message", DestructiveHint: new(bool)},
	}, func(_ context.Context, _ *mcp.CallToolRequest, in SendMessageInput) (*mcp.CallToolResult, any, error) {
		return a.sendMessage(app, in)
	})
	return server
}

func (a *MCPAPI) sendMessage(app *model.Application, in SendMessageInput) (*mcp.CallToolResult, any, error) {
	if in.Message == "" {
		return nil, nil, errors.New("message must not be empty")
	}
	msg := &model.CreateMessage{
		Title:    in.Title,
		Message:  in.Message,
		Priority: in.Priority,
		Extras:   mcpExtras(in),
	}
	created, err := a.messages.createMessage(app, msg)
	if err != nil {
		log.Error().Err(err).Uint("app_id", app.ID).Msg("MCP: could not create message")
		return nil, nil, errors.New("failed to send message")
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Message sent (id %d).", created.ID)}},
	}, nil, nil
}

func mcpExtras(in SendMessageInput) map[string]any {
	extras := map[string]any{}
	if in.Markdown {
		extras["client::display"] = map[string]any{"contentType": "text/markdown"}
	}
	notification := map[string]any{}
	if in.ClickURL != "" {
		notification["click"] = map[string]any{"url": in.ClickURL}
	}
	if in.BigImageURL != "" {
		notification["bigImageUrl"] = in.BigImageURL
	}
	if len(notification) > 0 {
		extras["client::notification"] = notification
	}
	if len(extras) == 0 {
		return nil
	}
	return extras
}
