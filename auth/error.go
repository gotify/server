package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// AuthenticationError is an interface for authentication-related errors with custom HTTP status code
type AuthenticationError interface {
	error
	Code() int
}

func abortContextWithAuthenticationError(ctx *gin.Context, err error) {
	if authError, ok := err.(AuthenticationError); ok {
		ctx.AbortWithError(authError.Code(), authError)
	} else {
		ctx.AbortWithError(401, err)
	}
}

// NoAuthProviderError is returned then no authentication provider can handle this request
type NoAuthProviderError struct {
	DesignatedAuthenticator string
}

// Error implements AuthenticationError
func (e NoAuthProviderError) Error() string {
	if e.DesignatedAuthenticator == "" {
		e.DesignatedAuthenticator = "all available authenticators"
	}
	return fmt.Sprintf("%s failed to authenticate this request", e.DesignatedAuthenticator)
}

// Code implements AuthenticationError
func (e NoAuthProviderError) Code() int {
	return 401
}

// ProviderNotFoundError is returned when the designated authentication provider cannot be found
type ProviderNotFoundError struct{}

// Error implements AuthenticationError
func (e ProviderNotFoundError) Error() string {
	return "the designated authenticator is not loaded"
}

// Code implements AuthenticationError
func (e ProviderNotFoundError) Code() int {
	return 400
}

// TokenRequiredError is returned when a token is required for this operation
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

// NotAdminError is returned when admin priviledge is required for this operation
type NotAdminError struct{}

// Error implements AuthenticationError
func (e NotAdminError) Error() string {
	return "you do not have sufficient priviledge to access this api"
}

// Code implements AuthenticationError
func (e NotAdminError) Code() int {
	return 403
}
