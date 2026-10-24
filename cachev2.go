package cachev2

import (
	"sync"
	"time"
)

type item struct {
	value     any
	expiresAt time.Time
}

type Memory struct {
	mu    sync.RWMutex
	items map[string]item
	stop  chan struct{}
}

func New() *Memory {
	m := &Memory{
		items: make(map[string]item),
		stop:  make(chan struct{}),
	}

	return m
}

func (m *Memory) Set(key string, value any, ttl time.Duration) {
	var exp time.Time

	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	m.mu.Lock()
	m.items[key] = item{value: value, expiresAt: exp}
	m.mu.Unlock()
}

func (m *Memory) Get(key string) (any, bool) {
	m.mu.RLock()
	it, ok := m.items[key]
	m.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if !it.expiresAt.IsZero() && time.Now().After(it.expiresAt) {
		m.mu.Lock()
		if cur, ok := m.items[key]; ok && cur.expiresAt.Equal(it.expiresAt) {
			delete(m.items, key)
		}
		m.mu.Unlock()
		return nil, false
	}

	return it.value, true
}

func (m *Memory) Delete(key string) {
	m.mu.Lock()
	delete(m.items, key)
	m.mu.Unlock()
}
