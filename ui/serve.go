package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
)

var box = packr.NewBox("../ui/build")

// Register registers the ui on the root path.
func Register(r *gin.Engine) {
	r.GET("/", serveFile("index.html", "text/html"))
	r.GET("/index.html", serveFile("index.html", "text/html"))
	r.GET("/manifest.json", serveFile("manifest.json", "application/json"))
	r.GET("/service-worker.js", serveFile("service-worker.js", "text/javascript"))
	r.GET("/assets-manifest.json", serveFile("asserts-manifest.json", "application/json"))
	r.GET("/static/*any", gin.WrapH(http.FileServer(box)))
}

func serveFile(name, contentType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", contentType)
		ctx.String(200, box.String(name))
	}
}
