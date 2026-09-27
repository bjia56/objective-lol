"use strict";

const isNode = typeof process !== "undefined" && !!(process.versions && process.versions.node);

// Standard Go and TinyGo each ship a wasm_exec.js that defines an
// incompatible `globalThis.Go` runtime class under the same name (their
// import objects differ - e.g. TinyGo's wasm needs "wasi_snapshot_preview1"
// imports the standard Go runtime doesn't provide). So "Go is already
// defined" is only safe to treat as "already loaded" when it's the *same*
// wasm_exec.js as last time; switching variants within one process must
// force a reload rather than trusting whatever's already on globalThis.
let lastLoadedPath = null;

async function ensureGoRuntime(wasmExecPath) {
  if (lastLoadedPath === wasmExecPath && typeof globalThis.Go === "function") return;

  if (isNode) {
    // wasm_exec.js is a plain (non-module) script that assigns
    // `globalThis.Go` as a side effect - requiring it for that effect is
    // exactly what Go's own docs recommend running under Node. Evict any
    // cached copy first so switching variants re-executes it instead of
    // getting a no-op from require()'s module cache.
    const resolved = require.resolve(wasmExecPath);
    delete require.cache[resolved];
    require(resolved);
    lastLoadedPath = wasmExecPath;
    return;
  }

  if (typeof document !== "undefined") {
    await new Promise((resolve, reject) => {
      const script = document.createElement("script");
      script.src = wasmExecPath;
      script.onload = () => resolve();
      script.onerror = () => reject(new Error(`objective-lol: failed to load ${wasmExecPath}`));
      document.head.appendChild(script);
    });
    lastLoadedPath = wasmExecPath;
    return;
  }

  throw new Error(
    "objective-lol: don't know how to load wasm_exec.js in this environment; " +
      "load it yourself first so that globalThis.Go is defined (e.g. via a <script> " +
      "tag or importScripts() in a worker)."
  );
}

// Instantiates the Objective-LOL wasm module and runs it to completion of
// its setup phase (main() registers `globalThis.objectivelol` and then
// blocks forever, so go.run()'s promise intentionally is never awaited
// here - only the exported bridge object is).
async function instantiateWasm({ wasmBytes, wasmModule, wasmExecPath }) {
  if (!wasmModule && !wasmBytes) {
    throw new Error("objective-lol: instantiateWasm needs either wasmBytes or wasmModule");
  }

  await ensureGoRuntime(wasmExecPath);

  const go = new globalThis.Go();
  const result = wasmModule
    ? await WebAssembly.instantiate(wasmModule, go.importObject)
    : await WebAssembly.instantiate(wasmBytes, go.importObject);

  const instance = wasmModule ? result : result.instance;
  const ready = new Promise((resolve) => {
    // objectivelol is set synchronously near the top of main(), well
    // before the blocking select{}, but we still poll a tick to be safe
    // against future reordering.
    const check = () => {
      if (globalThis.objectivelol) resolve();
      else setTimeout(check, 0);
    };
    check();
  });

  go.run(instance).catch((err) => {
    // go.run's promise only settles if the wasm program exits/traps;
    // since main() blocks forever this should never fire in practice,
    // but surface it loudly if it ever does (e.g. a panic).
    console.error("objective-lol: wasm runtime exited unexpectedly:", err);
  });

  await ready;
  return globalThis.objectivelol;
}

module.exports = { instantiateWasm };
