package lru

import (
	"fmt"
	"testing"
)

func TestCapacity(t *testing.T) {
	cache := New[string, int](10)

	for i := 0; i < 100; i++ {
		i := i

		value, _, err := cache.Fetch(fmt.Sprintf("%d", i), func() (int, error) {
			return i, nil
		})

		if err != nil {
			t.Fatalf("err wasn't nil; %s", err)
		}

		if i != value {
			t.Fatalf("%d != %d", i, value)
		}

		if len(cache.Entries) > cache.Capacity {
			t.Fatalf("too many entries; %d", len(cache.Entries))
		}
	}
}

func TestCacheHitAndMiss(t *testing.T) {
	cache := New[string, int](100)

	for i := 0; i < 100; i++ {
		i := i

		_, hit, _ := cache.Fetch(fmt.Sprintf("%d", i), func() (int, error) {
			return i, nil
		})

		if hit != false {
			t.Fatal("cache returned a hit in a cold start")
		}
	}

	for i := 0; i < 100; i++ {
		i := i

		_, hit, _ := cache.Fetch(fmt.Sprintf("%d", i), func() (int, error) {
			return i, nil
		})

		if hit != true {
			t.Fatal("cache returned a miss when the cache was pre-loaded")
		}
	}

	_, hit, _ := cache.Fetch("missing-no", func() (int, error) {
		return 1337, nil
	})

	if hit == true {
		t.Fatal("a known miss was a hit")
	}
}
