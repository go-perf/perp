package perp_test

import (
	"fmt"
	"hash/maphash"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/go-perf/perp"
)

type cache[K comparable, V any] interface {
	Load(key K) (V, bool)
	Store(key K, value V)
}

var workerSets = []int{1, 2, 4, 8, 12, 16, 22, 28, 32}

func BenchmarkCacheLoadStore(b *testing.B) {
	for _, r := range workerSets {
		cache, err := perp.NewCache[string, string](80)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("w-%v", r), func(b *testing.B) {
			benchLoadStore(b, r, cache)
		})
	}
}

func BenchmarkMutexCacheLoadStore(b *testing.B) {
	for _, r := range workerSets {
		cache := newMutexCache[string, string](80)
		b.Run(fmt.Sprintf("w-%v", r), func(b *testing.B) {
			benchLoadStore(b, r, cache)
		})
	}
}

func BenchmarkStripedCacheLoadStore(b *testing.B) {
	for _, r := range workerSets {
		cache := newStripedMapCache[string, string](80)
		b.Run(fmt.Sprintf("w-%v", r), func(b *testing.B) {
			benchLoadStore(b, r, cache)
		})
	}
}

func benchLoadStore(b *testing.B, workers int, cache cache[string, string]) {
	keys := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		keys = append(keys, strconv.Itoa(i))
	}

	count := uint64(0)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			idx := 0
			hits := uint64(0)
			for i := 0; i < b.N; i++ {
				idx++
				if idx == len(keys) {
					idx = 0
				}
				cache.Store(keys[len(keys)-1-idx], keys[idx])
				_, ok := cache.Load(keys[idx])
				if ok {
					hits++
				}
			}
			atomic.AddUint64(&count, hits)
			wg.Done()
		}()
	}
	wg.Wait()
	if rand.Float32() > 2 {
		b.Logf("hits: %v", atomic.LoadUint64(&count))
	}
}

type stripe[V any] struct {
	m     sync.RWMutex
	store map[string]V
	_     [32]byte
}

type stripedMapCache[V any] struct {
	sizePerGoroutine int
	stripes          []*stripe[V]
}

func newStripedMapCache[K comparable, V any](sizePerGoroutine uint) *stripedMapCache[V] {
	cache := &stripedMapCache[V]{
		sizePerGoroutine: int(sizePerGoroutine),
		stripes:          make([]*stripe[V], 0, 64),
	}
	for i := 0; i < 64; i++ {
		cache.stripes = append(cache.stripes, &stripe[V]{
			m:     sync.RWMutex{},
			store: make(map[string]V),
		})
	}
	return cache
}

func (m *stripedMapCache[V]) Load(key string) (V, bool) {
	var h maphash.Hash
	_, _ = h.WriteString(key)
	idx := uint64(64-1) & h.Sum64()
	stripe := m.stripes[idx]
	stripe.m.RLock()
	defer stripe.m.RUnlock()
	val, ok := stripe.store[key]
	return val, ok
}

func (m *stripedMapCache[V]) Store(key string, value V) {
	var h maphash.Hash
	_, _ = h.WriteString(key)
	stripeIdx := uint64(64-1) & h.Sum64()
	stripe := m.stripes[stripeIdx]
	stripe.m.Lock()
	defer stripe.m.Unlock()

	stripe.store[key] = value
	if len(stripe.store) <= m.sizePerGoroutine {
		return
	}
	for k := range stripe.store {
		delete(stripe.store, k)
		break
	}
}

type mutexMapCache[K comparable, V any] struct {
	m                sync.RWMutex
	sizePerGoroutine int
	store            map[K]V
}

func newMutexCache[K comparable, V any](sizePerGoroutine uint) *mutexMapCache[K, V] {
	return &mutexMapCache[K, V]{
		m:                sync.RWMutex{},
		sizePerGoroutine: int(sizePerGoroutine),
		store:            make(map[K]V),
	}
}

func (m *mutexMapCache[K, V]) Load(key K) (V, bool) {
	m.m.RLock()
	defer m.m.RUnlock()
	val, ok := m.store[key]
	return val, ok
}

func (m *mutexMapCache[K, V]) Store(key K, value V) {
	m.m.Lock()
	defer m.m.Unlock()
	m.store[key] = value
	if len(m.store) <= m.sizePerGoroutine {
		return
	}
	for k := range m.store {
		delete(m.store, k)
		break
	}
}
