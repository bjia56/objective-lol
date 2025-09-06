# Control Flow

This guide covers conditional statements, loops, and exception handling in Objective-LOL.

## Block Scoping

All control flow structures create their own **block scope**. Variables declared inside these blocks are only accessible within that block and are automatically cleaned up when the block exits.

## Conditional Statements

### Basic IF Statement

```lol
IZ <condition>?
    BTW statements here
KTHX
```

**Example:**
```lol
I HAS A VARIABLE AGE TEH INTEGR ITZ 18

IZ AGE BIGGR THAN 17?
    I HAS A VARIABLE STATUS TEH STRIN ITZ "adult"  BTW Block-scoped variable
    SAYZ WIT "You are an adult!"
KTHX
BTW STATUS is not accessible here
```

### IF-ELSE Statement

```lol
IZ <condition>?
    BTW if block
NOPE
    BTW else block
KTHX
```

**Example:**
```lol
I HAS A VARIABLE SCORE TEH INTEGR ITZ 85

IZ SCORE BIGGR THAN 89?
    SAYZ WIT "Grade A"
NOPE
    SAYZ WIT "Grade B or lower"
KTHX
```

### MEBBE (Else-If) Statement

Use `MEBBE` to chain multiple conditions together:

```lol
IZ <condition>?
    BTW if block
MEBBE <condition>?
    BTW else-if block
MEBBE <condition>?
    BTW another else-if block
NOPE
    BTW else block
KTHX
```

**Example:**
```lol
I HAS A VARIABLE SCORE TEH INTEGR ITZ 85

IZ SCORE BIGGR THAN 89?
    SAYZ WIT "Grade A"
MEBBE SCORE BIGGR THAN 79?
    SAYZ WIT "Grade B"
MEBBE SCORE BIGGR THAN 69?
    SAYZ WIT "Grade C"
MEBBE SCORE BIGGR THAN 59?
    SAYZ WIT "Grade D"
NOPE
    SAYZ WIT "Grade F"
KTHX
```

**Advanced Example with Nested Conditions:**
```lol
I HAS A VARIABLE AGE TEH INTEGR ITZ 25
I HAS A VARIABLE HAS_LICENSE TEH BOOL ITZ YEZ

IZ AGE SMALLR THAN 16?
    SAYZ WIT "Too young to drive"
MEBBE AGE SMALLR THAN 18?
    IZ HAS_LICENSE?
        SAYZ WIT "Can drive with restrictions"
    NOPE
        SAYZ WIT "Need license first"
    KTHX
MEBBE AGE SMALLR THAN 21?
    SAYZ WIT "Can drive, no alcohol"
NOPE
    SAYZ WIT "Full driving privileges"
KTHX
```

### Complex Conditions

```lol
I HAS A VARIABLE X TEH INTEGR ITZ 5
I HAS A VARIABLE Y TEH INTEGR ITZ 10
I HAS A VARIABLE IS_VALID TEH BOOL ITZ YEZ

BTW Using boolean variables
IZ IS_VALID?
    SAYZ WIT "Valid input"
KTHX

BTW Using logical operators with parentheses
IZ (X SMALLR THAN Y) AN (Y BIGGR THAN 5)?
    SAYZ WIT "Both conditions are true"
KTHX
```

## Loops

### WHILE Loop

```lol
WHILE <condition>
    BTW statements here
KTHX
```

**Example with Block Scoping:**
```lol
I HAS A VARIABLE COUNTER TEH INTEGR ITZ 5

WHILE COUNTER BIGGR THAN 0
    I HAS A VARIABLE MSG TEH STRIN ITZ "Counting down"  BTW Block-scoped
    SAYZ WIT MSG
    SAYZ WIT COUNTER
    COUNTER ITZ COUNTER LES 1
KTHX
BTW MSG is not accessible here
```

### Nested Loops

```lol
I HAS A VARIABLE I TEH INTEGR ITZ 1
I HAS A VARIABLE J TEH INTEGR

WHILE I SMALLR THAN 4
    I HAS A VARIABLE OUTER_VAR TEH STRIN ITZ "outer"  BTW Outer block scope
    J ITZ 1

    WHILE J SMALLR THAN 3
        I HAS A VARIABLE INNER_VAR TEH STRIN ITZ "inner"  BTW Inner block scope
        SAY WIT I
        SAY WIT ", "
        SAYZ WIT J
        J ITZ J MOAR 1
    KTHX
    BTW INNER_VAR not accessible here

    I ITZ I MOAR 1
KTHX
BTW Neither OUTER_VAR nor INNER_VAR accessible here
```

