package docs

import (
	"embed"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotify/location"
)

//go:embed spec.json
var box embed.FS

// Serve serves the documentation.
func Serve(ctx *gin.Context) {
	base := location.Get(ctx).Host
	if basePathFromQuery := ctx.Query("base"); basePathFromQuery != "" {
		base = basePathFromQuery
	}
	ctx.Writer.WriteString(get(base))
}

func get(base string) string {
	spec, err := box.ReadFile("spec.json")
	if err != nil {
		panic(err)
	}
	return strings.Replace(string(spec), "localhost", base, 1)
}
