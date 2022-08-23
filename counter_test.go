package perp

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestCounter(t *testing.T) {
	c := NewCounter()

	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < n; i++ {
				c.Add(1)
			}
		}()
	}
	wg.Wait()

	mustEqual(t, c.Load(), int64(n*n))
}

func TestCounterReset(t *testing.T) {
	c := NewCounter()

	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)

	var sum int64

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < n; i++ {
				c.Add(1)
				if i%32 == 0 {
					atomic.AddInt64(&sum, c.Reset())
				}
			}
		}()
	}

	wg.Wait()

	sum += c.Reset()
	mustEqual(t, sum, int64(n*n))
	mustEqual(t, c.Load(), int64(0))
}

func BenchmarkCounter(b *testing.B) {
	c := NewCounter()

	for i := 0; i < b.N; i++ {
		c.Add(1)
	}
	mustEqual(b, c.Load(), int64(b.N))
}

func BenchmarkCounterMutex(b *testing.B) {
	var c mutexCounter

	for i := 0; i < b.N; i++ {
		c.Add(1)
	}
	mustEqual(b, c.Load(), int64(b.N))
}

func BenchmarkCounterAtomic(b *testing.B) {
	var c atomicCounter

	for i := 0; i < b.N; i++ {
		c.Add(1)
	}
	mustEqual(b, c.Load(), int64(b.N))
}

func BenchmarkCounterPaddedAtomic(b *testing.B) {
	var c paddedAtomicCounter

	for i := 0; i < b.N; i++ {
		c.Add(1)
	}
	mustEqual(b, c.Load(), int64(b.N))
}

func BenchmarkParallelCounter(b *testing.B) {
	c := NewCounter()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Add(1)
		}
	})
	mustEqual(b, c.Load(), int64(b.N))
}

func BenchmarkParallelCounterMutex(b *testing.B) {
	var c mutexCounter

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Add(1)
		}
	})
	mustEqual(b, c.Load(), int64(b.N))
}

func BenchmarkParallelCounterAtomic(b *testing.B) {
	var c atomicCounter

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Add(1)
		}
	})
	mustEqual(b, c.Load(), int64(b.N))
}

func BenchmarkParallelCounterPaddedAtomic(b *testing.B) {
	var c paddedAtomicCounter

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Add(1)
		}
	})
	mustEqual(b, c.Load(), int64(b.N))
}

type mutexCounter struct {
	mu sync.Mutex
	n  int64
}

func (c *mutexCounter) Add(n int64) {
	c.mu.Lock()
	c.n += n
	c.mu.Unlock()
}

func (c *mutexCounter) Load() int64 {
	c.mu.Lock()
	v := c.n
	c.mu.Unlock()
	return v
}

type atomicCounter struct {
	n int64
}

func (c *atomicCounter) Add(n int64) {
	atomic.AddInt64(&c.n, n)
}

func (c *atomicCounter) Load() int64 {
	return atomic.LoadInt64(&c.n)
}

type paddedAtomicCounter struct {
	_ [56]byte
	n int64
}

func (c *paddedAtomicCounter) Add(n int64) {
	atomic.AddInt64(&c.n, n)
}

func (c *paddedAtomicCounter) Load() int64 {
	return atomic.LoadInt64(&c.n)
}
