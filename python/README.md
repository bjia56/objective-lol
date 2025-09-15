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