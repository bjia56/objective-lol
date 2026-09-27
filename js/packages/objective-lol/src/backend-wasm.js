"use strict";

const { instantiateWasm } = require("./wasm-runtime");

class WasmBackend {
  constructor(bridge) {
    this._bridge = bridge;
  }

  static async create(opts = {}) {
    const bridge = await instantiateWasm(opts);
    return new WasmBackend(bridge);
  }

  async newVM(optionsJSON) {
    return this._bridge.newVM(optionsJSON || "");
  }

  async freeVM(handle) {
    return this._bridge.freeVM(handle);
  }

  async execute(handle, code) {
    return this._bridge.execute(handle, code);
  }

  async call(handle, name, argsJSON) {
    return this._bridge.call(handle, name, argsJSON);
  }

  async callMethod(handle, objectJSON, name, argsJSON) {
    return this._bridge.callMethod(handle, objectJSON, name, argsJSON);
  }

  async defineVariable(handle, name, valueJSON, constant) {
    return this._bridge.defineVariable(handle, name, valueJSON, !!constant);
  }

  async setVariable(handle, name, valueJSON) {
    return this._bridge.setVariable(handle, name, valueJSON);
  }

  async getVariable(handle, name) {
    return this._bridge.getVariable(handle, name);
  }

  async newObjectInstance(handle, className) {
    return this._bridge.newObjectInstance(handle, className);
  }

  // No native registration step is needed under wasm: Go holds the JS
  // function reference directly (as a js.Value) for as long as the
  // Objective-LOL function that wraps it exists.
  registerCallback(nativeCallback) {
    return { handle: nativeCallback, release: () => {} };
  }

  async defineFunction(handle, name, argc, callbackHandle) {
    return this._bridge.defineFunction(handle, name, argc, callbackHandle);
  }
}

module.exports = { WasmBackend };
