// Package bridge implements a JSON-in/JSON-out calling convention over
// pkg/api, shared by the wasm (syscall/js) and native (cgo c-shared)
// entry points under js/go. Keeping this logic here means the two
// platform-specific main packages only have to adapt argument/return
// marshaling to their respective FFI, not reimplement VM semantics.
//
// Every exported function returns a single JSON "envelope" string:
//
//	{"ok": true, "value": <json>}
//	{"ok": false, "error": {"type": "runtime", "message": "..."}}
//
// so that callers on both sides only ever need a string in the return
// type and never need to construct language-specific error objects
// from the FFI plumbing itself.
package bridge

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bjia56/objective-lol/pkg/api"
)

type envelope struct {
	OK    bool        `json:"ok"`
	Value interface{} `json:"value,omitempty"`
	Error *errorInfo  `json:"error,omitempty"`
}

type errorInfo struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func ok(value interface{}) string {
	b, err := json.Marshal(envelope{OK: true, Value: value})
	if err != nil {
		return failString(fmt.Sprintf("failed to marshal result: %v", err))
	}
	return string(b)
}

func fail(err error) string {
	errType := "runtime"
	if vmErr, ok := err.(*api.VMError); ok {
		errType = string(vmErr.Type)
	}
	return failString2(errType, err.Error())
}

func failString(message string) string {
	return failString2("internal", message)
}

func failString2(errType, message string) string {
	b, _ := json.Marshal(envelope{OK: false, Error: &errorInfo{Type: errType, Message: message}})
	return string(b)
}

// --- VM registry ---

var (
	vmMu      sync.RWMutex
	vms       = make(map[int64]*api.VM)
	nextHandl int64
)

func getVM(handle int64) (*api.VM, error) {
	vmMu.RLock()
	defer vmMu.RUnlock()
	vm, found := vms[handle]
	if !found {
		return nil, fmt.Errorf("no VM found for handle %d", handle)
	}
	return vm, nil
}

// vmOptions is the JSON-decodable subset of api.VMConfig that makes sense
// to expose across the FFI boundary. Stdout/Stdin are intentionally left
// to the Go-side defaults (os.Stdout/os.Stdin); redirecting I/O across
// the boundary is not supported in this first cut.
type vmOptions struct {
	WorkingDirectory string `json:"workingDirectory"`
	TimeoutMS        int64  `json:"timeoutMs"`
}

// NewVM creates a new VM instance and returns its integer handle (as the
// envelope's "value") for use in subsequent calls.
func NewVM(optionsJSON string) string {
	// "/" rather than "." because api.VM.initialize() resolves this via
	// filepath.Abs, which calls os.Getwd() for a relative path - and
	// Getwd is unimplemented under GOOS=js (wasm/browser). An absolute
	// default sidesteps that call on every backend, not just wasm.
	opts := vmOptions{WorkingDirectory: "/"}
	if optionsJSON != "" {
		if err := json.Unmarshal([]byte(optionsJSON), &opts); err != nil {
			return failString2("config", fmt.Sprintf("invalid options JSON: %v", err))
		}
	}

	cfg := api.DefaultConfig()
	if opts.WorkingDirectory != "" {
		cfg.WorkingDirectory = opts.WorkingDirectory
	}
	if opts.TimeoutMS > 0 {
		cfg.Timeout = time.Duration(opts.TimeoutMS) * time.Millisecond
	}

	vm, err := api.NewVM(cfg)
	if err != nil {
		return fail(err)
	}

	handle := atomic.AddInt64(&nextHandl, 1)
	vmMu.Lock()
	vms[handle] = vm
	vmMu.Unlock()

	return ok(handle)
}

// FreeVM releases a VM handle. VMs cannot be reused after Execute anyway
// (see api.VM.Execute), so this mainly just drops the reference so it can
// be garbage collected.
func FreeVM(handle int64) string {
	vmMu.Lock()
	delete(vms, handle)
	vmMu.Unlock()
	return ok(nil)
}

// Execute runs Objective-LOL source code on the given VM.
func Execute(handle int64, code string) string {
	vm, err := getVM(handle)
	if err != nil {
		return fail(err)
	}
	result, err := vm.Execute(code)
	if err != nil {
		return fail(err)
	}
	return ok(result.Value)
}

// Call invokes a top-level Objective-LOL function. argsJSON must decode to
// a JSON array.
func Call(handle int64, name string, argsJSON string) string {
	vm, err := getVM(handle)
	if err != nil {
		return fail(err)
	}
	args, err := decodeArgs(argsJSON)
	if err != nil {
		return fail(err)
	}
	result, err := vm.Call(name, args)
	if err != nil {
		return fail(err)
	}
	return ok(result)
}