## Exception Handling

Objective-LOL features comprehensive exception handling with try-catch-finally blocks.
Exceptions are thrown and caught as STRIN types.

### Basic Try-Catch

```lol
MAYB
    BTW try block - code that might throw exceptions
OOPSIE <variable_name>
    BTW catch block - handle the exception
KTHX
```

**Example:**
```lol
MAYB
    SAYZ WIT "Attempting risky operation"
    OOPS "Something went wrong!"
    SAYZ WIT "This line won't execute"
OOPSIE ERROR_MESSAGE
    SAYZ WIT "Caught exception: "
    SAYZ WIT ERROR_MESSAGE
KTHX
```

### Try-Catch-Finally

```lol
MAYB
    BTW try block
OOPSIE <variable_name>
    BTW catch block
ALWAYZ
    BTW finally block - always executes
KTHX
```

**Example with Block Scoping:**
```lol
MAYB
    I HAS A VARIABLE FILE_HANDLE TEH STRIN ITZ "file.txt"  BTW Try scope
    OOPS "File not found!"
OOPSIE FILE_ERROR
    I HAS A VARIABLE ERROR_CODE TEH INTEGR ITZ 404         BTW Catch scope
    SAYZ WIT "File error: "
    SAYZ WIT FILE_ERROR
ALWAYZ
    I HAS A VARIABLE CLEANUP_MSG TEH STRIN ITZ "All done"  BTW Finally scope
    SAYZ WIT "Cleaning up resources"
KTHX
BTW None of the block-scoped variables are accessible here
```

### Throwing Exceptions

Use `OOPS` to throw exceptions:

```lol
HAI ME TEH FUNCSHUN VALIDATE_AGE WIT AGE TEH INTEGR
    IZ AGE SMALLR THAN 0?
        OOPS "Age cannot be negative!"
    KTHX
    IZ AGE BIGGR THAN 150?
        OOPS "Age seems unrealistic!"
    KTHX
    GIVEZ AGE
KTHXBAI
```

### Built-in Exceptions

The interpreter automatically throws exceptions for common errors:

```lol
MAYB
    BTW Division by zero
    I HAS A VARIABLE RESULT TEH DUBBLE ITZ 10.0 DIVIDEZ 0.0
OOPSIE MATH_ERROR
    SAYZ WIT MATH_ERROR  BTW "Division by zero"
KTHX

MAYB
    BTW Type casting errors
    I HAS A VARIABLE NUM TEH INTEGR ITZ "not_a_number" AS INTEGR
OOPSIE CAST_ERROR
    SAYZ WIT CAST_ERROR  BTW "cannot cast string 'not_a_number' to INTEGR"
KTHX

MAYB
    BTW Array bounds errors
    I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT
    I HAS A VARIABLE ITEM TEH INTEGR ITZ ARR DO AT WIT 10
OOPSIE BOUNDS_ERROR
    SAYZ WIT BOUNDS_ERROR  BTW "Array index 10 out of bounds (size 0)"
KTHX
```

### Exception Propagation

Exceptions automatically propagate up the call stack until caught:

```lol
HAI ME TEH FUNCSHUN RISKY_OPERATION
    OOPS "Operation failed!"
KTHXBAI

HAI ME TEH FUNCSHUN CALLER
    RISKY_OPERATION  BTW Exception propagates from here
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    MAYB
        CALLER
    OOPSIE ERR
        SAYZ WIT "Caught propagated exception: "
        SAYZ WIT ERR
    KTHX
KTHXBAI
```

## Quick Reference

| Construct | Syntax |
|-----------|--------|
| If | `IZ condition? ... KTHX` |
| If-Else | `IZ condition? ... NOPE ... KTHX` |
| While | `WHILE condition ... KTHX` |
| Try-Catch | `MAYB ... OOPSIE var ... KTHX` |
| Try-Finally | `MAYB ... ALWAYZ ... KTHX` |
| Throw | `OOPS "message"` |

## Next Steps

- [Functions](functions.md) - Function declaration and scoping