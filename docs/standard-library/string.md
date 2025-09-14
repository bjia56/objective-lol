# STRING Module

## Import

```lol
BTW Full import
I CAN HAS STRING?

BTW Selective import examples
I CAN HAS REPLACE FROM STRING?
I CAN HAS INDEX_OF FROM STRING?
```

### CAPITALIZE

Capitalizes the first character of a STRIN and makes the rest lowercase.
Returns a new STRIN with the first letter capitalized and the rest in lower case.

**Returns:** STRIN

### CONCAT

Concatenates multiple values into a single STRIN.

**Returns:** STRIN

### CONTAINS

Checks if STR contains the substring SUBSTR.
Returns TRUE if SUBSTR is found within STR, otherwise FALSE.

**Returns:** BOOL

### INDEX_OF

Finds the index of the first occurrence of SUBSTR in STR.
Returns the zero-based index of SUBSTR in STR, or -1 if not found.

**Returns:** INTEGR

### LEN

Returns the length of a STRIN in characters.
Counts the number of UTF-8 characters in the STRIN.

**Returns:** INTEGR

### LOWER

Converts all characters in a STRIN to lowercase.
Returns a new STRIN with all letters converted to lower case.

**Returns:** STRIN

### LTRIM

Removes whitespace from the left end of a STRIN.
Trims spaces, tabs, newlines, and carriage returns.

**Returns:** STRIN

### REPEAT

Repeats a STRIN a specified number of times.
Returns a new STRIN consisting of the original STRIN repeated COUNT times.

**Returns:** STRIN

### REPLACE

Replaces the first occurrence of OLD substring with NEW substring in STR.
Returns a new STRIN with the first occurrence of OLD replaced by NEW.

**Returns:** STRIN

### REPLACE_ALL

Replaces all occurrences of OLD substring with NEW substring in STR.
Returns a new STRIN with all occurrences of OLD replaced by NEW.

**Returns:** STRIN

### RTRIM

Removes whitespace from the right end of a STRIN.
Trims spaces, tabs, newlines, and carriage returns.

**Returns:** STRIN

### SPLIT

Splits a STRIN into a BUKKIT array using the specified separator.
Returns array of substrings divided by the separator string.

**Returns:** BUKKIT

### SUBSTR

Extracts a substring from a STRIN starting at the given position.
Returns substring from START index for LENGTH characters. Bounds are checked.

**Returns:** STRIN

### TITLE

Converts the first character of each word to uppercase.
Returns a new STRIN with the first letter of each word capitalized.

**Returns:** STRIN

### TRIM

Removes whitespace from both ends of a STRIN.
Trims spaces, tabs, newlines, and carriage returns.

**Returns:** STRIN

### UPPER

Converts all characters in a STRIN to uppercase.
Returns a new STRIN with all letters converted to upper case.

**Returns:** STRIN

