package plugin

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func requirePluginEnabled(id uint, db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		conf, err := db.GetPluginConfByID(id)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		if conf == nil || !conf.Enabled {
			c.AbortWithError(400, errors.New("plugin is disabled"))
		}
	}
}
