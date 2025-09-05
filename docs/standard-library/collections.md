# Collections - BUKKIT and BASKIT

This guide covers the built-in collection types: BUKKIT (arrays) and BASKIT (maps/dictionaries). These are built-in classes that don't require imports.

## BUKKIT - Dynamic Arrays

BUKKIT is a dynamic array type that can hold any combination of values and provides comprehensive array functionality.

### Creating BUKKIT Arrays

```lol
BTW Create empty array
I HAS A VARIABLE NUMS TEH BUKKIT ITZ NEW BUKKIT
SAYZ WIT NUMS SIZ                    BTW Output: 0

BTW Arrays are dynamic and can hold mixed types
I HAS A VARIABLE MIXED TEH BUKKIT ITZ NEW BUKKIT
MIXED DO PUSH WIT 42
MIXED DO PUSH WIT "hello"
MIXED DO PUSH WIT YEZ
MIXED DO PUSH WIT 3.14
SAYZ WIT MIXED SIZ                   BTW Output: 4
```

### BUKKIT Methods

#### Adding Elements

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT

BTW Add to end
ARR DO PUSH WIT 10
ARR DO PUSH WIT 20
ARR DO PUSH WIT 30
SAYZ WIT ARR SIZ                     BTW 3

BTW Add to beginning
ARR DO UNSHIFT WIT 5
SAYZ WIT ARR SIZ                     BTW 4
SAYZ WIT ARR DO AT WIT 0             BTW 5
```

#### Accessing Elements

```lol
I HAS A VARIABLE NUMBERS TEH BUKKIT ITZ NEW BUKKIT
NUMBERS DO PUSH WIT 100
NUMBERS DO PUSH WIT 200
NUMBERS DO PUSH WIT 300

BTW Get elements by index (0-based)
SAYZ WIT NUMBERS DO AT WIT 0         BTW 100
SAYZ WIT NUMBERS DO AT WIT 1         BTW 200
SAYZ WIT NUMBERS DO AT WIT 2         BTW 300

BTW Modify elements
NUMBERS DO SET WIT 1 AN WIT 999      BTW Set index 1 to 999
SAYZ WIT NUMBERS DO AT WIT 1         BTW 999
```

#### Removing Elements

```lol
I HAS A VARIABLE STACK TEH BUKKIT ITZ NEW BUKKIT
STACK DO PUSH WIT "first"
STACK DO PUSH WIT "second"
STACK DO PUSH WIT "third"

BTW Remove from end (stack behavior)
I HAS A VARIABLE LAST TEH STRIN ITZ STACK DO POP
SAYZ WIT LAST                        BTW "third"
SAYZ WIT STACK SIZ                   BTW 2

BTW Remove from beginning (queue behavior)
I HAS A VARIABLE FIRST TEH STRIN ITZ STACK DO SHIFT
SAYZ WIT FIRST                       BTW "first"
SAYZ WIT STACK SIZ                   BTW 1
```

#### Array Operations

```lol
I HAS A VARIABLE WORDS TEH BUKKIT ITZ NEW BUKKIT
WORDS DO PUSH WIT "objective"
WORDS DO PUSH WIT "lol"
WORDS DO PUSH WIT "programming"

BTW Join elements into string
I HAS A VARIABLE SENTENCE TEH STRIN ITZ WORDS DO JOIN WIT " "
SAYZ WIT SENTENCE                    BTW "objective lol programming"

BTW Sort array in place
WORDS DO SORT
I HAS A VARIABLE SORTED TEH STRIN ITZ WORDS DO JOIN WIT ", "
SAYZ WIT SORTED                      BTW "lol, objective, programming"

BTW Reverse array in place
WORDS DO REVERSE
I HAS A VARIABLE REVERSED TEH STRIN ITZ WORDS DO JOIN WIT " | "
SAYZ WIT REVERSED                    BTW "programming | objective | lol"
```

#### Search Operations

```lol
I HAS A VARIABLE FRUITS TEH BUKKIT ITZ NEW BUKKIT
FRUITS DO PUSH WIT "apple"
FRUITS DO PUSH WIT "banana"
FRUITS DO PUSH WIT "orange"

