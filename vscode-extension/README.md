# Objective-LOL Language Support for VS Code

This extension provides language support for Objective-LOL, a LOLCODE-inspired programming language with modern features.

## Features

- **Syntax highlighting** for Objective-LOL files (`.olol`)
- **Language Server Protocol (LSP) support** with:
  - Syntax error diagnostics
  - Code completion (keywords, variables, functions)
  - Hover information
  - Go-to-definition
- **Auto-closing pairs** for parentheses and quotes
- **Code folding** for functions, classes, and control structures
- **Comment support** with `BTW` line comments

## Requirements

To use the LSP features, you need to have the Objective-LOL LSP server installed:

1. Build the LSP server:
   ```bash
   go build -o olol-lsp cmd/olol-lsp/main.go
   ```

2. Make sure `olol-lsp` is in your PATH, or configure the path in VS Code settings.

## Extension Settings

This extension contributes the following settings:

* `objective-lol.lspPath`: Path to the Objective-LOL LSP server executable (default: `olol-lsp`)
* `objective-lol.enableLSP`: Enable/disable Language Server Protocol support (default: `true`)
* `objective-lol.trace.server`: Traces the communication between VS Code and the language server (`off`, `messages`, `verbose`)

## Usage

1. Create a file with the `.olol` extension
2. Start coding in Objective-LOL:

```objective-lol
HAI ME TEH FUNCSHUN MAIN

TEH VARIABLE greeting ITZ STRIN WIT "Hello, World!"
SAYZ greeting

FUNCSHUN add WIT INTEGR a AN INTEGR b GIVEZ INTEGR
    GIVEZ UP a MOAR b
KTHX

SAYZ add WIT 5 AN 3

KTHXBAI
```

## Language Features

### Syntax Highlighting

The extension provides comprehensive syntax highlighting for:
- Keywords (`HAI`, `KTHXBAI`, `TEH`, `FUNCSHUN`, etc.)
- Types (`INTEGR`, `DUBBLE`, `STRIN`, `BOOL`)
- Operators (`MOAR`, `LES`, `TIEMZ`, `SAEM`)
- Comments (`BTW`)
- Strings and numbers

### LSP Features

When the LSP server is running, you get:
- **Diagnostics**: Real-time syntax error detection
- **Auto-completion**: Smart suggestions for keywords and symbols
- **Hover**: Type information and documentation
- **Go-to-definition**: Navigate to symbol definitions

## Commands

- `Objective-LOL: Restart Language Server` - Restarts the LSP server connection

## Known Issues

- The LSP server is still in early development
- Some advanced language features may not be fully supported yet

## Release Notes

### 0.1.0

Initial release with basic language support and LSP integration.

## Contributing

This extension is part of the Objective-LOL language project. Contributions are welcome!

## License

This extension is provided as-is for the Objective-LOL programming language.