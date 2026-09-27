#!/usr/bin/env node
// Builds the native c-shared library from js/go/native for the *current*
// host platform/arch and stages it under native/<platform>-<arch>/. This
// is meant to be run once per target platform (e.g. in CI, once per
// runner OS) rather than cross-compiled, since cgo needs a matching C
// toolchain for the target.
const { execFileSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const pkgRoot = path.resolve(__dirname, "..");
const repoRoot = path.resolve(pkgRoot, "..", "..", "..");
const goNativePkg = path.join(repoRoot, "js", "go", "native");

const platform = process.platform; // 'linux' | 'darwin' | 'win32'
const arch = process.arch; // 'x64' | 'arm64' | ...

const libNames = {
  linux: "libolol.so",
  darwin: "libolol.dylib",
  win32: "olol.dll",
};
const libName = libNames[platform];
if (!libName) {
  throw new Error(`unsupported platform for native build: ${platform}`);
}

const outDir = path.join(pkgRoot, "native", `${platform}-${arch}`);
fs.mkdirSync(outDir, { recursive: true });
const outFile = path.join(outDir, libName);

console.log(`building ${goNativePkg} -> ${outFile}`);
execFileSync(
  "go",
  ["build", "-buildmode=c-shared", "-o", outFile, "./js/go/native"],
  {
    cwd: repoRoot,
    env: { ...process.env, CGO_ENABLED: "1" },
    stdio: "inherit",
  }
);

// go build -buildmode=c-shared also emits a .h header we don't need in
// the published package.
const header = outFile.replace(/\.(so|dylib|dll)$/, ".h");
if (fs.existsSync(header)) fs.unlinkSync(header);

console.log("native build complete:", outFile);
