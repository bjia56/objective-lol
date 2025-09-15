# RANDOM Module

## Import

```lol
BTW Full import
I CAN HAS RANDOM?

BTW Selective import examples
I CAN HAS SEED_TIME FROM RANDOM?
I CAN HAS RANDOM_INT FROM RANDOM?
```

## Seeding

### SEED

Sets the random number generator seed for reproducible results.
Using the same seed will produce the same sequence of random numbers.

**Syntax:** `SEED WIT <seed>`
**Returns:** 

**Parameters:**
- `seed` (INTEGR): The seed value for the random number generator

**Example: Set seed for reproducible results**

```lol
SEED WIT 42
I HAS A VARIABLE NUM1 TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100
I HAS A VARIABLE NUM2 TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100
BTW NUM1 and NUM2 will be the same every time with same seed
```

**Example: Different seeds produce different sequences**

```lol
SEED WIT 123
I HAS A VARIABLE SEQ1 TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 10
SEED WIT 456
I HAS A VARIABLE SEQ2 TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 10
BTW SEQ1 and SEQ2 will likely be different
```

**Example: Reset to same sequence**

```lol
SEED WIT 999
I HAS A VARIABLE W TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100
I HAS A VARIABLE X TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100
SEED WIT 999
I HAS A VARIABLE Y TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100
I HAS A VARIABLE Z TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100
BTW W will equal Y, X will equal Z
```

**Note:** Use SEED_TIME for non-deterministic behavior

**Note:** Same seed always produces same sequence

**Note:** Useful for testing and debugging

**Note:** Affects all subsequent random operations

**See also:** SEED_TIME, RANDOM_INT

### SEED_TIME

Seeds the random number generator with the current time.
Provides different random sequences on each program run.

**Syntax:** `SEED_TIME`
**Returns:** 

**Example: Initialize for random behavior**

```lol
SEED_TIME
I HAS A VARIABLE DICE TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 7
SAYZ WIT "You rolled: " MOAR DICE
```

**Example: Use at program start**

```lol
BTW Always call this first for truly random behavior
SEED_TIME
WHILE YEZ
I HAS A VARIABLE MOVE TEH STRIN ITZ RANDOM_CHOICE WIT MOVES
SAYZ WIT "Computer plays: " MOAR MOVE
I HAS A VARIABLE CONTINUE TEH BOOL ITZ RANDOM_BOOL
IZ CONTINUE?
OUTTA HERE
KTHX
KTHX
```

**Example: Multiple runs produce different results**

```lol
SEED_TIME
I HAS A VARIABLE SESSION_ID TEH STRIN ITZ RANDOM_STRING WIT 8 AN WIT "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
SAYZ WIT "Session ID: " MOAR SESSION_ID
BTW Each program run will have different session ID
```

**Example: Games and simulations**

```lol
SEED_TIME
I HAS A VARIABLE WEATHER TEH STRIN
I HAS A VARIABLE ROLL TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100
IZ ROLL BIGGR THAN 70?
WEATHER ITZ "Sunny"
NOPE IZ ROLL BIGGR THAN 40?
WEATHER ITZ "Cloudy"
NOPE
WEATHER ITZ "Rainy"
KTHX
SAYZ WIT "Today's weather: " MOAR WEATHER
```

**Note:** Should be called once at program start

**Note:** Uses nanosecond precision for uniqueness

**Note:** Subsequent calls will change the sequence

**Note:** Perfect for games, simulations, and one-time use

**See also:** SEED, RANDOM_INT

## Random Numbers

### RANDOM_FLOAT

Returns a random floating-point number between 0.0 (inclusive) and 1.0 (exclusive).
Useful for probability calculations and random selection.

**Syntax:** `RANDOM_FLOAT`
**Returns:** DUBBLE

**Example: Basic random float**

```lol
I HAS A VARIABLE PROBABILITY TEH DUBBLE ITZ RANDOM_FLOAT
SAYZ WIT "Random probability: " MOAR PROBABILITY
```

**Example: Probability-based decisions**

```lol
I HAS A VARIABLE CHANCE TEH DUBBLE ITZ RANDOM_FLOAT
IZ CHANCE BIGGR THAN 0.5?
SAYZ WIT "Heads"
NOPE
SAYZ WIT "Tails"
KTHX
```

