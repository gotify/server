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

func generateRandomToken(prefix string) string {
	b := make([]rune, randomTokenLength)
	for i := range b {
		b[i] = tokenCharacters[rand.Intn(len(tokenCharacters))]
	}
	return prefix + string(b)
}
