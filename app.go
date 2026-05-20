package main

import (
	"os"
	"time"

	"github.com/gotify/server/v2/config"
	"github.com/gotify/server/v2/database"
	"github.com/gotify/server/v2/mode"
	"github.com/gotify/server/v2/model"
	"github.com/gotify/server/v2/router"
	"github.com/gotify/server/v2/runner"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Version the version of Gotify.
	Version = "unknown"
	// Commit the git commit hash of this version.
	Commit = "unknown"
	// BuildDate the date on which this binary was build.
	BuildDate = "unknown"
	// Mode the build mode.
	Mode = mode.Dev
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: noColor()})

	vInfo := &model.VersionInfo{Version: Version, Commit: Commit, BuildDate: BuildDate}
	mode.Set(Mode)

	log.Info().Str("version", vInfo.Version).Str("build_date", BuildDate).Msg("Gotify")
	conf := config.Get()

	if conf.PluginsDir != "" {
		if err := os.MkdirAll(conf.PluginsDir, 0o755); err != nil {
			panic(err)
		}
	}
	if err := os.MkdirAll(conf.UploadedImagesDir, 0o755); err != nil {
		panic(err)
	}

	db, err := database.New(conf.Database.Dialect, conf.Database.Connection, conf.DefaultUser.Name, conf.DefaultUser.Pass, conf.PassStrength, true, time.Now)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	engine, closeable := router.Create(db, vInfo, conf)
	defer closeable()

	if err := runner.Run(engine, conf); err != nil {
		log.Error().Err(err).Msg("Server error")
		os.Exit(1)
	}
}

func noColor() bool {
	// https://no-color.org/
	if os.Getenv("NO_COLOR") == "1" {
		return true
	}

	isTTY := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())
	return !isTTY
}
