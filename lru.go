// Package lru implements a generic least recently used cache.
package lru

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	sync.RWMutex
	Capacity int
	Size     int
	entries  map[K]V
	head     *entry[K]
	tail     *entry[K]
}

type entry[K comparable] struct {
	Key  K
	Next *entry[K]
	Prev *entry[K]
}

// Fetch attempts to retrieve a value from the cache using the provided key. If
// the key isn't in the cache, the fetcher function will execute to get and store
// the value for the key.
func (c *Cache[K, V]) Fetch(key K, fetcher func() (V, error)) (V, bool, error) {
	c.RLock()

	/* cache hit */

	if val, ok := c.entries[key]; ok {
		c.RUnlock()
		return val, ok, nil
	}

	c.RUnlock()

	/* cache miss */

	val, err := fetcher()

	if err != nil {
		return val, false, err
	}

	c.Lock()

	defer c.Unlock()

	if c.head == nil {
		e := entry[K]{Key: key}

		c.head = &e
		c.tail = &e

		c.entries[key] = val

		c.Size++
	} else {
		if c.Size == c.Capacity {
			if c.tail.Next != nil {
				delete(c.entries, c.tail.Key)

				c.tail = c.tail.Next
				c.tail.Prev = nil

				c.Size--
			}
		}

		e := entry[K]{Key: key, Prev: c.head}

		c.head.Next = &e
		c.head = &e

		c.entries[key] = val

		c.Size++
	}

	return val, false, nil
}

// New initialzes a new LRU cache.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		Capacity: capacity,
		entries:  map[K]V{},
	}
}
