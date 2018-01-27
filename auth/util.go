package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/model"
)

func RegisterAuthentication(ctx *gin.Context, user *model.User, userId uint) {
	ctx.Set("user", user)
	ctx.Set("userid", userId)
}

func GetUserId(ctx *gin.Context) uint {
	user := ctx.MustGet("user").(*model.User)
	if user == nil {
		userId := ctx.MustGet("userid").(uint)
		if userId == 0 {
			panic("token and user may not be null")
		}
		return userId
	}

	return user.Id
}
