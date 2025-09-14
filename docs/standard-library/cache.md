# CACHE Module

## Import

```lol
BTW Full import
I CAN HAS CACHE?

BTW Selective import examples
```

## Cache Interface

### STASH Class

Abstract base class for all cache implementations in the CACHE module.
Provides a common interface for cache operations like storing, retrieving, and managing cached data.

**Methods:**

#### CLEAR

Removes all items from the cache.
Resets the cache to empty state.

**Syntax:** `<cache> DO CLEAR`
**Example: Clear entire cache**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
BTW ... add many items ...
CACHE DO CLEAR BTW Remove all items
SAYZ WIT "Cache cleared, size: " MOAR CACHE SIZ
```

**Example: Reset between tests**

```lol
I HAS A VARIABLE TEST_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60
BTW ... run test that adds items ...
TEST_CACHE DO CLEAR BTW Clean up for next test
```

**Example: Emergency cleanup**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000
BTW ... cache grows too large ...
CACHE DO CLEAR BTW Emergency memory cleanup
```

**Note:** Completely empties the cache

**Note:** Cannot be undone - use with caution

**Note:** Thread-safe in concrete implementations

**Note:** Size becomes 0 after clearing

#### CONTAINS

Checks if a key exists in the cache.
Returns YEZ if exists and not expired (in TTL caches), NO otherwise.

**Syntax:** `<cache> DO CONTAINS WIT <key>`
**Parameters:**
- `key` (STRIN): The cache key to check

**Example: Check key existence**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "session" AN WIT "abc123"
IZ CACHE DO CONTAINS WIT "session"?
SAYZ WIT "Session exists"
KTHX
```

**Example: Safe operations**

```lol
IZ CACHE DO CONTAINS WIT "user_data"?
I HAS A VARIABLE DATA TEH STRIN ITZ CACHE DO GET WIT "user_data"
BTW ... process data ...
NOPE
SAYZ WIT "No cached user data"
KTHX
```

**Example: Cache statistics**

```lol
I HAS A VARIABLE TOTAL_KEYS TEH NUMBR ITZ 0
BTW ... iterate through possible keys ...
IZ CACHE DO CONTAINS WIT "key1"?
TOTAL_KEYS ITZ TOTAL_KEYS MOAR 1
KTHX
```

**Note:** Does not affect cache ordering (unlike GET)

**Note:** In TTL caches, automatically cleans up expired entries

**Note:** Thread-safe in concrete implementations

#### DELETE

Removes a key-value pair from the cache.
Returns YEZ if deleted, NO if key not found.

**Syntax:** `<cache> DO DELETE WIT <key>`
**Parameters:**
- `key` (STRIN): The cache key to remove

**Example: Remove a cached item**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "temp" AN WIT "temporary_data"
I HAS A VARIABLE DELETED TEH BOOL ITZ CACHE DO DELETE WIT "temp"
IZ DELETED?
SAYZ WIT "Item removed successfully"
KTHX
```

**Example: Cleanup expired sessions**

```lol
I HAS A VARIABLE SESSION_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 3600
BTW ... store session data ...
SESSION_CACHE DO DELETE WIT "session_123" BTW Force logout
```

**Example: Batch cleanup**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
BTW ... store multiple items ...
CACHE DO DELETE WIT "item1"
CACHE DO DELETE WIT "item2"
CACHE DO DELETE WIT "item3"
```

**Note:** Returns NO if key doesn't exist (not an error)

**Note:** In TTL caches, works regardless of expiration status

**Note:** Thread-safe in concrete implementations

#### GET

Retrieves a value by key from the cache.
Returns NOTHIN if key not found or expired (in TTL caches).

**Syntax:** `<cache> DO GET WIT <key>`
**Parameters:**
- `key` (STRIN): The cache key to retrieve

**Example: Retrieve a value**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "user" AN WIT "alice"
I HAS A VARIABLE USER TEH STRIN ITZ CACHE DO GET WIT "user"
SAYZ WIT "User: " MOAR USER
```

**Example: Handle missing keys**

```lol
I HAS A VARIABLE VALUE TEH STRIN ITZ CACHE DO GET WIT "nonexistent"
IZ VALUE SAEM AS NOTHIN?
SAYZ WIT "Key not found"
KTHX
```

**Example: Safe retrieval pattern**

