package docs

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/gotify/location"
)

var box = packr.New("docs", "./")

// Serve serves the documentation.
func Serve(ctx *gin.Context) {
	base := location.Get(ctx).Host
	if basePathFromQuery := ctx.Query("base"); basePathFromQuery != "" {
		base = basePathFromQuery
	}
	ctx.Writer.WriteString(get(base))
}

func get(base string) string {
	spec, err := box.FindString("spec.json")
	if err != nil {
		panic(err)
	}
	return strings.Replace(spec, "localhost", base, 1)
}
