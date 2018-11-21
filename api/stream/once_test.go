package stream

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Execute(t *testing.T) {
	executeOnce := once{}
	execution := make(chan struct{})
	fExecute := func() {
		execution <- struct{}{}
	}
	go executeOnce.Do(fExecute)
	go executeOnce.Do(fExecute)

	select {
	case <-execution:
		// expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("fExecute should be executed once")
	}

	select {
	case <-execution:
		t.Fatal("should only execute once")
	case <-time.After(100 * time.Millisecond):
		// expected
	}

	assert.False(t, executeOnce.mayExecute())

	go executeOnce.Do(fExecute)

	select {
	case <-execution:
		t.Fatal("should only execute once")
	case <-time.After(100 * time.Millisecond):
		// expected
	}
}
