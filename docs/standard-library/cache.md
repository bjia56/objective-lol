# CACHE Module - Caching Classes

The CACHE module provides high-performance caching classes for storing and retrieving key-value pairs with different eviction strategies.

## Import

```lol
BTW Full import
I CAN HAS CACHE?

BTW Selective import examples
I CAN HAS MEMSTASH FROM CACHE?
I CAN HAS TIMESTASH FROM CACHE?
I CAN HAS STASH AN MEMSTASH FROM CACHE?
```

## Cache Classes

### STASH - Base Cache Class

Abstract base class for all cache implementations. Provides common interface for cache operations.

**Classes:** Base class for `MEMSTASH` and `TIMESTASH`

### MEMSTASH - LRU Cache

An LRU (Least Recently Used) cache with fixed capacity. When the cache reaches capacity, the least recently used items are evicted.

**Constructor:** `NEW MEMSTASH WIT <capacity>`
**Parameters:**
- `capacity` (INTEGR) - Maximum number of items to store

**Properties:**
- `SIZ` (INTEGR) - Current number of items in cache

```lol
I CAN HAS MEMSTASH FROM CACHE?

I HAS A VARIABLE LRU_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
SAYZ WIT LRU_CACHE SIZ    BTW 0
```

#### PUT - Store Key-Value Pair

Stores a key-value pair in the cache. Updates existing keys and moves them to most recently used position.

**Syntax:** `cache DO PUT WIT <key> AN WIT <value>`
**Parameters:**
- `key` (STRIN) - Cache key
- `value` (STRIN) - Cache value

**Returns:** NOTHIN

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 3

CACHE DO PUT WIT "user:123" AN WIT "John Doe"
CACHE DO PUT WIT "user:456" AN WIT "Jane Smith"
CACHE DO PUT WIT "user:789" AN WIT "Bob Wilson"
```

#### GET - Retrieve Value

Retrieves a value by key and marks it as recently used.

**Syntax:** `cache DO GET WIT <key>`
**Parameters:**
- `key` (STRIN) - Cache key to retrieve

**Returns:** STRIN or NOTHIN if key not found

```lol
I HAS A VARIABLE USER_NAME TEH STRIN ITZ CACHE DO GET WIT "user:123"
IZ USER_NAME SAEM AS NOTHIN?
    SAYZ WIT "User not found"
NOPE
    SAYZ WIT USER_NAME
KTHX
```

#### CONTAINS - Check Key Existence

Checks if a key exists in the cache without affecting its position.

**Syntax:** `cache DO CONTAINS WIT <key>`
**Parameters:**
- `key` (STRIN) - Cache key to check

**Returns:** BOOL (YEZ if exists, NO otherwise)

```lol
IZ CACHE DO CONTAINS WIT "user:123"?
    SAYZ WIT "User exists in cache"
NOPE
    SAYZ WIT "User not cached"
KTHX
```

#### DELETE - Remove Key

Removes a key-value pair from the cache.

**Syntax:** `cache DO DELETE WIT <key>`
**Parameters:**
- `key` (STRIN) - Cache key to remove

**Returns:** BOOL (YEZ if deleted, NO if key not found)

```lol
I HAS A VARIABLE DELETED TEH BOOL ITZ CACHE DO DELETE WIT "user:123"
IZ DELETED?
    SAYZ WIT "User removed from cache"
NOPE
    SAYZ WIT "User not found in cache"
KTHX
```

#### CLEAR - Remove All Items

Removes all items from the cache.

**Syntax:** `cache DO CLEAR`
**Returns:** NOTHIN

```lol
CACHE DO CLEAR
SAYZ WIT CACHE SIZ    BTW 0
```

### TIMESTASH - TTL Cache

A TTL (Time To Live) cache where items expire after a specified duration.

**Constructor:** `NEW TIMESTASH WIT <ttl_seconds>`
**Parameters:**
- `ttl_seconds` (INTEGR) - Time to live in seconds

**Properties:**
- `SIZ` (INTEGR) - Current number of non-expired items

```lol
I CAN HAS TIMESTASH FROM CACHE?

BTW Cache items for 5 minutes (300 seconds)
I HAS A VARIABLE TTL_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
```

TIMESTASH supports the same methods as MEMSTASH (`PUT`, `GET`, `CONTAINS`, `DELETE`, `CLEAR`) with identical syntax. The key difference is automatic expiration based on time rather than capacity.

```lol
I HAS A VARIABLE SESSION_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60

BTW Store session data
SESSION_CACHE DO PUT WIT "session:abc123" AN WIT "user_id:789"
SESSION_CACHE DO PUT WIT "session:def456" AN WIT "user_id:456"

BTW Retrieve session (works if not expired)
I HAS A VARIABLE SESSION_DATA TEH STRIN ITZ SESSION_CACHE DO GET WIT "session:abc123"
```

## Usage Examples

### User Session Cache

```lol
I CAN HAS CACHE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN CREATE_SESSION_CACHE
    BTW 30 minute sessions
    GIVEZ NEW TIMESTASH WIT 1800
KTHXBAI

HAI ME TEH FUNCSHUN CACHE_USER_SESSION WIT CACHE TEH TIMESTASH AN WIT SESSION_ID TEH STRIN AN WIT USER_ID TEH STRIN
    I HAS A VARIABLE KEY TEH STRIN ITZ "session:" SMOOSH SESSION_ID
    CACHE DO PUT WIT KEY AN WIT USER_ID
    SAYZ WIT "Session cached"
