# SYSTEM Module

## Import

```lol
BTW Full import
I CAN HAS SYSTEM?

BTW Selective import examples
I CAN HAS ENV FROM SYSTEM?
I CAN HAS OS FROM SYSTEM?
I CAN HAS ARCH FROM SYSTEM?
```

## Global Variables

### ENV

Global environment variable manager. Pre-initialized ENVBASKIT instance containing all current environment variables.

**Type:** ENVBASKIT

```lol
I CAN HAS ENV FROM MODULE?

I HAS A VARIABLE HOME_DIR TEH STRIN ITZ ENV DO GET WIT "HOME"
I HAS A VARIABLE USER_NAME TEH STRIN ITZ ENV DO GET WIT "USER"
SAYZ WIT "Hello " MOAR USER_NAME MOAR ", your home is " MOAR HOME_DIR
```

```lol
I CAN HAS ENV FROM MODULE?

I HAS A VARIABLE PATH_VAR TEH STRIN ITZ ENV DO GET WIT "PATH"
I HAS A VARIABLE SHELL_VAR TEH STRIN ITZ ENV DO GET WIT "SHELL"
SAYZ WIT "Using shell: " MOAR SHELL_VAR
```

```lol
I CAN HAS ENV FROM MODULE?

ENV DO PUT WIT "MY_APP_CONFIG" AN WIT "/etc/myapp/config.json"
ENV DO PUT WIT "MY_APP_DEBUG" AN WIT YEZ
```

```lol
I CAN HAS ENV FROM MODULE?

ENV DO PUT WIT "TEMP_FILE" AN WIT "/tmp/temp.txt"
BTW ... use temp file ...
ENV DO REMOVE WIT "TEMP_FILE" BTW Clean up
```

## Environment Management

### ENVBASKIT Class

Special BASKIT type to provide integration with system environment variables.
Automatically syncs with the actual process environment and provides enhanced functionality.

**Methods:**

#### CLEAR

Clears all environment variables that are tracked in the internal map and unsets them from the actual environment.
Only removes variables that were previously accessed or set through this ENVBASKIT instance.

**Syntax:** `<envbaskit> DO CLEAR`
**Example: Clear all tracked variables**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
ENV DO PUT WIT "VAR1" AN WIT "value1"
ENV DO PUT WIT "VAR2" AN WIT "value2"
SAYZ WIT "Before clear: " MOAR ENV SIZ
ENV DO CLEAR
SAYZ WIT "After clear: " MOAR ENV SIZ
```

**Example: Fresh start with environment**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
BTW ... modify some environment variables ...
ENV DO CLEAR BTW Reset to clean state
```

**Note:** Only clears variables that were accessed through this ENVBASKIT instance

**Note:** Does not affect environment variables not accessed through this instance

**Note:** After clearing, the instance starts fresh and will reload from environment on next access

#### CONTAINS

Checks if an environment variable exists in either the internal map or actual environment.

**Syntax:** `<envbaskit> DO CONTAINS WIT <key>`
**Parameters:**
- `key` (STRIN): The environment variable name to check

**Example: Check if variable exists**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
IZ ENV DO CONTAINS WIT "HOME"?
SAYZ WIT "HOME variable exists"
KTHX
```

**Example: Check before accessing**

```lol
IZ ENV DO CONTAINS WIT "MY_CONFIG"?
I HAS A VARIABLE CONFIG TEH STRIN ITZ ENV DO GET WIT "MY_CONFIG"
SAYZ WIT "Config: " MOAR CONFIG
NOPE
SAYZ WIT "MY_CONFIG not set, using default"
KTHX
```

**Example: Check system variables**

```lol
IZ ENV DO CONTAINS WIT "PATH"?
I HAS A VARIABLE PATH_VAR TEH STRIN ITZ ENV DO GET WIT "PATH"
SAYZ WIT "PATH is set to: " MOAR PATH_VAR
KTHX
```

**Note:** Searches both internal map and actual environment variables

**Note:** Environment variables found are automatically cached in the internal map

#### COPY

Creates a new ENVBASKIT instance with the same data as the current instance.
The copy is independent and does not sync changes back to the original.

**Syntax:** `<envbaskit> DO COPY`
**Example: Create independent copy**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
ENV DO PUT WIT "MY_VAR" AN WIT "original_value"
I HAS A VARIABLE ENV_COPY TEH ENVBASKIT ITZ ENV DO COPY
ENV_COPY DO PUT WIT "MY_VAR" AN WIT "modified_value"
SAYZ WIT "Original: " MOAR (ENV DO GET WIT "MY_VAR")
SAYZ WIT "Copy: " MOAR (ENV_COPY DO GET WIT "MY_VAR")
```

**Example: Backup current state**

```lol
I HAS A VARIABLE BACKUP TEH ENVBASKIT ITZ ENV DO COPY
BTW ... modify environment ...
ENV DO CLEAR
ENV DO MERGE WIT BACKUP BTW Restore from backup
```

