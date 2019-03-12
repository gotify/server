package auth

import (
	"crypto/rand"
	"math"
)

var (
	tokenCharacters   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_")
	randomTokenLength = 14
	applicationPrefix = "A"
	clientPrefix      = "C"
	pluginPrefix      = "P"

	randReader = rand.Reader
)

func randIntn(n int) int {
	requiredBytes := (int(math.Ceil(math.Log2(float64(n)))) + 7) / 8
	for {
		buf := make([]byte, requiredBytes)
		if _, err := randReader.Read(buf); err != nil {
			panic("crypto rand is unavailable")
		}
		res := 0
		for _, n := range buf {
			res <<= 8
			res += int(n)
		}
		if res < n {
			return res
		}
	}
}

// GenerateNotExistingToken receives a token generation func and a func to check whether the token exists, returns a unique token.
func GenerateNotExistingToken(generateToken func() string, tokenExists func(token string) bool) string {
	for {
		token := generateToken()
		if !tokenExists(token) {
			return token
		}
	}
}

// GenerateApplicationToken generates an application token.
func GenerateApplicationToken() string {
	return generateRandomToken(applicationPrefix)
}

// GenerateClientToken generates a client token.
func GenerateClientToken() string {
	return generateRandomToken(clientPrefix)
}

// GeneratePluginToken generates a plugin token.
func GeneratePluginToken() string {
	return generateRandomToken(pluginPrefix)
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
		index := randIntn(len(tokenCharacters))
		b[i] = tokenCharacters[index]
	}
	return string(b)
}

func init() {
	randIntn(1)
}
