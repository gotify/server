package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/model"
)

func RegisterAuthentication(ctx *gin.Context, user *model.User, token *model.Token) {
	ctx.Set("user", user)
	ctx.Set("token", token)
}

func GetUserID(ctx *gin.Context) uint {
	user := ctx.MustGet("user").(*model.User);
	if user == nil {
		token := GetToken(ctx)
		if token == nil {
			panic("token and user may not be null")
		}
		return token.UserID
	}

	return user.ID
}

func GetToken(ctx *gin.Context) *model.Token {
	return ctx.MustGet("token").(*model.Token);
}
