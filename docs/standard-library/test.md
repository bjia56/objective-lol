# TEST Module - Testing and Assertions

The TEST module provides testing functionality for verifying program behavior and creating test suites.

## Importing TEST Module

```lol
BTW Import entire module
I CAN HAS TEST?

BTW Selective import
I CAN HAS ASSERT FROM TEST?
```

## ASSERT Function

The ASSERT function verifies that a condition is true and throws an exception if it's false. This is the primary testing tool for validating program behavior.

### Syntax

```lol
ASSERT WIT <condition>
```

**Parameters:**
- **condition**: Any type - The condition to test for truthiness

**Behavior:**
- If condition is truthy: Function returns normally
- If condition is falsy: Throws "Assertion failed" exception

### Truthiness Rules

The ASSERT function uses Objective-LOL's truthiness rules:

| Type | Falsy Values | Truthy Values |
|------|-------------|---------------|
| `BOOL` | `NO` | `YEZ` |
| `INTEGR` | `0` | Any non-zero number |
| `DUBBLE` | `0.0` | Any non-zero number |
| `STRIN` | `""` (empty string) | Any non-empty string |
| `NOTHIN` | Always falsy | Never truthy |
| Objects | Never falsy | Always truthy |

## Basic Testing Examples

### Simple Assertions

```lol
I CAN HAS TEST?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN TEST_BASIC_ASSERTIONS
    SAYZ WIT "=== Basic Assertion Tests ==="
    
    BTW Test boolean values
    ASSERT WIT YEZ              BTW Passes
    SAYZ WIT "✓ Boolean true assertion passed"
    
    BTW Test integer values
    ASSERT WIT 42               BTW Passes (non-zero)
    SAYZ WIT "✓ Positive integer assertion passed"
    
    ASSERT WIT -5               BTW Passes (non-zero)
    SAYZ WIT "✓ Negative integer assertion passed"
    
    BTW Test string values
    ASSERT WIT "Hello"          BTW Passes (non-empty)
    SAYZ WIT "✓ Non-empty string assertion passed"
    
    SAYZ WIT "All basic assertions completed!"
KTHXBAI
```

### Assertion Failures

```lol
I CAN HAS TEST?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN TEST_ASSERTION_FAILURES
    SAYZ WIT "=== Testing Assertion Failures ==="
    
    MAYB
        ASSERT WIT NO           BTW This will fail
        SAYZ WIT "This should not print"
    OOPSIE ASSERT_ERROR
        SAYZ WIT "✓ Boolean false assertion failed as expected"
    KTHX
    
    MAYB
        ASSERT WIT 0            BTW This will fail
        SAYZ WIT "This should not print"
    OOPSIE ASSERT_ERROR
        SAYZ WIT "✓ Zero integer assertion failed as expected"
    KTHX
    
    MAYB
        ASSERT WIT ""           BTW This will fail
        SAYZ WIT "This should not print"
    OOPSIE ASSERT_ERROR
        SAYZ WIT "✓ Empty string assertion failed as expected"
    KTHX
    
    MAYB
        ASSERT WIT NOTHIN       BTW This will fail
        SAYZ WIT "This should not print"
    OOPSIE ASSERT_ERROR
        SAYZ WIT "✓ NOTHIN assertion failed as expected"
    KTHX
    
    SAYZ WIT "All failure tests completed!"
KTHXBAI
```

## Testing Functions

### Unit Test Pattern

