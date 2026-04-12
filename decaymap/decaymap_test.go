package decaymap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecayMap(t *testing.T) {
	epoch := time.Now()
	dm := NewDecayMap[string, string](epoch, 10*time.Second)

	dm.Set(epoch.Add(1*time.Second), "11", "value11")
	dm.Set(epoch.Add(2*time.Second), "12", "value12")
	dm.Set(epoch.Add(3*time.Second), "13", "value13")

	value, ok := dm.Pop("11")
	assert.True(t, ok)
	assert.Equal(t, "value11", value)

	_, ok = dm.Pop("11")
	assert.False(t, ok)

	dm.Set(epoch.Add(11*time.Second), "21", "value21")
	dm.Set(epoch.Add(12*time.Second), "22", "value22")
	dm.Set(epoch.Add(13*time.Second), "23", "value23")

	value, ok = dm.Pop("21")
	assert.True(t, ok)
	assert.Equal(t, "value21", value)

	value, ok = dm.Pop("12")
	assert.True(t, ok)
	assert.Equal(t, "value12", value)

	dm.Set(epoch.Add(21*time.Second), "31", "value31")

	_, ok = dm.Pop("13")
	assert.False(t, ok)
}
