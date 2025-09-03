# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Objective-LOL is a programming language interpreter written in Go that implements a LOLCODE-inspired language with modern language features. The codebase follows a traditional interpreter architecture with lexer, parser, AST, and tree-walking interpreter components.

## Build and Development Commands

### Building the Interpreter
```bash
go build -o olol cmd/olol/main.go
```

### Running Tests
```bash
# Run all Go unit tests
go test ./...

# Run with verbose output to see detailed test results
go test -v ./...

# Run specific package tests
go test ./pkg/parser
go test ./pkg/interpreter
go test ./pkg/integration

# Run individual .olol test files
./olol pkg/integration/tests/18_precedence_test.olol
./olol pkg/integration/tests/26_simple_exception_test.olol
```

### Code Quality
```bash
# Format Go code
go fmt ./...

# Build and check for compilation errors
go build ./...

# Run with specific test file for debugging
./olol path/to/test.olol
```

## Architecture

The codebase is organized into several key packages:

### Core Components
- **cmd/olol/main.go**: Main entry point (canonical build: `go build -o olol cmd/olol/main.go`)
- **pkg/cli/cli.go**: CLI package with Run() and executeFile() functions used by both main and tests
- **pkg/parser**: Lexer and recursive descent parser that converts source code to AST
- **pkg/ast**: AST node definitions using visitor pattern for tree traversal
- **pkg/interpreter**: Tree-walking interpreter that executes AST nodes
- **pkg/environment**: Runtime environment and scoping system
- **pkg/modules**: Cross-platform module resolution and caching system for file imports
- **pkg/types**: Value types and object system for the language

### Standard Library
Located in `pkg/stdlib/`:
- **stdio.go**: I/O functions (SAYZ, SAY, GIMME)
- **math.go**: Mathematical functions (ABS, MAX, MIN, SQRT, POW, SIN, COS, RANDOM)
- **time.go**: Time-related functions using DATE class (DATE class with YEAR, MONTH, DAY, HOUR, MINUTE, SECOND, MILLISECOND, NANOSECOND, FORMAT methods) and global SLEEP function
- **test.go**: Testing and assertion functions (ASSERT)
- **arrays.go**: BUKKIT array type with methods (PUSH, POP, AT, SET, SIZ, SORT, REVERSE, JOIN, FIND, CONTAINS)
- **maps.go**: BASKIT map type with methods (PUT, GET, CONTAINS, REMOVE, KEYS, VALUES, PAIRS, MERGE, COPY, CLEAR)
- **string.go**: String utility functions (LEN, CONCAT)
- **io.go**: Advanced I/O classes (READER, WRITER, BUFFERED_READER, BUFFERED_WRITER)

### Language Features
The interpreter supports:
- Variables with type declarations (INTEGR, DUBBLE, STRIN, BOOL, NOTHIN, BUKKIT, BASKIT)
- Functions with parameters and return values
- Classes with inheritance (KITTEH OF), visibility modifiers (EVRYONE/MAHSELF), constructor methods
- **Module system with file imports**: Cross-platform POSIX path resolution, caching, circular import detection
- **Exception handling system**: MAYB/OOPS/OOPSIE/ALWAYZ try-catch-finally blocks with string-based exceptions
- Control flow (IZ/NOPE conditionals, WHILE loops)
- Arithmetic, comparison, and logical operators
- Parentheses for expression grouping and precedence override
- Type casting with AS operator
- Object instantiation and method calls

### Parser Architecture
Uses recursive descent parsing with:
- Token-based lexer in `pkg/parser/lexer.go`
- Parser maintains current/peek tokens for lookahead
- Error collection system for parse errors
- Case-insensitive keyword handling (converted to uppercase internally)
- Expression parsing with operator precedence and parentheses support

#### Expression Parsing and Operator Precedence
The parser implements a precedence-climbing algorithm with the following precedence levels (lowest to highest):
1. **OR** - Logical OR operations
2. **AN** - Logical AND operations
3. **SAEM AS** - Equality comparisons
4. **BIGGR THAN, SMALLR THAN** - Relational comparisons
5. **MOAR, LES** - Addition and subtraction
6. **TIEMZ, DIVIDEZ** - Multiplication and division
7. **AS** - Type casting
8. **Primary expressions** - Literals, identifiers, function calls, parenthesized expressions