BTW Check if array contains value
I HAS A VARIABLE HAS_APPLE TEH BOOL ITZ FRUITS DO CONTAINS WIT "apple"
SAYZ WIT HAS_APPLE                   BTW YEZ

BTW Find index of value (-1 if not found)
I HAS A VARIABLE BANANA_INDEX TEH INTEGR ITZ FRUITS DO FIND WIT "banana"
SAYZ WIT BANANA_INDEX               BTW 1

I HAS A VARIABLE GRAPE_INDEX TEH INTEGR ITZ FRUITS DO FIND WIT "grape"
SAYZ WIT GRAPE_INDEX                BTW -1 (not found)
```

#### Array Slicing

```lol
I HAS A VARIABLE NUMBERS TEH BUKKIT ITZ NEW BUKKIT
NUMBERS DO PUSH WIT 1
NUMBERS DO PUSH WIT 2
NUMBERS DO PUSH WIT 3
NUMBERS DO PUSH WIT 4
NUMBERS DO PUSH WIT 5

BTW Create sub-array from index 1 to 4 (exclusive)
I HAS A VARIABLE SLICE TEH BUKKIT ITZ NUMBERS DO SLICE WIT 1 AN WIT 4
SAYZ WIT SLICE DO JOIN WIT ", "     BTW "2, 3, 4"
```

#### Utility Operations

```lol
I HAS A VARIABLE DATA TEH BUKKIT ITZ NEW BUKKIT
DATA DO PUSH WIT "a"
DATA DO PUSH WIT "b"
DATA DO PUSH WIT "c"

BTW Clear all elements
DATA DO CLEAR
SAYZ WIT DATA SIZ                   BTW 0
```

## BASKIT - Maps/Dictionaries

BASKIT is a key-value storage type that provides comprehensive map functionality.

### Creating BASKIT Maps

```lol
BTW Create empty map
I HAS A VARIABLE SCORES TEH BASKIT ITZ NEW BASKIT
SAYZ WIT SCORES SIZ                 BTW 0

BTW Maps can hold any type of values with string keys
I HAS A VARIABLE PLAYER_DATA TEH BASKIT ITZ NEW BASKIT
PLAYER_DATA DO PUT WIT "name" AN WIT "Alice"
PLAYER_DATA DO PUT WIT "score" AN WIT 1500
PLAYER_DATA DO PUT WIT "active" AN WIT YEZ
SAYZ WIT PLAYER_DATA SIZ            BTW 3
```

### BASKIT Methods

#### Adding and Accessing Values

```lol
I HAS A VARIABLE INVENTORY TEH BASKIT ITZ NEW BASKIT

BTW Add key-value pairs
INVENTORY DO PUT WIT "sword" AN WIT 1
INVENTORY DO PUT WIT "potion" AN WIT 5
INVENTORY DO PUT WIT "gold" AN WIT 250

BTW Access values by key
SAYZ WIT INVENTORY DO GET WIT "gold"     BTW 250
SAYZ WIT INVENTORY DO GET WIT "potion"   BTW 5

BTW Update existing key
INVENTORY DO PUT WIT "gold" AN WIT 300
SAYZ WIT INVENTORY DO GET WIT "gold"     BTW 300
```

#### Key Operations

```lol
I HAS A VARIABLE CONFIG TEH BASKIT ITZ NEW BASKIT
CONFIG DO PUT WIT "debug" AN WIT YEZ
CONFIG DO PUT WIT "port" AN WIT 8080
CONFIG DO PUT WIT "host" AN WIT "localhost"

BTW Check if key exists
I HAS A VARIABLE HAS_DEBUG TEH BOOL ITZ CONFIG DO CONTAINS WIT "debug"
SAYZ WIT HAS_DEBUG                      BTW YEZ

I HAS A VARIABLE HAS_SSL TEH BOOL ITZ CONFIG DO CONTAINS WIT "ssl"
SAYZ WIT HAS_SSL                        BTW NO

