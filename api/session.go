package api

import (
	"errors"
	"time"

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
	SecureCookie  bool
}

// swagger:operation POST /auth/local/login auth localLogin
//
// Authenticate via basic auth and create a session.
//
//	---
//	consumes: [application/x-www-form-urlencoded]
//	produces: [application/json]
//	security:
//	- basicAuth: []
//	parameters:
//	- name: name
//	  in: formData
//	  description: the client name to create
//	  required: true
//	  type: string
//	responses:
//	  200:
//	    description: Ok
//	    schema:
//	        $ref: "#/definitions/CurrentUser"
//	    headers:
//	      Set-Cookie:
//	        type: string
//	        description: session cookie
//	  401:
//	    description: Unauthorized
//	    schema:
//	        $ref: "#/definitions/Error"
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

	elevatedUntil := time.Now().Add(model.DefaultElevationDuration)
	client := model.Client{
		Name:          clientParams.Name,
		Token:         auth.GenerateNotExistingToken(generateClientToken, a.clientExists),
		UserID:        user.ID,
		ElevatedUntil: &elevatedUntil,
	}
	if success := successOrAbort(ctx, 500, a.DB.CreateClient(&client)); !success {
		return
	}

	auth.SetCookie(ctx.Writer, client.Token, auth.CookieMaxAge, a.SecureCookie)

	ctx.JSON(200, &model.CurrentUserExternal{
		ID:            user.ID,
		Name:          user.Name,
		Admin:         user.Admin,
		ClientID:      client.ID,
		ElevatedUntil: client.ElevatedUntil,
	})
}

// swagger:operation POST /auth/logout auth logout
//
// End the current session.
//
// Clears the session cookie and deletes the associated client.
//
//	---
//	produces: [application/json]
//	security:
//	- clientTokenHeader: []
//	- clientTokenQuery: []
//	- basicAuth: []
//	responses:
//	  200:
//	    description: Ok
//	    headers:
//	      Set-Cookie:
//	        type: string
//	        description: cleared session cookie
//	  400:
//	    description: Bad Request
//	    schema:
//	        $ref: "#/definitions/Error"
func (a *SessionAPI) Logout(ctx *gin.Context) {
	auth.SetCookie(ctx.Writer, "", -1, a.SecureCookie)

	client := auth.GetClient(ctx)
	if client == nil {
		ctx.AbortWithError(403, errors.New("no client auth provided"))
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