Parentheses `()` can be used to override operator precedence and group sub-expressions:
- `2 TIEMZ 3 MOAR 4` evaluates to `10` (multiplication first)
- `2 TIEMZ (3 MOAR 4)` evaluates to `14` (parentheses override precedence)
- Nested parentheses are fully supported: `((2 MOAR 3) TIEMZ 4)`

### Interpreter Architecture
Implements visitor pattern:
- AST nodes implement Accept() method taking Visitor
- Interpreter implements all Visit*() methods for each node type
- Environment stack for variable scoping
- Runtime environment manages global functions and classes
- Object instance tracking for method calls and member access

### Exception Handling Architecture
**Exception System** (`pkg/ast/nodes.go`):
- **TryStatementNode**: AST node for MAYB/OOPSIE/ALWAYZ blocks
- **ThrowStatementNode**: AST node for OOPS statements
- **Exception Type**: String-based exception type implementing Go's error interface
- **Exception Propagation**: Exceptions propagate via error return values until caught

**Built-in Exception Integration**:
- Division by zero in arithmetic operations
- Type casting failures in value conversion
- Array bounds violations in BUKKIT operations
- Key not found errors in BASKIT operations
- Undefined variable/function access

### Module System
**File Import Architecture** (`pkg/modules/`):
- **ModuleResolver**: Cross-platform path resolution with POSIX-to-native conversion
- **AST Caching**: Each `.olol` file parsed once and cached by absolute path
- **Environment Caching**: Executed module environments cached to prevent re-execution
- **Circular Import Detection**: Execution-level tracking prevents infinite recursion
- **Context-Aware Resolution**: Import paths resolved relative to importing file's directory

**Import Syntax**:
- Built-in modules: `I CAN HAS STDIO?`
- File modules: `I CAN HAS "math_utils"?` (auto-appends `.olol` extension)
- Selective imports: `I CAN HAS FUNC1 AN FUNC2 FROM "module"?`
- Private declarations: Functions/classes/variables starting with `_` are not importable

**Path Handling**:
- Source code uses POSIX paths: `"utils/helpers"` → `utils/helpers.olol`
- Runtime conversion to native paths (Windows: `utils\helpers.olol`)
- Supports relative: `"../shared"`, absolute: `"/project/modules"`, and nested: `"dir/subdir/module"`

### Testing
Comprehensive test suite with both Go unit tests and .olol integration tests:
- **Go unit tests**: Located in `*_test.go` files throughout the codebase, run with `go test ./...`
- **Integration tests**: Located in `pkg/integration/tests/` directory with `.olol` extension
- **Functional tests**: `pkg/integration/functional_test.go` automatically runs all `.olol` files in tests directory
- Each .olol test file is self-documenting with expected outputs
- Module import tests in `pkg/integration/tests/test_modules/` demonstrate file import capabilities
- Exception handling integration tests validate the exception system

## Development Notes

### File Extensions
- Source files use `.olol` extension (enforced by interpreter)
- Go source follows standard Go conventions

### Error Handling
- **Exception System**: Comprehensive exception handling with MAYB/OOPS/OOPSIE/ALWAYZ syntax
- **Built-in Exceptions**: Automatic exception throwing for division by zero, type casting errors, array bounds errors, undefined variables
- **Exception Propagation**: Exceptions propagate up the call stack until caught
- **Parser Error Collection**: Parser collects all errors before failing
- **Runtime Error Messages**: Interpreter provides detailed error messages
- **Test Error Reporting**: Test runner shows specific failure reasons

### Code Style
- Follows standard Go formatting (use `go fmt`)
- AST nodes use descriptive names ending in "Node"
- Visitor methods follow "Visit" + node type naming

### Adding New Features
1. Add tokens to lexer if needed
2. Create AST node in `pkg/ast/nodes.go`
3. Add visitor method to Visitor interface
4. Implement parsing logic in parser
5. Implement execution logic in interpreter
6. Add comprehensive tests
