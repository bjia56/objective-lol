# TEST Module

## Import

```lol
BTW Full import
I CAN HAS TEST?

BTW Selective import examples
I CAN HAS ASSERT FROM TEST?
```

## Assertions

### ASSERT

Asserts that a condition is truthy, throwing an exception if the condition evaluates to NO.
Accepts any type and evaluates truthiness according to Objective-LOL truthiness rules.

**Syntax:** `ASSERT WIT <condition>`
**Returns:** 

**Parameters:**
- `condition` (): Any value to test for truthiness

**Example: Basic assertion**

```lol
ASSERT WIT YEZ
SAYZ WIT "Test passed!"
```

**Example: Assert with variables**

```lol
I HAS A VARIABLE COUNT TEH NUMBR ITZ 5
ASSERT WIT COUNT
```

**Example: Assert comparison result**

```lol
ASSERT WIT 2 SAEM AS 2
```

**Note:** Truthiness: NO, 0, 0.0, "", empty arrays, NOTHIN are falsy

**Note:** All other values are truthy

**Note:** Throws "Assertion failed" when condition is falsy

