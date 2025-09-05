# RANDOM Module - Random Number Generation

The RANDOM module provides comprehensive random number generation functionality for various data types and use cases.

## Import

```lol
BTW Full import
I CAN HAS RANDOM?

BTW Selective import examples
I CAN HAS RANDOM_FLOAT FROM RANDOM?
I CAN HAS RANDOM_INT FROM RANDOM?
I CAN HAS SEED_TIME AN UUID FROM RANDOM?
I CAN HAS RANDOM_BOOL AN RANDOM_STRING FROM RANDOM?
```

## Seeding Functions

### SEED - Set Random Seed

Sets the random number generator seed to a specific value for reproducible results.

**Syntax:** `SEED WIT <seed_value>`
**Parameters:** seed_value (INTEGR)

```lol
I CAN HAS SEED FROM RANDOM?

SEED WIT 12345    BTW Set specific seed
BTW All subsequent random calls will be deterministic
```

### SEED_TIME - Seed with Current Time

Seeds the random number generator with the current time for non-deterministic results.

**Syntax:** `SEED_TIME`

```lol
I CAN HAS SEED_TIME FROM RANDOM?

SEED_TIME    BTW Seed with current timestamp
BTW All subsequent random calls will vary between runs
```

## Basic Random Functions

### RANDOM_FLOAT - Random Float

Returns a random floating-point number between 0 (inclusive) and 1 (exclusive).

**Syntax:** `RANDOM_FLOAT`
**Returns:** DUBBLE

```lol
I CAN HAS RANDOM_FLOAT FROM RANDOM?

I HAS A VARIABLE RAND1 TEH DUBBLE ITZ RANDOM_FLOAT    BTW e.g., 0.7234
I HAS A VARIABLE RAND2 TEH DUBBLE ITZ RANDOM_FLOAT    BTW e.g., 0.1849

BTW Generate random number in range [0, 100)
I HAS A VARIABLE PERCENT TEH DUBBLE ITZ RANDOM_FLOAT TIEMZ 100
```

### RANDOM_INT - Random Integer

Returns a random integer in the range [min, max).

**Syntax:** `RANDOM_INT WIT <min> AN WIT <max>`
**Returns:** INTEGR

```lol
I CAN HAS RANDOM_INT FROM RANDOM?

BTW Dice roll (1-6)
I HAS A VARIABLE DICE TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 7

BTW Random index for array of size 10
I HAS A VARIABLE INDEX TEH INTEGR ITZ RANDOM_INT WIT 0 AN WIT 10

BTW Random number between -10 and 10 (inclusive)
I HAS A VARIABLE SIGNED TEH INTEGR ITZ RANDOM_INT WIT -10 AN WIT 11
```

### RANDOM_BOOL - Random Boolean

Returns a random boolean value.

**Syntax:** `RANDOM_BOOL`
**Returns:** BOOL

```lol
I CAN HAS RANDOM_BOOL FROM RANDOM?

I HAS A VARIABLE COIN_FLIP TEH BOOL ITZ RANDOM_BOOL    BTW YEZ or NO

IZ COIN_FLIP?
    SAYZ WIT "Heads!"
NOPE
    SAYZ WIT "Tails!"
KTHX
```

## String Generation

### RANDOM_STRING - Random String

Generates a random string of specified length using characters from a given charset.

**Syntax:** `RANDOM_STRING WIT <length> AN WIT <charset>`
**Returns:** STRIN

```lol
I CAN HAS RANDOM_STRING FROM RANDOM?

BTW Generate random alphanumeric string
I HAS A VARIABLE CHARSET TEH STRIN ITZ "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
I HAS A VARIABLE RANDOM_ID TEH STRIN ITZ RANDOM_STRING WIT 8 AN WIT CHARSET

BTW Generate random numeric string
I HAS A VARIABLE NUMERIC_CHARSET TEH STRIN ITZ "0123456789"
I HAS A VARIABLE RANDOM_CODE TEH STRIN ITZ RANDOM_STRING WIT 6 AN WIT NUMERIC_CHARSET

BTW Generate random password-like string
I HAS A VARIABLE PASS_CHARSET TEH STRIN ITZ "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"
I HAS A VARIABLE PASSWORD TEH STRIN ITZ RANDOM_STRING WIT 12 AN WIT PASS_CHARSET
```

