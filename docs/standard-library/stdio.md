# STDIO Module

## Import

```lol
BTW Full import
I CAN HAS STDIO?

BTW Selective import examples
I CAN HAS SAY FROM STDIO?
I CAN HAS SAYZ FROM STDIO?
```

## Output

### SAY

Prints a value to standard output without a newline.
Accepts any type and converts it to STRIN representation.

**Syntax:** `SAY WIT <value>`
**Returns:** 

**Parameters:**
- `value` (ANY): Any value to print (INTEGR, DUBBLE, STRIN, BOOL, etc.)

**Example: Print string without newline**

```lol
SAY WIT "Hello "
SAY WIT "World"
SAY WIT "!"
BTW Output: Hello World!
```

**Example: Print numbers**

```lol
SAY WIT 42
SAY WIT " is the answer"
BTW Output: 42 is the answer
```

**Example: Print boolean values**

```lol
SAY WIT YEZ
SAY WIT " and "
SAY WIT NO
BTW Output: YEZ and NO
```

**Note:** Does not add a newline character

**Note:** Accepts any type and converts to string representation

**See also:** SAYZ

### SAYZ

Prints a value to standard output followed by a newline.
Accepts any type and converts it to STRIN representation.

**Syntax:** `SAYZ WIT <value>`
**Returns:** 

**Parameters:**
- `value` (ANY): Any value to print (INTEGR, DUBBLE, STRIN, BOOL, etc.)

**Example: Print lines of text**

```lol
SAYZ WIT "First line"
SAYZ WIT "Second line"
SAYZ WIT 42
BTW Output:
BTW First line
BTW Second line
BTW 42
```

**Example: Print variables**

```lol
I HAS A VARIABLE NAME TEH STRIN ITZ "Alice"
I HAS A VARIABLE AGE TEH INTEGR ITZ 25
SAYZ WIT NAME
SAYZ WIT AGE
BTW Output:
BTW Alice
BTW 25
```

**Note:** Automatically adds a newline character

**Note:** Accepts any type and converts to string representation

**See also:** SAY

## Input

### GIMME

Reads a line of input from standard input.
Returns the input as a STRIN with trailing newline removed.

**Syntax:** `GIMME`
**Returns:** STRIN

**Example: Read user input**

```lol
SAYZ WIT "Enter your name: "
I HAS A VARIABLE USER_NAME TEH STRIN ITZ GIMME
SAY WIT "Hello, "
SAYZ WIT USER_NAME
BTW If user enters "Alice", output: Hello, Alice
```

**Example: Interactive calculator**

```lol
SAYZ WIT "Enter first number: "
I HAS A VARIABLE NUM1_STR TEH STRIN ITZ GIMME
I HAS A VARIABLE NUM1 TEH INTEGR ITZ NUM1_STR AS INTEGR
SAYZ WIT "Enter second number: "
I HAS A VARIABLE NUM2_STR TEH STRIN ITZ GIMME
I HAS A VARIABLE NUM2 TEH INTEGR ITZ NUM2_STR AS INTEGR
I HAS A VARIABLE SUM TEH INTEGR ITZ NUM1 MOAR NUM2
SAY WIT "Sum: "
SAYZ WIT SUM
```

**Example: Simple quiz**

```lol
SAYZ WIT "What is 5 + 3?"
I HAS A VARIABLE ANSWER TEH STRIN ITZ GIMME
IZ ANSWER SAEM AS "8"?
SAYZ WIT "Correct!"
NOPE
SAYZ WIT "Wrong! The answer is 8."
KTHX
```

**Note:** Waits for user to press Enter

**Note:** Removes trailing newline and carriage return characters

**Note:** Returns empty string on EOF

**See also:** SAY, SAYZ