**Example: Isolated environment testing**

```lol
I HAS A VARIABLE TEST_ENV TEH ENVBASKIT ITZ ENV DO COPY
TEST_ENV DO PUT WIT "TEST_VAR" AN WIT "test_value"
BTW ... run tests with TEST_ENV ...
BTW Changes don't affect original ENV
```

**Note:** The copy is completely independent of the original

**Note:** Changes to the copy don't affect the original and vice versa

**Note:** Both copies will sync with the actual environment independently

#### GET

Gets the value of an environment variable.
Checks the internal map first, then the actual environment, throws exception if not found.

**Syntax:** `<envbaskit> DO GET WIT <key>`
**Parameters:**
- `key` (STRIN): The environment variable name to retrieve

**Example: Get environment variable**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
I HAS A VARIABLE HOME_DIR TEH STRIN ITZ ENV DO GET WIT "HOME"
SAYZ WIT "Home directory: "
SAYZ WIT HOME_DIR
```

**Example: Get custom variable**

```lol
ENV DO PUT WIT "MY_APP_CONFIG" AN WIT "/etc/myapp.conf"
I HAS A VARIABLE CONFIG_PATH TEH STRIN ITZ ENV DO GET WIT "MY_APP_CONFIG"
```

**Example: Handle missing variable**

```lol
IZ ENV DO CONTAINS WIT "NON_EXISTENT_VAR"?
I HAS A VARIABLE VALUE TEH STRIN ITZ ENV DO GET WIT "NON_EXISTENT_VAR"
SAYZ WIT VALUE
KTHX
```

**Note:** Searches internal map first, then actual environment variables

**Note:** Values from environment are cached in the internal map after first access

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

Sets an environment variable both in the internal map and the actual process environment.

**Syntax:** `<envbaskit> DO PUT WIT <key> AN WIT <value>`
**Parameters:**
- `key` (STRIN): The environment variable name
- `value` (): The value to set (converted to string)

**Example: Set environment variable**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
ENV DO PUT WIT "MY_VAR" AN WIT "my_value"
```

**Example: Set numeric value**

```lol
ENV DO PUT WIT "PORT" AN WIT 8080
```

**Example: Override existing variable**

```lol
ENV DO PUT WIT "PATH" AN WIT "/usr/local/bin:/usr/bin"
```

**Note:** Changes are reflected in both the ENVBASKIT instance and the actual process environment

**Note:** Values are automatically converted to strings

#### REFRESH

Refreshes the internal map with all current environment variables.
Discards any previous state and reloads from actual environment.

**Syntax:** `<envbaskit> DO REFRESH`
**Example: Refresh after external changes**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
BTW External process changes environment...
ENV DO REFRESH BTW Reload current environment state
I HAS A VARIABLE NEW_VAR TEH STRIN ITZ ENV DO GET WIT "NEW_VAR"
```

**Example: Reset to clean environment state**

```lol
ENV DO PUT WIT "TEMP1" AN WIT "value1"
ENV DO PUT WIT "TEMP2" AN WIT "value2"
ENV DO REFRESH BTW Discard changes, reload from environment
```

**Example: Periodic refresh**

```lol
WHILE YEZ
ENV DO REFRESH
I HAS A VARIABLE STATUS TEH STRIN ITZ ENV DO GET WIT "STATUS_VAR"
BTW ... process status ...
I HAS A VARIABLE DELAY TEH NUMBR ITZ 5
SLEEPZ WIT DELAY
KTHX
```

**Note:** Discards all previous changes made through this ENVBASKIT instance

**Note:** Reloads all current environment variables into the internal map

**Note:** Useful when external processes have modified the environment

#### REMOVE

Removes an environment variable from both the internal map and actual process environment.
Returns the previous value, throws exception if not found.

**Syntax:** `<envbaskit> DO REMOVE WIT <key>`
**Parameters:**
- `key` (STRIN): The environment variable name to remove

**Example: Remove environment variable**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
ENV DO PUT WIT "TEMP_VAR" AN WIT "temporary_value"
I HAS A VARIABLE OLD_VALUE TEH STRIN ITZ ENV DO REMOVE WIT "TEMP_VAR"
SAYZ WIT "Removed value was: " MOAR OLD_VALUE
```

**Example: Clean up after use**

```lol
ENV DO PUT WIT "SESSION_ID" AN WIT "abc123"
BTW ... use session ...
ENV DO REMOVE WIT "SESSION_ID" BTW Clean up
```

**Example: Handle removal errors**

```lol
IZ ENV DO CONTAINS WIT "NON_EXISTENT"?
ENV DO REMOVE WIT "NON_EXISTENT"
NOPE
SAYZ WIT "Variable doesn't exist, can't remove"
KTHX
```

**Note:** Removes from both internal map and actual process environment

