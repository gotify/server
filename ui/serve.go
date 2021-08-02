package ui

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/gotify/server/v2/model"
)

var box = packr.New("ui", "../ui/build")

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
	ui := r.Group("/", gzip.Gzip(gzip.DefaultCompression))
	ui.GET("/", serveFile("index.html", "text/html", func(content string) string {
		return strings.Replace(content, "%CONFIG%", string(uiConfigBytes), 1)
	}))
	ui.GET("/index.html", serveFile("index.html", "text/html", noop))
	ui.GET("/manifest.json", serveFile("manifest.json", "application/json", noop))
	ui.GET("/asset-manifest.json", serveFile("asset-manifest.json", "application/json", noop))
	ui.GET("/static/*any", gin.WrapH(http.FileServer(box)))
}

func noop(s string) string {
	return s
}

func serveFile(name, contentType string, convert func(string) string) gin.HandlerFunc {
	content, err := box.FindString(name)
	if err != nil {
		panic(err)
	}
	converted := convert(content)
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", contentType)
		ctx.String(200, converted)
	}
}
