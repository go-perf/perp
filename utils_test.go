package perp

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_getProcID(t *testing.T) {
	seen := make([]int64, runtime.GOMAXPROCS(0))
	const iter = 1_000_000

	startCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(seen))

	for range seen {
		go func() {
			defer wg.Done()
			<-startCh
			for i := 0; i < iter; i++ {
				atomic.AddInt64(&seen[getProcID()], 1)
			}
		}()
	}

	close(startCh)
	wg.Wait()

	var sum int64
	for _, n := range seen {
		sum += n
		if n == 0 {
			t.Fatal()
		}
	}

	mustEqual(t, sum, int64(iter*len(seen)))
}
