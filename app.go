package main

import (
	"math/rand"
	"time"

	"fmt"

	"github.com/gotify/server/config"
	"github.com/gotify/server/database"
	"github.com/gotify/server/model"
	"github.com/gotify/server/router"
	"github.com/gotify/server/runner"
	"github.com/gotify/server/mode"
)

var (
	// Version the version of Gotify.
	Version = "unknown"
	// Commit the git commit hash of this version.
	Commit = "unknown"
	// BuildDate the date on which this binary was build.
	BuildDate = "unknown"
	// Mode the build mode
	Mode = mode.Dev
)

func main() {
	vInfo := &model.VersionInfo{Version: Version, Commit: Commit, BuildDate: BuildDate}
	mode.Set(Mode);

	fmt.Println("Starting Gotify version", vInfo.Version+"@"+BuildDate)
	rand.Seed(time.Now().UnixNano())
	conf := config.Get()
	db, err := database.New(conf.Database.Dialect, conf.Database.Connection, conf.DefaultUser.Name, conf.DefaultUser.Pass, conf.PassStrength)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	engine, closeable := router.Create(db, vInfo, conf)
	defer closeable()

	runner.Run(engine, conf)
}
