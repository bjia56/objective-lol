package api

import (
	"runtime"
	"testing"
)

// With caller-thread callbacks enabled, a callback must run on the goroutine
// that entered the VM (pinned to one OS thread), not on a fresh one.
func TestCallbacksRunOnCallerThread(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		vm, err := NewVM(DefaultConfig())
		if err != nil {
			t.Fatal(err)
		}
		vm.SetCallbacksOnCallerThread(enabled)

		ranOnCaller := make(chan bool, 1)
		var callerStack string
		shim := vm.GetCompatibilityShim()
		err = shim.DefineFunction("id", "PING", 0, func(id, args string) string {
			ranOnCaller <- stackID() == callerStack
			return `{"result": 1}`
		})
		if err != nil {
			t.Fatal(err)
		}

		runtime.LockOSThread()
		callerStack = stackID()
		_, err = vm.Execute("I CAN HAS STDIO?\nHAI ME TEH FUNCSHUN MAIN\n SAYZ WIT PING\nKTHXBAI\n")
		runtime.UnlockOSThread()
		if err != nil {
			t.Fatal(err)
		}
		if got := <-ranOnCaller; got != enabled {
			t.Errorf("enabled=%v: callback ran on caller goroutine = %v", enabled, got)
		}
	}
}

// stackID returns the current goroutine's id line ("goroutine N [running]:").
func stackID() string {
	buf := make([]byte, 64)
	buf = buf[:runtime.Stack(buf, false)]
	for i, b := range buf {
		if b == '[' {
			return string(buf[:i])
		}
	}
	return string(buf)
}