// CallMethod invokes a method on an Objective-LOL object. objectJSON must
// decode to an object reference (as previously returned in a "value",
// i.e. containing "__GoValue_id"), and argsJSON must decode to a JSON
// array.
func CallMethod(handle int64, objectJSON string, name string, argsJSON string) string {
	vm, err := getVM(handle)
	if err != nil {
		return fail(err)
	}
	object, err := decodeValue(objectJSON)
	if err != nil {
		return fail(err)
	}
	args, err := decodeArgs(argsJSON)
	if err != nil {
		return fail(err)
	}
	result, err := vm.CallMethod(object, name, args)
	if err != nil {
		return fail(err)
	}
	return ok(result)
}

// DefineVariable defines a global variable in the VM.
func DefineVariable(handle int64, name string, valueJSON string, constant bool) string {
	vm, err := getVM(handle)
	if err != nil {
		return fail(err)
	}
	value, err := decodeValue(valueJSON)
	if err != nil {
		return fail(err)
	}
	if err := vm.DefineVariable(name, value, constant); err != nil {
		return fail(err)
	}
	return ok(nil)
}

// SetVariable sets (or creates) a global variable in the VM.
func SetVariable(handle int64, name string, valueJSON string) string {
	vm, err := getVM(handle)
	if err != nil {
		return fail(err)
	}
	value, err := decodeValue(valueJSON)
	if err != nil {
		return fail(err)
	}
	if err := vm.SetVariable(name, value); err != nil {
		return fail(err)
	}
	return ok(nil)
}

// GetVariable reads a global variable from the VM.
func GetVariable(handle int64, name string) string {
	vm, err := getVM(handle)
	if err != nil {
		return fail(err)
	}
	value, err := vm.GetVariable(name)
	if err != nil {
		return fail(err)
	}
	return ok(value)
}

// NewObjectInstance constructs a new instance of an Objective-LOL class.
func NewObjectInstance(handle int64, className string) string {
	vm, err := getVM(handle)
	if err != nil {
		return fail(err)
	}
	value, err := vm.NewObjectInstance(className)
	if err != nil {
		return fail(err)
	}
	return ok(value)
}

// NativeCallback is the calling convention a platform adapter must
// implement in order to expose a host-language function to
// DefineFunction. It receives a JSON-encoded array of arguments and must
// return an envelope string (as produced by ok()/fail() in this
// package) so that host-side errors and panics propagate as ordinary
// Objective-LOL exceptions rather than crashing the interpreter.
type NativeCallback func(argsJSON string) string

// DefineFunction defines a global function in the VM backed by a
// host-language callback.
func DefineFunction(handle int64, name string, argc int, callback NativeCallback) string {
	vm, err := getVM(handle)
	if err != nil {
		return fail(err)
	}

	err = vm.DefineFunction(name, argc, func(args []api.GoValue) (api.GoValue, error) {
		argsJSON, err := json.Marshal(args)
		if err != nil {
			return api.WrapAny(nil), fmt.Errorf("error encoding arguments: %v", err)
		}

		resultEnvelope := callback(string(argsJSON))

		var env envelope
		if err := json.Unmarshal([]byte(resultEnvelope), &env); err != nil {
			return api.WrapAny(nil), fmt.Errorf("host callback returned invalid envelope: %v", err)
		}
		if !env.OK {
			msg := "host callback failed"
			if env.Error != nil {
				msg = env.Error.Message
			}
			return api.WrapAny(nil), fmt.Errorf("%s", msg)
		}
		return api.WrapAny(env.Value), nil
	})
	if err != nil {
		return fail(err)
	}
	return ok(nil)
}

func decodeArgs(argsJSON string) ([]api.GoValue, error) {
	if argsJSON == "" {
		return nil, nil
	}
	var raw []interface{}
	if err := json.Unmarshal([]byte(argsJSON), &raw); err != nil {
		return nil, fmt.Errorf("invalid arguments JSON: %v", err)
	}
	args := make([]api.GoValue, len(raw))
	for i, v := range raw {
		args[i] = api.WrapAny(v)
	}
	return args, nil
}

func decodeValue(valueJSON string) (api.GoValue, error) {
	if valueJSON == "" || valueJSON == "null" {
		return api.WrapAny(nil), nil
	}
	var raw interface{}
	if err := json.Unmarshal([]byte(valueJSON), &raw); err != nil {
		return api.WrapAny(nil), fmt.Errorf("invalid value JSON: %v", err)
	}
	return api.WrapAny(raw), nil
}
