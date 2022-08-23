package perp

import (
	"reflect"
	"testing"
)

func failIfOk(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fail()
	}
}

func failIfErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func mustEqual(t testing.TB, have, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(have, want) {
		t.Fatalf("\nhave: %+v\nwant: %+v\n", have, want)
	}
}
