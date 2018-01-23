package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/model"
	"github.com/jmattheis/memo/auth"
	"fmt"
	"math/rand"
)

var (
	tokenCharacters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_")
	tokenLength     = 15
)

type TokenDatabase interface {
	CreateToken(token *model.Token) error
	GetTokenById(id string) *model.Token
	GetTokensByUser(userId uint) []*model.Token
	DeleteToken(id string) error
}

type TokenApi struct {
	DB TokenDatabase
}

func (a *TokenApi) CreateToken(ctx *gin.Context) {
	token := model.Token{}
	if err := ctx.Bind(&token); err == nil {
		for ok := true; ok; ok = a.DB.GetTokenById(token.Id) != nil {
			token.Id = randToken()
		}
		token.UserID = auth.GetUserID(ctx)
		a.DB.CreateToken(&token)
		ctx.JSON(200, token)
	}
}

func (a *TokenApi) GetTokens(ctx *gin.Context) {
	userId := auth.GetUserID(ctx)
	tokens := a.DB.GetTokensByUser(userId)
	ctx.JSON(200, tokens)
}

func (a *TokenApi) DeleteToken(ctx *gin.Context) {
	tokenId := ctx.Param("id")
	if token := a.DB.GetTokenById(tokenId); token != nil {
		a.DB.DeleteToken(tokenId)
	} else {
		ctx.AbortWithError(404, fmt.Errorf("token with id %s doesn't exists", tokenId))
	}
}

func randToken() string {
	b := make([]rune, tokenLength)
	for i := range b {
		b[i] = tokenCharacters[rand.Intn(len(tokenCharacters))]
	}
	return string(b)
}