```lol
IZ CACHE DO CONTAINS WIT "config"?
I HAS A VARIABLE CONFIG TEH STRIN ITZ CACHE DO GET WIT "config"
BTW ... use config ...
KTHX
```

**Note:** Returns NOTHIN for missing keys

**Note:** In TTL caches, also returns NOTHIN for expired entries

**Note:** May affect cache ordering (LRU caches move accessed items to front)

#### PUT

Stores a key-value pair in the cache.
Overwrites existing keys with new values.

**Syntax:** `<cache> DO PUT WIT <key> AN WIT <value>`
**Parameters:**
- `key` (STRIN): The cache key to store
- `value` (STRIN): The value to associate with the key

**Example: Store a value**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "username" AN WIT "john_doe"
```

**Example: Update existing value**

```lol
CACHE DO PUT WIT "counter" AN WIT "1"
CACHE DO PUT WIT "counter" AN WIT "2" BTW Overwrites previous value
```

**Example: Store configuration**

```lol
CACHE DO PUT WIT "db_host" AN WIT "localhost"
CACHE DO PUT WIT "db_port" AN WIT "5432"
```

**Note:** Overwrites existing keys without warning

**Note:** Implementation depends on cache type (LRU vs TTL behavior)

**Note:** Thread-safe in concrete implementations

**Member Variables:**

#### SIZ

Read-only property that returns the current number of items in the cache.


**Example: Check cache size**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "key1" AN WIT "value1"
CACHE DO PUT WIT "key2" AN WIT "value2"
SAYZ WIT "Cache size: " MOAR CACHE SIZ
```

**Example: Monitor cache growth**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
BTW ... add items over time ...
IZ CACHE SIZ BIGGR THAN 50?
SAYZ WIT "Cache getting large, consider cleanup"
KTHX
```

**Example: Empty cache check**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
IZ CACHE SIZ SAEM AS 0?
SAYZ WIT "Cache is empty"
KTHX
```

**Note:** For TTL caches, this may trigger cleanup of expired entries

**Note:** Always reflects current state (may change between calls)

**Note:** Thread-safe in concrete implementations

**Example: Basic cache usage pattern**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "key1" AN WIT "value1"
I HAS A VARIABLE VALUE TEH STRIN ITZ CACHE DO GET WIT "key1"
SAYZ WIT "Retrieved: " MOAR VALUE
```

**Example: Cache operations**

```lol
IZ CACHE DO CONTAINS WIT "key1"?
SAYZ WIT "Key exists"
KTHX
CACHE DO DELETE WIT "key1"
CACHE DO CLEAR
```

## Lru Cache

### MEMSTASH Class

An LRU (Least Recently Used) cache with fixed capacity.
When the cache reaches capacity, the least recently used item is evicted automatically.
Thread-safe implementation with proper concurrency support.

**Methods:**

#### CLEAR

Removes all items from the MEMSTASH.
Resets the cache to empty state while preserving capacity.

**Syntax:** `<memstash> DO CLEAR`
**Example: Clear entire cache**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
BTW ... add many items ...
SAYZ WIT "Size before clear: " MOAR CACHE SIZ
CACHE DO CLEAR
SAYZ WIT "Size after clear: " MOAR CACHE SIZ
```

**Example: Reset between operations**

```lol
I HAS A VARIABLE TEMP_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 50
BTW ... use for temporary calculations ...
TEMP_CACHE DO CLEAR BTW Clean slate for next operation
```

**Example: Emergency cleanup**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000
BTW ... cache grew too large or corrupted ...
CACHE DO CLEAR BTW Emergency reset
```

**Example: Memory management**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
BTW ... long-running application ...
IZ CACHE SIZ BIGGR THAN 80?
CACHE DO CLEAR BTW Prevent memory bloat
SAYZ WIT "Cache cleared to prevent memory issues"
KTHX
```

**Note:** Completely empties the cache

**Note:** Preserves the original capacity

**Note:** Cannot be undone - use with caution

**Note:** Thread-safe for concurrent operations

**Note:** Size becomes 0 after clearing

#### CONTAINS

Checks if a key exists in the MEMSTASH without affecting its position.
Does not move the item in the LRU ordering (unlike GET).
Returns YEZ if exists, NO otherwise.

**Syntax:** `<memstash> DO CONTAINS WIT <key>`
**Parameters:**
- `key` (STRIN): The cache key to check

