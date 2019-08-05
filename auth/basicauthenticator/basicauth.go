package basicauthenticator

import (
	"net/http"

	"github.com/gotify/server/auth/basicauthenticator/password"
	"github.com/gotify/server/model"
)

type Database interface {
	GetUserByName(name string) (*model.User, error)
}

type AuthProvider struct {
	DB Database
}

func (a *AuthProvider) Authenticate(req *http.Request) (user *model.User, err error) {
	if name, pass, ok := req.BasicAuth(); ok {
		if user, err := a.DB.GetUserByName(name); err != nil {
			return nil, err
		} else if user != nil && password.ComparePassword(user.Pass, []byte(pass)) {
			return user, nil
		}
	}
	return nil, nil
}
