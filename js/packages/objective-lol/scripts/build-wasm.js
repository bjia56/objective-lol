#!/usr/bin/env node
// Builds the browser wasm artifact from js/go/wasm and stages wasm_exec.js
// alongside it. Requires a Go toolchain on PATH (no cgo needed for this
// target).
const { execFileSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const pkgRoot = path.resolve(__dirname, "..");
const repoRoot = path.resolve(pkgRoot, "..", "..", "..");
const goWasmPkg = path.join(repoRoot, "js", "go", "wasm");
const outDir = path.join(pkgRoot, "wasm");

fs.mkdirSync(outDir, { recursive: true });

const goroot = execFileSync("go", ["env", "GOROOT"], { encoding: "utf8" }).trim();
const wasmExecCandidates = [
  path.join(goroot, "lib", "wasm", "wasm_exec.js"), // Go >= 1.24
  path.join(goroot, "misc", "wasm", "wasm_exec.js"), // Go < 1.24
];
const wasmExecSrc = wasmExecCandidates.find((p) => fs.existsSync(p));
if (!wasmExecSrc) {
  throw new Error(`could not find wasm_exec.js under GOROOT (${goroot})`);
}
fs.copyFileSync(wasmExecSrc, path.join(outDir, "wasm_exec.js"));

console.log(`building ${goWasmPkg} -> ${path.join(outDir, "objectivelol.wasm")}`);
execFileSync("go", ["build", "-o", path.join(outDir, "objectivelol.wasm"), "./js/go/wasm"], {
  cwd: repoRoot,
  env: { ...process.env, GOOS: "js", GOARCH: "wasm" },
  stdio: "inherit",
});

console.log("wasm build complete");
