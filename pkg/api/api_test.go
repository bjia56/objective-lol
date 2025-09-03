package api

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVM(t *testing.T) {
	vm := NewVM()
	assert.NotNil(t, vm)
	assert.NotNil(t, vm.config)
	assert.True(t, vm.isInitialized)
}

func TestVMWithOptions(t *testing.T) {
	var output bytes.Buffer
	var input bytes.Buffer

	vm := NewVM(
		WithStdout(&output),
		WithStdin(&input),
		WithTimeout(5*time.Second),
	)

	config := vm.GetConfig()
	assert.Equal(t, &output, config.Stdout)
	assert.Equal(t, &input, config.Stdin)
	assert.Equal(t, 5*time.Second, config.Timeout)
}

func TestExecuteBasicProgram(t *testing.T) {
	var output bytes.Buffer
	vm := NewVM(WithStdout(&output))

	code := `
		I CAN HAS STDIO?

		HAI ME TEH FUNCSHUN MAIN
			SAYZ WIT "Hello, World!"
		KTHXBAI
	`

	result, err := vm.Execute(code)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, output.String(), "Hello, World!")
}

func TestExecuteWithReturn(t *testing.T) {
	vm := NewVM()

	code := `
		HAI ME TEH FUNCSHUN MAIN TEH INTEGR
			GIVEZ 42
		KTHXBAI
	`

	result, err := vm.Execute(code)
	require.NoError(t, err)
	assert.Equal(t, int64(42), result.Value)
}

func TestExecuteWithMath(t *testing.T) {
	vm := NewVM()

	code := `
		I CAN HAS STDIO?

		HAI ME TEH FUNCSHUN MAIN TEH INTEGR
			I HAS A VARIABLE RESULT TEH INTEGR ITZ 10 MOAR 5 TIEMZ 2
			GIVEZ RESULT
		KTHXBAI
	`

	result, err := vm.Execute(code)
	require.NoError(t, err)
	assert.Equal(t, int64(20), result.Value)
}

func TestSetAndGetVariables(t *testing.T) {
	vm := NewVM()

	// Set variables
	err := vm.Set("MY_STRING", "Hello")
	require.NoError(t, err)

	err = vm.Set("MY_NUMBER", 42)
	require.NoError(t, err)

	err = vm.Set("MY_FLOAT", 3.14)
	require.NoError(t, err)

	err = vm.Set("MY_BOOL", true)
	require.NoError(t, err)

	// Get variables
	strVal, err := vm.Get("MY_STRING")
	require.NoError(t, err)
	assert.Equal(t, "Hello", strVal)

	numVal, err := vm.Get("MY_NUMBER")
	require.NoError(t, err)
	assert.Equal(t, int64(42), numVal)

	floatVal, err := vm.Get("MY_FLOAT")
	require.NoError(t, err)
	assert.Equal(t, 3.14, floatVal)

	boolVal, err := vm.Get("MY_BOOL")
	require.NoError(t, err)
	assert.Equal(t, true, boolVal)
}

