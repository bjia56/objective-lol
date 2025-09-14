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

// moduleSystemCategories defines the order that categories should be rendered in documentation
var moduleSystemCategories = []string{
	"global-variables",
	"environment-management",
	"environment-variables",
	"environment-properties",
	"system-information",
}

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
					"Automatically syncs with the actual process environment and provides enhanced functionality.",
					"",
					"@class ENVBASKIT",
					"@inherits MAPS.BASKIT",
					"@example Basic environment variable access",
					"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
					"I HAS A VARIABLE HOME_DIR TEH STRIN ITZ ENV DO GET WIT \"HOME\"",
					"SAYZ WIT \"Home directory: \" MOAR HOME_DIR",
					"@example Set and get custom variables",
					"ENV DO PUT WIT \"MY_APP_CONFIG\" AN WIT \"/etc/myapp.conf\"",
					"ENV DO PUT WIT \"MY_APP_PORT\" AN WIT 8080",
					"I HAS A VARIABLE CONFIG TEH STRIN ITZ ENV DO GET WIT \"MY_APP_CONFIG\"",
					"I HAS A VARIABLE PORT TEH NUMBR ITZ ENV DO GET WIT \"MY_APP_PORT\"",
					"@example Check variable existence",
					"IZ ENV DO CONTAINS WIT \"PATH\"?",
					"    I HAS A VARIABLE PATH_VAR TEH STRIN ITZ ENV DO GET WIT \"PATH\"",
					"    SAYZ WIT \"PATH: \" MOAR PATH_VAR",
					"KTHX",
					"@example Environment variable cleanup",
					"ENV DO PUT WIT \"TEMP_SESSION\" AN WIT \"session123\"",
					"BTW ... use session ...",
					"ENV DO REMOVE WIT \"TEMP_SESSION\" BTW Clean up",
					"@example Refresh from external changes",
					"BTW External process modifies environment...",
					"ENV DO REFRESH BTW Reload current environment state",
					"I HAS A VARIABLE NEW_VAR TEH STRIN ITZ ENV DO GET WIT \"EXTERNAL_VAR\"",
					"@example Create isolated environment copy",
					"I HAS A VARIABLE TEST_ENV TEH ENVBASKIT ITZ ENV DO COPY",
					"TEST_ENV DO PUT WIT \"TEST_VAR\" AN WIT \"test_value\"",
					"BTW ... run tests ...",
					"BTW Original ENV unaffected by TEST_ENV changes",
					"@note Inherits all BASKIT methods (KEYS, VALUES, PAIRS, MERGE, etc.)",
					"@note Automatically syncs with actual process environment variables",
					"@note Environment variables are cached after first access for performance",
					"@note Changes made through ENVBASKIT affect both internal state and actual environment",
					"@see BASKIT, ENV (global variable)",
					"@category environment-management",
				},
				PublicFunctions: map[string]*environment.Function{
					// Override PUT to also set environment variable
					"PUT": {
						Name:       "PUT",
						Parameters: []environment.Parameter{{Name: "KEY", Type: "STRIN"}, {Name: "VALUE", Type: ""}},
						Documentation: []string{
							"Sets an environment variable both in the internal map and the actual process environment.",
							"",
							"@syntax <envbaskit> DO PUT WIT <key> AN WIT <value>",
							"@param {STRIN} key - The environment variable name",
							"@param value - The value to set (converted to string)",
							"@returns {NOTHIN}",
							"@example Set environment variable",
							"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
							"ENV DO PUT WIT \"MY_VAR\" AN WIT \"my_value\"",
							"@example Set numeric value",
							"ENV DO PUT WIT \"PORT\" AN WIT 8080",
							"@example Override existing variable",
							"ENV DO PUT WIT \"PATH\" AN WIT \"/usr/local/bin:/usr/bin\"",
							"@note Changes are reflected in both the ENVBASKIT instance and the actual process environment",
							"@note Values are automatically converted to strings",
							"@see GET, CONTAINS, REMOVE",
							"@category environment-variables",
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
							"",
							"@syntax <envbaskit> DO GET WIT <key>",
							"@param {STRIN} key - The environment variable name to retrieve",
							"@returns The value of the environment variable",
							"@throws Exception if the key is not found in either the internal map or actual environment",
							"@example Get environment variable",
							"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
							"I HAS A VARIABLE HOME_DIR TEH STRIN ITZ ENV DO GET WIT \"HOME\"",
							"SAYZ WIT \"Home directory: \"",
							"SAYZ WIT HOME_DIR",
							"@example Get custom variable",
							"ENV DO PUT WIT \"MY_APP_CONFIG\" AN WIT \"/etc/myapp.conf\"",
							"I HAS A VARIABLE CONFIG_PATH TEH STRIN ITZ ENV DO GET WIT \"MY_APP_CONFIG\"",
							"@example Handle missing variable",
							"IZ ENV DO CONTAINS WIT \"NON_EXISTENT_VAR\"?",
							"    I HAS A VARIABLE VALUE TEH STRIN ITZ ENV DO GET WIT \"NON_EXISTENT_VAR\"",
							"    SAYZ WIT VALUE",
							"KTHX",
							"@note Searches internal map first, then actual environment variables",
							"@note Values from environment are cached in the internal map after first access",
							"@see PUT, CONTAINS",
							"@category environment-variables",
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
							"",
							"@syntax <envbaskit> DO CONTAINS WIT <key>",
							"@param {STRIN} key - The environment variable name to check",
							"@returns {BOOL} YEZ if the variable exists, NO otherwise",
							"@example Check if variable exists",
							"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
							"IZ ENV DO CONTAINS WIT \"HOME\"?",
							"    SAYZ WIT \"HOME variable exists\"",
							"KTHX",
							"@example Check before accessing",
							"IZ ENV DO CONTAINS WIT \"MY_CONFIG\"?",
							"    I HAS A VARIABLE CONFIG TEH STRIN ITZ ENV DO GET WIT \"MY_CONFIG\"",
							"    SAYZ WIT \"Config: \" MOAR CONFIG",
							"NOPE",
							"    SAYZ WIT \"MY_CONFIG not set, using default\"",
							"KTHX",
							"@example Check system variables",
							"IZ ENV DO CONTAINS WIT \"PATH\"?",
							"    I HAS A VARIABLE PATH_VAR TEH STRIN ITZ ENV DO GET WIT \"PATH\"",
							"    SAYZ WIT \"PATH is set to: \" MOAR PATH_VAR",
							"KTHX",
							"@note Searches both internal map and actual environment variables",
							"@note Environment variables found are automatically cached in the internal map",
							"@see GET, PUT",
							"@category environment-variables",
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
							"",
							"@syntax <envbaskit> DO REMOVE WIT <key>",
							"@param {STRIN} key - The environment variable name to remove",
							"@returns The previous value of the environment variable",
							"@throws Exception if the key is not found",
							"@example Remove environment variable",
							"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
							"ENV DO PUT WIT \"TEMP_VAR\" AN WIT \"temporary_value\"",
							"I HAS A VARIABLE OLD_VALUE TEH STRIN ITZ ENV DO REMOVE WIT \"TEMP_VAR\"",
							"SAYZ WIT \"Removed value was: \" MOAR OLD_VALUE",
							"@example Clean up after use",
							"ENV DO PUT WIT \"SESSION_ID\" AN WIT \"abc123\"",
							"BTW ... use session ...",
							"ENV DO REMOVE WIT \"SESSION_ID\" BTW Clean up",
							"@example Handle removal errors",
							"IZ ENV DO CONTAINS WIT \"NON_EXISTENT\"?",
							"    ENV DO REMOVE WIT \"NON_EXISTENT\"",
							"NOPE",
							"    SAYZ WIT \"Variable doesn't exist, can't remove\"",
							"KTHX",
							"@note Removes from both internal map and actual process environment",
							"@note Returns the value that was removed",
							"@note Throws exception if variable doesn't exist in either location",
							"@see PUT, GET, CLEAR",
							"@category environment-variables",
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
							"",
							"@syntax <envbaskit> DO CLEAR",
							"@returns {NOTHIN}",
							"@example Clear all tracked variables",
							"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
							"ENV DO PUT WIT \"VAR1\" AN WIT \"value1\"",
							"ENV DO PUT WIT \"VAR2\" AN WIT \"value2\"",
							"SAYZ WIT \"Before clear: \" MOAR ENV SIZ",
							"ENV DO CLEAR",
							"SAYZ WIT \"After clear: \" MOAR ENV SIZ",
							"@example Fresh start with environment",
							"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
							"BTW ... modify some environment variables ...",
							"ENV DO CLEAR BTW Reset to clean state",
							"@note Only clears variables that were accessed through this ENVBASKIT instance",
							"@note Does not affect environment variables not accessed through this instance",
							"@note After clearing, the instance starts fresh and will reload from environment on next access",
							"@see REFRESH, REMOVE",
							"@category environment-variables",
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
							"",
							"@syntax <envbaskit> DO REFRESH",
							"@returns {NOTHIN}",
							"@example Refresh after external changes",
							"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
							"BTW External process changes environment...",
							"ENV DO REFRESH BTW Reload current environment state",
							"I HAS A VARIABLE NEW_VAR TEH STRIN ITZ ENV DO GET WIT \"NEW_VAR\"",
							"@example Reset to clean environment state",
							"ENV DO PUT WIT \"TEMP1\" AN WIT \"value1\"",
							"ENV DO PUT WIT \"TEMP2\" AN WIT \"value2\"",
							"ENV DO REFRESH BTW Discard changes, reload from environment",
							"@example Periodic refresh",
							"WHILE YEZ",
							"    ENV DO REFRESH",
							"    I HAS A VARIABLE STATUS TEH STRIN ITZ ENV DO GET WIT \"STATUS_VAR\"",
							"    BTW ... process status ...",
							"    I HAS A VARIABLE DELAY TEH NUMBR ITZ 5",
							"    SLEEPZ WIT DELAY",
							"KTHX",
							"@note Discards all previous changes made through this ENVBASKIT instance",
							"@note Reloads all current environment variables into the internal map",
							"@note Useful when external processes have modified the environment",
							"@see CLEAR, GET",
							"@category environment-variables",
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
							"",
							"@syntax <envbaskit> DO COPY",
							"@returns {ENVBASKIT} A new ENVBASKIT instance with copied data",
							"@example Create independent copy",
							"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
							"ENV DO PUT WIT \"MY_VAR\" AN WIT \"original_value\"",
							"I HAS A VARIABLE ENV_COPY TEH ENVBASKIT ITZ ENV DO COPY",
							"ENV_COPY DO PUT WIT \"MY_VAR\" AN WIT \"modified_value\"",
							"SAYZ WIT \"Original: \" MOAR (ENV DO GET WIT \"MY_VAR\")",
							"SAYZ WIT \"Copy: \" MOAR (ENV_COPY DO GET WIT \"MY_VAR\")",
							"@example Backup current state",
							"I HAS A VARIABLE BACKUP TEH ENVBASKIT ITZ ENV DO COPY",
							"BTW ... modify environment ...",
							"ENV DO CLEAR",
							"ENV DO MERGE WIT BACKUP BTW Restore from backup",
							"@example Isolated environment testing",
							"I HAS A VARIABLE TEST_ENV TEH ENVBASKIT ITZ ENV DO COPY",
							"TEST_ENV DO PUT WIT \"TEST_VAR\" AN WIT \"test_value\"",
							"BTW ... run tests with TEST_ENV ...",
							"BTW Changes don't affect original ENV",
							"@note The copy is completely independent of the original",
							"@note Changes to the copy don't affect the original and vice versa",
							"@note Both copies will sync with the actual environment independently",
							"@see MERGE, CLEAR",
							"@category environment-variables",
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
							Name:     "SIZ",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Number of environment variables currently tracked in this ENVBASKIT instance.",
								"",
								"@var {INTEGR} SIZ",
								"@readonly",
								"@example Check number of variables",
								"I HAS A VARIABLE ENV TEH ENVBASKIT ITZ NEW ENVBASKIT",
								"SAYZ WIT \"Environment variables: \" MOAR ENV SIZ",
								"@example Monitor variable count",
								"ENV DO PUT WIT \"VAR1\" AN WIT \"value1\"",
								"ENV DO PUT WIT \"VAR2\" AN WIT \"value2\"",
								"SAYZ WIT \"After adding: \" MOAR ENV SIZ",
								"ENV DO REMOVE WIT \"VAR1\"",
								"SAYZ WIT \"After removing: \" MOAR ENV SIZ",
								"@note Only counts variables accessed through this ENVBASKIT instance",
								"@note Does not include all environment variables, only those tracked internally",
								"@see KEYS, VALUES",
								"@category environment-properties",
							},
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
				Name:     "ENV",
				Type:     "ENVBASKIT",
				Value:    getGlobalENV(),
				IsLocked: true,
				IsPublic: true,
				Documentation: []string{
					"Global environment variable manager. Pre-initialized ENVBASKIT instance containing all current environment variables.",
					"",
					"@var {ENVBASKIT} ENV",
					"@global",
					"@readonly",
					"@example Access environment variables",
					"I HAS A VARIABLE HOME_DIR TEH STRIN ITZ ENV DO GET WIT \"HOME\"",
					"I HAS A VARIABLE USER_NAME TEH STRIN ITZ ENV DO GET WIT \"USER\"",
					"SAYZ WIT \"Hello \" MOAR USER_NAME MOAR \", your home is \" MOAR HOME_DIR",
					"@example Check system paths",
					"I HAS A VARIABLE PATH_VAR TEH STRIN ITZ ENV DO GET WIT \"PATH\"",
					"I HAS A VARIABLE SHELL_VAR TEH STRIN ITZ ENV DO GET WIT \"SHELL\"",
					"SAYZ WIT \"Using shell: \" MOAR SHELL_VAR",
					"@example Set custom environment variables",
					"ENV DO PUT WIT \"MY_APP_CONFIG\" AN WIT \"/etc/myapp/config.json\"",
					"ENV DO PUT WIT \"MY_APP_DEBUG\" AN WIT YEZ",
					"@example Environment variable cleanup",
					"ENV DO PUT WIT \"TEMP_FILE\" AN WIT \"/tmp/temp.txt\"",
					"BTW ... use temp file ...",
					"ENV DO REMOVE WIT \"TEMP_FILE\" BTW Clean up",
					"@note Automatically populated with all current environment variables",
					"@note Changes made through ENV affect the actual process environment",
					"@note Use ENV DO REFRESH to reload from external environment changes",
					"@see ENVBASKIT, OS, ARCH",
					"@category global-variables",
				},
			},
			"OS": {
				Name:     "OS",
				Type:     "STRIN",
				Value:    environment.StringValue(goRuntime.GOOS),
				IsLocked: true,
				IsPublic: true,
				Documentation: []string{
					"Operating system name (e.g. windows, linux, darwin).",
					"",
					"@var {STRIN} OS",
					"@global",
					"@readonly",
					"@example Check operating system",
					"IZ OS SAEM AS \"linux\"?",
					"    SAYZ WIT \"Running on Linux\"",
					"KTHX",
					"@example OS-specific file paths",
					"I HAS A VARIABLE PATH_SEP TEH STRIN",
					"IZ OS SAEM AS \"windows\"?",
					"    PATH_SEP ITZ \"\\\\\"",
					"NOPE",
					"    PATH_SEP ITZ \"/\"",
					"KTHX",
					"SAYZ WIT \"Path separator: \" MOAR PATH_SEP",
					"@example OS-specific behavior",
					"IZ OS SAEM AS \"darwin\"?",
					"    SAYZ WIT \"Running on macOS\"",
					"KTHX",
					"@example Cross-platform scripting",
					"I HAS A VARIABLE CMD TEH STRIN",
					"IZ OS SAEM AS \"windows\"?",
					"    CMD ITZ \"dir\"",
					"NOPE",
					"    CMD ITZ \"ls -la\"",
					"KTHX",
					"@note Common values: 'linux', 'darwin', 'windows', 'freebsd'",
					"@note Use with ARCH to determine full platform information",
					"@note Value is determined at runtime by the Go runtime",
					"@see ARCH, ENV",
					"@category system-information",
				},
			},
			"ARCH": {
				Name:     "ARCH",
				Type:     "STRIN",
				Value:    environment.StringValue(goRuntime.GOARCH),
				IsLocked: true,
				IsPublic: true,
				Documentation: []string{
					"System architecture (e.g. amd64, 386, arm64).",
					"",
					"@var {STRIN} ARCH",
					"@global",
					"@readonly",
					"@example Check system architecture",
					"IZ ARCH SAEM AS \"amd64\"?",
					"    SAYZ WIT \"Running on 64-bit x86\"",
					"KTHX",
					"@example Architecture-specific operations",
					"IZ ARCH SAEM AS \"arm64\"?",
					"    SAYZ WIT \"Running on ARM 64-bit\"",
					"KTHX",
					"@example Platform detection",
					"SAYZ WIT \"Platform: \" MOAR OS MOAR \"-\" MOAR ARCH",
					"@example Memory model detection",
					"I HAS A VARIABLE IS_64BIT TEH BOOL",
					"IZ ARCH SAEM AS \"amd64\" OR ARCH SAEM AS \"arm64\"?",
					"    IS_64BIT ITZ YEZ",
					"NOPE",
					"    IS_64BIT ITZ NO",
					"KTHX",
					"@example Architecture-specific file selection",
					"I HAS A VARIABLE LIB_PATH TEH STRIN",
					"LIB_PATH ITZ \"/usr/lib/\" MOAR ARCH",
					"SAYZ WIT \"Library path: \" MOAR LIB_PATH",
					"@note Common values: 'amd64', '386', 'arm64', 'arm', 'ppc64', 's390x'",
					"@note Use with OS to determine full platform information",
					"@note Value is determined at runtime by the Go runtime",
					"@see OS, ENV",
					"@category system-information",
				},
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
