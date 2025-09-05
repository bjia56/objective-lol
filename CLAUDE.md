# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Objective-LOL is a tree-walking interpreter for a programming language inspired by LOLCODE with modern features. The language includes strong typing, object-oriented programming, modules, and exception handling, with LOLCODE-inspired syntax.

## Build Commands

### Main Interpreter
```bash
# Build the main interpreter
go build -o olol cmd/olol/main.go

# Build LSP server (optional)
go build -o olol-lsp cmd/olol-lsp/main.go

# Build all packages (check compilation)
go build ./...

# Format code
go fmt ./...
```

### Testing
```bash
# Run all unit tests
go test ./...

# Run with verbose output
go test -v ./...

# Run a single integration test
./olol pkg/integration/tests/01_basic_syntax.olol

# Run integration tests programmatically
go test ./pkg/integration/
```

### Development
```bash
# Run interpreter on a file
./olol program.olol

# Check if binary exists and build if needed
[ -f olol ] || go build -o olol cmd/olol/main.go
```

## Architecture Overview

### Core Components

**Parser & Lexer** (`pkg/parser/`):
- Lexer tokenizes LOLCODE-style syntax (`HAI`, `KTHXBAI`, etc.)
- Parser builds AST from tokens with support for expressions, functions, classes

**Interpreter** (`pkg/interpreter/`):
- Tree-walking interpreter with context support
- Tracks current environment, class context, and object instances
- Handles module imports and scope management

**Environment** (`pkg/environment/`):
- Multi-layered scoping system (global, module, local, class)
- Value system with 7 types: INTEGR, DUBBLE, STRIN, BOOL, NOTHIN, BUKKIT, BASKIT
- Object-oriented features with inheritance and visibility

**Standard Library** (`pkg/stdlib/`):
- Modular system with STDIO, MATH, TIME, STRING, TEST, IO, THREAD, FILE
- Each module provides domain-specific functionality
- Thread-safe implementation with proper concurrency support

**Module System** (`pkg/modules/`):
- Resolves both built-in and file-based modules
- Supports selective imports (`I CAN HAS FUNC FROM MODULE?`)
- Handles relative path resolution and module caching

**LSP Support** (`pkg/lsp/`):
- Language Server Protocol implementation
- Semantic analysis and diagnostics
- Workspace management for IDE integration

### Key Design Patterns

**Visitor Pattern**: AST nodes implement visitor pattern for interpretation
**Environment Chain**: Nested scopes with lexical scoping rules
**Module Resolution**: Hierarchical module loading with caching
**Error Propagation**: Go-style error handling with LOLCODE exception syntax

### File Organization

- `/cmd/` - Entry points for binaries (interpreter and LSP)
- `/pkg/` - Core library code organized by functionality
- `/docs/` - Comprehensive language documentation
- `/pkg/integration/tests/` - Integration test suite with `.olol` files

## Language Specifics

**File Extension**: All source files must use `.olol` extension
**Entry Point**: Programs must have a `MAIN` function
**Case Sensitivity**: Keywords are case-insensitive
**Syntax Style**: LOLCODE-inspired (`HAI ME TEH FUNCSHUN`, `KTHXBAI`)

## Testing Strategy

**Unit Tests**: Each package has `*_test.go` files for Go unit tests
**Integration Tests**: `.olol` files in `pkg/integration/tests/` test end-to-end functionality
**Test Categories**: Basic syntax, stdlib modules, OOP features, exception handling, threading

## Language Example


```lol
BTW Calculator with classes, functions, and exception handling
I CAN HAS STDIO?
I CAN HAS MATH?

BTW Global variable
HAI ME TEH VARIABLE PI TEH DUBBLE ITZ 3.14159

HAI ME TEH CLAS CALCULATOR
    EVRYONE
    DIS TEH VARIABLE RESULT TEH DUBBLE ITZ 0.0

    DIS TEH FUNCSHUN ADD TEH DUBBLE WIT X TEH DUBBLE AN WIT Y TEH DUBBLE
        RESULT ITZ X MOAR Y
        GIVEZ RESULT
    KTHX

    DIS TEH FUNCSHUN DIVIDE TEH DUBBLE WIT X TEH DUBBLE AN WIT Y TEH DUBBLE
        MAYB
            IZ Y SAEM AS 0.0?
                OOPS "Division by zero!"
            KTHX
            RESULT ITZ X DIVIDEZ Y
            GIVEZ RESULT
        OOPSIE ERR
            SAYZ WIT ERR
            GIVEZ 0.0
        KTHX
    KTHX

    DIS TEH FUNCSHUN CIRCLE_AREA TEH DUBBLE WIT RADIUS TEH DUBBLE
        GIVEZ PI TIEMZ RADIUS TIEMZ RADIUS
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE CALC TEH CALCULATOR ITZ NEW CALCULATOR

    BTW Addition and arrays
    I HAS A VARIABLE NUMS TEH BUKKIT ITZ NEW BUKKIT
    NUMS DO PUSH WIT 5.0
    NUMS DO PUSH WIT 3.0

    I HAS A VARIABLE SUM TEH DUBBLE ITZ CALC DO ADD WIT (NUMS DO AT WIT 0) AN WIT (NUMS DO AT WIT 1)
    SAY WIT "5 + 3 = "
    SAYZ WIT SUM

    BTW Global variable usage
    I HAS A VARIABLE AREA TEH DUBBLE ITZ CALC DO CIRCLE_AREA WIT 2.0
    SAY WIT "Circle area (r=2): "
    SAYZ WIT AREA

    BTW Test exception handling
    CALC DO DIVIDE WIT 10.0 AN WIT 0.0
KTHXBAI
```

## Common Development Tasks

**Adding New Standard Library Functions**: Extend modules in `pkg/stdlib/`
**Language Features**: Modify parser, AST nodes, and interpreter
**LSP Features**: Update semantic analysis in `pkg/lsp/analyzer/`
**Testing**: Add both Go unit tests and `.olol` integration tests