package api

import (
	"github.com/gin-gonic/gin"
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
// ---
// produces: [application/json]
// responses:
//   200:
//     description: Ok
//     schema:
//         $ref: "#/definitions/Health"
//   500:
//     description: Ok
//     schema:
//         $ref: "#/definitions/Health"
func (a *HealthAPI) Health(ctx *gin.Context) {
	if err := a.DB.Ping(); err != nil {
		ctx.JSON(500, model.Health{
			Health:   model.StatusOrange,
			Database: model.StatusRed,
		})
		return
	}
	ctx.JSON(200, model.Health{
		Health:   model.StatusGreen,
		Database: model.StatusGreen,
	})
}
