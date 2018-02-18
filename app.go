package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/config"
	"github.com/gotify/server/database"
	"github.com/gotify/server/router"
	"github.com/gotify/server/runner"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	conf := config.Get()
	db, err := database.New(conf.Database.Dialect, conf.Database.Connection, conf.DefaultUser.Name, conf.DefaultUser.Pass)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	engine, closeable := router.Create(db)
	defer closeable()

	runner.Run(engine, conf)
}
