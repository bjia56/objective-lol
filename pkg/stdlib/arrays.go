package stdlib

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleArraysCategories defines the order that categories should be rendered in documentation
var moduleArraysCategories = []string{
	"array-creation",
	"array-access",
	"array-modification",
	"array-inspection",
	"array-transformation",
}

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
					"Provides methods for adding, removing, accessing, and manipulating elements.",
					"",
					"@class BUKKIT",
					"@example Create empty array",
					"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT",
					"BTW Creates an empty BUKKIT",
					"@example Create array with initial values",
					"I HAS A VARIABLE NUMS TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
					"BTW Creates BUKKIT with [1, 2, 3]",
					"@example Mixed type array",
					"I HAS A VARIABLE MIXED TEH BUKKIT ITZ NEW BUKKIT WIT \"hello\" AN WIT 42 AN WIT YEZ",
					"BTW Creates BUKKIT with [\"hello\", 42, YEZ]",
					"@note Can store any combination of types (INTEGR, DUBBLE, STRIN, BOOL, etc.)",
					"@note Dynamic size - grows and shrinks as needed",
					"@see BASKIT",
					"@category array-creation",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BUKKIT": {
						Name: "BUKKIT",
						Documentation: []string{
							"Initializes a BUKKIT array with an optional list of initial elements.",
							"Creates a new dynamic array that can grow and shrink as needed.",
							"",
							"@syntax NEW BUKKIT [WIT element1 AN WIT element2 ...]",
							"@param {...ANY} elements - Optional initial elements of any type",
							"@returns {BUKKIT} A new BUKKIT instance",
							"@example Create empty array",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT",
							"BTW ARR is now empty []",
							"@example Create with initial values",
							"I HAS A VARIABLE NUMS TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
							"BTW NUMS is now [1, 2, 3]",
							"@note Accepts variable number of arguments",
							"@see PUSH, SET",
							"@category array-creation",
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
							"",
							"@syntax array DO AT WIT <index>",
							"@param {INTEGR} index - Zero-based index of element to retrieve",
							"@returns {ANY} The element at the specified index",
							"@example Access elements",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT \"a\" AN WIT \"b\" AN WIT \"c\"",
							"I HAS A VARIABLE FIRST TEH STRIN ITZ ARR DO AT WIT 0",
							"I HAS A VARIABLE SECOND TEH STRIN ITZ ARR DO AT WIT 1",
							"BTW FIRST = \"a\", SECOND = \"b\"",
							"@throws Index out of bounds exception if index < 0 or index >= size",
							"@note Uses 0-based indexing",
							"@see SET, SIZ",
							"@category array-access",
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
							"",
							"@syntax array DO SET WIT <index> AN WIT <value>",
							"@param {INTEGR} index - Zero-based index of element to modify",
							"@param {ANY} value - New value to assign",
							"@returns {NOTHIN} No return value",
							"@example Modify elements",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
							"ARR DO SET WIT 1 AN WIT 99",
							"BTW ARR is now [1, 99, 3]",
							"@example Change types",
							"ARR DO SET WIT 0 AN WIT \"hello\"",
							"BTW ARR is now [\"hello\", 99, 3]",
							"@throws Index out of bounds exception if index < 0 or index >= size",
							"@note Can change element type",
							"@see AT, SIZ",
							"@category array-modification",
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
							"",
							"@syntax array DO PUSH WIT <element1> [AN WIT <element2> ...]",
							"@param {...ANY} elements - One or more elements to add",
							"@returns {INTEGR} The new size of the array",
							"@example Add single element",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2",
							"I HAS A VARIABLE NEW_SIZE TEH INTEGR ITZ ARR DO PUSH WIT 3",
							"BTW ARR is now [1, 2, 3], NEW_SIZE = 3",
							"@example Add multiple elements",
							"ARR DO PUSH WIT 4 AN WIT 5 AN WIT 6",
							"BTW ARR is now [1, 2, 3, 4, 5, 6]",
							"@note Accepts variable number of arguments",
							"@see POP, UNSHIFT",
							"@category array-modification",
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
							"",
							"@syntax array DO POP",
							"@returns {ANY} The removed last element",
							"@example Remove last element",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
							"I HAS A VARIABLE LAST TEH INTEGR ITZ ARR DO POP",
							"BTW LAST = 3, ARR is now [1, 2]",
							"@example Stack behavior",
							"ARR DO PUSH WIT \"top\"",
							"I HAS A VARIABLE POPPED TEH STRIN ITZ ARR DO POP",
							"BTW POPPED = \"top\", ARR is back to [1, 2]",
							"@throws Exception if array is empty",
							"@note Modifies the original array",
							"@see PUSH, SHIFT",
							"@category array-modification",
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
							"",
							"@syntax array DO SHIFT",
							"@returns {ANY} The removed first element",
							"@example Remove first element",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
							"I HAS A VARIABLE FIRST TEH INTEGR ITZ ARR DO SHIFT",
							"BTW FIRST = 1, ARR is now [2, 3]",
							"@example Queue behavior",
							"ARR DO PUSH WIT 4",
							"I HAS A VARIABLE NEXT TEH INTEGR ITZ ARR DO SHIFT",
							"BTW NEXT = 2, ARR is now [3, 4]",
							"@throws Exception if array is empty",
							"@note Modifies the original array",
							"@see UNSHIFT, POP",
							"@category array-modification",
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
							"",
							"@syntax array DO UNSHIFT WIT <element>",
							"@param {ANY} element - Element to add at the beginning",
							"@returns {INTEGR} The new size of the array",
							"@example Add to beginning",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 2 AN WIT 3",
							"I HAS A VARIABLE NEW_SIZE TEH INTEGR ITZ ARR DO UNSHIFT WIT 1",
							"BTW NEW_SIZE = 3, ARR is now [1, 2, 3]",
							"@example Queue behavior",
							"ARR DO UNSHIFT WIT \"first\"",
							"BTW ARR is now [\"first\", 1, 2, 3]",
							"@note Shifts all existing elements to higher indices",
							"@see SHIFT, PUSH",
							"@category array-modification",
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
							"Resets the array to have zero elements.",
							"",
							"@syntax array DO CLEAR",
							"@returns {NOTHIN} No return value",
							"@example Clear array",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
							"SAYZ WIT ARR SIZ",
							"BTW Output: 3",
							"ARR DO CLEAR",
							"SAYZ WIT ARR SIZ",
							"BTW Output: 0",
							"@note Removes all elements but keeps the array object",
							"@note More efficient than creating a new array",
							"@see SIZ",
							"@category array-modification",
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
							"",
							"@syntax array DO REVERSE",
							"@returns {BUKKIT} The same array (for method chaining)",
							"@example Reverse array",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
							"ARR DO REVERSE",
							"BTW ARR is now [3, 2, 1]",
							"@example Method chaining",
							"I HAS A VARIABLE RESULT TEH BUKKIT ITZ ARR DO REVERSE DO SORT",
							"BTW RESULT is [1, 2, 3] (reversed then sorted)",
							"@note Modifies the original array",
							"@note Returns self for method chaining",
							"@see SORT",
							"@category array-transformation",
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
							"",
							"@syntax array DO SORT",
							"@returns {BUKKIT} The same array (for method chaining)",
							"@example Sort numbers",
							"I HAS A VARIABLE NUMS TEH BUKKIT ITZ NEW BUKKIT WIT 3 AN WIT 1 AN WIT 2",
							"NUMS DO SORT",
							"BTW NUMS is now [1, 2, 3]",
							"@example Sort strings",
							"I HAS A VARIABLE WORDS TEH BUKKIT ITZ NEW BUKKIT WIT \"banana\" AN WIT \"apple\" AN WIT \"cherry\"",
							"WORDS DO SORT",
							"BTW WORDS is now [\"apple\", \"banana\", \"cherry\"]",
							"@example Mixed types",
							"I HAS A VARIABLE MIXED TEH BUKKIT ITZ NEW BUKKIT WIT 2 AN WIT \"1\" AN WIT 3",
							"MIXED DO SORT",
							"BTW Sorts by string representation when types differ",
							"@note Modifies the original array",
							"@note Numbers sort numerically, strings alphabetically",
							"@note Mixed types fall back to string comparison",
							"@see REVERSE",
							"@category array-transformation",
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
							"",
							"@syntax array DO JOIN WIT <separator>",
							"@param {STRIN} separator - String to place between elements",
							"@returns {STRIN} All elements joined into a single string",
							"@example Join with comma",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
							"I HAS A VARIABLE RESULT TEH STRIN ITZ ARR DO JOIN WIT \", \"",
							"BTW RESULT = \"1, 2, 3\"",
							"@example Join words",
							"I HAS A VARIABLE WORDS TEH BUKKIT ITZ NEW BUKKIT WIT \"hello\" AN WIT \"world\"",
							"I HAS A VARIABLE SENTENCE TEH STRIN ITZ WORDS DO JOIN WIT \" \"",
							"BTW SENTENCE = \"hello world\"",
							"@example Empty array",
							"I HAS A VARIABLE EMPTY TEH BUKKIT ITZ NEW BUKKIT",
							"I HAS A VARIABLE EMPTY_STR TEH STRIN ITZ EMPTY DO JOIN WIT \",\"",
							"BTW EMPTY_STR = \"\"",
							"@note Converts all elements to strings",
							"@see SPLIT (in STRING module)",
							"@category array-transformation",
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
							"",
							"@syntax array DO SLICE WIT <start> AN WIT <end>",
							"@param {INTEGR} start - Starting index (inclusive)",
							"@param {INTEGR} end - Ending index (exclusive)",
							"@returns {BUKKIT} New array containing the sliced elements",
							"@example Basic slicing",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 0 AN WIT 1 AN WIT 2 AN WIT 3 AN WIT 4",
							"I HAS A VARIABLE SUB TEH BUKKIT ITZ ARR DO SLICE WIT 1 AN WIT 3",
							"BTW SUB = [1, 2] (indices 1 and 2, excluding 3)",
							"@example Negative indices",
							"I HAS A VARIABLE LAST_TWO TEH BUKKIT ITZ ARR DO SLICE WIT -2 AN WIT -0",
							"BTW LAST_TWO = [3, 4] (last two elements)",
							"@example Copy array",
							"I HAS A VARIABLE COPY TEH BUKKIT ITZ ARR DO SLICE WIT 0 AN WIT ARR SIZ",
							"BTW COPY is a shallow copy of ARR",
							"@throws Index out of bounds exception for invalid indices",
							"@note Creates a new array, doesn't modify original",
							"@note Negative indices count from end (-1 = last element)",
							"@see AT, SIZ",
							"@category array-transformation",
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
							"",
							"@syntax array DO FIND WIT <value>",
							"@param {ANY} value - Value to search for",
							"@returns {INTEGR} Index of first occurrence, or -1 if not found",
							"@example Find number",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 10 AN WIT 20 AN WIT 30",
							"I HAS A VARIABLE INDEX TEH INTEGR ITZ ARR DO FIND WIT 20",
							"BTW INDEX = 1",
							"@example Find string",
							"I HAS A VARIABLE WORDS TEH BUKKIT ITZ NEW BUKKIT WIT \"cat\" AN WIT \"dog\" AN WIT \"cat\"",
							"I HAS A VARIABLE FIRST_CAT TEH INTEGR ITZ WORDS DO FIND WIT \"cat\"",
							"BTW FIRST_CAT = 0 (finds first occurrence)",
							"@example Not found",
							"I HAS A VARIABLE NOT_FOUND TEH INTEGR ITZ ARR DO FIND WIT 999",
							"BTW NOT_FOUND = -1",
							"@note Uses equality comparison between values",
							"@note Returns index of first match only",
							"@see CONTAINS, AT",
							"@category array-inspection",
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
							"",
							"@syntax array DO CONTAINS WIT <value>",
							"@param {ANY} value - Value to search for",
							"@returns {BOOL} YEZ if value is found, NO otherwise",
							"@example Check for number",
							"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
							"I HAS A VARIABLE HAS_TWO TEH BOOL ITZ ARR DO CONTAINS WIT 2",
							"BTW HAS_TWO = YEZ",
							"@example Check for string",
							"I HAS A VARIABLE PETS TEH BUKKIT ITZ NEW BUKKIT WIT \"cat\" AN WIT \"dog\"",
							"I HAS A VARIABLE HAS_BIRD TEH BOOL ITZ PETS DO CONTAINS WIT \"bird\"",
							"BTW HAS_BIRD = NO",
							"@example Use in conditional",
							"IZ ARR DO CONTAINS WIT 5?",
							"    SAYZ WIT \"Found 5!\"",
							"NOPE",
							"    SAYZ WIT \"5 not found\"",
							"KTHX",
							"@note More convenient than FIND when you only need to know if value exists",
							"@see FIND",
							"@category array-inspection",
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
								"Automatically updated when elements are added or removed.",
								"",
								"@property {INTEGR} SIZ",
								"@readonly",
								"@example Check array size",
								"I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT WIT 1 AN WIT 2 AN WIT 3",
								"SAYZ WIT ARR SIZ",
								"BTW Output: 3",
								"@example Empty array",
								"I HAS A VARIABLE EMPTY TEH BUKKIT ITZ NEW BUKKIT",
								"SAYZ WIT EMPTY SIZ",
								"BTW Output: 0",
								"@example Dynamic sizing",
								"ARR DO PUSH WIT 4",
								"SAYZ WIT ARR SIZ",
								"BTW Output: 4",
								"ARR DO POP",
								"SAYZ WIT ARR SIZ",
								"BTW Output: 3",
								"@note Always reflects current element count",
								"@note Cannot be modified directly",
								"@see PUSH, POP, CLEAR",
								"@category array-inspection",
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