### UUID - Generate UUID

Generates a universally unique identifier (UUID) string.

**Syntax:** `UUID`
**Returns:** STRIN

```lol
I CAN HAS UUID FROM RANDOM?

I HAS A VARIABLE UNIQUE_ID TEH STRIN ITZ UUID
SAYZ WIT UNIQUE_ID    BTW e.g., "f47ac10b-58cc-4372-a567-0e02b2c3d479"

BTW Use UUID for unique identifiers
I HAS A VARIABLE SESSION_ID TEH STRIN ITZ UUID
I HAS A VARIABLE REQUEST_ID TEH STRIN ITZ UUID
```

## Collection Functions

### RANDOM_CHOICE - Random Array Element

Selects a random element from a BUKKIT array.

**Syntax:** `RANDOM_CHOICE WIT <array>`
**Returns:** Value from array

```lol
I CAN HAS RANDOM_CHOICE FROM RANDOM?

I HAS A VARIABLE COLORS TEH BUKKIT ITZ NEW BUKKIT
COLORS DO PUSH WIT "red"
COLORS DO PUSH WIT "blue"
COLORS DO PUSH WIT "green"
COLORS DO PUSH WIT "yellow"

I HAS A VARIABLE CHOSEN_COLOR TEH STRIN ITZ RANDOM_CHOICE WIT COLORS
SAYZ WIT CHOSEN_COLOR    BTW One of the colors
```

### SHUFFLE - Shuffle Array

Randomly rearranges the elements in a BUKKIT array in-place.

**Syntax:** `SHUFFLE WIT <array>`

```lol
I CAN HAS SHUFFLE FROM RANDOM?

I HAS A VARIABLE DECK TEH BUKKIT ITZ NEW BUKKIT
DECK DO PUSH WIT "Ace"
DECK DO PUSH WIT "King"
DECK DO PUSH WIT "Queen"
DECK DO PUSH WIT "Jack"

SHUFFLE WIT DECK    BTW Randomly reorder the elements
BTW DECK now contains the same elements in random order
```

## Usage Examples

### Dice Rolling Simulator

```lol
I CAN HAS RANDOM?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN ROLL_DICE WIT SIDES TEH INTEGR
    GIVEZ RANDOM_INT WIT 1 AN WIT SIDES MOAR 1
KTHXBAI

HAI ME TEH FUNCSHUN SIMULATE_DICE_ROLLS WIT NUM_ROLLS TEH INTEGR
    SEED_TIME    BTW Ensure different results each time
    
    I HAS A VARIABLE TOTAL TEH INTEGR ITZ 0
    I HAS A VARIABLE ROLL_COUNT TEH INTEGR ITZ 0

    WHILE ROLL_COUNT SMALLR THAN NUM_ROLLS
        I HAS A VARIABLE ROLL TEH INTEGR ITZ ROLL_DICE WIT 6
        TOTAL ITZ TOTAL MOAR ROLL
        SAY WIT "Roll "
        SAY WIT ROLL_COUNT MOAR 1
        SAY WIT ": "
        SAYZ WIT ROLL
        ROLL_COUNT ITZ ROLL_COUNT MOAR 1
    KTHX

    SAY WIT "Average: "
    SAYZ WIT TOTAL DIVIDEZ NUM_ROLLS AS INTEGR
KTHXBAI
```

### Password Generator

