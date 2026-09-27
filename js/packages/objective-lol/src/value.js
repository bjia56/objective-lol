"use strict";

const GO_VALUE_ID_KEY = "__GoValue_id";
const REF_SYMBOL = Symbol("objective-lol.ref");

class ObjectiveLOLError extends Error {
  constructor(type, message) {
    super(message);
    this.name = "ObjectiveLOLError";
    this.type = type;
  }
}

// Parses a js/go/bridge envelope string ({"ok":true,"value":...} or
// {"ok":false,"error":{...}}), throwing ObjectiveLOLError on failure.
function parseEnvelope(envelopeJSON) {
  let envelope;
  try {
    envelope = JSON.parse(envelopeJSON);
  } catch (err) {
    throw new ObjectiveLOLError("internal", `malformed response from VM: ${err.message}`);
  }
  if (!envelope.ok) {
    const info = envelope.error || {};
    throw new ObjectiveLOLError(info.type || "runtime", info.message || "unknown error");
  }
  return envelope.value;
}

// Wraps a {"__GoValue_id": "..."} object reference in a Proxy so that
// `ref.SOME_METHOD(a, b)` on the JS side transparently becomes a
// vm.callMethod(ref, "SOME_METHOD", [a, b]) round trip, mirroring the
// dynamic proxy classes the Python bindings build with a metaclass.
function createObjectProxy(vm, ref) {
  const handler = {
    get(_target, prop) {
      if (prop === REF_SYMBOL) return ref;
      if (typeof prop === "symbol" || prop === "then") return undefined;
      return (...args) => vm.callMethod(ref, prop, args);
    },
    has() {
      return true;
    },
  };
  return new Proxy(Object.create(null), handler);
}

function isObjectProxy(value) {
  return value !== null && typeof value === "object" && value[REF_SYMBOL] !== undefined;
}

// Converts a value decoded from a bridge envelope into the JS value
// handed back to callers: primitives pass through, arrays/objects
// recurse, and object references become proxies.
function fromWire(vm, value) {
  if (value === null || value === undefined) return null;
  if (Array.isArray(value)) return value.map((v) => fromWire(vm, v));
  if (typeof value === "object") {
    if (Object.prototype.hasOwnProperty.call(value, GO_VALUE_ID_KEY)) {
      return createObjectProxy(vm, value);
    }
    const result = {};
    for (const key of Object.keys(value)) result[key] = fromWire(vm, value[key]);
    return result;
  }
  return value;
}

// Converts a JS value supplied by a caller into something JSON-safe to
// send across the bridge: object proxies unwrap back to their
// {"__GoValue_id"} reference, everything else recurses structurally.
function toWire(value) {
  if (value === null || value === undefined) return null;
  if (isObjectProxy(value)) return value[REF_SYMBOL];
  if (Array.isArray(value)) return value.map(toWire);
  if (typeof value === "object") {
    const result = {};
    for (const key of Object.keys(value)) result[key] = toWire(value[key]);
    return result;
  }
  return value;
}

module.exports = {
  GO_VALUE_ID_KEY,
  ObjectiveLOLError,
  parseEnvelope,
  fromWire,
  toWire,
  isObjectProxy,
};
