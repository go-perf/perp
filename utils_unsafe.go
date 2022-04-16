//go:build !purego

package perp

import (
	_ "unsafe" // for go:linkname
)

func getProcID() int {
	pid := runtimeProcPin()
	runtimeProcUnpin()
	return pid
}

//go:linkname runtimeProcPin runtime.procPin
func runtimeProcPin() int

//go:linkname runtimeProcUnpin runtime.procUnpin
func runtimeProcUnpin()
