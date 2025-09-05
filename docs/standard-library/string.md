# STRING Module - String Utility Functions

The STRING module provides essential string manipulation and analysis functions for working with text data.

## Importing STRING Module

```lol
BTW Import entire module
I CAN HAS STRING?

BTW Selective imports
I CAN HAS LEN FROM STRING?
I CAN HAS CONCAT FROM STRING?
```

## String Functions

### LEN - String Length

Returns the length of a string (number of characters).

**Syntax:**
```lol
LEN WIT <string>
```

**Parameters:**
- **string**: STRIN - The string to measure

**Returns:** INTEGR - The length of the string

**Examples:**
```lol
I CAN HAS STRING?
I CAN HAS STDIO?

BTW Basic string length
I HAS A VARIABLE MESSAGE TEH STRIN ITZ "Hello, World!"
I HAS A VARIABLE LENGTH TEH INTEGR ITZ LEN WIT MESSAGE
SAYZ WIT LENGTH                    BTW Output: 13

BTW Empty string
I HAS A VARIABLE EMPTY TEH STRIN ITZ ""
SAYZ WIT LEN WIT EMPTY            BTW Output: 0

BTW String with special characters
I HAS A VARIABLE SPECIAL TEH STRIN ITZ "Line 1\nLine 2\t\tTabbed"
SAYZ WIT LEN WIT SPECIAL          BTW Output: 19
```

### CONCAT - String Concatenation

Combines two strings into a single string.

**Syntax:**
```lol
CONCAT WIT <string1> AN WIT <string2>
```

**Parameters:**
- **string1**: STRIN - The first string
- **string2**: STRIN - The second string

**Returns:** STRIN - The concatenated string

**Examples:**
```lol
I CAN HAS STRING?
I CAN HAS STDIO?

BTW Basic concatenation
I HAS A VARIABLE FIRST TEH STRIN ITZ "Hello"
I HAS A VARIABLE SECOND TEH STRIN ITZ " World"
I HAS A VARIABLE RESULT TEH STRIN ITZ CONCAT WIT FIRST AN WIT SECOND
SAYZ WIT RESULT                   BTW Output: Hello World

BTW Building longer strings
I HAS A VARIABLE NAME TEH STRIN ITZ "Alice"
I HAS A VARIABLE GREETING TEH STRIN ITZ CONCAT WIT "Hello, " AN WIT NAME
I HAS A VARIABLE FULL_MSG TEH STRIN ITZ CONCAT WIT GREETING AN WIT "!"
SAYZ WIT FULL_MSG                 BTW Output: Hello, Alice!

BTW Concatenating with empty strings
I HAS A VARIABLE TEXT TEH STRIN ITZ "Important"
I HAS A VARIABLE WITH_EMPTY TEH STRIN ITZ CONCAT WIT TEXT AN WIT ""
SAYZ WIT WITH_EMPTY               BTW Output: Important
```

## Practical Examples

### String Builder Pattern

```lol
I CAN HAS STRING?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN BUILD_MESSAGE WIT NAME TEH STRIN AN WIT AGE TEH INTEGR
    BTW Convert age to string for concatenation
    I HAS A VARIABLE AGE_STR TEH STRIN ITZ AGE AS STRIN
    
    BTW Build message step by step
    I HAS A VARIABLE MSG TEH STRIN ITZ CONCAT WIT "Name: " AN WIT NAME
    MSG ITZ CONCAT WIT MSG AN WIT ", Age: "
    MSG ITZ CONCAT WIT MSG AN WIT AGE_STR
    MSG ITZ CONCAT WIT MSG AN WIT " years old"
    
    GIVEZ MSG
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE PERSON_INFO TEH STRIN ITZ BUILD_MESSAGE WIT "Bob" AN WIT 25
    SAYZ WIT PERSON_INFO          BTW Output: Name: Bob, Age: 25 years old
KTHXBAI
```

### Text Analysis

