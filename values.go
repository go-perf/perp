package perp

import "runtime"

// Values that are sharded across processors in Go runtime.
// Changing runtime.GOMAXPROCS might cause panics.
type Values[T any] struct {
	shards []T
}

// NewValues creates a new Values of type T.
func NewValues[T any](setDefault func() T) *Values[T] {
	shards := make([]T, runtime.GOMAXPROCS(0))
	for i := range shards {
		shards[i] = setDefault()
	}
	return &Values[T]{shards: shards}
}

// Get the value for current processor.
func (vs *Values[T]) Get() T {
	return vs.shards[getProcID()]
}

// Iter will iterate over all the values/processors.
func (vs *Values[T]) Iter(fn func(T)) {
	for _, pv := range vs.shards {
		fn(pv)
	}
}