**Example: Random selection weights**

```lol
I HAS A VARIABLE WEIGHT TEH DUBBLE ITZ RANDOM_FLOAT
IZ WEIGHT BIGGR THAN 0.7?
I HAS A VARIABLE RARITY TEH STRIN ITZ "Legendary"
NOPE IZ WEIGHT BIGGR THAN 0.4?
RARITY ITZ "Rare"
NOPE
RARITY ITZ "Common"
KTHX
SAYZ WIT "Got a " MOAR RARITY MOAR " item!"
```

**Example: Animation timing**

```lol
WHILE YEZ
I HAS A VARIABLE DELAY TEH DUBBLE ITZ RANDOM_FLOAT MOAR 0.1
BTW Add small random delay for natural feel
I HAS A VARIABLE FRAME TEH DUBBLE ITZ RANDOM_FLOAT
IZ FRAME BIGGR THAN 0.95?
OUTTA HERE
KTHX
KTHX
```

**Note:** Always returns value >= 0.0 and < 1.0

**Note:** Uses high-quality random number generation

**Note:** Thread-safe for concurrent use

**Note:** Perfect for probabilities and normalized random values

**See also:** RANDOM_RANGE, RANDOM_INT

### RANDOM_INT

Returns a random integer within the specified range.
Range is [min, max) - includes min but excludes max. Min must be less than max.

**Syntax:** `RANDOM_INT WIT <min> AN WIT <max>`
**Returns:** INTEGR

**Parameters:**
- `min` (INTEGR): Minimum value (inclusive)
- `max` (INTEGR): Maximum value (exclusive)

**Example: Dice roll**

```lol
I HAS A VARIABLE DICE TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 7
SAYZ WIT "You rolled a " MOAR DICE
```

**Example: Random index**

```lol
I HAS A VARIABLE INDEX TEH INTEGR ITZ RANDOM_INT WIT 0 AN WIT 10
SAYZ WIT "Selected index: " MOAR INDEX
```

**Example: Random ID generation**

```lol
I HAS A VARIABLE USER_ID TEH INTEGR ITZ RANDOM_INT WIT 1000 AN WIT 10000
SAYZ WIT "New user ID: " MOAR USER_ID
```

**Example: Game mechanics**

```lol
I HAS A VARIABLE DAMAGE TEH INTEGR ITZ RANDOM_INT WIT 10 AN WIT 21
I HAS A VARIABLE CRIT_CHANCE TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 101
IZ CRIT_CHANCE BIGGR THAN 90?
DAMAGE ITZ DAMAGE UP 10
SAYZ WIT "Critical hit!"
KTHX
SAYZ WIT "Dealt " MOAR DAMAGE MOAR " damage"
```

**Example: Array shuffling indices**

```lol
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN 5
I HAS A VARIABLE POS TEH INTEGR ITZ RANDOM_INT WIT 0 AN WIT 5
BTW Use for Fisher-Yates shuffle
IDX ITZ IDX MOAR 1
KTHX
```

**Note:** Min must be less than max

**Note:** Returns integers >= min and < max

**Note:** Useful for discrete random selections

**Note:** Thread-safe for concurrent use

**See also:** RANDOM_RANGE, RANDOM_CHOICE

### RANDOM_RANGE

Returns a random floating-point number within the specified range.
Range is [min, max) - includes min but excludes max. Min must be less than max.

**Syntax:** `RANDOM_RANGE WIT <min> AN WIT <max>`
**Returns:** DUBBLE

**Parameters:**
- `min` (DUBBLE): Minimum value (inclusive)
- `max` (DUBBLE): Maximum value (exclusive)

**Example: Random position**

```lol
I HAS A VARIABLE X_POS TEH DUBBLE ITZ RANDOM_RANGE WIT 0.0 AN WIT 100.0
I HAS A VARIABLE Y_POS TEH DUBBLE ITZ RANDOM_RANGE WIT 0.0 AN WIT 50.0
SAYZ WIT "Position: (" MOAR X_POS MOAR ", " MOAR Y_POS MOAR ")"
```

**Example: Random damage**

```lol
I HAS A VARIABLE DAMAGE TEH DUBBLE ITZ RANDOM_RANGE WIT 5.0 AN WIT 15.0
SAYZ WIT "Dealt " MOAR DAMAGE MOAR " damage!"
```

