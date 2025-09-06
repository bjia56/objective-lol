# Modules and Code Organization

This guide covers the module system, file imports, and code organization in Objective-LOL.

## Module System Overview

Objective-LOL features a comprehensive module system that supports:

- **Built-in standard library modules** (STDIO, MATH, TIME, etc.)
- **Custom file modules** (import from other `.olol` files)
- **Selective imports** (import specific functions/classes)
- **Cross-platform path resolution**
- **Circular import detection**
- **Function-scoped imports**

## Import Syntax

All imports use the `I CAN HAS` syntax:

### Full Import

Import all public declarations from a module:

```lol
I CAN HAS STDIO?              BTW Built-in module (all I/O functions)
I CAN HAS "math_utils"?       BTW File module (all functions/classes/variables)
```

### Selective Import

Import only specific declarations:

```lol
I CAN HAS SAY AN SAYZ FROM STDIO?                    BTW Built-in module
I CAN HAS CALCULATE AN PI FROM "math_utils"?         BTW File module
```

## File Module Imports

### File Path Rules

- **Extension**: The `.olol` extension is automatically appended
- **Path Style**: Always use POSIX-style forward slashes `/`
- **Relative Paths**: Resolved relative to the importing file's directory
- **Cross-Platform**: Automatically converted for Windows/Linux/macOS

### Path Examples

```lol
BTW Import from same directory
I CAN HAS "math_helpers"?

BTW Import from subdirectory
I CAN HAS "utils/string_helpers"?

BTW Import from parent directory
I CAN HAS "../shared/common"?

BTW Import from nested structure
I CAN HAS "algorithms/sorting/quick_sort"?
```

### File Import Example

**math_utils.olol:**
```lol
HAI ME TEH FUNCSHUN SQUARE WIT N TEH INTEGR
    GIVEZ N TIEMZ N
KTHXBAI

HAI ME TEH FUNCSHUN CUBE WIT N TEH INTEGR
    GIVEZ N TIEMZ N TIEMZ N
KTHXBAI

HAI ME TEH VARIABLE PI TEH DUBBLE ITZ 3.14159

HAI ME TEH FUNCSHUN _PRIVATE_HELPER
    BTW This cannot be imported (starts with _)
    GIVEZ "secret"
KTHXBAI
```

**main.olol:**
```lol
I CAN HAS "math_utils"?

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE RESULT TEH INTEGR ITZ SQUARE WIT 5    BTW 25
    I HAS A VARIABLE VOLUME TEH INTEGR ITZ CUBE WIT 3      BTW 27
    I HAS A VARIABLE CIRCLE TEH DUBBLE ITZ PI TIEMZ 2      BTW 6.28318
KTHXBAI
```

### Selective File Import

```lol
BTW Import only specific declarations
I CAN HAS SQUARE AN PI FROM "math_utils"?

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE AREA TEH DUBBLE ITZ PI TIEMZ SQUARE WIT 5
    BTW CUBE function is not available (not imported)
KTHXBAI
```

## Import Scoping

### Function-Scoped Imports

**Imports are scoped to the function where they appear:**

```lol
BTW Global imports - available everywhere
I CAN HAS SAYZ FROM STDIO?
I CAN HAS TIME?

HAI ME TEH FUNCSHUN OUTER_FUNCTION
    I CAN HAS ABS AN MAX FROM MATH?    BTW Local to this function

    HAI ME TEH FUNCSHUN INNER_FUNCTION
        BTW Can access: SAYZ (global), TIME (global), ABS/MAX (parent)
        SAYZ WIT "Hello from inner function"
        I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT -42.5
    KTHXBAI

    INNER_FUNCTION
KTHXBAI

HAI ME TEH FUNCSHUN SEPARATE_FUNCTION
    BTW Can access SAYZ and TIME (global) but NOT MATH
    SAYZ WIT "This works"
    BTW I HAS A VARIABLE X TEH DUBBLE ITZ ABS WIT -5  BTW This would fail!
KTHXBAI
```

### Scoping Rules

1. **Lexical Scoping**: Functions inherit imports from their calling context
2. **No Leakage**: Imports in one function don't affect sibling functions
3. **Scope Chain**: Lookup walks up the parent scope chain

## Module Features

### Private Declarations

Functions, classes, and variables starting with underscore `_` are private:

```lol
HAI ME TEH FUNCSHUN PUBLIC_FUNCTION
    GIVEZ 42
KTHXBAI

HAI ME TEH FUNCSHUN _PRIVATE_HELPER
    BTW This cannot be imported by other modules
    GIVEZ "secret"
KTHXBAI
```

### Transitive Imports

Modules can import other modules:

**string_utils.olol:**
```lol
I CAN HAS "math_utils"?    BTW Import another module

HAI ME TEH FUNCSHUN REPEAT_STRING WIT TEXT TEH STRIN AN TIMES TEH INTEGR
    BTW Can use functions from math_utils here
    GIVEZ TEXT
KTHXBAI
```

### Circular Import Detection

The module system prevents infinite recursion:

```
Error: circular import detected during execution: circular_module_a
```

## Standard Library Modules

### Available Modules

| Module | Description | Key Items |
|--------|-------------|----------|
| `STDIO` | I/O functions | `SAY`, `SAYZ`, `GIMME` |
| `MATH` | Mathematical functions | `ABS`, `MAX`, `SQRT`, `RANDOM` |
| `TIME` | Time and date | `DATE` class, `SLEEP` function |
| `STRING` | String utilities | `LEN`, `CONCAT` |
| `TEST` | Testing functions | `ASSERT` |
| `IO` | Advanced I/O | `READER`, `WRITER`, buffered classes |
| `THREAD` | Concurrency | `YARN`, `KNOT` classes |
| `FILE` | File utilities | `DOCUMENT` class |
| `HTTP` | HTTP client | `INTERWEB`, `RESPONSE` classes |

### Standard Library Examples

```lol
BTW Full module imports
I CAN HAS STDIO?
I CAN HAS MATH?

BTW Selective imports
I CAN HAS SAY AN SAYZ FROM STDIO?
I CAN HAS ABS AN MAX FROM MATH?
I CAN HAS DATE FROM TIME?
```

## Quick Reference

| Operation | Syntax |
|-----------|--------|
| Full Import | `I CAN HAS module?` |
| Selective Import | `I CAN HAS item1 AN item2 FROM module?` |
| File Import | `I CAN HAS "path/to/module"?` |
| Built-in Import | `I CAN HAS STDIO?` |

## Next Steps

- [Standard Library Overview](../standard-library/overview.md) - Complete standard library reference
- [Examples](../examples/) - Real-world module usage examples