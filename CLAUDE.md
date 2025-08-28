# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Objective-LOL is a programming language interpreter written in Go that implements a humorous, LOLCODE-inspired language. The codebase follows a traditional interpreter architecture with lexer, parser, AST, and tree-walking interpreter components.

## Build and Development Commands

### Building the Interpreter
```bash
go build -o olol cmd/olol/main.go
```

### Running Tests
```bash
# Run all tests with the provided script
./run_tests.sh

# Run with verbose output to see detailed test results
./run_tests.sh -v

# Run individual test files
./olol tests/01_basic_syntax.olol
./olol tests/15_comprehensive.olol
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
- **cmd/olol/main.go**: Entry point that creates interpreter, registers stdlib, and executes .olol files
- **pkg/parser**: Lexer and recursive descent parser that converts source code to AST
- **pkg/ast**: AST node definitions using visitor pattern for tree traversal
- **pkg/interpreter**: Tree-walking interpreter that executes AST nodes
- **pkg/environment**: Runtime environment and scoping system
- **pkg/types**: Value types and object system for the language

### Standard Library
Located in `pkg/stdlib/`:
- **stdio.go**: I/O functions (VISIBLEZ, VISIBLE, GIMME)
- **math.go**: Mathematical functions (ABS, MAX, MIN, SQRT, POW, SIN, COS, RANDOM)
- **tiem.go**: Time-related functions (YEAR, MONTH, DAY, NOW, FORMAT_TIME, SLEEP)

### Language Features
The interpreter supports:
- Variables with type declarations (INTEGR, DUBBLE, STRIN, BOOL, NOTHIN)
- Functions with parameters and return values
- Classes with inheritance (KITTEH OF), visibility modifiers (EVRYONE/MAHSELF)
- Control flow (IZ/NOPE conditionals, WHILE loops)
- Arithmetic, comparison, and logical operators
- Type casting with AS operator
- Object instantiation and method calls

### Parser Architecture
Uses recursive descent parsing with:
- Token-based lexer in `pkg/parser/lexer.go`
- Parser maintains current/peek tokens for lookahead
- Error collection system for parse errors
- Case-insensitive keyword handling (converted to uppercase internally)

### Interpreter Architecture
Implements visitor pattern:
- AST nodes implement Accept() method taking Visitor
- Interpreter implements all Visit*() methods for each node type
- Environment stack for variable scoping
- Runtime environment manages global functions and classes
- Object instance tracking for method calls and member access

### Testing
Comprehensive test suite in `tests/` directory:
- Files numbered 01-15 covering basic to advanced features
- Each test file is self-documenting with expected outputs
- Tests use `.olol` extension to distinguish from examples
- `run_tests.sh` provides automated test execution with pass/fail reporting

## Development Notes

### File Extensions
- Source files use `.olol` extension (enforced by interpreter)
- Go source follows standard Go conventions

### Error Handling
- Parser collects all errors before failing
- Interpreter returns detailed error messages
- Test runner shows specific failure reasons

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