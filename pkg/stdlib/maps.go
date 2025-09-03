package stdlib

import (
	"fmt"
	"sort"
	"sync"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

type BaskitMap map[string]types.Value

func NewBaskitInstance() *environment.ObjectInstance {
	class := getMapClasses()["BASKIT"]
	return &environment.ObjectInstance{
		Class:      class,
		NativeData: make(BaskitMap),
		Hierarchy:  []string{"BASKIT"},
		Variables: map[string]*environment.Variable{
			"SIZ": {
				Name:     "SIZ",
				Type:     "INTEGR",
				Value:    types.IntegerValue(0),
				IsLocked: true,
				IsPublic: true,
			},
		},
	}
}

// updateSIZ updates the SIZ variable in the object instance based on map length
func updateBaskitSIZ(obj *environment.ObjectInstance, baskitMap BaskitMap) {
	if sizVar, exists := obj.Variables["SIZ"]; exists {
		sizVar.Value = types.IntegerValue(len(baskitMap))
	}
}

// Global BASKIT class definition - created once and reused
var mapClassesOnce = sync.Once{}
var mapClasses map[string]*environment.Class

func getMapClasses() map[string]*environment.Class {
	mapClassesOnce.Do(func() {
		mapClasses = map[string]*environment.Class{
			"BASKIT": {
				Name: "BASKIT",
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BASKIT": {
						Name:       "BASKIT",
						Parameters: []environment.Parameter{}, // Empty constructor - no arguments
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							// Create empty map and store in NativeData
							baskitMap := make(BaskitMap)
							this.NativeData = baskitMap
							return types.NOTHIN, nil
						},
					},
					// Map methods
					"PUT": {
						Name:       "PUT",
						Parameters: []environment.Parameter{{Name: "KEY", Type: ""}, {Name: "VALUE", Type: ""}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							key, value := args[0], args[1]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Convert key to string
								keyStr := key.String()
								baskitMap[keyStr] = value
								updateBaskitSIZ(this, baskitMap)
								return types.NOTHIN, nil
							}
							return types.NOTHIN, fmt.Errorf("PUT: invalid context")
						},
					},
					"GET": {
						Name:       "GET",
						Parameters: []environment.Parameter{{Name: "KEY", Type: ""}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							key := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Convert key to string
								keyStr := key.String()
								if value, exists := baskitMap[keyStr]; exists {
									return value, nil
								}
								return nil, ast.Exception{Message: fmt.Sprintf("Key '%s' not found in BASKIT", keyStr)}
							}
							return types.NOTHIN, fmt.Errorf("GET: invalid context")
						},
					},
					"CONTAINS": {
						Name:       "CONTAINS",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{{Name: "KEY", Type: ""}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							key := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Convert key to string
								keyStr := key.String()
								_, exists := baskitMap[keyStr]
								if exists {
									return types.YEZ, nil
								}
								return types.NO, nil
							}
							return types.NOTHIN, fmt.Errorf("CONTAINS: invalid context")
						},
					},
					"REMOVE": {
						Name:       "REMOVE",
						Parameters: []environment.Parameter{{Name: "KEY", Type: ""}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							key := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Convert key to string
								keyStr := key.String()
								if value, exists := baskitMap[keyStr]; exists {
									delete(baskitMap, keyStr)
									updateBaskitSIZ(this, baskitMap)
									return value, nil
								}
								return nil, ast.Exception{Message: fmt.Sprintf("Key '%s' not found in BASKIT", keyStr)}
							}
							return types.NOTHIN, fmt.Errorf("REMOVE: invalid context")
						},
					},
					"CLEAR": {
						Name:       "CLEAR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							newMap := make(BaskitMap)
							this.NativeData = newMap
							updateBaskitSIZ(this, newMap)
							return types.NOTHIN, nil
						},
					},
					"KEYS": {
						Name:       "KEYS",
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
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
									bukkitSlice = append(bukkitSlice, types.StringValue(k))
								}
								bukkitObj.NativeData = bukkitSlice
								updateSIZ(bukkitObj, bukkitSlice)
								return types.NewObjectValue(bukkitObj, "BUKKIT"), nil
							}
							return types.NOTHIN, fmt.Errorf("KEYS: invalid context")
						},
					},
					"VALUES": {
						Name:       "VALUES",
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
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
								updateSIZ(bukkitObj, bukkitSlice)
								return types.NewObjectValue(bukkitObj, "BUKKIT"), nil
							}
							return types.NOTHIN, fmt.Errorf("VALUES: invalid context")
						},
					},
					"PAIRS": {
						Name:       "PAIRS",
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
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
									pairSlice := BukkitSlice{types.StringValue(k), baskitMap[k]}
									pairObj.NativeData = pairSlice
									updateSIZ(pairObj, pairSlice)
									bukkitSlice = append(bukkitSlice, types.NewObjectValue(pairObj, "BUKKIT"))
								}
								bukkitObj.NativeData = bukkitSlice
								updateSIZ(bukkitObj, bukkitSlice)
								return types.NewObjectValue(bukkitObj, "BUKKIT"), nil
							}
							return types.NOTHIN, fmt.Errorf("PAIRS: invalid context")
						},
					},
					"MERGE": {
						Name:       "MERGE",
						Parameters: []environment.Parameter{{Name: "OTHER", Type: "BASKIT"}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							other := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Extract other BASKIT
								otherObj, ok := other.(types.ObjectValue)
								if !ok {
									return nil, fmt.Errorf("MERGE expects BASKIT argument, got %s", other.Type())
								}
								otherInstance, ok := otherObj.Instance.(*environment.ObjectInstance)
								if !ok {
									return nil, fmt.Errorf("MERGE expects BASKIT argument")
								}
								otherBaskitMap, ok := otherInstance.NativeData.(BaskitMap)
								if !ok {
									return nil, fmt.Errorf("MERGE expects BASKIT argument")
								}

								// Merge other into this
								for k, v := range otherBaskitMap {
									baskitMap[k] = v
								}
								updateBaskitSIZ(this, baskitMap)
								return types.NOTHIN, nil
							}
							return types.NOTHIN, fmt.Errorf("MERGE: invalid context")
						},
					},
					"COPY": {
						Name:       "COPY",
						ReturnType: "BASKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Create new BASKIT with copied data
								newBaskitObj := NewBaskitInstance()
								newBaskitMap := make(BaskitMap)
								for k, v := range baskitMap {
									newBaskitMap[k] = v
								}
								newBaskitObj.NativeData = newBaskitMap
								updateBaskitSIZ(newBaskitObj, newBaskitMap)
								return types.NewObjectValue(newBaskitObj, "BASKIT"), nil
							}
							return types.NOTHIN, fmt.Errorf("COPY: invalid context")
						},
					},
				},
				PublicVariables: map[string]*environment.Variable{
					"SIZ": {
						Name:     "SIZ",
						Type:     "INTEGR",
						Value:    types.IntegerValue(0),
						IsLocked: true,
						IsPublic: true,
					},
				},
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
