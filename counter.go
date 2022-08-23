package perp

import (
	"sync/atomic"
)

// Counter is sharded across processors in Go runtime.
type Counter struct {
	vs *Values[*counter]
}

// NewCounter returns a new Counter initialized to zero.
func NewCounter() *Counter {
	vs := NewValues(func() *counter {
		return &counter{}
	})
	return &Counter{vs: vs}
}

// Add n to the counter.
func (c *Counter) Add(n int64) {
	c.vs.Get().Add(n)
}

// Load the total counter value.
func (c *Counter) Load() int64 {
	var sum int64
	c.vs.Iter(func(v *counter) {
		sum += v.Load()
	})
	return sum
}

// Reset the counter to zero and return the old value.
func (c *Counter) Reset() int64 {
	var sum int64
	c.vs.Iter(func(v *counter) {
		sum += v.Reset()
	})
	return sum
}

type counter struct {
	n int64
	_ [56]byte // cache line aligment
}

func (c *counter) Add(n int64) int64 {
	return atomic.AddInt64(&c.n, n)
}

func (c *counter) Load() int64 {
	return atomic.LoadInt64(&c.n)
}

func (c *counter) Reset() int64 {
	return atomic.SwapInt64(&c.n, 0)
}
