package stdlib

import (
	"fmt"
	"sort"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleMapsCategories defines the order that categories should be rendered in documentation
var moduleMapsCategories = []string{
	"map-creation",
	"map-access",
	"map-modification",
	"map-inspection",
	"map-transformation",
}

type BaskitMap map[string]environment.Value

func NewBaskitInstance() *environment.ObjectInstance {
	class := getMapClasses()["BASKIT"]
	env := environment.NewEnvironment(nil)
	env.DefineClass(class)
	obj := &environment.ObjectInstance{
		Environment: env,
		Class:       class,
		NativeData:  make(BaskitMap),
		Variables:   make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(obj)
	return obj
}

// Global BASKIT class definition - created once and reused
var mapClassesOnce = sync.Once{}
var mapClasses map[string]*environment.Class

func getMapClasses() map[string]*environment.Class {
	mapClassesOnce.Do(func() {
		mapClasses = map[string]*environment.Class{
			"BASKIT": {
				Name: "BASKIT",
				Documentation: []string{
					"A dynamic map (dictionary) that stores key-value pairs.",
					"Keys are STRIN type and values can be any type.",
					"Provides fast lookup, insertion, and deletion of key-value pairs.",
					"",
					"@class BASKIT",
					"@example Create empty map",
					"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
					"BTW Creates an empty dictionary",
					"@example Store and retrieve values",
					"MAP DO PUT WIT \"name\" AN WIT \"Alice\"",
					"MAP DO PUT WIT \"age\" AN WIT 25",
					"I HAS A VARIABLE NAME TEH STRIN ITZ MAP DO GET WIT \"name\"",
					"BTW NAME = \"Alice\"",
					"@example Mixed value types",
					"MAP DO PUT WIT \"count\" AN WIT 100",
					"MAP DO PUT WIT \"active\" AN WIT YEZ",
					"MAP DO PUT WIT \"items\" AN WIT NEW BUKKIT",
					"BTW Values can be any type: STRIN, INTEGR, BOOL, BUKKIT, etc.",
					"@note Keys are always strings, values can be any type",
					"@note Fast O(1) average case lookup performance",
					"@note Keys are case-sensitive",
					"@see BUKKIT",
					"@category map-creation",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BASKIT": {
						Name: "BASKIT",
						Documentation: []string{
							"Initializes an empty BASKIT map.",
							"Creates a new dictionary with no key-value pairs.",
							"",
							"@syntax NEW BASKIT",
							"@returns {BASKIT} A new empty BASKIT instance",
							"@example Create empty map",
							"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
							"SAYZ WIT MAP SIZ",
							"BTW Output: 0",
							"@note Creates an empty dictionary ready for key-value pairs",
							"@see PUT, GET",
							"@category map-creation",
						},
						Parameters: []environment.Parameter{}, // Empty constructor - no arguments
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							// Create empty map and store in NativeData
							baskitMap := make(BaskitMap)
							this.NativeData = baskitMap
							return environment.NOTHIN, nil
						},
					},
					// Map methods
					"PUT": {
						Name: "PUT",
						Documentation: []string{
							"Stores a key-value pair in the BASKIT.",
							"If key already exists, overwrites the existing value.",
							"",
							"@syntax map DO PUT WIT <key> AN WIT <value>",
							"@param {STRIN} key - The key to store (converted to string)",
							"@param {ANY} value - The value to associate with the key",
							"@returns {NOTHIN} No return value",
							"@example Add new entries",
							"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
							"MAP DO PUT WIT \"name\" AN WIT \"Alice\"",
							"MAP DO PUT WIT \"age\" AN WIT 25",
							"BTW MAP now contains name->Alice and age->25",
							"@example Overwrite existing key",
							"MAP DO PUT WIT \"age\" AN WIT 26",
							"BTW age is now 26 (was 25)",
							"@example Mixed value types",
							"MAP DO PUT WIT \"active\" AN WIT YEZ",
							"MAP DO PUT WIT \"items\" AN WIT NEW BUKKIT",
							"BTW Values can be any type",
							"@note Key is converted to string if not already",
							"@note Overwrites existing values for the same key",
							"@see GET, CONTAINS",
							"@category map-modification",
						},
						Parameters: []environment.Parameter{{Name: "KEY", Type: "STRIN"}, {Name: "VALUE", Type: ""}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							key, value := args[0], args[1]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Convert key to string
								keyStr := key.String()
								baskitMap[keyStr] = value
								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "PUT: invalid context"}
						},
					},
					"GET": {
						Name: "GET",
						Documentation: []string{
							"Retrieves the value associated with the specified key.",
							"Throws an exception if the key is not found.",
							"",
							"@syntax map DO GET WIT <key>",
							"@param {STRIN} key - The key to look up",
							"@returns {ANY} The value associated with the key",
							"@example Retrieve values",
							"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
							"MAP DO PUT WIT \"name\" AN WIT \"Bob\"",
							"MAP DO PUT WIT \"score\" AN WIT 95",
							"I HAS A VARIABLE NAME TEH STRIN ITZ MAP DO GET WIT \"name\"",
							"I HAS A VARIABLE SCORE TEH INTEGR ITZ MAP DO GET WIT \"score\"",
							"BTW NAME = \"Bob\", SCORE = 95",
							"@example Handle missing key",
							"MAYB",
							"    I HAS A VARIABLE MISSING TEH STRIN ITZ MAP DO GET WIT \"missing\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Key not found!\"",
							"KTHX",
							"@throws Key not found exception if key doesn't exist",
							"@note Use CONTAINS to check if key exists first",
							"@see PUT, CONTAINS",
							"@category map-access",
						},
						Parameters: []environment.Parameter{{Name: "KEY", Type: "STRIN"}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							key := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Convert key to string
								keyStr := key.String()
								if value, exists := baskitMap[keyStr]; exists {
									return value, nil
								}
								return nil, runtime.Exception{Message: fmt.Sprintf("BASKIT: key %q not found", keyStr)}
							}
							return environment.NOTHIN, runtime.Exception{Message: "GET: invalid context"}
						},
					},
					"CONTAINS": {
						Name: "CONTAINS",
						Documentation: []string{
							"Checks if the specified key exists in the BASKIT.",
							"Returns YEZ if key exists, NO otherwise.",
							"",
							"@syntax map DO CONTAINS WIT <key>",
							"@param {STRIN} key - The key to check for existence",
							"@returns {BOOL} YEZ if key exists, NO otherwise",
							"@example Check for keys",
							"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
							"MAP DO PUT WIT \"name\" AN WIT \"Carol\"",
							"I HAS A VARIABLE HAS_NAME TEH BOOL ITZ MAP DO CONTAINS WIT \"name\"",
							"I HAS A VARIABLE HAS_AGE TEH BOOL ITZ MAP DO CONTAINS WIT \"age\"",
							"BTW HAS_NAME = YEZ, HAS_AGE = NO",
							"@example Use in conditional",
							"IZ MAP DO CONTAINS WIT \"score\"?",
							"    I HAS A VARIABLE SCORE TEH INTEGR ITZ MAP DO GET WIT \"score\"",
							"NOPE",
							"    SAYZ WIT \"Score not set\"",
							"KTHX",
							"@note Safer than GET when you only need to check existence",
							"@see GET, KEYS",
							"@category map-inspection",
						},
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{{Name: "KEY", Type: "STRIN"}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							key := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Convert key to string
								keyStr := key.String()
								_, exists := baskitMap[keyStr]
								if exists {
									return environment.YEZ, nil
								}
								return environment.NO, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "CONTAINS: invalid context"}
						},
					},
					"REMOVE": {
						Name: "REMOVE",
						Documentation: []string{
							"Removes a key-value pair from the BASKIT and returns the value.",
							"Throws an exception if the key is not found.",
							"",
							"@syntax map DO REMOVE WIT <key>",
							"@param {STRIN} key - The key to remove",
							"@returns {ANY} The value that was removed",
							"@example Remove entries",
							"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
							"MAP DO PUT WIT \"temp\" AN WIT 42",
							"MAP DO PUT WIT \"keep\" AN WIT \"important\"",
							"I HAS A VARIABLE REMOVED TEH INTEGR ITZ MAP DO REMOVE WIT \"temp\"",
							"BTW REMOVED = 42, MAP now only contains \"keep\"",
							"@example Safe removal",
							"IZ MAP DO CONTAINS WIT \"old_key\"?",
							"    I HAS A VARIABLE OLD_VAL TEH STRIN ITZ MAP DO REMOVE WIT \"old_key\"",
							"KTHX",
							"@throws Key not found exception if key doesn't exist",
							"@note Use CONTAINS to check before removing if unsure",
							"@see PUT, CONTAINS, CLEAR",
							"@category map-modification",
						},
						Parameters: []environment.Parameter{{Name: "KEY", Type: "STRIN"}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							key := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Convert key to string
								keyStr := key.String()
								if value, exists := baskitMap[keyStr]; exists {
									delete(baskitMap, keyStr)
									return value, nil
								}
								return nil, runtime.Exception{Message: fmt.Sprintf("BASKIT: key %q not found", keyStr)}
							}
							return environment.NOTHIN, runtime.Exception{Message: "REMOVE: invalid context"}
						},
					},
					"CLEAR": {
						Name: "CLEAR",
						Documentation: []string{
							"Removes all key-value pairs from the BASKIT.",
							"After clearing, the map size will be 0.",
							"",
							"@syntax map DO CLEAR",
							"@returns {NOTHIN} No return value",
							"@example Clear all entries",
							"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
							"MAP DO PUT WIT \"a\" AN WIT 1",
							"MAP DO PUT WIT \"b\" AN WIT 2",
							"SAYZ WIT MAP SIZ",
							"BTW Output: 2",
							"MAP DO CLEAR",
							"SAYZ WIT MAP SIZ",
							"BTW Output: 0",
							"@example Reset for reuse",
							"MAP DO CLEAR",
							"MAP DO PUT WIT \"fresh\" AN WIT \"start\"",
							"BTW MAP is now empty and ready for new data",
							"@note More efficient than creating a new BASKIT",
							"@note Keeps the same object but removes all contents",
							"@see SIZ, REMOVE",
							"@category map-modification",
						},
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							newMap := make(BaskitMap)
							this.NativeData = newMap
							return environment.NOTHIN, nil
						},
					},
					"KEYS": {
						Name: "KEYS",
						Documentation: []string{
							"Returns a BUKKIT containing all keys in the BASKIT.",
							"Keys are sorted alphabetically for consistent ordering.",
							"",
							"@syntax map DO KEYS",
							"@returns {BUKKIT} Array of all keys as strings",
							"@example Get all keys",
							"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
							"MAP DO PUT WIT \"zebra\" AN WIT 1",
							"MAP DO PUT WIT \"apple\" AN WIT 2",
							"MAP DO PUT WIT \"banana\" AN WIT 3",
							"I HAS A VARIABLE KEYS TEH BUKKIT ITZ MAP DO KEYS",
							"BTW KEYS = [\"apple\", \"banana\", \"zebra\"] (alphabetical)",
							"@example Iterate over keys",
							"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
							"WHILE IDX SMALLR THAN KEYS SIZ",
							"    I HAS A VARIABLE KEY TEH STRIN ITZ KEYS DO AT WIT IDX",
							"    I HAS A VARIABLE VALUE TEH INTEGR ITZ MAP DO GET WIT KEY",
							"    SAYZ WIT KEY MOAR \": \" MOAR VALUE",
							"    IDX ITZ IDX MOAR 1",
							"KTHX",
							"@note Keys are always sorted alphabetically",
							"@note Returns empty BUKKIT if map is empty",
							"@see VALUES, PAIRS",
							"@category map-inspection",
						},
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Create a new BUKKIT with all keys
								keys := make([]string, 0, len(baskitMap))
								for k := range baskitMap {
									keys = append(keys, k)
								}
								// Sort keys for consistent ordering
								sort.Strings(keys)

								// Create BUKKIT instance
								bukkitObj := NewBukkitInstance()
								bukkitSlice := make(BukkitSlice, 0, len(keys))
								for _, k := range keys {
									bukkitSlice = append(bukkitSlice, environment.StringValue(k))
								}
								bukkitObj.NativeData = bukkitSlice
								return bukkitObj, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "KEYS: invalid context"}
						},
					},
					"VALUES": {
						Name: "VALUES",
						Documentation: []string{
							"Returns a BUKKIT containing all values in the BASKIT.",
							"Values are ordered according to their keys' alphabetical order.",
							"",
							"@syntax map DO VALUES",
							"@returns {BUKKIT} Array of all values in key-alphabetical order",
							"@example Get all values",
							"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
							"MAP DO PUT WIT \"c\" AN WIT \"third\"",
							"MAP DO PUT WIT \"a\" AN WIT \"first\"",
							"MAP DO PUT WIT \"b\" AN WIT \"second\"",
							"I HAS A VARIABLE VALUES TEH BUKKIT ITZ MAP DO VALUES",
							"BTW VALUES = [\"first\", \"second\", \"third\"] (by key order: a, b, c)",
							"@example Process all values",
							"I HAS A VARIABLE TOTAL TEH INTEGR ITZ 0",
							"I HAS A VARIABLE SCORES TEH BUKKIT ITZ MAP DO VALUES",
							"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
							"WHILE IDX SMALLR THAN SCORES SIZ",
							"    TOTAL ITZ TOTAL MOAR (SCORES DO AT WIT IDX)",
							"    IDX ITZ IDX MOAR 1",
							"KTHX",
							"@note Values are ordered by their keys' alphabetical order",
							"@note Returns empty BUKKIT if map is empty",
							"@see KEYS, PAIRS",
							"@category map-inspection",
						},
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Create a new BUKKIT with all values
								// Sort keys first for consistent ordering
								keys := make([]string, 0, len(baskitMap))
								for k := range baskitMap {
									keys = append(keys, k)
								}
								sort.Strings(keys)

								// Create BUKKIT instance with values in key order
								bukkitObj := NewBukkitInstance()
								bukkitSlice := make(BukkitSlice, 0, len(baskitMap))
								for _, k := range keys {
									bukkitSlice = append(bukkitSlice, baskitMap[k])
								}
								bukkitObj.NativeData = bukkitSlice
								return bukkitObj, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "VALUES: invalid context"}
						},
					},
					"PAIRS": {
						Name: "PAIRS",
						Documentation: []string{
							"Returns a BUKKIT of key-value pairs as BUKKITs containing [key, value].",
							"Useful for iterating over both keys and values simultaneously.",
							"",
							"@syntax map DO PAIRS",
							"@returns {BUKKIT} Array of [key, value] pairs as BUKKITs",
							"@example Get key-value pairs",
							"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
							"MAP DO PUT WIT \"name\" AN WIT \"David\"",
							"MAP DO PUT WIT \"age\" AN WIT 30",
							"I HAS A VARIABLE PAIRS TEH BUKKIT ITZ MAP DO PAIRS",
							"BTW PAIRS = [[\"age\", 30], [\"name\", \"David\"]] (by key order)",
							"@example Iterate over pairs",
							"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
							"WHILE IDX SMALLR THAN PAIRS SIZ",
							"    I HAS A VARIABLE PAIR TEH BUKKIT ITZ PAIRS DO AT WIT IDX",
							"    I HAS A VARIABLE KEY TEH STRIN ITZ PAIR DO AT WIT 0",
							"    I HAS A VARIABLE VALUE TEH STRIN ITZ PAIR DO AT WIT 1",
							"    SAYZ WIT KEY MOAR \": \" MOAR VALUE",
							"    IDX ITZ IDX MOAR 1",
							"KTHX",
							"@note Each pair is a BUKKIT with [key, value]",
							"@note Pairs are ordered by key alphabetically",
							"@see KEYS, VALUES",
							"@category map-inspection",
						},
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Create a new BUKKIT with key-value pairs as BUKKIT objects
								keys := make([]string, 0, len(baskitMap))
								for k := range baskitMap {
									keys = append(keys, k)
								}
								sort.Strings(keys)

								// Create BUKKIT instance with pairs
								bukkitObj := NewBukkitInstance()
								bukkitSlice := make(BukkitSlice, 0, len(baskitMap))
								for _, k := range keys {
									// Create a pair as a BUKKIT with [key, value]
									pairObj := NewBukkitInstance()
									pairSlice := BukkitSlice{environment.StringValue(k), baskitMap[k]}
									pairObj.NativeData = pairSlice
									bukkitSlice = append(bukkitSlice, pairObj)
								}
								bukkitObj.NativeData = bukkitSlice
								return bukkitObj, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "PAIRS: invalid context"}
						},
					},
					"MERGE": {
						Name: "MERGE",
						Documentation: []string{
							"Merges another BASKIT's key-value pairs into this BASKIT.",
							"Existing keys are overwritten with values from the other BASKIT.",
							"",
							"@syntax map DO MERGE WIT <other_baskit>",
							"@param {BASKIT} other - Another BASKIT to merge from",
							"@returns {NOTHIN} No return value",
							"@example Merge maps",
							"I HAS A VARIABLE MAP1 TEH BASKIT ITZ NEW BASKIT",
							"I HAS A VARIABLE MAP2 TEH BASKIT ITZ NEW BASKIT",
							"MAP1 DO PUT WIT \"a\" AN WIT 1",
							"MAP1 DO PUT WIT \"b\" AN WIT 2",
							"MAP2 DO PUT WIT \"b\" AN WIT 99",
							"MAP2 DO PUT WIT \"c\" AN WIT 3",
							"MAP1 DO MERGE WIT MAP2",
							"BTW MAP1 now contains: a->1, b->99, c->3 (b was overwritten)",
							"@example Configuration merging",
							"I HAS A VARIABLE DEFAULTS TEH BASKIT ITZ NEW BASKIT",
							"I HAS A VARIABLE USER_CONFIG TEH BASKIT ITZ NEW BASKIT",
							"DEFAULTS DO PUT WIT \"timeout\" AN WIT 30",
							"DEFAULTS DO PUT WIT \"retries\" AN WIT 3",
							"USER_CONFIG DO PUT WIT \"timeout\" AN WIT 60",
							"DEFAULTS DO MERGE WIT USER_CONFIG",
							"BTW DEFAULTS now has user's timeout but default retries",
							"@note Modifies the original BASKIT",
							"@note Overwrites existing keys with new values",
							"@see COPY, PUT",
							"@category map-transformation",
						},
						Parameters: []environment.Parameter{{Name: "OTHER", Type: "BASKIT"}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							other := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Extract other BASKIT
								otherInstance, ok := other.(*environment.ObjectInstance)
								if !ok {
									return nil, runtime.Exception{Message: fmt.Sprintf("MERGE expects BASKIT argument, got %s", other.Type())}
								}
								otherBaskitMap, ok := otherInstance.NativeData.(BaskitMap)
								if !ok {
									return nil, runtime.Exception{Message: fmt.Sprintf("MERGE expects BASKIT argument, got %s", other.Type())}
								}

								// Merge other into this
								for k, v := range otherBaskitMap {
									baskitMap[k] = v
								}
								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "MERGE: invalid context"}
						},
					},
					"COPY": {
						Name: "COPY",
						Documentation: []string{
							"Creates a shallow copy of the BASKIT with all current key-value pairs.",
							"Changes to the copy do not affect the original BASKIT.",
							"",
							"@syntax map DO COPY",
							"@returns {BASKIT} A new BASKIT with the same key-value pairs",
							"@example Create independent copy",
							"I HAS A VARIABLE ORIGINAL TEH BASKIT ITZ NEW BASKIT",
							"ORIGINAL DO PUT WIT \"shared\" AN WIT \"value\"",
							"I HAS A VARIABLE COPY TEH BASKIT ITZ ORIGINAL DO COPY",
							"COPY DO PUT WIT \"new\" AN WIT \"item\"",
							"BTW ORIGINAL doesn't have \"new\" key, COPY does",
							"@example Backup before modification",
							"I HAS A VARIABLE BACKUP TEH BASKIT ITZ ORIGINAL DO COPY",
							"ORIGINAL DO PUT WIT \"temp\" AN WIT \"data\"",
							"BTW Can restore from BACKUP if needed",
							"@note Creates a shallow copy (references to objects are shared)",
							"@note Independent of original for key-value structure",
							"@see MERGE",
							"@category map-transformation",
						},
						ReturnType: "BASKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Create new BASKIT with copied data
								newBaskitObj := NewBaskitInstance()
								newBaskitMap := make(BaskitMap)
								for k, v := range baskitMap {
									newBaskitMap[k] = v
								}
								newBaskitObj.NativeData = newBaskitMap
								return newBaskitObj, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "COPY: invalid context"}
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"SIZ": {
						Variable: environment.Variable{
							Name: "SIZ",
							Documentation: []string{
								"The number of key-value pairs currently stored in the BASKIT.",
								"Read-only property that updates automatically.",
								"",
								"@property {INTEGR} SIZ",
								"@readonly",
								"@example Check map size",
								"I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT",
								"SAYZ WIT MAP SIZ",
								"BTW Output: 0",
								"MAP DO PUT WIT \"key1\" AN WIT \"value1\"",
								"MAP DO PUT WIT \"key2\" AN WIT \"value2\"",
								"SAYZ WIT MAP SIZ",
								"BTW Output: 2",
								"@example Empty after clear",
								"MAP DO CLEAR",
								"SAYZ WIT MAP SIZ",
								"BTW Output: 0",
								"@note Always reflects current number of key-value pairs",
								"@note Cannot be modified directly",
								"@see PUT, REMOVE, CLEAR",
								"@category map-inspection",
							},
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								return environment.IntegerValue(len(baskitMap)), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "SIZ: invalid context"}
						},
						NativeSet: nil, // Read-only
					},
				},
				QualifiedName:    "stdlib:MAPS.BASKIT",
				ModulePath:       "stdlib:MAPS",
				ParentClasses:    []string{},
				MRO:              []string{"stdlib:MAPS.BASKIT"},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
		}
	})
	return mapClasses
}

// RegisterMapsInEnv registers BASKIT class in the given environment
func RegisterMapsInEnv(env *environment.Environment, _ ...string) error {
	// Register the BASKIT class
	for _, class := range getMapClasses() {
		env.DefineClass(class)
	}
	return nil
}
