//go:build !race

package perp_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/go-perf/perp"
)

func TestConcurrentHitRatioTestNoRace(t *testing.T) {
	expectedRationWithoutRaceDetector := 0.8
	ratio := runConcurrentHitRatioTest(t)
	if math.Abs(ratio-expectedRationWithoutRaceDetector) > 0.01 {
		t.Logf("act ratio: %+v", ratio)
		t.Fatalf("actual ratio is too small: %v", ratio)
	}
}

func TestCache(t *testing.T) {
	cache, err := perp.NewCache[string, any](10)
	if err != nil {
		t.Fatal(err)
	}

	{
		v, ok := cache.Load("k1")
		eq(t, false, ok)
		eq(t, nil, v)
	}

	{
		cache.Store("k1", "value1")
		v, ok := cache.Load("k1")
		eq(t, true, ok)
		eq(t, "value1", v)
	}
}

func eq(t *testing.T, exp interface{}, act interface{}) {
	t.Helper()
	if reflect.DeepEqual(exp, act) {
		return
	}
	t.Logf("exp: %T=`%+v`", exp, exp)
	t.Logf("act: %T=`%+v`", act, act)
	t.Fatalf("assert failed")
}
