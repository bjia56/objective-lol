# MAPS Module

## Import

```lol
BTW Full import
I CAN HAS MAPS?

BTW Selective import examples
```

## Map Creation

### BASKIT Class

A dynamic map (dictionary) that stores key-value pairs.
Keys are STRIN type and values can be any type.
Provides fast lookup, insertion, and deletion of key-value pairs.

**Methods:**

#### BASKIT

Initializes an empty BASKIT map.
Creates a new dictionary with no key-value pairs.

**Syntax:** `NEW BASKIT`
**Example: Create empty map**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
SAYZ WIT MAP SIZ
BTW Output: 0
```

**Note:** Creates an empty dictionary ready for key-value pairs

#### CLEAR

Removes all key-value pairs from the BASKIT.
After clearing, the map size will be 0.

**Syntax:** `map DO CLEAR`
**Example: Clear all entries**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
MAP DO PUT WIT "a" AN WIT 1
MAP DO PUT WIT "b" AN WIT 2
SAYZ WIT MAP SIZ
BTW Output: 2
MAP DO CLEAR
SAYZ WIT MAP SIZ
BTW Output: 0
```

**Example: Reset for reuse**

```lol
MAP DO CLEAR
MAP DO PUT WIT "fresh" AN WIT "start"
BTW MAP is now empty and ready for new data
```

**Note:** More efficient than creating a new BASKIT

**Note:** Keeps the same object but removes all contents

#### CONTAINS

Checks if the specified key exists in the BASKIT.
Returns YEZ if key exists, NO otherwise.

**Syntax:** `map DO CONTAINS WIT <key>`
**Parameters:**
- `key` (STRIN): The key to check for existence

**Example: Check for keys**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
MAP DO PUT WIT "name" AN WIT "Carol"
I HAS A VARIABLE HAS_NAME TEH BOOL ITZ MAP DO CONTAINS WIT "name"
I HAS A VARIABLE HAS_AGE TEH BOOL ITZ MAP DO CONTAINS WIT "age"
BTW HAS_NAME = YEZ, HAS_AGE = NO
```

**Example: Use in conditional**

```lol
IZ MAP DO CONTAINS WIT "score"?
I HAS A VARIABLE SCORE TEH INTEGR ITZ MAP DO GET WIT "score"
NOPE
SAYZ WIT "Score not set"
KTHX
```

**Note:** Safer than GET when you only need to check existence

#### COPY

Creates a shallow copy of the BASKIT with all current key-value pairs.
Changes to the copy do not affect the original BASKIT.

**Syntax:** `map DO COPY`
**Example: Create independent copy**

```lol
I HAS A VARIABLE ORIGINAL TEH BASKIT ITZ NEW BASKIT
ORIGINAL DO PUT WIT "shared" AN WIT "value"
I HAS A VARIABLE COPY TEH BASKIT ITZ ORIGINAL DO COPY
COPY DO PUT WIT "new" AN WIT "item"
BTW ORIGINAL doesn't have "new" key, COPY does
```

**Example: Backup before modification**

```lol
I HAS A VARIABLE BACKUP TEH BASKIT ITZ ORIGINAL DO COPY
ORIGINAL DO PUT WIT "temp" AN WIT "data"
BTW Can restore from BACKUP if needed
```

**Note:** Creates a shallow copy (references to objects are shared)

**Note:** Independent of original for key-value structure

#### GET

Retrieves the value associated with the specified key.
Throws an exception if the key is not found.

**Syntax:** `map DO GET WIT <key>`
**Parameters:**
- `key` (STRIN): The key to look up

**Example: Retrieve values**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
MAP DO PUT WIT "name" AN WIT "Bob"
MAP DO PUT WIT "score" AN WIT 95
I HAS A VARIABLE NAME TEH STRIN ITZ MAP DO GET WIT "name"
I HAS A VARIABLE SCORE TEH INTEGR ITZ MAP DO GET WIT "score"
BTW NAME = "Bob", SCORE = 95
```

**Example: Handle missing key**

