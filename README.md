# Objective-LOL

A programming language interpreter inspired by LOLCODE, implemented in Go. Objective-LOL combines the syntax of LOLCODE with modern language features including object-oriented programming, a rich type system, and a comprehensive standard library.

## Features

- **Strong Type System**: Six built-in types including arrays with automatic type conversion and explicit casting
- **Object-Oriented Programming**: Classes with inheritance, constructors, visibility modifiers, and method overriding
- **Functions**: Support for parameters, return values, and recursion with lexical scoping
- **Control Flow**: Conditional statements and loops with intuitive syntax
- **Module System**: Import system supporting both built-in and file modules with selective imports
- **Exception Handling**: Comprehensive try-catch-finally blocks with built-in and custom exceptions
- **Standard Library**: I/O operations, mathematical functions, time utilities, string functions, and buffered I/O
- **Case-Insensitive**: Keywords automatically converted to uppercase for consistency

## Quick Start

### Installation

```bash
git clone https://github.com/bjia56/objective-lol.git
cd objective-lol
go build -o olol cmd/olol/main.go
```

### Hello World

Create `hello.olol`:

```lol
BTW This is a comment - Hello World program

HAI ME TEH FUNCSHUN MAIN
    SAYZ WIT "Hello, World!"
KTHXBAI
```

Run it:

```bash
./olol hello.olol
```

## Language Overview

### Data Types

- **INTEGR**: 64-bit signed integers (supports decimal and hex: `42`, `0xFF`)
- **DUBBLE**: 64-bit floating point numbers (`3.14159`, `-2.5`)
- **STRIN**: Strings with escape sequences (`"Hello \"World\"!"`)
- **BOOL**: Boolean values (`YEZ` for true, `NO` for false)
- **BUKKIT**: Dynamic arrays with rich manipulation methods
- **NOTHIN**: Null/void type

### Variables

```lol
I HAS A VARIABLE NAME TEH STRIN ITZ "Alice"
I HAS A VARIABLE AGE TEH INTEGR ITZ 25
I HAS A LOCKD VARIABLE PI TEH DUBBLE ITZ 3.14159  BTW Constant
```

### Functions

```lol
HAI ME TEH FUNCSHUN GREET WIT NAME TEH STRIN
    SAY WIT "Hello, "
    SAYZ WIT NAME
KTHXBAI

HAI ME TEH FUNCSHUN ADD TEH INTEGR WIT X TEH INTEGR AN WIT Y TEH INTEGR
    GIVEZ X MOAR Y
KTHXBAI
```

### Classes and Objects

```lol
HAI ME TEH CLAS PERSON
    EVRYONE    BTW Public
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"

    DIS TEH FUNCSHUN INTRODUCE
        SAY WIT "Hi, I'm "
        SAYZ WIT NAME
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE PERSON1 TEH PERSON ITZ NEW PERSON
    PERSON1 NAME ITZ "Bob"
    PERSON1 DO INTRODUCE
KTHXBAI
```

### Control Flow

```lol
IZ AGE BIGGR THAN 17?
    SAYZ WIT "You are an adult!"
NOPE
    SAYZ WIT "You are a minor"
KTHX

I HAS A VARIABLE COUNTER TEH INTEGR ITZ 5
WHILE COUNTER BIGGR THAN 0
    SAYZ WIT COUNTER
    COUNTER ITZ COUNTER LES 1
KTHX
```

## Building and Testing

### Build Commands

```bash
# Build the interpreter
go build -o olol cmd/olol/main.go

# Format Go code
go fmt ./...

# Check for compilation errors
go build ./...
```

### Running Tests

```bash
# Run all tests
./run_tests.sh

# Run with verbose output
./run_tests.sh -v

# Run individual test file
./olol pkg/integration/tests/01_basic_syntax.olol
```

The test suite includes comprehensive test files covering:

- Basic syntax and variables
- Arithmetic and comparison operations
- Control flow and functions
- Object-oriented features with constructors
- Standard library functions across all modules
- Module system with file imports and selective imports
- Exception handling with try-catch-finally blocks
- Array operations with BUKKIT
- Buffered I/O operations
- Error handling and edge cases

## Documentation

- **[Language Reference](LANGUAGE_REFERENCE.md)**: Complete language specification with detailed examples

## Architecture

The interpreter follows a traditional architecture:

- **Lexer** (`pkg/parser/lexer.go`): Tokenizes source code
- **Parser** (`pkg/parser/parser.go`): Recursive descent parser generating AST
- **AST** (`pkg/ast/nodes.go`): Abstract syntax tree nodes with visitor pattern
- **Interpreter** (`pkg/interpreter/interpreter.go`): Tree-walking interpreter
- **Environment** (`pkg/environment/environment.go`): Variable scoping and runtime environment
- **Standard Library** (`pkg/stdlib/`): Built-in functions for I/O, math, and time

## Language Features

### Operators

| Operation | Syntax | Example |
|-----------|--------|---------|
| Addition | `MOAR` | `5 MOAR 3` → 8 |
| Subtraction | `LES` | `10 LES 4` → 6 |
| Multiplication | `TIEMZ` | `6 TIEMZ 7` → 42 |
| Division | `DIVIDEZ` | `15 DIVIDEZ 3` → 5.0 |
| Greater than | `BIGGR THAN` | `5 BIGGR THAN 3` → YEZ |
| Less than | `SMALLR THAN` | `3 SMALLR THAN 5` → YEZ |
| Equal to | `SAEM AS` | `5 SAEM AS 5` → YEZ |
| Logical AND | `AN` | `YEZ AN NO` → NO |
| Logical OR | `OR` | `YEZ OR NO` → YEZ |

### Module System and Standard Library

