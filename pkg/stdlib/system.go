package stdlib

import (
	"fmt"
	"os"
	goRuntime "runtime"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// NewEnvbaskitInstance creates a new ENVBASKIT instance (internal use only)
func NewEnvbaskitInstance() *environment.ObjectInstance {
	class := getSystemClasses()["ENVBASKIT"]
	env := environment.NewEnvironment(nil)
	env.DefineClass(class)

	// Also define the parent BASKIT class
	parentClass := getMapClasses()["BASKIT"]
	env.DefineClass(parentClass)

	obj := &environment.ObjectInstance{
		Environment: env,
		Class:       class,
		NativeData:  make(BaskitMap), // Use same native data type as BASKIT
		Variables:   make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(obj)
	return obj
}

// Global SYSTEM classes definition - created once and reused
var systemClassesOnce = sync.Once{}
var systemClasses map[string]*environment.Class

func getSystemClasses() map[string]*environment.Class {
	systemClassesOnce.Do(func() {
		// Get the parent BASKIT class
		baskitClass := getMapClasses()["BASKIT"]

		systemClasses = map[string]*environment.Class{
			"ENVBASKIT": {
				Name:          "ENVBASKIT",
				ParentClasses: []string{"stdlib:MAPS.BASKIT"}, // Inherit from BASKIT
				Documentation: []string{
					"Special BASKIT type to provide integration with system environment variables.",
					"Automatically syncs with the actual process environment.",
				},
				PublicFunctions: map[string]*environment.Function{
					// Override PUT to also set environment variable
					"PUT": {
						Name:       "PUT",
						Parameters: []environment.Parameter{{Name: "KEY", Type: "STRIN"}, {Name: "VALUE", Type: ""}},
						Documentation: []string{
							"Sets an environment variable both in the internal map and the actual process environment.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							key, value := args[0], args[1]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Convert key and value to strings
								keyStr := key.String()
								valueStr := value.String()

								// Set in BASKIT
								baskitMap[keyStr] = environment.StringValue(valueStr)

								// Also set as actual environment variable
								if err := os.Setenv(keyStr, valueStr); err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("PUT: failed to set environment variable %s: %v", keyStr, err)}
								}

								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "PUT: invalid context"}
						},
					},
					// Override GET to also check environment variables
					"GET": {
						Name:       "GET",
						Parameters: []environment.Parameter{{Name: "KEY", Type: "STRIN"}},
						Documentation: []string{
							"Gets the value of an environment variable.",
							"Checks the internal map first, then the actual environment, throws exception if not found.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							key := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								keyStr := key.String()

								// First try to get from BASKIT
								if value, exists := baskitMap[keyStr]; exists {
									return value, nil
								}

								// If not found in BASKIT, try to get from actual environment
								if envValue, exists := os.LookupEnv(keyStr); exists {
									// Update BASKIT with the environment value
									baskitMap[keyStr] = environment.StringValue(envValue)
									return environment.StringValue(envValue), nil
								}

								return nil, runtime.Exception{Message: fmt.Sprintf("GET: key %q not found in environment", keyStr)}
							}
							return environment.NOTHIN, runtime.Exception{Message: "GET: invalid context"}
						},
					},
					// Override CONTAINS to also check environment variables
					"CONTAINS": {
						Name:       "CONTAINS",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{{Name: "KEY", Type: "STRIN"}},
						Documentation: []string{
							"Checks if an environment variable exists in either the internal map or actual environment.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							key := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								keyStr := key.String()

								// Check BASKIT first
								if _, exists := baskitMap[keyStr]; exists {
									return environment.YEZ, nil
								}

								// Check actual environment
								if _, exists := os.LookupEnv(keyStr); exists {
									// Update BASKIT with the environment value
									if envValue, ok := os.LookupEnv(keyStr); ok {
										baskitMap[keyStr] = environment.StringValue(envValue)
									}
									return environment.YEZ, nil
								}

								return environment.NO, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "CONTAINS: invalid context"}
						},
					},
					// Override REMOVE to also unset environment variable
					"REMOVE": {
						Name:       "REMOVE",
						Parameters: []environment.Parameter{{Name: "KEY", Type: "STRIN"}},
						Documentation: []string{
							"Removes an environment variable from both the internal map and actual process environment.",
							"Returns the previous value, throws exception if not found.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							key := args[0]

							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								keyStr := key.String()

								// Get current value before removing
								var value environment.Value = environment.NOTHIN
								if val, exists := baskitMap[keyStr]; exists {
									value = val
								} else if envValue, exists := os.LookupEnv(keyStr); exists {
									value = environment.StringValue(envValue)
								} else {
									return nil, runtime.Exception{Message: fmt.Sprintf("REMOVE: key %q not found in environment", keyStr)}
								}

								// Remove from BASKIT
								delete(baskitMap, keyStr)

								// Also unset environment variable
								if err := os.Unsetenv(keyStr); err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("REMOVE: failed to unset environment variable %s: %v", keyStr, err)}
								}

								return value, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "REMOVE: invalid context"}
						},
					},
					// Override CLEAR to also clear environment variables managed by this instance
					"CLEAR": {
						Name:       "CLEAR",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Clears all environment variables that are tracked in the internal map and unsets them from the actual environment.",
							"Only removes variables that were previously accessed or set through this ENVBASKIT instance.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Clear all environment variables that are in this BASKIT
								for key := range baskitMap {
									os.Unsetenv(key)
								}

								// Clear the BASKIT
								newMap := make(BaskitMap)
								this.NativeData = newMap
								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "CLEAR: invalid context"}
						},
					},
					// REFRESH method to sync with current environment variables
					"REFRESH": {
						Name:       "REFRESH",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Refreshes the internal map with all current environment variables.",
							"Discards any previous state and reloads from actual environment.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Clear current data
								for k := range baskitMap {
									delete(baskitMap, k)
								}

								// Reload from environment
								for _, envVar := range os.Environ() {
									for i, char := range envVar {
										if char == '=' {
											key := envVar[:i]
											value := envVar[i+1:]
											baskitMap[key] = environment.StringValue(value)
											break
										}
									}
								}

								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "REFRESH: invalid context"}
						},
					},
					// Inherit other methods from BASKIT by delegating to parent implementation
					// KEYS, VALUES, PAIRS work the same as BASKIT
					"KEYS":   baskitClass.PublicFunctions["KEYS"],
					"VALUES": baskitClass.PublicFunctions["VALUES"],
					"PAIRS":  baskitClass.PublicFunctions["PAIRS"],
					"MERGE":  baskitClass.PublicFunctions["MERGE"],
					"COPY": {
						Name:       "COPY",
						ReturnType: "ENVBASKIT",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Creates a new ENVBASKIT instance with the same data as the current instance.",
							"The copy is independent and does not sync changes back to the original.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								// Create new ENVBASKIT with copied data
								newEnvbaskitObj := NewEnvbaskitInstance()
								newBaskitMap := make(BaskitMap)
								for k, v := range baskitMap {
									newBaskitMap[k] = v
								}
								newEnvbaskitObj.NativeData = newBaskitMap
								return newEnvbaskitObj, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "COPY: invalid context"}
						},
					},
				},
				// PRIVATE constructor - cannot be called from LOLCODE
				PrivateFunctions: map[string]*environment.Function{
					"ENVBASKIT": {
						Name:       "ENVBASKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							// Initialize parent BASKIT data
							baskitMap := make(BaskitMap)
							this.NativeData = baskitMap

							// Load all current environment variables into the BASKIT
							for _, envVar := range os.Environ() {
								for i, char := range envVar {
									if char == '=' {
										key := envVar[:i]
										value := envVar[i+1:]
										baskitMap[key] = environment.StringValue(value)
										break
									}
								}
							}

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"SIZ": {
						Variable: environment.Variable{
							Name:          "SIZ",
							Type:          "INTEGR",
							IsLocked:      true,
							IsPublic:      true,
							Documentation: []string{"Number of environment variables."},
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
				QualifiedName:    "stdlib:SYSTEM.ENVBASKIT",
				ModulePath:       "stdlib:SYSTEM",
				MRO:              []string{"stdlib:SYSTEM.ENVBASKIT", "stdlib:MAPS.BASKIT"},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
		}
	})
	return systemClasses
}

