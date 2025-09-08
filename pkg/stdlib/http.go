package stdlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// InterwebData stores the internal state of an INTERWEB client
type InterwebData struct {
	Client *http.Client
}

// ResponseData stores the internal state of a RESPONSE
type ResponseData struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// Global HTTP class definitions - created once and reused
var httpClassesOnce = sync.Once{}
var httpClasses map[string]*environment.Class

func getHTTPClasses() map[string]*environment.Class {
	httpClassesOnce.Do(func() {
		httpClasses = map[string]*environment.Class{
			"INTERWEB": {
				Name:          "INTERWEB",
				QualifiedName: "stdlib:HTTP.INTERWEB",
				ModulePath:    "stdlib:HTTP",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:HTTP.INTERWEB"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"INTERWEB": {
						Name:       "INTERWEB",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							// Initialize HTTP client
							client := &http.Client{
								Timeout: 30 * time.Second,
							}

							interwebData := &InterwebData{
								Client: client,
							}
							this.NativeData = interwebData

							// Initialize empty headers BASKIT
							headers := NewBaskitInstance()
							this.Variables["HEADERS"].Value = headers

							return environment.NOTHIN, nil
						},
					},
					// GET method
					"GET": {
						Name:       "GET",
						ReturnType: "RESPONSE",
						Parameters: []environment.Parameter{
							{Name: "url", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return executeHTTPRequest(this, "GET", args[0], environment.StringValue(""))
						},
					},
					// POST method
					"POST": {
						Name:       "POST",
						ReturnType: "RESPONSE",
						Parameters: []environment.Parameter{
							{Name: "url", Type: "STRIN"},
							{Name: "data", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return executeHTTPRequest(this, "POST", args[0], args[1])
						},
					},
					// PUT method
					"PUT": {
						Name:       "PUT",
						ReturnType: "RESPONSE",
						Parameters: []environment.Parameter{
							{Name: "url", Type: "STRIN"},
							{Name: "data", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return executeHTTPRequest(this, "PUT", args[0], args[1])
						},
					},
					// DELETE method
					"DELETE": {
						Name:       "DELETE",
						ReturnType: "RESPONSE",
						Parameters: []environment.Parameter{
							{Name: "url", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return executeHTTPRequest(this, "DELETE", args[0], environment.StringValue(""))
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"TIMEOUT": {
						Variable: environment.Variable{
							Name:     "TIMEOUT",
							Type:     "INTEGR",
							IsLocked: false,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if interwebData, ok := this.NativeData.(*InterwebData); ok {
								return environment.IntegerValue(int(interwebData.Client.Timeout.Seconds())), nil
							}
							return environment.IntegerValue(0), fmt.Errorf("invalid context for TIMEOUT")
						},
						NativeSet: func(this *environment.ObjectInstance, value environment.Value) error {
							if interwebData, ok := this.NativeData.(*InterwebData); ok {
								if intVal, ok := value.(environment.IntegerValue); ok {
									interwebData.Client.Timeout = time.Duration(int(intVal)) * time.Second
									return nil
								}
								return fmt.Errorf("TIMEOUT expects INTEGR value, got %s", value.Type())
							}
							return fmt.Errorf("invalid context for TIMEOUT")
						},
					},
					"HEADERS": {
						Variable: environment.Variable{
							Name:     "HEADERS",
							Type:     "BASKIT",
							Value:    NewBaskitInstance(),
							IsLocked: false,
							IsPublic: true,
						},
						// Since we are returning an object (BASKIT), we cannot
						// use NativeGet/NativeSet.
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"RESPONSE": {
				Name:          "RESPONSE",
				QualifiedName: "stdlib:HTTP.RESPONSE",
				ModulePath:    "stdlib:HTTP",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:HTTP.RESPONSE"},
				PublicFunctions: map[string]*environment.Function{
					// TO_JSON method
					"TO_JSON": {
						Name:       "TO_JSON",
						ReturnType: "BASKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							bodyVar, exists := this.Variables["BODY"]
							if !exists {
								return environment.NOTHIN, runtime.Exception{Message: "TO_JSON: BODY variable not found"}
							}

							bodyVal, ok := bodyVar.Value.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "TO_JSON: BODY is not a string"}
							}

							var jsonData interface{}
							err := json.Unmarshal([]byte(bodyVal), &jsonData)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("TO_JSON: failed to parse JSON: %v", err)}
							}

							// Convert JSON to BASKIT
							result := jsonToBaskit(jsonData)
							return result, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"STATUS": {
						Variable: environment.Variable{
							Name:     "STATUS",
							Type:     "INTEGR",
							Value:    environment.IntegerValue(0),
							IsLocked: true,
							IsPublic: true,
						},
					},
					"BODY": {
						Variable: environment.Variable{
							Name:     "BODY",
							Type:     "STRIN",
							Value:    environment.StringValue(""),
							IsLocked: true,
							IsPublic: true,
						},
					},
					"HEADERS": {
						Variable: environment.Variable{
							Name:     "HEADERS",
							Type:     "BASKIT",
							Value:    NewBaskitInstance(),
							IsLocked: true,
							IsPublic: true,
						},
					},
					"IS_SUCCESS": {
						Variable: environment.Variable{
							Name:     "IS_SUCCESS",
							Type:     "BOOL",
							Value:    environment.NO,
							IsLocked: true,
							IsPublic: true,
						},
					},
					"IS_ERROR": {
						Variable: environment.Variable{
							Name:     "IS_ERROR",
							Type:     "BOOL",
							Value:    environment.NO,
							IsLocked: true,
							IsPublic: true,
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
		}
	})
	return httpClasses
}

// executeHTTPRequest performs the actual HTTP request
func executeHTTPRequest(this *environment.ObjectInstance, method string, urlArg environment.Value, dataArg environment.Value) (environment.Value, error) {
	// Get URL
	urlVal, ok := urlArg.(environment.StringValue)
	if !ok {
		return environment.NOTHIN, fmt.Errorf("%s expects STRIN url, got %s", method, urlArg.Type())
	}

	// Get data (for POST/PUT)
	var dataVal environment.StringValue
	if dataArg != nil {
		var ok bool
		dataVal, ok = dataArg.(environment.StringValue)
		if !ok {
			return environment.NOTHIN, fmt.Errorf("%s expects STRIN data, got %s", method, dataArg.Type())
		}
	}

	// Get INTERWEB data
	interwebData, ok := this.NativeData.(*InterwebData)
	if !ok {
		return environment.NOTHIN, fmt.Errorf("%s: invalid context", method)
	}

	// Create request
	var req *http.Request
	var err error

	if method == "POST" || method == "PUT" {
		req, err = http.NewRequest(method, string(urlVal), bytes.NewReader([]byte(dataVal)))
	} else {
		req, err = http.NewRequest(method, string(urlVal), nil)
	}

	if err != nil {
		return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("HTTP request failed: %v", err)}
	}

	// Apply headers from HEADERS baskit
	headersVar, exists := this.Variables["HEADERS"]
	if exists {
		if headersBaskit, ok := headersVar.Value.(*environment.ObjectInstance); ok {
			if baskitMap, ok := headersBaskit.NativeData.(BaskitMap); ok {
				for key, value := range baskitMap {
					if valueStr, ok := value.(environment.StringValue); ok {
						req.Header.Set(key, string(valueStr))
					}
				}
			}
		}
	}

	// Execute request
	resp, err := interwebData.Client.Do(req)
	if err != nil {
		return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("HTTP request failed: %v", err)}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to read response body: %v", err)}
	}

	// Create RESPONSE object
	responseClass := getHTTPClasses()["RESPONSE"]
	env := environment.NewEnvironment(nil)
	env.DefineClass(responseClass)
	responseInstance := &environment.ObjectInstance{
		Environment: env,
		Class:       responseClass,
		Variables:   make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(responseInstance)

	// Set response data
	responseInstance.Variables["STATUS"].Value = environment.IntegerValue(resp.StatusCode)
	responseInstance.Variables["BODY"].Value = environment.StringValue(string(body))

	// Convert response headers to BASKIT
	headerBaskit := NewBaskitInstance()
	headerMap := headerBaskit.NativeData.(BaskitMap)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headerMap[key] = environment.StringValue(values[0])
		}
	}
	responseInstance.Variables["HEADERS"].Value = headerBaskit

	// Set success/error flags
	isSuccess := resp.StatusCode >= 200 && resp.StatusCode < 300
	isError := resp.StatusCode >= 400

	responseInstance.Variables["IS_SUCCESS"].Value = environment.BoolValue(isSuccess)
	responseInstance.Variables["IS_ERROR"].Value = environment.BoolValue(isError)

	return responseInstance, nil
}

// jsonToBaskit converts a JSON interface{} to a BASKIT
func jsonToBaskit(data interface{}) environment.Value {
	switch v := data.(type) {
	case map[string]interface{}:
		baskit := NewBaskitInstance()
		baskitMap := baskit.NativeData.(BaskitMap)
		for key, value := range v {
			baskitMap[key] = jsonToBaskit(value)
		}
		return baskit
	case []interface{}:
		bukkit := NewBukkitInstance()
		bukkitSlice := bukkit.NativeData.(BukkitSlice)
		for _, value := range v {
			bukkitSlice = append(bukkitSlice, jsonToBaskit(value))
		}
		bukkit.NativeData = bukkitSlice
		return bukkit
	case string:
		return environment.StringValue(v)
	case float64:
		// JSON numbers are float64, convert based on whether it's an integer
		if v == float64(int(v)) {
			return environment.IntegerValue(int(v))
		}
		return environment.DoubleValue(v)
	case bool:
		return environment.BoolValue(v)
	case nil:
		return environment.NOTHIN
	default:
		return environment.StringValue(fmt.Sprintf("%v", v))
	}
}

// RegisterHTTPInEnv registers HTTP classes in the given environment
func RegisterHTTPInEnv(env *environment.Environment, declarations ...string) error {
	httpClasses := getHTTPClasses()

	// If declarations is empty, import all classes
	if len(declarations) == 0 {
		for _, class := range httpClasses {
			env.DefineClass(class)
		}
		return nil
	}

	// Otherwise, import only specified classes
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if class, exists := httpClasses[declUpper]; exists {
			env.DefineClass(class)
			// If importing INTERWEB, also import RESPONSE (required dependency)
			if declUpper == "INTERWEB" {
				if responseClass, exists := httpClasses["RESPONSE"]; exists {
					env.DefineClass(responseClass)
				}
			}
		} else {
			return fmt.Errorf("unknown HTTP class: %s", decl)
		}
	}

	return nil
}
