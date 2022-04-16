//go:build purego

package perp

import (
	"runtime"
	"sync/atomic"
)

var (
	goMaxProcs = runtime.GOMAXPROCS(0)

	procIndex int32
)

// NOTE: not as good as unsafe (not purego) version
// at least we will spread data across all shards.
func getProcID() int {
	return int(atomic.AddInt32(&procIndex, 1)) % goMaxProcs
}
