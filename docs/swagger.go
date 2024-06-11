package docs

import (
	_ "embed"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotify/location"
)

//go:embed spec.json
var spec string

// Serve serves the documentation.
func Serve(ctx *gin.Context) {
	base := location.Get(ctx).Host
	if basePathFromQuery := ctx.Query("base"); basePathFromQuery != "" {
		base = basePathFromQuery
	}
	ctx.Writer.WriteString(getSwaggerJSON(base))
}

func getSwaggerJSON(base string) string {
	return strings.Replace(spec, "localhost", base, 1)
}
