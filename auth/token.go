package auth

import (
	"crypto/rand"
	"io"
	"math/big"
)

var (
	tokenCharacters   = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_")
	randomTokenLength = 14
	applicationPrefix = "A"
	clientPrefix      = "C"
	pluginPrefix      = "P"

	randSource RandSource = &RandSourceFromReader{rand.Reader}
)

// RandSourceFromReader is randomization source from a Reader to random data
type RandSourceFromReader struct {
	Source io.Reader
}

// Token implements RandSource
func (r *RandSourceFromReader) Token(length int, chars []byte) string {
	b := make([]byte, length)
	for i := range b {
		index := randIntn(r.Source, len(chars))
		b[i] = chars[index]
	}
	return string(b)
}

// RandSource is an abstraction of randomization provider
type RandSource interface {
	Token(len int, chars []byte) string
}

func randIntn(randReader io.Reader, n int) int {
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
	return randSource.Token(length, tokenCharacters)
}

func init() {
	randSource.Token(1, tokenCharacters)
}
