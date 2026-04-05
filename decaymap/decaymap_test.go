
package decaymap

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecayMap2(t *testing.T) {
	epoch := time.Now()
	dm := NewDecayMap[string, string](epoch, 10*time.Millisecond)

	now := epoch
	for ts := 0; ts < 100; ts++ {
		dm.Set(now, fmt.Sprintf("key%d", ts), fmt.Sprintf("value%d", ts))
		for backts := 0; backts < ts; backts++ {
			res, ok := dm.Get(fmt.Sprintf("key%d", backts))
			if ts-backts <= 10 {
				assert.True(t, ok)
				assert.Equal(t, fmt.Sprintf("value%d", backts), res)
			} else if ts-backts >= 20 {
				assert.False(t, ok)
				assert.Equal(t, "", res)
			}
		}
		now = now.Add(1 * time.Millisecond)
	}

	now = now.Add(20 * time.Millisecond)
	dm.Set(now, "dummy", "dummy") // rachet internal state
	for ts := 0; ts < 100; ts++ {
		res, ok := dm.Get(fmt.Sprintf("key%d", ts))
		assert.False(t, ok)
		assert.Equal(t, "", res)
	}
}
