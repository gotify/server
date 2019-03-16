package auth

import (
	"crypto/rand"
	"math/big"
)

var (
	tokenCharacters   = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_")
	randomTokenLength = 14
	applicationPrefix = "A"
	clientPrefix      = "C"
	pluginPrefix      = "P"

	randReader = rand.Reader
)

func randIntn(n int) int {
	max := big.NewInt(int64(n))
	res, err := rand.Int(randReader, max)
	if err != nil {
		panic("random source is not available")
	}
	return int(res.Int64())
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
	res := make([]byte, length)
	for i := range res {
		index := randIntn(len(tokenCharacters))
		res[i] = tokenCharacters[index]
	}
	return string(res)
}

func init() {
	randIntn(2)
}
