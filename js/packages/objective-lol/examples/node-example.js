// Run with: node examples/node-example.js
// (after `npm run build` from the package root to produce the native
// addon and, optionally, the wasm artifact)
"use strict";

const { createVM } = require("..");

async function main() {
  const vm = await createVM(); // native backend by default

  await vm.execute(`
    I CAN HAS STDIO?

    HAI ME TEH CLAS COUNTER
        EVRYONE
        DIS TEH VARIABLE N TEH INTEGR ITZ 0

        DIS TEH FUNCSHUN BUMP TEH INTEGR WIT AMOUNT TEH INTEGR
            N ITZ N MOAR AMOUNT
            GIVEZ N
        KTHX
    KTHXBAI

    HAI ME TEH FUNCSHUN MAIN
      SAYZ WIT "hello from Objective-LOL"
    KTHXBAI
  `);

  // Call a JS function from Objective-LOL and vice versa.
  await vm.defineFunction("DOUBLE", 1, (x) => x * 2);
  console.log("DOUBLE(21) =", await vm.call("DOUBLE", 21));

  // Objects cross the boundary as proxies; method calls round-trip
  // through the VM.
  const counter = await vm.newObjectInstance("COUNTER");
  console.log("BUMP(5) =", await counter.BUMP(5));
  console.log("BUMP(10) =", await counter.BUMP(10));

  await vm.free();
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
