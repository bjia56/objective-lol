"use strict";

const path = require("path");
const fs = require("fs");

const { ObjectiveLOLVM } = require("./vm");
const { NodeBackend } = require("./backend-node");
const { WasmBackend } = require("./backend-wasm");
const { ObjectiveLOLError } = require("./value");

// Two wasm builds ship side by side (see scripts/build-wasm.js): the
// standard Go compiler's output, and a smaller/narrower TinyGo one.
// Picking one here only matters for testing wasm under Node - a real
// browser deployment picks a variant simply by which wasm/wasm_exec URL
// its own code references (see index.browser.js and the README), so a
// bundler's asset step naturally never ships the other one.
const WASM_VARIANTS = {
  go: {
    wasmPath: path.join(__dirname, "..", "wasm", "objectivelol.wasm"),
    wasmExecPath: path.join(__dirname, "..", "wasm", "wasm_exec.js"),
  },
  tinygo: {
    wasmPath: path.join(__dirname, "..", "wasm", "objectivelol-tinygo.wasm"),
    wasmExecPath: path.join(__dirname, "..", "wasm", "wasm_exec-tinygo.js"),
  },
};

// Creates a ready-to-use VM. By default this uses the native addon
// backend (fastest, full FILE/THREAD stdlib support). Pass
// `{ backend: "wasm" }` to run the same wasm build used in the browser
// instead - useful for testing it under Node, or on a platform without a
// prebuilt native library. `wasmVariant` picks "go" (default) or
// "tinygo"; `wasmPath`/`wasmExecPath` override either explicitly.
async function createVM(options = {}) {
  let backend;
  if (options.backend === "wasm") {
    const variant = WASM_VARIANTS[options.wasmVariant || "go"];
    if (!variant) {
      throw new Error(`createVM: unknown wasmVariant "${options.wasmVariant}" (expected "go" or "tinygo")`);
    }
    backend = await WasmBackend.create({
      wasmBytes: fs.readFileSync(options.wasmPath || variant.wasmPath),
      wasmExecPath: options.wasmExecPath || variant.wasmExecPath,
    });
  } else {
    backend = new NodeBackend();
  }
  return ObjectiveLOLVM.create(backend, options);
}

module.exports = { ObjectiveLOLVM, NodeBackend, WasmBackend, ObjectiveLOLError, createVM };
