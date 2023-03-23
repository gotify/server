package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/model"
)

// RegisterAuthentication registers the user id, user and or token.
func RegisterAuthentication(ctx *gin.Context, user *model.User, userID uint, tokenID string, minPriority int) {
	ctx.Set("user", user)
	ctx.Set("userid", userID)
	ctx.Set("tokenid", tokenID)
	ctx.Set("min_priority", minPriority)
}

// GetUserID returns the user id which was previously registered by RegisterAuthentication.
func GetUserID(ctx *gin.Context) uint {
	id := TryGetUserID(ctx)
	if id == nil {
		panic("token and user may not be null")
	}
	return *id
}

// TryGetUserID returns the user id or nil if one is not set.
func TryGetUserID(ctx *gin.Context) *uint {
	user := ctx.MustGet("user").(*model.User)
	if user == nil {
		userID := ctx.MustGet("userid").(uint)
		if userID == 0 {
			return nil
		}
		return &userID
	}

	return &user.ID
}

// GetTokenID returns the tokenID.
func GetTokenID(ctx *gin.Context) string {
	return ctx.MustGet("tokenid").(string)
}

// GetMinPriority returns the minimum priority for messages for this client
func GetMinPriority(ctx *gin.Context) int {
	return ctx.MustGet("min_priority").(int)
}
