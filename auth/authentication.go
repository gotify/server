package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/model"
	"strings"
)

const (
	headerName   = "Authorization"
	headerSchema = "ApiKey "
)

// The Database interface for encapsulating database access.
type Database interface {
	GetApplicationByID(id string) *model.Application
	GetClientByID(id string) *model.Client
	GetUserByName(name string) *model.User
	GetUserByID(id uint) *model.User
}

// Auth is the provider for authentication middleware
type Auth struct {
	DB Database
}

type authenticate func(tokenID string, user *model.User) (success bool, userId uint)

// RequireAdmin returns a gin middleware which requires a client token or basic authentication header to be supplied
// with the request. Also the authenticated user must be an administrator.
func (a *Auth) RequireAdmin() gin.HandlerFunc {
	return a.requireToken(func(tokenID string, user *model.User) (bool, uint) {
		if user != nil {
			return user.Admin, user.ID
		}
		if token := a.DB.GetClientByID(tokenID); token != nil {
			return a.DB.GetUserByID(token.UserID).Admin, token.UserID
		}
		return false, 0
	})
}

// RequireClient returns a gin middleware which requires a client token or basic authentication header to be supplied
// with the request.
func (a *Auth) RequireClient() gin.HandlerFunc {
	return a.requireToken(func(tokenID string, user *model.User) (bool, uint) {
		if user != nil {
			return true, user.ID
		}
		if token := a.DB.GetClientByID(tokenID); token != nil {
			return true, token.UserID
		}
		return false, 0
	})
}

// RequireApplicationToken returns a gin middleware which requires an application token to be supplied with the request.
func (a *Auth) RequireApplicationToken() gin.HandlerFunc {
	return a.requireToken(func(tokenID string, user *model.User) (bool, uint) {
		if user != nil {
			return false, 0
		}
		if token := a.DB.GetApplicationByID(tokenID); token != nil {
			return true, token.UserID
		}
		return false, 0
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
	if header := ctx.Request.Header.Get(headerName); header != "" && strings.HasPrefix(header, headerSchema) {
		return strings.TrimPrefix(header, headerSchema)
	}
	return ""
}

func (a *Auth) userFromBasicAuth(ctx *gin.Context) *model.User {
	if name, pass, ok := ctx.Request.BasicAuth(); ok {
		if user := a.DB.GetUserByName(name); user != nil && ComparePassword(user.Pass, []byte(pass)) {
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
			if ok, userID := auth(token, user); ok {
				RegisterAuthentication(ctx, user, userID)
				ctx.Next()
				return
			}
		}
		ctx.AbortWithError(401, errors.New("could not authenticate"))
	}
}
