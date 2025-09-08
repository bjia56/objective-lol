package stdlib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
)

func TestHTTPModule(t *testing.T) {
	// Register HTTP module
	env := environment.NewEnvironment(nil)
	err := RegisterHTTPInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register HTTP module: %v", err)
	}

	// Check that classes are registered
	if _, err := env.GetClass("INTERWEB"); err != nil {
		t.Error("INTERWEB class not registered")
	}
	if _, err := env.GetClass("RESPONSE"); err != nil {
		t.Error("RESPONSE class not registered")
	}
}

func TestINTERWEBConstructor(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterHTTPInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register HTTP module: %v", err)
	}

	interwebClass, err := env.GetClass("INTERWEB")
	if err != nil {
		t.Fatalf("Failed to get INTERWEB class: %v", err)
	}
	constructor := interwebClass.PublicFunctions["INTERWEB"]

	// Create instance
	instance := &environment.ObjectInstance{
		Class:     interwebClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	_, err = constructor.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Check default values
	timeoutVar, exists := instance.Variables["TIMEOUT"]
	if !exists {
		t.Error("TIMEOUT variable not set")
	} else {
		val, err := timeoutVar.Get(instance)
		if err != nil {
			t.Errorf("Failed to get TIMEOUT value: %v", err)
		}
		if int(val.(environment.IntegerValue)) != 30 {
			t.Errorf("Expected TIMEOUT=30, got %v", val)
		}
	}

	headersVar, exists := instance.Variables["HEADERS"]
	if !exists {
		t.Error("HEADERS variable not set")
	} else if _, ok := headersVar.Value.(*environment.ObjectInstance); !ok {
		t.Error("HEADERS should be a BASKIT")
	}
}

func TestHTTPGETRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	// Setup HTTP module
	env := environment.NewEnvironment(nil)
	err := RegisterHTTPInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register HTTP module: %v", err)
	}

	// Create INTERWEB instance
	interwebClass, err := env.GetClass("INTERWEB")
	if err != nil {
		t.Fatalf("Failed to get INTERWEB class: %v", err)
	}
	instance := &environment.ObjectInstance{
		Class:     interwebClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := interwebClass.PublicFunctions["INTERWEB"]
	_, err = constructor.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Execute GET request
	getMethod := interwebClass.PublicFunctions["GET"]
	result, err := getMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue(server.URL),
	})
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	// Check response
	response, ok := result.(*environment.ObjectInstance)
	if !ok {
		t.Fatal("Expected ObjectInstance response")
	}

	// Check status
	statusVar := response.Variables["STATUS"]
	if statusVal, ok := statusVar.Value.(environment.IntegerValue); !ok || int(statusVal) != 200 {
		t.Errorf("Expected status 200, got %v", statusVar.Value)
	}

	// Check body
	bodyVar := response.Variables["BODY"]
	if bodyVal, ok := bodyVar.Value.(environment.StringValue); !ok || string(bodyVal) != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got %v", bodyVar.Value)
	}

	// Check success flag
	successVar := response.Variables["IS_SUCCESS"]
	if successVal, ok := successVar.Value.(environment.BoolValue); !ok || !bool(successVal) {
		t.Error("Expected IS_SUCCESS to be true")
	}

	// Check headers
	headersVar := response.Variables["HEADERS"]
	if headersBaskit, ok := headersVar.Value.(*environment.ObjectInstance); !ok {
		t.Error("Expected HEADERS to be a BASKIT")
	} else if baskitMap, ok := headersBaskit.NativeData.(BaskitMap); !ok {
		t.Error("Expected HEADERS to have BaskitMap data")
	} else {
		contentType, exists := baskitMap["Content-Type"]
		if !exists {
			t.Error("Content-Type header not found")
		} else if ctVal, ok := contentType.(environment.StringValue); !ok || string(ctVal) != "text/plain" {
			t.Errorf("Expected Content-Type 'text/plain', got %v", contentType)
		}
	}
}

func TestHTTPPOSTRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Read body
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		if string(body) != "test data" {
			t.Errorf("Expected body 'test data', got %s", string(body))
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}))
	defer server.Close()

	// Setup HTTP module
	env := environment.NewEnvironment(nil)
	err := RegisterHTTPInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register HTTP module: %v", err)
	}

	// Create INTERWEB instance
	interwebClass, err := env.GetClass("INTERWEB")
	if err != nil {
		t.Fatalf("Failed to get INTERWEB class: %v", err)
	}
	instance := &environment.ObjectInstance{
		Class:     interwebClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := interwebClass.PublicFunctions["INTERWEB"]
	_, err = constructor.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Execute POST request
	postMethod := interwebClass.PublicFunctions["POST"]
	result, err := postMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue(server.URL),
		environment.StringValue("test data"),
	})
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	// Check response
	response, ok := result.(*environment.ObjectInstance)
	if !ok {
		t.Fatal("Expected ObjectInstance response")
	}

	// Check status
	statusVar := response.Variables["STATUS"]
	if statusVal, ok := statusVar.Value.(environment.IntegerValue); !ok || int(statusVal) != 201 {
		t.Errorf("Expected status 201, got %v", statusVar.Value)
	}

	// Check body
	bodyVar := response.Variables["BODY"]
	if bodyVal, ok := bodyVar.Value.(environment.StringValue); !ok || string(bodyVal) != "Created" {
		t.Errorf("Expected body 'Created', got %v", bodyVar.Value)
	}
}

