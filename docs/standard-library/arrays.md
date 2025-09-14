# ARRAYS Module

## Import

```lol
BTW Full import
I CAN HAS ARRAYS?

BTW Selective import examples
```

## Array Creation

### BUKKIT Class

A dynamic array that can hold any combination of values and types.
Provides methods for adding, removing, accessing, and manipulating elements.

**Methods:**

#### AT

Gets the element at the specified index (0-based).
Throws an exception if the index is out of bounds.

**Syntax:** `array DO AT WIT <index>`
**Parameters:**
- `index` (INTEGR): Zero-based index of element to retrieve

**Example: Access elements**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT "a" AN WIT "b" AN WIT "c"
I HAS A VARIABLE FIRST TEH STRIN ITZ ARR DO AT WIT 0
I HAS A VARIABLE SECOND TEH STRIN ITZ ARR DO AT WIT 1
BTW FIRST = "a", SECOND = "b"
```

**Note:** Uses 0-based indexing

#### BUKKIT

Initializes a BUKKIT array with an optional list of initial elements.
Creates a new dynamic array that can grow and shrink as needed.

**Syntax:** `NEW BUKKIT [WIT element1 AN WIT element2 ...]`
**Parameters:**
- `elements` (...ANY): Optional initial elements of any type

**Example: Create empty array**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT
BTW ARR is now empty []
```

**Example: Create with initial values**

```lol
I HAS A VARIABLE NUMS TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
BTW NUMS is now [1, 2, 3]
```

**Note:** Accepts variable number of arguments

#### CLEAR

Removes all elements from the BUKKIT, making it empty.
Resets the array to have zero elements.

**Syntax:** `array DO CLEAR`
**Example: Clear array**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
SAYZ WIT ARR SIZ
BTW Output: 3
ARR DO CLEAR
SAYZ WIT ARR SIZ
BTW Output: 0
```

**Note:** Removes all elements but keeps the array object

**Note:** More efficient than creating a new array

#### CONTAINS

Checks if the BUKKIT contains the specified value.
Returns YEZ if found, NO otherwise.

**Syntax:** `array DO CONTAINS WIT <value>`
**Parameters:**
- `value` (ANY): Value to search for

**Example: Check for number**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
I HAS A VARIABLE HAS_TWO TEH BOOL ITZ ARR DO CONTAINS WIT 2
BTW HAS_TWO = YEZ
```

**Example: Check for string**

```lol
I HAS A VARIABLE PETS TEH BUKKIT ITZ NEW BUKKIT WIT "cat" AN WIT "dog"
I HAS A VARIABLE HAS_BIRD TEH BOOL ITZ PETS DO CONTAINS WIT "bird"
BTW HAS_BIRD = NO
```

**Example: Use in conditional**

```lol
IZ ARR DO CONTAINS WIT 5?
SAYZ WIT "Found 5!"
NOPE
SAYZ WIT "5 not found"
KTHX
```

**Note:** More convenient than FIND when you only need to know if value exists

#### FIND

Finds the first index of the specified value in the BUKKIT.
Returns -1 if the value is not found.

**Syntax:** `array DO FIND WIT <value>`
**Parameters:**
- `value` (ANY): Value to search for

**Example: Find number**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 10 AN WIT 20 AN WIT 30
I HAS A VARIABLE INDEX TEH INTEGR ITZ ARR DO FIND WIT 20
BTW INDEX = 1
```

**Example: Find string**

```lol
I HAS A VARIABLE WORDS TEH BUKKIT ITZ NEW BUKKIT WIT "cat" AN WIT "dog" AN WIT "cat"
I HAS A VARIABLE FIRST_CAT TEH INTEGR ITZ WORDS DO FIND WIT "cat"
BTW FIRST_CAT = 0 (finds first occurrence)
```

**Example: Not found**

```lol
I HAS A VARIABLE NOT_FOUND TEH INTEGR ITZ ARR DO FIND WIT 999
BTW NOT_FOUND = -1
```

**Note:** Uses equality comparison between values

**Note:** Returns index of first match only

#### JOIN

Joins all elements in the BUKKIT into a single string using the specified separator.
Returns an empty string if the BUKKIT is empty.

**Syntax:** `array DO JOIN WIT <separator>`
**Parameters:**
- `separator` (STRIN): String to place between elements

**Example: Join with comma**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
I HAS A VARIABLE RESULT TEH STRIN ITZ ARR DO JOIN WIT ", "
BTW RESULT = "1, 2, 3"
```

