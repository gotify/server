package auth

import (
	"log"
	"math/rand"
)

// UseMathRand changes the random generator to math/rand, used only in testing environment
func UseMathRand() {
	log.Println("warn: using math/rand for randomness")
	randReader = rand.New(rand.NewSource(1))
}