// Global ENV variable
var globalENV *environment.ObjectInstance
var globalENVOnce sync.Once

func getGlobalENV() *environment.ObjectInstance {
	globalENVOnce.Do(func() {
		globalENV = NewEnvbaskitInstance()
		// Initialize it by calling the private constructor
		envbaskitClass := getSystemClasses()["ENVBASKIT"]
		constructor := envbaskitClass.PrivateFunctions["ENVBASKIT"]
		constructor.NativeImpl(nil, globalENV, []environment.Value{})
	})
	return globalENV
}

// Global SYSTEM variables - created once and reused
var systemVariablesOnce = sync.Once{}
var systemVariables map[string]*environment.Variable

func getSystemVariables() map[string]*environment.Variable {
	systemVariablesOnce.Do(func() {
		systemVariables = map[string]*environment.Variable{
			"ENV": {
				Name:          "ENV",
				Type:          "ENVBASKIT",
				Value:         getGlobalENV(),
				IsLocked:      true,
				IsPublic:      true,
				Documentation: []string{"Global environment variable manager.", "Pre-initialized ENVBASKIT instance containing all current environment variables."},
			},
			"OS": {
				Name:          "OS",
				Type:          "STRIN",
				Value:         environment.StringValue(goRuntime.GOOS),
				IsLocked:      true,
				IsPublic:      true,
				Documentation: []string{"Operating system name (e.g. windows, linux, darwin)."},
			},
			"ARCH": {
				Name:          "ARCH",
				Type:          "STRIN",
				Value:         environment.StringValue(goRuntime.GOARCH),
				IsLocked:      true,
				IsPublic:      true,
				Documentation: []string{"System architecture (e.g. amd64, 386, arm64)."},
			},
		}
	})
	return systemVariables
}

