"use strict";

const { ObjectiveLOLVM } = require("./vm");
const { WasmBackend } = require("./backend-wasm");
const { ObjectiveLOLError } = require("./value");

// Creates a ready-to-use VM backed by WebAssembly, the only option in a
// browser. `wasmURL` should point at the published objectivelol.wasm
// asset (e.g. resolved by your bundler, such as
// `new URL('@objective-lol/core/wasm/objectivelol.wasm', import.meta.url)`);
// `wasmExecURL` likewise for wasm/wasm_exec.js. Both are fetched with the
// platform's global `fetch`/`<script>` loading, so no bundler-specific
// asset plugin is required.
async function createVM(options = {}) {
  if (!options.wasmURL) {
    throw new Error(
      "createVM: pass { wasmURL, wasmExecURL } pointing at the published " +
        "objectivelol.wasm / wasm_exec.js assets"
    );
  }

  const response = await fetch(options.wasmURL);
  if (!response.ok) {
    throw new Error(`objective-lol: failed to fetch ${options.wasmURL}: ${response.status}`);
  }
  const wasmBytes = await response.arrayBuffer();

  const backend = await WasmBackend.create({
    wasmBytes,
    wasmExecPath: options.wasmExecURL,
  });
  return ObjectiveLOLVM.create(backend, options);
}

module.exports = { ObjectiveLOLVM, WasmBackend, ObjectiveLOLError, createVM };
