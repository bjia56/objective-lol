# Operators and Precedence Reference

Complete reference for all operators, precedence rules, and expression evaluation in Objective-LOL.

## Operator Precedence

Operators are evaluated in the following order, from highest to lowest precedence:

| Precedence | Operators | Associativity | Description |
|------------|-----------|---------------|-------------|
| 1 (Highest) | `( )` | N/A | Parentheses for grouping |
| 2 | `AS` | Left-to-right | Type casting |
| 3 | `TIEMZ` `DIVIDEZ` | Left-to-right | Multiplication and division |
| 4 | `MOAR` `LES` | Left-to-right | Addition and subtraction |
| 5 | `BIGGR THAN` `SMALLR THAN` | Left-to-right | Relational comparisons |
| 6 | `SAEM AS` | Left-to-right | Equality comparison |
| 7 | `AN` | Left-to-right | Logical AND |
| 8 (Lowest) | `OR` | Left-to-right | Logical OR |

## Arithmetic Operators

### Addition - MOAR

**Syntax:** `operand1 MOAR operand2`

**Supported Types:**
- INTEGR + INTEGR → INTEGR
- DUBBLE + DUBBLE → DUBBLE
- INTEGR + DUBBLE → DUBBLE
- STRIN + STRIN → STRIN (concatenation)

```lol
I HAS A VARIABLE A TEH INTEGR ITZ 5 MOAR 3     BTW 8
I HAS A VARIABLE B TEH DUBBLE ITZ 2.5 MOAR 1.0 BTW 3.5
I HAS A VARIABLE C TEH STRIN ITZ "Hello" MOAR " World"  BTW "Hello World"
I HAS A VARIABLE D TEH DUBBLE ITZ 10 MOAR 2.5  BTW 12.5 (mixed types)
```

### Subtraction - LES

**Syntax:** `operand1 LES operand2`

**Supported Types:**
- INTEGR - INTEGR → INTEGR
- DUBBLE - DUBBLE → DUBBLE
- INTEGR - DUBBLE → DUBBLE

```lol
I HAS A VARIABLE A TEH INTEGR ITZ 10 LES 3     BTW 7
I HAS A VARIABLE B TEH DUBBLE ITZ 5.5 LES 2.0  BTW 3.5
I HAS A VARIABLE C TEH DUBBLE ITZ 15 LES 2.5   BTW 12.5 (mixed types)
```

### Multiplication - TIEMZ

**Syntax:** `operand1 TIEMZ operand2`

**Supported Types:**
- INTEGR * INTEGR → INTEGR
- DUBBLE * DUBBLE → DUBBLE
- INTEGR * DUBBLE → DUBBLE

```lol
I HAS A VARIABLE A TEH INTEGR ITZ 6 TIEMZ 7    BTW 42
I HAS A VARIABLE B TEH DUBBLE ITZ 3.0 TIEMZ 2.5 BTW 7.5
I HAS A VARIABLE C TEH DUBBLE ITZ 4 TIEMZ 2.5  BTW 10.0 (mixed types)
```

### Division - DIVIDEZ

**Syntax:** `operand1 DIVIDEZ operand2`

**Supported Types:**
- INTEGR / INTEGR → DUBBLE (always returns DUBBLE)
- DUBBLE / DUBBLE → DUBBLE
- INTEGR / DUBBLE → DUBBLE

**Exception:** Division by zero throws exception

```lol
I HAS A VARIABLE A TEH DUBBLE ITZ 15 DIVIDEZ 3   BTW 5.0
I HAS A VARIABLE B TEH DUBBLE ITZ 10.0 DIVIDEZ 4.0 BTW 2.5
I HAS A VARIABLE C TEH DUBBLE ITZ 7 DIVIDEZ 2    BTW 3.5 (integer division = double)

BTW Division by zero
MAYB
    I HAS A VARIABLE BAD TEH DUBBLE ITZ 10.0 DIVIDEZ 0.0
OOPSIE MATH_ERROR
    SAYZ WIT "Division by zero error!"
KTHX
```

## Comparison Operators

### Equality - SAEM AS

**Syntax:** `operand1 SAEM AS operand2`

**Returns:** BOOL

```lol
I HAS A VARIABLE A TEH BOOL ITZ 5 SAEM AS 5      BTW YEZ
I HAS A VARIABLE B TEH BOOL ITZ "hello" SAEM AS "hello" BTW YEZ
I HAS A VARIABLE C TEH BOOL ITZ 3.14 SAEM AS 3.14   BTW YEZ
I HAS A VARIABLE D TEH BOOL ITZ 1 SAEM AS 2      BTW NO

BTW Mixed type comparison (automatic conversion)
I HAS A VARIABLE E TEH BOOL ITZ 5 SAEM AS 5.0    BTW YEZ
I HAS A VARIABLE F TEH BOOL ITZ 0 SAEM AS NO     BTW YEZ (0 converts to false)
```

### Greater Than - BIGGR THAN

**Syntax:** `operand1 BIGGR THAN operand2`

**Returns:** BOOL