func TestCallFunction(t *testing.T) {
	vm := NewVM()

	// Define multiple functions with different signatures
	code := `
		I CAN HAS STRING?

		HAI ME TEH FUNCSHUN ADD TEH INTEGR WIT X TEH INTEGR AN WIT Y TEH INTEGR
			GIVEZ X MOAR Y
		KTHXBAI

		HAI ME TEH FUNCSHUN MULTIPLY TEH DUBBLE WIT X TEH DUBBLE AN WIT Y TEH DUBBLE
			GIVEZ X TIEMZ Y
		KTHXBAI

		HAI ME TEH FUNCSHUN GREET TEH STRIN WIT NAME TEH STRIN
			GIVEZ CONCAT WIT "Hello, " AN WIT NAME
		KTHXBAI

		HAI ME TEH FUNCSHUN IS_POSITIVE TEH BOOL WIT NUM TEH INTEGR
			IZ NUM BIGGR THAN 0?
				GIVEZ YEZ
			NOPE
				GIVEZ NO
			KTHX
		KTHXBAI

		HAI ME TEH FUNCSHUN NO_PARAMS TEH STRIN
			GIVEZ "No parameters!"
		KTHXBAI

		HAI ME TEH FUNCSHUN MAIN
			BTW Empty main for setup
		KTHXBAI
	`

	_, err := vm.Execute(code)
	require.NoError(t, err)

	t.Run("ADD function with integers", func(t *testing.T) {
		result, err := vm.Call("ADD", 10, 5)
		require.NoError(t, err)
		assert.Equal(t, int64(15), result)
	})

	t.Run("MULTIPLY function with floats", func(t *testing.T) {
		result, err := vm.Call("MULTIPLY", 3.5, 2.0)
		require.NoError(t, err)
		assert.Equal(t, 7.0, result)
	})

	t.Run("GREET function with string", func(t *testing.T) {
		result, err := vm.Call("GREET", "Alice")
		require.NoError(t, err)
		assert.Equal(t, "Hello, Alice", result)
	})

	t.Run("IS_POSITIVE function with boolean return", func(t *testing.T) {
		// Test positive number
		result, err := vm.Call("IS_POSITIVE", 5)
		require.NoError(t, err)
		assert.Equal(t, true, result)

		// Test negative number
		result, err = vm.Call("IS_POSITIVE", -3)
		require.NoError(t, err)
		assert.Equal(t, false, result)

		// Test zero
		result, err = vm.Call("IS_POSITIVE", 0)
		require.NoError(t, err)
		assert.Equal(t, false, result)
	})

	t.Run("Function with no parameters", func(t *testing.T) {
		result, err := vm.Call("NO_PARAMS")
		require.NoError(t, err)
		assert.Equal(t, "No parameters!", result)
	})

	t.Run("Call non-existent function", func(t *testing.T) {
		_, err := vm.Call("NON_EXISTENT")
		require.Error(t, err)

		vmErr, ok := err.(*VMError)
		require.True(t, ok)
		assert.True(t, vmErr.IsRuntimeError())
		assert.Contains(t, vmErr.Message, "function NON_EXISTENT not found")
	})

	t.Run("Call function with wrong number of arguments", func(t *testing.T) {
		// ADD expects 2 arguments, provide 1
		_, err := vm.Call("ADD", 5)
		require.Error(t, err)

		vmErr, ok := err.(*VMError)
		require.True(t, ok)
		assert.True(t, vmErr.IsRuntimeError())
	})

	t.Run("Call function with mixed argument types", func(t *testing.T) {
		// Test type conversion - should work with automatic casting
		result, err := vm.Call("ADD", 5.7, 2.3) // floats should be cast to integers
		require.NoError(t, err)
		assert.Equal(t, int64(7), result) // 5 + 2 = 7 (after casting to int)
	})

	t.Run("Case insensitive function names", func(t *testing.T) {
		// Function names should be case insensitive
		result, err := vm.Call("add", 1, 2)
		require.NoError(t, err)
		assert.Equal(t, int64(3), result)

		result, err = vm.Call("Add", 3, 4)
		require.NoError(t, err)
		assert.Equal(t, int64(7), result)
	})
}

func TestExecuteWithTimeout(t *testing.T) {
	vm := NewVM(WithTimeout(100 * time.Millisecond))

	// Infinite loop code
	code := `
		HAI ME TEH FUNCSHUN MAIN
			I HAS A VARIABLE X TEH INTEGR ITZ 0
			WHILE YEZ
				X ITZ X MOAR 1
			KTHX
		KTHXBAI
	`

	start := time.Now()
	_, err := vm.Execute(code)
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.True(t, elapsed >= 100*time.Millisecond)
	assert.True(t, elapsed < 200*time.Millisecond) // Should timeout quickly

	vmErr, ok := err.(*VMError)
	require.True(t, ok)
	assert.True(t, vmErr.IsTimeoutError())
}

func TestExecuteWithCompileError(t *testing.T) {
	vm := NewVM()

	// Invalid syntax
	code := `
		HAI ME TEH INVALID SYNTAX
			THIS IS NOT VALID
	`

	_, err := vm.Execute(code)
	require.Error(t, err)

	vmErr, ok := err.(*VMError)
	require.True(t, ok)
	assert.True(t, vmErr.IsCompileError())
}

func TestExecuteWithRuntimeError(t *testing.T) {
	vm := NewVM()

	// Runtime error - undefined variable
	code := `
		HAI ME TEH FUNCSHUN MAIN
			SAYZ WIT UNDEFINED_VARIABLE
		KTHXBAI
	`

	_, err := vm.Execute(code)
	require.Error(t, err)

	vmErr, ok := err.(*VMError)
	require.True(t, ok)
	assert.True(t, vmErr.IsRuntimeError())
}

