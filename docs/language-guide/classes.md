# Classes and Object-Oriented Programming

This guide covers classes, objects, inheritance, constructors, and visibility modifiers in Objective-LOL.

## Class Declaration

Classes are declared using `HAI ME TEH CLAS`:

```lol
HAI ME TEH CLAS <name> [KITTEH OF <parent> [AN OF <parent2> ...]]
    BTW class body
KTHXBAI
```

### Basic Class Example

```lol
HAI ME TEH CLAS PERSON
    EVRYONE    BTW Public visibility (default)
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
    DIS TEH VARIABLE AGE TEH INTEGR ITZ 0

    DIS TEH FUNCSHUN GET_NAME TEH STRIN
        GIVEZ NAME
    KTHX

    DIS TEH FUNCSHUN SET_AGE WIT NEW_AGE TEH INTEGR
        AGE ITZ NEW_AGE
    KTHX
KTHXBAI
```

## Visibility Modifiers

Control member access with visibility keywords:

| Modifier | Description | Access |
|----------|-------------|---------|
| `EVRYONE` | Public | Accessible from outside the class |
| `MAHSELF` | Private | Only accessible within the class |

### Visibility Example

```lol
HAI ME TEH CLAS BANK_ACCOUNT
    EVRYONE
    DIS TEH VARIABLE OWNER TEH STRIN ITZ "Anonymous"

    MAHSELF
    DIS TEH VARIABLE BALANCE TEH DUBBLE ITZ 0.0    BTW Private member

    EVRYONE
    DIS TEH FUNCSHUN DEPOSIT WIT AMOUNT TEH DUBBLE
        BALANCE ITZ BALANCE MOAR AMOUNT    BTW Can access private member
    KTHX

    DIS TEH FUNCSHUN GET_BALANCE TEH DUBBLE
        GIVEZ BALANCE
    KTHX
KTHXBAI
```

## Object Creation and Usage

### Creating Objects

```lol
BTW Without constructor
I HAS A VARIABLE PERSON1 TEH PERSON ITZ NEW PERSON

BTW With constructor (covered in next section)
I HAS A VARIABLE POINT1 TEH POINT ITZ NEW POINT WIT 10 AN WIT 20
```

### Accessing Members

```lol
HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE PERSON1 TEH PERSON ITZ NEW PERSON

    BTW Access member variables directly
    PERSON1 NAME ITZ "Alice"
    PERSON1 AGE ITZ 25

    BTW Read member variables
    SAYZ WIT PERSON1 NAME
    SAYZ WIT PERSON1 AGE

    BTW Call methods with DO
    PERSON1 DO SET_AGE WIT 26

    BTW Call methods with return values
    I HAS A VARIABLE CURRENT_NAME TEH STRIN ITZ PERSON1 DO GET_NAME
    SAYZ WIT CURRENT_NAME
KTHXBAI
```

## Constructors

Constructor methods are special methods that initialize objects:

- Must have the same name as the class (case-insensitive)
- Should not declare a return type (treated as void)
- Can accept parameters for initialization
- Called automatically during object instantiation

### Constructor Examples

```lol
HAI ME TEH CLAS POINT
    EVRYONE
    DIS TEH VARIABLE X TEH INTEGR ITZ 0
    DIS TEH VARIABLE Y TEH INTEGR ITZ 0

    BTW Constructor method - same name as class
    DIS TEH FUNCSHUN POINT WIT X_VAL TEH INTEGR AN WIT Y_VAL TEH INTEGR
        X ITZ X_VAL
        Y ITZ Y_VAL
        SAYZ WIT "Point created!"
    KTHX

    DIS TEH FUNCSHUN DISPLAY
        SAY WIT "Point("
        SAY WIT X
        SAY WIT ", "
        SAY WIT Y
        SAYZ WIT ")"
    KTHX
KTHXBAI

HAI ME TEH CLAS RECTANGLE
    EVRYONE
    DIS TEH VARIABLE WIDTH TEH INTEGR ITZ 0
    DIS TEH VARIABLE HEIGHT TEH INTEGR ITZ 0
    DIS TEH VARIABLE COLOR TEH STRIN ITZ "white"

    BTW Constructor with multiple parameters
    DIS TEH FUNCSHUN RECTANGLE WIT W TEH INTEGR AN WIT H TEH INTEGR AN WIT C TEH STRIN
        WIDTH ITZ W
        HEIGHT ITZ H
        COLOR ITZ C
    KTHX

    DIS TEH FUNCSHUN GET_AREA TEH INTEGR
        GIVEZ WIDTH TIEMZ HEIGHT
    KTHX
KTHXBAI
```

### Using Constructors

