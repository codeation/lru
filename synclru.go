// Package lru implements asynchronous LRU cache.
package lru

import (
	"container/list"
	"sync"
)

type oncePair[K comparable, V any] struct {
	key   K
	value V
	err   error
	once  sync.Once
}

// SyncLRU is a LRU cache for concurrent use.
type SyncLRU[K comparable, V any] struct {
	queue    *list.List
	elems    map[K]*list.Element
	mutex    sync.Mutex
	capacity int
	fn       func(K) (V, error)
}

// NewSyncLRU creates an LRU cache with the specified capacity.
// f is a callback function for getting a value by key, which is called if there is no value in the cache.
// The callback function will be called once for concurrent Get requests with the same key.
func NewSyncLRU[K comparable, V any](capacity int, fn func(K) (V, error)) *SyncLRU[K, V] {
	return &SyncLRU[K, V]{
		queue:    list.New(),
		elems:    make(map[K]*list.Element, capacity),
		capacity: capacity,
		fn:       fn,
	}
}

// Reset resets cache contents.
func (c *SyncLRU[K, V]) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.queue.Init()
	c.elems = make(map[K]*list.Element, c.capacity)
}

// Get returns the cached value for the key, or waits for the callback function to return a value.
func (c *SyncLRU[K, V]) Get(key K) (V, error) {
	c.mutex.Lock()
	var p *oncePair[K, V]
	e, ok := c.elems[key]
	if ok {
		p = e.Value.(*oncePair[K, V])
		c.queue.MoveToFront(e)
	} else {
		for c.queue.Len() >= c.capacity {
			backElem := c.queue.Back()
			c.queue.Remove(backElem)
			delete(c.elems, backElem.Value.(*oncePair[K, V]).key)
		}
		p = &oncePair[K, V]{
			key: key,
		}
		e = c.queue.PushFront(p)
		c.elems[key] = e
	}
	c.mutex.Unlock()
	p.once.Do(func() {
		p.value, p.err = c.fn(key)
		if p.err != nil {
			c.mutex.Lock()
			c.queue.Remove(e)
			delete(c.elems, key)
			c.mutex.Unlock()
		}
	})
	return p.value, p.err
}
