package docs

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/gotify/location"
)

// Serve serves the documentation.
func Serve(ctx *gin.Context) {
	base := location.Get(ctx).Host
	if bQ := ctx.Query("base"); bQ != "" {
		base = bQ
	}
	ctx.Writer.WriteString(get(base))
}

func get(base string) string {
	box := packr.NewBox("./")
	return strings.Replace(box.String("spec.json"), "localhost", base, 1)
}
