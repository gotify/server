package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/model"
)

const authKey = "auth"

type authentication struct {
	client *model.Client
	app    *model.Application
	user   *model.User
}

// RegisterUser stores the authenticated user on the gin context.
func RegisterUser(ctx *gin.Context, user *model.User) {
	ctx.Set(authKey, &authentication{user: user})
}

// RegisterClient stores the authenticated client on the gin context.
func RegisterClient(ctx *gin.Context, client *model.Client) {
	ctx.Set(authKey, &authentication{client: client})
}

// RegisterApplication stores the authenticated application on the gin context.
func RegisterApplication(ctx *gin.Context, app *model.Application) {
	ctx.Set(authKey, &authentication{app: app})
}

func getInfo(ctx *gin.Context) *authentication {
	if v, ok := ctx.Get(authKey); ok {
		return v.(*authentication)
	}
	return &authentication{}
}

// GetUserID returns the user id which was previously registered by one of the Register* functions.
func GetUserID(ctx *gin.Context) uint {
	id := TryGetUserID(ctx)
	if id == nil {
		panic("token and user may not be null")
	}
	return *id
}

// TryGetUserID returns the user id or nil if one is not set.
func TryGetUserID(ctx *gin.Context) *uint {
	info := getInfo(ctx)
	switch {
	case info.user != nil:
		return &info.user.ID
	case info.client != nil:
		return &info.client.UserID
	case info.app != nil:
		return &info.app.UserID
	default:
		return nil
	}
}

// GetApplication returns the authenticated application or nil if no application
// was registered.
func GetApplication(ctx *gin.Context) *model.Application {
	return getInfo(ctx).app
}

// GetClient returns the authenticated client or nil if no client was registered.
func GetClient(ctx *gin.Context) *model.Client {
	return getInfo(ctx).client
}
