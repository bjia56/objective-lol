# MATH Module - Mathematical Functions

The MATH module provides essential mathematical functions for numerical computations.

## Import

```lol
BTW Full import
I CAN HAS MATH?

BTW Selective import examples
I CAN HAS ABS FROM MATH?
I CAN HAS ABS AN MAX AN MIN FROM MATH?
I CAN HAS SQRT AN POW FROM MATH?
I CAN HAS RANDOM AN RANDINT FROM MATH?
I CAN HAS SIN AN COS FROM MATH?
```

## Basic Mathematical Functions

### ABS - Absolute Value

Returns the absolute value of a number.

**Syntax:** `ABS WIT <number>`
**Returns:** DUBBLE

```lol
I CAN HAS ABS FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ ABS WIT -5.5      BTW 5.5
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ ABS WIT 42        BTW 42.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ ABS WIT 0         BTW 0.0
```

### MAX - Maximum of Two Values

Returns the larger of two values.

**Syntax:** `MAX WIT <value1> AN WIT <value2>`
**Returns:** DUBBLE

```lol
I CAN HAS MAX FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ MAX WIT 10.5 AN WIT 7.2    BTW 10.5
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ MAX WIT -3 AN WIT -8       BTW -3.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ MAX WIT 0 AN WIT 0         BTW 0.0
```

### MIN - Minimum of Two Values

Returns the smaller of two values.

**Syntax:** `MIN WIT <value1> AN WIT <value2>`
**Returns:** DUBBLE

```lol
I CAN HAS MIN FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ MIN WIT 10.5 AN WIT 7.2    BTW 7.2
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ MIN WIT -3 AN WIT -8       BTW -8.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ MIN WIT 5 AN WIT 5         BTW 5.0
```

## Advanced Mathematical Functions

### SQRT - Square Root

Returns the square root of a number.

**Syntax:** `SQRT WIT <number>`
**Returns:** DUBBLE

```lol
I CAN HAS SQRT FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ SQRT WIT 16.0     BTW 4.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ SQRT WIT 2.0      BTW 1.414...
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ SQRT WIT 0.0      BTW 0.0
```

### POW - Power Function

Returns base raised to the power of exponent.

**Syntax:** `POW WIT <base> AN WIT <exponent>`
**Returns:** DUBBLE

```lol
I CAN HAS POW FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ POW WIT 2.0 AN WIT 3.0     BTW 8.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ POW WIT 10.0 AN WIT 2.0    BTW 100.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ POW WIT 4.0 AN WIT 0.5     BTW 2.0 (square root)
```

## Trigonometric Functions

### SIN - Sine Function

Returns the sine of an angle (in radians).

**Syntax:** `SIN WIT <angle_radians>`
**Returns:** DUBBLE

```lol
I CAN HAS SIN FROM MATH?

I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159
I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ SIN WIT 0.0           BTW 0.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ SIN WIT PI DIVIDEZ 2  BTW 1.0 (sin(π/2))
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ SIN WIT PI            BTW ≈0.0 (sin(π))
```

### COS - Cosine Function

Returns the cosine of an angle (in radians).

**Syntax:** `COS WIT <angle_radians>`
**Returns:** DUBBLE

```lol
I CAN HAS COS FROM MATH?

I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159
I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ COS WIT 0.0           BTW 1.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ COS WIT PI DIVIDEZ 2  BTW ≈0.0 (cos(π/2))
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ COS WIT PI            BTW -1.0 (cos(π))
```

## Random Number Functions

### RANDOM - Random Float

Returns a random floating-point number between 0 (inclusive) and 1 (exclusive).

**Syntax:** `RANDOM`
**Returns:** DUBBLE

```lol
I CAN HAS RANDOM FROM MATH?

I HAS A VARIABLE RAND1 TEH DUBBLE ITZ RANDOM    BTW e.g., 0.7234
I HAS A VARIABLE RAND2 TEH DUBBLE ITZ RANDOM    BTW e.g., 0.1849

BTW Generate random number in range [0, 100)
I HAS A VARIABLE PERCENT TEH DUBBLE ITZ RANDOM TIEMZ 100
```

### RANDINT - Random Integer

Returns a random integer in the range [min, max).

**Syntax:** `RANDINT WIT <min> AN WIT <max>`
**Returns:** INTEGR

```lol
I CAN HAS RANDINT FROM MATH?

BTW Dice roll (1-6)
I HAS A VARIABLE DICE TEH INTEGR ITZ RANDINT WIT 1 AN WIT 7

BTW Random index for array of size 10
I HAS A VARIABLE INDEX TEH INTEGR ITZ RANDINT WIT 0 AN WIT 10

BTW Random number between -10 and 10
I HAS A VARIABLE SIGNED TEH INTEGR ITZ RANDINT WIT -10 AN WIT 11
```

