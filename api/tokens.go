package api

import (
	"github.com/gotify/server/auth"
)

var generateApplicationToken = func() string {
	return auth.GenerateApplicationToken()
}

var generateClientToken = func() string {
	return auth.GenerateClientToken()
}

var generateImageName = func() string {
	return auth.GenerateImageName()
}
