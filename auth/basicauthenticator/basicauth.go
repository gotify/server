package basicauthenticator

import (
	"net/http"

	"github.com/gotify/server/auth/basicauthenticator/password"
	"github.com/gotify/server/model"
)

// The Database interface for encapsulating database access.
type Database interface {
	GetUserByName(name string) (*model.User, error)
}

// AuthProvider is the auth provider for HTTP Basic Auth authentication using internal credential managemant
type AuthProvider struct {
	DB Database
}

// Authenticate implements auth.AuthenticationProvider
func (a *AuthProvider) Authenticate(req *http.Request) (user *model.User, err error) {
	if name, pass, ok := req.BasicAuth(); ok {
		if user, err := a.DB.GetUserByName(name); err != nil {
			return nil, nil
		} else if user != nil && password.ComparePassword(user.Pass, []byte(pass)) {
			return user, nil
		}
	}
	return nil, nil
}

// Name implements auth.AuthenticationProvider
func (a *AuthProvider) Name() string {
	return "internal"
}