```lol
MAYB
I HAS A VARIABLE MISSING TEH STRIN ITZ MAP DO GET WIT "missing"
OOPSIE ERR
SAYZ WIT "Key not found!"
KTHX
```

**Note:** Use CONTAINS to check if key exists first

#### KEYS

Returns a BUKKIT containing all keys in the BASKIT.
Keys are sorted alphabetically for consistent ordering.

**Syntax:** `map DO KEYS`
**Example: Get all keys**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
MAP DO PUT WIT "zebra" AN WIT 1
MAP DO PUT WIT "apple" AN WIT 2
MAP DO PUT WIT "banana" AN WIT 3
I HAS A VARIABLE KEYS TEH BUKKIT ITZ MAP DO KEYS
BTW KEYS = ["apple", "banana", "zebra"] (alphabetical)
```

**Example: Iterate over keys**

```lol
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN KEYS SIZ
I HAS A VARIABLE KEY TEH STRIN ITZ KEYS DO AT WIT IDX
I HAS A VARIABLE VALUE TEH INTEGR ITZ MAP DO GET WIT KEY
SAYZ WIT KEY MOAR ": " MOAR VALUE
IDX ITZ IDX MOAR 1
KTHX
```

**Note:** Keys are always sorted alphabetically

**Note:** Returns empty BUKKIT if map is empty

#### MERGE

Merges another BASKIT's key-value pairs into this BASKIT.
Existing keys are overwritten with values from the other BASKIT.

**Syntax:** `map DO MERGE WIT <other_baskit>`
**Parameters:**
- `other` (BASKIT): Another BASKIT to merge from

**Example: Merge maps**

```lol
I HAS A VARIABLE MAP1 TEH BASKIT ITZ NEW BASKIT
I HAS A VARIABLE MAP2 TEH BASKIT ITZ NEW BASKIT
MAP1 DO PUT WIT "a" AN WIT 1
MAP1 DO PUT WIT "b" AN WIT 2
MAP2 DO PUT WIT "b" AN WIT 99
MAP2 DO PUT WIT "c" AN WIT 3
MAP1 DO MERGE WIT MAP2
BTW MAP1 now contains: a->1, b->99, c->3 (b was overwritten)
```

**Example: Configuration merging**

```lol
I HAS A VARIABLE DEFAULTS TEH BASKIT ITZ NEW BASKIT
I HAS A VARIABLE USER_CONFIG TEH BASKIT ITZ NEW BASKIT
DEFAULTS DO PUT WIT "timeout" AN WIT 30
DEFAULTS DO PUT WIT "retries" AN WIT 3
USER_CONFIG DO PUT WIT "timeout" AN WIT 60
DEFAULTS DO MERGE WIT USER_CONFIG
BTW DEFAULTS now has user's timeout but default retries
```

**Note:** Modifies the original BASKIT

**Note:** Overwrites existing keys with new values

#### PAIRS

Returns a BUKKIT of key-value pairs as BUKKITs containing [key, value].
Useful for iterating over both keys and values simultaneously.

**Syntax:** `map DO PAIRS`
**Example: Get key-value pairs**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
MAP DO PUT WIT "name" AN WIT "David"
MAP DO PUT WIT "age" AN WIT 30
I HAS A VARIABLE PAIRS TEH BUKKIT ITZ MAP DO PAIRS
BTW PAIRS = [["age", 30], ["name", "David"]] (by key order)
```

**Example: Iterate over pairs**

```lol
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN PAIRS SIZ
I HAS A VARIABLE PAIR TEH BUKKIT ITZ PAIRS DO AT WIT IDX
I HAS A VARIABLE KEY TEH STRIN ITZ PAIR DO AT WIT 0
I HAS A VARIABLE VALUE TEH STRIN ITZ PAIR DO AT WIT 1
SAYZ WIT KEY MOAR ": " MOAR VALUE
IDX ITZ IDX MOAR 1
KTHX
```

**Note:** Each pair is a BUKKIT with [key, value]

**Note:** Pairs are ordered by key alphabetically

#### PUT

Stores a key-value pair in the BASKIT.
If key already exists, overwrites the existing value.

