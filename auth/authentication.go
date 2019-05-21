package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/auth/password"
	"github.com/gotify/server/model"
)

const (
	headerName = "X-Gotify-Key"
)

// The Database interface for encapsulating database access.
type Database interface {
	GetApplicationByToken(token string) (*model.Application, error)
	GetClientByToken(token string) (*model.Client, error)
	GetPluginConfByToken(token string) (*model.PluginConf, error)
	GetUserByName(name string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
}

// Auth is the provider for authentication middleware
type Auth struct {
	DB Database
}

type authenticate func(tokenID string, user *model.User) (authenticated bool, success bool, userId uint)

// RequireAdmin returns a gin middleware which requires a client token or basic authentication header to be supplied
// with the request. Also the authenticated user must be an administrator.
func (a *Auth) RequireAdmin() gin.HandlerFunc {
	return a.requireToken(func(tokenID string, user *model.User) (bool, bool, uint) {
		if user != nil {
			return true, user.Admin, user.ID
		}
		if token, _ := a.DB.GetClientByToken(tokenID); token != nil {
			user, _ := a.DB.GetUserByID(token.UserID)
			return true, user.Admin, token.UserID
		}
		return false, false, 0
	})
}

// RequireClient returns a gin middleware which requires a client token or basic authentication header to be supplied
// with the request.
func (a *Auth) RequireClient() gin.HandlerFunc {
	return a.requireToken(func(tokenID string, user *model.User) (bool, bool, uint) {
		if user != nil {
			return true, true, user.ID
		}
		if token, _ := a.DB.GetClientByToken(tokenID); token != nil {
			return true, true, token.UserID
		}
		return false, false, 0
	})
}

// RequireApplicationToken returns a gin middleware which requires an application token to be supplied with the request.
func (a *Auth) RequireApplicationToken() gin.HandlerFunc {
	return a.requireToken(func(tokenID string, user *model.User) (bool, bool, uint) {
		if user != nil {
			return true, false, 0
		}
		if token, _ := a.DB.GetApplicationByToken(tokenID); token != nil {
			return true, true, token.UserID
		}
		return false, false, 0
	})
}

func (a *Auth) tokenFromQueryOrHeader(ctx *gin.Context) string {
	if token := a.tokenFromQuery(ctx); token != "" {
		return token
	} else if token := a.tokenFromHeader(ctx); token != "" {
		return token
	}
	return ""
}

func (a *Auth) tokenFromQuery(ctx *gin.Context) string {
	return ctx.Request.URL.Query().Get("token")
}

func (a *Auth) tokenFromHeader(ctx *gin.Context) string {
	return ctx.Request.Header.Get(headerName)
}

func (a *Auth) userFromBasicAuth(ctx *gin.Context) *model.User {
	if name, pass, ok := ctx.Request.BasicAuth(); ok {
		if user, _ := a.DB.GetUserByName(name); user != nil && password.ComparePassword(user.Pass, []byte(pass)) {
			return user
		}
	}
	return nil
}

func (a *Auth) requireToken(auth authenticate) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := a.tokenFromQueryOrHeader(ctx)
		user := a.userFromBasicAuth(ctx)

		if user != nil || token != "" {
			if authenticated, ok, userID := auth(token, user); ok {
				RegisterAuthentication(ctx, user, userID, token)
				ctx.Next()
				return
			} else if authenticated {
				ctx.AbortWithError(403, errors.New("you are not allowed to access this api"))
				return
			}
		}
		ctx.AbortWithError(401, errors.New("you need to provide a valid access token or user credentials to access this api"))
	}
}