**Example: Check without affecting LRU**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "data" AN WIT "value"
IZ CACHE DO CONTAINS WIT "data"?
SAYZ WIT "Key exists (not moved to front)"
KTHX
```

**Example: Safe access pattern**

```lol
IZ CACHE DO CONTAINS WIT "user_session"?
I HAS A VARIABLE SESSION TEH STRIN ITZ CACHE DO GET WIT "user_session"
BTW ... use session ...
NOPE
SAYZ WIT "Session not cached"
KTHX
```

**Example: Cache statistics**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
BTW ... add various items ...
I HAS A VARIABLE HIT_COUNT TEH NUMBR ITZ 0
I HAS A VARIABLE MISS_COUNT TEH NUMBR ITZ 0
IZ CACHE DO CONTAINS WIT "key1"?
HIT_COUNT ITZ HIT_COUNT MOAR 1
NOPE
MISS_COUNT ITZ MISS_COUNT MOAR 1
KTHX
```

**Note:** Does not affect LRU ordering

**Note:** Use GET if you want to access and update LRU position

**Note:** Thread-safe for concurrent operations

**Note:** Faster than GET for existence checks only

#### DELETE

Removes a key-value pair from the MEMSTASH.
Returns YEZ if deleted, NO if key not found.

**Syntax:** `<memstash> DO DELETE WIT <key>`
**Parameters:**
- `key` (STRIN): The cache key to remove

**Example: Remove specific item**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "temp" AN WIT "temporary_data"
I HAS A VARIABLE REMOVED TEH BOOL ITZ CACHE DO DELETE WIT "temp"
IZ REMOVED?
SAYZ WIT "Item removed successfully"
KTHX
```

**Example: Cleanup expired sessions**

```lol
I HAS A VARIABLE SESSIONS TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000
BTW ... store session data ...
SESSIONS DO DELETE WIT "expired_session_123"
```

**Example: Selective cache clearing**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
BTW ... cache has various data types ...
CACHE DO DELETE WIT "debug_info" BTW Remove debug data
CACHE DO DELETE WIT "temp_calc" BTW Remove temporary results
```

**Example: Handle deletion result**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
I HAS A VARIABLE WAS_DELETED TEH BOOL ITZ CACHE DO DELETE WIT "nonexistent"
IZ WAS_DELETED?
SAYZ WIT "Item was deleted"
NOPE
SAYZ WIT "Item didn't exist"
KTHX
```

**Note:** Returns NO if key doesn't exist (not an error)

**Note:** Thread-safe for concurrent operations

**Note:** Freed space can be used for new items

**Note:** Does not affect LRU ordering of remaining items

#### GET

Retrieves a value by key and marks it as recently used.
Moves the accessed item to the front (most recently used position).
Returns NOTHIN if key not found.

**Syntax:** `<memstash> DO GET WIT <key>`
**Parameters:**
- `key` (STRIN): The cache key to retrieve

**Example: Access with LRU update**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "data" AN WIT "value"
I HAS A VARIABLE VALUE TEH STRIN ITZ CACHE DO GET WIT "data" BTW Moves to front
SAYZ WIT "Got: " MOAR VALUE
```

**Example: LRU behavior demonstration**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 3
CACHE DO PUT WIT "a" AN WIT "1"
CACHE DO PUT WIT "b" AN WIT "2"
CACHE DO PUT WIT "c" AN WIT "3"
CACHE DO GET WIT "a" BTW 'a' now most recently used
CACHE DO PUT WIT "d" AN WIT "4" BTW Evicts 'b' (least recently used)
IZ CACHE DO CONTAINS WIT "b"?
SAYZ WIT "'b' still exists"
NOPE
SAYZ WIT "'b' was evicted, 'a' was saved by access"
KTHX
```

**Example: Cache warming**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000
BTW Pre-load frequently accessed items
CACHE DO PUT WIT "hot1" AN WIT "value1"
CACHE DO PUT WIT "hot2" AN WIT "value2"
BTW ... application runs, accessing items ...
CACHE DO GET WIT "hot1" BTW Keep this item fresh
```

**Note:** Moves accessed items to most recently used position

**Note:** Affects which items get evicted when capacity is reached

**Note:** Thread-safe for concurrent access

**Note:** Use CONTAINS to check existence without affecting LRU order

#### MEMSTASH

Initializes a MEMSTASH LRU cache with the specified capacity.
Capacity must be positive integer representing maximum number of items.

**Syntax:** `NEW MEMSTASH WIT <capacity>`
**Parameters:**
- `capacity` (INTEGR): Maximum number of items the cache can hold