**Note:** Returns the value that was removed

**Note:** Throws exception if variable doesn't exist in either location

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

Number of environment variables currently tracked in this ENVBASKIT instance.


**Example: Check number of variables**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
SAYZ WIT "Environment variables: " MOAR ENV SIZ
```

**Example: Monitor variable count**

```lol
ENV DO PUT WIT "VAR1" AN WIT "value1"
ENV DO PUT WIT "VAR2" AN WIT "value2"
SAYZ WIT "After adding: " MOAR ENV SIZ
ENV DO REMOVE WIT "VAR1"
SAYZ WIT "After removing: " MOAR ENV SIZ
```

**Note:** Only counts variables accessed through this ENVBASKIT instance

**Note:** Does not include all environment variables, only those tracked internally

**Example: Basic environment variable access**

```lol
I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT
I HAS A VARIABLE HOME_DIR TEH STRIN ITZ ENV DO GET WIT "HOME"
SAYZ WIT "Home directory: " MOAR HOME_DIR
```

**Example: Set and get custom variables**

```lol
ENV DO PUT WIT "MY_APP_CONFIG" AN WIT "/etc/myapp.conf"
ENV DO PUT WIT "MY_APP_PORT" AN WIT 8080
I HAS A VARIABLE CONFIG TEH STRIN ITZ ENV DO GET WIT "MY_APP_CONFIG"
I HAS A VARIABLE PORT TEH NUMBR ITZ ENV DO GET WIT "MY_APP_PORT"
```

**Example: Check variable existence**

```lol
IZ ENV DO CONTAINS WIT "PATH"?
I HAS A VARIABLE PATH_VAR TEH STRIN ITZ ENV DO GET WIT "PATH"
SAYZ WIT "PATH: " MOAR PATH_VAR
KTHX
```

**Example: Environment variable cleanup**

```lol
ENV DO PUT WIT "TEMP_SESSION" AN WIT "session123"
BTW ... use session ...
ENV DO REMOVE WIT "TEMP_SESSION" BTW Clean up
```

**Example: Refresh from external changes**

```lol
BTW External process modifies environment...
ENV DO REFRESH BTW Reload current environment state
I HAS A VARIABLE NEW_VAR TEH STRIN ITZ ENV DO GET WIT "EXTERNAL_VAR"
```

**Example: Create isolated environment copy**

```lol
I HAS A VARIABLE TEST_ENV TEH ENVBASKIT ITZ ENV DO COPY
TEST_ENV DO PUT WIT "TEST_VAR" AN WIT "test_value"
BTW ... run tests ...
BTW Original ENV unaffected by TEST_ENV changes
```

## System Information

### ARCH

System architecture (e.g. amd64, 386, arm64).

**Type:** STRIN

```lol
I CAN HAS ARCH FROM MODULE?

IZ ARCH SAEM AS "amd64"?
SAYZ WIT "Running on 64-bit x86"
KTHX
```

```lol
I CAN HAS ARCH FROM MODULE?

IZ ARCH SAEM AS "arm64"?
SAYZ WIT "Running on ARM 64-bit"
KTHX
```

```lol
I CAN HAS ARCH FROM MODULE?

SAYZ WIT "Platform: " MOAR OS MOAR "-" MOAR ARCH
```

```lol
I CAN HAS ARCH FROM MODULE?

I HAS A VARIABLE IS_64BIT TEH BOOL
IZ ARCH SAEM AS "amd64" OR ARCH SAEM AS "arm64"?
IS_64BIT ITZ YEZ
NOPE
IS_64BIT ITZ NO
KTHX
```

```lol
I CAN HAS ARCH FROM MODULE?

I HAS A VARIABLE LIB_PATH TEH STRIN
LIB_PATH ITZ "/usr/lib/" MOAR ARCH
SAYZ WIT "Library path: " MOAR LIB_PATH
```

### OS

Operating system name (e.g. windows, linux, darwin).

**Type:** STRIN

```lol
I CAN HAS OS FROM MODULE?

IZ OS SAEM AS "linux"?
SAYZ WIT "Running on Linux"
KTHX
```

```lol
I CAN HAS OS FROM MODULE?

I HAS A VARIABLE PATH_SEP TEH STRIN
IZ OS SAEM AS "windows"?
PATH_SEP ITZ "\\"
NOPE
PATH_SEP ITZ "/"
KTHX
SAYZ WIT "Path separator: " MOAR PATH_SEP
```

```lol
I CAN HAS OS FROM MODULE?

IZ OS SAEM AS "darwin"?
SAYZ WIT "Running on macOS"
KTHX
```

```lol
I CAN HAS OS FROM MODULE?

I HAS A VARIABLE CMD TEH STRIN
IZ OS SAEM AS "windows"?
CMD ITZ "dir"
NOPE
CMD ITZ "ls -la"
KTHX
```

## Miscellaneous

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

