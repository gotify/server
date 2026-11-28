package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/gotify/server/v3/auth/password"
	"github.com/gotify/server/v3/database"
	"github.com/gotify/server/v3/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserCLI(t *testing.T) *database.GormDatabase {
	dbFile := filepath.Join(t.TempDir(), "gotify.db")
	t.Setenv("GOTIFY_CONFIG_FILE", filepath.Join(t.TempDir(), "missing.env"))
	t.Setenv("GOTIFY_DATABASE_DIALECT", "sqlite3")
	t.Setenv("GOTIFY_DATABASE_CONNECTION", dbFile)
	t.Setenv("GOTIFY_PASSSTRENGTH", "4")

	db, err := database.New("sqlite3", dbFile, "admin", "admin", 4, true, time.Now)
	require.NoError(t, err)
	t.Cleanup(db.Close)
	return db
}

func fakeTerminal(t *testing.T, isTerminal bool, inputs ...string) {
	oldIsTerminal, oldRead := stdinIsTerminal, readPassword
	t.Cleanup(func() { stdinIsTerminal, readPassword = oldIsTerminal, oldRead })
	stdinIsTerminal = func() bool { return isTerminal }
	readPassword = func() ([]byte, error) {
		if len(inputs) == 0 {
			return nil, errors.New("unexpected prompt")
		}
		next := inputs[0]
		inputs = inputs[1:]
		return []byte(next), nil
	}
}

func getUser(t *testing.T, db *database.GormDatabase, name string) *model.User {
	user, err := db.GetUserByName(name)
	require.NoError(t, err)
	return user
}

func TestRunUser(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantCode int
		stdout   string
		stderr   string
	}{
		{"no command", []string{"user"}, 2, "", "Usage: gotify user"},
		{"help", []string{"user", "help"}, 0, "Usage: gotify user", ""},
		{"unknown command", []string{"user", "bogus"}, 2, "", `unknown command "bogus"`},
		{"create without name", []string{"user", "create"}, 2, "", "Usage: gotify user"},
		{"create too many args", []string{"user", "create", "a", "b", "c"}, 2, "", "Usage: gotify user"},
		{"create unknown flag", []string{"user", "create", "--nope", "a"}, 2, "", "not defined"},
		{"create existing", []string{"user", "create", "admin", "pw"}, 1, "", `user "admin" already exists`},
		{"set-password without name", []string{"user", "set-password"}, 2, "", "Usage: gotify user"},
		{"set-password missing user", []string{"user", "set-password", "nobody", "pw"}, 1, "", `user "nobody" does not exist`},
		{"set-admin without value", []string{"user", "set-admin", "admin"}, 2, "", "Usage: gotify user"},
		{"set-admin invalid value", []string{"user", "set-admin", "admin", "yes"}, 2, "", "Usage: gotify user"},
		{"set-admin missing user", []string{"user", "set-admin", "nobody", "true"}, 1, "", `user "nobody" does not exist`},
		{"password too long", []string{"user", "create", "alice", string(make([]byte, 73))}, 1, "", "password length exceeds 72 bytes"},
	}
	setupUserCLI(t)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(c.args, &stdout, &stderr)
			assert.Equal(t, c.wantCode, code)
			assert.Contains(t, stdout.String(), c.stdout)
			assert.Contains(t, stderr.String(), c.stderr)
		})
	}
}

func TestRunUser_create(t *testing.T) {
	db := setupUserCLI(t)

	var stdout, stderr bytes.Buffer
	assert.Equal(t, 0, run([]string{"user", "create", "alice", "secret"}, &stdout, &stderr), stderr.String())
	alice := getUser(t, db, "alice")
	require.NotNil(t, alice)
	assert.False(t, alice.Admin)
	assert.True(t, password.ComparePassword(alice.Pass, []byte("secret")))

	assert.Equal(t, 0, run([]string{"user", "create", "--admin", "bob", "secret"}, &stdout, &stderr), stderr.String())
	bob := getUser(t, db, "bob")
	require.NotNil(t, bob)
	assert.True(t, bob.Admin)
}

