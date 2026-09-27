#!/usr/bin/env node
// Builds the browser wasm artifact(s) from js/go/wasm and stages the
// matching wasm_exec.js glue alongside them.
//
// Two variants can be built, side by side, so a consumer picks one by
// which URL their own code references (see README.md - a bundler only
// copies the asset it's actually pointed at, so the other variant never
// ends up in a deployed bundle):
//
//   node scripts/build-wasm.js            standard Go compiler
//     -> wasm/objectivelol.wasm, wasm/wasm_exec.js
//   node scripts/build-wasm.js --tinygo    TinyGo (smaller, narrower stdlib)
//     -> wasm/objectivelol-tinygo.wasm, wasm/wasm_exec-tinygo.js
//
// The standard variant only needs a Go toolchain on PATH. The TinyGo
// variant additionally needs `tinygo` on PATH, backed by a Go 1.19-1.23
// toolchain specifically (TinyGo does not yet support newer Go - point
// PATH at one if your default `go` is newer).
const { execFileSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const pkgRoot = path.resolve(__dirname, "..");
const repoRoot = path.resolve(pkgRoot, "..", "..", "..");
const goWasmPkg = path.join(repoRoot, "js", "go", "wasm");
const outDir = path.join(pkgRoot, "wasm");

fs.mkdirSync(outDir, { recursive: true });

const useTinygo = process.argv.includes("--tinygo");

if (useTinygo) {
  buildTinygo();
} else {
  buildStandard();
}

function buildStandard() {
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

  const outFile = path.join(outDir, "objectivelol.wasm");
  console.log(`building ${goWasmPkg} -> ${outFile} (standard Go)`);
  execFileSync("go", ["build", "-o", outFile, "./js/go/wasm"], {
    cwd: repoRoot,
    env: { ...process.env, GOOS: "js", GOARCH: "wasm" },
    stdio: "inherit",
  });

  console.log("wasm build complete (standard)");
}

function buildTinygo() {
  const goVersionOut = execFileSync("go", ["version"], { encoding: "utf8" });
  const match = goVersionOut.match(/go(\d+)\.(\d+)/);
  const [, major, minor] = match ? match.map(Number) : [null, 0, 0];
  if (!(major === 1 && minor >= 19 && minor <= 23)) {
    throw new Error(
      `TinyGo needs a Go 1.19-1.23 toolchain on PATH to build with (found: ${goVersionOut.trim()}). ` +
        `Put a matching \`go\` earlier on PATH and retry - this is independent of the Go version used ` +
        `for the native backend or the standard wasm build.`
    );
  }

  const tinygoRoot = execFileSync("tinygo", ["env", "TINYGOROOT"], { encoding: "utf8" }).trim();
  const wasmExecSrc = path.join(tinygoRoot, "targets", "wasm_exec.js");
  if (!fs.existsSync(wasmExecSrc)) {
    throw new Error(`could not find wasm_exec.js under TINYGOROOT (${tinygoRoot})`);
  }
  fs.copyFileSync(wasmExecSrc, path.join(outDir, "wasm_exec-tinygo.js"));

  const outFile = path.join(outDir, "objectivelol-tinygo.wasm");
  console.log(`building ${goWasmPkg} -> ${outFile} (TinyGo)`);
  execFileSync("tinygo", ["build", "-o", outFile, "-target", "wasm", "./js/go/wasm"], {
    cwd: repoRoot,
    stdio: "inherit",
  });

  console.log("wasm build complete (tinygo)");
}
