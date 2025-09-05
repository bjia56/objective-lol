# Objective-LOL Quick Reference

**Essential Syntax Lookup** · [📚 Complete Docs](README.md) · [🔍 Keywords](reference/keywords.md)

## Hello World

```lol
HAI ME TEH FUNCSHUN MAIN
    I CAN HAS STDIO?
    SAYZ WIT "Hello, World!"
KTHXBAI
```

## Data Types

| Type | Literal Examples | Usage |
|------|------------------|-------|
| `INTEGR` | `42`, `0xFF`, `-17` | 64-bit integer |
| `DUBBLE` | `3.14`, `-2.5` | 64-bit float |
| `STRIN` | `"Hello"`, `"Line\n"` | UTF-8 string |
| `BOOL` | `YEZ`, `NO` | Boolean values |
| `NOTHIN` | `NOTHIN` | Null/void |
| `BUKKIT` | `NEW BUKKIT` | Dynamic array |
| `BASKIT` | `NEW BASKIT` | Key-value map |

## Variables & Assignment

### Local

```lol
I HAS A VARIABLE NAME TEH STRIN ITZ "Alice"    BTW Declaration + init
NAME ITZ "Bob"                                 BTW Assignment
I HAS A LOCKD VARIABLE PI TEH DUBBLE ITZ 3.14  BTW Constant
```

### Global

```lol
HAI ME TEH VARIABLE NAME TEH STRIN ITZ "Charlie"
HAI ME TEH LOCKD VARIABLE E TEH DUBBLE ITZ 2.72
```

## Operators

| Category | Operators | Examples |
|----------|-----------|----------|
| **Arithmetic** | `MOAR` `LES` `TIEMZ` `DIVIDEZ` | `5 MOAR 3` → 8 |
| **Comparison** | `BIGGR THAN` `SMALLR THAN` `SAEM AS` | `X BIGGR THAN 5` |
| **Logical** | `AN` `OR` | `YEZ AN NO` → NO |
| **Casting** | `AS` | `"42" AS INTEGR` |
| **Grouping** | `( )` | `(2 MOAR 3) TIEMZ 4` |

**Precedence:** `()` > `AS` > `*/` > `+-` > `<>` > `==` > `AN` > `OR`

## Control Flow

| Pattern | Syntax | Example |
|---------|--------|---------|
| **If/Else** | `IZ condition? ... NOPE ... KTHX` | `IZ AGE BIGGR THAN 17? ... KTHX` |
| **Loop** | `WHILE condition ... KTHX` | `WHILE COUNT BIGGR THAN 0 ... KTHX` |
| **Exception** | `MAYB ... OOPSIE var ... ALWAYZ ... KTHX` | `MAYB ... OOPSIE ERR ... KTHX` |

```lol
IZ AGE BIGGR THAN 17?
    SAYZ WIT "Adult"
NOPE
    SAYZ WIT "Minor"
KTHX

WHILE COUNT BIGGR THAN 0
    COUNT ITZ COUNT LES 1
KTHX

MAYB
    RESULT ITZ 10.0 DIVIDEZ 0.0
OOPSIE ERR
    SAYZ WIT ERR
KTHX
```

## Functions

```lol
BTW Function definition
HAI ME TEH FUNCSHUN ADD TEH INTEGR WIT A TEH INTEGR AN WIT B TEH INTEGR
    GIVEZ A MOAR B
KTHXBAI

BTW Function call
I HAS A VARIABLE SUM TEH INTEGR ITZ ADD WIT 5 AN WIT 3
```

## Classes

```lol
HAI ME TEH CLAS PERSON
    EVRYONE BTW Start of public members
    DIS TEH VARIABLE NAME TEH STRIN ITZ "Unknown"

    MAHSELF BTW Start of private members
    DIS TEH VARIABLE AGE TEH INTEGR ITZ 0

    EVRYONE BTW Switch back to public

    BTW Constructor
    DIS TEH FUNCSHUN PERSON WIT N TEH STRIN AN WIT A TEH INTEGR
        NAME ITZ N
        AGE ITZ A
    KTHX

    DIS TEH FUNCSHUN GREET
        SAYZ WIT NAME
    KTHX
KTHXBAI

I HAS A VARIABLE P TEH PERSON ITZ NEW PERSON WIT "Alice" AN WIT 25
P DO GREET
```

## Collections

| Type | Common Methods | Usage |
|------|----------------|-------|
| **BUKKIT** | `PUSH` `POP` `AT` `SET` `SIZ` `JOIN` | `ARR DO PUSH WIT 10` |
| **BASKIT** | `PUT` `GET` `CONTAINS` `KEYS` `SIZ` | `MAP DO PUT WIT "key" AN WIT "val"` |

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT
ARR DO PUSH WIT 10
SAYZ WIT ARR DO AT WIT 0    BTW 10

I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
MAP DO PUT WIT "name" AN WIT "Alice"
SAYZ WIT MAP DO GET WIT "name"    BTW Alice
```

## Imports

| Import Type | Syntax | Example |
|-------------|--------|---------|
| **Built-in** | `I CAN HAS MODULE?` | `I CAN HAS STDIO?` |
| **Selective** | `I CAN HAS FUNC FROM MODULE?` | `I CAN HAS SAY FROM STDIO?` |
| **File** | `I CAN HAS "filename"?` | `I CAN HAS "utils"?` |

**Standard Modules:**
- `STDIO` → `SAY` `SAYZ` `GIMME`
- `MATH` → `ABS` `MAX` `SQRT` `RANDOM` `SIN` `COS`
- `TIME` → `DATE` class, `SLEEP` function
- `STRING` → `LEN` `CONCAT`
- `TEST` → `ASSERT`
- `IO` → `READER` `WRITER` `READERWRITER` `BUFFERED_READER` `BUFFERED_WRITER`
- `THREAD` → `YARN` `KNOT`
- `FILE` → `DOCUMENT`

## Common Patterns

**Input/Output**
```lol
I CAN HAS STDIO?
SAYZ WIT "Enter name: "
I HAS A VARIABLE NAME TEH STRIN ITZ GIMME
SAYZ WIT NAME
```

**Type Conversion**
```lol
I HAS A VARIABLE NUM TEH INTEGR ITZ "42" AS INTEGR
I HAS A VARIABLE STR TEH STRIN ITZ 42 AS STRIN
```

**Math & Random**
```lol
I CAN HAS MATH?
I HAS A VARIABLE ROOT TEH DUBBLE ITZ SQRT WIT 16.0
I HAS A VARIABLE DICE TEH INTEGR ITZ RANDINT WIT 1 AN WIT 7
```

## Key Rules

- **File Extension**: Must use `.olol` (not `.lol`)
- **Case Insensitive**: Keywords work in any case
- **Entry Point**: `MAIN` function required
- **Block Endings**: `KTHX` ends blocks, `KTHXBAI` ends functions/classes

## Build & Run

```bash
go build -o olol cmd/olol/main.go    # Build interpreter
./olol program.olol                  # Run program
```

---

**📚 Complete Docs:** [docs/README.md](docs/README.md) • **🔍 All Operators:** [docs/reference/operators.md](docs/reference/operators.md) • **💡 Examples:** [docs/examples/](docs/examples/)