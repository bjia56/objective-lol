package stdlib

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// Global RANDOM function definitions - created once and reused
var randomFuncOnce = sync.Once{}
var randomFunctions map[string]*environment.Function

func getRandomFunctions() map[string]*environment.Function {
	randomFuncOnce.Do(func() {
		randomFunctions = map[string]*environment.Function{
			"SEED": {
				Name: "SEED",
				Documentation: []string{
					"Sets the random number generator seed for reproducible results.",
					"Using the same seed will produce the same sequence of random numbers.",
				},
				Parameters: []environment.Parameter{
					{Name: "seed", Type: "INTEGR"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					seed := args[0]

					if seedVal, ok := seed.(environment.IntegerValue); ok {
						rand.Seed(int64(seedVal))
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
				},
				Parameters: []environment.Parameter{},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					rand.Seed(time.Now().UnixNano())
					return environment.NOTHIN, nil
				},
			},
			"RANDOM_FLOAT": {
				Name: "RANDOM_FLOAT",
				Documentation: []string{
					"Returns a random floating-point number between 0.0 (inclusive) and 1.0 (exclusive).",
					"Useful for probability calculations and random selection.",
				},
				ReturnType: "DUBBLE",
				Parameters: []environment.Parameter{},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					return environment.DoubleValue(rand.Float64()), nil
				},
			},
			"RANDOM_RANGE": {
				Name: "RANDOM_RANGE",
				Documentation: []string{
					"Returns a random floating-point number within the specified range.",
					"Range is [min, max) - includes min but excludes max. Min must be less than max.",
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
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("RANDOM_RANGE: min must be less than max")}
							}
							result := float64(minVal) + rand.Float64()*(float64(maxVal)-float64(minVal))
							return environment.DoubleValue(result), nil
						}
					}

					return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("RANDOM_RANGE: invalid numeric arguments")}
				},
			},
			"RANDOM_INT": {
				Name: "RANDOM_INT",
				Documentation: []string{
					"Returns a random integer within the specified range.",
					"Range is [min, max) - includes min but excludes max. Min must be less than max.",
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
							result := rand.Int63n(int64(maxVal-minVal)) + int64(minVal)
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
				},
				ReturnType: "BOOL",
				Parameters: []environment.Parameter{},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					if rand.Float64() < 0.5 {
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
							index := rand.Intn(len(slice))
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
								j := rand.Intn(i + 1)
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
								result[i] = charsetStr[rand.Intn(len(charsetStr))]
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
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					// Simple UUID v4 implementation
					uuid := make([]byte, 16)
					rand.Read(uuid)

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
