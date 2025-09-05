# Objective-LOL Documentation

Welcome to the comprehensive documentation for Objective-LOL, a programming language inspired by LOLCODE with modern features like strong typing, object-oriented programming, and a module system.

## Quick Start

**New to Objective-LOL?** Start here:

1. [**Getting Started**](language-guide/getting-started.md) - Install, build, and write your first program
2. [**Syntax Basics**](language-guide/syntax-basics.md) - Learn data types, variables, and operators
3. [**Control Flow**](language-guide/control-flow.md) - Conditionals, loops, and exception handling

## Documentation Structure

### 📚 Language Guide

**Core language concepts and tutorials:**

- [**Getting Started**](language-guide/getting-started.md) - Installation, building, first program
- [**Syntax Basics**](language-guide/syntax-basics.md) - Types, variables, operators, expressions
- [**Control Flow**](language-guide/control-flow.md) - If/while statements, exception handling
- [**Functions**](language-guide/functions.md) - Function declaration, parameters, scoping
- [**Classes**](language-guide/classes.md) - Object-oriented programming, inheritance
- [**Modules**](language-guide/modules.md) - Import system, file organization

### 📖 Standard Library

**Built-in modules and functions:**

- [**Overview**](standard-library/overview.md) - Import system, scoping rules
- [**Collections**](standard-library/collections.md) - `BUKKIT` arrays and `BASKIT` maps
- [**STDIO**](standard-library/stdio.md) - Input/output functions
- [**MATH**](standard-library/math.md) - Mathematical functions
- [**RANDOM**](standard-library/random.md) - Random number generator
- [**TIME**](standard-library/time.md) - Date/time functionality
- [**STRING**](standard-library/string.md) - String utilities
- [**TEST**](standard-library/test.md) - Testing and assertions
- [**IO**](standard-library/io.md) - Advanced I/O classes (buffered readers/writers)
- [**THREAD**](standard-library/threading.md) - Concurrency
- [**FILE**](standard-library/file.md) - File system operations

### 🔍 Reference

**Quick lookup and reference materials:**

- [**Keywords**](reference/keywords.md) - Complete keyword reference
- [**Operators**](reference/operators.md) - Operator precedence and usage
- [**Exceptions**](reference/exceptions.md) - Exception handling and built-in errors

### 💡 Examples

**Real-world programming examples:**

- [**Calculator**](examples/calculator.md) - Simple calculator program
- [**File Processing**](examples/file-processing.md) - Data processing, logs, and file I/O

## Language Overview

### Key Features

- **Strong Type System**: INTEGR, DUBBLE, STRIN, BOOL, NOTHIN, BUKKIT, BASKIT
- **Object-Oriented**: Classes with inheritance, constructors, visibility modifiers
- **Functions**: First-class functions with lexical scoping
- **Module System**: File imports with selective/full import syntax
- **Exception Handling**: Try-catch-finally blocks with MAYB/OOPS/OOPSIE/ALWAYZ
- **Collections**: Dynamic arrays (BUKKIT) and maps (BASKIT) with rich methods
- **Concurrency**: Threading with YARN and synchronization with KNOT
- **Cross-Platform**: Runs on Windows, macOS, and Linux

### Syntax Highlights

```lol
BTW Hello World
HAI ME TEH FUNCSHUN MAIN
    I CAN HAS STDIO?
    SAYZ WIT "Hello, World!"
KTHXBAI
```

```lol
BTW Object-oriented programming
HAI ME TEH CLAS PERSON
    EVRYONE
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
    DIS TEH VARIABLE AGE TEH INTEGR ITZ 0

    DIS TEH FUNCSHUN PERSON WIT NAME_VAL TEH STRIN AN WIT AGE_VAL TEH INTEGR
        NAME ITZ NAME_VAL
        AGE ITZ AGE_VAL
    KTHX
KTHXBAI
```

```lol
BTW Exception handling
MAYB
    I HAS A VARIABLE RESULT TEH DUBBLE ITZ 10.0 DIVIDEZ 0.0
OOPSIE ERROR_MSG
    SAYZ WIT "Caught error: "
    SAYZ WIT ERROR_MSG
KTHX
```

## Learning Path

### Beginner Path
1. **Getting Started** → **Syntax Basics** → **Control Flow**
2. Practice with **STDIO** and **MATH** modules
3. Try the **Calculator Example**

### Intermediate Path
1. **Functions** → **Classes** → **Modules**
2. Explore **Collections** and **STRING** modules
3. Work through **RPG Character** and **Data Processing** examples

### Advanced Path
1. **Advanced I/O** and **Threading** modules
2. **Exception Handling** patterns
3. **Concurrent Programs** examples
4. Use **Reference** section for quick lookup

## Quick Reference

### Common Patterns

```lol
BTW Variable declaration
I HAS A VARIABLE NAME TEH TYPE [ITZ VALUE]

BTW Function declaration
HAI ME TEH FUNCSHUN NAME [TEH RETURN_TYPE] [WIT PARAMS]
    BTW body
KTHXBAI

BTW Class declaration
HAI ME TEH CLAS NAME [KITTEH OF PARENT]
    BTW members
KTHXBAI

BTW Import syntax
I CAN HAS MODULE?                          BTW Full import
I CAN HAS FUNC1 AN FUNC2 FROM MODULE?      BTW Selective import

BTW Exception handling
MAYB
    BTW try block
OOPSIE ERROR_VAR
    BTW catch block
ALWAYZ
    BTW finally block
KTHX
```

### Essential Imports

```lol
I CAN HAS STDIO?    BTW For SAY, SAYZ, GIMME
I CAN HAS MATH?     BTW For ABS, MAX, SQRT, RANDOM
I CAN HAS TIME?     BTW For DATE class, SLEEP function
I CAN HAS FILE?     BTW For DOCUMENT class, file operations
I CAN HAS TEST?     BTW For ASSERT function, testing
```

## File Extensions

All Objective-LOL source files must use the `.olol` extension.

## Community and Support

- **Issues**: [GitHub Issues](https://github.com/bjia56/objective-lol/issues)
- **Source Code**: [GitHub Repository](https://github.com/bjia56/objective-lol)
- **Documentation**: You're reading it!