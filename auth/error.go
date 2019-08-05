package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthenticationError interface {
	error
	Code() int
}

func abortContextWithAuthenticaionError(ctx *gin.Context, err error) {
	if authError, ok := err.(AuthenticationError); ok {
		ctx.AbortWithError(authError.Code(), authError)
	} else {
		ctx.AbortWithError(401, err)
	}
}

type NoAuthProviderError struct {
	DesignatedAuthenticator string
}

func (e NoAuthProviderError) Error() string {
	if e.DesignatedAuthenticator == "" {
		e.DesignatedAuthenticator = "all available authenticators"
	}
	return fmt.Sprintf("%s failed to authenticate this request", e.DesignatedAuthenticator)
}

func (e NoAuthProviderError) Code() int {
	return 401
}

type AuthProviderNotFoundError struct{}

func (a AuthProviderNotFoundError) Error() string {
	return "the designated authenticator is not loaded"
}

func (e AuthProviderNotFoundError) Code() int {
	return 400
}

type TokenRequiredError struct {
	TokenType string
}

// Error implements AuthenticationError
func (e TokenRequiredError) Error() string {
	return fmt.Sprintf("%s token is required to access this api", e.TokenType)
}

// Code implements AuthenticationError
func (e TokenRequiredError) Code() int {
	return 400
}
