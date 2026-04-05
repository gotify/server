package decaymap

import (
	"sync"
	"time"
)

// DecayMap is a coarse-grained paired map that bulk-frees outdated entries.
type DecayMap[K comparable, V any] struct {
	mu         sync.Mutex
	maps       [2]map[K]V
	epoch      time.Time
	lastInsert time.Time
	generation uint64
	period     time.Duration
}

func NewDecayMap[K comparable, V any](epoch time.Time, period time.Duration) *DecayMap[K, V] {
	return &DecayMap[K, V]{
		maps:       [2]map[K]V{make(map[K]V), make(map[K]V)},
		epoch:      epoch,
		generation: 0,
		period:     period,
	}
}

// Attempts to retrieve and remove a value from the store, does not guarantee the value is unexpired.
func (m *DecayMap[K, V]) Pop(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	res, ok := m.maps[0][key]
	if ok {
		delete(m.maps[0], key)
		return res, ok
	}
	res, ok = m.maps[1][key]
	if ok {
		delete(m.maps[1], key)
		return res, ok
	}
	return res, ok
}

// Sets a value in the store, overwriting the existing value if exists.
//
// Inserting values in non monotonic order is a no op.
func (m *DecayMap[K, V]) Set(now time.Time, key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if now.Before(m.lastInsert) {
		return
	}
	m.lastInsert = now
	curGen := uint64(now.Sub(m.epoch) / m.period)
	genDelta := curGen - m.generation
	if genDelta >= 2 {
		// both are outdated, clear both
		m.maps[0] = make(map[K]V)
		m.maps[1] = make(map[K]V)
		m.generation = curGen
	} else if genDelta == 1 {
		// one is outdated, clear the outdated one
		m.maps[curGen&1] = make(map[K]V)
		m.generation = curGen
	}
	delete(m.maps[(curGen&1)^1], key)
	m.maps[curGen&1][key] = value
}
