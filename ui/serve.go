package ui

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

var box = packr.New("ui", "../ui/build")

// Register registers the ui on the root path.
func Register(r *gin.Engine) {
	ui := r.Group("/", gzip.Gzip(gzip.DefaultCompression))
	ui.GET("/", serveFile("index.html", "text/html"))
	ui.GET("/index.html", serveFile("index.html", "text/html"))
	ui.GET("/manifest.json", serveFile("manifest.json", "application/json"))
	ui.GET("/service-worker.js", serveFile("service-worker.js", "text/javascript"))
	ui.GET("/assets-manifest.json", serveFile("asserts-manifest.json", "application/json"))
	ui.GET("/static/*any", gin.WrapH(http.FileServer(box)))
}

func serveFile(name, contentType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", contentType)
		content, err := box.FindString(name)
		if err != nil {
			panic(err)
		}
		ctx.String(200, content)
	}
}
