# Exceptions Reference

Complete reference for exception handling, built-in exceptions, and error patterns in Objective-LOL.

## Exception Handling Syntax

### Try-Catch Block

```lol
MAYB
    BTW Code that might throw exceptions
OOPSIE <variable_name>
    BTW Exception handling code
    BTW <variable_name> contains the exception message as STRIN
KTHX
```

### Try-Catch-Finally Block

```lol
MAYB
    BTW Code that might throw exceptions
OOPSIE <variable_name>
    BTW Exception handling code
ALWAYZ
    BTW Code that always runs (cleanup)
KTHX
```

### Try-Finally Block (No Catch)

```lol
MAYB
    BTW Code that might throw exceptions
ALWAYZ
    BTW Cleanup code
KTHX
```

### Throwing Exceptions

```lol
OOPS "Exception message"
```

## Exception Types

All exceptions in Objective-LOL are **string-based**. The exception message is always of type STRIN.

```lol
MAYB
    OOPS "This is a custom exception message"
OOPSIE ERROR_MESSAGE
    BTW ERROR_MESSAGE is type STRIN
    SAYZ WIT ERROR_MESSAGE  BTW "This is a custom exception message"
KTHX
```

## Exception Propagation

Exceptions propagate up the call stack until caught:

```lol
HAI ME TEH FUNCSHUN LEVEL_3
    OOPS "Error from level 3"
KTHXBAI

HAI ME TEH FUNCSHUN LEVEL_2
    LEVEL_3  BTW Exception propagates through here
KTHXBAI

HAI ME TEH FUNCSHUN LEVEL_1
    LEVEL_2  BTW Exception propagates through here
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    MAYB
        LEVEL_1
    OOPSIE PROPAGATED_ERROR
        SAYZ WIT "Caught at top level: "
        SAYZ WIT PROPAGATED_ERROR  BTW "Error from level 3"
    KTHX
KTHXBAI
```

## Exception Handling Patterns

### Resource Cleanup Pattern

```lol
I CAN HAS FILE?

HAI ME TEH FUNCSHUN SAFE_FILE_PROCESSING WIT FILENAME TEH STRIN
    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NOTHIN

    MAYB
        DOC ITZ NEW DOCUMENT WIT FILENAME AN WIT "RW"
        DOC DO OPEN

        BTW Process file
        DOC DO WRITE WIT "Processing..."

    OOPSIE FILE_ERROR
        SAYZ WIT "File error: "
        SAYZ WIT FILE_ERROR

    ALWAYZ
        BTW Always clean up resources
        IZ DOC SAEM AS NOTHIN SAEM AS NO?
            IZ DOC IS_OPEN?
                DOC DO CLOSE
                SAYZ WIT "File closed in cleanup"
            KTHX
        KTHX
    KTHX
KTHXBAI
```

### Retry Pattern

```lol
HAI ME TEH FUNCSHUN RETRY_OPERATION WIT MAX_ATTEMPTS TEH INTEGR
    I HAS A VARIABLE ATTEMPT TEH INTEGR ITZ 1

    WHILE ATTEMPT SMALLR THAN MAX_ATTEMPTS MOAR 1
        MAYB
            BTW Potentially failing operation
            RISKY_OPERATION
            SAYZ WIT "Operation succeeded!"
            GIVEZ UP  BTW Exit function on success

        OOPSIE OPERATION_ERROR
            SAY WIT "Attempt "
            SAY WIT ATTEMPT
            SAY WIT " failed: "
            SAYZ WIT OPERATION_ERROR

            ATTEMPT ITZ ATTEMPT MOAR 1
        KTHX
    KTHX

    SAYZ WIT "All retry attempts failed"
KTHXBAI
```

### Exception Translation Pattern

```lol
HAI ME TEH FUNCSHUN HIGH_LEVEL_OPERATION
    MAYB
        BTW Low-level operations that might fail
        COMPLEX_FILE_PROCESSING

    OOPSIE LOW_LEVEL_ERROR
        BTW Translate low-level errors to high-level ones
        IZ LOW_LEVEL_ERROR SAEM AS "File not found"?
            OOPS "Configuration file is missing"
        NOPE
            IZ LOW_LEVEL_ERROR SAEM AS "Permission denied"?
                OOPS "Insufficient privileges to access configuration"
            NOPE
                BTW Re-throw unknown errors
                OOPS LOW_LEVEL_ERROR
            KTHX
        KTHX
    KTHX
KTHXBAI
```

### Validation Pattern

```lol
HAI ME TEH FUNCSHUN VALIDATE_INPUT WIT VALUE TEH INTEGR
    IZ VALUE SMALLR THAN 0?
        OOPS "Value cannot be negative"
    KTHX

    IZ VALUE BIGGR THAN 100?
        OOPS "Value cannot exceed 100"
    KTHX

    BTW Value is valid
    GIVEZ VALUE
KTHXBAI

HAI ME TEH FUNCSHUN PROCESS_WITH_VALIDATION WIT INPUT TEH INTEGR
    MAYB
        I HAS A VARIABLE VALID_VALUE TEH INTEGR ITZ VALIDATE_INPUT WIT INPUT
        SAYZ WIT "Processing valid value: "
        SAYZ WIT VALID_VALUE

    OOPSIE VALIDATION_ERROR
        SAYZ WIT "Validation failed: "
        SAYZ WIT VALIDATION_ERROR
    KTHX
KTHXBAI
```

### Multiple Exception Types Pattern

```lol
HAI ME TEH FUNCSHUN HANDLE_DIFFERENT_ERRORS
    MAYB
        COMPLEX_OPERATION_THAT_CAN_FAIL_MULTIPLE_WAYS

    OOPSIE ERROR_MSG
        BTW Check error message to determine type
        IZ ERROR_MSG SAEM AS "Division by zero"?
            SAYZ WIT "Mathematical error detected"
        NOPE
            BTW Check if it contains certain keywords
            BTW In a real implementation, you'd need string search functions
            SAYZ WIT "Unknown error: "
            SAYZ WIT ERROR_MSG
        KTHX
    KTHX
KTHXBAI
```

## Exception Reference Quick Guide

### Syntax Summary

```lol
BTW Throw exception
OOPS "message"

BTW Try-catch
MAYB
    code
OOPSIE var
    handler
KTHX

BTW Try-finally
MAYB
    code
ALWAYZ
    cleanup
KTHX

BTW Try-catch-finally
MAYB
    code
OOPSIE var
    handler
ALWAYZ
    cleanup
KTHX
```

## Related

- [Control Flow](../language-guide/control-flow.md) - Exception handling basics
- [TEST Module](../standard-library/test.md) - Testing exception scenarios
- [FILE Module](../standard-library/file.md) - File I/O exceptions