```lol
HAI ME TEH FUNCSHUN MAIN
    BTW Create objects with constructor arguments
    I HAS A VARIABLE ORIGIN TEH POINT ITZ NEW POINT WIT 0 AN WIT 0
    I HAS A VARIABLE CORNER TEH POINT ITZ NEW POINT WIT 10 AN WIT 5

    ORIGIN DO DISPLAY    BTW Point(0, 0)
    CORNER DO DISPLAY    BTW Point(10, 5)

    BTW Constructor with multiple parameters
    I HAS A VARIABLE RECT TEH RECTANGLE ITZ NEW RECTANGLE WIT 20 AN WIT 15 AN WIT "blue"
    SAYZ WIT RECT DO GET_AREA    BTW 300
KTHXBAI
```

### Parameterless Constructors

```lol
HAI ME TEH CLAS COUNTER
    EVRYONE
    DIS TEH VARIABLE VALUE TEH INTEGR ITZ 0

    BTW Parameterless constructor - called automatically
    DIS TEH FUNCSHUN COUNTER
        VALUE ITZ 1
        SAYZ WIT "Counter initialized!"
    KTHX

    DIS TEH FUNCSHUN GET_VALUE TEH INTEGR
        GIVEZ VALUE
    KTHX
KTHXBAI

BTW Usage
I HAS A VARIABLE COUNTER1 TEH COUNTER ITZ NEW COUNTER
BTW Output: "Counter initialized!"
SAYZ WIT COUNTER1 DO GET_VALUE    BTW 1
```

## Class Inheritance

### Single Inheritance

Use `KITTEH OF` for inheritance:

```lol
BTW Base class
HAI ME TEH CLAS ANIMAL
    EVRYONE
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"
    DIS TEH VARIABLE SPECIES TEH STRIN ITZ "Unknown"

    DIS TEH FUNCSHUN MAKE_SOUND
        SAYZ WIT "Some generic animal sound"
    KTHX
KTHXBAI

BTW Derived class
HAI ME TEH CLAS DOG KITTEH OF ANIMAL
    EVRYONE
    DIS TEH VARIABLE BREED TEH STRIN ITZ "Mixed"

    BTW Override parent method
    DIS TEH FUNCSHUN MAKE_SOUND
        SAYZ WIT "Woof!"
    KTHX

    DIS TEH FUNCSHUN WAG_TAIL
        SAYZ WIT "Wagging tail happily!"
    KTHX
KTHXBAI
```

### Multiple Inheritance

Objective-LOL supports multiple inheritance using `AN OF` to specify additional parent classes:

```lol
BTW Base classes
HAI ME TEH CLAS FLYER
    EVRYONE
    DIS TEH FUNCSHUN FLY
        SAYZ WIT "Flying through the air!"
    KTHX
KTHXBAI

HAI ME TEH CLAS SWIMMER
    EVRYONE
    DIS TEH FUNCSHUN SWIM
        SAYZ WIT "Swimming in the water!"
    KTHX
KTHXBAI

BTW Multiple inheritance - inherits from both ANIMAL and FLYER and SWIMMER
HAI ME TEH CLAS DUCK KITTEH OF ANIMAL AN OF FLYER AN OF SWIMMER
    EVRYONE
    DIS TEH FUNCSHUN MAKE_SOUND
        SAYZ WIT "Quack!"
    KTHX

    DIS TEH FUNCSHUN DIVE
        SAYZ WIT "Diving underwater!"
    KTHX
KTHXBAI
```

### Method Resolution Order (MRO)

When multiple parent classes have methods with the same name, Objective-LOL uses **Method Resolution Order** to determine which method to call:

**MRO Rules:**
1. **Depth-First**: Search the inheritance hierarchy in depth-first order
2. **Left-to-Right**: Among multiple parents, search from left to right
3. **Child First**: Child class methods override parent methods

```lol
HAI ME TEH CLAS BASE
    EVRYONE
    DIS TEH FUNCSHUN COMMON_METHOD
        SAYZ WIT "From BASE"
    KTHX
KTHXBAI

HAI ME TEH CLAS MIXIN_A KITTEH OF BASE
    EVRYONE
    DIS TEH FUNCSHUN COMMON_METHOD
        SAYZ WIT "From MIXIN_A"
    KTHX
KTHXBAI

HAI ME TEH CLAS MIXIN_B KITTEH OF BASE
    EVRYONE
    DIS TEH FUNCSHUN COMMON_METHOD
        SAYZ WIT "From MIXIN_B"
    KTHX
KTHXBAI

BTW Multiple inheritance with method conflicts
HAI ME TEH CLAS CHILD KITTEH OF MIXIN_A AN OF MIXIN_B
    EVRYONE
    BTW No override - which COMMON_METHOD gets called?
    BTW MRO: CHILD -> MIXIN_A -> BASE -> MIXIN_B -> BASE
    BTW Result: MIXIN_A's method is called (leftmost parent wins)
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_MRO
    I HAS A VARIABLE OBJ TEH CHILD ITZ NEW CHILD
    OBJ DO COMMON_METHOD  BTW Output: "From MIXIN_A"
KTHXBAI
```

