//go:build js && wasm

// Command wasm builds the browser-facing Objective-LOL entry point via
// GOOS=js GOARCH=wasm. It exposes a global `objectivelol` object whose
// methods mirror js/go/bridge, adapted to JS calling conventions
// (arguments/returns are plain JS strings; callbacks are ordinary JS
// functions rather than C function pointers).
package main

import (
	"syscall/js"

	"github.com/bjia56/objective-lol/js/go/bridge"
)

// wrap adapts a Go function taking JSON-string args to a js.Func. Every
// exported bridge function already returns a JSON envelope string, so
// the adapter just has to shuttle JS strings/numbers in and out.
func wrap(fn func(args []js.Value) string) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return fn(args)
	})
}

func str(v js.Value) string {
	if v.IsUndefined() || v.IsNull() {
		return ""
	}
	return v.String()
}

func main() {
	exports := js.Global().Get("Object").New()

	exports.Set("newVM", wrap(func(args []js.Value) string {
		return bridge.NewVM(str(args[0]))
	}))

	exports.Set("freeVM", wrap(func(args []js.Value) string {
		return bridge.FreeVM(int64(args[0].Int()))
	}))

	exports.Set("execute", wrap(func(args []js.Value) string {
		return bridge.Execute(int64(args[0].Int()), str(args[1]))
	}))

	exports.Set("call", wrap(func(args []js.Value) string {
		return bridge.Call(int64(args[0].Int()), str(args[1]), str(args[2]))
	}))

	exports.Set("callMethod", wrap(func(args []js.Value) string {
		return bridge.CallMethod(int64(args[0].Int()), str(args[1]), str(args[2]), str(args[3]))
	}))

	exports.Set("defineVariable", wrap(func(args []js.Value) string {
		return bridge.DefineVariable(int64(args[0].Int()), str(args[1]), str(args[2]), args[3].Bool())
	}))

	exports.Set("setVariable", wrap(func(args []js.Value) string {
		return bridge.SetVariable(int64(args[0].Int()), str(args[1]), str(args[2]))
	}))

	exports.Set("getVariable", wrap(func(args []js.Value) string {
		return bridge.GetVariable(int64(args[0].Int()), str(args[1]))
	}))

	exports.Set("newObjectInstance", wrap(func(args []js.Value) string {
		return bridge.NewObjectInstance(int64(args[0].Int()), str(args[1]))
	}))

	// defineFunction(handle, name, argc, jsCallback) where jsCallback is a
	// *synchronous* JS function: (argsJSON: string) => string, returning a
	// bridge envelope. Synchronous only, because Go's wasm scheduler can
	// suspend a goroutine to let the JS event loop run, but the native
	// call stack that invoked this callback (deep inside the tree-walking
	// interpreter) has no way to "come back later" with a result -
	// there's no continuation to resume into. An async/Promise-returning
	// callback would need the interpreter call chain itself to be
	// async-aware, which pkg/interpreter is not.
	exports.Set("defineFunction", wrap(func(args []js.Value) string {
		handle := int64(args[0].Int())
		name := str(args[1])
		argc := args[2].Int()
		jsCallback := args[3]

		callback := bridge.NativeCallback(func(argsJSON string) string {
			result := jsCallback.Invoke(argsJSON)
			return result.String()
		})

		return bridge.DefineFunction(handle, name, argc, callback)
	}))

	js.Global().Set("objectivelol", exports)

	// Keep the wasm program alive; all work happens through the exported
	// callbacks above, driven by the JS host.
	select {}
}
