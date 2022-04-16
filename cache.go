package perp

import (
	"errors"
	"sync"
)

type Cache[K comparable, V any] struct {
	maxSizePerP int
	pool        *sync.Pool
}

type cacheStripe[K comparable, V any] struct {
	cache map[K]V
}

func NewCache[K comparable, V any](maxSizePerP int) (*Cache[K, V], error) {
	if maxSizePerP <= 0 {
		return nil, errors.New("maxSizePerP must be greater than zero")
	}

	c := &Cache[K, V]{
		maxSizePerP: maxSizePerP,
		pool: &sync.Pool{
			New: func() interface{} {
				return &cacheStripe[K, V]{
					cache: make(map[K]V),
				}
			},
		},
	}
	return c, nil
}

func (c *Cache[K, V]) Load(key K) (V, bool) {
	stripe := c.pool.Get().(*cacheStripe[K, V])
	defer c.pool.Put(stripe)
	value, ok := stripe.cache[key]
	return value, ok
}

func (c *Cache[K, V]) Store(key K, value V) {
	stripe := c.pool.Get().(*cacheStripe[K, V])
	defer c.pool.Put(stripe)

	stripe.cache[key] = value
	if len(stripe.cache) > c.maxSizePerP {
		for k := range stripe.cache {
			delete(stripe.cache, k)
			break
		}
	}
}
