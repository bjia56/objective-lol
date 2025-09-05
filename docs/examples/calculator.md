# Calculator Example

A simple calculator program demonstrating basic arithmetic operations, functions, and error handling.

## Basic Calculator

```lol
BTW Simple calculator program
I CAN HAS STDIO?
I CAN HAS MATH?

HAI ME TEH FUNCSHUN CALCULATE TEH DUBBLE WIT X TEH DUBBLE AN WIT OP TEH STRIN AN WIT Y TEH DUBBLE
    IZ OP SAEM AS "+"?
        GIVEZ X MOAR Y
    KTHX

    IZ OP SAEM AS "-"?
        GIVEZ X LES Y
    KTHX

    IZ OP SAEM AS "*"?
        GIVEZ X TIEMZ Y
    KTHX

    IZ OP SAEM AS "/"?
        IZ Y SAEM AS 0.0?
            OOPS "Division by zero!"
        KTHX
        GIVEZ X DIVIDEZ Y
    KTHX

    OOPS "Unknown operator"
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    SAYZ WIT "=== Simple Calculator ==="
    SAYZ WIT ""

    BTW Get first number
    SAYZ WIT "Enter first number: "
    I HAS A VARIABLE NUM1_STR TEH STRIN ITZ GIMME

    BTW Get operator
    SAYZ WIT "Enter operator (+, -, *, /): "
    I HAS A VARIABLE OPERATOR TEH STRIN ITZ GIMME

    BTW Get second number
    SAYZ WIT "Enter second number: "
    I HAS A VARIABLE NUM2_STR TEH STRIN ITZ GIMME

    BTW Convert strings to numbers and calculate
    MAYB
        I HAS A VARIABLE NUM1 TEH DUBBLE ITZ NUM1_STR AS DUBBLE
        I HAS A VARIABLE NUM2 TEH DUBBLE ITZ NUM2_STR AS DUBBLE
        I HAS A VARIABLE RESULT TEH DUBBLE ITZ CALCULATE WIT NUM1 AN WIT OPERATOR AN WIT NUM2

        SAY WIT NUM1
        SAY WIT " "
        SAY WIT OPERATOR
        SAY WIT " "
        SAY WIT NUM2
        SAY WIT " = "
        SAYZ WIT RESULT

    OOPSIE ERROR_MSG
        SAYZ WIT "Error: "
        SAYZ WIT ERROR_MSG
    KTHX
KTHXBAI
```

## Advanced Calculator with Menu

