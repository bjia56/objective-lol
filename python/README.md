# Objective-LOL Python Bindings

Python bindings for the Objective-LOL programming language - a modern, strongly-typed language inspired by LOLCODE.

## Installation

```bash
pip install objective-lol
```

## Quick Start

```python
import objective_lol as olol

# Create a VM instance
vm = olol.ObjectiveLOLVM()

# Execute Objective-LOL code
code = """
I CAN HAS STDIO?
HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE X TEH INTEGR ITZ 42
    SAYZ WIT X
KTHXBAI
"""

vm.execute(code)  # Prints: 42
```

## Features

- **Execute Objective-LOL code**: Run complete programs or code snippets
- **Python integration**: Call Python functions from Objective-LOL and vice versa
- **Type conversion**: Automatic conversion between Python and Objective-LOL types
- **Module system**: Import and use custom modules
- **Class definitions**: Define and use classes from Python

## Advanced Usage

### Defining Python Functions for Objective-LOL

```python
vm = olol.ObjectiveLOLVM()

def add_numbers(a, b):
    return a + b

vm.define_function("add_numbers", add_numbers)

code = """
I CAN HAS STDIO?
HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE RESULT TEH INTEGR ITZ add_numbers WIT 10 AN WIT 20
    SAYZ WIT RESULT
KTHXBAI
"""

vm.execute(code)  # Prints: 30
```

### Working with Classes

```python
vm = olol.ObjectiveLOLVM()

class Calculator:
    def add(self, x, y):
        return x + y

    def multiply(self, x, y):
        return x * y

vm.define_class(Calculator)

code = """
I CAN HAS STDIO?
HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE CALC TEH Calculator ITZ NEW Calculator
    I HAS A VARIABLE SUM TEH INTEGR ITZ CALC DO add WIT 5 AN WIT 3
    SAYZ WIT SUM
KTHXBAI
"""

vm.execute(code)  # Prints: 8
```

## Type Mapping

| Objective-LOL Type | Python Type |
|-------------------|-------------|
| INTEGR            | int         |
| DUBBLE            | float       |
| STRIN             | str         |
| BOOL              | bool        |
| NOTHIN            | None        |
| BUKKIT (array)    | list        |
| BASKIT (map)      | dict        |

## Links

- [Main Project](https://github.com/bjia56/objective-lol)
- [Language Documentation](https://github.com/bjia56/objective-lol/tree/main/docs)

## License

MIT License

## Changelog

### [0.0.7] - 2026-04-30

#### Changed
- Proxy methods now determine their async/sync interface based on whether the VM was initialized with an asyncio loop, rather than inspecting the wrapped method's own async signature
- Error messages from Python callbacks now include the exception type (e.g. `ValueError: ...` instead of just the message text)

#### Fixed
- Member variable getter errors are now propagated as interpreter exceptions instead of being silently swallowed
- Function and method call failure messages now include the function/method name for easier debugging

### [0.0.6] - 2026-04-30

#### Changed
- Refactored asyncio threading shims into factory methods (`_run_in_loop`, `_run_blocking`, `_make_loop_method_wrapper`, `_make_loop_wrapper`) to eliminate duplicated thread+Future patterns
- Getter and setter registrations now route through the asyncio event loop when one is provided, consistent with method dispatch

### [0.0.5] - 2026-04-29

#### Changed
- Replaced `prefer_async_loop` with an explicit `asyncio_loop` to specify which event loop to use
- Changed interop of sync Python methods to be called within the asyncio loop, if one is provided
- Changed interop of async Python methods to require an asyncio loop

### [0.0.4] - 2026-04-29

#### Added
- Added a mechanism to specify the working directory of the Objective-LOL interpreter from Python at initialization time

#### Fixed
- Fixed interop of setting primitive values on Python objects from within Objective-LOL

### [0.0.3] - 2026-03-24

#### Fixed
- Fixed interop of nested lists and dicts between Python and Objective-LOL

### [0.0.2] - 2025-09-16

#### Changed
- Modified octal prefix from `0` to `0o`

#### Fixed
- Fixed a stack corruption crash on Windows by using compatibility function calls with separate goroutines

### [0.0.1] - 2025-09-14

#### Added
- Initial Python bindings for Objective-LOL interpreter
- Pre-built wheels for Windows, macOS, and Linux for Python 3.9+
