package ui

import (
	"embed"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/model"
)

//go:embed build/*
var box embed.FS

type uiConfig struct {
	Register bool              `json:"register"`
	Version  model.VersionInfo `json:"version"`
}

type AddPrefix struct {
	Prefix       string
	InnerHandler http.Handler
}

func (a *AddPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = a.Prefix + r.URL.Path
	a.InnerHandler.ServeHTTP(w, r)
}

// Register registers the ui on the root path.
func Register(r *gin.Engine, version model.VersionInfo, register bool) {
	uiConfigBytes, err := json.Marshal(uiConfig{Version: version, Register: register})
	if err != nil {
		panic(err)
	}
	ui := r.Group("/", gzip.Gzip(gzip.DefaultCompression))
	ui.GET("/", serveFile("index.html", "text/html", func(content string) string {
		return strings.Replace(content, "%CONFIG%", string(uiConfigBytes), 1)
	}))
	ui.GET("/index.html", serveFile("index.html", "text/html", noop))
	ui.GET("/manifest.json", serveFile("manifest.json", "application/json", noop))
	ui.GET("/asset-manifest.json", serveFile("asset-manifest.json", "application/json", noop))
	ui.GET("/static/*any", gin.WrapH(&AddPrefix{
		Prefix:       "/build",
		InnerHandler: http.FileServer(http.FS(box)),
	}))
}

func noop(s string) string {
	return s
}

func serveFile(name, contentType string, convert func(string) string) gin.HandlerFunc {
	content, err := box.ReadFile("build/" + name)
	if err != nil {
		panic(err)
	}
	converted := convert(string(content))
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", contentType)
		ctx.String(200, converted)
	}
}
