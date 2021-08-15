// Package lru implements asynchronous LRU cache.
package lru

import (
	"container/list"
	"fmt"
	"reflect"
	"sync"
)

type pair struct {
	once  *sync.Once
	value interface{}
	err   error
	elem  *list.Element
}

// Cache is a LRU cache.
type Cache struct {
	values  map[string]*pair
	mutex   sync.Mutex
	maxSize int
	list    *list.List
}

// NewCache creates an LRU cache with the specified size.
func NewCache(maxSize int) *Cache {
	return &Cache{
		values:  map[string]*pair{},
		maxSize: maxSize,
		list:    list.New(),
	}
}

// Reset resets cache contents.
func (c *Cache) Reset() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.values = map[string]*pair{}
	c.list.Init()
	return nil
}

// Get returns the cached value for the key, or waits until f returns a value.
func (c *Cache) Get(key string, f func(key string) (interface{}, error), value interface{}) error {
	c.mutex.Lock()
	p, ok := c.values[key]
	if ok {
		c.list.MoveToFront(p.elem)
	} else {
		for c.list.Len() >= c.maxSize {
			e := c.list.Back()
			delete(c.values, e.Value.(string))
			c.list.Remove(e)
		}
		p = &pair{
			once: new(sync.Once),
			elem: c.list.PushFront(key),
		}
		c.values[key] = p
	}
	c.mutex.Unlock()
	p.once.Do(func() {
		p.value, p.err = f(key)
	})
	if p.err != nil {
		return p.err
	}
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("value must be a non-nil pointer")
	}
	if !reflect.TypeOf(p.value).AssignableTo(v.Elem().Type()) {
		return fmt.Errorf("value must be a %T pointer", p.value)
	}
	v.Elem().Set(reflect.ValueOf(p.value))
	return nil
}
