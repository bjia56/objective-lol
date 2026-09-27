"use strict";

const { parseEnvelope, fromWire, toWire } = require("./value");

// ObjectiveLOLVM is the platform-agnostic API surface. It's constructed
// with a `backend` implementing the low-level operations (see
// backend-node.js / backend-wasm.js), so the same class works whether
// the interpreter is running as a native addon or as wasm.
class ObjectiveLOLVM {
  constructor(backend) {
    this._backend = backend;
    this._handle = null;
    this._registeredCallbacks = [];
  }

  static async create(backend, options = {}) {
    const vm = new ObjectiveLOLVM(backend);
    const optionsJSON = JSON.stringify({
      workingDirectory: options.workingDirectory,
      timeoutMs: options.timeoutMs,
    });
    vm._handle = parseEnvelope(await backend.newVM(optionsJSON));
    return vm;
  }

  async execute(code) {
    return fromWire(this, parseEnvelope(await this._backend.execute(this._handle, code)));
  }

  async call(name, ...args) {
    const argsJSON = JSON.stringify(args.map(toWire));
    return fromWire(this, parseEnvelope(await this._backend.call(this._handle, name, argsJSON)));
  }

  async callMethod(ref, methodName, args = []) {
    const objectJSON = JSON.stringify(toWire(ref));
    const argsJSON = JSON.stringify(args.map(toWire));
    return fromWire(
      this,
      parseEnvelope(await this._backend.callMethod(this._handle, objectJSON, methodName, argsJSON))
    );
  }

  async defineVariable(name, value, constant = false) {
    parseEnvelope(
      await this._backend.defineVariable(this._handle, name, JSON.stringify(toWire(value)), constant)
    );
  }

  async setVariable(name, value) {
    parseEnvelope(await this._backend.setVariable(this._handle, name, JSON.stringify(toWire(value))));
  }

  async getVariable(name) {
    return fromWire(this, parseEnvelope(await this._backend.getVariable(this._handle, name)));
  }

  async newObjectInstance(className) {
    return fromWire(this, parseEnvelope(await this._backend.newObjectInstance(this._handle, className)));
  }

  // Defines a global Objective-LOL function backed by `fn`. `fn` MUST be
  // synchronous (return a plain value, not a Promise): both the native
  // and wasm backends invoke it as a direct, blocking call from deep
  // inside the tree-walking interpreter, with no way to suspend and
  // resume that call later with an async result.
  async defineFunction(name, argc, fn) {
    const nativeCallback = (argsJSON) => {
      try {
        const args = JSON.parse(argsJSON).map((v) => fromWire(this, v));
        const result = fn(...args);
        if (result && typeof result.then === "function") {
          throw new Error(
            `defineFunction("${name}") callback returned a Promise; callbacks must be synchronous`
          );
        }
        return JSON.stringify({ ok: true, value: toWire(result === undefined ? null : result) });
      } catch (err) {
        return JSON.stringify({ ok: false, error: { type: "runtime", message: String(err && err.message ? err.message : err) } });
      }
    };

    const registered = this._backend.registerCallback(nativeCallback);
    this._registeredCallbacks.push(registered);
    parseEnvelope(
      await this._backend.defineFunction(this._handle, name, argc, registered.handle)
    );
  }

  async free() {
    if (this._handle == null) return;
    parseEnvelope(await this._backend.freeVM(this._handle));
    this._handle = null;
    for (const cb of this._registeredCallbacks) {
      if (cb.release) cb.release();
    }
    this._registeredCallbacks = [];
  }
}

module.exports = { ObjectiveLOLVM };
