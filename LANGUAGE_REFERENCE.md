# Objective-LOL Language Reference

**Version:** 1.0
**Language Extension:** `.olol`

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Program Structure](#program-structure)
4. [Data Types](#data-types)
5. [Variables](#variables)
6. [Operators](#operators)
7. [Control Flow](#control-flow)
8. [Functions](#functions)
9. [Object-Oriented Programming](#object-oriented-programming)
10. [Standard Library](#standard-library)
11. [Type System](#type-system)
12. [Examples](#examples)
13. [Complete Language Reference](#complete-language-reference)
14. [Error Handling](#error-handling)

---

## Introduction

Objective-LOL is a programming language inspired by LOLCODE, implemented in Go. It features:

- Strong type system with automatic casting
- Object-oriented programming with classes and inheritance
- Functions with parameters and return values
- Standard library for I/O, mathematics, and time operations
- Control flow structures (conditionals and loops)

All Objective-LOL source files must use the `.olol` file extension.

---

## Getting Started

### Installation and Building

```bash
# Clone the repository
git clone https://github.com/bjia56/objective-lol.git
cd objective-lol

# Build the interpreter
go build -o olol cmd/olol/main.go

# Run a program
./olol program.olol
```

### Your First Program

Create a file named `hello.olol`:

```lol
BTW This is a comment - Hello World program

HAI ME TEH FUNCSHUN MAIN
    VISIBLEZ WIT "Hello, World!"
KTHXBAI
```

Run it:
```bash
./olol hello.olol
```

---

## Program Structure

### Basic Structure

Every Objective-LOL program must have a `MAIN` function as the entry point:

```lol
BTW Program entry point
HAI ME TEH FUNCSHUN MAIN
    BTW Your code goes here
KTHXBAI
```

### Comments

Comments start with `BTW` and continue to the end of the line:

```lol
BTW This is a single-line comment
VISIBLEZ WIT "Hello"  BTW This is also a comment
```

### Case Sensitivity

Keywords and variables are **case-insensitive** and are automatically converted to uppercase. These are equivalent:

```lol
hai me teh funcshun main    BTW Same as below
HAI ME TEH FUNCSHUN MAIN    BTW Keywords converted to uppercase
```

---

## Data Types

Objective-LOL has five built-in data types:

### INTEGR (Integer)

64-bit signed integers. Supports decimal and hexadecimal notation:

```lol
I HAS A VARIABLE NUM1 TEH INTEGR ITZ 42
I HAS A VARIABLE NUM2 TEH INTEGR ITZ -17
I HAS A VARIABLE HEX_NUM TEH INTEGR ITZ 0xFF    BTW Hexadecimal 255
I HAS A VARIABLE NEG_HEX TEH INTEGR ITZ -0x10   BTW Negative hex -16
```

### DUBBLE (Double-precision floating point)

64-bit floating point numbers:

```lol
I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159
I HAS A VARIABLE NEGATIVE TEH DUBBLE ITZ -2.5
I HAS A VARIABLE ZERO TEH DUBBLE ITZ 0.0
```

### STRIN (String)

String literals enclosed in double quotes. Supports escape sequences:

```lol
I HAS A VARIABLE NAME TEH STRIN ITZ "Alice"
I HAS A VARIABLE MESSAGE TEH STRIN ITZ "Hello \"World\"!"
I HAS A VARIABLE NEWLINE TEH STRIN ITZ "Line 1\nLine 2"
I HAS A VARIABLE TAB TEH STRIN ITZ "Column 1\tColumn 2"
```

**Escape Sequences:**
- `\"` - Double quote
- `\\` - Backslash
- `\n` - Newline
- `\t` - Tab
- `\r` - Carriage return

### BOOL (Boolean)

Boolean values with special keywords:

```lol
I HAS A VARIABLE IS_READY TEH BOOL ITZ YEZ    BTW true
I HAS A VARIABLE IS_DONE TEH BOOL ITZ NO      BTW false
```

### NOTHIN (Null/Void)

Represents absence of value:

```lol
I HAS A VARIABLE EMPTY_VAR TEH STRIN ITZ NOTHIN
VISIBLEZ WIT NOTHIN  BTW Prints "NOTHIN"
```

---

## Variables

### Variable Declaration

Use `I HAS A VARIABLE` syntax for local (function-scope) variable declaration:

```lol
BTW Basic declaration with initialization
I HAS A VARIABLE X TEH INTEGR ITZ 42
I HAS A VARIABLE NAME TEH STRIN ITZ "Bob"
I HAS A VARIABLE READY TEH BOOL ITZ YEZ

BTW Declaration without initialization (defaults to type's zero value)
I HAS A VARIABLE EMPTY TEH STRIN        BTW Becomes empty string ""
I HAS A VARIABLE ZERO TEH INTEGR        BTW Becomes 0
I HAS A VARIABLE FALSY TEH BOOL         BTW Becomes NO
```

For global variables, use `HAI ME TEH VARIABLE`:

```lol
HAI ME TEH VARIABLE GLOBAL_VAR TEH INTEGR
```

### Variable Assignment

Assign new values using `ITZ`:

```lol
I HAS A VARIABLE COUNT TEH INTEGR ITZ 0
COUNT ITZ 10        BTW Assign new value
COUNT ITZ COUNT MOAR 1  BTW Increment by 1
```

### Locked Variables (Constants)

Use `LOCKD` to create immutable variables:

```lol
I HAS A LOCKD VARIABLE PI TEH DUBBLE ITZ 3.14159
I HAS A LOCKD VARIABLE MAX_SIZE TEH INTEGR ITZ 100

BTW PI ITZ 3.14  BTW This would cause a runtime error!
```

---

## Operators

### Arithmetic Operators

| Operator | Meaning | Example |
|----------|---------|---------|
| `MOAR` | Addition (+) | `5 MOAR 3` → 8 |
| `LES` | Subtraction (-) | `10 LES 4` → 6 |
| `TIEMZ` | Multiplication (*) | `6 TIEMZ 7` → 42 |
| `DIVIDEZ` | Division (/) | `15 DIVIDEZ 3` → 5.0 |

```lol
I HAS A VARIABLE A TEH INTEGR ITZ 10
I HAS A VARIABLE B TEH INTEGR ITZ 3

VISIBLEZ WIT A MOAR B      BTW 13
VISIBLEZ WIT A LES B       BTW 7
VISIBLEZ WIT A TIEMZ B     BTW 30
VISIBLEZ WIT A DIVIDEZ B   BTW 3.333... (always returns DUBBLE)
```

### Comparison Operators

| Operator | Meaning | Example |
|----------|---------|---------|
| `BIGGR THAN` | Greater than (>) | `5 BIGGR THAN 3` → YEZ |
| `SMALLR THAN` | Less than (<) | `3 SMALLR THAN 5` → YEZ |
| `SAEM AS` | Equal to (==) | `5 SAEM AS 5` → YEZ |

```lol
I HAS A VARIABLE X TEH INTEGR ITZ 10
I HAS A VARIABLE Y TEH INTEGR ITZ 5

IZ X BIGGR THAN Y?         BTW YEZ
    VISIBLEZ WIT "X is bigger"
KTHX

IZ X SAEM AS Y?            BTW NO
    VISIBLEZ WIT "Equal"
NOPE
    VISIBLEZ WIT "Not equal"
KTHX
```

### Logical Operators

| Operator | Meaning | Example |
|----------|---------|---------|
| `AN` | Logical AND | `YEZ AN NO` → NO |
| `OR` | Logical OR | `YEZ OR NO` → YEZ |

```lol
I HAS A VARIABLE X TEH BOOL ITZ YEZ
I HAS A VARIABLE Y TEH BOOL ITZ NO

I HAS A VARIABLE RESULT1 TEH BOOL ITZ X AN Y  BTW NO
I HAS A VARIABLE RESULT2 TEH BOOL ITZ X OR Y  BTW YEZ
```

### Expression Grouping (Parentheses)

Parentheses `()` can be used to override operator precedence and group sub-expressions:

| Usage | Example | Description |
|-------|---------|-------------|
| `(expression)` | `(5 MOAR 3) TIEMZ 2` | Groups sub-expressions |
| Nested parentheses | `((2 MOAR 3) LES 1)` | Multiple levels of grouping |

#### Operator Precedence (highest to lowest)

1. **Parentheses** - `()`
2. **Type Casting** - `AS`
3. **Multiplication/Division** - `TIEMZ`, `DIVIDEZ`
4. **Addition/Subtraction** - `MOAR`, `LES`
5. **Comparisons** - `BIGGR THAN`, `SMALLR THAN`
6. **Equality** - `SAEM AS`
7. **Logical AND** - `AN`
8. **Logical OR** - `OR`

```lol
BTW Without parentheses - follows precedence rules
I HAS A VARIABLE RESULT1 TEH INTEGR ITZ 2 MOAR 3 TIEMZ 4  BTW 14 (3*4 + 2)

BTW With parentheses - override precedence
I HAS A VARIABLE RESULT2 TEH INTEGR ITZ (2 MOAR 3) TIEMZ 4  BTW 20 (5 * 4)

BTW Complex expressions with nested parentheses
I HAS A VARIABLE RESULT3 TEH INTEGR ITZ ((2 MOAR 3) TIEMZ (6 LES 2))  BTW 20 (5 * 4)

BTW Parentheses with logical operators
I HAS A VARIABLE A TEH INTEGR ITZ 5
I HAS A VARIABLE B TEH INTEGR ITZ 3
I HAS A VARIABLE C TEH INTEGR ITZ 10
I HAS A VARIABLE CONDITION TEH BOOL ITZ (A BIGGR THAN B) AN (C BIGGR THAN A)  BTW YEZ
```

---

## Control Flow

### Conditional Statements (IZ/NOPE)

#### Simple IF Statement

```lol
I HAS A VARIABLE AGE TEH INTEGR ITZ 18

IZ AGE BIGGR THAN 17?
    VISIBLEZ WIT "You are an adult!"
KTHX
```

#### IF-ELSE Statement

```lol
I HAS A VARIABLE SCORE TEH INTEGR ITZ 85

IZ SCORE BIGGR THAN 89?
    VISIBLEZ WIT "Grade A"
NOPE
    VISIBLEZ WIT "Grade B or lower"
KTHX
```

#### Complex Conditions

```lol
I HAS A VARIABLE X TEH INTEGR ITZ 5
I HAS A VARIABLE Y TEH INTEGR ITZ 10
I HAS A VARIABLE IS_VALID TEH BOOL ITZ YEZ

BTW Using boolean variables in conditions
IZ IS_VALID?
    VISIBLEZ WIT "Valid input"
KTHX

BTW Using logical operators
I HAS A VARIABLE TEST1 TEH BOOL ITZ X SMALLR THAN Y
I HAS A VARIABLE TEST2 TEH BOOL ITZ Y BIGGR THAN 5
IZ TEST1 AN TEST2?
    VISIBLEZ WIT "Both conditions are true"
KTHX
```

### Loops (WHILE)

```lol
BTW Countdown loop
I HAS A VARIABLE COUNTER TEH INTEGR ITZ 5
WHILE COUNTER BIGGR THAN 0
    VISIBLEZ WIT COUNTER
    COUNTER ITZ COUNTER LES 1
KTHX

BTW Nested loops
I HAS A VARIABLE I TEH INTEGR ITZ 1
I HAS A VARIABLE J TEH INTEGR
WHILE I SMALLR THAN 4
    J ITZ 1
    WHILE J SMALLR THAN 3
        VISIBLE WIT I
        VISIBLE WIT ", "
        VISIBLEZ WIT J
        J ITZ J MOAR 1
    KTHX
    I ITZ I MOAR 1
KTHX
```

---

## Functions

### Function Declaration

Functions are declared using `HAI ME TEH FUNCSHUN`:

```lol
BTW Function with no parameters or return value
HAI ME TEH FUNCSHUN SAY_HELLO
    VISIBLEZ WIT "Hello from function!"
KTHXBAI

BTW Function with parameters
HAI ME TEH FUNCSHUN GREET WIT NAME TEH STRIN
    VISIBLE WIT "Hello, "
    VISIBLEZ WIT NAME
KTHXBAI

BTW Function with return value
HAI ME TEH FUNCSHUN ADD TEH INTEGR WIT X TEH INTEGR AN WIT Y TEH INTEGR
    GIVEZ X MOAR Y
KTHXBAI
```

### Function Parameters

Multiple parameters use `AN WIT`:

```lol
HAI ME TEH FUNCSHUN CALCULATE TEH DUBBLE WIT A TEH DUBBLE AN WIT B TEH DUBBLE AN WIT C TEH DUBBLE
    I HAS A VARIABLE RESULT TEH DUBBLE ITZ A MOAR B TIEMZ C
    GIVEZ RESULT
KTHXBAI
```

### Return Statements

Use `GIVEZ` to return values:

```lol
HAI ME TEH FUNCSHUN GET_MAX TEH INTEGR WIT A TEH INTEGR AN WIT B TEH INTEGR
    IZ A BIGGR THAN B?
        GIVEZ A
    NOPE
        GIVEZ B
    KTHX
KTHXBAI

BTW Early return
HAI ME TEH FUNCSHUN CHECK_POSITIVE TEH STRIN WIT NUM TEH INTEGR
    IZ NUM SMALLR THAN 1?
        GIVEZ "Not positive"  BTW Early return
    KTHX
    GIVEZ "Positive"
KTHXBAI

BTW Return nothing (void function)
HAI ME TEH FUNCSHUN PRINT_INFO
    VISIBLEZ WIT "Information printed"
    GIVEZ UP    BTW Explicit void return (optional)
KTHXBAI
```

### Function Calls

```lol
HAI ME TEH FUNCSHUN MAIN
    BTW Call function with no parameters
    SAY_HELLO

    BTW Call function with parameters
    GREET WIT "Alice"

    BTW Call function and use return value
    I HAS A VARIABLE SUM TEH INTEGR ITZ ADD WIT 10 AN WIT 5
    VISIBLEZ WIT SUM

    BTW Function call as expression
    VISIBLEZ WIT GET_MAX WIT 15 AN WIT 23
KTHXBAI
```

### Recursive Functions

```lol
HAI ME TEH FUNCSHUN FACTORIAL TEH INTEGR WIT N TEH INTEGR
    IZ N SMALLR THAN 2?
        GIVEZ 1
    NOPE
        GIVEZ N TIEMZ FACTORIAL WIT N LES 1
    KTHX
KTHXBAI
```

### Function Scoping and Environment

Functions in Objective-LOL follow **lexical scoping** similar to bash functions:

#### Variable and Function Lookup

- Functions can access variables and functions from their **calling context**
- Lookup walks up the parent environment chain: current scope → caller scope → caller's caller → etc.
- Each function call creates a new environment with the calling environment as its parent

#### Module Import Scoping

- **Module imports are function-scoped** - they only affect the function where they appear
- Functions inherit imports from their calling context through the parent chain
- No import leakage between sibling functions

```lol
I HAS A VARIABLE GLOBAL_VAR TEH STRIN ITZ "Available everywhere"
I CAN HAS STDIO?  BTW Global import

HAI ME TEH FUNCSHUN OUTER
    I HAS A VARIABLE LOCAL_VAR TEH STRIN ITZ "Available to inner functions"
    I CAN HAS MATH?  BTW Local import
    
    HAI ME TEH FUNCSHUN INNER
        BTW Can access: GLOBAL_VAR, LOCAL_VAR, STDIO, MATH
        VISIBLEZ WIT GLOBAL_VAR
        VISIBLEZ WIT LOCAL_VAR
        I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT -42
    KTHXBAI
    
    INNER
KTHXBAI
```

This scoping behavior enables powerful composition patterns while maintaining clear import boundaries.

---

## Object-Oriented Programming

### Class Declaration

Classes are declared using `HAI ME TEH CLAS`:

```lol
BTW Simple class with member variables and methods
HAI ME TEH CLAS PERSON
    EVRYONE    BTW Public visibility (default)
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
    DIS TEH VARIABLE AGE TEH INTEGR ITZ 0

    DIS TEH FUNCSHUN GET_NAME TEH STRIN
        GIVEZ NAME
    KTHX

    DIS TEH FUNCSHUN SET_NAME WIT NEW_NAME TEH STRIN
        NAME ITZ NEW_NAME
    KTHX

    DIS TEH FUNCSHUN INTRODUCE
        VISIBLE WIT "Hi, I'm "
        VISIBLE WIT NAME
        VISIBLE WIT " and I'm "
        VISIBLE WIT AGE
        VISIBLEZ WIT " years old."
    KTHX
KTHXBAI
```

### Visibility Modifiers

- `EVRYONE` - Public (accessible from outside the class)
- `MAHSELF` - Private (only accessible within the class)

```lol
HAI ME TEH CLAS BANK_ACCOUNT
    EVRYONE
    DIS TEH VARIABLE OWNER TEH STRIN ITZ "Anonymous"

    MAHSELF
    DIS TEH VARIABLE BALANCE TEH DUBBLE ITZ 0.0

    EVRYONE
    DIS TEH FUNCSHUN DEPOSIT WIT AMOUNT TEH DUBBLE
        BALANCE ITZ BALANCE MOAR AMOUNT  BTW Can access private member
    KTHX

    DIS TEH FUNCSHUN GET_BALANCE TEH DUBBLE
        GIVEZ BALANCE
    KTHX
KTHXBAI
```

### Object Creation and Usage

```lol
HAI ME TEH FUNCSHUN MAIN
    BTW Create new object
    I HAS A VARIABLE PERSON1 TEH PERSON ITZ NEW PERSON

    BTW Access member variables directly
    PERSON1 NAME ITZ "Alice"
    PERSON1 AGE ITZ 25

    BTW Access member variables
    VISIBLEZ WIT PERSON1 NAME
    VISIBLEZ WIT PERSON1 AGE

    BTW Call methods with DO
    PERSON1 DO INTRODUCE
    PERSON1 DO SET_NAME WIT "Bob"

    BTW Call methods with return values
    I HAS A VARIABLE CURRENT_NAME TEH STRIN ITZ PERSON1 DO GET_NAME
    VISIBLEZ WIT CURRENT_NAME
KTHXBAI
```

### Class Inheritance

Use `KITTEH OF` for inheritance:

```lol
BTW Base class
HAI ME TEH CLAS ANIMAL
    EVRYONE
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
    DIS TEH VARIABLE SPECIES TEH STRIN ITZ "Unknown"

    DIS TEH FUNCSHUN MAKE_SOUND
        VISIBLEZ WIT "Some generic animal sound"
    KTHX
KTHXBAI

BTW Derived class
HAI ME TEH CLAS DOG KITTEH OF ANIMAL
    EVRYONE
    DIS TEH VARIABLE BREED TEH STRIN ITZ "Mixed"

    BTW Override parent method
    DIS TEH FUNCSHUN MAKE_SOUND
        VISIBLEZ WIT "Woof!"
    KTHX

    DIS TEH FUNCSHUN WAG_TAIL
        VISIBLEZ WIT "Wagging tail happily!"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE MY_DOG TEH DOG ITZ NEW DOG
    MY_DOG NAME ITZ "Buddy"
    MY_DOG SPECIES ITZ "Canine"    BTW Inherited from ANIMAL
    MY_DOG BREED ITZ "Golden Retriever"

    MY_DOG DO MAKE_SOUND    BTW Calls overridden method: "Woof!"
    MY_DOG DO WAG_TAIL      BTW Calls dog-specific method
KTHXBAI
```

### Member Access

Access object member variables in expressions directly and call methods with `DO`:

```lol
BTW Alternative syntax for member access
I HAS A VARIABLE PERSON_NAME TEH STRIN ITZ PERSON1 NAME
I HAS A VARIABLE PERSON_AGE TEH INTEGR ITZ PERSON1 AGE

BTW Method calls with IN
I HAS A VARIABLE GREETING TEH STRIN ITZ PERSON1 DO GET_NAME
```

---

## Standard Library

### Module Import System

Standard library functions must be explicitly imported using the `I CAN HAS <module>?` syntax before they can be used.

#### Import Syntax

```lol
I CAN HAS STDIO?    BTW Import I/O functions
I CAN HAS MATH?     BTW Import mathematical functions  
I CAN HAS TIEM?     BTW Import time functions
```

#### Function-Scoped Imports

**Imports are scoped to the function where they appear.** This means:

- Imports inside a function are only available within that function
- Functions can access imports from their calling context (parent scopes)
- Each function maintains its own import scope
- No import leakage between sibling functions

```lol
BTW Global import - available everywhere
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN OUTER_FUNCTION
    I CAN HAS MATH?  BTW Math available in OUTER_FUNCTION and its callees
    
    HAI ME TEH FUNCSHUN INNER_FUNCTION
        BTW Can access STDIO (from global) and MATH (from OUTER_FUNCTION)
        VISIBLEZ WIT "Hello from inner function"
        I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT -42.5
    KTHXBAI
    
    INNER_FUNCTION
KTHXBAI

HAI ME TEH FUNCSHUN SEPARATE_FUNCTION  
    BTW Can access STDIO (global) but NOT MATH (not imported here)
    VISIBLEZ WIT "This works"
    BTW I HAS A VARIABLE X TEH DUBBLE ITZ ABS WIT -5  BTW This would fail!
KTHXBAI
```

**Scoping Rules:**

1. **Lexical Scoping**: Functions inherit imports from their calling context
2. **No Leakage**: Imports in one function don't affect sibling functions  
3. **Parent Access**: Child functions can access parent function imports
4. **Similar to Bash**: Variable and function lookup walks up the parent environment chain

#### Available Modules

- **STDIO**: I/O functions (`VISIBLE`, `VISIBLEZ`, `GIMME`)
- **MATH**: Mathematical functions (`ABS`, `MAX`, `MIN`, `SQRT`, `POW`, `RANDOM`, `SIN`, `COS`, etc.)
- **TIEM**: Time functions (`NOW`, `YEAR`, `MONTH`, `DAY`, `HOUR`, `MINUTE`, `SECOND`, etc.)

### I/O Functions (STDIO)

#### VISIBLE - Print without newline

```lol
VISIBLE WIT "Hello "
VISIBLE WIT "World"
VISIBLE WIT "!"
BTW Output: Hello World!
```

#### VISIBLEZ - Print with newline

```lol
VISIBLEZ WIT "First line"
VISIBLEZ WIT "Second line"
BTW Output:
BTW First line
BTW Second line
```

#### GIMME - Read input from user

```lol
VISIBLEZ WIT "Enter your name: "
I HAS A VARIABLE USER_INPUT TEH STRIN ITZ GIMME
VISIBLE WIT "Hello, "
VISIBLEZ WIT USER_INPUT
```

### Math Functions

#### Basic Math

```lol
BTW ABS - Absolute value
I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ ABS WIT -5.5  BTW 5.5

BTW MAX/MIN - Maximum/minimum of two values
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ MAX WIT 10.5 AN WIT 7.2  BTW 10.5
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ MIN WIT 10.5 AN WIT 7.2  BTW 7.2

BTW SQRT - Square root
I HAS A VARIABLE RESULT4 TEH DUBBLE ITZ SQRT WIT 16.0  BTW 4.0

BTW POW - Power (base^exponent)
I HAS A VARIABLE RESULT5 TEH DUBBLE ITZ POW WIT 2.0 AN WIT 3.0  BTW 8.0
```

#### Trigonometric Functions

```lol
BTW SIN/COS - Trigonometric functions (input in radians)
I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159
I HAS A VARIABLE SINE_VAL TEH DUBBLE ITZ SIN WIT PI DIVIDEZ 2.0    BTW sin(π/2) = 1.0
I HAS A VARIABLE COSINE_VAL TEH DUBBLE ITZ COS WIT 0.0             BTW cos(0) = 1.0
```

#### Random Numbers

```lol
BTW RANDOM - Random number between 0 and 1
I HAS A VARIABLE RAND_NUM TEH DUBBLE ITZ RANDOM

BTW RANDINT - Random integer in range [min, max)
I HAS A VARIABLE DICE_ROLL TEH INTEGR ITZ RANDINT WIT 1 AN WIT 7  BTW 1-6
```

### Time Functions (TIEM)

#### Current Time

```lol
BTW NOW - Current Unix timestamp (seconds)
I HAS A VARIABLE TIMESTAMP TEH INTEGR ITZ NOW

BTW MILLIS - Current Unix timestamp (milliseconds)
I HAS A VARIABLE TIMESTAMP_MS TEH INTEGR ITZ MILLIS

BTW Date/time components
I HAS A VARIABLE CURRENT_YEAR TEH INTEGR ITZ YEAR      BTW e.g., 2024
I HAS A VARIABLE CURRENT_MONTH TEH INTEGR ITZ MONTH    BTW 1-12
I HAS A VARIABLE CURRENT_DAY TEH INTEGR ITZ DAY        BTW 1-31
I HAS A VARIABLE CURRENT_HOUR TEH INTEGR ITZ HOUR      BTW 0-23
I HAS A VARIABLE CURRENT_MIN TEH INTEGR ITZ MINUTE     BTW 0-59
I HAS A VARIABLE CURRENT_SEC TEH INTEGR ITZ SECOND     BTW 0-59
```

#### Time Formatting and Delays

```lol
BTW FORMAT_TIME - Format Unix timestamp as string
I HAS A VARIABLE FORMATTED TEH STRIN ITZ FORMAT_TIME WIT NOW AN WIT "2006-01-02 15:04:05"

BTW SLEEP - Pause execution for specified seconds
VISIBLEZ WIT "Waiting..."
SLEEP WIT 2.5  BTW Sleep for 2.5 seconds
VISIBLEZ WIT "Done waiting!"
```

---

## Type System

### Type Casting with AS

Explicit type conversion using the `AS` operator:

```lol
BTW Numeric conversions
I HAS A VARIABLE INT_VAL TEH INTEGR ITZ 42
I HAS A VARIABLE DOUBLE_VAL TEH DUBBLE ITZ INT_VAL AS DUBBLE  BTW 42.0

I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159
I HAS A VARIABLE TRUNCATED TEH INTEGR ITZ PI AS INTEGR       BTW 3

BTW String conversions
I HAS A VARIABLE NUM_STR TEH STRIN ITZ "123"
I HAS A VARIABLE NUM TEH INTEGR ITZ NUM_STR AS INTEGR        BTW 123

I HAS A VARIABLE FLOAT_STR TEH STRIN ITZ "45.67"
I HAS A VARIABLE FLOAT_NUM TEH DUBBLE ITZ FLOAT_STR AS DUBBLE  BTW 45.67

BTW Boolean conversions
I HAS A VARIABLE BOOL_STR TEH STRIN ITZ "YEZ"
I HAS A VARIABLE BOOL_VAL TEH BOOL ITZ BOOL_STR AS BOOL      BTW YEZ

I HAS A VARIABLE ZERO TEH INTEGR ITZ 0
I HAS A VARIABLE IS_ZERO TEH BOOL ITZ ZERO AS BOOL          BTW NO (false)

I HAS A VARIABLE NONZERO TEH INTEGR ITZ 5
I HAS A VARIABLE IS_NONZERO TEH BOOL ITZ NONZERO AS BOOL    BTW YEZ (true)
```

### Automatic Type Conversions

The type system automatically handles compatible conversions:

```lol
BTW Integer + Double = Double
I HAS A VARIABLE INT_NUM TEH INTEGR ITZ 5
I HAS A VARIABLE DOUBLE_NUM TEH DUBBLE ITZ 2.5
I HAS A VARIABLE RESULT TEH DUBBLE ITZ INT_NUM MOAR DOUBLE_NUM  BTW 7.5

BTW Division always returns Double
I HAS A VARIABLE A TEH INTEGR ITZ 10
I HAS A VARIABLE B TEH INTEGR ITZ 3
I HAS A VARIABLE DIVISION TEH DUBBLE ITZ A DIVIDEZ B  BTW 3.333...
```

### Truth Values

Each type has a truth value for boolean contexts:

- `INTEGR`: 0 is NO (false), any other value is YEZ (true)
- `DUBBLE`: 0.0 is NO (false), any other value is YEZ (true)
- `STRIN`: Empty string "" is NO (false), any non-empty string is YEZ (true)
- `BOOL`: YEZ is true, NO is false
- `NOTHIN`: Always NO (false)

---

## Examples

### Complete Programs

#### Calculator Program

```lol
BTW Simple calculator program

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
            VISIBLEZ WIT "Error: Division by zero!"
            GIVEZ 0.0
        KTHX
        GIVEZ A DIVIDEZ B
    KTHX

    VISIBLEZ WIT "Error: Unknown operator"
    GIVEZ 0.0
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    VISIBLEZ WIT "=== Calculator ==="

    I HAS A VARIABLE NUM1 TEH DUBBLE ITZ 10.5
    I HAS A VARIABLE NUM2 TEH DUBBLE ITZ 3.2

    VISIBLE WIT NUM1
    VISIBLE WIT " + "
    VISIBLE WIT NUM2
    VISIBLE WIT " = "
    VISIBLEZ WIT CALCULATE WIT NUM1 AN WIT "+" AN WIT NUM2

    VISIBLE WIT NUM1
    VISIBLE WIT " * "
    VISIBLE WIT NUM2
    VISIBLE WIT " = "
    VISIBLEZ WIT CALCULATE WIT NUM1 AN WIT "*" AN WIT NUM2
KTHXBAI
```

#### Class-Based Game Character

```lol
BTW RPG Character system

HAI ME TEH CLAS CHARACTER
    EVRYONE
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
    DIS TEH VARIABLE LEVEL TEH INTEGR ITZ 1

    MAHSELF
    DIS TEH VARIABLE HP TEH INTEGR ITZ 100
    DIS TEH VARIABLE MAX_HP TEH INTEGR ITZ 100

    EVRYONE
    DIS TEH FUNCSHUN SET_NAME WIT NEW_NAME TEH STRIN
        NAME ITZ NEW_NAME
    KTHX

    DIS TEH FUNCSHUN GET_HP TEH INTEGR
        GIVEZ HP
    KTHX

    DIS TEH FUNCSHUN TAKE_DAMAGE WIT DAMAGE TEH INTEGR
        HP ITZ HP LES DAMAGE
        IZ HP SMALLR THAN 0?
            HP ITZ 0
        KTHX

        VISIBLE WIT NAME
        VISIBLE WIT " takes "
        VISIBLE WIT DAMAGE
        VISIBLE WIT " damage! HP: "
        VISIBLEZ WIT HP
    KTHX

    DIS TEH FUNCSHUN HEAL WIT AMOUNT TEH INTEGR
        HP ITZ HP MOAR AMOUNT
        IZ HP BIGGR THAN MAX_HP?
            HP ITZ MAX_HP
        KTHX

        VISIBLE WIT NAME
        VISIBLE WIT " heals "
        VISIBLE WIT AMOUNT
        VISIBLE WIT " HP! Current HP: "
        VISIBLEZ WIT HP
    KTHX

    DIS TEH FUNCSHUN IS_ALIVE TEH BOOL
        GIVEZ HP BIGGR THAN 0
    KTHX
KTHXBAI

HAI ME TEH CLAS WARRIOR KITTEH OF CHARACTER
    EVRYONE
    DIS TEH FUNCSHUN ATTACK WIT TARGET TEH CHARACTER
        I HAS A VARIABLE DAMAGE TEH INTEGR ITZ RANDINT WIT 15 AN WIT 26  BTW 15-25 damage
        VISIBLE WIT NAME
        VISIBLE WIT " attacks with sword for "
        VISIBLE WIT DAMAGE
        VISIBLEZ WIT " damage!"
        TARGET DO TAKE_DAMAGE WIT DAMAGE
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    VISIBLEZ WIT "=== Battle Simulation ==="

    I HAS A VARIABLE HERO TEH WARRIOR ITZ NEW WARRIOR
    HERO DO SET_NAME WIT "Sir Lancelot"

    I HAS A VARIABLE ENEMY TEH CHARACTER ITZ NEW CHARACTER
    ENEMY DO SET_NAME WIT "Goblin"

    VISIBLEZ WIT "Battle begins!"

    WHILE HERO DO IS_ALIVE AN ENEMY DO IS_ALIVE
        HERO DO ATTACK WIT ENEMY

        IZ ENEMY DO IS_ALIVE?
            I HAS A VARIABLE ENEMY_DAMAGE TEH INTEGR ITZ RANDINT WIT 5 AN WIT 16  BTW 5-15 damage
            VISIBLE WIT "Goblin strikes back for "
            VISIBLE WIT ENEMY_DAMAGE
            VISIBLEZ WIT " damage!"
            HERO DO TAKE_DAMAGE WIT ENEMY_DAMAGE
        KTHX

        VISIBLEZ WIT "---"
        SLEEP WIT 1.0  BTW Dramatic pause
    KTHX

    IZ HERO DO IS_ALIVE?
        VISIBLEZ WIT "Hero wins!"
    NOPE
        VISIBLEZ WIT "Hero defeated!"
    KTHX
KTHXBAI
```

---

## Complete Language Reference

### Keywords

#### Program Structure
- `HAI` - Start function/class declaration
- `ME` - Part of declaration syntax
- `TEH` - Type declaration keyword
- `KTHXBAI` - End function/class declaration
- `KTHX` - End block (if/while)
- `BTW` - Comment marker

#### Function Declaration
- `FUNCSHUN` - Function keyword
- `WIT` - Parameter declaration
- `AN` - Parameter separator
- `GIVEZ` - Return statement
- `UP` - Void return
- `NATIV` - Native function marker

#### Class Declaration
- `CLAS` - Class keyword
- `KITTEH OF` - Inheritance
- `DIS` - Member declaration
- `EVRYONE` - Public visibility
- `MAHSELF` - Private visibility
- `SHARD` - Shared/static member
- `LOCKD` - Locked/constant member

#### Variables
- `I HAS A` - Variable declaration prefix
- `VARIABLE` - Variable keyword
- `ITZ` - Assignment/initialization
- `NEW` - Object instantiation

#### Control Flow
- `IZ` - If statement
- `NOPE` - Else clause
- `WHILE` - While loop
- `DO` - Method call

#### Data Types
- `INTEGR` - Integer type
- `DUBBLE` - Double/float type
- `STRIN` - String type
- `BOOL` - Boolean type
- `NOTHIN` - Null/void type

#### Operators
- `MOAR` - Addition (+)
- `LES` - Subtraction (-)
- `TIEMZ` - Multiplication (*)
- `DIVIDEZ` - Division (/)
- `BIGGR THAN` - Greater than (>)
- `SMALLR THAN` - Less than (<)
- `SAEM AS` - Equal to (==)
- `AN` - Logical AND
- `OR` - Logical OR
- `AS` - Type casting
- `IN` - Member access

#### Boolean Values
- `YEZ` - True
- `NO` - False

#### Special
- `?` - Question mark (used in conditionals)
- `(` `)` - Parentheses (expression grouping)

### Built-in Functions

#### I/O Functions
- `VISIBLE WIT <value>` - Print value without newline
- `VISIBLEZ WIT <value>` - Print value with newline
- `GIMME` - Read line from input → STRIN

#### Math Functions
- `ABS WIT <number>` - Absolute value → DUBBLE
- `MAX WIT <a> AN WIT <b>` - Maximum of two values → DUBBLE
- `MIN WIT <a> AN WIT <b>` - Minimum of two values → DUBBLE
- `SQRT WIT <number>` - Square root → DUBBLE
- `POW WIT <base> AN WIT <exp>` - Power → DUBBLE
- `SIN WIT <angle>` - Sine (radians) → DUBBLE
- `COS WIT <angle>` - Cosine (radians) → DUBBLE
- `RANDOM` - Random number [0,1) → DUBBLE
- `RANDINT WIT <min> AN WIT <max>` - Random integer [min,max) → INTEGR

#### Time Functions
- `NOW` - Current Unix timestamp → INTEGR
- `MILLIS` - Current timestamp in milliseconds → INTEGR
- `YEAR` - Current year → INTEGR
- `MONTH` - Current month (1-12) → INTEGR
- `DAY` - Current day of month → INTEGR
- `HOUR` - Current hour (0-23) → INTEGR
- `MINUTE` - Current minute (0-59) → INTEGR
- `SECOND` - Current second (0-59) → INTEGR
- `FORMAT_TIME WIT <timestamp> AN WIT <format>` - Format timestamp → STRIN
- `SLEEP WIT <seconds>` - Sleep for duration → NOTHIN

### Syntax Patterns

#### Variable Declaration
```
I HAS A [LOCKD] VARIABLE <name> TEH <type> [ITZ <value>]
```

#### Function Declaration
```
HAI ME TEH FUNCSHUN <name> [TEH <return_type>] [WIT <param> TEH <type> [AN WIT <param> TEH <type>]...]
    <statements>
KTHXBAI
```

#### Class Declaration
```
HAI ME TEH CLAS <name> [KITTEH OF <parent>]
    [EVRYONE|MAHSELF]
    [DIS TEH VARIABLE <name> TEH <type> [ITZ <value>]]
    [DIS TEH FUNCSHUN <name> ...]
KTHXBAI
```

#### Control Flow
```
IZ <condition>?
    <statements>
[NOPE
    <statements>]
KTHX

WHILE <condition>
    <statements>
KTHX
```

---

## Error Handling

### Common Errors

#### Syntax Errors
```lol
BTW Missing KTHXBAI
HAI ME TEH FUNCSHUN TEST
    VISIBLEZ WIT "Hello"
BTW Error: Expected KTHXBAI

BTW Missing TEH in variable declaration
I HAS A VARIABLE X INTEGR ITZ 5
BTW Error: Expected TEH after VARIABLE X
```

#### Type Errors
```lol
BTW Cannot cast incompatible types
I HAS A VARIABLE STR TEH STRIN ITZ "not a number"
I HAS A VARIABLE NUM TEH INTEGR ITZ STR AS INTEGR
BTW Runtime Error: cannot cast string 'not a number' to INTEGR

BTW Division by zero
I HAS A VARIABLE RESULT TEH DUBBLE ITZ 10.0 DIVIDEZ 0.0
BTW Returns 0.0 (handled gracefully)
```

#### Runtime Errors
```lol
BTW Accessing undefined variable
VISIBLEZ WIT UNDEFINED_VAR
BTW Runtime Error: undefined variable 'UNDEFINED_VAR'

BTW Modifying locked variable
I HAS A LOCKD VARIABLE CONSTANT TEH INTEGR ITZ 42
CONSTANT ITZ 100
BTW Runtime Error: cannot modify locked variable 'CONSTANT'
```

### Debugging Tips

1. **Use VISIBLEZ for debugging**: Add temporary output statements to trace program flow
2. **Check variable types**: Use explicit casting when mixing types
3. **Verify function signatures**: Ensure parameter types match function definitions
4. **Test boundary conditions**: Check division by zero, negative square roots, etc.
5. **Use the test runner**: Run `./run_tests.sh` to verify interpreter behavior

### File Extensions

All Objective-LOL source files **must** use the `.olol` extension. The interpreter enforces this requirement.

---

*This reference guide covers all major features of Objective-LOL. For additional examples, see the test files in the `tests/` directory of the repository.*