**Example: Join words**

```lol
I HAS A VARIABLE WORDS TEH BUKKIT ITZ NEW BUKKIT WIT "hello" AN WIT "world"
I HAS A VARIABLE SENTENCE TEH STRIN ITZ WORDS DO JOIN WIT " "
BTW SENTENCE = "hello world"
```

**Example: Empty array**

```lol
I HAS A VARIABLE EMPTY TEH BUKKIT ITZ NEW BUKKIT
I HAS A VARIABLE EMPTY_STR TEH STRIN ITZ EMPTY DO JOIN WIT ","
BTW EMPTY_STR = ""
```

**Note:** Converts all elements to strings

#### POP

Removes and returns the last element from the BUKKIT.
Throws an exception if the BUKKIT is empty.

**Syntax:** `array DO POP`
**Example: Remove last element**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
I HAS A VARIABLE LAST TEH INTEGR ITZ ARR DO POP
BTW LAST = 3, ARR is now [1, 2]
```

**Example: Stack behavior**

```lol
ARR DO PUSH WIT "top"
I HAS A VARIABLE POPPED TEH STRIN ITZ ARR DO POP
BTW POPPED = "top", ARR is back to [1, 2]
```

**Note:** Modifies the original array

#### PUSH

Adds elements to the end of the BUKKIT.
Returns the BUKKIT's new size.

**Syntax:** `array DO PUSH WIT <element1> [AN WIT <element2> ...]`
**Parameters:**
- `elements` (...ANY): One or more elements to add

**Example: Add single element**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2
I HAS A VARIABLE NEW_SIZE TEH INTEGR ITZ ARR DO PUSH WIT 3
BTW ARR is now [1, 2, 3], NEW_SIZE = 3
```

**Example: Add multiple elements**

```lol
ARR DO PUSH WIT 4 AN WIT 5 AN WIT 6
BTW ARR is now [1, 2, 3, 4, 5, 6]
```

**Note:** Accepts variable number of arguments

#### REVERSE

Reverses the order of elements in the BUKKIT in place.
Returns the BUKKIT itself.

**Syntax:** `array DO REVERSE`
**Example: Reverse array**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
ARR DO REVERSE
BTW ARR is now [3, 2, 1]
```

**Example: Method chaining**

```lol
I HAS A VARIABLE RESULT TEH BUKKIT ITZ ARR DO REVERSE DO SORT
BTW RESULT is [1, 2, 3] (reversed then sorted)
```

**Note:** Modifies the original array

**Note:** Returns self for method chaining

#### SET

Sets the element at the specified index to the given value.
Throws an exception if the index is out of bounds.

**Syntax:** `array DO SET WIT <index> AN WIT <value>`
**Parameters:**
- `index` (INTEGR): Zero-based index of element to modify
- `value` (ANY): New value to assign

**Example: Modify elements**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
ARR DO SET WIT 1 AN WIT 99
BTW ARR is now [1, 99, 3]
```

**Example: Change types**

```lol
ARR DO SET WIT 0 AN WIT "hello"
BTW ARR is now ["hello", 99, 3]
```

**Note:** Can change element type

#### SHIFT

Removes and returns the first element from the BUKKIT.
Throws an exception if the BUKKIT is empty.

**Syntax:** `array DO SHIFT`
**Example: Remove first element**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
I HAS A VARIABLE FIRST TEH INTEGR ITZ ARR DO SHIFT
BTW FIRST = 1, ARR is now [2, 3]
```

**Example: Queue behavior**

```lol
ARR DO PUSH WIT 4
I HAS A VARIABLE NEXT TEH INTEGR ITZ ARR DO SHIFT
BTW NEXT = 2, ARR is now [3, 4]
```

**Note:** Modifies the original array

#### SLICE

Creates a new BUKKIT containing elements from START index to END index (exclusive).
Supports negative indices to count from the end.
Throws an exception if indices are out of bounds.

**Syntax:** `array DO SLICE WIT <start> AN WIT <end>`
**Parameters:**
- `start` (INTEGR): Starting index (inclusive)
- `end` (INTEGR): Ending index (exclusive)

**Example: Basic slicing**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 0 AN WIT 1 AN WIT 2 AN WIT 3 AN WIT 4
I HAS A VARIABLE SUB TEH BUKKIT ITZ ARR DO SLICE WIT 1 AN WIT 3
BTW SUB = [1, 2] (indices 1 and 2, excluding 3)
```

