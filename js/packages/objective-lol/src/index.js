"use strict";

const path = require("path");
const fs = require("fs");

const { ObjectiveLOLVM } = require("./vm");
const { NodeBackend } = require("./backend-node");
const { WasmBackend } = require("./backend-wasm");
const { ObjectiveLOLError } = require("./value");

const DEFAULT_WASM_PATH = path.join(__dirname, "..", "wasm", "objectivelol.wasm");
const DEFAULT_WASM_EXEC_PATH = path.join(__dirname, "..", "wasm", "wasm_exec.js");

// Creates a ready-to-use VM. By default this uses the native addon
// backend (fastest, full FILE/THREAD stdlib support). Pass
// `{ backend: "wasm" }` to run the same wasm build used in the browser
// instead - useful for testing it under Node, or on a platform without a
// prebuilt native library.
async function createVM(options = {}) {
  let backend;
  if (options.backend === "wasm") {
    backend = await WasmBackend.create({
      wasmBytes: fs.readFileSync(options.wasmPath || DEFAULT_WASM_PATH),
      wasmExecPath: options.wasmExecPath || DEFAULT_WASM_EXEC_PATH,
    });
  } else {
    backend = new NodeBackend();
  }
  return ObjectiveLOLVM.create(backend, options);
}

module.exports = { ObjectiveLOLVM, NodeBackend, WasmBackend, ObjectiveLOLError, createVM };