**Example: Create small cache**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 10
```

**Example: Create large cache**

```lol
I HAS A VARIABLE BIG_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 10000
```

**Example: Handle invalid capacity**

```lol
BTW This would throw an exception:
BTW I HAS A VARIABLE BAD_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 0
```

**Example: Memory-conscious caching**

```lol
I HAS A VARIABLE MEM_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 50
BTW Small cache for memory-limited environments
```

**Note:** Capacity cannot be changed after creation

**Note:** Capacity must be positive (> 0)

**Note:** Memory usage scales with capacity

**Note:** Consider available memory when setting capacity

#### PUT

Stores a key-value pair in the MEMSTASH. Updates existing keys and moves them to most recently used position.
If at capacity, removes the least recently used item before adding the new one.

**Syntax:** `<memstash> DO PUT WIT <key> AN WIT <value>`
**Parameters:**
- `key` (STRIN): The cache key to store
- `value` (STRIN): The value to associate with the key

**Example: Basic storage**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "config" AN WIT "{\"host\":\"localhost\"}"
```

**Example: Update with LRU movement**

```lol
CACHE DO PUT WIT "user" AN WIT "alice"
CACHE DO PUT WIT "user" AN WIT "bob" BTW Moves to front, no eviction
```

**Example: Capacity management**

```lol
I HAS A VARIABLE SMALL_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 2
SMALL_CACHE DO PUT WIT "a" AN WIT "1"
SMALL_CACHE DO PUT WIT "b" AN WIT "2"
SMALL_CACHE DO PUT WIT "c" AN WIT "3" BTW Evicts 'a'
IZ SMALL_CACHE DO CONTAINS WIT "a"?
SAYZ WIT "'a' still exists"
NOPE
SAYZ WIT "'a' was evicted"
KTHX
```

**Note:** Updates existing keys without eviction

**Note:** Moves updated items to most recently used position

**Note:** May trigger eviction if at capacity

**Note:** Thread-safe for concurrent operations

**Member Variables:**

#### SIZ

Read-only property that returns the current number of items in the MEMSTASH.


**Example: Monitor cache usage**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "item1" AN WIT "value1"
CACHE DO PUT WIT "item2" AN WIT "value2"
SAYZ WIT "Current size: " MOAR CACHE SIZ MOAR "/100"
```

**Example: Capacity management**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 50
WHILE YEZ
BTW ... add items ...
IZ CACHE SIZ BIGGR THAN 45?
OUTTA HERE BTW Near capacity
KTHX
KTHX
```

**Example: Empty cache detection**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
IZ CACHE SIZ SAEM AS 0?
SAYZ WIT "Cache is empty"
KTHX
```

**Example: Cache efficiency metrics**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000
I HAS A VARIABLE HITS TEH NUMBR ITZ 0
I HAS A VARIABLE MISSES TEH NUMBR ITZ 0
BTW ... process requests ...
IZ CACHE DO CONTAINS WIT "request_key"?
HITS ITZ HITS MOAR 1
I HAS A VARIABLE DATA TEH STRIN ITZ CACHE DO GET WIT "request_key"
NOPE
MISSES ITZ MISSES MOAR 1
BTW ... fetch from source ...
KTHX
I HAS A VARIABLE HIT_RATE TEH NUMBR ITZ HITS / (HITS + MISSES)
```

**Note:** Always reflects current item count

**Note:** Maximum value is the cache capacity

**Note:** Thread-safe for concurrent access

**Note:** Use for monitoring and capacity planning