```lol
BTW Advanced calculator with menu system
I CAN HAS STDIO?
I CAN HAS MATH?

HAI ME TEH FUNCSHUN SHOW_MENU
    SAYZ WIT ""
    SAYZ WIT "=== Advanced Calculator ==="
    SAYZ WIT "1. Basic arithmetic"
    SAYZ WIT "2. Power function"
    SAYZ WIT "3. Square root"
    SAYZ WIT "4. Trigonometry"
    SAYZ WIT "5. Random number"
    SAYZ WIT "6. Exit"
    SAY WIT "Choose option (1-6): "
KTHXBAI

HAI ME TEH FUNCSHUN BASIC_ARITHMETIC
    SAYZ WIT "Enter first number: "
    I HAS A VARIABLE X TEH DUBBLE ITZ GIMME AS DUBBLE

    SAYZ WIT "Enter operator (+, -, *, /): "
    I HAS A VARIABLE OP TEH STRIN ITZ GIMME

    SAYZ WIT "Enter second number: "
    I HAS A VARIABLE Y TEH DUBBLE ITZ GIMME AS DUBBLE

    MAYB
        IZ OP SAEM AS "+"?
            SAY WIT X
            SAY WIT " + "
            SAY WIT Y
            SAY WIT " = "
            SAYZ WIT X MOAR Y
        KTHX
        IZ OP SAEM AS "-"?
            SAY WIT X
            SAY WIT " - "
            SAY WIT Y
            SAY WIT " = "
            SAYZ WIT X LES Y
        KTHX
        IZ OP SAEM AS "*"?
            SAY WIT X
            SAY WIT " * "
            SAY WIT Y
            SAY WIT " = "
            SAYZ WIT X TIEMZ Y
        KTHX
        IZ OP SAEM AS "/"?
            IZ Y SAEM AS 0.0?
                OOPS "Division by zero!"
            KTHX
            SAY WIT X
            SAY WIT " / "
            SAY WIT Y
            SAY WIT " = "
            SAYZ WIT X DIVIDEZ Y
        KTHX
        SAYZ WIT "Unknown operator!"
    OOPSIE ERR
        SAYZ WIT "Error: "
        SAYZ WIT ERR
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN POWER_FUNCTION
    SAYZ WIT "Enter base: "
    I HAS A VARIABLE BASE TEH DUBBLE ITZ GIMME AS DUBBLE

    SAYZ WIT "Enter exponent: "
    I HAS A VARIABLE EXP TEH DUBBLE ITZ GIMME AS DUBBLE

    I HAS A VARIABLE RESULT TEH DUBBLE ITZ POW WIT BASE AN WIT EXP
    SAY WIT BASE
    SAY WIT " ^ "
    SAY WIT EXP
    SAY WIT " = "
    SAYZ WIT RESULT
KTHXBAI

HAI ME TEH FUNCSHUN SQUARE_ROOT
    SAYZ WIT "Enter number: "
    I HAS A VARIABLE NUM TEH DUBBLE ITZ GIMME AS DUBBLE

    MAYB
        IZ NUM SMALLR THAN 0.0?
            OOPS "Cannot calculate square root of negative number!"
        KTHX
        I HAS A VARIABLE RESULT TEH DUBBLE ITZ SQRT WIT NUM
        SAY WIT "√"
        SAY WIT NUM
        SAY WIT " = "
        SAYZ WIT RESULT
    OOPSIE ERR
        SAYZ WIT "Error: "
        SAYZ WIT ERR
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN TRIGONOMETRY
    SAYZ WIT "Enter angle in degrees: "
    I HAS A VARIABLE DEGREES TEH DUBBLE ITZ GIMME AS DUBBLE

    BTW Convert to radians (π/180 ≈ 0.017453)
    I HAS A VARIABLE RADIANS TEH DUBBLE ITZ DEGREES TIEMZ 0.017453292519943295

    I HAS A VARIABLE SINE_VAL TEH DUBBLE ITZ SIN WIT RADIANS
    I HAS A VARIABLE COSINE_VAL TEH DUBBLE ITZ COS WIT RADIANS

    SAY WIT "sin("
    SAY WIT DEGREES
    SAY WIT "°) = "
    SAYZ WIT SINE_VAL

    SAY WIT "cos("
    SAY WIT DEGREES
    SAY WIT "°) = "
    SAYZ WIT COSINE_VAL
KTHXBAI

HAI ME TEH FUNCSHUN RANDOM_NUMBER
    SAYZ WIT "1. Random float (0-1)"
    SAYZ WIT "2. Random integer (custom range)"
    SAY WIT "Choose: "
    I HAS A VARIABLE CHOICE TEH STRIN ITZ GIMME

    IZ CHOICE SAEM AS "1"?
        I HAS A VARIABLE RAND TEH DUBBLE ITZ RANDOM
        SAY WIT "Random number: "
        SAYZ WIT RAND
    KTHX

    IZ CHOICE SAEM AS "2"?
        SAYZ WIT "Enter minimum: "
        I HAS A VARIABLE MIN_VAL TEH INTEGR ITZ GIMME AS INTEGR

        SAYZ WIT "Enter maximum: "
        I HAS A VARIABLE MAX_VAL TEH INTEGR ITZ GIMME AS INTEGR

        I HAS A VARIABLE RAND_INT TEH INTEGR ITZ RANDINT WIT MIN_VAL AN WIT MAX_VAL
        SAY WIT "Random integer between "
        SAY WIT MIN_VAL
        SAY WIT " and "
        SAY WIT MAX_VAL LES 1
        SAY WIT ": "
        SAYZ WIT RAND_INT
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE RUNNING TEH BOOL ITZ YEZ

    WHILE RUNNING
        SHOW_MENU
        I HAS A VARIABLE CHOICE TEH STRIN ITZ GIMME

        IZ CHOICE SAEM AS "1"?
            BASIC_ARITHMETIC
        KTHX
        IZ CHOICE SAEM AS "2"?
            POWER_FUNCTION
        KTHX
        IZ CHOICE SAEM AS "3"?
            SQUARE_ROOT
        KTHX
        IZ CHOICE SAEM AS "4"?
            TRIGONOMETRY
        KTHX
        IZ CHOICE SAEM AS "5"?
            RANDOM_NUMBER
        KTHX
        IZ CHOICE SAEM AS "6"?
            SAYZ WIT "Goodbye!"
            RUNNING ITZ NO
        KTHX

        IZ (CHOICE SAEM AS "1") OR (CHOICE SAEM AS "2") OR (CHOICE SAEM AS "3") OR (CHOICE SAEM AS "4") OR (CHOICE SAEM AS "5") OR (CHOICE SAEM AS "6")?
            BTW Valid choice, do nothing
        NOPE
            SAYZ WIT "Invalid choice! Please select 1-6."
        KTHX
    KTHX
KTHXBAI
```

