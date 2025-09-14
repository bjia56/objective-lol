package stdlib

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleRandomCategories defines the order that categories should be rendered in documentation
var moduleRandomCategories = []string{
	"seeding",
	"random-numbers",
	"random-selection",
	"random-generation",
	"random-utilities",
}

// Global RANDOM function definitions - created once and reused
var randomFuncOnce = sync.Once{}
var randomFunctions map[string]*environment.Function

// Global random source for deterministic seeding
var globalRand *rand.Rand
var randMutex sync.Mutex

func init() {
	// Initialize with a random seed for default behavior
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func getRandomFunctions() map[string]*environment.Function {
	randomFuncOnce.Do(func() {
		randomFunctions = map[string]*environment.Function{
			"SEED": {
				Name: "SEED",
				Documentation: []string{
					"Sets the random number generator seed for reproducible results.",
					"Using the same seed will produce the same sequence of random numbers.",
					"",
					"@syntax SEED WIT <seed>",
					"@param {INTEGR} seed - The seed value for the random number generator",
					"@returns {NOTHIN}",
					"@example Set seed for reproducible results",
					"SEED WIT 42",
					"I HAS A VARIABLE NUM1 TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100",
					"I HAS A VARIABLE NUM2 TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100",
					"BTW NUM1 and NUM2 will be the same every time with same seed",
					"@example Different seeds produce different sequences",
					"SEED WIT 123",
					"I HAS A VARIABLE SEQ1 TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 10",
					"SEED WIT 456",
					"I HAS A VARIABLE SEQ2 TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 10",
					"BTW SEQ1 and SEQ2 will likely be different",
					"@example Reset to same sequence",
					"SEED WIT 999",
					"I HAS A VARIABLE W TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100",
					"I HAS A VARIABLE X TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100",
					"SEED WIT 999",
					"I HAS A VARIABLE Y TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100",
					"I HAS A VARIABLE Z TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100",
					"BTW W will equal Y, X will equal Z",
					"@note Use SEED_TIME for non-deterministic behavior",
					"@note Same seed always produces same sequence",
					"@note Useful for testing and debugging",
					"@note Affects all subsequent random operations",
					"@see SEED_TIME, RANDOM_INT",
					"@category seeding",
				},
				Parameters: []environment.Parameter{
					{Name: "seed", Type: "INTEGR"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					seed := args[0]

					if seedVal, ok := seed.(environment.IntegerValue); ok {
						randMutex.Lock()
						globalRand = rand.New(rand.NewSource(int64(seedVal)))
						randMutex.Unlock()
						return environment.NOTHIN, nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "SEED: invalid seed type"}
				},
			},
			"SEED_TIME": {
				Name: "SEED_TIME",
				Documentation: []string{
					"Seeds the random number generator with the current time.",
					"Provides different random sequences on each program run.",
					"",
					"@syntax SEED_TIME",
					"@returns {NOTHIN}",
					"@example Initialize for random behavior",
					"SEED_TIME",
					"I HAS A VARIABLE DICE TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 7",
					"SAYZ WIT \"You rolled: \" MOAR DICE",
					"@example Use at program start",
					"BTW Always call this first for truly random behavior",
					"SEED_TIME",
					"WHILE YEZ",
					"    I HAS A VARIABLE MOVE TEH STRIN ITZ RANDOM_CHOICE WIT MOVES",
					"    SAYZ WIT \"Computer plays: \" MOAR MOVE",
					"    I HAS A VARIABLE CONTINUE TEH BOOL ITZ RANDOM_BOOL",
					"    IZ CONTINUE?",
					"        GTFO",
					"    KTHX",
					"KTHX",
					"@example Multiple runs produce different results",
					"SEED_TIME",
					"I HAS A VARIABLE SESSION_ID TEH STRIN ITZ RANDOM_STRING WIT 8 AN WIT \"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789\"",
					"SAYZ WIT \"Session ID: \" MOAR SESSION_ID",
					"BTW Each program run will have different session ID",
					"@example Games and simulations",
					"SEED_TIME",
					"I HAS A VARIABLE WEATHER TEH STRIN",
					"I HAS A VARIABLE ROLL TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 100",
					"IZ ROLL BIGGR THAN 70?",
					"    WEATHER ITZ \"Sunny\"",
					"NOPE IZ ROLL BIGGR THAN 40?",
					"    WEATHER ITZ \"Cloudy\"",
					"NOPE",
					"    WEATHER ITZ \"Rainy\"",
					"KTHX",
					"SAYZ WIT \"Today's weather: \" MOAR WEATHER",
					"@note Should be called once at program start",
					"@note Uses nanosecond precision for uniqueness",
					"@note Subsequent calls will change the sequence",
					"@note Perfect for games, simulations, and one-time use",
					"@see SEED, RANDOM_INT",
					"@category seeding",
				},
				Parameters: []environment.Parameter{},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					randMutex.Lock()
					globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
					randMutex.Unlock()
					return environment.NOTHIN, nil
				},
			},
			"RANDOM_FLOAT": {
				Name: "RANDOM_FLOAT",
				Documentation: []string{
					"Returns a random floating-point number between 0.0 (inclusive) and 1.0 (exclusive).",
					"Useful for probability calculations and random selection.",
					"",
					"@syntax RANDOM_FLOAT",
					"@returns {DUBBLE} Random float between 0.0 and 1.0",
					"@example Basic random float",
					"I HAS A VARIABLE PROBABILITY TEH DUBBLE ITZ RANDOM_FLOAT",
					"SAYZ WIT \"Random probability: \" MOAR PROBABILITY",
					"@example Probability-based decisions",
					"I HAS A VARIABLE CHANCE TEH DUBBLE ITZ RANDOM_FLOAT",
					"IZ CHANCE BIGGR THAN 0.5?",
					"    SAYZ WIT \"Heads\"",
					"NOPE",
					"    SAYZ WIT \"Tails\"",
					"KTHX",
					"@example Random selection weights",
					"I HAS A VARIABLE WEIGHT TEH DUBBLE ITZ RANDOM_FLOAT",
					"IZ WEIGHT BIGGR THAN 0.7?",
					"    I HAS A VARIABLE RARITY TEH STRIN ITZ \"Legendary\"",
					"NOPE IZ WEIGHT BIGGR THAN 0.4?",
					"    RARITY ITZ \"Rare\"",
					"NOPE",
					"    RARITY ITZ \"Common\"",
					"KTHX",
					"SAYZ WIT \"Got a \" MOAR RARITY MOAR \" item!\"",
					"@example Animation timing",
					"WHILE YEZ",
					"    I HAS A VARIABLE DELAY TEH DUBBLE ITZ RANDOM_FLOAT MOAR 0.1",
					"    BTW Add small random delay for natural feel",
					"    I HAS A VARIABLE FRAME TEH DUBBLE ITZ RANDOM_FLOAT",
					"    IZ FRAME BIGGR THAN 0.95?",
					"        GTFO",
					"    KTHX",
					"KTHX",
					"@note Always returns value >= 0.0 and < 1.0",
					"@note Uses high-quality random number generation",
					"@note Thread-safe for concurrent use",
					"@note Perfect for probabilities and normalized random values",
					"@see RANDOM_RANGE, RANDOM_INT",
					"@category random-numbers",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					randMutex.Lock()
					result := globalRand.Float64()
					randMutex.Unlock()
					return environment.DoubleValue(result), nil
				},
			},
			"RANDOM_RANGE": {
				Name: "RANDOM_RANGE",
				Documentation: []string{
					"Returns a random floating-point number within the specified range.",
					"Range is [min, max) - includes min but excludes max. Min must be less than max.",
					"",
					"@syntax RANDOM_RANGE WIT <min> AN WIT <max>",
					"@param {DUBBLE} min - Minimum value (inclusive)",
					"@param {DUBBLE} max - Maximum value (exclusive)",
					"@returns {DUBBLE} Random float in range [min, max)",
					"@example Random position",
					"I HAS A VARIABLE X_POS TEH DUBBLE ITZ RANDOM_RANGE WIT 0.0 AN WIT 100.0",
					"I HAS A VARIABLE Y_POS TEH DUBBLE ITZ RANDOM_RANGE WIT 0.0 AN WIT 50.0",
					"SAYZ WIT \"Position: (\" MOAR X_POS MOAR \", \" MOAR Y_POS MOAR \")\"",
					"@example Random damage",
					"I HAS A VARIABLE DAMAGE TEH DUBBLE ITZ RANDOM_RANGE WIT 5.0 AN WIT 15.0",
					"SAYZ WIT \"Dealt \" MOAR DAMAGE MOAR \" damage!\"",
					"@example Temperature simulation",
					"I HAS A VARIABLE TEMP TEH DUBBLE ITZ RANDOM_RANGE WIT -10.0 AN WIT 40.0",
					"IZ TEMP BIGGR THAN 30.0?",
					"    SAYZ WIT \"Hot day!\"",
					"NOPE IZ TEMP LSS THAN 0.0?",
					"    SAYZ WIT \"Freezing!\"",
					"NOPE",
					"    SAYZ WIT \"Nice weather\"",
					"KTHX",
					"@example Random delays",
					"I HAS A VARIABLE DELAY TEH DUBBLE ITZ RANDOM_RANGE WIT 0.5 AN WIT 2.0",
					"BTW Wait between 0.5 and 2.0 seconds",
					"SAYZ WIT \"Waiting \" MOAR DELAY MOAR \" seconds...\"",
					"@example Statistical distribution",
					"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
					"WHILE IDX SMALLR THAN 100",
					"    I HAS A VARIABLE VALUE TEH DUBBLE ITZ RANDOM_RANGE WIT -1.0 AN WIT 1.0",
					"    BTW Collect samples for statistical analysis",
					"    IDX ITZ IDX MOAR 1",
					"KTHX",
					"@note Min must be less than max",
					"@note Returns values >= min and < max",
					"@note Useful for continuous random values",
					"@note Thread-safe for concurrent use",
					"@see RANDOM_FLOAT, RANDOM_INT",
					"@category random-numbers",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{
					{Name: "min", Type: "DUBBLE"},
					{Name: "max", Type: "DUBBLE"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					min, max := args[0], args[1]

					if minVal, ok := min.(environment.DoubleValue); ok {
						if maxVal, ok := max.(environment.DoubleValue); ok {
							if minVal >= maxVal {
								return environment.NOTHIN, runtime.Exception{Message: "RANDOM_RANGE: min must be less than max"}
							}
							randMutex.Lock()
							randomVal := globalRand.Float64()
							randMutex.Unlock()
							result := float64(minVal) + randomVal*(float64(maxVal)-float64(minVal))
							return environment.DoubleValue(result), nil
						}
					}

					return environment.NOTHIN, runtime.Exception{Message: "RANDOM_RANGE: invalid numeric arguments"}
				},
			},
			"RANDOM_INT": {
				Name: "RANDOM_INT",
				Documentation: []string{
					"Returns a random integer within the specified range.",
					"Range is [min, max) - includes min but excludes max. Min must be less than max.",
					"",
					"@syntax RANDOM_INT WIT <min> AN WIT <max>",
					"@param {INTEGR} min - Minimum value (inclusive)",
					"@param {INTEGR} max - Maximum value (exclusive)",
					"@returns {INTEGR} Random integer in range [min, max)",
					"@example Dice roll",
					"I HAS A VARIABLE DICE TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 7",
					"SAYZ WIT \"You rolled a \" MOAR DICE",
					"@example Random index",
					"I HAS A VARIABLE INDEX TEH INTEGR ITZ RANDOM_INT WIT 0 AN WIT 10",
					"SAYZ WIT \"Selected index: \" MOAR INDEX",
					"@example Random ID generation",
					"I HAS A VARIABLE USER_ID TEH INTEGR ITZ RANDOM_INT WIT 1000 AN WIT 10000",
					"SAYZ WIT \"New user ID: \" MOAR USER_ID",
					"@example Game mechanics",
					"I HAS A VARIABLE DAMAGE TEH INTEGR ITZ RANDOM_INT WIT 10 AN WIT 21",
					"I HAS A VARIABLE CRIT_CHANCE TEH INTEGR ITZ RANDOM_INT WIT 1 AN WIT 101",
					"IZ CRIT_CHANCE BIGGR THAN 90?",
					"    DAMAGE ITZ DAMAGE UP 10",
					"    SAYZ WIT \"Critical hit!\"",
					"KTHX",
					"SAYZ WIT \"Dealt \" MOAR DAMAGE MOAR \" damage\"",
					"@example Array shuffling indices",
					"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
					"WHILE IDX SMALLR THAN 5",
					"    I HAS A VARIABLE POS TEH INTEGR ITZ RANDOM_INT WIT 0 AN WIT 5",
					"    BTW Use for Fisher-Yates shuffle",
					"    IDX ITZ IDX MOAR 1",
					"KTHX",
					"@note Min must be less than max",
					"@note Returns integers >= min and < max",
					"@note Useful for discrete random selections",
					"@note Thread-safe for concurrent use",
					"@see RANDOM_RANGE, RANDOM_CHOICE",
					"@category random-numbers",
				},
				ReturnType: "INTEGR",
				Parameters: []environment.Parameter{
					{Name: "min", Type: "INTEGR"},
					{Name: "max", Type: "INTEGR"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					min, max := args[0], args[1]

					if minVal, ok := min.(environment.IntegerValue); ok {
						if maxVal, ok := max.(environment.IntegerValue); ok {
							if minVal >= maxVal {
								return environment.NOTHIN, runtime.Exception{Message: "RANDOM_INT: min must be less than max"}
							}
							randMutex.Lock()
							result := globalRand.Int63n(int64(maxVal-minVal)) + int64(minVal)
							randMutex.Unlock()
							return environment.IntegerValue(result), nil
						}
					}

					return environment.NOTHIN, runtime.Exception{Message: "RANDOM_INT: invalid integer arguments"}
				},
			},
			"RANDOM_BOOL": {
				Name: "RANDOM_BOOL",
				Documentation: []string{
					"Returns a random boolean value (YEZ or NO).",
					"Each value has a 50% probability of being returned.",
					"",
					"@syntax RANDOM_BOOL",
					"@returns {BOOL} Random boolean (YEZ or NO)",
					"@example Coin flip",
					"I HAS A VARIABLE COIN TEH BOOL ITZ RANDOM_BOOL",
					"IZ COIN?",
					"    SAYZ WIT \"Heads\"",
					"NOPE",
					"    SAYZ WIT \"Tails\"",
					"KTHX",
					"@example Random events",
					"I HAS A VARIABLE RAIN TEH BOOL ITZ RANDOM_BOOL",
					"IZ RAIN?",
					"    SAYZ WIT \"It's raining today\"",
					"NOPE",
					"    SAYZ WIT \"Sunny weather\"",
					"KTHX",
					"@example Game mechanics",
					"I HAS A VARIABLE SUCCESS TEH BOOL ITZ RANDOM_BOOL",
					"IZ SUCCESS?",
					"    SAYZ WIT \"Action succeeded!\"",
					"NOPE",
					"    SAYZ WIT \"Action failed!\"",
					"KTHX",
					"@example Random spawning",
					"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
					"WHILE IDX SMALLR THAN 10",
					"    I HAS A VARIABLE SPAWN_ENEMY TEH BOOL ITZ RANDOM_BOOL",
					"    IZ SPAWN_ENEMY?",
					"        SAYZ WIT \"Enemy spawned!\"",
					"    KTHX",
					"    IDX ITZ IDX MOAR 1",
					"KTHX",
					"@example A/B testing",
					"I HAS A VARIABLE USE_NEW_FEATURE TEH BOOL ITZ RANDOM_BOOL",
					"IZ USE_NEW_FEATURE?",
					"    SAYZ WIT \"Using new feature version\"",
					"NOPE",
					"    SAYZ WIT \"Using old feature version\"",
					"KTHX",
					"@note Each outcome has exactly 50% probability",
					"@note Useful for binary random decisions",
					"@note Thread-safe for concurrent use",
					"@note Perfect for yes/no scenarios",
					"@see RANDOM_INT, RANDOM_CHOICE",
					"@category random-selection",
				},
				ReturnType: "BOOL",
				Parameters: []environment.Parameter{},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					randMutex.Lock()
					randomVal := globalRand.Float64()
					randMutex.Unlock()
					if randomVal < 0.5 {
						return environment.NO, nil
					}
					return environment.YEZ, nil
				},
			},
			"RANDOM_CHOICE": {
				Name: "RANDOM_CHOICE",
				Documentation: []string{
					"Returns a randomly selected element from a BUKKIT array.",
					"Array must not be empty. Each element has equal probability of selection.",
					"",
					"@syntax RANDOM_CHOICE WIT <array>",
					"@param {BUKKIT} array - The array to select from",
					"@returns Random element from the array",
					"@example Random card from deck",
					"I HAS A VARIABLE CARDS TEH BUKKIT ITZ BUKKIT WIT \"Ace\" AN \"King\" AN \"Queen\" AN \"Jack\"",
					"I HAS A VARIABLE CARD TEH STRIN ITZ RANDOM_CHOICE WIT CARDS",
					"SAYZ WIT \"Drew: \" MOAR CARD",
					"@example Random enemy type",
					"I HAS A VARIABLE ENEMIES TEH BUKKIT ITZ BUKKIT WIT \"Goblin\" AN \"Orc\" AN \"Troll\" AN \"Dragon\"",
					"I HAS A VARIABLE ENEMY TEH STRIN ITZ RANDOM_CHOICE WIT ENEMIES",
					"SAYZ WIT \"Encountered a \" MOAR ENEMY",
					"@example Random quote",
					"I HAS A VARIABLE QUOTES TEH BUKKIT ITZ BUKKIT WIT \"Hello World\" AN \"LOL\" AN \"Objective-C\"",
					"I HAS A VARIABLE QUOTE TEH STRIN ITZ RANDOM_CHOICE WIT QUOTES",
					"SAYZ QUOTE",
					"@example Game loot table",
					"I HAS A VARIABLE LOOT TEH BUKKIT ITZ BUKKIT WIT \"Sword\" AN \"Shield\" AN \"Potion\" AN \"Gold\" AN \"Key\"",
					"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
					"WHILE IDX SMALLR THAN 3",
					"    I HAS A VARIABLE ITEM TEH STRIN ITZ RANDOM_CHOICE WIT LOOT",
					"    SAYZ WIT \"Found: \" MOAR ITEM",
					"    IDX ITZ IDX MOAR 1",
					"KTHX",
					"@example Random color",
					"I HAS A VARIABLE COLORS TEH BUKKIT ITZ BUKKIT WIT \"Red\" AN \"Blue\" AN \"Green\" AN \"Yellow\" AN \"Purple\"",
					"I HAS A VARIABLE COLOR TEH STRIN ITZ RANDOM_CHOICE WIT COLORS",
					"SAYZ WIT \"Selected color: \" MOAR COLOR",
					"@note Array must not be empty",
					"@note All elements have equal probability",
					"@note Returns element by reference",
					"@note Thread-safe for concurrent use",
					"@see SHUFFLE, RANDOM_INT",
					"@category random-selection",
				},
				Parameters: []environment.Parameter{
					{Name: "array", Type: "BUKKIT"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					array := args[0]

					if arrayObj, ok := array.(*environment.ObjectInstance); ok {
						if slice, ok := arrayObj.NativeData.(BukkitSlice); ok {
							if len(slice) == 0 {
								return environment.NOTHIN, runtime.Exception{Message: "RANDOM_CHOICE: empty array"}
							}
							randMutex.Lock()
							index := globalRand.Intn(len(slice))
							randMutex.Unlock()
							return slice[index], nil
						}
					}

					return environment.NOTHIN, runtime.Exception{Message: "RANDOM_CHOICE: invalid array argument"}
				},
			},
			"SHUFFLE": {
				Name: "SHUFFLE",
				Documentation: []string{
					"Returns a new BUKKIT with elements randomly shuffled.",
					"Uses Fisher-Yates algorithm. Original array is not modified.",
					"",
					"@syntax SHUFFLE WIT <array>",
					"@param {BUKKIT} array - The array to shuffle",
					"@returns {BUKKIT} New shuffled array",
					"@example Shuffle a deck of cards",
					"I HAS A VARIABLE DECK TEH BUKKIT ITZ BUKKIT WIT \"A\" AN \"2\" AN \"3\" AN \"4\" AN \"5\"",
					"I HAS A VARIABLE SHUFFLED TEH BUKKIT ITZ SHUFFLE WIT DECK",
					"SAYZ WIT \"Original: \" MOAR DECK",
					"SAYZ WIT \"Shuffled: \" MOAR SHUFFLED",
					"@example Randomize playlist",
					"I HAS A VARIABLE SONGS TEH BUKKIT ITZ BUKKIT WIT \"Song1\" AN \"Song2\" AN \"Song3\" AN \"Song4\"",
					"I HAS A VARIABLE PLAYLIST TEH BUKKIT ITZ SHUFFLE WIT SONGS",
					"SAYZ WIT \"Random playlist: \" MOAR PLAYLIST",
					"@example Random game order",
					"I HAS A VARIABLE PLAYERS TEH BUKKIT ITZ BUKKIT WIT \"Alice\" AN \"Bob\" AN \"Charlie\" AN \"Dave\"",
					"I HAS A VARIABLE ORDER TEH BUKKIT ITZ SHUFFLE WIT PLAYERS",
					"SAYZ WIT \"Playing order: \" MOAR ORDER",
					"@example Randomize questions",
					"I HAS A VARIABLE QUESTIONS TEH BUKKIT ITZ BUKKIT WIT \"Q1\" AN \"Q2\" AN \"Q3\" AN \"Q4\" AN \"Q5\"",
					"I HAS A VARIABLE QUIZ TEH BUKKIT ITZ SHUFFLE WIT QUESTIONS",
					"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
					"WHILE IDX SMALLR THAN 5",
					"    SAYZ WIT \"Question: \" MOAR QUIZ DO AT IDX",
					"    IDX ITZ IDX MOAR 1",
					"KTHX",
					"@example Randomize positions",
					"I HAS A VARIABLE POSITIONS TEH BUKKIT ITZ BUKKIT WIT 1 AN 2 AN 3 AN 4 AN 5",
					"I HAS A VARIABLE RANDOM_POS TEH BUKKIT ITZ SHUFFLE WIT POSITIONS",
					"SAYZ WIT \"Random positions: \" MOAR RANDOM_POS",
					"@note Original array remains unchanged",
					"@note Returns new BUKKIT instance",
					"@note Uses Fisher-Yates shuffle algorithm",
					"@note Thread-safe for concurrent use",
					"@see RANDOM_CHOICE, RANDOM_INT",
					"@category random-utilities",
				},
				ReturnType: "BUKKIT",
				Parameters: []environment.Parameter{
					{Name: "array", Type: "BUKKIT"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					array := args[0]

					if arrayObj, ok := array.(*environment.ObjectInstance); ok {
						if slice, ok := arrayObj.NativeData.(BukkitSlice); ok {
							// Create a copy to avoid modifying original
							shuffled := make(BukkitSlice, len(slice))
							copy(shuffled, slice)

							// Fisher-Yates shuffle
							for i := len(shuffled) - 1; i > 0; i-- {
								randMutex.Lock()
								j := globalRand.Intn(i + 1)
								randMutex.Unlock()
								shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
							}

							// Create new BUKKIT with shuffled data
							newObj := NewBukkitInstance()
							newObj.NativeData = shuffled
							return newObj, nil
						}
					}

					return environment.NOTHIN, runtime.Exception{Message: "SHUFFLE: invalid array argument"}
				},
			},
			"RANDOM_STRING": {
				Name: "RANDOM_STRING",
				Documentation: []string{
					"Generates a random string of specified length using given character set.",
					"Each character is randomly selected from the charset. Charset must not be empty.",
					"",
					"@syntax RANDOM_STRING WIT <length> AN WIT <charset>",
					"@param {INTEGR} length - Length of the generated string",
					"@param {STRIN} charset - Characters to choose from",
					"@returns {STRIN} Random string of specified length",
					"@example Random password",
					"I HAS A VARIABLE PASSWORD TEH STRIN ITZ RANDOM_STRING WIT 8 AN WIT \"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789\"",
					"SAYZ WIT \"Generated password: \" MOAR PASSWORD",
					"@example Random hex color",
					"I HAS A VARIABLE COLOR TEH STRIN ITZ RANDOM_STRING WIT 6 AN WIT \"0123456789ABCDEF\"",
					"SAYZ WIT \"Random color: #\" MOAR COLOR",
					"@example Random ID",
					"I HAS A VARIABLE ID TEH STRIN ITZ RANDOM_STRING WIT 10 AN WIT \"abcdefghijklmnopqrstuvwxyz0123456789\"",
					"SAYZ WIT \"Random ID: \" MOAR ID",
					"@example Random letters only",
					"I HAS A VARIABLE WORD TEH STRIN ITZ RANDOM_STRING WIT 5 AN WIT \"ABCDEFGHIJKLMNOPQRSTUVWXYZ\"",
					"SAYZ WIT \"Random word: \" MOAR WORD",
					"@example Random binary string",
					"I HAS A VARIABLE BINARY TEH STRIN ITZ RANDOM_STRING WIT 16 AN WIT \"01\"",
					"SAYZ WIT \"Random binary: \" MOAR BINARY",
					"@example Random emoji",
					"I HAS A VARIABLE EMOJI_CHARSET TEH STRIN ITZ \"😀😃😄😁😆😅😂🤣😊😇🙂🙃😉😌😍🥰😘😗😙😚😋😛😝😜🤪🤨🧐🤓😎🤩🥳😏😒😞😔😟😕🙁☹️😣😖😫😩🥺😢😭😤😠😡🤬🤯😳🥵🥶😱😨😰😥😓🤗🤔🤭🤫🤥😶😐😑😬🙄😯😦😧😮😲🥱😴🤤😪😵🤐🥴🤢🤮🤧😷🤒🤕🤑🤠😈👿👹👺🤡💩👻💀☠️👽👾🤖🎃😺😸😹😻😼😽🙀😿😾\"",
					"I HAS A VARIABLE RANDOM_EMOJI TEH STRIN ITZ RANDOM_STRING WIT 1 AN WIT EMOJI_CHARSET",
					"SAYZ WIT \"Random emoji: \" MOAR RANDOM_EMOJI",
					"@note Length must be >= 0",
					"@note Charset must not be empty",
					"@note Each character has equal probability",
					"@note Thread-safe for concurrent use",
					"@see RANDOM_INT, UUID",
					"@category random-generation",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{
					{Name: "length", Type: "INTEGR"},
					{Name: "charset", Type: "STRIN"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					length := args[0]
					charset := args[1]

					if lengthVal, ok := length.(environment.IntegerValue); ok {
						if charsetVal, ok := charset.(environment.StringValue); ok {
							if lengthVal <= 0 {
								return environment.StringValue(""), nil
							}

							charsetStr := string(charsetVal)
							if len(charsetStr) == 0 {
								return environment.NOTHIN, runtime.Exception{Message: "RANDOM_STRING: empty charset"}
							}

							result := make([]byte, lengthVal)
							for i := range result {
								randMutex.Lock()
								randomIndex := globalRand.Intn(len(charsetStr))
								randMutex.Unlock()
								result[i] = charsetStr[randomIndex]
							}

							return environment.StringValue(string(result)), nil
						}
					}

					return environment.NOTHIN, runtime.Exception{Message: "RANDOM_STRING: invalid arguments"}
				},
			},
			"UUID": {
				Name: "UUID",
				Documentation: []string{
					"Generates a random UUID (Universally Unique Identifier) version 4.",
					"Returns a string in format: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx.",
					"",
					"@syntax UUID",
					"@returns {STRIN} Random UUID v4 string",
					"@example Generate unique ID",
					"I HAS A VARIABLE ID TEH STRIN ITZ UUID",
					"SAYZ WIT \"Generated UUID: \" MOAR ID",
					"@example User session ID",
					"I HAS A VARIABLE SESSION_ID TEH STRIN ITZ UUID",
					"SAYZ WIT \"Session ID: \" MOAR SESSION_ID",
					"@example Database primary key",
					"I HAS A VARIABLE PRIMARY_KEY TEH STRIN ITZ UUID",
					"SAYZ WIT \"New record ID: \" MOAR PRIMARY_KEY",
					"@example File naming",
					"I HAS A VARIABLE FILENAME TEH STRIN ITZ UUID MOAR \".txt\"",
					"SAYZ WIT \"Creating file: \" MOAR FILENAME",
					"@example Object identification",
					"I HAS A VARIABLE OBJECT_ID TEH STRIN ITZ UUID",
					"SAYZ WIT \"Object created with ID: \" MOAR OBJECT_ID",
					"@example API key generation",
					"I HAS A VARIABLE API_KEY TEH STRIN ITZ UUID",
					"SAYZ WIT \"Generated API key: \" MOAR API_KEY",
					"@example Transaction ID",
					"I HAS A VARIABLE TXN_ID TEH STRIN ITZ UUID",
					"SAYZ WIT \"Transaction ID: \" MOAR TXN_ID",
					"@note Uses cryptographically secure random generation",
					"@note Globally unique with extremely high probability",
					"@note Version 4 format (random)",
					"@note Thread-safe for concurrent use",
					"@see RANDOM_STRING, RANDOM_INT",
					"@category random-generation",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					// Simple UUID v4 implementation
					uuid := make([]byte, 16)
					_, err := cryptorand.Read(uuid)
					if err != nil {
						return environment.NOTHIN, runtime.Exception{Message: "UUID: failed to generate random bytes"}
					}

					// Set version (4) and variant bits
					uuid[6] = (uuid[6] & 0x0f) | 0x40
					uuid[8] = (uuid[8] & 0x3f) | 0x80

					result := fmt.Sprintf("%x-%x-%x-%x-%x",
						uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])

					return environment.StringValue(result), nil
				},
			},
		}
	})
	return randomFunctions
}

// RegisterRANDOMInEnv registers RANDOM functions in the given environment
// declarations: empty slice means import all, otherwise import only specified functions
func RegisterRANDOMInEnv(env *environment.Environment, declarations ...string) error {
	randomFunctions := getRandomFunctions()

	// If declarations is empty, import all functions
	if len(declarations) == 0 {
		for _, fn := range randomFunctions {
			env.DefineFunction(fn)
		}
		return nil
	}

	// Otherwise, import only specified functions
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if fn, exists := randomFunctions[declUpper]; exists {
			env.DefineFunction(fn)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown RANDOM declaration: %s", decl)}
		}
	}

	return nil
}
