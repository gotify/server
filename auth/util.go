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

// TryGetTokenID returns the tokenID or an empty string if no token-based
// authentication was registered.
func TryGetTokenID(ctx *gin.Context) string {
	info := getInfo(ctx)
	switch {
	case info.client != nil:
		return info.client.Token
	case info.app != nil:
		return info.app.Token
	}
	return ""
}