**Example: Basic LRU cache usage**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100
CACHE DO PUT WIT "user1" AN WIT "John"
CACHE DO PUT WIT "user2" AN WIT "Jane"
I HAS A VARIABLE USER TEH STRIN ITZ CACHE DO GET WIT "user1" BTW Moves to front
SAYZ WIT "User: " MOAR USER
```

**Example: LRU eviction behavior**

```lol
I HAS A VARIABLE SMALL_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 2
SMALL_CACHE DO PUT WIT "first" AN WIT "value1"
SMALL_CACHE DO PUT WIT "second" AN WIT "value2"
SMALL_CACHE DO PUT WIT "third" AN WIT "value3" BTW Evicts 'first'
I HAS A VARIABLE MISSING TEH STRIN ITZ SMALL_CACHE DO GET WIT "first"
IZ MISSING SAEM AS NOTHIN?
SAYZ WIT "First item was evicted"
KTHX
```

**Example: Access pattern optimization**

```lol
I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000
BTW Frequently accessed items stay in cache
WHILE YEZ
CACHE DO GET WIT "hot_data" BTW Keeps this item fresh
CACHE DO PUT WIT "new_data" AN WIT "value"
IZ CACHE SIZ BIGGR THAN 900?
OUTTA HERE BTW Cache getting full
KTHX
KTHX
```

## Ttl Cache

### TIMESTASH Class

A TTL (Time To Live) cache where items expire after a specified duration.
Items are automatically removed when accessed after their expiration time.
Thread-safe implementation with automatic cleanup of expired entries.

**Methods:**

#### CLEAR

Removes all items from the TIMESTASH regardless of expiration.
Resets the cache to empty state while preserving TTL setting.

**Syntax:** `<timestash> DO CLEAR`
**Example: Clear entire cache**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
BTW ... add many items ...
SAYZ WIT "Size before clear: " MOAR CACHE SIZ
CACHE DO CLEAR
SAYZ WIT "Size after clear: " MOAR CACHE SIZ
```

**Example: Reset between operations**

```lol
I HAS A VARIABLE TEMP_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60
BTW ... use for temporary batch processing ...
TEMP_CACHE DO CLEAR BTW Clean slate for next batch
```

**Example: Emergency cleanup**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600
BTW ... cache corrupted or needs refresh ...
CACHE DO CLEAR BTW Emergency reset
```

**Example: Memory management**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
BTW ... long-running application ...
IZ CACHE SIZ BIGGR THAN 1000?
CACHE DO CLEAR BTW Prevent memory bloat
SAYZ WIT "Cache cleared to manage memory"
KTHX
```

**Example: Application restart simulation**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800
BTW ... simulate application restart ...
CACHE DO CLEAR BTW Clear all sessions
SAYZ WIT "All sessions cleared - fresh start"
```

**Note:** Completely empties the cache

**Note:** Preserves the original TTL setting

**Note:** Cannot be undone - use with caution

**Note:** Thread-safe for concurrent operations

**Note:** Size becomes 0 after clearing

#### CONTAINS

Checks if a non-expired key exists in the TIMESTASH.
Returns YEZ if exists and not expired, NO otherwise. Automatically removes expired entries.

**Syntax:** `<timestash> DO CONTAINS WIT <key>`
**Parameters:**
- `key` (STRIN): The cache key to check

**Example: Check existence without retrieval**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
CACHE DO PUT WIT "data" AN WIT "value"
IZ CACHE DO CONTAINS WIT "data"?
SAYZ WIT "Data exists and is fresh"
KTHX
```

**Example: Safe access pattern**

```lol
IZ CACHE DO CONTAINS WIT "user_session"?
I HAS A VARIABLE SESSION TEH STRIN ITZ CACHE DO GET WIT "user_session"
BTW ... use session ...
NOPE
SAYZ WIT "Session not available"
KTHX
```

**Example: Expiration detection**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 1
CACHE DO PUT WIT "temp" AN WIT "value"
BTW ... wait 2 seconds ...
IZ CACHE DO CONTAINS WIT "temp"?
SAYZ WIT "Still fresh"
NOPE
SAYZ WIT "Expired and cleaned up"
KTHX
```

**Example: Cache statistics**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60
I HAS A VARIABLE FRESH_COUNT TEH NUMBR ITZ 0
BTW ... check multiple keys ...
IZ CACHE DO CONTAINS WIT "key1"?
FRESH_COUNT ITZ FRESH_COUNT MOAR 1
KTHX
SAYZ WIT "Fresh items: " MOAR FRESH_COUNT
```

**Note:** Automatically cleans up expired entries during check

**Note:** Does not extend TTL (unlike GET in some cache systems)

**Note:** Thread-safe for concurrent operations

**Note:** Use GET if you need the actual value

#### DELETE

Removes a key-value pair from the TIMESTASH regardless of expiration.
Returns YEZ if deleted, NO if key not found.

**Syntax:** `<timestash> DO DELETE WIT <key>`
**Parameters:**
- `key` (STRIN): The cache key to remove

