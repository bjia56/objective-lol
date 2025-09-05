# Functions

This guide covers function declaration, calling, parameters, return values, and scoping in Objective-LOL.

## Function Declaration

Functions are declared using the `HAI ME TEH FUNCSHUN` syntax:

```lol
HAI ME TEH FUNCSHUN <name> [TEH <return_type>] [WIT <param> TEH <type> [AN WIT <param> TEH <type>]...]
    BTW function body
KTHXBAI
```

### Basic Function Examples

```lol
BTW Function with no parameters or return value
HAI ME TEH FUNCSHUN SAY_HELLO
    SAYZ WIT "Hello from function!"
KTHXBAI

BTW Function with parameters
HAI ME TEH FUNCSHUN GREET WIT NAME TEH STRIN
    SAY WIT "Hello, "
    SAYZ WIT NAME
KTHXBAI

BTW Function with return value
HAI ME TEH FUNCSHUN ADD TEH INTEGR WIT X TEH INTEGR AN WIT Y TEH INTEGR
    GIVEZ X MOAR Y
KTHXBAI
```

## Function Parameters

### Multiple Parameters

Use `AN WIT` to separate multiple parameters:

```lol
HAI ME TEH FUNCSHUN CALCULATE TEH DUBBLE WIT A TEH DUBBLE AN WIT B TEH DUBBLE AN WIT C TEH DUBBLE
    I HAS A VARIABLE RESULT TEH DUBBLE ITZ A MOAR B TIEMZ C
    GIVEZ RESULT
KTHXBAI
```

### Parameter Scope

Parameters are local to the function and follow the same scoping rules as local variables.

## Return Statements

### Returning Values

Use `GIVEZ` to return values:

```lol
HAI ME TEH FUNCSHUN GET_MAX TEH INTEGR WIT A TEH INTEGR AN WIT B TEH INTEGR
    IZ A BIGGR THAN B?
        GIVEZ A
    NOPE
        GIVEZ B
    KTHX
KTHXBAI
```

### Early Returns

```lol
HAI ME TEH FUNCSHUN CHECK_POSITIVE TEH STRIN WIT NUM TEH INTEGR
    IZ NUM SMALLR THAN 1?
        GIVEZ "Not positive"  BTW Early return
    KTHX
    GIVEZ "Positive"
KTHXBAI
```

### Void Functions

```lol
BTW Explicit void return (optional)
HAI ME TEH FUNCSHUN PRINT_INFO
    SAYZ WIT "Information printed"
    GIVEZ UP    BTW Optional explicit void return
KTHXBAI

BTW Implicit void return
HAI ME TEH FUNCSHUN ANOTHER_PRINT
    SAYZ WIT "No explicit return needed"
KTHXBAI
```

## Function Calls

### Basic Calls

```lol
HAI ME TEH FUNCSHUN MAIN
    BTW Call function with no parameters
    SAY_HELLO

    BTW Call function with parameters
    GREET WIT "Alice"

    BTW Call function and use return value
    I HAS A VARIABLE SUM TEH INTEGR ITZ ADD WIT 10 AN WIT 5
    SAYZ WIT SUM

    BTW Function call as expression
    SAYZ WIT GET_MAX WIT 15 AN WIT 23
KTHXBAI
```

## Recursive Functions

Functions can call themselves recursively:

```lol
HAI ME TEH FUNCSHUN FACTORIAL TEH INTEGR WIT N TEH INTEGR
    IZ N SMALLR THAN 2?
        GIVEZ 1
    NOPE
        GIVEZ N TIEMZ FACTORIAL WIT N LES 1
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN FIBONACCI TEH INTEGR WIT N TEH INTEGR
    IZ N SMALLR THAN 2?
        GIVEZ N
    NOPE
        GIVEZ FIBONACCI WIT N LES 1 MOAR FIBONACCI WIT N LES 2
    KTHX
KTHXBAI
```

## Function Scoping

