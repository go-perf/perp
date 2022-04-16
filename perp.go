package perp

import (
	"runtime"
)

type Values[T any] struct {
	shards []T
}

func NewValues[T any](setDefault func() T) *Values[T] {
	shards := make([]T, runtime.GOMAXPROCS(0))
	for i := range shards {
		shards[i] = setDefault()
	}
	return &Values[T]{shards: shards}
}

func (vs *Values[T]) Get() T {
	return vs.shards[getProcID()]
}