```lol
I CAN HAS RANDOM?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN GENERATE_PASSWORD WIT LENGTH TEH INTEGR
    I HAS A VARIABLE LOWERCASE TEH STRIN ITZ "abcdefghijklmnopqrstuvwxyz"
    I HAS A VARIABLE UPPERCASE TEH STRIN ITZ "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    I HAS A VARIABLE NUMBERS TEH STRIN ITZ "0123456789"
    I HAS A VARIABLE SYMBOLS TEH STRIN ITZ "!@#$%^&*()_+-=[]{}|;:,.<>?"
    
    BTW Combine all character sets
    I HAS A VARIABLE CHARSET TEH STRIN ITZ LOWERCASE MOAR UPPERCASE MOAR NUMBERS MOAR SYMBOLS
    
    SEED_TIME
    GIVEZ RANDOM_STRING WIT LENGTH AN WIT CHARSET
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE PASSWORD TEH STRIN ITZ GENERATE_PASSWORD WIT 16
    SAY WIT "Generated password: "
    SAYZ WIT PASSWORD
KTHXBAI
```

### Random Selection Game

```lol
I CAN HAS RANDOM?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN RANDOM_GAME
    SEED_TIME
    
    I HAS A VARIABLE PRIZES TEH BUKKIT ITZ NEW BUKKIT
    PRIZES DO PUSH WIT "Car"
    PRIZES DO PUSH WIT "Vacation"
    PRIZES DO PUSH WIT "Cash"
    PRIZES DO PUSH WIT "Electronics"
    PRIZES DO PUSH WIT "Gift Card"
    
    SAY WIT "You won: "
    I HAS A VARIABLE PRIZE TEH STRIN ITZ RANDOM_CHOICE WIT PRIZES
    SAYZ WIT PRIZE
    
    BTW Show probability
    I HAS A VARIABLE CHANCE TEH DUBBLE ITZ RANDOM_FLOAT TIEMZ 100
    SAY WIT "Lucky number: "
    SAYZ WIT ROUND WIT CHANCE
KTHXBAI
```

### Statistical Sampling

```lol
I CAN HAS RANDOM?
I CAN HAS STDIO?
I CAN HAS MATH?

HAI ME TEH FUNCSHUN SAMPLE_NORMAL_DISTRIBUTION WIT COUNT TEH INTEGR
    SEED_TIME
    
    I HAS A VARIABLE SAMPLES TEH BUKKIT ITZ NEW BUKKIT
    I HAS A VARIABLE I TEH INTEGR ITZ 0
    
    WHILE I SMALLR THAN COUNT
        BTW Box-Muller transform for normal distribution
        I HAS A VARIABLE U1 TEH DUBBLE ITZ RANDOM_FLOAT
        I HAS A VARIABLE U2 TEH DUBBLE ITZ RANDOM_FLOAT
        
        I HAS A VARIABLE Z TEH DUBBLE ITZ SQRT WIT -2.0 TIEMZ LOG WIT U1
        Z ITZ Z TIEMZ COS WIT 2.0 TIEMZ PI TIEMZ U2
        
        SAMPLES DO PUSH WIT Z
        I ITZ I MOAR 1
    KTHX
    
    GIVEZ SAMPLES
KTHXBAI
```

## Function Summary

### Seeding Functions
| Function | Parameters | Description |
|----------|------------|-------------|
| `SEED` | seed (INTEGR) | Set specific random seed |
| `SEED_TIME` | none | Seed with current time |

### Basic Random Functions
| Function | Parameters | Returns | Description |
|----------|------------|---------|-------------|
| `RANDOM_FLOAT` | none | DUBBLE | Random float [0,1) |
| `RANDOM_INT` | min, max (INTEGR) | INTEGR | Random integer [min,max) |
| `RANDOM_BOOL` | none | BOOL | Random boolean |

### String Generation
| Function | Parameters | Returns | Description |
|----------|------------|---------|-------------|
| `RANDOM_STRING` | length (INTEGR), charset (STRIN) | STRIN | Random string from charset |
| `UUID` | none | STRIN | Generate UUID |

### Collection Functions
| Function | Parameters | Returns | Description |
|----------|------------|---------|-------------|
| `RANDOM_CHOICE` | array (BUKKIT) | Value | Random element from array |
| `SHUFFLE` | array (BUKKIT) | none | Shuffle array in-place |

## Related

- [MATH Module](math.md) - For mathematical functions used with random numbers
- [STRING Module](string.md) - For string manipulation with generated strings
- [STDIO Module](stdio.md) - For displaying random results