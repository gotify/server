package ui

import (
	"embed"
	"encoding/json"
	"io/fs"
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

// Register registers the ui on the root path.
func Register(r *gin.Engine, version model.VersionInfo, register bool) {
	uiConfigBytes, err := json.Marshal(uiConfig{Version: version, Register: register})
	if err != nil {
		panic(err)
	}

	replaceConfig := func(content string) string {
		return strings.Replace(content, "%CONFIG%", string(uiConfigBytes), 1)
	}

	ui := r.Group("/", gzip.Gzip(gzip.DefaultCompression))
	ui.GET("/", serveFile("index.html", "text/html", replaceConfig))
	ui.GET("/index.html", serveFile("index.html", "text/html", replaceConfig))
	ui.GET("/manifest.json", serveFile("manifest.json", "application/json", noop))
	ui.GET("/asset-manifest.json", serveFile("asset-manifest.json", "application/json", noop))

	subBox, err := fs.Sub(box, "build")
	if err != nil {
		panic(err)
	}
	ui.GET("/static/*any", gin.WrapH(http.FileServer(http.FS(subBox))))
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
