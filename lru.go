// Package lru implements asynchronous LRU cache.
package lru

import (
	"container/list"
	"sync"
)

type pair[V any] struct {
	once  *sync.Once
	value V
	err   error
	elem  *list.Element
}

// Cache is a LRU cache.
type Cache[K comparable, V any] struct {
	values   map[K]*pair[V]
	mutex    sync.Mutex
	capacity int
	keys     *list.List
	f        func(K) (V, error)
}

// NewCache creates an LRU cache with the specified capacity;
// f - function to get value by key, which is called if there is no value in the cache
func NewCache[K comparable, V any](capacity int, f func(K) (V, error)) *Cache[K, V] {
	return &Cache[K, V]{
		values:   make(map[K]*pair[V], capacity),
		capacity: capacity,
		keys:     list.New(),
		f:        f,
	}
}

// Reset resets cache contents.
func (c *Cache[K, V]) Reset() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.values = make(map[K]*pair[V], c.capacity)
	c.keys.Init()
	return nil
}

// Get returns the cached value for the key, or waits until f returns a value.
func (c *Cache[K, V]) Get(key K) (V, error) {
	c.mutex.Lock()
	p, ok := c.values[key]
	if ok {
		c.keys.MoveToFront(p.elem)
	} else {
		for c.keys.Len() >= c.capacity {
			e := c.keys.Back()
			delete(c.values, e.Value.(K))
			c.keys.Remove(e)
		}
		p = &pair[V]{
			once: new(sync.Once),
			elem: c.keys.PushFront(key),
		}
		c.values[key] = p
	}
	c.mutex.Unlock()
	p.once.Do(func() {
		p.value, p.err = c.f(key)
		if p.err != nil {
			c.mutex.Lock()
			p.once = new(sync.Once)
			c.mutex.Unlock()
		}
	})
	return p.value, p.err
}
