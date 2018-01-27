package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/model"
)

// RegisterAuthentication registers the user or the user id; The id can later be obtained by GetUserID.
func RegisterAuthentication(ctx *gin.Context, user *model.User, userID uint) {
	ctx.Set("user", user)
	ctx.Set("userid", userID)
}

// GetUserID returns the user id which was previously registered by RegisterAuthentication.
func GetUserID(ctx *gin.Context) uint {
	user := ctx.MustGet("user").(*model.User)
	if user == nil {
		userID := ctx.MustGet("userid").(uint)
		if userID == 0 {
			panic("token and user may not be null")
		}
		return userID
	}

	return user.ID
}
