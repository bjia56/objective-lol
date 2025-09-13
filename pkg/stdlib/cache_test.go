package stdlib

import (
	"testing"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
)

func TestMemStashBasicOperations(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterCACHEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register CACHE module: %v", err)
	}

	// Create MEMSTASH with capacity 3
	memStashClass, exists := env.GetClass("MEMSTASH")
	if exists != nil {
		t.Fatal("MEMSTASH class not found")
	}

	instance := &environment.ObjectInstance{
		Class:     memStashClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)
	constructor := memStashClass.PublicFunctions["MEMSTASH"]

	// Test constructor
	_, err = constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(3)})
	if err != nil {
		t.Fatalf("MEMSTASH constructor failed: %v", err)
	}

	// Test PUT
	putFunc := memStashClass.PublicFunctions["PUT"]
	_, err = putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
		environment.StringValue("value1"),
	})
	if err != nil {
		t.Fatalf("PUT failed: %v", err)
	}

	// Test GET
	getFunc := memStashClass.PublicFunctions["GET"]
	result, err := getFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if result != environment.StringValue("value1") {
		t.Fatalf("Expected 'value1', got %v", result)
	}

	// Test CONTAINS
	containsFunc := memStashClass.PublicFunctions["CONTAINS"]
	result, err = containsFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	if err != nil {
		t.Fatalf("CONTAINS failed: %v", err)
	}
	if result != environment.YEZ {
		t.Fatalf("Expected YES, got %v", result)
	}

	// Test SIZ property
	sizVar := memStashClass.PublicVariables["SIZ"]
	sizeResult, err := sizVar.NativeGet(instance)
	if err != nil {
		t.Fatalf("SIZ property failed: %v", err)
	}
	if sizeResult != environment.IntegerValue(1) {
		t.Fatalf("Expected size 1, got %v", sizeResult)
	}

	// Test DELETE
	deleteFunc := memStashClass.PublicFunctions["DELETE"]
	result, err = deleteFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	if err != nil {
		t.Fatalf("DELETE failed: %v", err)
	}
	if result != environment.YEZ {
		t.Fatalf("Expected YES (deleted), got %v", result)
	}

	// Verify deletion
	sizeResult, err = sizVar.NativeGet(instance)
	if err != nil {
		t.Fatalf("SIZ property failed: %v", err)
	}
	if sizeResult != environment.IntegerValue(0) {
		t.Fatalf("Expected size 0 after delete, got %v", sizeResult)
	}
}

func TestMemStashLRUEviction(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterCACHEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register CACHE module: %v", err)
	}

	// No interpreter needed for direct function calls

	memStashClass, _ := env.GetClass("MEMSTASH")
	instance := &environment.ObjectInstance{
		Class:     memStashClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)
	constructor := memStashClass.PublicFunctions["MEMSTASH"]

	// Create MEMSTASH with capacity 2
	_, err = constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(2)})
	if err != nil {
		t.Fatalf("MEMSTASH constructor failed: %v", err)
	}

	putFunc := memStashClass.PublicFunctions["PUT"]
	getFunc := memStashClass.PublicFunctions["GET"]
	containsFunc := memStashClass.PublicFunctions["CONTAINS"]
	sizVar := memStashClass.PublicVariables["SIZ"]

	// Add two items
	putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"), environment.StringValue("value1"),
	})
	putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key2"), environment.StringValue("value2"),
	})

	// Verify size is 2
	sizeResult, _ := sizVar.NativeGet(instance)
	if sizeResult != environment.IntegerValue(2) {
		t.Fatalf("Expected size 2, got %v", sizeResult)
	}

	// Access key1 to make it most recently used
	getFunc.NativeImpl(nil, instance, []environment.Value{environment.StringValue("key1")})

	// Add third item, should evict key2 (LRU)
	putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key3"), environment.StringValue("value3"),
	})

	// Verify size is still 2
	sizeResult, _ = sizVar.NativeGet(instance)
	if sizeResult != environment.IntegerValue(2) {
		t.Fatalf("Expected size 2 after eviction, got %v", sizeResult)
	}

	// key1 should still exist (was accessed recently)
	result, _ := containsFunc.NativeImpl(nil, instance, []environment.Value{environment.StringValue("key1")})
	if result != environment.YEZ {
		t.Fatal("key1 should still exist")
	}

	// key2 should be evicted
	result, _ = containsFunc.NativeImpl(nil, instance, []environment.Value{environment.StringValue("key2")})
	if result != environment.NO {
		t.Fatal("key2 should be evicted")
	}

	// key3 should exist
	result, _ = containsFunc.NativeImpl(nil, instance, []environment.Value{environment.StringValue("key3")})
	if result != environment.YEZ {
		t.Fatal("key3 should exist")
	}
}