```lol
I CAN HAS STRING?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN ANALYZE_TEXT WIT TEXT TEH STRIN
    I HAS A VARIABLE LENGTH TEH INTEGR ITZ LEN WIT TEXT
    
    SAY WIT "Text: \""
    SAY WIT TEXT
    SAYZ WIT "\""
    
    SAY WIT "Length: "
    SAYZ WIT LENGTH
    
    BTW Categorize by length
    IZ LENGTH SAEM AS 0?
        SAYZ WIT "Category: Empty string"
    NOPE
        IZ LENGTH SMALLR THAN 10?
            SAYZ WIT "Category: Short text"
        NOPE
            IZ LENGTH SMALLR THAN 50?
                SAYZ WIT "Category: Medium text"
            NOPE
                SAYZ WIT "Category: Long text"
            KTHX
        KTHX
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    ANALYZE_TEXT WIT ""
    ANALYZE_TEXT WIT "Hi"
    ANALYZE_TEXT WIT "This is a medium length sentence."
    ANALYZE_TEXT WIT "This is a much longer piece of text that contains many words and would be categorized as long text."
KTHXBAI
```

### String Comparison Helper

```lol
I CAN HAS STRING?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN COMPARE_STRINGS WIT STR1 TEH STRIN AN WIT STR2 TEH STRIN
    I HAS A VARIABLE LEN1 TEH INTEGR ITZ LEN WIT STR1
    I HAS A VARIABLE LEN2 TEH INTEGR ITZ LEN WIT STR2
    
    SAY WIT "String 1: \""
    SAY WIT STR1
    SAY WIT "\" (length: "
    SAY WIT LEN1
    SAYZ WIT ")"
    
    SAY WIT "String 2: \""
    SAY WIT STR2
    SAY WIT "\" (length: "
    SAY WIT LEN2
    SAYZ WIT ")"
    
    IZ STR1 SAEM AS STR2?
        SAYZ WIT "Result: Strings are identical"
    NOPE
        IZ LEN1 SAEM AS LEN2?
            SAYZ WIT "Result: Different content, same length"
        NOPE
            IZ LEN1 BIGGR THAN LEN2?
                SAYZ WIT "Result: First string is longer"
            NOPE
                SAYZ WIT "Result: Second string is longer"
            KTHX
        KTHX
    KTHX
    SAYZ WIT ""
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    COMPARE_STRINGS WIT "hello" AN WIT "hello"
    COMPARE_STRINGS WIT "cat" AN WIT "dog"
    COMPARE_STRINGS WIT "short" AN WIT "much longer string"
KTHXBAI
```

### Data Formatter

```lol
I CAN HAS STRING?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN FORMAT_TABLE_ROW WIT COL1 TEH STRIN AN WIT COL2 TEH STRIN AN WIT COL3 TEH STRIN
    BTW Create table separator
    I HAS A VARIABLE SEPARATOR TEH STRIN ITZ " | "
    
    BTW Build the row
    I HAS A VARIABLE ROW TEH STRIN ITZ CONCAT WIT COL1 AN WIT SEPARATOR
    ROW ITZ CONCAT WIT ROW AN WIT COL2
    ROW ITZ CONCAT WIT ROW AN WIT SEPARATOR
    ROW ITZ CONCAT WIT ROW AN WIT COL3
    
    GIVEZ ROW
KTHXBAI

HAI ME TEH FUNCSHUN CREATE_REPORT
    SAYZ WIT "Employee Report"
    SAYZ WIT "==============="
    
    BTW Header
    I HAS A VARIABLE HEADER TEH STRIN ITZ FORMAT_TABLE_ROW WIT "Name" AN WIT "Department" AN WIT "Years"
    SAYZ WIT HEADER
    
    BTW Data rows
    I HAS A VARIABLE ROW1 TEH STRIN ITZ FORMAT_TABLE_ROW WIT "Alice" AN WIT "Engineering" AN WIT "5"
    SAYZ WIT ROW1
    
    I HAS A VARIABLE ROW2 TEH STRIN ITZ FORMAT_TABLE_ROW WIT "Bob" AN WIT "Marketing" AN WIT "3"
    SAYZ WIT ROW2
    
    I HAS A VARIABLE ROW3 TEH STRIN ITZ FORMAT_TABLE_ROW WIT "Carol" AN WIT "Design" AN WIT "7"
    SAYZ WIT ROW3
KTHXBAI
```

### Password Validator