```lol
I CAN HAS TEST?
I CAN HAS STDIO?

BTW Function to test
HAI ME TEH FUNCSHUN ADD TEH INTEGR WIT A TEH INTEGR AN WIT B TEH INTEGR
    GIVEZ A MOAR B
KTHXBAI

BTW Function to test
HAI ME TEH FUNCSHUN IS_EVEN TEH BOOL WIT NUM TEH INTEGR
    GIVEZ (NUM TIEMZ 2) DIVIDEZ 2 SAEM AS NUM
KTHXBAI

HAI ME TEH FUNCSHUN TEST_ADD_FUNCTION
    SAYZ WIT "=== Testing ADD Function ==="
    
    BTW Test positive numbers
    I HAS A VARIABLE RESULT1 TEH INTEGR ITZ ADD WIT 5 AN WIT 3
    ASSERT WIT RESULT1 SAEM AS 8
    SAYZ WIT "✓ ADD(5, 3) = 8"
    
    BTW Test negative numbers
    I HAS A VARIABLE RESULT2 TEH INTEGR ITZ ADD WIT -2 AN WIT 7
    ASSERT WIT RESULT2 SAEM AS 5
    SAYZ WIT "✓ ADD(-2, 7) = 5"
    
    BTW Test zeros
    I HAS A VARIABLE RESULT3 TEH INTEGR ITZ ADD WIT 0 AN WIT 0
    ASSERT WIT RESULT3 SAEM AS 0
    SAYZ WIT "✓ ADD(0, 0) = 0"
    
    SAYZ WIT "ADD function tests passed!"
KTHXBAI

HAI ME TEH FUNCSHUN TEST_IS_EVEN_FUNCTION
    SAYZ WIT "=== Testing IS_EVEN Function ==="
    
    BTW Test even numbers
    ASSERT WIT IS_EVEN WIT 2
    SAYZ WIT "✓ 2 is even"
    
    ASSERT WIT IS_EVEN WIT 0
    SAYZ WIT "✓ 0 is even"
    
    ASSERT WIT IS_EVEN WIT -4
    SAYZ WIT "✓ -4 is even"
    
    BTW Test odd numbers (should be false, so use NOT logic)
    I HAS A VARIABLE IS_3_EVEN TEH BOOL ITZ IS_EVEN WIT 3
    ASSERT WIT IS_3_EVEN SAEM AS NO
    SAYZ WIT "✓ 3 is not even"
    
    I HAS A VARIABLE IS_7_EVEN TEH BOOL ITZ IS_EVEN WIT 7
    ASSERT WIT IS_7_EVEN SAEM AS NO
    SAYZ WIT "✓ 7 is not even"
    
    SAYZ WIT "IS_EVEN function tests passed!"
KTHXBAI
```

## Testing with Collections

### Array Testing

```lol
I CAN HAS TEST?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN TEST_BUKKIT_OPERATIONS
    SAYZ WIT "=== Testing BUKKIT Operations ==="
    
    I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT
    
    BTW Test initial state
    ASSERT WIT ARR SIZ SAEM AS 0
    SAYZ WIT "✓ New BUKKIT has size 0"
    
    BTW Test PUSH operation
    ARR DO PUSH WIT 10
    ASSERT WIT ARR SIZ SAEM AS 1
    SAYZ WIT "✓ After PUSH, size is 1"
    
    ASSERT WIT ARR DO AT WIT 0 SAEM AS 10
    SAYZ WIT "✓ First element is 10"
    
    BTW Test multiple elements
    ARR DO PUSH WIT 20
    ARR DO PUSH WIT 30
    ASSERT WIT ARR SIZ SAEM AS 3
    SAYZ WIT "✓ After adding more elements, size is 3"
    
    BTW Test POP operation
    I HAS A VARIABLE POPPED TEH INTEGR ITZ ARR DO POP
    ASSERT WIT POPPED SAEM AS 30
    SAYZ WIT "✓ POP returned 30"
    
    ASSERT WIT ARR SIZ SAEM AS 2
    SAYZ WIT "✓ After POP, size is 2"
    
    SAYZ WIT "BUKKIT tests passed!"
KTHXBAI
```

### Map Testing

```lol
I CAN HAS TEST?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN TEST_BASKIT_OPERATIONS
    SAYZ WIT "=== Testing BASKIT Operations ==="
    
    I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
    
    BTW Test initial state
    ASSERT WIT MAP SIZ SAEM AS 0
    SAYZ WIT "✓ New BASKIT has size 0"
    
    BTW Test PUT and GET operations
    MAP DO PUT WIT "key1" AN WIT "value1"
    ASSERT WIT MAP SIZ SAEM AS 1
    SAYZ WIT "✓ After PUT, size is 1"
    
    I HAS A VARIABLE VALUE TEH STRIN ITZ MAP DO GET WIT "key1"
    ASSERT WIT VALUE SAEM AS "value1"
    SAYZ WIT "✓ GET returns correct value"
    
    BTW Test CONTAINS operation
    I HAS A VARIABLE HAS_KEY TEH BOOL ITZ MAP DO CONTAINS WIT "key1"
    ASSERT WIT HAS_KEY
    SAYZ WIT "✓ CONTAINS returns true for existing key"
    
    I HAS A VARIABLE NO_KEY TEH BOOL ITZ MAP DO CONTAINS WIT "nonexistent"
    ASSERT WIT NO_KEY SAEM AS NO
    SAYZ WIT "✓ CONTAINS returns false for non-existing key"
    
    SAYZ WIT "BASKIT tests passed!"
KTHXBAI
```

