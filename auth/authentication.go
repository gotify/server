package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/model"
	"strings"
)

const (
	headerName   = "Authorization"
	headerSchema = "ApiKey "
)

type Database interface {
	GetApplicationById(id string) *model.Application
	GetClientById(id string) *model.Client
	GetUserByName(name string) *model.User
	GetUserById(id uint) *model.User
}

type Auth struct {
	DB Database
}

type authenticate func(tokenId string, user *model.User) (success bool, userId uint)

func (a *Auth) RequireAdmin() gin.HandlerFunc {
	return a.requireToken(func(tokenId string, user *model.User) (bool, uint) {
		if user != nil {
			return user.Admin, user.Id
		}
		if token := a.DB.GetClientById(tokenId); token != nil {
			return a.DB.GetUserById(token.UserId).Admin, token.UserId
		}
		return false, 0
	})
}

func (a *Auth) RequireAll() gin.HandlerFunc {
	return a.requireToken(func(tokenId string, user *model.User) (bool, uint) {
		if user != nil {
			return true, user.Id
		}
		if token := a.DB.GetClientById(tokenId); token != nil {
			return true, token.UserId
		}
		return false, 0
	})
}

func (a *Auth) RequireWrite() gin.HandlerFunc {
	return a.requireToken(func(tokenId string, user *model.User) (bool, uint) {
		if user != nil {
			return false, 0
		}
		if token := a.DB.GetApplicationById(tokenId); token != nil {
			return true, token.UserId
		}
		return false, 0
	})
}

func (a *Auth) tokenFromQueryOrHeader(ctx *gin.Context) string {
	if token := a.tokenFromQuery(ctx); token != "" {
		return token
	} else if token := a.tokenFromHeader(ctx); token != "" {
		return token
	}
	return ""
}

func (a *Auth) tokenFromQuery(ctx *gin.Context) string {
	return ctx.Request.URL.Query().Get("token")
}

func (a *Auth) tokenFromHeader(ctx *gin.Context) string {
	if header := ctx.Request.Header.Get(headerName); header != "" && strings.HasPrefix(header, headerSchema) {
		return strings.TrimPrefix(header, headerSchema)
	}
	return ""
}

func (a *Auth) userFromBasicAuth(ctx *gin.Context) *model.User {
	if name, pass, ok := ctx.Request.BasicAuth(); ok {
		if user := a.DB.GetUserByName(name); user != nil && ComparePassword(user.Pass, []byte(pass)) {
			return user
		}
	}
	return nil
}

func (a *Auth) requireToken(auth authenticate) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := a.tokenFromQueryOrHeader(ctx)
		user := a.userFromBasicAuth(ctx)

		if user != nil || token != "" {
			if ok, _ := auth(token, user); ok {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithError(401, errors.New("could not authenticate"))
	}
}