```lol
I HAS A VARIABLE A TEH BOOL ITZ 10 BIGGR THAN 5   BTW YEZ
I HAS A VARIABLE B TEH BOOL ITZ 3.5 BIGGR THAN 3.7 BTW NO
I HAS A VARIABLE C TEH BOOL ITZ 5 BIGGR THAN 5    BTW NO (not strictly greater)

BTW String comparison (lexicographic)
I HAS A VARIABLE D TEH BOOL ITZ "b" BIGGR THAN "a" BTW YEZ
I HAS A VARIABLE E TEH BOOL ITZ "apple" BIGGR THAN "application" BTW NO
```

### Less Than - SMALLR THAN

**Syntax:** `operand1 SMALLR THAN operand2`

**Returns:** BOOL

```lol
I HAS A VARIABLE A TEH BOOL ITZ 3 SMALLR THAN 8   BTW YEZ
I HAS A VARIABLE B TEH BOOL ITZ 2.7 SMALLR THAN 2.5 BTW NO
I HAS A VARIABLE C TEH BOOL ITZ 5 SMALLR THAN 5   BTW NO (not strictly less)

BTW String comparison
I HAS A VARIABLE D TEH BOOL ITZ "apple" SMALLR THAN "banana" BTW YEZ
```

## Logical Operators

### Logical AND - AN

**Syntax:** `operand1 AN operand2`

**Returns:** BOOL

**Short-circuit evaluation:** If first operand is falsy, second is not evaluated.

```lol
I HAS A VARIABLE A TEH BOOL ITZ YEZ AN YEZ       BTW YEZ
I HAS A VARIABLE B TEH BOOL ITZ YEZ AN NO        BTW NO
I HAS A VARIABLE C TEH BOOL ITZ NO AN YEZ        BTW NO

BTW With truthiness conversion
I HAS A VARIABLE D TEH BOOL ITZ 5 AN "hello"     BTW YEZ (both truthy)
I HAS A VARIABLE E TEH BOOL ITZ 0 AN YEZ         BTW NO (0 is falsy)
I HAS A VARIABLE F TEH BOOL ITZ "" AN "text"     BTW NO (empty string is falsy)
```

### Logical OR - OR

**Syntax:** `operand1 OR operand2`

**Returns:** BOOL

**Short-circuit evaluation:** If first operand is truthy, second is not evaluated.

```lol
I HAS A VARIABLE A TEH BOOL ITZ YEZ OR NO        BTW YEZ
I HAS A VARIABLE B TEH BOOL ITZ NO OR YEZ        BTW YEZ
I HAS A VARIABLE C TEH BOOL ITZ NO OR NO         BTW NO

BTW With truthiness conversion
I HAS A VARIABLE D TEH BOOL ITZ 0 OR "hello"     BTW YEZ (second operand is truthy)
I HAS A VARIABLE E TEH BOOL ITZ "" OR 0          BTW NO (both falsy)
I HAS A VARIABLE F TEH BOOL ITZ 42 OR ""         BTW YEZ (first operand is truthy)
```

## Type Casting Operator

### Type Casting - AS

**Syntax:** `operand AS target_type`

**Supported Conversions:**

| From → To | Behavior | Example |
|-----------|----------|---------|
| INTEGR → DUBBLE | Direct conversion | `42 AS DUBBLE` → 42.0 |
| DUBBLE → INTEGR | Truncation | `3.7 AS INTEGR` → 3 |
| INTEGR → STRIN | String representation | `42 AS STRIN` → "42" |
| DUBBLE → STRIN | String representation | `3.14 AS STRIN` → "3.14" |
| STRIN → INTEGR | Parse integer | `"123" AS INTEGR` → 123 |
| STRIN → DUBBLE | Parse double | `"3.14" AS DUBBLE` → 3.14 |
| BOOL → STRIN | "YEZ" or "NO" | `YEZ AS STRIN` → "YEZ" |
| Any → BOOL | Truthiness | `0 AS BOOL` → NO |

```lol
BTW Numeric conversions
I HAS A VARIABLE A TEH DUBBLE ITZ 42 AS DUBBLE   BTW 42.0
I HAS A VARIABLE B TEH INTEGR ITZ 3.7 AS INTEGR  BTW 3

BTW String conversions
I HAS A VARIABLE C TEH STRIN ITZ 123 AS STRIN    BTW "123"
I HAS A VARIABLE D TEH INTEGR ITZ "456" AS INTEGR BTW 456

BTW Boolean conversions
I HAS A VARIABLE E TEH BOOL ITZ 5 AS BOOL        BTW YEZ (non-zero = true)
I HAS A VARIABLE F TEH BOOL ITZ "" AS BOOL       BTW NO (empty string = false)

BTW Invalid conversions throw exceptions
MAYB
    I HAS A VARIABLE BAD TEH INTEGR ITZ "not_a_number" AS INTEGR
OOPSIE CAST_ERROR
    SAYZ WIT "Type casting error!"
KTHX
```

## Grouping with Parentheses

Parentheses `()` override operator precedence:

