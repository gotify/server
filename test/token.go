package test

import "sync"

// Tokens returns a token generation function with takes a series of tokens and output them in order.
func Tokens(tokens ...string) func() string {
	var i int
	lock := sync.Mutex{}
	return func() string {
		lock.Lock()
		defer lock.Unlock()
		res := tokens[i%len(tokens)]
		i++
		return res
	}
}
