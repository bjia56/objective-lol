# MATH Module

## Import

```lol
BTW Full import
I CAN HAS MATH?

BTW Selective import examples
I CAN HAS MAX FROM MATH?
I CAN HAS LOG2 FROM MATH?
I CAN HAS PI FROM MATH?
```

## Mathematical Constants

### E

Euler's number e ≈ 2.71828.
The base of natural logarithms, fundamental mathematical constant.

**Type:** DUBBLE
**Value:** 2.71828182846

```lol
I CAN HAS E FROM MODULE?

I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG WIT E
BTW Result: 1.0 (ln(e) = 1)
```

```lol
I CAN HAS E FROM MODULE?

I HAS A VARIABLE RESULT TEH DUBBLE ITZ EXP WIT 2.0
BTW Result: 7.389056099 (e^2)
```

### PI

The mathematical constant π (pi) ≈ 3.14159.
Represents the ratio of a circle's circumference to its diameter.

**Type:** DUBBLE
**Value:** 3.14159265359

```lol
I CAN HAS PI FROM MODULE?

I HAS A VARIABLE RADIUS TEH DUBBLE ITZ 2.0
I HAS A VARIABLE AREA TEH DUBBLE ITZ PI TIEMZ RADIUS TIEMZ RADIUS
BTW Result: 12.566370614
```

```lol
I CAN HAS PI FROM MODULE?

I HAS A VARIABLE DEGREES TEH DUBBLE ITZ 180.0
I HAS A VARIABLE RADIANS TEH DUBBLE ITZ DEGREES TIEMZ PI DIVIDEZ 180.0
BTW Result: 3.14159265359 (π radians)
```

## Basic Math

### ABS

Returns the absolute value of a number.
Removes the sign and returns the positive magnitude.

**Syntax:** `ABS WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The number to get absolute value of

**Example: Basic absolute value**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT -5.5
BTW Result: 5.5
```

**Example: Positive number unchanged**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT 42.0
BTW Result: 42.0
```

**Example: Zero unchanged**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ABS WIT 0.0
BTW Result: 0.0
```

**Note:** Works with both positive and negative numbers

**See also:** MAX, MIN

### MAX

Returns the larger of two numbers.
Compares two values and returns the maximum.

**Syntax:** `MAX WIT <value1> AN WIT <value2>`
**Returns:** DUBBLE

**Parameters:**
- `a` (DUBBLE): First number to compare
- `b` (DUBBLE): Second number to compare

**Example: Compare positive numbers**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ MAX WIT 10.5 AN WIT 7.2
BTW Result: 10.5
```

**Example: Compare negative numbers**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ MAX WIT -3.0 AN WIT -8.0
BTW Result: -3.0
```

**Example: Equal values**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ MAX WIT 5.0 AN WIT 5.0
BTW Result: 5.0
```

**Note:** Returns the first value if both are equal

**See also:** MIN, ABS

### MIN

Returns the smaller of two numbers.
Compares two values and returns the minimum.

**Syntax:** `MIN WIT <value1> AN WIT <value2>`
**Returns:** DUBBLE

**Parameters:**
- `a` (DUBBLE): First number to compare
- `b` (DUBBLE): Second number to compare

**Example: Compare positive numbers**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ MIN WIT 10.5 AN WIT 7.2
BTW Result: 7.2
```

**Example: Compare negative numbers**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ MIN WIT -3.0 AN WIT -8.0
BTW Result: -8.0
```

**Example: Equal values**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ MIN WIT 5.0 AN WIT 5.0
BTW Result: 5.0
```

**Note:** Returns the first value if both are equal

**See also:** MAX, ABS

## Advanced Math

### POW

Returns base raised to the power of exponent (base^exponent).
Performs exponentiation using floating-point arithmetic.

**Syntax:** `POW WIT <base> AN WIT <exponent>`
**Returns:** DUBBLE

**Parameters:**
- `base` (DUBBLE): The base number
- `exponent` (DUBBLE): The power to raise the base to

**Example: Integer exponent**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ POW WIT 2.0 AN WIT 3.0
BTW Result: 8.0
```

**Example: Fractional exponent (square root)**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ POW WIT 4.0 AN WIT 0.5
BTW Result: 2.0
```

**Example: Power of ten**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ POW WIT 10.0 AN WIT 2.0
BTW Result: 100.0
```

**Example: Zero exponent**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ POW WIT 5.0 AN WIT 0.0
BTW Result: 1.0
```

**Note:** Any number to the power of 0 equals 1

**See also:** SQRT, EXP, LOG

### SQRT

Returns the square root of a number.
Input must be non-negative. Throws error for negative values.

**Syntax:** `SQRT WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The number to get square root of (must be ≥ 0)

**Example: Perfect square**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ SQRT WIT 16.0
BTW Result: 4.0
```

**Example: Decimal result**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ SQRT WIT 2.0
BTW Result: 1.4142135623
```

**Example: Zero input**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ SQRT WIT 0.0
BTW Result: 0.0
```

**Note:** Input must be non-negative