**Example: Remove specific item**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
CACHE DO PUT WIT "temp" AN WIT "temporary_data"
I HAS A VARIABLE REMOVED TEH BOOL ITZ CACHE DO DELETE WIT "temp"
IZ REMOVED?
SAYZ WIT "Item removed successfully"
KTHX
```

**Example: Force logout**

```lol
I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800
SESSIONS DO DELETE WIT "compromised_session_123"
```

**Example: Cleanup expired data**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60
BTW ... some items may have expired ...
CACHE DO DELETE WIT "expired_item" BTW Remove even if expired
```

**Example: Handle deletion result**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
I HAS A VARIABLE WAS_DELETED TEH BOOL ITZ CACHE DO DELETE WIT "nonexistent"
IZ WAS_DELETED?
SAYZ WIT "Item was deleted"
NOPE
SAYZ WIT "Item didn't exist"
KTHX
```

**Example: Selective cache management**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600
BTW ... cache has various data types ...
CACHE DO DELETE WIT "debug_info" BTW Remove debug data
CACHE DO DELETE WIT "temp_calc" BTW Remove temporary results
```

**Note:** Works regardless of expiration status

**Note:** Returns NO if key doesn't exist (not an error)

**Note:** Thread-safe for concurrent operations

**Note:** Useful for manual cache management

#### GET

Retrieves a non-expired value by key.
Returns NOTHIN if key not found or expired. Automatically removes expired entries.

**Syntax:** `<timestash> DO GET WIT <key>`
**Parameters:**
- `key` (STRIN): The cache key to retrieve

**Example: Retrieve fresh data**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
CACHE DO PUT WIT "data" AN WIT "fresh_value"
I HAS A VARIABLE VALUE TEH STRIN ITZ CACHE DO GET WIT "data"
IZ VALUE SAEM AS NOTHIN?
SAYZ WIT "Data not found or expired"
NOPE
SAYZ WIT "Got: " MOAR VALUE
KTHX
```

**Example: Handle expiration**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 1
CACHE DO PUT WIT "temp" AN WIT "temporary"
BTW ... wait 2 seconds ...
I HAS A VARIABLE EXPIRED TEH STRIN ITZ CACHE DO GET WIT "temp"
IZ EXPIRED SAEM AS NOTHIN?
SAYZ WIT "Data expired and was cleaned up"
KTHX
```

**Example: Session validation**

```lol
I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800
I HAS A VARIABLE SESSION_DATA TEH STRIN ITZ SESSIONS DO GET WIT "user_session"
IZ SESSION_DATA SAEM AS NOTHIN?
SAYZ WIT "Session expired, please login again"
KTHX
```

**Example: Cache hit/miss tracking**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60
I HAS A VARIABLE RESULT TEH STRIN ITZ CACHE DO GET WIT "api_data"
IZ RESULT SAEM AS NOTHIN?
SAYZ WIT "Cache miss - fetching from API"
BTW ... fetch from API ...
CACHE DO PUT WIT "api_data" AN WIT "fetched_data"
NOPE
SAYZ WIT "Cache hit - using cached data"
KTHX
```

**Note:** Automatically cleans up expired entries on access

**Note:** Returns NOTHIN for both missing keys and expired entries

**Note:** Thread-safe for concurrent operations

**Note:** Does not extend TTL on access (unlike some cache implementations)

#### PUT

Stores a key-value pair in the TIMESTASH with TTL expiration.
Overwrites existing keys with new expiration time.

**Syntax:** `<timestash> DO PUT WIT <key> AN WIT <value>`
**Parameters:**
- `key` (STRIN): The cache key to store
- `value` (STRIN): The value to associate with the key

**Example: Store with TTL**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
CACHE DO PUT WIT "user_data" AN WIT "{\"name\":\"John\",\"id\":123}"
```

**Example: Update existing key**

```lol
CACHE DO PUT WIT "config" AN WIT "version:1"
CACHE DO PUT WIT "config" AN WIT "version:2" BTW Resets expiration
```

**Example: Session storage**

```lol
I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800
SESSIONS DO PUT WIT "session_abc123" AN WIT "user_id:456"
```

**Example: API response caching**

```lol
I HAS A VARIABLE API_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600
API_CACHE DO PUT WIT "/api/users" AN WIT "[{\"id\":1,\"name\":\"Alice\"}]"
```

**Example: Temporary data storage**

```lol
I HAS A VARIABLE TEMP TEH TIMESTASH ITZ NEW TIMESTASH WIT 60
TEMP DO PUT WIT "calculation" AN WIT "result:42"
```

**Note:** Each PUT resets the expiration timer for that key

