package external

import (
	"fmt"
	"net/http"
	"plugin"

	"github.com/gotify/server/model"
)

// Compat is external authenticator compat
type Compat struct {
	name               string
	authenticationfunc func(req *http.Request) (user *model.User, err error)
}

// Name implements Authenticator
func (c *Compat) Name() string {
	return c.name
}

// Authenticate implements Authenticator
func (c *Compat) Authenticate(req *http.Request) (user *model.User, err error) {
	return c.authenticationfunc(req)
}

type pluginOpenFailedError struct {
	innerErr error
}

func (e pluginOpenFailedError) Error() string {
	return fmt.Sprintf("externalauth: plugin open: %v", e.innerErr)
}

type symbolNotFoundError struct {
	symbol string
	path   string
}

func (s symbolNotFoundError) Error() string {
	return fmt.Sprintf("externalauth: symbol '%s' is not found in plugin %s", s.symbol, s.path)
}

type symbolTypeError struct {
	symbol   string
	path     string
	got      interface{}
	expected interface{}
}

func (s symbolTypeError) Error() string {
	return fmt.Sprintf("externalauth: unexpected type of symbol '%s' in plugin %s (got %t expected %t)", s.symbol, s.path, s.got, s.expected)
}

// LoadAuthenticatorPlugin loads an external authentication plugin with given name
func LoadAuthenticatorPlugin(path string, name string) (*Compat, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, pluginOpenFailedError{err}
	}
	authenticateSym, err := p.Lookup("Authenticate")
	if err != nil {
		return nil, symbolNotFoundError{"Authenticate", path}
	}
	authenticateFunc, ok := authenticateSym.(func(req *http.Request) (user *model.User, err error))
	if !ok {
		return nil, symbolTypeError{"Authenticate", path, authenticateSym, (func(req *http.Request) (user *model.User, err error))(nil)}
	}
	return &Compat{
		name:               name,
		authenticationfunc: authenticateFunc,
	}, nil
}
