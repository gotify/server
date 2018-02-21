package main

import (
	"math/rand"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/config"
	"github.com/gotify/server/database"
	"github.com/gotify/server/model"
	"github.com/gotify/server/router"
	"github.com/gotify/server/runner"
)

var (
	// Version the version of Gotify.
	Version = "unknown"
	// Commit the git commit hash of this version.
	Commit = "unknown"
	// BuildDate the date on which this binary was build.
	BuildDate = "unknown"
	// Branch the git branch of this version.
	Branch = "unknown"
)

func main() {
	vInfo := &model.VersionInfo{Version: Version, Commit: Commit, BuildDate: BuildDate, Branch: Branch}

	fmt.Println("Starting Gotify version", vInfo.Version+"@"+BuildDate)
	rand.Seed(time.Now().UnixNano())
	conf := config.Get()
	db, err := database.New(conf.Database.Dialect, conf.Database.Connection, conf.DefaultUser.Name, conf.DefaultUser.Pass)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	engine, closeable := router.Create(db, vInfo)
	defer closeable()

	runner.Run(engine, conf)
}