func TestRunUser_create_prompt(t *testing.T) {
	db := setupUserCLI(t)
	fakeTerminal(t, true, "prompted", "prompted")

	var stdout, stderr bytes.Buffer
	assert.Equal(t, 0, run([]string{"user", "create", "alice"}, &stdout, &stderr), stderr.String())
	assert.Contains(t, stderr.String(), "Confirm password: ")
	alice := getUser(t, db, "alice")
	require.NotNil(t, alice)
	assert.True(t, password.ComparePassword(alice.Pass, []byte("prompted")))
}

func TestRunUser_prompt_mismatch(t *testing.T) {
	db := setupUserCLI(t)
	fakeTerminal(t, true, "one", "two")

	var stdout, stderr bytes.Buffer
	assert.Equal(t, 1, run([]string{"user", "create", "alice"}, &stdout, &stderr))
	assert.Contains(t, stderr.String(), "passwords do not match")
	assert.Nil(t, getUser(t, db, "alice"))
}

func TestRunUser_prompt_empty(t *testing.T) {
	setupUserCLI(t)
	fakeTerminal(t, true, "", "")

	var stdout, stderr bytes.Buffer
	assert.Equal(t, 1, run([]string{"user", "set-password", "admin"}, &stdout, &stderr))
	assert.Contains(t, stderr.String(), "password must not be empty")
}

func TestRunUser_prompt_noTerminal(t *testing.T) {
	setupUserCLI(t)
	fakeTerminal(t, false)

	var stdout, stderr bytes.Buffer
	assert.Equal(t, 1, run([]string{"user", "set-password", "admin"}, &stdout, &stderr))
	assert.Contains(t, stderr.String(), "stdin is not a terminal")
}

func TestRunUser_setPassword(t *testing.T) {
	db := setupUserCLI(t)

	var stdout, stderr bytes.Buffer
	assert.Equal(t, 0, run([]string{"user", "set-password", "admin", "newpw"}, &stdout, &stderr), stderr.String())
	admin := getUser(t, db, "admin")
	assert.True(t, password.ComparePassword(admin.Pass, []byte("newpw")))
	assert.False(t, password.ComparePassword(admin.Pass, []byte("admin")))
	assert.True(t, admin.Admin)

	fakeTerminal(t, true, "prompted", "prompted")
	assert.Equal(t, 0, run([]string{"user", "set-password", "admin"}, &stdout, &stderr), stderr.String())
	admin = getUser(t, db, "admin")
	assert.True(t, password.ComparePassword(admin.Pass, []byte("prompted")))
}

func TestRunUser_setAdmin(t *testing.T) {
	db := setupUserCLI(t)
	require.NoError(t, db.CreateUser(&model.User{Name: "alice"}))

	var stdout, stderr bytes.Buffer
	assert.Equal(t, 0, run([]string{"user", "set-admin", "alice", "true"}, &stdout, &stderr), stderr.String())
	assert.True(t, getUser(t, db, "alice").Admin)

	assert.Equal(t, 0, run([]string{"user", "set-admin", "alice", "false"}, &stdout, &stderr), stderr.String())
	assert.False(t, getUser(t, db, "alice").Admin)
}

func TestRunUser_doesNotCreateDefaultUser(t *testing.T) {
	dbFile := filepath.Join(t.TempDir(), "gotify.db")
	t.Setenv("GOTIFY_CONFIG_FILE", filepath.Join(t.TempDir(), "missing.env"))
	t.Setenv("GOTIFY_DATABASE_DIALECT", "sqlite3")
	t.Setenv("GOTIFY_DATABASE_CONNECTION", dbFile)
	t.Setenv("GOTIFY_PASSSTRENGTH", "4")

	var stdout, stderr bytes.Buffer
	assert.Equal(t, 0, run([]string{"user", "create", "--admin", "bob", "secret"}, &stdout, &stderr), stderr.String())

	db, err := database.New("sqlite3", dbFile, "admin", "admin", 4, false, time.Now)
	require.NoError(t, err)
	defer db.Close()
	count, err := db.CountUser()
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Nil(t, getUser(t, db, "admin"))
}
