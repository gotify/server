package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/auth"
	"github.com/gotify/server/v2/auth/password"
	"github.com/gotify/server/v2/model"
)

// SessionDatabase is the interface for session-related database access.
type SessionDatabase interface {
	GetUserByName(name string) (*model.User, error)
	CreateClient(client *model.Client) error
	GetClientByToken(token string) (*model.Client, error)
	DeleteClientByID(id uint) error
}

// SessionAPI provides handlers for cookie-based session authentication.
type SessionAPI struct {
	DB            SessionDatabase
	NotifyDeleted func(uint, string)
}

// Login authenticates via basic auth, creates a client, sets an HttpOnly cookie, and returns user info.
func (a *SessionAPI) Login(ctx *gin.Context) {
	name, pass, ok := ctx.Request.BasicAuth()
	if !ok {
		ctx.AbortWithError(401, errors.New("basic auth required"))
		return
	}

	user, err := a.DB.GetUserByName(name)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	if user == nil || !password.ComparePassword(user.Pass, []byte(pass)) {
		ctx.AbortWithError(401, errors.New("invalid credentials"))
		return
	}

	clientParams := ClientParams{}
	if err := ctx.Bind(&clientParams); err != nil {
		return
	}

	client := model.Client{
		Name:   clientParams.Name,
		Token:  auth.GenerateNotExistingToken(generateClientToken, a.clientExists),
		UserID: user.ID,
	}
	if success := successOrAbort(ctx, 500, a.DB.CreateClient(&client)); !success {
		return
	}

	auth.SetCookie(ctx.Writer, client.Token, auth.CookieMaxAge)

	ctx.JSON(200, &model.UserExternal{
		ID:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
	})
}

// Logout deletes the client for the current session and clears the cookie.
func (a *SessionAPI) Logout(ctx *gin.Context) {
	auth.SetCookie(ctx.Writer, "", -1)

	tokenID := auth.TryGetTokenID(ctx)
	if tokenID == "" {
		ctx.AbortWithError(400, errors.New("no client auth provided"))
		return
	}
	client, err := a.DB.GetClientByToken(tokenID)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	if client == nil {
		ctx.Status(200)
		return
	}

	a.NotifyDeleted(client.UserID, client.Token)
	if success := successOrAbort(ctx, 500, a.DB.DeleteClientByID(client.ID)); !success {
		return
	}

	ctx.Status(200)
}

func (a *SessionAPI) clientExists(token string) bool {
	client, _ := a.DB.GetClientByToken(token)
	return client != nil
}