## Scientific Calculator Functions

```lol
BTW Scientific calculator with helper functions
I CAN HAS MATH?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEGREES_TO_RADIANS TEH DUBBLE WIT DEGREES TEH DUBBLE
    I HAS A VARIABLE PI TEH DUBBLE ITZ 3.141592653589793
    GIVEZ DEGREES TIEMZ PI DIVIDEZ 180.0
KTHXBAI

HAI ME TEH FUNCSHUN RADIANS_TO_DEGREES TEH DUBBLE WIT RADIANS TEH DUBBLE
    I HAS A VARIABLE PI TEH DUBBLE ITZ 3.141592653589793
    GIVEZ RADIANS TIEMZ 180.0 DIVIDEZ PI
KTHXBAI

HAI ME TEH FUNCSHUN FACTORIAL TEH INTEGR WIT N TEH INTEGR
    IZ N SMALLR THAN 0?
        OOPS "Cannot calculate factorial of negative number!"
    KTHX
    IZ N SMALLR THAN 2?
        GIVEZ 1
    NOPE
        GIVEZ N TIEMZ FACTORIAL WIT N LES 1
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN CALCULATE_HYPOTENUSE TEH DUBBLE WIT X TEH DUBBLE AN WIT Y TEH DUBBLE
    I HAS A VARIABLE X_SQ TEH DUBBLE ITZ X TIEMZ X
    I HAS A VARIABLE Y_SQ TEH DUBBLE ITZ Y TIEMZ Y
    I HAS A VARIABLE Z_SQ TEH DUBBLE ITZ X_SQ MOAR Y_SQ
    GIVEZ SQRT WIT Z_SQ
KTHXBAI

HAI ME TEH FUNCSHUN SCIENTIFIC_DEMO
    SAYZ WIT "=== Scientific Calculator Demo ==="
    SAYZ WIT ""

    BTW Demonstrate various calculations
    SAY WIT "sin(30°) = "
    I HAS A VARIABLE RADIANS TEH DUBBLE ITZ DEGREES_TO_RADIANS WIT 30.0
    SAYZ WIT SIN WIT RADIANS

    SAY WIT "cos(60°) = "
    I HAS A VARIABLE RAD60 TEH DUBBLE ITZ DEGREES_TO_RADIANS WIT 60.0
    SAYZ WIT COS WIT RAD60

    SAY WIT "5! = "
    SAYZ WIT FACTORIAL WIT 5

    SAY WIT "√144 = "
    SAYZ WIT SQRT WIT 144.0

    SAY WIT "2^10 = "
    SAYZ WIT POW WIT 2.0 AN WIT 10.0

    SAY WIT "Hypotenuse of triangle with sides 3,4: "
    SAYZ WIT CALCULATE_HYPOTENUSE WIT 3.0 AN WIT 4.0

    SAY WIT "Random percentage: "
    I HAS A VARIABLE PERCENT TEH DUBBLE ITZ RANDOM TIEMZ 100.0
    SAY WIT PERCENT
    SAYZ WIT "%"

    SAY WIT "Random dice roll: "
    SAYZ WIT RANDINT WIT 1 AN WIT 7
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    SCIENTIFIC_DEMO
KTHXBAI
```

## Running the Examples

Save any of the examples to a `.olol` file and run:

```bash
# Basic calculator
./olol basic_calculator.olol

# Advanced calculator
./olol advanced_calculator.olol

# Scientific demo
./olol scientific_calculator.olol
```

## Related Documentation

- [MATH Module](../standard-library/math.md) - Mathematical functions reference
- [STDIO Module](../standard-library/stdio.md) - Input/output functions
- [Functions](../language-guide/functions.md) - Function declaration and usage
- [Control Flow](../language-guide/control-flow.md) - Conditionals and loops