## Usage Examples

### Calculator Functions

```lol
I CAN HAS MATH?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN CALCULATE_HYPOTENUSE WIT A TEH DUBBLE AN WIT B TEH DUBBLE
    I HAS A VARIABLE A_SQUARED TEH DUBBLE ITZ POW WIT A AN WIT 2
    I HAS A VARIABLE B_SQUARED TEH DUBBLE ITZ POW WIT B AN WIT 2
    I HAS A VARIABLE C_SQUARED TEH DUBBLE ITZ A_SQUARED MOAR B_SQUARED
    GIVEZ SQRT WIT C_SQUARED
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE HYPOTENUSE TEH DUBBLE ITZ CALCULATE_HYPOTENUSE WIT 3.0 AN WIT 4.0
    SAYZ WIT HYPOTENUSE    BTW 5.0
KTHXBAI
```

### Statistics Functions

```lol
I CAN HAS MATH?

HAI ME TEH FUNCSHUN FIND_RANGE WIT NUMBERS TEH BUKKIT
    I HAS A VARIABLE MIN_VAL TEH DUBBLE ITZ NUMBERS DO AT WIT 0
    I HAS A VARIABLE MAX_VAL TEH DUBBLE ITZ NUMBERS DO AT WIT 0
    I HAS A VARIABLE INDEX TEH INTEGR ITZ 1

    WHILE INDEX SMALLR THAN NUMBERS SIZ
        I HAS A VARIABLE CURRENT TEH DUBBLE ITZ NUMBERS DO AT WIT INDEX
        MIN_VAL ITZ MIN WIT MIN_VAL AN WIT CURRENT
        MAX_VAL ITZ MAX WIT MAX_VAL AN WIT CURRENT
        INDEX ITZ INDEX MOAR 1
    KTHX

    GIVEZ MAX_VAL LES MIN_VAL
KTHXBAI
```

### Random Number Generation

```lol
I CAN HAS MATH?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN SIMULATE_DICE_ROLLS WIT NUM_ROLLS TEH INTEGR
    I HAS A VARIABLE TOTAL TEH INTEGR ITZ 0
    I HAS A VARIABLE ROLL_COUNT TEH INTEGR ITZ 0

    WHILE ROLL_COUNT SMALLR THAN NUM_ROLLS
        I HAS A VARIABLE ROLL TEH INTEGR ITZ RANDINT WIT 1 AN WIT 7
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

## Trigonometry Example

```lol
I CAN HAS MATH?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEGREES_TO_RADIANS WIT DEGREES TEH DUBBLE
    I HAS A VARIABLE PI TEH DUBBLE ITZ 3.14159265359
    GIVEZ DEGREES TIEMZ PI DIVIDEZ 180.0
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    BTW Calculate sine and cosine of 45 degrees
    I HAS A VARIABLE ANGLE_DEG TEH DUBBLE ITZ 45.0
    I HAS A VARIABLE ANGLE_RAD TEH DUBBLE ITZ DEGREES_TO_RADIANS WIT ANGLE_DEG

    I HAS A VARIABLE SINE_VAL TEH DUBBLE ITZ SIN WIT ANGLE_RAD
    I HAS A VARIABLE COSINE_VAL TEH DUBBLE ITZ COS WIT ANGLE_RAD

    SAY WIT "sin(45°) = "
    SAYZ WIT SINE_VAL      BTW ≈ 0.7071
    SAY WIT "cos(45°) = "
    SAYZ WIT COSINE_VAL    BTW ≈ 0.7071
KTHXBAI
```

## Function Summary

| Function | Parameters | Returns | Description |
|----------|------------|---------|-------------|
| `ABS` | number | DUBBLE | Absolute value |
| `MAX` | value1, value2 | DUBBLE | Maximum of two values |
| `MIN` | value1, value2 | DUBBLE | Minimum of two values |
| `SQRT` | number | DUBBLE | Square root |
| `POW` | base, exponent | DUBBLE | Power function |
| `SIN` | angle_radians | DUBBLE | Sine function |
| `COS` | angle_radians | DUBBLE | Cosine function |
| `RANDOM` | none | DUBBLE | Random float [0,1) |
| `RANDINT` | min, max | INTEGR | Random integer [min,max) |

## Related

- [STDIO Module](stdio.md) - For displaying calculation results
- [Examples](../examples/calculator.md) - Mathematical calculation examples