**See also:** POW, ABS

## Trigonometry

### ACOS

Returns the arc cosine (inverse cosine) of a value in radians.
Input must be in range [-1, 1]. Result is in range [0, π].

**Syntax:** `ACOS WIT <value>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): Input value (must be between -1 and 1)

**Example: Arc cosine of 1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ACOS WIT 1.0
BTW Result: 0.0
```

**Example: Arc cosine of 0**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ACOS WIT 0.0
BTW Result: π/2 (≈1.5708)
```

**Example: Arc cosine of -1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ACOS WIT -1.0
BTW Result: π (≈3.1416)
```

**Note:** Input must be in range [-1, 1]

**See also:** COS, ASIN, ATAN, PI

### ASIN

Returns the arc sine (inverse sine) of a value in radians.
Input must be in range [-1, 1]. Result is in range [-π/2, π/2].

**Syntax:** `ASIN WIT <value>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): Input value (must be between -1 and 1)

**Example: Arc sine of 0**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ASIN WIT 0.0
BTW Result: 0.0
```

**Example: Arc sine of 1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ASIN WIT 1.0
BTW Result: π/2 (≈1.5708)
```

**Example: Arc sine of -1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ASIN WIT -1.0
BTW Result: -π/2 (≈-1.5708)
```

**Note:** Input must be in range [-1, 1]

**See also:** SIN, ACOS, ATAN, PI

### ATAN

Returns the arc tangent (inverse tangent) of a value in radians.
Result is in range [-π/2, π/2].

**Syntax:** `ATAN WIT <value>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): Input value (any real number)

**Example: Arc tangent of 0**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN WIT 0.0
BTW Result: 0.0
```

**Example: Arc tangent of 1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN WIT 1.0
BTW Result: π/4 (≈0.7854)
```

**Example: Arc tangent of -1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN WIT -1.0
BTW Result: -π/4 (≈-0.7854)
```

**Note:** Input can be any real number

**See also:** TAN, ATAN2, ASIN, ACOS, PI

### ATAN2

Returns the arc tangent of y/x in radians, considering quadrant.
Result is in range [-π, π]. More robust than ATAN for coordinate conversion.

**Syntax:** `ATAN2 WIT <y> AN WIT <x>`
**Returns:** DUBBLE

**Parameters:**
- `y` (DUBBLE): Y coordinate
- `x` (DUBBLE): X coordinate

**Example: Point in first quadrant**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN2 WIT 1.0 AN WIT 1.0
BTW Result: π/4 (≈0.7854)
```

**Example: Point on positive y-axis**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN2 WIT 1.0 AN WIT 0.0
BTW Result: π/2 (≈1.5708)
```

**Example: Point on positive x-axis**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ATAN2 WIT 0.0 AN WIT 1.0
BTW Result: 0.0
```

**Note:** More robust than ATAN for coordinate conversion

**Note:** Handles all quadrants correctly

**See also:** ATAN, TAN, PI

### COS

Returns the cosine of an angle in radians.
Input angle should be in radians, not degrees.

**Syntax:** `COS WIT <angle_radians>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The angle in radians

**Example: Cosine of 0**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ COS WIT 0.0
BTW Result: 1.0
```

**Example: Cosine of π/2 (90 degrees)**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ COS WIT PI DIVIDEZ 2
BTW Result: ≈0.0
```

**Example: Cosine of π (180 degrees)**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ COS WIT PI
BTW Result: -1.0
```

**Note:** Input must be in radians, not degrees

**Note:** Result is always between -1 and 1

**See also:** SIN, TAN, ACOS, PI

### SIN

Returns the sine of an angle in radians.
Input angle should be in radians, not degrees.

**Syntax:** `SIN WIT <angle_radians>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The angle in radians

**Example: Sine of 0**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ SIN WIT 0.0
BTW Result: 0.0
```

**Example: Sine of π/2 (90 degrees)**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ SIN WIT PI DIVIDEZ 2
BTW Result: 1.0
```

**Example: Sine of π (180 degrees)**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ SIN WIT PI
BTW Result: ≈0.0
```

**Note:** Input must be in radians, not degrees

**Note:** Result is always between -1 and 1

**See also:** COS, TAN, ASIN, PI

### TAN

Returns the tangent of an angle in radians.
Input angle should be in radians. Undefined at π/2 + nπ.

**Syntax:** `TAN WIT <angle_radians>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The angle in radians

**Example: Tangent of 0**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ TAN WIT 0.0
BTW Result: 0.0
```

**Example: Tangent of π/4 (45 degrees)**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ TAN WIT PI DIVIDEZ 4
BTW Result: 1.0
```

**Example: Tangent of π (180 degrees)**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ TAN WIT PI
BTW Result: ≈0.0
```

**Note:** Input must be in radians, not degrees

**Note:** Undefined at π/2 + nπ (90°, 270°, etc.)

**See also:** SIN, COS, ATAN, PI

## Logarithmic

### EXP

Returns e raised to the power of the given value (e^value).
The exponential function, inverse of natural logarithm.

**Syntax:** `EXP WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The exponent to raise e to

