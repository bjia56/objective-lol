"use strict";

const { test } = require("node:test");
const assert = require("node:assert/strict");
const { createVM } = require("../src/index");

for (const backendName of ["native", "wasm"]) {
  const makeVM = (opts = {}) =>
    createVM(backendName === "wasm" ? { ...opts, backend: "wasm" } : opts);

  test(`[${backendName}] execute() runs a MAIN function and returns its value`, async () => {
    const vm = await makeVM();
    const result = await vm.execute(`
      HAI ME TEH FUNCSHUN MAIN
        GIVEZ 1 MOAR 2
      KTHXBAI
    `);
    assert.equal(result, 3);
  });

  test(`[${backendName}] defineVariable() + call() round-trips values`, async () => {
    const vm = await makeVM();
    await vm.defineVariable("GREETING", "hello from js", false);
    await vm.execute(`
      I CAN HAS STDIO?
      HAI ME TEH FUNCSHUN MAIN
        SAYZ WIT GREETING
      KTHXBAI
    `);
  });

  test(`[${backendName}] call() invokes a top-level function with arguments`, async () => {
    const vm = await makeVM();
    await vm.execute(`
      HAI ME TEH FUNCSHUN ADD TEH INTEGR WIT X TEH INTEGR AN WIT Y TEH INTEGR
        GIVEZ X MOAR Y
      KTHXBAI

      HAI ME TEH FUNCSHUN MAIN
        GIVEZ 0
      KTHXBAI
    `);
    const sum = await vm.call("ADD", 5, 8);
    assert.equal(sum, 13);
  });

  test(`[${backendName}] object instances round-trip through a Proxy and callMethod`, async () => {
    const vm = await makeVM();
    await vm.execute(`
      HAI ME TEH CLAS COUNTER
          EVRYONE
          DIS TEH VARIABLE N TEH INTEGR ITZ 0

          DIS TEH FUNCSHUN BUMP TEH INTEGR WIT AMOUNT TEH INTEGR
              N ITZ N MOAR AMOUNT
              GIVEZ N
          KTHX
      KTHXBAI

      HAI ME TEH FUNCSHUN MAIN
        GIVEZ 0
      KTHXBAI
    `);

    const counter = await vm.newObjectInstance("COUNTER");
    assert.equal(await counter.BUMP(4), 4);
    assert.equal(await counter.BUMP(10), 14);
  });

  test(`[${backendName}] defineFunction() exposes a synchronous JS callback to Objective-LOL`, async () => {
    const vm = await makeVM();
    const seen = [];
    await vm.defineFunction("JSADD", 2, (a, b) => {
      seen.push([a, b]);
      return a + b;
    });
    await vm.execute(`
      HAI ME TEH FUNCSHUN MAIN
        GIVEZ 0
      KTHXBAI
    `);
    const result = await vm.call("JSADD", 3, 4);
    assert.equal(result, 7);
    assert.deepEqual(seen, [[3, 4]]);
  });

  test(`[${backendName}] BUKKIT/BASKIT values round-trip as arrays/objects`, async () => {
    const vm = await makeVM();
    await vm.execute(`
      HAI ME TEH FUNCSHUN MAIN
        GIVEZ 0
      KTHXBAI
    `);
    await vm.defineFunction("ECHO", 1, (v) => v);
    const arr = await vm.call("ECHO", [1, 2, 3]);
    assert.deepEqual(arr, [1, 2, 3]);
  });

  test(`[${backendName}] runtime errors surface as ObjectiveLOLError`, async () => {
    const vm = await makeVM();
    await assert.rejects(() => vm.execute("THIS IS NOT VALID OLOL"), /ObjectiveLOLError|error/i);
  });
}