func TestMemStashClear(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterCACHEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register CACHE module: %v", err)
	}

	// No interpreter needed for direct function calls

	memStashClass, _ := env.GetClass("MEMSTASH")
	instance := &environment.ObjectInstance{
		Class:     memStashClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)
	constructor := memStashClass.PublicFunctions["MEMSTASH"]

	// Create and populate cache
	constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(5)})

	putFunc := memStashClass.PublicFunctions["PUT"]
	putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"), environment.StringValue("value1"),
	})
	putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key2"), environment.StringValue("value2"),
	})

	// Verify items exist
	sizVar := memStashClass.PublicVariables["SIZ"]
	sizeResult, _ := sizVar.NativeGet(instance)
	if sizeResult != environment.IntegerValue(2) {
		t.Fatalf("Expected size 2 before clear, got %v", sizeResult)
	}

	// Clear cache
	clearFunc := memStashClass.PublicFunctions["CLEAR"]
	_, err = clearFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Fatalf("CLEAR failed: %v", err)
	}

	// Verify cache is empty
	sizeResult, _ = sizVar.NativeGet(instance)
	if sizeResult != environment.IntegerValue(0) {
		t.Fatalf("Expected size 0 after clear, got %v", sizeResult)
	}

	containsFunc := memStashClass.PublicFunctions["CONTAINS"]
	result, _ := containsFunc.NativeImpl(nil, instance, []environment.Value{environment.StringValue("key1")})
	if result != environment.NO {
		t.Fatal("key1 should not exist after clear")
	}
}

func TestTimeStashBasicOperations(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterCACHEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register CACHE module: %v", err)
	}

	// No interpreter needed for direct function calls

	// Create TIMESTASH with 60 second TTL
	timeStashClass, exists := env.GetClass("TIMESTASH")
	if exists != nil {
		t.Fatal("TIMESTASH class not found")
	}

	instance := &environment.ObjectInstance{
		Class:     timeStashClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)
	constructor := timeStashClass.PublicFunctions["TIMESTASH"]

	// Test constructor
	_, err = constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(60)})
	if err != nil {
		t.Fatalf("TIMESTASH constructor failed: %v", err)
	}

	// Test PUT
	putFunc := timeStashClass.PublicFunctions["PUT"]
	_, err = putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
		environment.StringValue("value1"),
	})
	if err != nil {
		t.Fatalf("PUT failed: %v", err)
	}

	// Test GET
	getFunc := timeStashClass.PublicFunctions["GET"]
	result, err := getFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if result != environment.StringValue("value1") {
		t.Fatalf("Expected 'value1', got %v", result)
	}

	// Test CONTAINS
	containsFunc := timeStashClass.PublicFunctions["CONTAINS"]
	result, err = containsFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	if err != nil {
		t.Fatalf("CONTAINS failed: %v", err)
	}
	if result != environment.YEZ {
		t.Fatalf("Expected YES, got %v", result)
	}

	// Test SIZ property
	sizVar := timeStashClass.PublicVariables["SIZ"]
	sizeResult, err := sizVar.NativeGet(instance)
	if err != nil {
		t.Fatalf("SIZ property failed: %v", err)
	}
	if sizeResult != environment.IntegerValue(1) {
		t.Fatalf("Expected size 1, got %v", sizeResult)
	}
}

func TestTimeStashTTLExpiration(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterCACHEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register CACHE module: %v", err)
	}

	// No interpreter needed for direct function calls

	timeStashClass, _ := env.GetClass("TIMESTASH")
	instance := &environment.ObjectInstance{
		Class:     timeStashClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)
	constructor := timeStashClass.PublicFunctions["TIMESTASH"]

	// Create TIMESTASH with 1 second TTL for quick testing
	_, err = constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(1)})
	if err != nil {
		t.Fatalf("TIMESTASH constructor failed: %v", err)
	}

	putFunc := timeStashClass.PublicFunctions["PUT"]
	getFunc := timeStashClass.PublicFunctions["GET"]
	containsFunc := timeStashClass.PublicFunctions["CONTAINS"]
	sizVar := timeStashClass.PublicVariables["SIZ"]

	// Add item
	putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"), environment.StringValue("value1"),
	})

	// Verify item exists
	result, _ := containsFunc.NativeImpl(nil, instance, []environment.Value{environment.StringValue("key1")})
	if result != environment.YEZ {
		t.Fatal("Item should exist before expiration")
	}

	// Wait for expiration
	time.Sleep(1100 * time.Millisecond) // Wait slightly longer than TTL

	// Verify item is expired
	result, _ = containsFunc.NativeImpl(nil, instance, []environment.Value{environment.StringValue("key1")})
	if result != environment.NO {
		t.Fatal("Item should be expired")
	}

	// GET should return NOTHIN for expired item
	getResult, _ := getFunc.NativeImpl(nil, instance, []environment.Value{environment.StringValue("key1")})
	if getResult != environment.NOTHIN {
		t.Fatalf("Expected NOTHIN for expired item, got %v", getResult)
	}

	// Size should be 0 after expiration cleanup
	sizeResult, _ := sizVar.NativeGet(instance)
	if sizeResult != environment.IntegerValue(0) {
		t.Fatalf("Expected size 0 after expiration cleanup, got %v", sizeResult)
	}
}