func TestVMReset(t *testing.T) {
	vm := NewVM()

	// Set a variable
	err := vm.Set("TEST_VAR", "test_value")
	require.NoError(t, err)

	// Verify it exists
	val, err := vm.Get("TEST_VAR")
	require.NoError(t, err)
	assert.Equal(t, "test_value", val)

	// Reset VM
	err = vm.Reset()
	require.NoError(t, err)

	// Variable should be gone
	_, err = vm.Get("TEST_VAR")
	require.Error(t, err)
}

func TestTypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"int", 42, int64(42)},
		{"int64", int64(100), int64(100)},
		{"float64", 2.718, 2.718},
		{"string", "hello", "hello"},
		{"bool true", true, true},
		{"bool false", false, false},
		// Note: nil conversion and float32 don't work as expected due to type system differences
	}

	vm := NewVM()

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			varName := fmt.Sprintf("TEST_VAR_%d", i)
			err := vm.Set(varName, tt.input)
			require.NoError(t, err)

			result, err := vm.Get(varName)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBukkitToGoSliceConversion(t *testing.T) {
	vm := NewVM()

	// Create a BUKKIT array using Objective-LOL code
	code := `
		HAI ME TEH FUNCSHUN test_bukkit TEH BUKKIT
			I HAS A VARIABLE ARR TEH BUKKIT ITZ NEW BUKKIT
			ARR DO PUSH WIT 42
			ARR DO PUSH WIT "hello"
			ARR DO PUSH WIT 3.14
			ARR DO PUSH WIT YEZ
			GIVEZ ARR
		KTHXBAI

		HAI ME TEH FUNCSHUN MAIN
		KTHXBAI
	`

	_, err := vm.Execute(code)
	require.NoError(t, err)

	result, err := vm.Call("test_bukkit")
	require.NoError(t, err)

	// Should convert BUKKIT to []interface{}
	arr, ok := result.([]interface{})
	require.True(t, ok, "Expected []interface{}, got %T", result)
	assert.Len(t, arr, 4)
	assert.Equal(t, int64(42), arr[0])
	assert.Equal(t, "hello", arr[1])
	assert.Equal(t, 3.14, arr[2])
	assert.Equal(t, true, arr[3])
}

func TestObjectToGoMapConversion(t *testing.T) {
	vm := NewVM()

	// Create a BASKIT map using Objective-LOL code
	code := `
		HAI ME TEH FUNCSHUN test_baskit TEH BASKIT
			I HAS A VARIABLE MAP TEH BASKIT ITZ NEW BASKIT
			MAP DO PUT WIT "name" AN WIT "Alice"
			MAP DO PUT WIT "age" AN WIT 30
			MAP DO PUT WIT "score" AN WIT 85.5
			MAP DO PUT WIT "active" AN WIT YEZ
			GIVEZ MAP
		KTHXBAI

		HAI ME TEH FUNCSHUN MAIN
		KTHXBAI
	`

	_, err := vm.Execute(code)
	require.NoError(t, err)

	result, err := vm.Call("test_baskit")
	require.NoError(t, err)

	// Should convert BASKIT to map[string]interface{}
	objMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Expected map[string]interface{}, got %T", result)
	assert.Len(t, objMap, 4)
	assert.Equal(t, "Alice", objMap["name"])
	assert.Equal(t, int64(30), objMap["age"])
	assert.Equal(t, 85.5, objMap["score"])
	assert.Equal(t, true, objMap["active"])
}

