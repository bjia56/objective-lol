package stdlib

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

type BukkitSlice []environment.Value

func NewBukkitInstance() *environment.ObjectInstance {
	class := getArrayClasses()["BUKKIT"]
	env := environment.NewEnvironment(nil)
	env.DefineClass(class)
	obj := &environment.ObjectInstance{
		Environment: env,
		Class:       class,
		NativeData:  make(BukkitSlice, 0),
		Variables:   make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(obj)
	return obj
}

// Global BUKKIT class definition - created once and reused
var arrayClassOnce = sync.Once{}
var arrayClasses map[string]*environment.Class

func getArrayClasses() map[string]*environment.Class {
	arrayClassOnce.Do(func() {
		arrayClasses = map[string]*environment.Class{
			"BUKKIT": {
				Name: "BUKKIT",
				Documentation: []string{
					"A dynamic array that can hold any combination of values and types.",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BUKKIT": {
						Name: "BUKKIT",
						Documentation: []string{
							"Initializes a BUKKIT array with an optional list of initial elements.",
						},
						IsVarargs: true,
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							// Create slice and store in NativeData
							slice := make(BukkitSlice, 0, len(args))
							slice = append(slice, args...)
							this.NativeData = slice
							return environment.NOTHIN, nil
						},
					},
					// Array methods
					"AT": {
						Name: "AT",
						Documentation: []string{
							"Gets the element at the specified index (0-based).",
							"Throws an exception if the index is out of bounds.",
						},
						Parameters: []environment.Parameter{{Name: "INDEX", Type: "INTEGR"}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							index := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								indexVal, ok := index.(environment.IntegerValue)
								if !ok {
									return nil, runtime.Exception{Message: fmt.Sprintf("AT expects INTEGR index, got %s", index.Type())}
								}
								idx := int(indexVal)
								if idx < 0 || idx >= len(slice) {
									return nil, runtime.Exception{Message: fmt.Sprintf("BUKKIT index %d out of bounds (size %d)", idx, len(slice))}
								}
								return slice[idx], nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "AT: invalid context"}
						},
					},
					"SET": {
						Name: "SET",
						Documentation: []string{
							"Sets the element at the specified index to the given value.",
							"Throws an exception if the index is out of bounds.",
						},
						Parameters: []environment.Parameter{{Name: "INDEX", Type: "INTEGR"}, {Name: "VALUE", Type: ""}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							index, value := args[0], args[1]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								indexVal, ok := index.(environment.IntegerValue)
								if !ok {
									return nil, runtime.Exception{Message: fmt.Sprintf("SET expects INTEGR index, got %s", index.Type())}
								}
								idx := int(indexVal)
								if idx < 0 || idx >= len(slice) {
									return nil, runtime.Exception{Message: fmt.Sprintf("BUKKIT index %d out of bounds (size %d)", idx, len(slice))}
								}
								slice[idx] = value
								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "SET: invalid context"}
						},
					},
					"PUSH": {
						Name: "PUSH",
						Documentation: []string{
							"Adds elements to the end of the BUKKIT.",
							"Returns the BUKKIT's new size.",
						},
						ReturnType: "INTEGR",
						IsVarargs:  true,
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								newSlice := append(slice, args...)
								this.NativeData = newSlice
								return environment.IntegerValue(len(newSlice)), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "PUSH: invalid context"}
						},
					},
					"POP": {
						Name: "POP",
						Documentation: []string{
							"Removes and returns the last element from the BUKKIT.",
							"Throws an exception if the BUKKIT is empty.",
						},
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								if len(slice) == 0 {
									return nil, runtime.Exception{Message: "cannot pop from empty BUKKIT"}
								}
								lastIndex := len(slice) - 1
								element := slice[lastIndex]
								newSlice := slice[:lastIndex]
								this.NativeData = newSlice
								return element, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "POP: invalid context"}
						},
					},
					"SHIFT": {
						Name: "SHIFT",
						Documentation: []string{
							"Removes and returns the first element from the BUKKIT.",
							"Throws an exception if the BUKKIT is empty.",
						},
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								if len(slice) == 0 {
									return nil, runtime.Exception{Message: "cannot shift from empty BUKKIT"}
								}
								element := slice[0]
								newSlice := slice[1:]
								this.NativeData = newSlice
								return element, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "SHIFT: invalid context"}
						},
					},
					"UNSHIFT": {
						Name: "UNSHIFT",
						Documentation: []string{
							"Adds an element to the beginning of the BUKKIT.",
							"Returns the new size of the BUKKIT.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							value := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								newSlice := append(BukkitSlice{value}, slice...)
								this.NativeData = newSlice
								return environment.IntegerValue(len(newSlice)), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "UNSHIFT: invalid context"}
						},
					},
					"CLEAR": {
						Name: "CLEAR",
						Documentation: []string{
							"Removes all elements from the BUKKIT, making it empty.",
						},
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							newSlice := make(BukkitSlice, 0)
							this.NativeData = newSlice
							return environment.NOTHIN, nil
						},
					},
					"REVERSE": {
						Name: "REVERSE",
						Documentation: []string{
							"Reverses the order of elements in the BUKKIT in place.",
							"Returns the BUKKIT itself.",
						},
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								// Reverse in-place
								for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
									slice[i], slice[j] = slice[j], slice[i]
								}
								return this, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "REVERSE: invalid context"}
						},
					},
					"SORT": {
						Name: "SORT",
						Documentation: []string{
							"Sorts the elements in the BUKKIT in place in ascending order.",
							"Handles different type combinations and converts to strings as fallback.",
							"Returns the BUKKIT itself.",
						},
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								sort.Slice(slice, func(i, j int) bool {
									left, right := slice[i], slice[j]

									// Handle different type combinations
									switch leftVal := left.(type) {
									case environment.IntegerValue:
										if rightInt, ok := right.(environment.IntegerValue); ok {
											return leftVal < rightInt
										} else if rightDouble, ok := right.(environment.DoubleValue); ok {
											return float64(leftVal) < float64(rightDouble)
										}
									case environment.DoubleValue:
										if rightDouble, ok := right.(environment.DoubleValue); ok {
											return leftVal < rightDouble
										} else if rightInt, ok := right.(environment.IntegerValue); ok {
											return float64(leftVal) < float64(rightInt)
										}
									case environment.StringValue:
										if rightStr, ok := right.(environment.StringValue); ok {
											return leftVal < rightStr
										}
									}

									// Default: convert both to strings and compare
									return left.String() < right.String()
								})
								return this, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "SORT: invalid context"}
						},
					},
					"JOIN": {
						Name: "JOIN",
						Documentation: []string{
							"Joins all elements in the BUKKIT into a single string using the specified separator.",
							"Returns an empty string if the BUKKIT is empty.",
						},
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{{Name: "SEPARATOR", Type: "STRIN"}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							separator := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								separatorVal, ok := separator.(environment.StringValue)
								if !ok {
									return nil, runtime.Exception{Message: fmt.Sprintf("JOIN expects STRIN separator, got %s", args[0].Type())}
								}

								if len(slice) == 0 {
									return environment.StringValue(""), nil
								}

								var parts []string
								for _, elem := range slice {
									parts = append(parts, elem.String())
								}
								return environment.StringValue(strings.Join(parts, string(separatorVal))), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "JOIN: invalid context"}
						},
					},
					"SLICE": {
						Name: "SLICE",
						Documentation: []string{
							"Creates a new BUKKIT containing elements from START index to END index (exclusive).",
							"Supports negative indices to count from the end.",
							"Throws an exception if indices are out of bounds.",
						},
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{{Name: "START", Type: "INTEGR"}, {Name: "END", Type: "INTEGR"}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							start, end := args[0], args[1]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								startVal, ok := start.(environment.IntegerValue)
								if !ok {
									return nil, runtime.Exception{Message: fmt.Sprintf("SLICE expects INTEGR start, got %s", start.Type())}
								}
								endVal, ok := end.(environment.IntegerValue)
								if !ok {
									return nil, runtime.Exception{Message: fmt.Sprintf("SLICE expects INTEGR end, got %s", end.Type())}
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
									return nil, runtime.Exception{Message: fmt.Sprintf("BUKKIT indices out of bounds: start=%d, end=%d, size=%d", startIdx, endIdx, size)}
								}

								newSlice := make(BukkitSlice, endIdx-startIdx)
								copy(newSlice, slice[startIdx:endIdx])

								// Create a new BUKKIT object with the sliced array
								newObject := NewBukkitInstance()
								newObject.NativeData = newSlice

								return newObject, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "SLICE: invalid context"}
						},
					},
					"FIND": {
						Name: "FIND",
						Documentation: []string{
							"Finds the first index of the specified value in the BUKKIT.",
							"Returns -1 if the value is not found.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							value := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								for i, elem := range slice {
									equal, err := elem.EqualTo(value)
									if err == nil && equal {
										return environment.IntegerValue(i), nil
									}
								}
								return environment.IntegerValue(-1), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "FIND: invalid context"}
						},
					},
					"CONTAINS": {
						Name: "CONTAINS",
						Documentation: []string{
							"Checks if the BUKKIT contains the specified value.",
							"Returns YEZ if found, NO otherwise.",
						},
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							value := args[0]

							if slice, ok := this.NativeData.(BukkitSlice); ok {
								for _, elem := range slice {
									equal, err := elem.EqualTo(value)
									if err == nil && equal {
										return environment.YEZ, nil
									}
								}
								return environment.NO, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "CONTAINS: invalid context"}
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"SIZ": {
						Variable: environment.Variable{
							Name:     "SIZ",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property that returns the current number of elements in the BUKKIT.",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if slice, ok := this.NativeData.(BukkitSlice); ok {
								return environment.IntegerValue(len(slice)), nil
							}
							return environment.IntegerValue(0), runtime.Exception{Message: "invalid context for SIZ"}
						},
						NativeSet: nil, // Read-only
					},
				},
				QualifiedName:    "stdlib:ARRAYS.BUKKIT",
				ModulePath:       "stdlib:ARRAYS",
				ParentClasses:    []string{},
				MRO:              []string{"stdlib:ARRAYS.BUKKIT"},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
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
