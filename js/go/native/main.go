// Command native builds a C shared library (via `go build -buildmode=c-shared`)
// exposing js/go/bridge over a plain C ABI, so it can be loaded from
// Node with a runtime FFI library (koffi) without requiring a compiled
// node-gyp/N-API addon.
//
// Every exported function returns `char*` - a JSON envelope string (see
// js/go/bridge's doc comment) - marshaled across the FFI boundary using
// koffi's automatic `char *` <-> JS string conversion (koffi copies the
// bytes into a JS string during the call). Node therefore never has to
// hold or free the raw pointer itself; ownership of the underlying
// malloc'd buffer stays on the Go side. Since koffi's copy already
// happened by the time the call returns, each exported function frees
// the *previous* call's buffer (via freeManaged, guarded by freeMu)
// before allocating its own - there is only ever one call in flight at a
// time because every FFI call here is synchronous, so this never frees
// a buffer koffi is still reading.
//
// DefineFunction is the one operation that needs a callback *into* the
// host. It takes a C function pointer (see the NativeCallback typedef
// below) plus an opaque userData token; olol_define_function invokes it
// synchronously whenever the Objective-LOL side calls the defined
// function, mirroring the synchronous-only restriction on the wasm side.
package main

/*
#include <stdlib.h>

// Matches the callback signature Node registers via koffi.register():
// (userData: void*, argsJSON: char*) -> char*. The returned buffer's
// memory is owned by the host FFI layer (koffi manages it internally for
// registered callbacks that return a string), not by Go.
typedef char* (*NativeCallback)(void* userData, char* argsJSON);

static char* callNativeCallback(NativeCallback cb, void* userData, char* argsJSON) {
    return cb(userData, argsJSON);
}
*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/bjia56/objective-lol/js/go/bridge"
)

func main() {}

var (
	freeMu  sync.Mutex
	pending []unsafe.Pointer
)

// returnManaged frees the buffer(s) returned by the previous call from
// this library before allocating a new one for s.
func returnManaged(s string) *C.char {
	freeMu.Lock()
	defer freeMu.Unlock()

	for _, p := range pending {
		C.free(p)
	}
	pending = pending[:0]

	cstr := C.CString(s)
	pending = append(pending, unsafe.Pointer(cstr))
	return cstr
}

//export olol_new_vm
func olol_new_vm(optionsJSON *C.char) *C.char {
	return returnManaged(bridge.NewVM(C.GoString(optionsJSON)))
}

//export olol_free_vm
func olol_free_vm(handle C.int64_t) *C.char {
	return returnManaged(bridge.FreeVM(int64(handle)))
}

//export olol_execute
func olol_execute(handle C.int64_t, code *C.char) *C.char {
	return returnManaged(bridge.Execute(int64(handle), C.GoString(code)))
}

//export olol_call
func olol_call(handle C.int64_t, name *C.char, argsJSON *C.char) *C.char {
	return returnManaged(bridge.Call(int64(handle), C.GoString(name), C.GoString(argsJSON)))
}

//export olol_call_method
func olol_call_method(handle C.int64_t, objectJSON *C.char, name *C.char, argsJSON *C.char) *C.char {
	return returnManaged(bridge.CallMethod(int64(handle), C.GoString(objectJSON), C.GoString(name), C.GoString(argsJSON)))
}

//export olol_define_variable
func olol_define_variable(handle C.int64_t, name *C.char, valueJSON *C.char, constant C.int) *C.char {
	return returnManaged(bridge.DefineVariable(int64(handle), C.GoString(name), C.GoString(valueJSON), constant != 0))
}

//export olol_set_variable
func olol_set_variable(handle C.int64_t, name *C.char, valueJSON *C.char) *C.char {
	return returnManaged(bridge.SetVariable(int64(handle), C.GoString(name), C.GoString(valueJSON)))
}

//export olol_get_variable
func olol_get_variable(handle C.int64_t, name *C.char) *C.char {
	return returnManaged(bridge.GetVariable(int64(handle), C.GoString(name)))
}

//export olol_new_object_instance
func olol_new_object_instance(handle C.int64_t, className *C.char) *C.char {
	return returnManaged(bridge.NewObjectInstance(int64(handle), C.GoString(className)))
}

//export olol_define_function
func olol_define_function(handle C.int64_t, name *C.char, argc C.int, cb C.NativeCallback, userData unsafe.Pointer) *C.char {
	callback := bridge.NativeCallback(func(argsJSON string) string {
		cArgsJSON := C.CString(argsJSON)
		defer C.free(unsafe.Pointer(cArgsJSON))

		cResult := C.callNativeCallback(cb, userData, cArgsJSON)
		if cResult == nil {
			return `{"ok":false,"error":{"type":"runtime","message":"native callback returned null"}}`
		}
		return C.GoString(cResult)
	})

	return returnManaged(bridge.DefineFunction(int64(handle), C.GoString(name), int(argc), callback))
}