**Example: Temperature simulation**

```lol
I HAS A VARIABLE TEMP TEH DUBBLE ITZ RANDOM_RANGE WIT -10.0 AN WIT 40.0
IZ TEMP BIGGR THAN 30.0?
SAYZ WIT "Hot day!"
NOPE IZ TEMP LSS THAN 0.0?
SAYZ WIT "Freezing!"
NOPE
SAYZ WIT "Nice weather"
KTHX
```

**Example: Random delays**

```lol
I HAS A VARIABLE DELAY TEH DUBBLE ITZ RANDOM_RANGE WIT 0.5 AN WIT 2.0
BTW Wait between 0.5 and 2.0 seconds
SAYZ WIT "Waiting " MOAR DELAY MOAR " seconds..."
```

**Example: Statistical distribution**

```lol
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN 100
I HAS A VARIABLE VALUE TEH DUBBLE ITZ RANDOM_RANGE WIT -1.0 AN WIT 1.0
BTW Collect samples for statistical analysis
IDX ITZ IDX MOAR 1
KTHX
```

**Note:** Min must be less than max

**Note:** Returns values >= min and < max

**Note:** Useful for continuous random values

**Note:** Thread-safe for concurrent use

**See also:** RANDOM_FLOAT, RANDOM_INT

## Random Selection

### RANDOM_BOOL

Returns a random boolean value (YEZ or NO).
Each value has a 50% probability of being returned.

**Syntax:** `RANDOM_BOOL`
**Returns:** BOOL

**Example: Coin flip**

```lol
I HAS A VARIABLE COIN TEH BOOL ITZ RANDOM_BOOL
IZ COIN?
SAYZ WIT "Heads"
NOPE
SAYZ WIT "Tails"
KTHX
```

**Example: Random events**

```lol
I HAS A VARIABLE RAIN TEH BOOL ITZ RANDOM_BOOL
IZ RAIN?
SAYZ WIT "It's raining today"
NOPE
SAYZ WIT "Sunny weather"
KTHX
```

**Example: Game mechanics**

```lol
I HAS A VARIABLE SUCCESS TEH BOOL ITZ RANDOM_BOOL
IZ SUCCESS?
SAYZ WIT "Action succeeded!"
NOPE
SAYZ WIT "Action failed!"
KTHX
```

**Example: Random spawning**

```lol
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN 10
I HAS A VARIABLE SPAWN_ENEMY TEH BOOL ITZ RANDOM_BOOL
IZ SPAWN_ENEMY?
SAYZ WIT "Enemy spawned!"
KTHX
IDX ITZ IDX MOAR 1
KTHX
```

**Example: A/B testing**

```lol
I HAS A VARIABLE USE_NEW_FEATURE TEH BOOL ITZ RANDOM_BOOL
IZ USE_NEW_FEATURE?
SAYZ WIT "Using new feature version"
NOPE
SAYZ WIT "Using old feature version"
KTHX
```

**Note:** Each outcome has exactly 50% probability

**Note:** Useful for binary random decisions

**Note:** Thread-safe for concurrent use

**Note:** Perfect for yes/no scenarios

**See also:** RANDOM_INT, RANDOM_CHOICE

### RANDOM_CHOICE

Returns a randomly selected element from a BUKKIT array.
Array must not be empty. Each element has equal probability of selection.

**Syntax:** `RANDOM_CHOICE WIT <array>`
**Returns:** 

**Parameters:**
- `array` (BUKKIT): The array to select from

**Example: Random card from deck**

```lol
I HAS A VARIABLE CARDS TEH BUKKIT ITZ BUKKIT WIT "Ace" AN "King" AN "Queen" AN "Jack"
I HAS A VARIABLE CARD TEH STRIN ITZ RANDOM_CHOICE WIT CARDS
SAYZ WIT "Drew: " MOAR CARD
```

**Example: Random enemy type**

```lol
I HAS A VARIABLE ENEMIES TEH BUKKIT ITZ BUKKIT WIT "Goblin" AN "Orc" AN "Troll" AN "Dragon"
I HAS A VARIABLE ENEMY TEH STRIN ITZ RANDOM_CHOICE WIT ENEMIES
SAYZ WIT "Encountered a " MOAR ENEMY
```

**Example: Random quote**

