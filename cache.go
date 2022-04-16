package perp

import (
	"errors"
	"sync"
)

type Cache[K comparable, V any] struct {
	values  *Values[safeCache[K, V]]
	maxSize int
}

func NewCache[K comparable, V any](maxSizePerP int) (*Cache[K, V], error) {
	if maxSizePerP <= 0 {
		return nil, errors.New("maxSizePerP must be greater than zero")
	}

	c := &Cache[K, V]{
		maxSize: maxSizePerP,
		values: NewValues(func() safeCache[K, V] {
			return safeCache[K, V]{
				data: make(map[K]V, maxSizePerP),
			}
		}),
	}
	return c, nil
}

func (c *Cache[K, V]) Load(key K) (V, bool) {
	cache := c.values.Get()
	return cache.Get(key)
}

func (c *Cache[K, V]) Store(key K, value V) {
	cache := c.values.Get()
	cache.Set(c.maxSize, key, value)
}

type safeCache[K comparable, V any] struct {
	sync.RWMutex
	data map[K]V
}

func (sc *safeCache[K, V]) Get(key K) (V, bool) {
	sc.RLock()
	defer sc.RUnlock()
	v, ok := sc.data[key]
	return v, ok
}

func (sc *safeCache[K, V]) Set(maxsize int, key K, value V) {
	sc.Lock()
	defer sc.Unlock()
	if len(sc.data) == maxsize {
		for k := range sc.data {
			delete(sc.data, k)
			break
		}
	}
	sc.data[key] = value
}