// RegisterSYSTEMInEnv registers SYSTEM module classes and variables in the given environment
// declarations: empty slice means import all, otherwise import only specified declarations
func RegisterSYSTEMInEnv(env *environment.Environment, declarations ...string) error {
	// First ensure MAPS module is registered (for BASKIT parent class)
	if err := RegisterMapsInEnv(env); err != nil {
		return runtime.Exception{Message: fmt.Sprintf("failed to register MAPS module dependency: %v", err)}
	}

	systemClasses := getSystemClasses()
	systemVariables := getSystemVariables()

	// If declarations is empty, import all classes and variables
	if len(declarations) == 0 {
		for _, class := range systemClasses {
			env.DefineClass(class)
		}
		for _, variable := range systemVariables {
			env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, variable.Documentation)
		}
		return nil
	}

	// Otherwise, import only specified declarations
	importedClasses := make(map[string]bool)
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)

		// Check if it's a class
		if class, exists := systemClasses[declUpper]; exists {
			env.DefineClass(class)
			importedClasses[declUpper] = true
		} else if variable, exists := systemVariables[declUpper]; exists {
			// Check if it's a variable
			// If importing ENV, we also need to import ENVBASKIT class since ENV is an instance of it
			if declUpper == "ENV" {
				if !importedClasses["ENVBASKIT"] {
					if envbaskitClass, exists := systemClasses["ENVBASKIT"]; exists {
						env.DefineClass(envbaskitClass)
						importedClasses["ENVBASKIT"] = true
					}
				}
			}

			env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, variable.Documentation)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown SYSTEM declaration: %s", decl)}
		}
	}

	return nil
}