KTHXBAI

HAI ME TEH FUNCSHUN GET_USER_FROM_SESSION WIT CACHE TEH TIMESTASH AN WIT SESSION_ID TEH STRIN
    I HAS A VARIABLE KEY TEH STRIN ITZ "session:" SMOOSH SESSION_ID
    GIVEZ CACHE DO GET WIT KEY
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE SESSION_CACHE TEH TIMESTASH ITZ CREATE_SESSION_CACHE
    
    CACHE_USER_SESSION WIT SESSION_CACHE AN WIT "abc123" AN WIT "user789"
    
    I HAS A VARIABLE USER_ID TEH STRIN ITZ GET_USER_FROM_SESSION WIT SESSION_CACHE AN WIT "abc123"
    IZ USER_ID SAEM AS NOTHIN?
        SAYZ WIT "Session expired or not found"
    NOPE
        SAY WIT "Active session for user: "
        SAYZ WIT USER_ID
    KTHX
KTHXBAI
```

### Data Processing Cache

```lol
I CAN HAS CACHE?
I CAN HAS STDIO?
I CAN HAS STRING?

HAI ME TEH FUNCSHUN PROCESS_WITH_CACHE WIT DATA_CACHE TEH MEMSTASH AN WIT INPUT TEH STRIN
    BTW Check cache first
    I HAS A VARIABLE CACHED_RESULT TEH STRIN ITZ DATA_CACHE DO GET WIT INPUT
    IZ CACHED_RESULT SAEM AS NOTHIN?
        BTW Cache miss - process data
        SAYZ WIT "Processing data..."
        I HAS A VARIABLE RESULT TEH STRIN ITZ UPPER WIT INPUT
        DATA_CACHE DO PUT WIT INPUT AN WIT RESULT
        GIVEZ RESULT
    NOPE
        BTW Cache hit
        SAYZ WIT "Found in cache!"
        GIVEZ CACHED_RESULT
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN MAIN
    BTW Create cache with capacity of 1000 items
    I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000
    
    BTW Process same data multiple times
    I HAS A VARIABLE RESULT1 TEH STRIN ITZ PROCESS_WITH_CACHE WIT CACHE AN WIT "hello"
    I HAS A VARIABLE RESULT2 TEH STRIN ITZ PROCESS_WITH_CACHE WIT CACHE AN WIT "hello"
    
    SAYZ WIT RESULT1    BTW HELLO
    SAYZ WIT RESULT2    BTW HELLO (from cache)
KTHXBAI
```

### Mixed Cache Strategy

```lol
I CAN HAS CACHE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN MAIN
    BTW LRU cache for frequently accessed data
    I HAS A VARIABLE FREQUENT_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 50
    
    BTW TTL cache for temporary data (1 hour)
    I HAS A VARIABLE TEMP_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 3600
    
    BTW Store different types of data
    FREQUENT_CACHE DO PUT WIT "config:theme" AN WIT "dark"
    TEMP_CACHE DO PUT WIT "download:temp123" AN WIT "processing"
    
    BTW Check both caches
    I HAS A VARIABLE THEME TEH STRIN ITZ FREQUENT_CACHE DO GET WIT "config:theme"
    I HAS A VARIABLE DOWNLOAD_STATUS TEH STRIN ITZ TEMP_CACHE DO GET WIT "download:temp123"
    
    SAYZ WIT THEME
    SAYZ WIT DOWNLOAD_STATUS
    
    BTW Show cache sizes
    SAY WIT "Frequent cache size: "
    SAYZ WIT FREQUENT_CACHE SIZ
    SAY WIT "Temp cache size: "
    SAYZ WIT TEMP_CACHE SIZ
KTHXBAI
```

## Class Summary

### MEMSTASH (LRU Cache)

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `NEW MEMSTASH` | capacity | MEMSTASH | Create LRU cache with fixed capacity |
| `PUT` | key, value | NOTHIN | Store key-value pair |
| `GET` | key | STRIN/NOTHIN | Retrieve value by key |
| `CONTAINS` | key | BOOL | Check if key exists |
| `DELETE` | key | BOOL | Remove key-value pair |
| `CLEAR` | none | NOTHIN | Remove all items |

| Property | Type | Description |
|----------|------|-------------|
| `SIZ` | INTEGR | Current number of items |

### TIMESTASH (TTL Cache)

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `NEW TIMESTASH` | ttl_seconds | TIMESTASH | Create TTL cache with expiration time |
| `PUT` | key, value | NOTHIN | Store key-value pair with TTL |
| `GET` | key | STRIN/NOTHIN | Retrieve non-expired value by key |
| `CONTAINS` | key | BOOL | Check if non-expired key exists |
| `DELETE` | key | BOOL | Remove key-value pair |
| `CLEAR` | none | NOTHIN | Remove all items |

| Property | Type | Description |
|----------|------|-------------|
| `SIZ` | INTEGR | Current number of non-expired items |

## Related

- [STRING Module](string.md) - For key manipulation
- [STDIO Module](stdio.md) - For cache debugging output
- [TIME Module](time.md) - For time-based operations