**Syntax:** `map DO PUT WIT <key> AN WIT <value>`
**Parameters:**
- `key` (STRIN): The key to store (converted to string)
- `value` (ANY): The value to associate with the key

**Example: Add new entries**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
MAP DO PUT WIT "name" AN WIT "Alice"
MAP DO PUT WIT "age" AN WIT 25
BTW MAP now contains name->Alice and age->25
```

**Example: Overwrite existing key**

```lol
MAP DO PUT WIT "age" AN WIT 26
BTW age is now 26 (was 25)
```

**Example: Mixed value types**

```lol
MAP DO PUT WIT "active" AN WIT YEZ
MAP DO PUT WIT "items" AN WIT NEW BUKKIT
BTW Values can be any type
```

**Note:** Key is converted to string if not already

**Note:** Overwrites existing values for the same key

#### REMOVE

Removes a key-value pair from the BASKIT and returns the value.
Throws an exception if the key is not found.

**Syntax:** `map DO REMOVE WIT <key>`
**Parameters:**
- `key` (STRIN): The key to remove

**Example: Remove entries**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
MAP DO PUT WIT "temp" AN WIT 42
MAP DO PUT WIT "keep" AN WIT "important"
I HAS A VARIABLE REMOVED TEH INTEGR ITZ MAP DO REMOVE WIT "temp"
BTW REMOVED = 42, MAP now only contains "keep"
```

**Example: Safe removal**

```lol
IZ MAP DO CONTAINS WIT "old_key"?
I HAS A VARIABLE OLD_VAL TEH STRIN ITZ MAP DO REMOVE WIT "old_key"
KTHX
```

**Note:** Use CONTAINS to check before removing if unsure

#### VALUES

Returns a BUKKIT containing all values in the BASKIT.
Values are ordered according to their keys' alphabetical order.

**Syntax:** `map DO VALUES`
**Example: Get all values**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
MAP DO PUT WIT "c" AN WIT "third"
MAP DO PUT WIT "a" AN WIT "first"
MAP DO PUT WIT "b" AN WIT "second"
I HAS A VARIABLE VALUES TEH BUKKIT ITZ MAP DO VALUES
BTW VALUES = ["first", "second", "third"] (by key order: a, b, c)
```

**Example: Process all values**

```lol
I HAS A VARIABLE TOTAL TEH INTEGR ITZ 0
I HAS A VARIABLE SCORES TEH BUKKIT ITZ MAP DO VALUES
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN SCORES SIZ
TOTAL ITZ TOTAL MOAR (SCORES DO AT WIT IDX)
IDX ITZ IDX MOAR 1
KTHX
```

**Note:** Values are ordered by their keys' alphabetical order

**Note:** Returns empty BUKKIT if map is empty

**Member Variables:**

#### SIZ

The number of key-value pairs currently stored in the BASKIT.
Read-only property that updates automatically.


**Example: Check map size**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
SAYZ WIT MAP SIZ
BTW Output: 0
MAP DO PUT WIT "key1" AN WIT "value1"
MAP DO PUT WIT "key2" AN WIT "value2"
SAYZ WIT MAP SIZ
BTW Output: 2
```

**Example: Empty after clear**

```lol
MAP DO CLEAR
SAYZ WIT MAP SIZ
BTW Output: 0
```

**Note:** Always reflects current number of key-value pairs

**Note:** Cannot be modified directly

**Example: Create empty map**

```lol
I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
BTW Creates an empty dictionary
```

**Example: Store and retrieve values**

```lol
MAP DO PUT WIT "name" AN WIT "Alice"
MAP DO PUT WIT "age" AN WIT 25
I HAS A VARIABLE NAME TEH STRIN ITZ MAP DO GET WIT "name"
BTW NAME = "Alice"
```

**Example: Mixed value types**

```lol
MAP DO PUT WIT "count" AN WIT 100
MAP DO PUT WIT "active" AN WIT YEZ
MAP DO PUT WIT "items" AN WIT NEW BUKKIT
BTW Values can be any type: STRIN, INTEGR, BOOL, BUKKIT, etc.
```

