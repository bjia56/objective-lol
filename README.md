# Objective-LOL

A programming language interpreter inspired by LOLCODE, implemented in Go. Objective-LOL combines the playful syntax of LOLCODE with modern language features including object-oriented programming, a rich type system, and a comprehensive standard library.

## Features

- **Strong Type System**: Five built-in types with automatic type conversion and explicit casting
- **Object-Oriented Programming**: Classes with inheritance, visibility modifiers, and method overriding
- **Functions**: Support for parameters, return values, and recursion
- **Control Flow**: Conditional statements and loops with intuitive syntax
- **Standard Library**: I/O operations, mathematical functions, and time utilities
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
    VISIBLEZ WIT "Hello, World!"
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
    VISIBLE WIT "Hello, "
    VISIBLEZ WIT NAME
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
        VISIBLE WIT "Hi, I'm "
        VISIBLEZ WIT NAME
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
    VISIBLEZ WIT "You are an adult!"
NOPE
    VISIBLEZ WIT "You are a minor"
KTHX

I HAS A VARIABLE COUNTER TEH INTEGR ITZ 5
WHILE COUNTER BIGGR THAN 0
    VISIBLEZ WIT COUNTER
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
./olol tests/01_basic_syntax.olol
```

The test suite includes comprehensive test files covering:

- Basic syntax and variables
- Arithmetic and comparison operations
- Control flow and functions
- Object-oriented features
- Standard library functions
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

### Standard Library

**I/O Functions:**
- `VISIBLE WIT <value>` - Print without newline
- `VISIBLEZ WIT <value>` - Print with newline
- `GIMME` - Read user input

**Math Functions:**
- `ABS`, `MAX`, `MIN`, `SQRT`, `POW`
- `SIN`, `COS` (trigonometric functions)
- `RANDOM`, `RANDINT` (random number generation)

**Time Functions:**
- `NOW`, `YEAR`, `MONTH`, `DAY`, `HOUR`, `MINUTE`, `SECOND`
- `FORMAT_TIME`, `SLEEP`

### Type Casting

```lol
I HAS A VARIABLE NUM_STR TEH STRIN ITZ "123"
I HAS A VARIABLE NUM TEH INTEGR ITZ NUM_STR AS INTEGR  BTW 123

I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159
I HAS A VARIABLE TRUNCATED TEH INTEGR ITZ PI AS INTEGR  BTW 3
```

## Examples

See the `tests/` directory for comprehensive examples ranging from basic syntax to complete programs with classes and inheritance.

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

See the repository for license information.