### Lexical Scoping

Functions in Objective-LOL follow **lexical scoping** similar to other modern languages:

- Functions can access variables and functions from their **calling context**
- Lookup walks up the parent environment chain: current scope → caller scope → caller's caller → etc.
- Each function call creates a new environment with the calling environment as its parent

### Scoping Example

```lol
HAI ME TEH VARIABLE GLOBAL_VAR TEH STRIN ITZ "Available everywhere"

HAI ME TEH FUNCSHUN OUTER
    I HAS A VARIABLE LOCAL_VAR TEH STRIN ITZ "Available to inner functions"

    INNER
    BTW INNER_VAR is not accessible here
KTHXBAI

HAI ME TEH FUNCSHUN INNER
    BTW Can access both GLOBAL_VAR and LOCAL_VAR
    SAYZ WIT GLOBAL_VAR
    SAYZ WIT LOCAL_VAR

    I HAS A VARIABLE INNER_VAR TEH STRIN ITZ "Only in inner"
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    BTW Call OUTER
    OUTER
KTHXBAI
```

### Module Import Scoping

- **Module imports are function-scoped** - they only affect the function where they appear
- Functions inherit imports from their calling context through the parent chain
- No import leakage between sibling functions

```lol
I CAN HAS STDIO?  BTW Global import

HAI ME TEH FUNCSHUN OUTER_FUNCTION
    I CAN HAS MATH?  BTW Local import

    HAI ME TEH FUNCSHUN INNER_FUNCTION
        BTW Can access: STDIO (global), MATH (parent)
        SAYZ WIT "Hello from inner function"
        I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT -42.5
    KTHXBAI

    INNER_FUNCTION
KTHXBAI

HAI ME TEH FUNCSHUN SEPARATE_FUNCTION
    BTW Can access STDIO (global) but NOT MATH (not imported here)
    SAYZ WIT "This works"
    BTW I HAS A VARIABLE X TEH DUBBLE ITZ ABS WIT -5  BTW This would fail!
KTHXBAI
```

## Examples

### Simple Calculator Function

```lol
HAI ME TEH FUNCSHUN CALCULATE TEH DUBBLE WIT A TEH DUBBLE AN WIT OP TEH STRIN AN WIT B TEH DUBBLE
    IZ OP SAEM AS "+"?
        GIVEZ A MOAR B
    KTHX
    IZ OP SAEM AS "-"?
        GIVEZ A LES B
    KTHX
    IZ OP SAEM AS "*"?
        GIVEZ A TIEMZ B
    KTHX
    IZ OP SAEM AS "/"?
        IZ B SAEM AS 0.0?
            OOPS "Division by zero"
        KTHX
        GIVEZ A DIVIDEZ B
    KTHX
    OOPS "Unknown operator"
KTHXBAI
```

### Helper Function Pattern

```lol
HAI ME TEH FUNCSHUN IS_EVEN TEH BOOL WIT NUM TEH INTEGR
    GIVEZ (NUM TIEMZ 2) SAEM AS NUM
KTHXBAI

HAI ME TEH FUNCSHUN PROCESS_NUMBER WIT NUM TEH INTEGR
    IZ IS_EVEN WIT NUM?
        SAYZ WIT "Even number"
    NOPE
        SAYZ WIT "Odd number"
    KTHX
KTHXBAI
```

## Quick Reference

| Concept | Syntax |
|---------|--------|
| Function Declaration | `HAI ME TEH FUNCSHUN name ... KTHXBAI` |
| Parameters | `WIT param TEH type AN WIT param2 TEH type` |
| Return Type | `HAI ME TEH FUNCSHUN name TEH return_type` |
| Return Statement | `GIVEZ value` or `GIVEZ UP` |
| Function Call | `function_name WIT arg1 AN WIT arg2` |

## Next Steps

- [Classes](classes.md) - Object-oriented programming
- [Modules](modules.md) - Code organization and imports