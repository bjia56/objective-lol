# Getting Started with Objective-LOL

Objective-LOL is a programming language inspired by LOLCODE. It features strong typing, object-oriented programming, functions, modules, and exception handling.

All Objective-LOL source files must use the `.olol` file extension.

## Installation and Building

### Prerequisites
- Go 1.21 or later

### Building the Interpreter

```bash
# Clone the repository
git clone https://github.com/bjia56/objective-lol.git
cd objective-lol

# Build the interpreter
go build -o olol cmd/olol/main.go
```

## Your First Program

Create a file named `hello.olol`:

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

Output:
```
Hello, World!
```

## Basic Program Structure

Every Objective-LOL program must have a `MAIN` function as the entry point:

```lol
BTW Program entry point
HAI ME TEH FUNCSHUN MAIN
    BTW Your code goes here
KTHXBAI
```

## Comments and Case Sensitivity

- Comments start with `BTW` and continue to the end of the line
- Keywords and variables are **case-insensitive** and automatically converted to uppercase

```lol
BTW This is a single-line comment
hai me teh funcshun main    BTW Same as HAI ME TEH FUNCSHUN MAIN
    SAYZ WIT "Hello"        BTW This is also a comment
kthxbai
```

## Next Steps

- [Syntax Basics](syntax-basics.md) - Learn about data types, variables, and operators
- [Control Flow](control-flow.md) - Conditionals, loops, and exception handling
- [Functions](functions.md) - Function declaration and usage
- [Classes](classes.md) - Object-oriented programming features
- [Modules](modules.md) - Import system and code organization