func TestTimeStashClear(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterCACHEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register CACHE module: %v", err)
	}

	// No interpreter needed for direct function calls

	timeStashClass, _ := env.GetClass("TIMESTASH")
	instance := &environment.ObjectInstance{
		Class:     timeStashClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)
	constructor := timeStashClass.PublicFunctions["TIMESTASH"]

	// Create and populate cache
	constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(60)})

	putFunc := timeStashClass.PublicFunctions["PUT"]
	putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"), environment.StringValue("value1"),
	})
	putFunc.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key2"), environment.StringValue("value2"),
	})

	// Verify items exist
	sizVar := timeStashClass.PublicVariables["SIZ"]
	sizeResult, _ := sizVar.NativeGet(instance)
	if sizeResult != environment.IntegerValue(2) {
		t.Fatalf("Expected size 2 before clear, got %v", sizeResult)
	}

	// Clear cache
	clearFunc := timeStashClass.PublicFunctions["CLEAR"]
	_, err = clearFunc.NativeImpl(nil, instance, []environment.Value{})
	if err != nil {
		t.Fatalf("CLEAR failed: %v", err)
	}

	// Verify cache is empty
	sizeResult, _ = sizVar.NativeGet(instance)
	if sizeResult != environment.IntegerValue(0) {
		t.Fatalf("Expected size 0 after clear, got %v", sizeResult)
	}

	containsFunc := timeStashClass.PublicFunctions["CONTAINS"]
	result, _ := containsFunc.NativeImpl(nil, instance, []environment.Value{environment.StringValue("key1")})
	if result != environment.NO {
		t.Fatal("key1 should not exist after clear")
	}
}

func TestCACHEModuleSelectiveImport(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test importing only MEMSTASH
	err := RegisterCACHEInEnv(env, "MEMSTASH")
	if err != nil {
		t.Fatalf("Failed to register MEMSTASH: %v", err)
	}

	// MEMSTASH should exist
	_, err = env.GetClass("MEMSTASH")
	if err != nil {
		t.Fatal("MEMSTASH class should be imported")
	}

	// TIMESTASH should not exist
	_, err = env.GetClass("TIMESTASH")
	if err == nil {
		t.Fatal("TIMESTASH class should not be imported")
	}

	// STASH should not exist (not imported)
	_, err = env.GetClass("STASH")
	if err == nil {
		t.Fatal("STASH class should not be imported")
	}
}

func TestCACHEModuleInvalidImport(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test importing non-existent class
	err := RegisterCACHEInEnv(env, "INVALID_CLASS")
	if err == nil {
		t.Fatal("Expected error for invalid class import")
	}

	expectedError := "unknown CACHE declaration: INVALID_CLASS"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestConstructorValidation(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterCACHEInEnv(env)
	if err != nil {
		t.Fatalf("Failed to register CACHE module: %v", err)
	}

	// No interpreter needed for direct function calls

	// Test MEMSTASH with invalid capacity
	memStashClass, _ := env.GetClass("MEMSTASH")
	instance := &environment.ObjectInstance{
		Class:     memStashClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)
	constructor := memStashClass.PublicFunctions["MEMSTASH"]

	_, err = constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(0)})
	if err == nil {
		t.Fatal("Expected error for zero capacity")
	}

	_, err = constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(-1)})
	if err == nil {
		t.Fatal("Expected error for negative capacity")
	}

	// Test TIMESTASH with invalid TTL
	timeStashClass, _ := env.GetClass("TIMESTASH")
	instance = &environment.ObjectInstance{
		Class:     timeStashClass,
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(instance)
	constructor = timeStashClass.PublicFunctions["TIMESTASH"]

	_, err = constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(0)})
	if err == nil {
		t.Fatal("Expected error for zero TTL")
	}

	_, err = constructor.NativeImpl(nil, instance, []environment.Value{environment.IntegerValue(-1)})
	if err == nil {
		t.Fatal("Expected error for negative TTL")
	}
}