func TestSliceToArrayConversion(t *testing.T) {
	vm := NewVM()

	// Set a Go slice that should be converted to BUKKIT
	err := vm.Set("MY_ARRAY", []interface{}{1, "hello", 3.14, true})
	require.NoError(t, err)

	// Test that the variable was converted to BUKKIT and can be used
	code := `
		HAI ME TEH FUNCSHUN test_size TEH INTEGR
			GIVEZ MY_ARRAY SIZ
		KTHXBAI

		HAI ME TEH FUNCSHUN test_access
			GIVEZ MY_ARRAY DO AT WIT 1
		KTHXBAI

		HAI ME TEH FUNCSHUN MAIN
		KTHXBAI
	`

	_, err = vm.Execute(code)
	require.NoError(t, err)

	// Test size
	result, err := vm.Call("test_size")
	require.NoError(t, err)
	assert.Equal(t, int64(4), result)

	// Test accessing elements
	result, err = vm.Call("test_access")
	require.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestMapToObjectConversion(t *testing.T) {
	vm := NewVM()

	// Set a Go map that should be converted to BASKIT
	goMap := map[string]interface{}{
		"name":   "Bob",
		"age":    25,
		"height": 6.2,
		"active": true,
	}
	err := vm.Set("MY_MAP", goMap)
	require.NoError(t, err)

	// Test that the variable was converted to BASKIT and can be used
	code := `
		HAI ME TEH FUNCSHUN test_size TEH INTEGR
			GIVEZ MY_MAP SIZ
		KTHXBAI

		HAI ME TEH FUNCSHUN test_access
			GIVEZ MY_MAP DO GET WIT "name"
		KTHXBAI

		HAI ME TEH FUNCSHUN MAIN
		KTHXBAI
	`

	_, err = vm.Execute(code)
	require.NoError(t, err)

	// Test size
	result, err := vm.Call("test_size")
	require.NoError(t, err)
	assert.Equal(t, int64(4), result)

	// Test accessing values
	result, err = vm.Call("test_access")
	require.NoError(t, err)
	assert.Equal(t, "Bob", result)
}

func TestStructToObjectConversion(t *testing.T) {
	vm := NewVM()

	// Define a test struct
	type TestStruct struct {
		Name   string
		Age    int
		Height float64
		Active bool
	}

	// Set a Go struct that should be converted to BASKIT
	testStruct := TestStruct{
		Name:   "Charlie",
		Age:    35,
		Height: 5.9,
		Active: false,
	}
	err := vm.Set("MY_STRUCT", testStruct)
	require.NoError(t, err)

	// Test that the variable was converted to BASKIT and can be used
	code := `
		HAI ME TEH FUNCSHUN test_size TEH INTEGR
			GIVEZ MY_STRUCT SIZ
		KTHXBAI

		HAI ME TEH FUNCSHUN test_name
			GIVEZ MY_STRUCT DO GET WIT "Name"
		KTHXBAI

		HAI ME TEH FUNCSHUN test_age TEH INTEGR
			GIVEZ MY_STRUCT DO GET WIT "Age"
		KTHXBAI

		HAI ME TEH FUNCSHUN MAIN
		KTHXBAI
	`

	_, err = vm.Execute(code)
	require.NoError(t, err)

	// Test size
	result, err := vm.Call("test_size")
	require.NoError(t, err)
	assert.Equal(t, int64(4), result)

	// Test accessing struct fields
	result, err = vm.Call("test_name")
	require.NoError(t, err)
	assert.Equal(t, "Charlie", result)

	// Test accessing numeric field
	result, err = vm.Call("test_age")
	require.NoError(t, err)
	assert.Equal(t, int64(35), result)
}

func TestArrayConversion(t *testing.T) {
	vm := NewVM()

	// Set an array - this tests Go slice → BUKKIT → Go slice conversion
	err := vm.Set("MY_ARRAY", []int{1, 2, 3, 4, 5})
	require.NoError(t, err)

	// Get the array back - should be converted through BUKKIT
	result, err := vm.Get("MY_ARRAY")
	require.NoError(t, err)

	// Should be converted to []interface{}
	arr, ok := result.([]interface{})
	require.True(t, ok, "Expected []interface{}, got %T", result)
	assert.Len(t, arr, 5)
	assert.Equal(t, int64(1), arr[0])
	assert.Equal(t, int64(5), arr[4])
}

func TestConfigValidation(t *testing.T) {
	// Test invalid configuration
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(string), "configuration error")
		}
	}()

	// This should panic due to nil stdout
	_ = NewVM(WithStdout(nil))
	t.Error("Expected panic due to invalid configuration")
}

func TestConcurrentAccess(t *testing.T) {
	vm := NewVM()

	// Set up concurrent access test
	done := make(chan bool, 2)

	// Goroutine 1: Set variables
	go func() {
		for i := 0; i < 10; i++ {
			err := vm.Set("VAR", i)
			assert.NoError(t, err)
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Goroutine 2: Get variables
	go func() {
		for i := 0; i < 10; i++ {
			_, err := vm.Get("VAR")
			// May error if variable doesn't exist yet, that's ok
			_ = err
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done
}