BTW Remove key-value pair (returns the removed value)
I HAS A VARIABLE OLD_PORT TEH INTEGR ITZ CONFIG DO REMOVE WIT "port"
SAYZ WIT OLD_PORT                       BTW 8080
SAYZ WIT CONFIG SIZ                     BTW 2
```

#### Getting Collections of Keys and Values

```lol
I HAS A VARIABLE PERSON TEH BASKIT ITZ NEW BASKIT
PERSON DO PUT WIT "firstName" AN WIT "John"
PERSON DO PUT WIT "lastName" AN WIT "Doe"
PERSON DO PUT WIT "age" AN WIT 30

BTW Get all keys as BUKKIT
I HAS A VARIABLE ALL_KEYS TEH BUKKIT ITZ PERSON DO KEYS
SAYZ WIT ALL_KEYS DO JOIN WIT ", "      BTW "firstName, lastName, age"

BTW Get all values as BUKKIT
I HAS A VARIABLE ALL_VALUES TEH BUKKIT ITZ PERSON DO VALUES
SAYZ WIT ALL_VALUES DO JOIN WIT " | "   BTW "John | Doe | 30"

BTW Get key-value pairs as BUKKIT of BUKKIT
I HAS A VARIABLE PAIRS TEH BUKKIT ITZ PERSON DO PAIRS
BTW Each pair is a BUKKIT with [key, value]
```

#### Map Operations

```lol
I HAS A VARIABLE DEFAULTS TEH BASKIT ITZ NEW BASKIT
DEFAULTS DO PUT WIT "theme" AN WIT "dark"
DEFAULTS DO PUT WIT "lang" AN WIT "en"

I HAS A VARIABLE USER_PREFS TEH BASKIT ITZ NEW BASKIT
USER_PREFS DO PUT WIT "theme" AN WIT "light"
USER_PREFS DO PUT WIT "fontSize" AN WIT 14

BTW Merge another BASKIT into this one
USER_PREFS DO MERGE WIT DEFAULTS
SAYZ WIT USER_PREFS SIZ                 BTW 3
SAYZ WIT USER_PREFS DO GET WIT "theme"  BTW "light" (not overwritten)
SAYZ WIT USER_PREFS DO GET WIT "lang"   BTW "en" (added from defaults)

BTW Create a copy
I HAS A VARIABLE BACKUP TEH BASKIT ITZ USER_PREFS DO COPY
SAYZ WIT BACKUP SIZ                     BTW 3

BTW Clear all entries
USER_PREFS DO CLEAR
SAYZ WIT USER_PREFS SIZ                 BTW 0
SAYZ WIT BACKUP SIZ                     BTW 3 (copy unaffected)
```

## Collection Usage Examples

### Data Processing with BUKKIT

```lol
I CAN HAS STDIO?
I CAN HAS MATH?

HAI ME TEH FUNCSHUN ANALYZE_SCORES WIT SCORES TEH BUKKIT
    BTW Find min, max, and average
    I HAS A VARIABLE MIN_SCORE TEH INTEGR ITZ SCORES DO AT WIT 0
    I HAS A VARIABLE MAX_SCORE TEH INTEGR ITZ SCORES DO AT WIT 0
    I HAS A VARIABLE TOTAL TEH INTEGR ITZ 0
    I HAS A VARIABLE INDEX TEH INTEGR ITZ 0

    WHILE INDEX SMALLR THAN SCORES SIZ
        I HAS A VARIABLE CURRENT TEH INTEGR ITZ SCORES DO AT WIT INDEX
        TOTAL ITZ TOTAL MOAR CURRENT
        MIN_SCORE ITZ MIN WIT MIN_SCORE AN WIT CURRENT AS DUBBLE AS INTEGR
        MAX_SCORE ITZ MAX WIT MAX_SCORE AN WIT CURRENT AS DUBBLE AS INTEGR
        INDEX ITZ INDEX MOAR 1
    KTHX

    I HAS A VARIABLE AVERAGE TEH INTEGR ITZ TOTAL DIVIDEZ SCORES SIZ AS INTEGR

    SAY WIT "Min: "
    SAY WIT MIN_SCORE
    SAY WIT ", Max: "
    SAY WIT MAX_SCORE
    SAY WIT ", Average: "
    SAYZ WIT AVERAGE
