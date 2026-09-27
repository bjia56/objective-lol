# @objective-lol/core

JavaScript/TypeScript interop for the [Objective-LOL](../../../README.md) interpreter, mirroring the design of the Python bindings under `python/` but split across two backends so the same API works in Node and in the browser:

- **Node**: a native shared library built from `js/go/native` (cgo `-buildmode=c-shared`), called via [koffi](https://koffi.dev) - no `node-gyp`/N-API compilation required on install.
- **Browser**: the same interpreter compiled to WebAssembly from `js/go/wasm` (`GOOS=js GOARCH=wasm`), run via Go's `wasm_exec.js` glue.

Both backends share one Go package, `js/go/bridge`, which exposes every VM operation as a single JSON-in/JSON-out call (see that package's doc comment). The JS side (`src/vm.js`, `src/value.js`) is entirely backend-agnostic; only `src/backend-node.js` and `src/backend-wasm.js` differ, adapting the bridge's JSON envelope calling convention to koffi's C ABI and to `syscall/js` respectively.

## Building

```bash
npm install
npm run build        # builds both backends for the current platform
# or individually:
npm run build:native
npm run build:wasm
```

`build:native` requires a matching C toolchain for the current platform (only builds for the host platform/arch - see `scripts/build-native.js`); `build:wasm` only needs a Go toolchain.

## Usage (Node)

```js
const { createVM } = require("@objective-lol/core");

const vm = await createVM(); // native backend

await vm.execute(`
  HAI ME TEH FUNCSHUN MAIN
    GIVEZ 1 MOAR 2
  KTHXBAI
`);

await vm.defineFunction("DOUBLE", 1, (x) => x * 2);
console.log(await vm.call("DOUBLE", 21)); // 42

const obj = await vm.newObjectInstance("SOME_CLASS");
await obj.SOME_METHOD(1, 2); // object references cross the boundary as
                              // Proxies; method calls route back through
                              // vm.callMethod()
```

## Usage (browser)

```js
import { createVM } from "@objective-lol/core/browser";

const vm = await createVM({
  wasmURL: new URL("@objective-lol/core/wasm/objectivelol.wasm", import.meta.url),
  wasmExecURL: new URL("@objective-lol/core/wasm/wasm_exec.js", import.meta.url),
});

await vm.execute(`...`);
```

See `examples/node-example.js` and `examples/browser-example.html` for runnable versions of both.

## API surface

`ObjectiveLOLVM` (returned by `createVM()`) exposes:

- `execute(code)` - run source, returns `MAIN`'s return value
- `call(name, ...args)` - call a top-level function
- `callMethod(ref, name, args)` - call a method on an object reference (usually called implicitly via the Proxy, e.g. `ref.NAME(...args)`)
- `defineVariable(name, value, constant)` / `setVariable(name, value)` / `getVariable(name)`
- `newObjectInstance(className)`
- `defineFunction(name, argc, fn)` - expose a **synchronous** JS function to Objective-LOL (see below)
- `free()`

Values marshal the same way as the Python bindings: numbers/strings/booleans pass through directly, arrays and plain objects become `BUKKIT`/`BASKIT`, and object instances cross the boundary as opaque references (wrapped in a JS `Proxy` on the way back in) rather than being copied.

### Why `defineFunction` callbacks must be synchronous

Both backends invoke a defined function as a direct, blocking call from deep inside the Go tree-walking interpreter - there's no continuation for either the interpreter or the FFI layer to suspend and resume later with an async result. A callback that returns a `Promise` throws instead of silently misbehaving.

## Native vs. wasm: what you lose in the browser

The wasm build is Go compiled with `GOOS=js`, which has no real OS threads (goroutines are cooperatively scheduled on the one JS thread) and no filesystem by default. Concretely: the `THREAD` stdlib module won't give you true parallelism under wasm, and `FILE`/relative-import module resolution has no disk to resolve against unless you provide your own virtual filesystem. Use the native backend under Node for full stdlib fidelity; reach for wasm only where running in a browser (or another environment without a native build) is the actual requirement.

Two spots in `pkg/stdlib` (`RWX`'s setter in `FILE`, and UDP `BIND` in `SOCKET`) call `os.Chmod`/`net.ListenPacket`, which aren't implemented for `GOOS=js` under either the standard wasm port or TinyGo. Both are isolated behind a `chmodFile`/`listenPacketConn` build-tag pair (`pkg/stdlib/chmod_{other,js}.go`, `pkg/stdlib/listenpacket_{other,js}.go`) so the modules still compile and register everywhere; only those two specific operations return an ordinary Objective-LOL exception under wasm instead of working, rather than being unavailable outright.

### Building with TinyGo instead

`scripts/build-wasm.js` uses the standard Go compiler, which produces a wasm binary on the order of 13 MB (~3 MB gzipped) because most of its bulk is the Go runtime/scheduler/GC compiled into linear memory - `-ldflags="-s -w"` barely moves that number. [TinyGo](https://tinygo.org) reimplements much of the runtime and stdlib for embedded/wasm targets and produces a substantially smaller binary for the same program - roughly 4.6 MB (~1.5 MB gzipped) for this interpreter, verified with:

```bash
tinygo build -o wasm/objectivelol.wasm -target wasm ./js/go/wasm
```

(TinyGo 0.34 requires a Go 1.19-1.23 toolchain to build *with*, not just Go itself - point `GOROOT`/`PATH` at one if your default `go` is newer.) This isn't wired into `build-wasm.js` as the default because TinyGo's stdlib coverage is narrower than upstream Go's (that's exactly why the chmod/listenPacket guards above exist) and its `syscall/js` finalizer support is a no-op (`syscall/js.finalizeRef not implemented`, printed to stderr - harmless for the request/response calling pattern used here, but worth knowing about before relying on it for anything long-running with heavy JS value churn). Swap it in if the size difference matters more than that tradeoff.
