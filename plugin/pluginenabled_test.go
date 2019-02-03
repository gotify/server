package plugin

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotify/server/model"
	"github.com/gotify/server/test/testdb"

	"github.com/gin-gonic/gin"
)

func TestRequirePluginEnabled(t *testing.T) {

	db := testdb.NewDBWithDefaultUser(t)
	conf := &model.PluginConf{
		ID:      1,
		UserID:  1,
		Enabled: true,
	}
	db.CreatePluginConf(conf)

	g := gin.New()

	mux := g.Group("/", requirePluginEnabled(1, db))

	mux.GET("/", func(c *gin.Context) {
		c.Status(200)
	})

	getCode := func() int {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		g.ServeHTTP(w, r)
		return w.Code
	}

	assert.Equal(t, 200, getCode())

	conf.Enabled = false
	db.UpdatePluginConf(conf)
	assert.Equal(t, 400, getCode())
}
