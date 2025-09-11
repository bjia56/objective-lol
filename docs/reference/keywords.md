# Keywords Reference

Complete reference of all Objective-LOL keywords organized by category.

## Program Structure

| Keyword | Description | Example |
|---------|-------------|----------|
| `HAI` | Start declaration | `HAI ME TEH FUNCSHUN MAIN` |
| `ME` | Declaration syntax | `HAI ME TEH FUNCSHUN` |
| `TEH` | Type declaration | `TEH INTEGR`, `TEH FUNCSHUN` |
| `KTHXBAI` | End function/class | `HAI ME TEH FUNCSHUN ... KTHXBAI` |
| `KTHX` | End block | `IZ condition? ... KTHX` |
| `BTW` | Comment marker | `BTW This is a comment` |

## Functions

| Keyword | Description | Example |
|---------|-------------|----------|
| `FUNCSHUN` | Function keyword | `HAI ME TEH FUNCSHUN NAME` |
| `WIT` | Parameter/argument | `WIT PARAM TEH TYPE` |
| `AN` | Parameter separator | `WIT A TEH INTEGR AN WIT B TEH DUBBLE` |
| `GIVEZ` | Return statement | `GIVEZ VALUE` |
| `UP` | Void return | `GIVEZ UP` |

## Classes

| Keyword | Description | Example |
|---------|-------------|----------|
| `CLAS` | Class keyword | `HAI ME TEH CLAS PERSON` |
| `KITTEH OF` | Inheritance | `HAI ME TEH CLAS DOG KITTEH OF ANIMAL` |
| `DIS` | Member declaration | `DIS TEH VARIABLE NAME TEH STRIN` |
| `EVRYONE` | Public visibility | `EVRYONE` (default) |
| `MAHSELF` | Private visibility | `MAHSELF` |
| `SHARD` | Shared/static | `DIS TEH SHARD VARIABLE` |
| `LOCKD` | Locked/constant | `DIS TEH LOCKD VARIABLE` |

## Variables

| Keyword | Description | Example |
|---------|-------------|----------|
| `I HAS A` | Variable declaration | `I HAS A VARIABLE NAME TEH TYPE` |
| `VARIABLE` | Variable keyword | `I HAS A VARIABLE` |
| `ITZ` | Assignment/initialization | `ITZ 42` |
| `NEW` | Object instantiation | `NEW BUKKIT`, `NEW POINT WIT 1 AN WIT 2` |
| `LOCKD` | Constant declaration | `I HAS A LOCKD VARIABLE` |

## Data Types

| Keyword | Description | Range/Notes |
|---------|-------------|-------------|
| `INTEGR` | 64-bit integer | -2^63 to 2^63-1 |
| `DUBBLE` | 64-bit float | IEEE 754 double precision |
| `STRIN` | String | UTF-8 strings with escape sequences |
| `BOOL` | Boolean | `YEZ` (true) or `NO` (false) |
| `BUKKIT` | Array/list | Dynamic arrays |
| `BASKIT` | Map/dictionary | Key-value pairs (string keys) |
| `NOTHIN` | Null/void | Absence of value |

## Control Flow

| Keyword | Description | Example |
|---------|-------------|----------|
| `IZ` | If statement | `IZ CONDITION?` |
| `NOPE` | Else clause | `IZ ... NOPE ...` |
| `WHILE` | While loop | `WHILE CONDITION` |
| `DO` | Method call | `OBJECT DO METHOD` |
| `?` | Conditional marker | `IZ X BIGGR THAN 5?` |

## Exception Handling

| Keyword | Description | Example |
|---------|-------------|----------|
| `MAYB` | Try block | `MAYB ... OOPSIE ... KTHX` |
| `OOPS` | Throw exception | `OOPS "Error message"` |
| `OOPSIE` | Catch exception | `OOPSIE ERROR_VAR` |
| `ALWAYZ` | Finally block | `ALWAYZ ... KTHX` |

## Operators

### Arithmetic
| Keyword | Operation | Example |
|---------|-----------|----------|
| `MOAR` | Addition (+) | `5 MOAR 3` |
| `LES` | Subtraction (-) | `10 LES 4` |
| `TIEMZ` | Multiplication (*) | `6 TIEMZ 7` |
| `DIVIDEZ` | Division (/) | `15 DIVIDEZ 3` |

### Comparison
| Keyword | Operation | Example |
|---------|-----------|----------|
| `BIGGR THAN` | Greater than (>) | `5 BIGGR THAN 3` |
| `SMALLR THAN` | Less than (<) | `3 SMALLR THAN 5` |
| `SAEM AS` | Equal to (==) | `5 SAEM AS 5` |

### Logical
| Keyword | Operation | Example |
|---------|-----------|----------|
| `AN` | Logical AND | `YEZ AN NO` |
| `OR` | Logical OR | `YEZ OR NO` |

### Type Casting
| Keyword | Operation | Example |
|---------|-----------|----------|
| `AS` | Type conversion | `"42" AS INTEGR` |

## Boolean Values

| Keyword | Value | Description |
|---------|-------|-------------|
| `YEZ` | True | Boolean true value |
| `NO` | False | Boolean false value |

## Module System

| Keyword | Description | Example |
|---------|-------------|----------|
| `I CAN HAS` | Import statement | `I CAN HAS STDIO?` |
| `FROM` | Selective import | `I CAN HAS SAY FROM STDIO?` |
| `?` | Import terminator | `I CAN HAS MATH?` |

## Special Syntax

| Keyword | Description | Example |
|---------|-------------|----------|
| `(` `)` | Expression grouping | `(2 MOAR 3) TIEMZ 4` |
| `"` `"` | String literals | `"Hello, World!"` |
| `0x` | Hexadecimal prefix | `0xFF` (255 in decimal) |
| `0` | Octal prefix | `0777` (511 in decimal) |

## Reserved for Future Use

These keywords are reserved but not currently implemented (subject to change):

- `BREAK` - Loop breaking
- `CONTINUE` - Loop continuation
- `SWITCH` - Switch statements
- `CASE` - Switch cases
- `DEFAULT` - Switch default case

## Case Sensitivity

**Important**: All keywords are case-insensitive and automatically converted to uppercase internally. These are equivalent:

```lol
hai me teh funcshun main
HAI ME TEH FUNCSHUN MAIN
Hai Me Teh Funcshun Main
```

For consistency and readability, the documentation uses uppercase for all keywords.

## Keyword Rules

1. **No Reserved Collisions**: You cannot use keywords as variable or function names
2. **Case Insensitive**: Keywords work in any case combination
3. **Context Sensitive**: Some keywords like `AN` and `WIT` are context-sensitive
4. **Multi-word Keywords**: Some operators are multi-word (e.g., `BIGGR THAN`)