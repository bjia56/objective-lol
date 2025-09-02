package stdlib

import (
	"testing"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterTIME(t *testing.T) {
	env := environment.NewEnvironment(nil)

	err := RegisterTIMEInEnv(env)
	require.NoError(t, err)

	// Test that TIME class is registered
	_, err = env.GetClass("DATE")
	assert.NoError(t, err, "DATE class should be registered")

	// Test that TIME function is registered
	_, err = env.GetFunction("SLEEP")
	assert.NoError(t, err, "SLEEP function should be registered")
}

func TestTimeDATEClass(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterTIMEInEnv(env)

	dateClass, err := env.GetClass("DATE")
	require.NoError(t, err)

	// Create a DATE instance
	instanceInterface, err := env.NewObjectInstance("DATE")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok, "Should return ObjectInstance")

	// Initialize the date (call constructor)
	constructor, exists := dateClass.PublicFunctions["DATE"]
	require.True(t, exists, "DATE constructor should exist")

	before := time.Now()
	_, err = constructor.NativeImpl(nil, instance, []types.Value{})
	require.NoError(t, err)
	after := time.Now()

	// Verify that the instance was initialized with a time
	dateTime, ok := instance.NativeData.(time.Time)
	require.True(t, ok, "NativeData should be time.Time")
	assert.True(t, dateTime.After(before.Add(-time.Second)), "Date should be recent")
	assert.True(t, dateTime.Before(after.Add(time.Second)), "Date should be recent")
}

func TestTimeDATEMethods(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterTIMEInEnv(env)

	dateClass, err := env.GetClass("DATE")
	require.NoError(t, err)

	instanceInterface, err := env.NewObjectInstance("DATE")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok, "Should return ObjectInstance")

	// Set a known date for testing
	knownDate := time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC)
	instance.NativeData = knownDate

	tests := []struct {
		method   string
		expected types.Value
	}{
		{"YEAR", types.IntegerValue(2023)},
		{"MONTH", types.IntegerValue(12)},
		{"DAY", types.IntegerValue(25)},
		{"HOUR", types.IntegerValue(15)},
		{"MINUTE", types.IntegerValue(30)},
		{"SECOND", types.IntegerValue(45)},
		{"MILLISECOND", types.IntegerValue(123)}, // 123456789 nanoseconds / 1e6
		{"NANOSECOND", types.IntegerValue(123456789)},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			method, exists := dateClass.PublicFunctions[test.method]
			require.True(t, exists, "Method %s should exist", test.method)

			result, err := method.NativeImpl(nil, instance, []types.Value{})
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestTimeDATEFormat(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterTIMEInEnv(env)

	dateClass, err := env.GetClass("DATE")
	require.NoError(t, err)

	instanceInterface, err := env.NewObjectInstance("DATE")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok, "Should return ObjectInstance")

	// Set a known date for testing
	knownDate := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	instance.NativeData = knownDate

	formatMethod, exists := dateClass.PublicFunctions["FORMAT"]
	require.True(t, exists, "FORMAT method should exist")

	tests := []struct {
		name     string
		layout   string
		expected string
	}{
		{
			"RFC3339 format",
			"2006-01-02T15:04:05Z07:00",
			"2023-12-25T15:30:45Z",
		},
		{
			"Simple format",
			"2006-01-02",
			"2023-12-25",
		},
		{
			"Time format",
			"15:04:05",
			"15:30:45",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{types.StringValue(test.layout)}
			result, err := formatMethod.NativeImpl(nil, instance, args)
			require.NoError(t, err)

			stringResult, ok := result.(types.StringValue)
			require.True(t, ok, "FORMAT should return StringValue")
			assert.Equal(t, test.expected, string(stringResult))
		})
	}
}

func TestTimeDATEErrorHandling(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterTIMEInEnv(env)

	dateClass, err := env.GetClass("DATE")
	require.NoError(t, err)

	instanceInterface, err := env.NewObjectInstance("DATE")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok, "Should return ObjectInstance")

	// Don't initialize the NativeData, should cause errors

	methods := []string{"YEAR", "MONTH", "DAY", "HOUR", "MINUTE", "SECOND", "MILLISECOND", "NANOSECOND"}

	for _, methodName := range methods {
		t.Run(methodName+" without init", func(t *testing.T) {
			method, exists := dateClass.PublicFunctions[methodName]
			require.True(t, exists)

			_, err := method.NativeImpl(nil, instance, []types.Value{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid context")
		})
	}
}

func TestTimeSLEEPFunction(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterTIMEInEnv(env)

	sleepFunc, err := env.GetFunction("SLEEP")
	require.NoError(t, err)

	// Test short sleep (1 second would be too long for tests)
	// We'll use 0 seconds to just test the function works
	before := time.Now()
	args := []types.Value{types.IntegerValue(0)}
	result, err := sleepFunc.NativeImpl(nil, nil, args)
	after := time.Now()

	require.NoError(t, err)
	assert.Equal(t, types.NOTHIN, result)

	// Should complete almost instantly for 0 seconds
	duration := after.Sub(before)
	assert.Less(t, duration, 100*time.Millisecond, "0 second sleep should be nearly instant")
}

func TestTimeSLEEPErrorHandling(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterTIMEInEnv(env)

	sleepFunc, err := env.GetFunction("SLEEP")
	require.NoError(t, err)

	// Test with invalid argument type
	args := []types.Value{types.StringValue("not a number")}
	_, err = sleepFunc.NativeImpl(nil, nil, args)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SLEEP: invalid argument type")
}

func TestTimeSelectiveImport(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Import only SLEEP function
	err := RegisterTIMEInEnv(env, "SLEEP")
	require.NoError(t, err)

	// SLEEP should be available
	_, err = env.GetFunction("SLEEP")
	assert.NoError(t, err)

	// DATE class should not be available
	_, err = env.GetClass("DATE")
	assert.Error(t, err)
}

func TestTimeSelectiveImportClass(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Import only DATE class
	err := RegisterTIMEInEnv(env, "DATE")
	require.NoError(t, err)

	// DATE class should be available
	_, err = env.GetClass("DATE")
	assert.NoError(t, err)

	// SLEEP function should not be available
	_, err = env.GetFunction("SLEEP")
	assert.Error(t, err)
}

func TestTimeInvalidDeclaration(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Try to import non-existent declaration
	err := RegisterTIMEInEnv(env, "INVALID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown TIME declaration: INVALID")
}

func TestTimeCaseInsensitive(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Import using lowercase
	err := RegisterTIMEInEnv(env, "date", "sleep")
	require.NoError(t, err)

	// Class and function should be available (stored in uppercase)
	_, err = env.GetClass("DATE")
	assert.NoError(t, err)

	_, err = env.GetFunction("SLEEP")
	assert.NoError(t, err)
}
