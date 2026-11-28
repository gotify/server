package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gotify/server/v3/auth/password"
	"github.com/gotify/server/v3/config"
	"github.com/gotify/server/v3/database"
	"github.com/gotify/server/v3/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

var (
	stdinIsTerminal = func() bool { return term.IsTerminal(int(os.Stdin.Fd())) }
	readPassword    = func() ([]byte, error) { return term.ReadPassword(int(os.Stdin.Fd())) }
)

func runUser(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUserUsage(stderr)
		return 2
	}

	switch args[0] {
	case "create":
		fs := flag.NewFlagSet("gotify user create", flag.ContinueOnError)
		fs.SetOutput(stderr)
		admin := fs.Bool("admin", false, "Create the user as admin.")
		fs.Usage = func() { printUserUsage(stderr) }
		if err := fs.Parse(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return 0
			}
			return 2
		}
		if fs.NArg() < 1 || fs.NArg() > 2 {
			printUserUsage(stderr)
			return 2
		}
		return withUserDB(stderr, func(db *database.GormDatabase, conf *config.Configuration) error {
			name := fs.Arg(0)
			existing, err := db.GetUserByName(name)
			if err != nil {
				return err
			}
			if existing != nil {
				return fmt.Errorf("user %q already exists", name)
			}
			pass, err := getNewPassword(fs.Arg(1), stderr, conf.PassStrength)
			if err != nil {
				return err
			}
			if err := db.CreateUser(&model.User{Name: name, Pass: pass, Admin: *admin}); err != nil {
				return err
			}
			fmt.Fprintf(stdout, "Created user %q (admin: %t)\n", name, *admin)
			return nil
		})
	case "set-password":
		if len(args) < 2 || len(args) > 3 {
			printUserUsage(stderr)
			return 2
		}
		return withUserDB(stderr, func(db *database.GormDatabase, conf *config.Configuration) error {
			user, err := getExistingUser(db, args[1])
			if err != nil {
				return err
			}
			pwArg := ""
			if len(args) == 3 {
				pwArg = args[2]
			}
			pass, err := getNewPassword(pwArg, stderr, conf.PassStrength)
			if err != nil {
				return err
			}
			user.Pass = pass
			if err := db.UpdateUser(user); err != nil {
				return err
			}
			fmt.Fprintf(stdout, "Updated password of user %q\n", user.Name)
			return nil
		})
	case "set-admin":
		if len(args) != 3 || (args[2] != "true" && args[2] != "false") {
			printUserUsage(stderr)
			return 2
		}
		admin := args[2] == "true"
		return withUserDB(stderr, func(db *database.GormDatabase, conf *config.Configuration) error {
			user, err := getExistingUser(db, args[1])
			if err != nil {
				return err
			}
			user.Admin = admin
			if err := db.UpdateUser(user); err != nil {
				return err
			}
			fmt.Fprintf(stdout, "Set admin of user %q to %t\n", user.Name, admin)
			return nil
		})
	case "help", "-h", "-help", "--help":
		printUserUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "gotify user: unknown command %q\n\n", args[0])
		printUserUsage(stderr)
		return 2
	}
}

// withUserDB loads the server configuration and opens the configured database.
// Unlike serve, it never creates the default user, so the command only does what was asked.
func withUserDB(stderr io.Writer, f func(db *database.GormDatabase, conf *config.Configuration) error) int {
	conf, futureLogs := config.Get()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: stderr, TimeFormat: time.RFC3339, NoColor: noColor(conf.NoColor)}).Level(zerolog.Level(conf.LogLevel))

	exit := false
	for _, futureLog := range futureLogs {
		log.WithLevel(futureLog.Level).Msg(futureLog.Msg)
		exit = exit || futureLog.Level == zerolog.FatalLevel || futureLog.Level == zerolog.PanicLevel
	}
	if exit {
		return 1
	}

	db, err := database.New(conf.Database.Dialect, conf.Database.Connection, conf.DefaultUser.Name, conf.DefaultUser.Pass, conf.PassStrength, false, time.Now)
	if err != nil {
		fmt.Fprintln(stderr, "Cannot initialize database:", err)
		return 1
	}
	defer db.Close()

	if err := f(db, conf); err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	return 0
}

func getExistingUser(db *database.GormDatabase, name string) (*model.User, error) {
	user, err := db.GetUserByName(name)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user %q does not exist", name)
	}
	return user, nil
}

// getNewPassword validates and hashes pw, prompting for it if it is empty.
func getNewPassword(pw string, stderr io.Writer, strength int) ([]byte, error) {
	if pw == "" {
		var err error
		if pw, err = promptPassword(stderr); err != nil {
			return nil, err
		}
	}
	if err := password.ValidateNewPassword(pw); err != nil {
		return nil, err
	}
	return password.CreatePassword(pw, strength)
}

func promptPassword(stderr io.Writer) (string, error) {
	if !stdinIsTerminal() {
		return "", errors.New("no password given and stdin is not a terminal, pass the password as argument or run the command in a terminal (e.g. docker exec -it)")
	}
	fmt.Fprint(stderr, "Password: ")
	pw, err := readPassword()
	fmt.Fprintln(stderr)
	if err != nil {
		return "", err
	}
	fmt.Fprint(stderr, "Confirm password: ")
	confirm, err := readPassword()
	fmt.Fprintln(stderr)
	if err != nil {
		return "", err
	}
	if string(pw) != string(confirm) {
		return "", errors.New("passwords do not match")
	}
	return string(pw), nil
}

func printUserUsage(w io.Writer) {
	fmt.Fprint(w, `Usage: gotify user <command> [arguments]

Manage users directly in the configured database, e.g. to recover a lost
admin account. If [password] is omitted, it is prompted for.

Commands:
  create [--admin] <name> [password]   Create a new user.
  set-password <name> [password]       Change the password of a user.
  set-admin <name> true|false          Grant or revoke admin permissions.
`)
}
