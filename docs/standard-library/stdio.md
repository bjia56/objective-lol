# STDIO Module - Input/Output Functions

The STDIO module provides basic input and output functionality for Objective-LOL programs.

## Import

```lol
BTW Full import
I CAN HAS STDIO?

BTW Selective import
I CAN HAS SAY FROM STDIO?
I CAN HAS SAY AN SAYZ FROM STDIO?
I CAN HAS SAY AN SAYZ AN GIMME FROM STDIO?
```

## Functions

### SAY - Print without newline

Prints a value to standard output without adding a newline character.

**Syntax:** `SAY WIT <value>`

**Parameters:**
- `value` - Any value to print (INTEGR, DUBBLE, STRIN, BOOL, etc.)

**Returns:** NOTHIN

**Examples:**
```lol
I CAN HAS SAY FROM STDIO?

SAY WIT "Hello "
SAY WIT "World"
SAY WIT "!"
BTW Output: Hello World!

BTW Print numbers
SAY WIT 42
SAY WIT " is the answer"
BTW Output: 42 is the answer

BTW Print boolean values
SAY WIT YEZ
SAY WIT " and "
SAY WIT NO
BTW Output: YEZ and NO
```

### SAYZ - Print with newline

Prints a value to standard output and adds a newline character.

**Syntax:** `SAYZ WIT <value>`

**Parameters:**
- `value` - Any value to print

**Returns:** NOTHIN

**Examples:**
```lol
I CAN HAS SAYZ FROM STDIO?

SAYZ WIT "First line"
SAYZ WIT "Second line"
SAYZ WIT 42
BTW Output:
BTW First line
BTW Second line
BTW 42

BTW Print variables
I HAS A VARIABLE NAME TEH STRIN ITZ "Alice"
I HAS A VARIABLE AGE TEH INTEGR ITZ 25
SAYZ WIT NAME
SAYZ WIT AGE
BTW Output:
BTW Alice
BTW 25
```

### GIMME - Read user input

Reads a line of input from the user (standard input).

**Syntax:** `GIMME`

**Parameters:** None

**Returns:** STRIN - The input line (without trailing newline)

**Examples:**
```lol
I CAN HAS STDIO?

SAYZ WIT "Enter your name: "
I HAS A VARIABLE USER_NAME TEH STRIN ITZ GIMME
SAY WIT "Hello, "
SAYZ WIT USER_NAME

BTW Interactive calculator
SAYZ WIT "Enter first number: "
I HAS A VARIABLE NUM1_STR TEH STRIN ITZ GIMME
I HAS A VARIABLE NUM1 TEH INTEGR ITZ NUM1_STR AS INTEGR

SAYZ WIT "Enter second number: "
I HAS A VARIABLE NUM2_STR TEH STRIN ITZ GIMME
I HAS A VARIABLE NUM2 TEH INTEGR ITZ NUM2_STR AS INTEGR

I HAS A VARIABLE SUM TEH INTEGR ITZ NUM1 MOAR NUM2
SAY WIT "Sum: "
SAYZ WIT SUM
```

## Usage Patterns

### Simple Output

```lol
I CAN HAS SAYZ FROM STDIO?

SAYZ WIT "Hello, World!"
SAYZ WIT "Welcome to Objective-LOL"
```

### Formatted Output

```lol
I CAN HAS SAY AN SAYZ FROM STDIO?

I HAS A VARIABLE NAME TEH STRIN ITZ "Bob"
I HAS A VARIABLE SCORE TEH INTEGR ITZ 95

SAY WIT "Player: "
SAY WIT NAME
SAY WIT ", Score: "
SAYZ WIT SCORE
BTW Output: Player: Bob, Score: 95
```

### Interactive Programs

```lol
I CAN HAS STDIO?

SAYZ WIT "=== Simple Quiz ==="
SAYZ WIT "What is 5 + 3?"
I HAS A VARIABLE ANSWER TEH STRIN ITZ GIMME

IZ ANSWER SAEM AS "8"?
    SAYZ WIT "Correct!"
NOPE
    SAYZ WIT "Wrong! The answer is 8."
KTHX
```

### Menu System

```lol
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN SHOW_MENU
    SAYZ WIT "=== Main Menu ==="
    SAYZ WIT "1. Option A"
    SAYZ WIT "2. Option B"
    SAYZ WIT "3. Exit"
    SAY WIT "Choose: "
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    SHOW_MENU
    I HAS A VARIABLE CHOICE TEH STRIN ITZ GIMME

    IZ CHOICE SAEM AS "1"?
        SAYZ WIT "You chose Option A"
    KTHX
    IZ CHOICE SAEM AS "2"?
        SAYZ WIT "You chose Option B"
    KTHX
    IZ CHOICE SAEM AS "3"?
        SAYZ WIT "Goodbye!"
    KTHX
KTHXBAI
```

## Type Handling

STDIO functions work with all Objective-LOL types:

```lol
I CAN HAS STDIO?

BTW Integer values
I HAS A VARIABLE NUM TEH INTEGR ITZ 42
SAYZ WIT NUM                    BTW Prints: 42

BTW Double values
I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159
SAYZ WIT PI                     BTW Prints: 3.14159

BTW String values
I HAS A VARIABLE TEXT TEH STRIN ITZ "Hello"
SAYZ WIT TEXT                   BTW Prints: Hello

BTW Boolean values
I HAS A VARIABLE FLAG TEH BOOL ITZ YEZ
SAYZ WIT FLAG                   BTW Prints: YEZ

BTW NOTHIN values
I HAS A VARIABLE EMPTY TEH STRIN ITZ NOTHIN
SAYZ WIT EMPTY                  BTW Prints: NOTHIN

BTW Collection info
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT
ARR DO PUSH WIT 1
ARR DO PUSH WIT 2
SAYZ WIT ARR SIZ                BTW Prints: 2
```

## Quick Reference

| Function | Purpose | Syntax | Returns |
|----------|---------|--------|----------|
| `SAY` | Print without newline | `SAY WIT value` | NOTHIN |
| `SAYZ` | Print with newline | `SAYZ WIT value` | NOTHIN |
| `GIMME` | Read user input | `GIMME` | STRIN |

## Related

- [String Module](string.md) - String manipulation functions
- [IO Module](io.md) - Advanced I/O classes for file operations
- [Examples](../examples/) - Real-world usage examples