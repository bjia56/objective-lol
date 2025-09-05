# MATH Module - Mathematical Functions

The MATH module provides essential mathematical functions for numerical computations.

## Import

```lol
BTW Full import
I CAN HAS MATH?

BTW Selective import examples
I CAN HAS ABS FROM MATH?
I CAN HAS SQRT AN POW FROM MATH?
I CAN HAS PI AN E FROM MATH?
```

## Mathematical Constants

### PI - Pi Constant

The mathematical constant π (pi) ≈ 3.14159.

**Type:** DUBBLE
**Value:** 3.14159265359

```lol
I CAN HAS PI FROM MATH?

I HAS A VARIABLE CIRCLE_AREA TEH DUBBLE ITZ PI TIEMZ RADIUS TIEMZ RADIUS
I HAS A VARIABLE RADIANS TEH DUBBLE ITZ DEGREES TIEMZ PI DIVIDEZ 180.0
```

### E - Euler's Number

The mathematical constant e (Euler's number) ≈ 2.71828.

**Type:** DUBBLE
**Value:** 2.71828182846

```lol
I CAN HAS E FROM MATH?

I HAS A VARIABLE NATURAL_LOG TEH DUBBLE ITZ LOG WIT E    BTW Result: 1.0
I HAS A VARIABLE EXPONENTIAL TEH DUBBLE ITZ EXP WIT 1.0  BTW Result: E
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
I CAN HAS SIN AN PI FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ SIN WIT 0.0           BTW 0.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ SIN WIT PI DIVIDEZ 2  BTW 1.0 (sin(π/2))
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ SIN WIT PI            BTW ≈0.0 (sin(π))
```

### COS - Cosine Function

Returns the cosine of an angle (in radians).

**Syntax:** `COS WIT <angle_radians>`
**Returns:** DUBBLE

```lol
I CAN HAS COS AN PI FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ COS WIT 0.0           BTW 1.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ COS WIT PI DIVIDEZ 2  BTW ≈0.0 (cos(π/2))
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ COS WIT PI            BTW -1.0 (cos(π))
```

### TAN - Tangent Function

Returns the tangent of an angle (in radians).

**Syntax:** `TAN WIT <angle_radians>`
**Returns:** DUBBLE

```lol
I CAN HAS TAN AN PI FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ TAN WIT 0.0           BTW 0.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ TAN WIT PI DIVIDEZ 4  BTW 1.0 (tan(π/4))
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ TAN WIT PI            BTW ≈0.0 (tan(π))
```

### ASIN - Arcsine Function

Returns the arcsine (inverse sine) of a value. Input must be in range [-1, 1].

**Syntax:** `ASIN WIT <value>`
**Returns:** DUBBLE (angle in radians)

```lol
I CAN HAS ASIN FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ ASIN WIT 0.0    BTW 0.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ ASIN WIT 1.0    BTW π/2
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ ASIN WIT -1.0   BTW -π/2
```

### ACOS - Arccosine Function

Returns the arccosine (inverse cosine) of a value. Input must be in range [-1, 1].

**Syntax:** `ACOS WIT <value>`
**Returns:** DUBBLE (angle in radians)

```lol
I CAN HAS ACOS FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ ACOS WIT 1.0    BTW 0.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ ACOS WIT 0.0    BTW π/2
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ ACOS WIT -1.0   BTW π
```

### ATAN - Arctangent Function

Returns the arctangent (inverse tangent) of a value.

**Syntax:** `ATAN WIT <value>`
**Returns:** DUBBLE (angle in radians)

```lol
I CAN HAS ATAN FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ ATAN WIT 0.0    BTW 0.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ ATAN WIT 1.0    BTW π/4
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ ATAN WIT -1.0   BTW -π/4
```

### ATAN2 - Two-Argument Arctangent

Returns the angle (in radians) from the x-axis to the point (x,y).

**Syntax:** `ATAN2 WIT <y> AN WIT <x>`
**Returns:** DUBBLE (angle in radians)

```lol
I CAN HAS ATAN2 FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ ATAN2 WIT 1.0 AN WIT 1.0    BTW π/4
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ ATAN2 WIT 1.0 AN WIT 0.0    BTW π/2
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ ATAN2 WIT 0.0 AN WIT 1.0    BTW 0.0
```

## Logarithmic and Exponential Functions

### LOG - Natural Logarithm

Returns the natural logarithm (base e) of a number.

**Syntax:** `LOG WIT <value>`
**Returns:** DUBBLE

```lol
I CAN HAS LOG AN E FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ LOG WIT E      BTW 1.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ LOG WIT 1.0    BTW 0.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ LOG WIT 10.0   BTW ≈2.3026
```

### LOG10 - Base-10 Logarithm

Returns the base-10 logarithm of a number.

**Syntax:** `LOG10 WIT <value>`
**Returns:** DUBBLE

```lol
I CAN HAS LOG10 FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ LOG10 WIT 10.0    BTW 1.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ LOG10 WIT 100.0   BTW 2.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ LOG10 WIT 1.0     BTW 0.0
```

### LOG2 - Base-2 Logarithm

Returns the base-2 logarithm of a number.

**Syntax:** `LOG2 WIT <value>`
**Returns:** DUBBLE

```lol
I CAN HAS LOG2 FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ LOG2 WIT 2.0     BTW 1.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ LOG2 WIT 8.0     BTW 3.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ LOG2 WIT 1.0     BTW 0.0
```

### EXP - Exponential Function

Returns e raised to the power of a number.

**Syntax:** `EXP WIT <value>`
**Returns:** DUBBLE

```lol
I CAN HAS EXP AN E FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ EXP WIT 1.0     BTW E (≈2.718)
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ EXP WIT 0.0     BTW 1.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ EXP WIT 2.0     BTW E² (≈7.389)
```

## Rounding Functions

### CEIL - Ceiling Function

Returns the smallest integer greater than or equal to the given value.

**Syntax:** `CEIL WIT <value>`
**Returns:** DUBBLE

```lol
I CAN HAS CEIL FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ CEIL WIT 3.2     BTW 4.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ CEIL WIT 3.0     BTW 3.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ CEIL WIT -3.7    BTW -3.0
```

### FLOOR - Floor Function

Returns the largest integer less than or equal to the given value.

**Syntax:** `FLOOR WIT <value>`
**Returns:** DUBBLE

```lol
I CAN HAS FLOOR FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ FLOOR WIT 3.7    BTW 3.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ FLOOR WIT 3.0    BTW 3.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ FLOOR WIT -3.2   BTW -4.0
```

### ROUND - Round Function

Returns the value rounded to the nearest integer.

**Syntax:** `ROUND WIT <value>`
**Returns:** DUBBLE

```lol
I CAN HAS ROUND FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ ROUND WIT 3.4    BTW 3.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ ROUND WIT 3.5    BTW 4.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ ROUND WIT 3.6    BTW 4.0
```

### TRUNC - Truncate Function

Returns the integer part of a number, removing any fractional digits.

**Syntax:** `TRUNC WIT <value>`
**Returns:** DUBBLE

```lol
I CAN HAS TRUNC FROM MATH?

I HAS A VARIABLE RESULT1 TEH DUBBLE ITZ TRUNC WIT 3.7    BTW 3.0
I HAS A VARIABLE RESULT2 TEH DUBBLE ITZ TRUNC WIT -3.7   BTW -3.0
I HAS A VARIABLE RESULT3 TEH DUBBLE ITZ TRUNC WIT 5.0    BTW 5.0
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

## Trigonometry Example

```lol
I CAN HAS MATH?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEGREES_TO_RADIANS WIT DEGREES TEH DUBBLE
    GIVEZ DEGREES TIEMZ PI DIVIDEZ 180.0
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    BTW Calculate sine and cosine of 45 degrees
    I HAS A VARIABLE ANGLE_DEG TEH DUBBLE ITZ 45.0
    I HAS A VARIABLE ANGLE_RAD TEH DUBBLE ITZ DEGREES_TO_RADIANS WIT ANGLE_DEG

    I HAS A VARIABLE SINE_VAL TEH DUBBLE ITZ SIN WIT ANGLE_RAD
    I HAS A VARIABLE COSINE_VAL TEH DUBBLE ITZ COS WIT ANGLE_RAD
    I HAS A VARIABLE TANGENT_VAL TEH DUBBLE ITZ TAN WIT ANGLE_RAD

    SAY WIT "sin(45°) = "
    SAYZ WIT SINE_VAL      BTW ≈ 0.7071
    SAY WIT "cos(45°) = "
    SAYZ WIT COSINE_VAL    BTW ≈ 0.7071
    SAY WIT "tan(45°) = "
    SAYZ WIT TANGENT_VAL   BTW ≈ 1.0
KTHXBAI
```

## Function Summary

### Constants

| Constant | Type | Value | Description |
|----------|------|-------|-------------|
| `PI` | DUBBLE | 3.14159... | Mathematical constant π |
| `E` | DUBBLE | 2.71828... | Euler's number |

### Functions

| Function | Parameters | Returns | Description |
|----------|------------|---------|-------------|
| `ABS` | number | DUBBLE | Absolute value |
| `MAX` | value1, value2 | DUBBLE | Maximum of two values |
| `MIN` | value1, value2 | DUBBLE | Minimum of two values |
| `SQRT` | number | DUBBLE | Square root |
| `POW` | base, exponent | DUBBLE | Power function |
| `SIN` | angle_radians | DUBBLE | Sine function |
| `COS` | angle_radians | DUBBLE | Cosine function |
| `TAN` | angle_radians | DUBBLE | Tangent function |
| `ASIN` | value | DUBBLE | Arcsine function |
| `ACOS` | value | DUBBLE | Arccosine function |
| `ATAN` | value | DUBBLE | Arctangent function |
| `ATAN2` | y, x | DUBBLE | Two-argument arctangent |
| `LOG` | value | DUBBLE | Natural logarithm |
| `LOG10` | value | DUBBLE | Base-10 logarithm |
| `LOG2` | value | DUBBLE | Base-2 logarithm |
| `EXP` | value | DUBBLE | Exponential function |
| `CEIL` | value | DUBBLE | Ceiling function |
| `FLOOR` | value | DUBBLE | Floor function |
| `ROUND` | value | DUBBLE | Rounding function |
| `TRUNC` | value | DUBBLE | Truncation function |

## Related

- [STDIO Module](stdio.md) - For displaying calculation results
- [RANDOM Module](random.md) - For random number generation
- [Examples](../examples/calculator.md) - Mathematical calculation examples