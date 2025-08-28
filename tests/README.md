# Objective-LOL Test Suite

This directory contains comprehensive tests for the Objective-LOL programming language implementation. The tests are organized to cover all major language features and edge cases.

## Test Files Overview

### Basic Language Features
1. **01_basic_syntax.olol** - Basic syntax, comments, literals, and output
2. **02_variables.olol** - Variable declarations, assignments, and data types
3. **03_arithmetic.olol** - Arithmetic operations (MOAR, LES, TIEMZ, DIVIDEZ)
4. **04_comparisons.olol** - Comparison operators (BIGGR THAN, SMALLR THAN, SAEM AS)
5. **05_logical_ops.olol** - Logical operations (AN, OR, NOT)
6. **06_control_flow.olol** - Control structures (IZ/NOPE, WHILE loops)

### Advanced Features
7. **07_functions.olol** - Function definitions, parameters, return values, recursion
8. **08_classes_basic.olol** - Basic classes, member variables, methods, visibility
9. **09_inheritance.olol** - Inheritance (KITTEH OF), method overriding, polymorphism
10. **10_stdlib_math.olol** - Math library functions (ABS, MAX, MIN, SQRT, POW, SIN, COS, RANDOM)
11. **11_stdlib_time.olol** - Time functions (YEAR, MONTH, DAY, NOW, FORMAT_TIME, SLEEP)
12. **12_stdlib_io.olol** - I/O functions (VISIBLEZ, VISIBLE, GIMME)

### Advanced Topics
13. **13_type_casting.olol** - Type conversions (AS operator)
14. **14_error_cases.olol** - Error handling and edge cases
15. **15_comprehensive.olol** - Complex integration test using all features

## Language Syntax Summary

### Basic Structure
```lol
HAI ME TEH FUNCSHUN MAIN
    BTW This is a comment
    VISIBLEZ WIT "Hello, World!"
KTHXBAI
```

### Variables
```lol
I HAS A VARIABLE NAME TEH STRIN ITZ "Alice"
HAI ME TEH VARIABLE AGE TEH INTEGR ITZ 25
HAI ME TEH LOCKD VARIABLE CONSTANT TEH DUBBLE ITZ 3.14159
```

### Data Types
- `INTEGR` - Integer numbers
- `DUBBLE` - Floating point numbers  
- `STRIN` - Text strings
- `BOOL` - Boolean values (`YEZ`/`NO`)
- `NOTHIN` - Null/empty value

### Arithmetic
- `MOAR` - Addition (+)
- `LES` - Subtraction (-)
- `TIEMZ` - Multiplication (*)
- `DIVIDEZ` - Division (/)

### Comparisons  
- `BIGGR THAN` - Greater than (>)
- `SMALLR THAN` - Less than (<)
- `SAEM AS` - Equality (==)

### Logical Operations
- `AN` - Logical AND
- `OR` - Logical OR  
- `NOT` - Logical NOT

### Control Flow
```lol
IZ condition QUESTION
    BTW then block
NOPE
    BTW else block
KTHX

WHILE condition
    BTW loop body
KTHX
```

### Functions
```lol
HAI ME TEH FUNCSHUN FUNCTION_NAME TEH RETURN_TYPE WIT PARAM TEH TYPE
    GIVEZ return_value
KTHXBAI
```

### Classes
```lol
HAI ME TEH CLAS CLASS_NAME KITTEH OF PARENT_CLASS
    EVRYONE  BTW public visibility
    DIS TEH VARIABLE MEMBER_VAR TEH TYPE ITZ value
    
    DIS TEH FUNCSHUN METHOD_NAME WIT PARAM TEH TYPE
        BTW method body
    KTHX
    
    MAHSELF  BTW private visibility
    DIS TEH VARIABLE PRIVATE_VAR TEH TYPE
KTHXBAI
```

### Object Usage
```lol
I HAS A VARIABLE OBJ TEH CLASS_NAME ITZ NEW CLASS_NAME
OBJ DO METHOD_NAME WIT argument
VISIBLEZ WIT OBJ MEMBER_VAR
OBJ MEMBER_VAR ITZ new_value
```

## Running Tests

To run the tests with the Objective-LOL interpreter:

```bash
# Compile the Go implementation
go build -o olol cmd/olol/main.go

# Run individual tests
./olol tests/01_basic_syntax.olol
./olol tests/02_variables.olol
# ... etc

# Run all tests
for test in tests/*.olol; do
    echo "Running $test"
    ./olol "$test"
    echo "---"
done
```

## Expected Output

Each test file produces detailed output showing:
- Test section headers (=== Test Name ===)
- Expected behavior demonstrations
- Variable values and computation results
- Function and method call results
- Error handling demonstrations

The tests are designed to be self-documenting, showing both correct usage patterns and expected outputs for each language feature.

## Test Coverage

These tests cover:
- ✅ All basic syntax elements
- ✅ All data types and literals
- ✅ All operators (arithmetic, comparison, logical)
- ✅ All control flow structures  
- ✅ Function definitions and calls
- ✅ Class definitions and inheritance
- ✅ Object creation and method calls
- ✅ All standard library functions
- ✅ Type casting and conversions
- ✅ Error conditions and edge cases
- ✅ Complex integration scenarios

## Notes

- Tests use `.olol` extension to distinguish from original examples
- All tests follow the case-insensitive convention (identifiers converted to uppercase internally)
- Tests include both positive test cases and error condition handling
- The comprehensive test demonstrates real-world usage patterns combining multiple features