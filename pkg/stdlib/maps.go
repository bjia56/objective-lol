package stdlib

import (
	"fmt"
	"sort"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

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
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BASKIT": {
						Name: "BASKIT",
						Documentation: []string{
							"Initializes an empty BASKIT map.",
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
