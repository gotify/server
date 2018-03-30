package auth

import (
	"math/rand"
)

var (
	tokenCharacters   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_")
	randomTokenLength = 14
	applicationPrefix = "A"
	clientPrefix      = "C"
)

// GenerateApplicationToken generates an application token.
func GenerateApplicationToken() string {
	return generateRandomToken(applicationPrefix)
}

// GenerateClientToken generates a client token.
func GenerateClientToken() string {
	return generateRandomToken(clientPrefix)
}

// GenerateImageName generates an image name.
func GenerateImageName() string {
	return generateRandomString(25)
}

func generateRandomToken(prefix string) string {
	return prefix + generateRandomString(randomTokenLength)
}

func generateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = tokenCharacters[rand.Intn(len(tokenCharacters))]
	}
	return string(b)
}