**Example: e to the power of 1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ EXP WIT 1.0
BTW Result: 2.718281828 (e)
```

**Example: e to the power of 0**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ EXP WIT 0.0
BTW Result: 1.0
```

**Example: e to the power of 2**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ EXP WIT 2.0
BTW Result: 7.389056099 (e²)
```

**Note:** The exponential function, inverse of natural logarithm

**Note:** EXP(LOG(x)) = x for positive x

**See also:** LOG, LOG10, LOG2, E, POW

### LOG

Returns the natural logarithm (base e) of a number.
Input must be positive. Throws error for zero or negative values.

**Syntax:** `LOG WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The number to get natural logarithm of (must be > 0)

**Example: Natural log of e**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG WIT E
BTW Result: 1.0
```

**Example: Natural log of 1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG WIT 1.0
BTW Result: 0.0
```

**Example: Natural log of 10**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG WIT 10.0
BTW Result: 2.302585093
```

**Note:** Input must be positive (greater than zero)

**See also:** LOG10, LOG2, EXP, E

### LOG10

Returns the base-10 logarithm of a number.
Input must be positive. Common logarithm for scientific calculations.

**Syntax:** `LOG10 WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The number to get base-10 logarithm of (must be > 0)

**Example: Log base 10 of 10**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG10 WIT 10.0
BTW Result: 1.0
```

**Example: Log base 10 of 100**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG10 WIT 100.0
BTW Result: 2.0
```

**Example: Log base 10 of 1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG10 WIT 1.0
BTW Result: 0.0
```

**Note:** Input must be positive (greater than zero)

**Note:** Common logarithm for scientific calculations

**See also:** LOG, LOG2, EXP

### LOG2

Returns the base-2 logarithm of a number.
Input must be positive. Useful for binary and computer science calculations.

**Syntax:** `LOG2 WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The number to get base-2 logarithm of (must be > 0)

**Example: Log base 2 of 2**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG2 WIT 2.0
BTW Result: 1.0
```

**Example: Log base 2 of 8**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG2 WIT 8.0
BTW Result: 3.0
```

**Example: Log base 2 of 1**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ LOG2 WIT 1.0
BTW Result: 0.0
```

**Note:** Input must be positive (greater than zero)

**Note:** Useful for binary and computer science calculations

**See also:** LOG, LOG10, EXP

## Rounding

### CEIL

Returns the smallest integer greater than or equal to the value (ceiling).
Rounds up to the next whole number.

**Syntax:** `CEIL WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The number to round up

**Example: Ceiling of positive decimal**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ CEIL WIT 3.2
BTW Result: 4.0
```

**Example: Ceiling of whole number**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ CEIL WIT 3.0
BTW Result: 3.0
```

**Example: Ceiling of negative number**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ CEIL WIT -3.7
BTW Result: -3.0
```

**Note:** Always rounds up to the next whole number

**Note:** For negative numbers, rounds towards zero

**See also:** FLOOR, ROUND, TRUNC

### FLOOR

Returns the largest integer less than or equal to the value (floor).
Rounds down to the previous whole number.

**Syntax:** `FLOOR WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The number to round down

**Example: Floor of positive decimal**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ FLOOR WIT 3.7
BTW Result: 3.0
```

**Example: Floor of whole number**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ FLOOR WIT 3.0
BTW Result: 3.0
```

**Example: Floor of negative number**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ FLOOR WIT -3.2
BTW Result: -4.0
```

**Note:** Always rounds down to the previous whole number

**Note:** For negative numbers, rounds away from zero

**See also:** CEIL, ROUND, TRUNC

### ROUND

Returns the value rounded to the nearest integer.
Rounds 0.5 up to the next integer (round half up).

**Syntax:** `ROUND WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The number to round

**Example: Round down**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ROUND WIT 3.4
BTW Result: 3.0
```

**Example: Round up**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ROUND WIT 3.6
BTW Result: 4.0
```

**Example: Round half up**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ROUND WIT 3.5
BTW Result: 4.0
```

**Example: Negative numbers**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ ROUND WIT -3.5
BTW Result: -3.0
```

**Note:** Uses "round half up" strategy for 0.5 values

**See also:** CEIL, FLOOR, TRUNC

### TRUNC

Returns the integer part of a number by removing the fractional part.
Truncates towards zero, different from floor for negative numbers.

**Syntax:** `TRUNC WIT <number>`
**Returns:** DUBBLE

**Parameters:**
- `value` (DUBBLE): The number to truncate

**Example: Truncate positive decimal**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ TRUNC WIT 3.7
BTW Result: 3.0
```

**Example: Truncate negative decimal**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ TRUNC WIT -3.7
BTW Result: -3.0
```

**Example: Truncate whole number**

```lol
I HAS A VARIABLE RESULT TEH DUBBLE ITZ TRUNC WIT 5.0
BTW Result: 5.0
```

**Note:** Always truncates towards zero

**Note:** Different from FLOOR for negative numbers

**See also:** CEIL, FLOOR, ROUND

