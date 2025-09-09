# Standard Library Overview

The Objective-LOL standard library provides essential functionality through a collection of built-in modules.

## Available Modules

### Core Modules

| Module | Description | Key Items |
|--------|-------------|-----------|
| [STDIO](stdio.md) | Input/Output functions | `SAY`, `SAYZ`, `GIMME` |
| [MATH](math.md) | Mathematical functions | `ABS`, `MAX`, `MIN`, `SQRT`, `POW`, `SIN`, `COS`, `LOG`, `PI`, `E` |
| [RANDOM](random.md) | Random number generation | `RANDOM_FLOAT`, `RANDOM_INT`, `RANDOM_BOOL`, `UUID` |
| [TIME](time.md) | Date and time functionality | `DATE` class, `SLEEP` function |
| [STRING](string.md) | String utility functions | `LEN`, `CONCAT`, `TRIM`, `UPPER`, `LOWER`, `SPLIT` |
| [TEST](test.md) | Testing and assertions | `ASSERT` |
| [IO](io.md) | Advanced I/O classes | `READER`, `WRITER`, `BUFFERED_READER`, `BUFFERED_WRITER` |
| [THREAD](threading.md) | Concurrency support | `YARN` thread class, `KNOT` mutex class |
| [FILE](file.md) | File system operations | `DOCUMENT` class with read/write operations |
| [CACHE](cache.md) | Caching classes | `MEMSTASH` LRU cache, `TIMESTASH` TTL cache |
| [HTTP](http.md) | HTTP client operations | `INTERWEB` client class, `RESPONSE` class |
| [SOCKET](socket.md) | Network socket operations | `SOKKIT` socket class, `WIRE` connection class |
| [PROCESS](process.md) | Process management operations | `MINION` process class, `PIPE` I/O class |

### Collection Types

Built-in collection types don't require imports:

| Type | Description | Usage |
|------|-------------|-------|
| [BUKKIT](collections.md) | Dynamic arrays | `NEW BUKKIT` |
| [BASKIT](collections.md) | Maps/dictionaries | `NEW BASKIT` |

## Import Patterns

### Common Usage Patterns

```lol
BTW Pattern 1: Global standard imports
I CAN HAS STDIO?
I CAN HAS MATH?

BTW Pattern 2: Selective global imports
I CAN HAS SAY AN SAYZ FROM STDIO?
I CAN HAS ABS AN MAX AN SQRT FROM MATH?

BTW Pattern 3: Function-specific imports
HAI ME TEH FUNCSHUN PROCESS_DATA
    I CAN HAS TIME?    BTW Only needed in this function
    I HAS A VARIABLE NOW_DATE TEH DATE ITZ NEW DATE
KTHXBAI

BTW Pattern 4: Mixed imports
I CAN HAS SAYZ FROM STDIO?           BTW Global selective
I CAN HAS MATH?                      BTW Global full

HAI ME TEH FUNCSHUN ADVANCED_CALC
    I CAN HAS DATE FROM TIME?        BTW Local selective
    BTW Can use SAYZ (global), all MATH (global), DATE (local)
KTHXBAI

BTW Pattern 5: Network operation example
I CAN HAS STDIO?                     BTW For output
I CAN HAS HTTP?                      BTW For web requests

HAI ME TEH FUNCSHUN FETCH_API_DATA WIT URL TEH STRIN
    I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
    I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL
    SAYZ WIT RESPONSE BODY
KTHXBAI
```
