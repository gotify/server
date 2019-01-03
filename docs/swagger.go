package docs

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/gotify/location"
)

// Serve serves the documentation.
func Serve(ctx *gin.Context) {
	url := location.Get(ctx)
	ctx.Writer.WriteString(get(url.Host))
}

func get(host string) string {
	box := packr.NewBox("./")
	return strings.Replace(box.String("spec.json"), "localhost", host, 1)
}