**Note:** Overwrites existing keys with new TTL

**Note:** Thread-safe for concurrent operations

**Note:** TTL applies from time of PUT, not from last access

#### TIMESTASH

Initializes a TIMESTASH TTL cache with the specified expiration time.
TTL must be positive integer representing seconds until expiration.

**Syntax:** `NEW TIMESTASH WIT <ttl_seconds>`
**Parameters:**
- `ttl_seconds` (INTEGR): Time-to-live in seconds for all cache items

**Example: Create 5-minute cache**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
```

**Example: Create session cache (30 minutes)**

```lol
I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800
```

**Example: Create short-lived cache (1 minute)**

```lol
I HAS A VARIABLE TEMP_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60
```

**Example: Handle invalid TTL**

```lol
BTW This would throw an exception:
BTW I HAS A VARIABLE BAD_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 0
```

**Example: API response caching**

```lol
I HAS A VARIABLE API_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600
BTW Cache API responses for 10 minutes
```

**Note:** TTL applies to all items in this cache instance

**Note:** TTL cannot be changed after creation

**Note:** TTL is in seconds (not milliseconds)

**Note:** Consider data freshness requirements when setting TTL

**Member Variables:**

#### SIZ

Read-only property that returns the current number of non-expired items in the TIMESTASH.
Automatically cleans up expired entries before returning the count.


**Example: Monitor active items**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
CACHE DO PUT WIT "item1" AN WIT "value1"
CACHE DO PUT WIT "item2" AN WIT "value2"
SAYZ WIT "Active items: " MOAR CACHE SIZ
```

**Example: Automatic cleanup trigger**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60
BTW ... items expire over time ...
I HAS A VARIABLE CURRENT_SIZE TEH NUMBR ITZ CACHE SIZ
SAYZ WIT "Fresh items after cleanup: " MOAR CURRENT_SIZE
```

**Example: Empty cache detection**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
IZ CACHE SIZ SAEM AS 0?
SAYZ WIT "Cache is empty"
KTHX
```

**Example: Cache efficiency monitoring**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600
I HAS A VARIABLE INITIAL_SIZE TEH NUMBR ITZ CACHE SIZ
BTW ... application runs and items expire ...
I HAS A VARIABLE FINAL_SIZE TEH NUMBR ITZ CACHE SIZ
I HAS A VARIABLE EXPIRED_COUNT TEH NUMBR ITZ INITIAL_SIZE - FINAL_SIZE
SAYZ WIT "Items expired: " MOAR EXPIRED_COUNT
```

**Example: Capacity planning**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800
WHILE YEZ
BTW ... monitor cache usage ...
I HAS A VARIABLE USAGE TEH NUMBR ITZ CACHE SIZ
IZ USAGE BIGGR THAN 1000?
SAYZ WIT "Cache usage high: " MOAR USAGE
KTHX
I HAS A VARIABLE DELAY TEH NUMBR ITZ 60
SLEEPZ WIT DELAY
KTHX
```

**Note:** Triggers automatic cleanup of expired entries

**Note:** Only counts non-expired items

**Note:** Thread-safe for concurrent access

**Note:** May return different values between calls due to expiration

**Note:** Use for monitoring cache health and planning

**Example: Basic TTL cache usage**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300
CACHE DO PUT WIT "session" AN WIT "user123"
I HAS A VARIABLE USER TEH STRIN ITZ CACHE DO GET WIT "session"
SAYZ WIT "User: " MOAR USER
BTW Wait 5+ minutes and session will expire
```

**Example: Session management**

```lol
I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800
SESSIONS DO PUT WIT "user_abc" AN WIT "{\"login\":\"2024-01-01\",\"role\":\"admin\"}"
BTW Session automatically expires in 30 minutes
```

**Example: Short-lived data caching**

```lol
I HAS A VARIABLE API_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60
API_CACHE DO PUT WIT "weather_nyc" AN WIT "{\"temp\":72,\"condition\":\"sunny\"}"
BTW Data expires in 1 minute
```

**Example: Automatic cleanup**

```lol
I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 10
CACHE DO PUT WIT "temp1" AN WIT "value1"
CACHE DO PUT WIT "temp2" AN WIT "value2"
BTW ... wait for expiration ...
I HAS A VARIABLE SIZE TEH NUMBR ITZ CACHE SIZ BTW Triggers cleanup
SAYZ WIT "Active items: " MOAR SIZE
```

