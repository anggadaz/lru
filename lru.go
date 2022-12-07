// Package lru implements a generic least recently used cache.
package lru

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	sync.Mutex
	Capacity int
	Entries  []Entry[K, V]
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Fetch attempts to retrieve a value from the cache using the provided key. If
// the key isn't in the cache, the fetcher function will execute to get and store
// the value for the key.
func (l *Cache[K, V]) Fetch(key K, fetcher func() (V, error)) (V, bool, error) {
	l.Lock()

	index := -1

	for i, item := range l.Entries {
		if item.Key == key {
			index = i
			break
		}
	}

	/* cache miss */

	if index == -1 {
		l.Unlock()

		v, err := fetcher()

		if err != nil {
			return v, false, err
		}

		l.Lock()

		defer l.Unlock()

		l.insert(Entry[K, V]{Key: key, Value: v})

		return v, false, nil
	}

	/* cache hit */

	defer l.Unlock()

	item := l.Entries[index]

	l.insert(item)

	return item.Value, true, nil
}

func (l *Cache[K, V]) insert(i Entry[K, V]) {
	entries := []Entry[K, V]{i}

	entries = append(entries, l.Entries...)

	l.Entries = entries

	if len(l.Entries) > l.Capacity {
		l.Entries = l.Entries[0 : len(l.Entries)-1]
	}
}

// New initialzes a new LRU cache.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		Capacity: capacity,
		Entries:  make([]Entry[K, V], capacity),
	}
}
