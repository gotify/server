package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
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

// Convert a Token to its hashed representation.
func HashToken(s string, salt []byte) (string, error) {
	saltHex := fmt.Sprintf("%x", salt)
	sha := sha256.New()
	_, err := sha.Write(salt)
	if err != nil {
		return "", err
	}
	_, err = sha.Write([]byte(s))
	if err != nil {
		return "", err
	}
	hashed := fmt.Sprintf("%x", sha.Sum(nil))
	return fmt.Sprintf("%s$%s", saltHex, hashed), nil
}

// CompareToken compares a token with a hashed representation, optionally upgrading the hash if necessary.
func CompareToken(s string, hashed string) (bool, *string, error) {
	if len(s) != randomTokenLength+1 /* prefix */ {
		return false, nil, errors.New("invalid token length")
	}

	split := strings.SplitN(hashed, "$", 2)

	// determine if we need to upgrade the hash
	if len(split) == 1 {
		match := s == hashed
		if match {
			var salt [16]byte
			_, err := io.ReadFull(randReader, salt[:])
			if err != nil {
				return false, nil, err
			}
			hashed, err := HashToken(s, salt[:])
			if err != nil {
				return false, nil, err
			}
			return true, &hashed, nil
		} else {
			return false, nil, nil
		}
	}

	if len(split) == 2 {
		salt, err := hex.DecodeString(split[0])
		if err != nil {
			return false, nil, err
		}
		inputHashed, err := HashToken(s, salt)
		if err != nil {
			return false, nil, err
		}
		return inputHashed == hashed, nil, nil
	}

	return false, nil, errors.New("invalid hash format")
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