func TestHTTPHeaders(t *testing.T) {
	// Create test server that checks headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer token123" {
			t.Errorf("Expected Authorization header 'Bearer token123', got %s", auth)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Headers OK"))
	}))
	defer server.Close()

	// Setup HTTP module
	env := environment.NewEnvironment(nil)
	err := RegisterHTTPInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register HTTP module: %v", err)
	}

	// Create INTERWEB instance
	interwebClass, err := env.GetClass("INTERWEB")
	if err != nil {
		t.Fatalf("Failed to get INTERWEB class: %v", err)
	}
	instance := &environment.ObjectInstance{
		Class:     interwebClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := interwebClass.PublicFunctions["INTERWEB"]
	_, err = constructor.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Set headers
	headersVar := instance.Variables["HEADERS"]
	headersBaskit := headersVar.Value.(*environment.ObjectInstance)
	baskitMap := headersBaskit.NativeData.(BaskitMap)
	baskitMap["Authorization"] = environment.StringValue("Bearer token123")
	baskitMap["Content-Type"] = environment.StringValue("application/json")

	// Execute GET request
	getMethod := interwebClass.PublicFunctions["GET"]
	_, err = getMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue(server.URL),
	})
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
}

func TestJSONParsing(t *testing.T) {
	// Create test server that returns JSON
	jsonData := map[string]interface{}{
		"name":   "Alice",
		"age":    30,
		"email":  "alice@example.com",
		"active": true,
	}
	jsonBytes, _ := json.Marshal(jsonData)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}))
	defer server.Close()

	// Setup HTTP module
	env := environment.NewEnvironment(nil)
	err := RegisterHTTPInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register HTTP module: %v", err)
	}

	// Create INTERWEB instance and make request
	interwebClass, err := env.GetClass("INTERWEB")
	if err != nil {
		t.Fatalf("Failed to get INTERWEB class: %v", err)
	}
	instance := &environment.ObjectInstance{
		Class:     interwebClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := interwebClass.PublicFunctions["INTERWEB"]
	_, err = constructor.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Execute GET request
	getMethod := interwebClass.PublicFunctions["GET"]
	result, err := getMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue(server.URL),
	})
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	response := result.(*environment.ObjectInstance)

	// Test TO_JSON method
	responseClass, err := env.GetClass("RESPONSE")
	if err != nil {
		t.Fatalf("Failed to get RESPONSE class: %v", err)
	}
	toJSONMethod := responseClass.PublicFunctions["TO_JSON"]
	jsonResult, err := toJSONMethod.NativeImpl(nil, response, []environment.Value{})
	if err != nil {
		t.Fatalf("TO_JSON failed: %v", err)
	}

	// Check parsed JSON
	jsonBaskit, ok := jsonResult.(*environment.ObjectInstance)
	if !ok {
		t.Fatal("Expected BASKIT from TO_JSON")
	}

	baskitData, ok := jsonBaskit.NativeData.(BaskitMap)
	if !ok {
		t.Fatal("Expected BaskitMap from TO_JSON")
	}

	// Check name field
	nameVal, exists := baskitData["name"]
	if !exists {
		t.Error("name field not found in parsed JSON")
	} else if name, ok := nameVal.(environment.StringValue); !ok || string(name) != "Alice" {
		t.Errorf("Expected name 'Alice', got %v", nameVal)
	}

	// Check age field
	ageVal, exists := baskitData["age"]
	if !exists {
		t.Error("age field not found in parsed JSON")
	} else if age, ok := ageVal.(environment.IntegerValue); !ok || int(age) != 30 {
		t.Errorf("Expected age 30, got %v", ageVal)
	}

	// Check active field
	activeVal, exists := baskitData["active"]
	if !exists {
		t.Error("active field not found in parsed JSON")
	} else if active, ok := activeVal.(environment.BoolValue); !ok || !bool(active) {
		t.Errorf("Expected active true, got %v", activeVal)
	}
}

func TestSelectiveImport(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Import only INTERWEB class
	err := RegisterHTTPInEnv(env, "INTERWEB")
	if err != nil {
		t.Fatalf("Failed to register INTERWEB: %v", err)
	}

	// INTERWEB should be available
	if _, err := env.GetClass("INTERWEB"); err != nil {
		t.Error("INTERWEB class not registered")
	}

	// RESPONSE should be available (required dependency)
	if _, err := env.GetClass("RESPONSE"); err != nil {
		t.Error("RESPONSE class not registered")
	}
}

func TestErrorHandling(t *testing.T) {
	// Setup HTTP module
	env := environment.NewEnvironment(nil)
	err := RegisterHTTPInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register HTTP module: %v", err)
	}

	// Create INTERWEB instance
	interwebClass, err := env.GetClass("INTERWEB")
	if err != nil {
		t.Fatalf("Failed to get INTERWEB class: %v", err)
	}
	instance := &environment.ObjectInstance{
		Class:     interwebClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)

	constructor := interwebClass.PublicFunctions["INTERWEB"]
	_, err = constructor.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Fatalf("Constructor failed: %v", err)
	}

	// Test invalid URL
	getMethod := interwebClass.PublicFunctions["GET"]
	_, err = getMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("invalid-url"),
	})
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}