```lol
I HAS A VARIABLE QUOTES TEH BUKKIT ITZ BUKKIT WIT "Hello World" AN "LOL" AN "Objective-C"
I HAS A VARIABLE QUOTE TEH STRIN ITZ RANDOM_CHOICE WIT QUOTES
SAYZ QUOTE
```

**Example: Game loot table**

```lol
I HAS A VARIABLE LOOT TEH BUKKIT ITZ BUKKIT WIT "Sword" AN "Shield" AN "Potion" AN "Gold" AN "Key"
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN 3
I HAS A VARIABLE ITEM TEH STRIN ITZ RANDOM_CHOICE WIT LOOT
SAYZ WIT "Found: " MOAR ITEM
IDX ITZ IDX MOAR 1
KTHX
```

**Example: Random color**

```lol
I HAS A VARIABLE COLORS TEH BUKKIT ITZ BUKKIT WIT "Red" AN "Blue" AN "Green" AN "Yellow" AN "Purple"
I HAS A VARIABLE COLOR TEH STRIN ITZ RANDOM_CHOICE WIT COLORS
SAYZ WIT "Selected color: " MOAR COLOR
```

**Note:** Array must not be empty

**Note:** All elements have equal probability

**Note:** Returns element by reference

**Note:** Thread-safe for concurrent use

**See also:** SHUFFLE, RANDOM_INT

## Random Generation

### RANDOM_STRING

Generates a random string of specified length using given character set.
Each character is randomly selected from the charset. Charset must not be empty.

**Syntax:** `RANDOM_STRING WIT <length> AN WIT <charset>`
**Returns:** STRIN

**Parameters:**
- `length` (INTEGR): Length of the generated string
- `charset` (STRIN): Characters to choose from

**Example: Random password**

```lol
I HAS A VARIABLE PASSWORD TEH STRIN ITZ RANDOM_STRING WIT 8 AN WIT "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
SAYZ WIT "Generated password: " MOAR PASSWORD
```

**Example: Random hex color**

```lol
I HAS A VARIABLE COLOR TEH STRIN ITZ RANDOM_STRING WIT 6 AN WIT "0123456789ABCDEF"
SAYZ WIT "Random color: #" MOAR COLOR
```

**Example: Random ID**

```lol
I HAS A VARIABLE ID TEH STRIN ITZ RANDOM_STRING WIT 10 AN WIT "abcdefghijklmnopqrstuvwxyz0123456789"
SAYZ WIT "Random ID: " MOAR ID
```

**Example: Random letters only**

```lol
I HAS A VARIABLE WORD TEH STRIN ITZ RANDOM_STRING WIT 5 AN WIT "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
SAYZ WIT "Random word: " MOAR WORD
```

**Example: Random binary string**

```lol
I HAS A VARIABLE BINARY TEH STRIN ITZ RANDOM_STRING WIT 16 AN WIT "01"
SAYZ WIT "Random binary: " MOAR BINARY
```

**Example: Random emoji**

```lol
I HAS A VARIABLE EMOJI_CHARSET TEH STRIN ITZ "😀😃😄😁😆😅😂🤣😊😇🙂🙃😉😌😍🥰😘😗😙😚😋😛😝😜🤪🤨🧐🤓😎🤩🥳😏😒😞😔😟😕🙁☹️😣😖😫😩🥺😢😭😤😠😡🤬🤯😳🥵🥶😱😨😰😥😓🤗🤔🤭🤫🤥😶😐😑😬🙄😯😦😧😮😲🥱😴🤤😪😵🤐🥴🤢🤮🤧😷🤒🤕🤑🤠😈👿👹👺🤡💩👻💀☠️👽👾🤖🎃😺😸😹😻😼😽🙀😿😾"
I HAS A VARIABLE RANDOM_EMOJI TEH STRIN ITZ RANDOM_STRING WIT 1 AN WIT EMOJI_CHARSET
SAYZ WIT "Random emoji: " MOAR RANDOM_EMOJI
```

**Note:** Length must be >= 0

**Note:** Charset must not be empty

**Note:** Each character has equal probability

**Note:** Thread-safe for concurrent use

**See also:** RANDOM_INT, UUID

### UUID

Generates a random UUID (Universally Unique Identifier) version 4.
Returns a string in format: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx.

**Syntax:** `UUID`
**Returns:** STRIN

