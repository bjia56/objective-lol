"use strict";

const fs = require("fs");
const path = require("path");
const koffi = require("koffi");

const LIB_NAMES = { linux: "libolol.so", darwin: "libolol.dylib", win32: "olol.dll" };

function resolveLibPath() {
  const platform = process.platform;
  const arch = process.arch;
  const libName = LIB_NAMES[platform];
  if (!libName) {
    throw new Error(`objective-lol: unsupported platform "${platform}" for the native backend`);
  }
  const libPath = path.join(__dirname, "..", "native", `${platform}-${arch}`, libName);
  if (!fs.existsSync(libPath)) {
    throw new Error(
      `objective-lol: native library not found at ${libPath}. Build it with ` +
        `\`npm run build:native\` (requires Go + a C toolchain for this platform), ` +
        `or use the browser/wasm entry point ("@objective-lol/core/browser") instead.`
    );
  }
  return libPath;
}

let state = null;

// Lazily loads the shared library and declares its C ABI. Deferred until
// first use so that requiring this module doesn't fail on a platform
// where only the wasm backend will actually be used.
function ensureLoaded() {
  if (state) return state;

  const lib = koffi.load(resolveLibPath());

  // Matches the `NativeCallback` typedef in js/go/native/main.go:
  // char* (*)(void *userData, char *argsJSON)
  const NativeCallbackProto = koffi.proto("NativeCallback", "str", ["void *", "str"]);

  const fns = {
    newVM: lib.func("char *olol_new_vm(str optionsJSON)"),
    freeVM: lib.func("char *olol_free_vm(int64 handle)"),
    execute: lib.func("char *olol_execute(int64 handle, str code)"),
    call: lib.func("char *olol_call(int64 handle, str name, str argsJSON)"),
    callMethod: lib.func("char *olol_call_method(int64 handle, str objectJSON, str name, str argsJSON)"),
    defineVariable: lib.func("char *olol_define_variable(int64 handle, str name, str valueJSON, int constant)"),
    setVariable: lib.func("char *olol_set_variable(int64 handle, str name, str valueJSON)"),
    getVariable: lib.func("char *olol_get_variable(int64 handle, str name)"),
    newObjectInstance: lib.func("char *olol_new_object_instance(int64 handle, str className)"),
    defineFunction: lib.func(
      "char *olol_define_function(int64 handle, str name, int argc, NativeCallback *cb, void *userData)"
    ),
  };

  state = { lib, fns, NativeCallbackProto };
  return state;
}

class NodeBackend {
  constructor() {
    const { fns } = ensureLoaded();
    this._fns = fns;
  }

  async newVM(optionsJSON) {
    return this._fns.newVM(optionsJSON || "");
  }

  async freeVM(handle) {
    return this._fns.freeVM(handle);
  }

  async execute(handle, code) {
    return this._fns.execute(handle, code);
  }

  async call(handle, name, argsJSON) {
    return this._fns.call(handle, name, argsJSON);
  }

  async callMethod(handle, objectJSON, name, argsJSON) {
    return this._fns.callMethod(handle, objectJSON, name, argsJSON);
  }

  async defineVariable(handle, name, valueJSON, constant) {
    return this._fns.defineVariable(handle, name, valueJSON, constant ? 1 : 0);
  }

  async setVariable(handle, name, valueJSON) {
    return this._fns.setVariable(handle, name, valueJSON);
  }

  async getVariable(handle, name) {
    return this._fns.getVariable(handle, name);
  }

  async newObjectInstance(handle, className) {
    return this._fns.newObjectInstance(handle, className);
  }

  // Registers a synchronous JS callback as a persistent (non-transient)
  // native function pointer, since Objective-LOL may call it long after
  // olol_define_function itself returns. Returns an opaque handle to
  // pass to defineFunction, plus a release() to free the slot (Koffi
  // allows at most 8192 concurrent registrations).
  registerCallback(nativeCallback) {
    const { NativeCallbackProto } = state;
    const pointer = koffi.register((_userData, argsJSON) => nativeCallback(argsJSON), koffi.pointer(NativeCallbackProto));
    return {
      handle: pointer,
      release: () => koffi.unregister(pointer),
    };
  }

  async defineFunction(handle, name, argc, callbackHandle) {
    return this._fns.defineFunction(handle, name, argc, callbackHandle, null);
  }
}

module.exports = { NodeBackend };
