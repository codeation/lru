// Package lru implements asynchronous LRU cache.
package lru

import (
	"container/list"
	"sync"
)

type pair[V any] struct {
	once    sync.Once
	keyElem *list.Element
	value   V
	err     error
}

// Cache is a LRU cache.
type Cache[K comparable, V any] struct {
	keys     *list.List
	values   map[K]*pair[V]
	mutex    sync.Mutex
	capacity int
	f        func(K) (V, error)
}

// NewCache creates an LRU cache with the specified capacity;
// f - function to get value by key, which is called if there is no value in the cache
func NewCache[K comparable, V any](capacity int, f func(K) (V, error)) *Cache[K, V] {
	return &Cache[K, V]{
		keys:     list.New(),
		values:   make(map[K]*pair[V], capacity),
		capacity: capacity,
		f:        f,
	}
}

// Reset resets cache contents.
func (c *Cache[K, V]) Reset() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.keys.Init()
	c.values = make(map[K]*pair[V], c.capacity)
	return nil
}

// Get returns the cached value for the key, or waits until f returns a value.
func (c *Cache[K, V]) Get(key K) (V, error) {
	c.mutex.Lock()
	p, ok := c.values[key]
	if ok {
		c.keys.MoveToFront(p.keyElem)
	} else {
		for c.keys.Len() >= c.capacity {
			keyElem := c.keys.Back()
			c.keys.Remove(keyElem)
			delete(c.values, keyElem.Value.(K))
		}
		p = &pair[V]{
			keyElem: c.keys.PushFront(key),
		}
		c.values[key] = p
	}
	c.mutex.Unlock()
	p.once.Do(func() {
		p.value, p.err = c.f(key)
		if p.err != nil {
			c.mutex.Lock()
			c.keys.Remove(p.keyElem)
			delete(c.values, key)
			c.mutex.Unlock()
		}
	})
	return p.value, p.err
}
