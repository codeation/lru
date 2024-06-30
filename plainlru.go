package lru

import (
	"container/list"
)

type pair[K comparable, V any] struct {
	key   K
	value V
}

// LRU is a plain LRU cache.
type LRU[K comparable, V any] struct {
	queue    *list.List
	elems    map[K]*list.Element
	capacity int
	fn       func(K) (V, error)
}

// New creates an LRU cache with the specified capacity.
// f is a callback function for getting a value by key, which is called if there is no value in the cache.
func New[K comparable, V any](capacity int, fn func(K) (V, error)) *LRU[K, V] {
	return &LRU[K, V]{
		queue:    list.New(),
		elems:    make(map[K]*list.Element, capacity),
		capacity: capacity,
		fn:       fn,
	}
}

// Reset resets cache contents.
func (c *LRU[K, V]) Reset() {
	c.queue.Init()
	c.elems = make(map[K]*list.Element, c.capacity)
}

// Get returns the cached value for the key, or waits for the callback function to return a value.
func (c *LRU[K, V]) Get(key K) (V, error) {
	e, ok := c.elems[key]
	if ok {
		c.queue.MoveToFront(e)
		return e.Value.(*pair[K, V]).value, nil
	}
	value, err := c.fn(key)
	if err != nil {
		return value, err
	}
	for c.queue.Len() >= c.capacity {
		e := c.queue.Back()
		c.queue.Remove(e)
		delete(c.elems, e.Value.(*pair[K, V]).key)
	}
	c.elems[key] = c.queue.PushFront(&pair[K, V]{
		key:   key,
		value: value,
	})
	return value, nil
}
