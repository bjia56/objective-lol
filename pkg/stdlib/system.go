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
				PublicFunctions: map[string]*environment.Function{
					// Override PUT to also set environment variable
					"PUT": {
						Name:       "PUT",
						Parameters: []environment.Parameter{{Name: "KEY", Type: ""}, {Name: "VALUE", Type: ""}},
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
									return environment.NOTHIN, fmt.Errorf("failed to set environment variable %s: %v", keyStr, err)
								}

								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, fmt.Errorf("PUT: invalid context")
						},
					},
					// Override GET to also check environment variables
					"GET": {
						Name:       "GET",
						Parameters: []environment.Parameter{{Name: "KEY", Type: ""}},
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

								return nil, runtime.Exception{Message: fmt.Sprintf("Key '%s' not found in environment", keyStr)}
							}
							return environment.NOTHIN, fmt.Errorf("GET: invalid context")
						},
					},
					// Override CONTAINS to also check environment variables
					"CONTAINS": {
						Name:       "CONTAINS",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{{Name: "KEY", Type: ""}},
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
							return environment.NOTHIN, fmt.Errorf("CONTAINS: invalid context")
						},
					},
					// Override REMOVE to also unset environment variable
					"REMOVE": {
						Name:       "REMOVE",
						Parameters: []environment.Parameter{{Name: "KEY", Type: ""}},
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
									return nil, runtime.Exception{Message: fmt.Sprintf("Key '%s' not found in environment", keyStr)}
								}

								// Remove from BASKIT
								delete(baskitMap, keyStr)

								// Also unset environment variable
								if err := os.Unsetenv(keyStr); err != nil {
									return environment.NOTHIN, fmt.Errorf("failed to unset environment variable %s: %v", keyStr, err)
								}

								return value, nil
							}
							return environment.NOTHIN, fmt.Errorf("REMOVE: invalid context")
						},
					},
					// Override CLEAR to also clear environment variables managed by this instance
					"CLEAR": {
						Name:       "CLEAR",
						Parameters: []environment.Parameter{},
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
							return environment.NOTHIN, fmt.Errorf("CLEAR: invalid context")
						},
					},
					// REFRESH method to sync with current environment variables
					"REFRESH": {
						Name:       "REFRESH",
						Parameters: []environment.Parameter{},
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
							return environment.NOTHIN, fmt.Errorf("REFRESH: invalid context")
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
							return environment.NOTHIN, fmt.Errorf("COPY: invalid context")
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
							Name:     "SIZ",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if baskitMap, ok := this.NativeData.(BaskitMap); ok {
								return environment.IntegerValue(len(baskitMap)), nil
							}
							return environment.NOTHIN, fmt.Errorf("SIZ: invalid context")
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
				Name:     "ENV",
				Type:     "ENVBASKIT",
				Value:    getGlobalENV(),
				IsLocked: true,
				IsPublic: true,
			},
			"OS": {
				Name:     "OS",
				Type:     "STRIN",
				Value:    environment.StringValue(goRuntime.GOOS),
				IsLocked: true,
				IsPublic: true,
			},
			"ARCH": {
				Name:     "ARCH",
				Type:     "STRIN",
				Value:    environment.StringValue(goRuntime.GOARCH),
				IsLocked: true,
				IsPublic: true,
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
		return fmt.Errorf("failed to register MAPS module dependency: %v", err)
	}

	systemClasses := getSystemClasses()
	systemVariables := getSystemVariables()

	// If declarations is empty, import all classes and variables
	if len(declarations) == 0 {
		for _, class := range systemClasses {
			env.DefineClass(class)
		}
		for _, variable := range systemVariables {
			err := env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, nil)
			if err != nil {
				return fmt.Errorf("failed to define SYSTEM variable %s: %v", variable.Name, err)
			}
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

			err := env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, nil)
			if err != nil {
				return fmt.Errorf("failed to define SYSTEM variable %s: %v", variable.Name, err)
			}
		} else {
			return fmt.Errorf("unknown SYSTEM declaration: %s", decl)
		}
	}

	return nil
}
