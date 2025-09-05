# Syntax Basics

This guide covers the fundamental syntax elements of Objective-LOL: data types, variables, and operators.

## Data Types

Objective-LOL has seven built-in data types:

### Primitive Types

| Type | Description | Example |
|------|-------------|---------|
| `INTEGR` | 64-bit signed integer | `42`, `-17`, `0xFF` |
| `DUBBLE` | 64-bit floating point | `3.14159`, `-2.5` |
| `STRIN` | String with escape sequences | `"Hello"`, `"Line 1\nLine 2"` |
| `BOOL` | Boolean values | `YEZ` (true), `NO` (false) |
| `NOTHIN` | Null/void value | `NOTHIN` |

### Collection Types

| Type | Description | Example |
|------|-------------|---------|
| `BUKKIT` | Dynamic array | `NEW BUKKIT` |
| `BASKIT` | Map/dictionary | `NEW BASKIT` |

### Basic Examples

```lol
BTW Integer types
I HAS A VARIABLE NUM TEH INTEGR ITZ 42
I HAS A VARIABLE HEX_NUM TEH INTEGR ITZ 0xFF        BTW Hexadecimal 255

BTW Floating point
I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159

BTW Strings with escape sequences
I HAS A VARIABLE NAME TEH STRIN ITZ "Alice"
I HAS A VARIABLE MESSAGE TEH STRIN ITZ "Hello \"World\"!\nNew line here"

BTW Boolean values
I HAS A VARIABLE IS_READY TEH BOOL ITZ YEZ
I HAS A VARIABLE IS_DONE TEH BOOL ITZ NO

BTW Null value
I HAS A VARIABLE EMPTY_VAR TEH STRIN ITZ NOTHIN
```

**String Escape Sequences:**
- `\"` - Double quote
- `\\` - Backslash
- `\n` - Newline
- `\t` - Tab
- `\r` - Carriage return

## Variables

### Variable Declaration

**Local Variables** (function scope):
```lol
I HAS A VARIABLE <name> TEH <type> [ITZ <value>]
```

**Global Variables**:
```lol
HAI ME TEH VARIABLE <name> TEH <type> [ITZ <value>]
```

**Constants** (immutable):
```lol
I HAS A LOCKD VARIABLE <name> TEH <type> ITZ <value>
```

### Examples

```lol
BTW Basic declaration with initialization
I HAS A VARIABLE X TEH INTEGR ITZ 42
I HAS A VARIABLE NAME TEH STRIN ITZ "Bob"

BTW Declaration without initialization (uses default values)
I HAS A VARIABLE EMPTY TEH STRIN        BTW Becomes ""
I HAS A VARIABLE ZERO TEH INTEGR        BTW Becomes 0
I HAS A VARIABLE FALSY TEH BOOL         BTW Becomes NO

BTW Constants
I HAS A LOCKD VARIABLE PI TEH DUBBLE ITZ 3.14159
I HAS A LOCKD VARIABLE MAX_SIZE TEH INTEGR ITZ 100

BTW Global variables
HAI ME TEH VARIABLE GLOBAL_COUNT TEH INTEGR ITZ 0
```

### Variable Assignment

```lol
I HAS A VARIABLE COUNT TEH INTEGR ITZ 0
COUNT ITZ 10                    BTW Assign new value
COUNT ITZ COUNT MOAR 1          BTW Increment by 1
```

## Operators

### Arithmetic Operators

| Operator | Meaning | Example | Result |
|----------|---------|---------|---------|
| `MOAR` | Addition (+) | `5 MOAR 3` | 8 |
| `LES` | Subtraction (-) | `10 LES 4` | 6 |
| `TIEMZ` | Multiplication (*) | `6 TIEMZ 7` | 42 |
| `DIVIDEZ` | Division (/) | `15 DIVIDEZ 3` | 5.0 (always DUBBLE) |

### Comparison Operators

| Operator | Meaning | Example | Result |
|----------|---------|---------|---------|
| `BIGGR THAN` | Greater than (>) | `5 BIGGR THAN 3` | YEZ |
| `SMALLR THAN` | Less than (<) | `3 SMALLR THAN 5` | YEZ |
| `SAEM AS` | Equal to (==) | `5 SAEM AS 5` | YEZ |

### Logical Operators

| Operator | Meaning | Example | Result |
|----------|---------|---------|---------|
| `AN` | Logical AND | `YEZ AN NO` | NO |
| `OR` | Logical OR | `YEZ OR NO` | YEZ |

### Type Casting

Use `AS` to explicitly convert between types:

```lol
BTW Numeric conversions
I HAS A VARIABLE INT_VAL TEH INTEGR ITZ 42
I HAS A VARIABLE DOUBLE_VAL TEH DUBBLE ITZ INT_VAL AS DUBBLE    BTW 42.0

BTW String to number
I HAS A VARIABLE NUM_STR TEH STRIN ITZ "123"
I HAS A VARIABLE NUM TEH INTEGR ITZ NUM_STR AS INTEGR          BTW 123

BTW Boolean conversions
I HAS A VARIABLE ZERO TEH INTEGR ITZ 0
I HAS A VARIABLE IS_ZERO TEH BOOL ITZ ZERO AS BOOL             BTW NO (false)
```

## Expression Grouping and Precedence

### Parentheses

Use parentheses `()` to override operator precedence:

```lol
BTW Without parentheses (multiplication first)
I HAS A VARIABLE RESULT1 TEH INTEGR ITZ 2 MOAR 3 TIEMZ 4       BTW 14

BTW With parentheses (addition first)
I HAS A VARIABLE RESULT2 TEH INTEGR ITZ (2 MOAR 3) TIEMZ 4     BTW 20

BTW Complex expressions
I HAS A VARIABLE COMPLEX TEH INTEGR ITZ ((2 MOAR 3) TIEMZ (6 LES 2))  BTW 20
```

### Operator Precedence (highest to lowest)

1. **Parentheses** - `()`
2. **Type Casting** - `AS`
3. **Multiplication/Division** - `TIEMZ`, `DIVIDEZ`
4. **Addition/Subtraction** - `MOAR`, `LES`
5. **Comparisons** - `BIGGR THAN`, `SMALLR THAN`
6. **Equality** - `SAEM AS`
7. **Logical AND** - `AN`
8. **Logical OR** - `OR`

## Automatic Type Conversions

The type system handles some conversions automatically:

```lol
BTW Integer + Double = Double
I HAS A VARIABLE MIXED TEH DUBBLE ITZ 5 MOAR 2.5      BTW 7.5

BTW Division always returns Double
I HAS A VARIABLE DIVISION TEH DUBBLE ITZ 10 DIVIDEZ 3 BTW 3.333...
```

## Truth Values

Each type has a truth value for boolean contexts:

| Type | Falsy Values | Truthy Values |
|------|-------------|---------------|
| `INTEGR` | `0` | Any non-zero |
| `DUBBLE` | `0.0` | Any non-zero |
| `STRIN` | `""` (empty) | Any non-empty |
| `BOOL` | `NO` | `YEZ` |
| `NOTHIN` | Always falsy | Never truthy |
| Collections | Never falsy | Always truthy |

## Next Steps

- [Control Flow](control-flow.md) - Learn about conditionals, loops, and exception handling
- [Collections](../standard-library/collections.md) - Working with BUKKIT and BASKIT types
- [Keywords Reference](../reference/keywords.md) - Complete keyword reference