### Complex Inheritance Hierarchy

```lol
BTW Diamond inheritance pattern
HAI ME TEH CLAS VEHICLE
    EVRYONE
    DIS TEH VARIABLE MAX_SPEED TEH INTEGR ITZ 0

    DIS TEH FUNCSHUN START_ENGINE
        SAYZ WIT "Engine starting..."
    KTHX
KTHXBAI

HAI ME TEH CLAS LAND_VEHICLE KITTEH OF VEHICLE
    EVRYONE
    DIS TEH VARIABLE WHEELS TEH INTEGR ITZ 4

    DIS TEH FUNCSHUN DRIVE
        SAYZ WIT "Driving on land"
    KTHX
KTHXBAI

HAI ME TEH CLAS WATER_VEHICLE KITTEH OF VEHICLE
    EVRYONE
    DIS TEH VARIABLE DISPLACEMENT TEH DUBBLE ITZ 0.0

    DIS TEH FUNCSHUN SAIL
        SAYZ WIT "Sailing on water"
    KTHX
KTHXBAI

BTW Amphibious vehicle inherits from both land and water vehicles
HAI ME TEH CLAS AMPHIBIOUS_VEHICLE KITTEH OF LAND_VEHICLE AN OF WATER_VEHICLE
    EVRYONE
    DIS TEH FUNCSHUN TRANSITION
        SAYZ WIT "Transitioning between land and water"
    KTHX

    BTW Can use methods from all parent classes
    DIS TEH FUNCSHUN DEMO_CAPABILITIES
        START_ENGINE    BTW From VEHICLE (via both paths)
        DRIVE           BTW From LAND_VEHICLE
        SAIL            BTW From WATER_VEHICLE
        TRANSITION      BTW Own method
    KTHX
KTHXBAI
```

### Using Inherited Classes

```lol
HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE MY_DOG TEH DOG ITZ NEW DOG
    MY_DOG NAME ITZ "Buddy"              BTW Inherited from ANIMAL
    MY_DOG SPECIES ITZ "Canine"          BTW Inherited from ANIMAL
    MY_DOG BREED ITZ "Golden Retriever"  BTW Specific to DOG

    MY_DOG DO MAKE_SOUND    BTW Calls overridden method: "Woof!"
    MY_DOG DO WAG_TAIL      BTW Calls dog-specific method
KTHXBAI
```

## Advanced Constructor Pattern

Constructor with validation and initialization:

```lol
HAI ME TEH CLAS BANK_ACCOUNT
    EVRYONE
    DIS TEH VARIABLE OWNER TEH STRIN ITZ "Anonymous"

    MAHSELF
    DIS TEH VARIABLE BALANCE TEH DUBBLE ITZ 0.0

    EVRYONE
    BTW Constructor with validation
    DIS TEH FUNCSHUN BANK_ACCOUNT WIT OWNER_NAME TEH STRIN AN WIT INITIAL_BALANCE TEH DUBBLE
        OWNER ITZ OWNER_NAME
        IZ INITIAL_BALANCE BIGGR THAN 0.0?
            BALANCE ITZ INITIAL_BALANCE
        NOPE
            BALANCE ITZ 0.0
            SAYZ WIT "Warning: Initial balance cannot be negative, set to 0"
        KTHX
    KTHX

    DIS TEH FUNCSHUN GET_BALANCE TEH DUBBLE
        GIVEZ BALANCE
    KTHX

    DIS TEH FUNCSHUN DEPOSIT WIT AMOUNT TEH DUBBLE
        IZ AMOUNT BIGGR THAN 0.0?
            BALANCE ITZ BALANCE MOAR AMOUNT
        KTHX
    KTHX
KTHXBAI
```

## Quick Reference

| Concept | Syntax |
|---------|--------|
| Class Declaration | `HAI ME TEH CLAS name ... KTHXBAI` |
| Single Inheritance | `HAI ME TEH CLAS child KITTEH OF parent` |
| Multiple Inheritance | `HAI ME TEH CLAS child KITTEH OF parent1 AN OF parent2` |
| Member Variable | `DIS TEH VARIABLE name TEH type` |
| Member Method | `DIS TEH FUNCSHUN name ...` |
| Constructor | `DIS TEH FUNCSHUN classname WIT params...` |
| Object Creation | `NEW classname [WIT args...]` |
| Method Call | `object DO method WIT args` |
| Member Access | `object member` |

## Next Steps

- [Modules](modules.md) - Code organization across files
- [Standard Library Collections](../standard-library/collections.md) - Built-in BUKKIT and BASKIT classes