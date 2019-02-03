package plugin

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func requirePluginEnabled(id uint, db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		if conf := db.GetPluginConfByID(id); conf == nil || !conf.Enabled {
			c.AbortWithError(400, errors.New("plugin is disabled"))
		}
	}
}