**Example: Negative indices**

```lol
I HAS A VARIABLE LAST_TWO TEH BUKKIT ITZ ARR DO SLICE WIT -2 AN WIT -0
BTW LAST_TWO = [3, 4] (last two elements)
```

**Example: Copy array**

```lol
I HAS A VARIABLE COPY TEH BUKKIT ITZ ARR DO SLICE WIT 0 AN WIT ARR SIZ
BTW COPY is a shallow copy of ARR
```

**Note:** Creates a new array, doesn't modify original

**Note:** Negative indices count from end (-1 = last element)

#### SORT

Sorts the elements in the BUKKIT in place in ascending order.
Handles different type combinations and converts to strings as fallback.
Returns the BUKKIT itself.

**Syntax:** `array DO SORT`
**Example: Sort numbers**

```lol
I HAS A VARIABLE NUMS TEH BUKKIT ITZ NEW BUKKIT WIT 3 AN WIT 1 AN WIT 2
NUMS DO SORT
BTW NUMS is now [1, 2, 3]
```

**Example: Sort strings**

```lol
I HAS A VARIABLE WORDS TEH BUKKIT ITZ NEW BUKKIT WIT "banana" AN WIT "apple" AN WIT "cherry"
WORDS DO SORT
BTW WORDS is now ["apple", "banana", "cherry"]
```

**Example: Mixed types**

```lol
I HAS A VARIABLE MIXED TEH BUKKIT ITZ NEW BUKKIT WIT 2 AN WIT "1" AN WIT 3
MIXED DO SORT
BTW Sorts by string representation when types differ
```

**Note:** Modifies the original array

**Note:** Numbers sort numerically, strings alphabetically

**Note:** Mixed types fall back to string comparison

#### UNSHIFT

Adds an element to the beginning of the BUKKIT.
Returns the new size of the BUKKIT.

**Syntax:** `array DO UNSHIFT WIT <element>`
**Parameters:**
- `element` (ANY): Element to add at the beginning

**Example: Add to beginning**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 2 AN WIT 3
I HAS A VARIABLE NEW_SIZE TEH INTEGR ITZ ARR DO UNSHIFT WIT 1
BTW NEW_SIZE = 3, ARR is now [1, 2, 3]
```

**Example: Queue behavior**

```lol
ARR DO UNSHIFT WIT "first"
BTW ARR is now ["first", 1, 2, 3]
```

**Note:** Shifts all existing elements to higher indices

**Member Variables:**

#### SIZ

Read-only property that returns the current number of elements in the BUKKIT.
Automatically updated when elements are added or removed.


**Example: Check array size**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
SAYZ WIT ARR SIZ
BTW Output: 3
```

**Example: Empty array**

```lol
I HAS A VARIABLE EMPTY TEH BUKKIT ITZ NEW BUKKIT
SAYZ WIT EMPTY SIZ
BTW Output: 0
```

**Example: Dynamic sizing**

```lol
ARR DO PUSH WIT 4
SAYZ WIT ARR SIZ
BTW Output: 4
ARR DO POP
SAYZ WIT ARR SIZ
BTW Output: 3
```

**Note:** Always reflects current element count

**Note:** Cannot be modified directly

**Example: Create empty array**

```lol
I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT
BTW Creates an empty BUKKIT
```

**Example: Create array with initial values**

```lol
I HAS A VARIABLE NUMS TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3
BTW Creates BUKKIT with [1, 2, 3]
```

**Example: Mixed type array**

```lol
I HAS A VARIABLE MIXED TEH BUKKIT ITZ NEW BUKKIT WIT "hello" AN WIT 42 AN WIT YEZ
BTW Creates BUKKIT with ["hello", 42, YEZ]
```