KTHXBAI
```

### Configuration Management with BASKIT

```lol
HAI ME TEH FUNCSHUN LOAD_CONFIG TEH BASKIT
    I HAS A VARIABLE CONFIG TEH BASKIT ITZ NEW BASKIT

    BTW Set default configuration
    CONFIG DO PUT WIT "server_port" AN WIT 3000
    CONFIG DO PUT WIT "debug_mode" AN WIT NO
    CONFIG DO PUT WIT "max_connections" AN WIT 100
    CONFIG DO PUT WIT "log_level" AN WIT "info"

    GIVEZ CONFIG
KTHXBAI

HAI ME TEH FUNCSHUN PRINT_CONFIG WIT CONFIG TEH BASKIT
    I HAS A VARIABLE KEYS TEH BUKKIT ITZ CONFIG DO KEYS
    I HAS A VARIABLE INDEX TEH INTEGR ITZ 0

    SAYZ WIT "=== Configuration ==="

    WHILE INDEX SMALLR THAN KEYS SIZ
        I HAS A VARIABLE KEY TEH STRIN ITZ KEYS DO AT WIT INDEX
        I HAS A VARIABLE VALUE TEH STRIN ITZ CONFIG DO GET WIT KEY AS STRIN
        SAY WIT KEY
        SAY WIT ": "
        SAYZ WIT VALUE
        INDEX ITZ INDEX MOAR 1
    KTHX
KTHXBAI
```

## Exception Handling

Both BUKKIT and BASKIT throw exceptions for invalid operations:

```lol
MAYB
    I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT
    I HAS A VARIABLE ITEM TEH INTEGR ITZ ARR DO AT WIT 5  BTW Index out of bounds
OOPSIE BOUNDS_ERROR
    SAYZ WIT BOUNDS_ERROR  BTW "Array index 5 out of bounds (size 0)"
KTHX

MAYB
    I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
    I HAS A VARIABLE VALUE TEH STRIN ITZ MAP DO GET WIT "nonexistent"
OOPSIE KEY_ERROR
    SAYZ WIT KEY_ERROR     BTW "Key 'nonexistent' not found in map"
KTHX
```

## Quick Reference

### BUKKIT Methods

| Method | Purpose | Returns |
|--------|---------|---------|
| `PUSH WIT value` | Add to end | New size |
| `POP` | Remove from end | Removed value |
| `UNSHIFT WIT value` | Add to beginning | New size |
| `SHIFT` | Remove from beginning | Removed value |
| `AT WIT index` | Get element | Element value |
| `SET WIT index AN WIT value` | Set element | NOTHIN |
| `JOIN WIT separator` | Join to string | STRIN |
| `SORT` | Sort in place | NOTHIN |
| `REVERSE` | Reverse in place | NOTHIN |
| `CONTAINS WIT value` | Check if contains | BOOL |
| `FIND WIT value` | Find index | INTEGR (-1 if not found) |
| `SLICE WIT start AN WIT end` | Create sub-array | BUKKIT |
| `CLEAR` | Remove all | NOTHIN |

### BASKIT Methods

| Method | Purpose | Returns |
|--------|---------|---------|
| `PUT WIT key AN WIT value` | Set key-value | NOTHIN |
| `GET WIT key` | Get value by key | Value (throws if not found) |
| `CONTAINS WIT key` | Check if key exists | BOOL |
| `REMOVE WIT key` | Remove key-value | Removed value |
| `KEYS` | Get all keys | BUKKIT |
| `VALUES` | Get all values | BUKKIT |
| `PAIRS` | Get key-value pairs | BUKKIT |
| `MERGE WIT other` | Merge another BASKIT | NOTHIN |
| `COPY` | Create shallow copy | BASKIT |
| `CLEAR` | Remove all entries | NOTHIN |

## Related

- [Examples](../examples/data-processing.md) - Complex collection usage examples
- [STDIO Module](stdio.md) - For displaying collection data