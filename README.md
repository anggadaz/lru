# LRU

A generic least recently used data structure.

## Usage

```golang
  import "github.com/prophittcorey/lru"

	cache := lru.New[string, int](10)

	integer, hit, err := cache.Fetch("some-unique-key", func() (int, error) {
		return getRandomIntExpensively(), nil
	})
```

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
