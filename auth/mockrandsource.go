package auth

import (
	"crypto/rand"
	"log"
)

// MockRandSource is random source that generated predefined tokens, used only in testing purposes
type MockRandSource struct {
	Tokens  []string
	counter int
}

// Token implements RandSource
func (r *MockRandSource) Token(length int, chars []byte) string {
	res := r.Tokens[r.counter%len(r.Tokens)]
	r.counter++
	return res
}

// UseCryptoRand resets the random source to crypto/rand
func UseCryptoRand() {
	randSource = &RandSourceFromReader{rand.Reader}
}

// UseTokenSource overrides the random source used for generating tokens, used only in testing environment
func UseTokenSource(source RandSource) {
	log.Println("warn: using custom random source")
	randSource = source
}
