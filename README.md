# LRU

[![Go Reference](https://pkg.go.dev/badge/github.com/prophittcorey/lru.svg)](https://pkg.go.dev/github.com/prophittcorey/lru)

A generic least recently used data structure. This data structure was
implemented with O(1) insertions and O(1) lookups at the expense of a little
extra memory usage.

## Usage

```golang
import "github.com/prophittcorey/lru"

cache := lru.New[string, int](10)

integer, hit, err := cache.Fetch("some-unique-key", func() (int, error) {
  return getExpensiveValue(), nil
})
```

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
