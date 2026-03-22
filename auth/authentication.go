package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/auth/password"
	"github.com/gotify/server/v2/model"
)

type authState int

const (
	authStateSkip authState = iota
	authStateForbidden
	authStateOk
)

const (
	headerName = "X-Gotify-Key"
	cookieName = "gotify-client-token"
)

// The Database interface for encapsulating database access.
type Database interface {
	GetApplicationByToken(token string) (*model.Application, error)
	GetClientByToken(token string) (*model.Client, error)
	GetUserByName(name string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateClientTokensLastUsed(tokens []string, t *time.Time) error
	UpdateApplicationTokenLastUsed(token string, t *time.Time) error
}

// Auth is the provider for authentication middleware.
type Auth struct {
	DB Database
}

// RequireAdmin returns a gin middleware which requires a client token or basic authentication header to be supplied
// with the request. Also the authenticated user must be an administrator.
func (a *Auth) RequireAdmin(ctx *gin.Context) {
	a.evaluateOr401(ctx, a.user(true), a.client(true))
}

// RequireClient returns a gin middleware which requires a client token or basic authentication header to be supplied
// with the request.
func (a *Auth) RequireClient(ctx *gin.Context) {
	a.evaluateOr401(ctx, a.user(false), a.client(false))
}

// RequireApplicationToken returns a gin middleware which requires an application token to be supplied with the request.
func (a *Auth) RequireApplicationToken(ctx *gin.Context) {
	if a.evaluate(ctx, a.application) {
		return
	}
	state, err := a.user(false)(ctx)
	if err != nil {
		ctx.AbortWithError(500, err)
	}
	if state != authStateSkip {
		// Return to the user that it's valid authentication, but we don't allow user auth for application endpoints.
		a.abort403(ctx)
		return
	}
	a.abort401(ctx)
}

func (a *Auth) Optional(ctx *gin.Context) {
	if !a.evaluate(ctx, a.user(false), a.client(false)) {
		RegisterAuthentication(ctx, nil, 0, "")
		ctx.Next()
	}
}

func (a *Auth) evaluate(ctx *gin.Context, funcs ...func(ctx *gin.Context) (authState, error)) bool {
	for _, fn := range funcs {
		state, err := fn(ctx)
		if err != nil {
			ctx.AbortWithError(500, err)
			return true
		}
		switch state {
		case authStateForbidden:
			a.abort403(ctx)
			return true
		case authStateOk:
			ctx.Next()
			return true
		case authStateSkip:
			continue
		}
	}
	return false
}

func (a *Auth) evaluateOr401(ctx *gin.Context, funcs ...func(ctx *gin.Context) (authState, error)) {
	if !a.evaluate(ctx, funcs...) {
		a.abort401(ctx)
	}
}

func (a *Auth) abort401(ctx *gin.Context) {
	ctx.AbortWithError(401, errors.New("you need to provide a valid access token or user credentials to access this api"))
}

func (a *Auth) abort403(ctx *gin.Context) {
	ctx.AbortWithError(403, errors.New("you are not allowed to access this api"))
}

func (a *Auth) user(requireAdmin bool) func(ctx *gin.Context) (authState, error) {
	return func(ctx *gin.Context) (authState, error) {
		if name, pass, ok := ctx.Request.BasicAuth(); ok {
			if user, err := a.DB.GetUserByName(name); err != nil {
				return authStateSkip, err
			} else if user != nil && password.ComparePassword(user.Pass, []byte(pass)) {
				RegisterAuthentication(ctx, user, user.ID, "")

				if requireAdmin && !user.Admin {
					return authStateForbidden, nil
				}
				return authStateOk, nil
			}
		}
		return authStateSkip, nil
	}
}

func (a *Auth) client(requireAdmin bool) func(ctx *gin.Context) (authState, error) {
	return func(ctx *gin.Context) (authState, error) {
		token := a.readTokenFromRequest(ctx)
		if token == "" {
			return authStateSkip, nil
		}
		client, err := a.DB.GetClientByToken(token)
		if err != nil {
			return authStateSkip, err
		}
		if client == nil {
			return authStateSkip, nil
		}
		RegisterAuthentication(ctx, nil, client.UserID, client.Token)

		now := time.Now()
		if client.LastUsed == nil || client.LastUsed.Add(5*time.Minute).Before(now) {
			if err := a.DB.UpdateClientTokensLastUsed([]string{client.Token}, &now); err != nil {
				return authStateSkip, err
			}
		}

		if requireAdmin {
			if user, err := a.DB.GetUserByID(client.UserID); err != nil {
				return authStateSkip, err
			} else if !user.Admin {
				return authStateForbidden, nil
			}
		}

		return authStateOk, nil
	}
}

func (a *Auth) application(ctx *gin.Context) (authState, error) {
	token := a.readTokenFromRequest(ctx)
	if token == "" {
		return authStateSkip, nil
	}
	app, err := a.DB.GetApplicationByToken(token)
	if err != nil {
		return authStateSkip, err
	}
	if app == nil {
		return authStateSkip, nil
	}
	RegisterAuthentication(ctx, nil, app.UserID, app.Token)

	now := time.Now()
	if app.LastUsed == nil || app.LastUsed.Add(5*time.Minute).Before(now) {
		if err := a.DB.UpdateApplicationTokenLastUsed(app.Token, &now); err != nil {
			return authStateSkip, err
		}
	}

	return authStateOk, nil
}

func (a *Auth) readTokenFromRequest(ctx *gin.Context) string {
	if token := a.tokenFromQuery(ctx); token != "" {
		return token
	} else if token := a.tokenFromXGotifyHeader(ctx); token != "" {
		return token
	} else if token := a.tokenFromAuthorizationHeader(ctx); token != "" {
		return token
	} else if token := a.tokenFromCookie(ctx); token != "" {
		return token
	}
	return ""
}

func (a *Auth) tokenFromCookie(ctx *gin.Context) string {
	token, err := ctx.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return token
}

func (a *Auth) tokenFromQuery(ctx *gin.Context) string {
	return ctx.Request.URL.Query().Get("token")
}

func (a *Auth) tokenFromXGotifyHeader(ctx *gin.Context) string {
	return ctx.Request.Header.Get(headerName)
}

func (a *Auth) tokenFromAuthorizationHeader(ctx *gin.Context) string {
	const prefix = "Bearer "

	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if len(authHeader) < len(prefix) || !strings.EqualFold(prefix, authHeader[:len(prefix)]) {
		return ""
	}

	return authHeader[len(prefix):]
}