**Example: Generate unique ID**

```lol
I HAS A VARIABLE ID TEH STRIN ITZ UUID
SAYZ WIT "Generated UUID: " MOAR ID
```

**Example: User session ID**

```lol
I HAS A VARIABLE SESSION_ID TEH STRIN ITZ UUID
SAYZ WIT "Session ID: " MOAR SESSION_ID
```

**Example: Database primary key**

```lol
I HAS A VARIABLE PRIMARY_KEY TEH STRIN ITZ UUID
SAYZ WIT "New record ID: " MOAR PRIMARY_KEY
```

**Example: File naming**

```lol
I HAS A VARIABLE FILENAME TEH STRIN ITZ UUID MOAR ".txt"
SAYZ WIT "Creating file: " MOAR FILENAME
```

**Example: Object identification**

```lol
I HAS A VARIABLE OBJECT_ID TEH STRIN ITZ UUID
SAYZ WIT "Object created with ID: " MOAR OBJECT_ID
```

**Example: API key generation**

```lol
I HAS A VARIABLE API_KEY TEH STRIN ITZ UUID
SAYZ WIT "Generated API key: " MOAR API_KEY
```

**Example: Transaction ID**

```lol
I HAS A VARIABLE TXN_ID TEH STRIN ITZ UUID
SAYZ WIT "Transaction ID: " MOAR TXN_ID
```

**Note:** Uses cryptographically secure random generation

**Note:** Globally unique with extremely high probability

**Note:** Version 4 format (random)

**Note:** Thread-safe for concurrent use

**See also:** RANDOM_STRING, RANDOM_INT

## Random Utilities

### SHUFFLE

Returns a new BUKKIT with elements randomly shuffled.
Uses Fisher-Yates algorithm. Original array is not modified.

**Syntax:** `SHUFFLE WIT <array>`
**Returns:** BUKKIT

**Parameters:**
- `array` (BUKKIT): The array to shuffle

**Example: Shuffle a deck of cards**

```lol
I HAS A VARIABLE DECK TEH BUKKIT ITZ BUKKIT WIT "A" AN "2" AN "3" AN "4" AN "5"
I HAS A VARIABLE SHUFFLED TEH BUKKIT ITZ SHUFFLE WIT DECK
SAYZ WIT "Original: " MOAR DECK
SAYZ WIT "Shuffled: " MOAR SHUFFLED
```

**Example: Randomize playlist**

```lol
I HAS A VARIABLE SONGS TEH BUKKIT ITZ BUKKIT WIT "Song1" AN "Song2" AN "Song3" AN "Song4"
I HAS A VARIABLE PLAYLIST TEH BUKKIT ITZ SHUFFLE WIT SONGS
SAYZ WIT "Random playlist: " MOAR PLAYLIST
```

**Example: Random game order**

```lol
I HAS A VARIABLE PLAYERS TEH BUKKIT ITZ BUKKIT WIT "Alice" AN "Bob" AN "Charlie" AN "Dave"
I HAS A VARIABLE ORDER TEH BUKKIT ITZ SHUFFLE WIT PLAYERS
SAYZ WIT "Playing order: " MOAR ORDER
```

**Example: Randomize questions**

```lol
I HAS A VARIABLE QUESTIONS TEH BUKKIT ITZ BUKKIT WIT "Q1" AN "Q2" AN "Q3" AN "Q4" AN "Q5"
I HAS A VARIABLE QUIZ TEH BUKKIT ITZ SHUFFLE WIT QUESTIONS
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN 5
SAYZ WIT "Question: " MOAR QUIZ DO AT IDX
IDX ITZ IDX MOAR 1
KTHX
```

**Example: Randomize positions**

```lol
I HAS A VARIABLE POSITIONS TEH BUKKIT ITZ BUKKIT WIT 1 AN 2 AN 3 AN 4 AN 5
I HAS A VARIABLE RANDOM_POS TEH BUKKIT ITZ SHUFFLE WIT POSITIONS
SAYZ WIT "Random positions: " MOAR RANDOM_POS
```

**Note:** Original array remains unchanged

**Note:** Returns new BUKKIT instance

**Note:** Uses Fisher-Yates shuffle algorithm

**Note:** Thread-safe for concurrent use

**See also:** RANDOM_CHOICE, RANDOM_INT