```lol
I CAN HAS STRING?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN VALIDATE_PASSWORD WIT PASSWORD TEH STRIN
    I HAS A VARIABLE LENGTH TEH INTEGR ITZ LEN WIT PASSWORD
    I HAS A VARIABLE IS_VALID TEH BOOL ITZ YEZ
    
    SAY WIT "Validating password: "
    SAY WIT PASSWORD
    SAY WIT " (length: "
    SAY WIT LENGTH
    SAYZ WIT ")"
    
    BTW Check minimum length
    IZ LENGTH SMALLR THAN 8?
        SAYZ WIT "❌ Password must be at least 8 characters"
        IS_VALID ITZ NO
    NOPE
        SAYZ WIT "✓ Length requirement met"
    KTHX
    
    BTW Check maximum length
    IZ LENGTH BIGGR THAN 128?
        SAYZ WIT "❌ Password must be no more than 128 characters"
        IS_VALID ITZ NO
    NOPE
        SAYZ WIT "✓ Length within limits"
    KTHX
    
    BTW Overall result
    IZ IS_VALID?
        SAYZ WIT "✓ Password is valid!"
    NOPE
        SAYZ WIT "❌ Password validation failed"
    KTHX
    
    SAYZ WIT ""
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    VALIDATE_PASSWORD WIT "short"
    VALIDATE_PASSWORD WIT "this_is_a_good_password"
    VALIDATE_PASSWORD WIT ""
KTHXBAI
```

## Working with Built-in String Operations

While the STRING module provides specialized functions, remember that Objective-LOL also has built-in string operations:

```lol
I CAN HAS STRING?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEMONSTRATE_STRING_OPERATIONS
    I HAS A VARIABLE TEXT1 TEH STRIN ITZ "Hello"
    I HAS A VARIABLE TEXT2 TEH STRIN ITZ "World"
    
    BTW Built-in concatenation with MOAR operator
    I HAS A VARIABLE COMBINED1 TEH STRIN ITZ TEXT1 MOAR " " MOAR TEXT2
    SAYZ WIT COMBINED1
    
    BTW STRING module concatenation
    I HAS A VARIABLE COMBINED2 TEH STRIN ITZ CONCAT WIT TEXT1 AN WIT " "
    COMBINED2 ITZ CONCAT WIT COMBINED2 AN WIT TEXT2
    SAYZ WIT COMBINED2
    
    BTW Both produce the same result
    IZ COMBINED1 SAEM AS COMBINED2?
        SAYZ WIT "Both methods produce identical results!"
    KTHX
KTHXBAI
```

## Error Handling

STRING functions will throw exceptions for invalid argument types:

```lol
I CAN HAS STRING?

MAYB
    BTW This will cause an error - passing integer instead of string
    I HAS A VARIABLE LENGTH TEH INTEGR ITZ LEN WIT 12345
OOPSIE TYPE_ERROR
    SAYZ WIT "Error: "
    SAYZ WIT TYPE_ERROR
KTHX

MAYB
    BTW This will cause an error - missing second argument
    I HAS A VARIABLE RESULT TEH STRIN ITZ CONCAT WIT "Hello" AN WIT 123
OOPSIE CONCAT_ERROR
    SAYZ WIT "Concatenation error: "
    SAYZ WIT CONCAT_ERROR
KTHX
```

## Quick Reference

### Functions

| Function | Parameters | Return Type | Description |
|----------|------------|-------------|-------------|
| `LEN WIT string` | string: STRIN | INTEGR | Get string length |
| `CONCAT WIT str1 AN WIT str2` | str1: STRIN, str2: STRIN | STRIN | Concatenate strings |

### Common Patterns

```lol
BTW Check if string is empty
IZ LEN WIT MY_STRING SAEM AS 0?
    SAYZ WIT "String is empty"
KTHX

BTW Build a sentence
I HAS A VARIABLE SENTENCE TEH STRIN ITZ CONCAT WIT SUBJECT AN WIT " "
SENTENCE ITZ CONCAT WIT SENTENCE AN WIT VERB
SENTENCE ITZ CONCAT WIT SENTENCE AN WIT " "
SENTENCE ITZ CONCAT WIT SENTENCE AN WIT OBJECT

BTW Length-based validation
IZ LEN WIT PASSWORD BIGGR THAN 7?
    SAYZ WIT "Password meets minimum length"
NOPE
    SAYZ WIT "Password too short"
KTHX
```

## Related

- [STDIO Module](stdio.md) - For string input/output operations
- [Syntax Basics](../language-guide/syntax-basics.md) - Built-in string operations
- [Collections](collections.md) - Working with BUKKIT arrays of strings