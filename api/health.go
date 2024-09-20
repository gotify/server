package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/gotify/server/v2/model"
)

// The HealthDatabase interface for encapsulating database access.
type HealthDatabase interface {
	Ping() error
}

// The HealthAPI provides handlers for the health information.
type HealthAPI struct {
	DB HealthDatabase
}

// Health returns health information.
// swagger:operation GET /health health getHealth
//
// Get health information.
//
//	---
//	produces: [application/json]
//	responses:
//	  200:
//	    description: Ok
//	    schema:
//	        $ref: "#/definitions/Health"
//	  500:
//	    description: Ok
//	    schema:
//	        $ref: "#/definitions/Health"
func (a *HealthAPI) Health(ctx *gin.Context) {
	var response *model.Health
	if err := a.DB.Ping(); err != nil {
		response = &model.Health{
			Health:   model.StatusOrange,
			Database: model.StatusRed,
		}
	} else {
		response = &model.Health{
			Health:   model.StatusGreen,
			Database: model.StatusGreen,
		}
	}

	status := 500
	if response.Database == model.StatusGreen || response.Health == model.StatusGreen {
		status = 200
	}

	renderer := render.JSON{Data: response}

	// in case this is called from a non http request
	if ctx.Request == nil {
		ctx.Render(status, renderer)
		return
	}

	switch ctx.Request.Method {
	case "HEAD":
		renderer.WriteContentType(ctx.Writer)
		ctx.Status(status)
	case "GET":
		ctx.Render(status, renderer)
	default:
		ctx.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}
