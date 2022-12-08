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

		if cache.Size > cache.Capacity {
			t.Fatalf("too many entries; %d", cache.Size)
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

		if hit {
			t.Fatal("cache returned a hit in a cold start")
		}
	}

	for i := 0; i < 100; i++ {
		i := i

		_, hit, _ := cache.Fetch(fmt.Sprintf("%d", i), func() (int, error) {
			return i, nil
		})

		if !hit {
			t.Fatal("cache returned a miss when the cache was pre-loaded")
		}
	}

	_, hit, _ := cache.Fetch("missing-no", func() (int, error) {
		return 1337, nil
	})

	if hit {
		t.Fatal("a known miss was a hit")
	}
}

func TestCacheEjection(t *testing.T) {
	cache := New[string, int](10)

	for i := 0; i < 10; i++ {
		i := i

		_, hit, _ := cache.Fetch(fmt.Sprintf("%d", i), func() (int, error) {
			return i, nil
		})

		if hit {
			t.Fatal("cache returned a hit in a cold start")
		}
	}

	/* "0" should be present */
	_, hit, _ := cache.Fetch("0", func() (int, error) {
		return 0, nil
	})

	if !hit {
		t.Fatal("a known hit was a miss")
	}

	if cache.tail.Key != "0" {
		t.Fatalf("known tail key was wrong; got %s", cache.tail.Key)
	}

	/* this new key should eject "0" */
	_, hit, _ = cache.Fetch("missing-no", func() (int, error) {
		return 1337, nil
	})

	if hit {
		t.Fatal("a known miss was a hit")
	}

	_, hit, _ = cache.Fetch("0", func() (int, error) {
		return 0, nil
	})

	if hit {
		t.Fatal("a known miss was a hit")
	}
}

func TestCacheClear(t *testing.T) {
	cache := New[string, int](10)

	for i := 0; i < 10; i++ {
		i := i

		_, hit, _ := cache.Fetch(fmt.Sprintf("%d", i), func() (int, error) {
			return i, nil
		})

		if hit {
			t.Fatal("cache returned a hit in a cold start")
		}
	}

	cache.Clear()

	for i := 0; i < 10; i++ {
		i := i

		_, hit, _ := cache.Fetch(fmt.Sprintf("%d", i), func() (int, error) {
			return i, nil
		})

		if hit {
			t.Fatal("cache returned a hit in a cold start")
		}
	}
}