## Test Suite Pattern

### Organized Test Suite

```lol
I CAN HAS TEST?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN RUN_TEST WIT TEST_NAME TEH STRIN AN WIT TEST_FUNCTION TEH STRIN
    SAY WIT "Running test: "
    SAYZ WIT TEST_NAME
    
    MAYB
        BTW Note: In a real implementation, you'd call the test function dynamically
        BTW For this example, we'll show the pattern
        SAYZ WIT "Test execution would happen here"
        SAYZ WIT "✓ Test passed!"
    OOPSIE TEST_ERROR
        SAY WIT "❌ Test failed: "
        SAYZ WIT TEST_ERROR
    KTHX
    
    SAYZ WIT ""
KTHXBAI

HAI ME TEH FUNCSHUN MATH_TEST_SUITE
    SAYZ WIT "======================================="
    SAYZ WIT "         MATH TEST SUITE"
    SAYZ WIT "======================================="
    SAYZ WIT ""
    
    I HAS A VARIABLE TOTAL_TESTS TEH INTEGR ITZ 0
    I HAS A VARIABLE PASSED_TESTS TEH INTEGR ITZ 0
    
    BTW Test arithmetic operations
    TOTAL_TESTS ITZ TOTAL_TESTS MOAR 1
    MAYB
        ASSERT WIT 2 MOAR 3 SAEM AS 5
        SAYZ WIT "✓ Addition test passed"
        PASSED_TESTS ITZ PASSED_TESTS MOAR 1
    OOPSIE ARITH_ERROR
        SAYZ WIT "❌ Addition test failed"
    KTHX
    
    TOTAL_TESTS ITZ TOTAL_TESTS MOAR 1
    MAYB
        ASSERT WIT 10 LES 3 SAEM AS 7
        SAYZ WIT "✓ Subtraction test passed"
        PASSED_TESTS ITZ PASSED_TESTS MOAR 1
    OOPSIE ARITH_ERROR
        SAYZ WIT "❌ Subtraction test failed"
    KTHX
    
    TOTAL_TESTS ITZ TOTAL_TESTS MOAR 1
    MAYB
        ASSERT WIT 4 TIEMZ 5 SAEM AS 20
        SAYZ WIT "✓ Multiplication test passed"
        PASSED_TESTS ITZ PASSED_TESTS MOAR 1
    OOPSIE ARITH_ERROR
        SAYZ WIT "❌ Multiplication test failed"
    KTHX
    
    BTW Report results
    SAYZ WIT ""
    SAYZ WIT "======================================="
    SAY WIT "Tests passed: "
    SAY WIT PASSED_TESTS
    SAY WIT "/"
    SAYZ WIT TOTAL_TESTS
    
    IZ PASSED_TESTS SAEM AS TOTAL_TESTS?
        SAYZ WIT "🎉 ALL TESTS PASSED!"
    NOPE
        SAYZ WIT "❌ Some tests failed"
    KTHX
    SAYZ WIT "======================================="
KTHXBAI
```

## Advanced Testing Patterns

### Testing with Custom Messages

```lol
I CAN HAS TEST?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN ASSERT_WITH_MESSAGE WIT CONDITION TEH BOOL AN WIT MESSAGE TEH STRIN
    MAYB
        ASSERT WIT CONDITION
        SAY WIT "✓ "
        SAYZ WIT MESSAGE
    OOPSIE ASSERT_ERROR
        SAY WIT "❌ "
        SAY WIT MESSAGE
        SAY WIT " - "
        SAYZ WIT ASSERT_ERROR
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN CUSTOM_MESSAGE_TESTS
    SAYZ WIT "=== Custom Message Tests ==="
    
    ASSERT_WITH_MESSAGE WIT YEZ AN WIT "Boolean true should pass"
    ASSERT_WITH_MESSAGE WIT 5 BIGGR THAN 3 AN WIT "5 should be greater than 3"
    ASSERT_WITH_MESSAGE WIT "hello" SAEM AS "hello" AN WIT "String equality should work"
    ASSERT_WITH_MESSAGE WIT NO AN WIT "This test should fail"
    
    SAYZ WIT "Custom message tests completed!"
KTHXBAI
```