Objective-LOL features a comprehensive module system supporting both built-in standard library modules and custom file modules.

#### Import Syntax

```lol
BTW Import entire modules
I CAN HAS STDIO?    BTW All I/O functions
I CAN HAS MATH?     BTW All mathematical functions
I CAN HAS TIME?     BTW Time functions and DATE class
I CAN HAS STRING?   BTW String utility functions
I CAN HAS IO?       BTW Advanced I/O classes
I CAN HAS TEST?     BTW Testing functions

BTW Selective imports
I CAN HAS SAY AN SAYZ FROM STDIO?
I CAN HAS ABS AN MAX FROM MATH?
I CAN HAS DATE AN SLEEP FROM TIME?

BTW File module imports
I CAN HAS "my_module"?              BTW Full import
I CAN HAS FUNC1 AN FUNC2 FROM "utils"?  BTW Selective import
```

#### Available Modules

**STDIO Module:**
- `SAY WIT <value>` - Print without newline
- `SAYZ WIT <value>` - Print with newline
- `GIMME` - Read user input

**MATH Module:**
- `ABS`, `MAX`, `MIN`, `SQRT`, `POW`
- `SIN`, `COS` (trigonometric functions)
- `RANDOM`, `RANDINT` (random number generation)

**TIME Module:**
- `DATE` class with methods: `YEAR`, `MONTH`, `DAY`, `HOUR`, `MINUTE`, `SECOND`, `MILLISECOND`, `NANOSECOND`, `FORMAT`
- `SLEEP` (global function)

**STRING Module:**
- `LEN WIT <string>` - Get string length
- `CONCAT WIT <str1> AN WIT <str2>` - Concatenate strings

**IO Module:**
- `READER`, `WRITER` - Base I/O classes
- `BUFFERED_READER`, `BUFFERED_WRITER` - Buffered I/O for performance

**TEST Module:**
- `ASSERT WIT <condition>` - Assertion for testing

**Built-in BUKKIT Arrays:**
- Dynamic arrays with methods like `PUSH`, `POP`, `AT`, `SET`, `SORT`, `REVERSE`, `JOIN`, `FIND`, `CONTAINS`

### Exception Handling

Objective-LOL supports comprehensive exception handling with try-catch-finally blocks:

```lol
MAYB
    I HAS A VARIABLE RESULT TEH DUBBLE ITZ 10.0 DIVIDEZ 0.0  BTW Throws exception
    SAYZ WIT "This won't print"
OOPSIE ERROR_MSG
    SAYZ WIT "Caught exception: "
    SAYZ WIT ERROR_MSG  BTW "Division by zero"
ALWAYZ
    SAYZ WIT "This always executes"
KTHX

BTW Throwing custom exceptions
HAI ME TEH FUNCSHUN VALIDATE_AGE WIT AGE TEH INTEGR
    IZ AGE SMALLR THAN 0?
        OOPS "Age cannot be negative!"
    KTHX
    GIVEZ AGE
KTHXBAI
```

**Built-in exceptions are automatically thrown for:**
- Division by zero
- Array bounds violations
- Type casting errors
- Undefined variables/functions

### Constructor Methods

Classes can have constructor methods with the same name as the class:

```lol
HAI ME TEH CLAS POINT
    EVRYONE
    DIS TEH VARIABLE X TEH INTEGR ITZ 0
    DIS TEH VARIABLE Y TEH INTEGR ITZ 0
    
    BTW Constructor - same name as class
    DIS TEH FUNCSHUN POINT WIT X_VAL TEH INTEGR AN WIT Y_VAL TEH INTEGR
        X ITZ X_VAL
        Y ITZ Y_VAL
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    BTW Create with constructor arguments
    I HAS A VARIABLE ORIGIN TEH POINT ITZ NEW POINT WIT 0 AN WIT 0
    I HAS A VARIABLE CORNER TEH POINT ITZ NEW POINT WIT 10 AN WIT 5
KTHXBAI
```

### Arrays (BUKKIT)

Dynamic arrays with rich functionality:

```lol
I HAS A VARIABLE NUMS TEH BUKKIT ITZ NEW BUKKIT
NUMS DO PUSH WIT 10
NUMS DO PUSH WIT 20
NUMS DO PUSH WIT 30

SAYZ WIT NUMS SIZ                    BTW 3
SAYZ WIT NUMS DO AT WIT 1           BTW 20
NUMS DO SET WIT 1 AN WIT 99         BTW Set element
I HAS A VARIABLE CSV TEH STRIN ITZ NUMS DO JOIN WIT ", "  BTW "10, 99, 30"

BTW Rich array operations
NUMS DO SORT                        BTW Sort in-place
NUMS DO REVERSE                     BTW Reverse in-place
SAYZ WIT NUMS DO FIND WIT 99        BTW Find index
SAYZ WIT NUMS DO CONTAINS WIT 10    BTW Check existence
```

### Type Casting

```lol
I HAS A VARIABLE NUM_STR TEH STRIN ITZ "123"
I HAS A VARIABLE NUM TEH INTEGR ITZ NUM_STR AS INTEGR  BTW 123

I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159
I HAS A VARIABLE TRUNCATED TEH INTEGR ITZ PI AS INTEGR  BTW 3
```

## Examples

See the `pkg/integration/tests/` directory for comprehensive examples ranging from basic syntax to complete programs with classes, inheritance, module imports, and exception handling.

## Requirements

- Go 1.21 or higher
- Source files must use `.olol` extension

## Contributing

This project follows standard Go conventions. Key guidelines:

- Use `go fmt` for code formatting
- Add comprehensive tests for new features
- Follow the existing AST visitor pattern for new language constructs
- Update documentation when adding new features

## License

MIT
