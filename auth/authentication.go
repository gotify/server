package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmattheis/memo/model"
	"strings"
)

const (
	headerName    = "Authorization"
	headerSchema  = "ApiKey "
	typeAdmin     = 0
	typeAll       = 1
	typeWriteOnly = 2
)

type Database interface {
	GetTokenById(id string) *model.Token
	GetUserByName(name string) *model.User
	GetUserById(id uint) *model.User
}

type Auth struct {
	DB Database
}

func (a *Auth) RequireAdmin() gin.HandlerFunc {
	return a.requireToken(typeAdmin)
}

func (a *Auth) RequireAll() gin.HandlerFunc {
	return a.requireToken(typeAll)
}

func (a *Auth) RequireWrite() gin.HandlerFunc {
	return a.requireToken(typeWriteOnly)
}

func (a *Auth) tokenFromQueryOrHeader(ctx *gin.Context) *model.Token {
	if token := a.tokenFromQuery(ctx); token != nil {
		return token
	} else if token := a.tokenFromHeader(ctx); token != nil {
		return token
	}
	return nil
}

func (a *Auth) tokenFromQuery(ctx *gin.Context) *model.Token {
	if token := ctx.Request.URL.Query().Get("token"); token != "" {
		return a.DB.GetTokenById(token)
	}
	return nil
}

func (a *Auth) tokenFromHeader(ctx *gin.Context) *model.Token {
	if header := ctx.Request.Header.Get(headerName); header != "" && strings.HasPrefix(header, headerSchema) {
		return a.DB.GetTokenById(strings.TrimPrefix(header, headerSchema))
	}
	return nil
}

func (a *Auth) userFromBasicAuth(ctx *gin.Context) *model.User {
	if name, pass, ok := ctx.Request.BasicAuth(); ok {
		if user := a.DB.GetUserByName(name); user != nil && ComparePassword(user.Pass, []byte(pass)) {
			return user
		}
	}
	return nil
}

func (a *Auth) isAuthenticated(checkType int, token *model.Token, user *model.User) bool {
	if token == nil && user == nil {
		return false
	}

	switch checkType {
	case typeWriteOnly:
		return true
	case typeAll:
		return user != nil || (token != nil && !token.WriteOnly)
	default:
		if user == nil {
			user = a.DB.GetUserById(token.UserID)
		}
		return user != nil && user.Admin
	}
}

func (a *Auth) requireToken(checkType int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := a.tokenFromQueryOrHeader(ctx)
		user := a.userFromBasicAuth(ctx)

		if a.isAuthenticated(checkType, token, user) {
			ctx.Next()
		} else {
			ctx.AbortWithError(401, errors.New("could not authenticate"))
		}
	}
}