### Performance Testing Pattern

```lol
I CAN HAS TEST?
I CAN HAS TIME?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN PERFORMANCE_TEST WIT TEST_NAME TEH STRIN
    SAY WIT "Performance test: "
    SAYZ WIT TEST_NAME
    
    I HAS A VARIABLE START_TIME TEH DATE ITZ NEW DATE
    
    BTW Simulate some work
    I HAS A VARIABLE COUNTER TEH INTEGR ITZ 0
    WHILE COUNTER SMALLR THAN 1000
        COUNTER ITZ COUNTER MOAR 1
    KTHX
    
    I HAS A VARIABLE END_TIME TEH DATE ITZ NEW DATE
    
    BTW Note: In a real implementation, you'd calculate the actual time difference
    SAYZ WIT "✓ Performance test completed"
KTHXBAI
```

## Error Handling in Tests

### Exception Testing

```lol
I CAN HAS TEST?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN TEST_EXPECTED_EXCEPTION
    SAYZ WIT "=== Testing Expected Exceptions ==="
    
    BTW Test that division by zero throws an exception
    I HAS A VARIABLE EXCEPTION_THROWN TEH BOOL ITZ NO
    
    MAYB
        I HAS A VARIABLE RESULT TEH DUBBLE ITZ 10.0 DIVIDEZ 0.0
        BTW If we get here, the exception wasn't thrown
        ASSERT WIT NO  BTW Force failure
    OOPSIE DIV_ERROR
        EXCEPTION_THROWN ITZ YEZ
    KTHX
    
    ASSERT WIT EXCEPTION_THROWN
    SAYZ WIT "✓ Division by zero correctly threw exception"
    
    BTW Test that normal division doesn't throw
    EXCEPTION_THROWN ITZ NO
    
    MAYB
        I HAS A VARIABLE RESULT TEH DUBBLE ITZ 10.0 DIVIDEZ 2.0
        ASSERT WIT RESULT SAEM AS 5.0
        SAYZ WIT "✓ Normal division completed successfully"
    OOPSIE UNEXPECTED_ERROR
        SAYZ WIT "❌ Unexpected error in normal division"
    KTHX
KTHXBAI
```

## Quick Reference

### ASSERT Function

| Function | Parameters | Behavior |
|----------|------------|----------|
| `ASSERT WIT condition` | condition: Any type | Throws "Assertion failed" if condition is falsy |

### Truthiness Reference

| Value | Truthiness |
|-------|------------|
| `YEZ` | Truthy |
| `NO` | Falsy |
| Non-zero numbers | Truthy |
| `0` or `0.0` | Falsy |
| Non-empty strings | Truthy |
| `""` | Falsy |
| `NOTHIN` | Falsy |
| Objects | Truthy |

### Common Test Patterns

```lol
BTW Basic assertion
ASSERT WIT ACTUAL SAEM AS EXPECTED

BTW Testing for inequality
ASSERT WIT ACTUAL SAEM AS EXPECTED SAEM AS NO

BTW Testing ranges
ASSERT WIT VALUE BIGGR THAN MIN_VALUE
ASSERT WIT VALUE SMALLR THAN MAX_VALUE

BTW Testing collections
ASSERT WIT ARRAY SIZ SAEM AS EXPECTED_SIZE
ASSERT WIT MAP DO CONTAINS WIT "key"

BTW Testing exceptions
MAYB
    RISKY_OPERATION
    ASSERT WIT NO  BTW Should not reach here
OOPSIE EXPECTED_ERROR
    BTW Exception was thrown as expected
    ASSERT WIT YEZ
KTHX
```

## Related

- [Exception Handling](../language-guide/control-flow.md#exception-handling) - Working with MAYB/OOPSIE
- [Functions](../language-guide/functions.md) - Creating testable functions
- [Collections](collections.md) - Testing BUKKIT and BASKIT operations