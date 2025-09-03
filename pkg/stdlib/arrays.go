package stdlib

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
	"github.com/bjia56/objective-lol/pkg/types"
)

type BukkitSlice []types.Value

func NewBukkitInstance() *environment.ObjectInstance {
	class := getArrayClasses()["BUKKIT"]
	return &environment.ObjectInstance{
		Class:      class,
		NativeData: make(BukkitSlice, 0),
		Hierarchy:  []string{"BUKKIT"},
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

// updateSIZ updates the SIZ variable in the object instance based on slice length
func updateSIZ(obj *environment.ObjectInstance, slice BukkitSlice) {
	if sizVar, exists := obj.Variables["SIZ"]; exists {
		sizVar.Value = types.IntegerValue(len(slice))
	}
}

// Global BUKKIT class definition - created once and reused
var arrayClassOnce = sync.Once{}
var arrayClasses map[string]*environment.Class

func getArrayClasses() map[string]*environment.Class {
	arrayClassOnce.Do(func() {
		arrayClasses = map[string]*environment.Class{
			"BUKKIT": {
				Name: "BUKKIT",
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BUKKIT": {
						Name:       "BUKKIT",
						Parameters: []environment.Parameter{}, // Empty constructor - no arguments
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							// Create empty slice and store in NativeData
							slice := make(BukkitSlice, 0)
							this.NativeData = slice
							return types.NOTHIN, nil
						},
					},
					// Array methods
					"AT": {
						Name:       "AT",
						Parameters: []environment.Parameter{{Name: "INDEX", Type: "INTEGR"}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							index := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								indexVal, ok := index.(types.IntegerValue)
								if !ok {
									return nil, fmt.Errorf("AT expects INTEGR index, got %s", index.Type())
								}
								idx := int(indexVal)
								if idx < 0 || idx >= len(slice) {
									return nil, runtime.Exception{Message: fmt.Sprintf("Array index %d out of bounds (size %d)", idx, len(slice))}
								}
								return slice[idx], nil
							}
							return types.NOTHIN, fmt.Errorf("AT: invalid context")
						},
					},
					"SET": {
						Name:       "SET",
						Parameters: []environment.Parameter{{Name: "INDEX", Type: "INTEGR"}, {Name: "VALUE", Type: ""}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							index, value := args[0], args[1]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								indexVal, ok := index.(types.IntegerValue)
								if !ok {
									return nil, fmt.Errorf("SET expects INTEGR index, got %s", index.Type())
								}
								idx := int(indexVal)
								if idx < 0 || idx >= len(slice) {
									return nil, runtime.Exception{Message: fmt.Sprintf("Array index %d out of bounds (size %d)", idx, len(slice))}
								}
								slice[idx] = value
								return types.NOTHIN, nil
							}
							return types.NOTHIN, fmt.Errorf("SET: invalid context")
						},
					},
					"PUSH": {
						Name:       "PUSH",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							value := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								newSlice := append(slice, value)
								this.NativeData = newSlice
								updateSIZ(this, newSlice)
								return types.IntegerValue(len(newSlice)), nil
							}
							return types.NOTHIN, fmt.Errorf("PUSH: invalid context")
						},
					},
					"POP": {
						Name:       "POP",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								if len(slice) == 0 {
									return nil, fmt.Errorf("cannot pop from empty array")
								}
								lastIndex := len(slice) - 1
								element := slice[lastIndex]
								newSlice := slice[:lastIndex]
								this.NativeData = newSlice
								updateSIZ(this, newSlice)
								return element, nil
							}
							return types.NOTHIN, fmt.Errorf("POP: invalid context")
						},
					},
					"SHIFT": {
						Name:       "SHIFT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								if len(slice) == 0 {
									return nil, fmt.Errorf("cannot shift from empty array")
								}
								element := slice[0]
								newSlice := slice[1:]
								this.NativeData = newSlice
								updateSIZ(this, newSlice)
								return element, nil
							}
							return types.NOTHIN, fmt.Errorf("SHIFT: invalid context")
						},
					},
					"UNSHIFT": {
						Name:       "UNSHIFT",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							value := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								newSlice := append(BukkitSlice{value}, slice...)
								this.NativeData = newSlice
								updateSIZ(this, newSlice)
								return types.IntegerValue(len(newSlice)), nil
							}
							return types.NOTHIN, fmt.Errorf("UNSHIFT: invalid context")
						},
					},
					"CLEAR": {
						Name:       "CLEAR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							newSlice := make(BukkitSlice, 0)
							this.NativeData = newSlice
							updateSIZ(this, newSlice)
							return types.NOTHIN, nil
						},
					},
					"REVERSE": {
						Name:       "REVERSE",
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								// Reverse in-place
								for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
									slice[i], slice[j] = slice[j], slice[i]
								}
								return types.NewObjectValue(this, "BUKKIT"), nil
							}
							return types.NOTHIN, fmt.Errorf("REVERSE: invalid context")
						},
					},
					"SORT": {
						Name:       "SORT",
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								sort.Slice(slice, func(i, j int) bool {
									left, right := slice[i], slice[j]

									// Handle different type combinations
									switch leftVal := left.(type) {
									case types.IntegerValue:
										if rightInt, ok := right.(types.IntegerValue); ok {
											return leftVal < rightInt
										} else if rightDouble, ok := right.(types.DoubleValue); ok {
											return float64(leftVal) < float64(rightDouble)
										}
									case types.DoubleValue:
										if rightDouble, ok := right.(types.DoubleValue); ok {
											return leftVal < rightDouble
										} else if rightInt, ok := right.(types.IntegerValue); ok {
											return float64(leftVal) < float64(rightInt)
										}
									case types.StringValue:
										if rightStr, ok := right.(types.StringValue); ok {
											return leftVal < rightStr
										}
									}

									// Default: convert both to strings and compare
									return left.String() < right.String()
								})
								return types.NewObjectValue(this, "BUKKIT"), nil
							}
							return types.NOTHIN, fmt.Errorf("SORT: invalid context")
						},
					},
					"JOIN": {
						Name:       "JOIN",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{{Name: "SEPARATOR", Type: "STRIN"}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							separator := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								separatorVal, ok := separator.(types.StringValue)
								if !ok {
									return nil, fmt.Errorf("JOIN expects STRIN separator, got %s", args[0].Type())
								}

								if len(slice) == 0 {
									return types.StringValue(""), nil
								}

								var parts []string
								for _, elem := range slice {
									parts = append(parts, elem.String())
								}
								return types.StringValue(strings.Join(parts, string(separatorVal))), nil
							}
							return types.NOTHIN, fmt.Errorf("JOIN: invalid context")
						},
					},
					"SLICE": {
						Name:       "SLICE",
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{{Name: "START", Type: "INTEGR"}, {Name: "END", Type: "INTEGR"}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							start, end := args[0], args[1]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								startVal, ok := start.(types.IntegerValue)
								if !ok {
									return nil, fmt.Errorf("SLICE expects INTEGR start, got %s", start.Type())
								}
								endVal, ok := end.(types.IntegerValue)
								if !ok {
									return nil, fmt.Errorf("SLICE expects INTEGR end, got %s", end.Type())
								}

								startIdx, endIdx := int(startVal), int(endVal)
								size := len(slice)

								// Handle negative indices
								if startIdx < 0 {
									startIdx = size + startIdx
								}
								if endIdx < 0 {
									endIdx = size + endIdx
								}

								// Bounds checking
								if startIdx < 0 || startIdx > size || endIdx < 0 || endIdx > size || startIdx > endIdx {
									return nil, runtime.Exception{Message: fmt.Sprintf("Slice indices out of bounds: start=%d, end=%d, size=%d", startIdx, endIdx, size)}
								}

								newSlice := make(BukkitSlice, endIdx-startIdx)
								copy(newSlice, slice[startIdx:endIdx])

								// Create a new BUKKIT object with the sliced array
								newObject := NewBukkitInstance()
								newObject.NativeData = newSlice
								updateSIZ(newObject, newSlice)

								return types.NewObjectValue(newObject, "BUKKIT"), nil
							}
							return types.NOTHIN, fmt.Errorf("SLICE: invalid context")
						},
					},
					"FIND": {
						Name:       "FIND",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							value := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								for i, elem := range slice {
									equal, err := elem.EqualTo(value)
									if err == nil && equal {
										return types.IntegerValue(i), nil
									}
								}
								return types.IntegerValue(-1), nil
							}
							return types.NOTHIN, fmt.Errorf("FIND: invalid context")
						},
					},
					"CONTAINS": {
						Name:       "CONTAINS",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							value := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								for _, elem := range slice {
									equal, err := elem.EqualTo(value)
									if err == nil && equal {
										return types.YEZ, nil
									}
								}
								return types.NO, nil
							}
							return types.NOTHIN, fmt.Errorf("CONTAINS: invalid context")
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
	return arrayClasses
}

// RegisterArrays registers BUKKIT class in the given environment
func RegisterArraysInEnv(env *environment.Environment, _ ...string) error {
	// Register the BUKKIT class
	for _, class := range getArrayClasses() {
		env.DefineClass(class)
	}
	return nil
}
