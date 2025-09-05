# Objective-LOL

A programming language inspired by LOLCODE.

## Key Features

- **Strong Type System**: Seven built-in types with automatic and explicit casting
- **Object-Oriented Programming**: Classes with inheritance, constructors, and visibility modifiers
- **Functions & Scoping**: Lexical scoping with parameters, return values, and recursion
- **Module System**: Import built-in and file modules with selective import syntax
- **Exception Handling**: Comprehensive try-catch-finally blocks with custom exceptions
- **Rich Collections**: Dynamic arrays and maps with extensive methods
- **Standard Library**: File I/O, mathematics, time, string utilities, threading, and buffered I/O

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
    I CAN HAS STDIO?
    SAYZ WIT "Hello, World!"
KTHXBAI
```

Run it:

```bash
./olol hello.olol
```

### Next Steps

- [Getting Started Guide](docs/language-guide/getting-started.md)
- [Quick Reference](docs/QUICK_REFERENCE.md)
- [Full Documentation](docs/README.md)

## Building and Development

### Build Commands

```bash
# Build the interpreter
go build -o olol cmd/olol/main.go

# Build LSP server (optional)
go build -o olol-lsp cmd/olol-lsp/main.go

# Format code
go fmt ./...

# Check compilation
go build ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run integration tests
./olol pkg/integration/tests/01_basic_syntax.olol
```

## License

MIT
