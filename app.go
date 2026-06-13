package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"time"

	"github.com/gotify/server/v2/config"
	"github.com/gotify/server/v2/config/migrate"
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
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	vInfo := &model.VersionInfo{Version: Version, Commit: Commit, BuildDate: BuildDate}
	fs := flag.NewFlagSet("gotify", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { printUsage(stderr) }
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	command := fs.Arg(0)
	switch command {
	case "serve", "":
		return serve(vInfo)
	case "version":
		fmt.Fprintln(stdout, "Version:", vInfo.Version)
		fmt.Fprintln(stdout, "Commit:", vInfo.Commit)
		fmt.Fprintln(stdout, "Build Date:", vInfo.BuildDate)
		fmt.Fprintln(stdout, "Go Build Info:")
		b, ok := debug.ReadBuildInfo()
		if ok {
			fmt.Fprintln(stdout, b)
		}
		return 0
	case "migrate-config":
		content, err := migrate.Config(fs.Arg(1))
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, content)
		return 0
	default:
		if command != "" {
			fmt.Fprintf(stderr, "gotify: unknown command %q\n\n", command)
		}
		printUsage(stderr)
		return 2
	}
}

func serve(vInfo *model.VersionInfo) int {
	mode.Set(Mode)

	conf, futureLogs := config.Get()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: noColor(conf.NoColor)}).Level(zerolog.Level(conf.LogLevel))
	log.Info().Str("version", vInfo.Version).Str("build_date", BuildDate).Msg("Gotify")

	exit := false
	for _, futureLog := range futureLogs {
		log.WithLevel(futureLog.Level).Msg(futureLog.Msg)
		exit = exit || futureLog.Level == zerolog.FatalLevel || futureLog.Level == zerolog.PanicLevel
	}
	if exit {
		return 1
	}

	if conf.PluginsDir != "" {
		if err := os.MkdirAll(conf.PluginsDir, 0o755); err != nil {
			log.Error().Err(err).Str("dir", conf.PluginsDir).Msg("Cannot create plugins directory")
			return 1
		}
	}
	if err := os.MkdirAll(conf.UploadedImagesDir, 0o755); err != nil {
		log.Error().Err(err).Str("dir", conf.UploadedImagesDir).Msg("Cannot create uploaded images directory")
		return 1
	}

	db, err := database.New(conf.Database.Dialect, conf.Database.Connection, conf.DefaultUser.Name, conf.DefaultUser.Pass, conf.PassStrength, true, time.Now)
	if err != nil {
		log.Error().Err(err).Msg("Cannot initialize database")
		return 1
	}
	defer db.Close()

	engine, closeable := router.Create(db, vInfo, conf)
	defer closeable()

	if err := runner.Run(engine, conf); err != nil {
		log.Error().Err(err).Msg("Server error")
		return 1
	}
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `Usage: gotify [flags] <command> [arguments]

Commands:
  serve                       Start the Gotify server.
  migrate-config <file.yml>   Convert an old YAML config file to the new env
                              format and print it to stdout.
  version                     Show version information
`)
}

func noColor(noColorEnv string) bool {
	// https://no-color.org/
	if noColorEnv == "1" {
		return true
	}

	isTTY := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())
	return !isTTY
}