```lol
BTW Without parentheses - multiplication first
I HAS A VARIABLE A TEH INTEGR ITZ 2 MOAR 3 TIEMZ 4  BTW 14 (2 + 12)

BTW With parentheses - addition first
I HAS A VARIABLE B TEH INTEGR ITZ (2 MOAR 3) TIEMZ 4 BTW 20 (5 * 4)

BTW Nested parentheses
I HAS A VARIABLE C TEH INTEGR ITZ ((2 MOAR 3) TIEMZ (4 LES 1)) BTW 15 (5 * 3)

BTW Complex expression with multiple operators
I HAS A VARIABLE D TEH BOOL ITZ (5 BIGGR THAN 3) AN ((2 MOAR 2) SAEM AS 4) BTW YEZ
```

## Truthiness Rules

Used by logical operators and conditional statements:

| Type | Falsy Values | Truthy Values |
|------|-------------|---------------|
| `BOOL` | `NO` | `YEZ` |
| `INTEGR` | `0` | Any non-zero |
| `DUBBLE` | `0.0` | Any non-zero |
| `STRIN` | `""` (empty) | Any non-empty |
| `NOTHIN` | Always falsy | Never truthy |
| Objects | Never falsy | Always truthy |

```lol
BTW Testing truthiness
I HAS A VARIABLE A TEH BOOL ITZ NO AS BOOL        BTW NO
I HAS A VARIABLE B TEH BOOL ITZ 0 AS BOOL         BTW NO
I HAS A VARIABLE C TEH BOOL ITZ "" AS BOOL        BTW NO
I HAS A VARIABLE D TEH BOOL ITZ NOTHIN AS BOOL    BTW NO

I HAS A VARIABLE E TEH BOOL ITZ YEZ AS BOOL       BTW YEZ
I HAS A VARIABLE F TEH BOOL ITZ 42 AS BOOL        BTW YEZ
I HAS A VARIABLE G TEH BOOL ITZ "hello" AS BOOL   BTW YEZ
```

## Short-Circuit Evaluation

### AND Operator

If the first operand is falsy, the second operand is not evaluated:

```lol
BTW This function will not be called
HAI ME TEH FUNCSHUN SIDE_EFFECT TEH BOOL
    SAYZ WIT "This should not print in AND example"
    GIVEZ YEZ
KTHXBAI

BTW False AND anything = False (short-circuit)
I HAS A VARIABLE RESULT TEH BOOL ITZ NO AN SIDE_EFFECT
BTW SIDE_EFFECT was not called
SAYZ WIT RESULT  BTW NO
```

### OR Operator

If the first operand is truthy, the second operand is not evaluated:

```lol
BTW This function will not be called
HAI ME TEH FUNCSHUN ANOTHER_SIDE_EFFECT TEH BOOL
    SAYZ WIT "This should not print in OR example"
    GIVEZ NO
KTHXBAI

BTW True OR anything = True (short-circuit)
I HAS A VARIABLE RESULT TEH BOOL ITZ YEZ OR ANOTHER_SIDE_EFFECT
BTW ANOTHER_SIDE_EFFECT was not called
SAYZ WIT RESULT  BTW YEZ
```

## Complex Expression ExamplesV

### Chained Comparisons

```lol
BTW Check if value is in range [1, 10]
I HAS A VARIABLE VALUE TEH INTEGR ITZ 5
I HAS A VARIABLE IN_RANGE TEH BOOL ITZ VALUE BIGGR THAN 0 AN VALUE SMALLR THAN 11
SAYZ WIT IN_RANGE  BTW YEZ
```

### Conditional Assignment Pattern

```lol
BTW Set default value if input is invalid
I HAS A VARIABLE INPUT TEH INTEGR ITZ 0
I HAS A VARIABLE VALUE TEH INTEGR ITZ INPUT BIGGR THAN 0 AN INPUT SMALLR THAN 101 AS BOOL AS INTEGR TIEMZ INPUT MOAR (INPUT SMALLR THAN 1 AS INTEGR TIEMZ 50)
BTW Complex ternary-like behavior using arithmetic and casting
```

### Mathematical Expressions

```lol
BTW Quadratic formula: (-b + sqrt(b²-4ac)) / 2a
I CAN HAS MATH?
I HAS A VARIABLE A TEH DUBBLE ITZ 1.0
I HAS A VARIABLE B TEH DUBBLE ITZ -5.0
I HAS A VARIABLE C TEH DUBBLE ITZ 6.0

I HAS A VARIABLE DISCRIMINANT TEH DUBBLE ITZ B TIEMZ B LES 4.0 TIEMZ A TIEMZ C
I HAS A VARIABLE ROOT TEH DUBBLE ITZ (B TIEMZ -1.0 MOAR SQRT WIT DISCRIMINANT) DIVIDEZ (2.0 TIEMZ A)
```

## Related

- [Syntax Basics](../language-guide/syntax-basics.md) - Basic operator usage
- [Control Flow](../language-guide/control-flow.md) - Using operators in conditions
- [Functions](../language-guide/functions.md) - Operators in function expressions