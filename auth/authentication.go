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
	authStateNotElevated
	authStateOk
)

const (
	headerName = "X-Gotify-Key"
	cookieName = "gotify-client-token"
)

var timeNow = time.Now

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
	DB           Database
	SecureCookie bool
}

// RequireAdmin requires an elevated client token or basic auth, the user must be an admin.
func (a *Auth) RequireAdmin(ctx *gin.Context) {
	a.evaluateOr401(ctx, a.handleUser(a.checkUserAdmin), a.handleClient(a.checkClientAdmin, a.checkClientElevated))
}

// RequireClient returns a gin middleware which requires a client token or basic authentication header to be supplied
// with the request.
func (a *Auth) RequireClient(ctx *gin.Context) {
	a.evaluateOr401(ctx, a.handleUser(), a.handleClient())
}

// RequireElevatedClient requires an elevated client token or basic auth.
func (a *Auth) RequireElevatedClient(ctx *gin.Context) {
	a.evaluateOr401(ctx, a.handleUser(), a.handleClient(a.checkClientElevated))
}

// RequireApplicationToken returns a gin middleware which requires an application token to be supplied with the request.
func (a *Auth) RequireApplicationToken(ctx *gin.Context) {
	if a.evaluate(ctx, a.handleApplication) {
		return
	}
	state, err := a.handleUser()(ctx)
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
	if !a.evaluate(ctx, a.handleUser(), a.handleClient()) {
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
		case authStateNotElevated:
			ctx.AbortWithError(403, errors.New("session not elevated, use basic auth or call /client:elevate"))
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

func (a *Auth) handleUser(checks ...func(*model.User) (authState, error)) func(ctx *gin.Context) (authState, error) {
	return func(ctx *gin.Context) (authState, error) {
		if name, pass, ok := ctx.Request.BasicAuth(); ok {
			if user, err := a.DB.GetUserByName(name); err != nil {
				return authStateSkip, err
			} else if user != nil && password.ComparePassword(user.Pass, []byte(pass)) {
				RegisterUser(ctx, user)

				for _, check := range checks {
					if state, err := check(user); err != nil || state != authStateOk {
						return state, err
					}
				}

				return authStateOk, nil
			}
		}
		return authStateSkip, nil
	}
}

func (a *Auth) handleClient(checks ...func(*model.Client) (authState, error)) func(ctx *gin.Context) (authState, error) {
	return func(ctx *gin.Context) (authState, error) {
		token, isCookie := a.readTokenFromRequest(ctx)
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
		RegisterClient(ctx, client)

		now := timeNow()
		if client.LastUsed == nil || client.LastUsed.Add(5*time.Minute).Before(now) {
			if err := a.DB.UpdateClientTokensLastUsed([]string{client.Token}, &now); err != nil {
				return authStateSkip, err
			}
			if isCookie {
				SetCookie(ctx.Writer, client.Token, CookieMaxAge, a.SecureCookie)
			}
		}

		for _, check := range checks {
			if state, err := check(client); err != nil || state != authStateOk {
				return state, err
			}
		}

		return authStateOk, nil
	}
}

func (a *Auth) handleApplication(ctx *gin.Context) (authState, error) {
	token, isCookie := a.readTokenFromRequest(ctx)
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
	RegisterApplication(ctx, app)

	now := timeNow()
	if app.LastUsed == nil || app.LastUsed.Add(5*time.Minute).Before(now) {
		if err := a.DB.UpdateApplicationTokenLastUsed(app.Token, &now); err != nil {
			return authStateSkip, err
		}
		if isCookie {
			SetCookie(ctx.Writer, app.Token, CookieMaxAge, a.SecureCookie)
		}
	}

	return authStateOk, nil
}

func (a *Auth) readTokenFromRequest(ctx *gin.Context) (string, bool) {
	if token := a.tokenFromQuery(ctx); token != "" {
		return token, false
	} else if token := a.tokenFromXGotifyHeader(ctx); token != "" {
		return token, false
	} else if token := a.tokenFromAuthorizationHeader(ctx); token != "" {
		return token, false
	} else if token := a.tokenFromCookie(ctx); token != "" {
		return token, true
	}
	return "", false
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

func (a *Auth) checkClientAdmin(client *model.Client) (authState, error) {
	if user, err := a.DB.GetUserByID(client.UserID); err != nil {
		return authStateSkip, err
	} else if !user.Admin {
		return authStateForbidden, nil
	}
	return authStateOk, nil
}

func (a *Auth) checkClientElevated(client *model.Client) (authState, error) {
	if client.ElevatedUntil == nil || !timeNow().Before(*client.ElevatedUntil) {
		return authStateNotElevated, nil
	}
	return authStateOk, nil
}

func (a *Auth) checkUserAdmin(user *model.User) (authState, error) {
	if !user.Admin {
		return authStateForbidden, nil
	}
	return authStateOk, nil
}
