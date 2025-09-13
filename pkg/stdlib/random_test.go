package stdlib

import (
	"math/rand"
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterRANDOM(t *testing.T) {
	env := environment.NewEnvironment(nil)

	err := RegisterRANDOMInEnv(env) // Empty slice imports all
	require.NoError(t, err)

	// Test that RANDOM functions are registered
	randomFunctions := []string{"SEED", "SEED_TIME", "RANDOM_FLOAT", "RANDOM_RANGE", "RANDOM_INT", "RANDOM_BOOL", "RANDOM_CHOICE", "SHUFFLE", "RANDOM_STRING", "UUID"}

	for _, funcName := range randomFunctions {
		_, err := env.GetFunction(funcName)
		assert.NoError(t, err, "Function %s should be registered", funcName)
	}
}

func TestRANDOMSEED(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)

	seedFunc, err := env.GetFunction("SEED")
	require.NoError(t, err)

	// Test seeding with a specific value
	args := []environment.Value{environment.IntegerValue(42)}
	result, err := seedFunc.NativeImpl(nil, nil, args)
	require.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)

	// Test error case: invalid seed type
	args = []environment.Value{environment.StringValue("invalid")}
	_, err = seedFunc.NativeImpl(nil, nil, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SEED: invalid seed type")
}

func TestRANDOMSEED_TIME(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)

	seedTimeFunc, err := env.GetFunction("SEED_TIME")
	require.NoError(t, err)

	args := []environment.Value{}
	result, err := seedTimeFunc.NativeImpl(nil, nil, args)
	require.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)
}

func TestRANDOMFLOAT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)

	randomFloatFunc, err := env.GetFunction("RANDOM_FLOAT")
	require.NoError(t, err)

	// Test that RANDOM_FLOAT returns values in expected range
	for i := 0; i < 100; i++ {
		args := []environment.Value{}
		result, err := randomFloatFunc.NativeImpl(nil, nil, args)
		require.NoError(t, err)

		doubleResult, ok := result.(environment.DoubleValue)
		require.True(t, ok, "RANDOM_FLOAT should return DoubleValue")

		val := float64(doubleResult)
		assert.GreaterOrEqual(t, val, 0.0, "RANDOM_FLOAT should be >= 0")
		assert.Less(t, val, 1.0, "RANDOM_FLOAT should be < 1")
	}
}

func TestRANDOMRANGE(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)

	randomRangeFunc, err := env.GetFunction("RANDOM_RANGE")
	require.NoError(t, err)

	// Test that RANDOM_RANGE returns values in expected range
	for i := 0; i < 50; i++ {
		args := []environment.Value{environment.DoubleValue(5.0), environment.DoubleValue(10.0)}
		result, err := randomRangeFunc.NativeImpl(nil, nil, args)
		require.NoError(t, err)

		doubleResult, ok := result.(environment.DoubleValue)
		require.True(t, ok, "RANDOM_RANGE should return DoubleValue")

		val := float64(doubleResult)
		assert.GreaterOrEqual(t, val, 5.0, "RANDOM_RANGE should be >= min")
		assert.Less(t, val, 10.0, "RANDOM_RANGE should be < max")
	}

	// Test error case: min >= max
	args := []environment.Value{environment.DoubleValue(10.0), environment.DoubleValue(5.0)}
	_, err = randomRangeFunc.NativeImpl(nil, nil, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RANDOM_RANGE: min must be less than max")
}

func TestRANDOMINT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)

	randomIntFunc, err := env.GetFunction("RANDOM_INT")
	require.NoError(t, err)

	// Test that RANDOM_INT returns values in expected range
	for i := 0; i < 50; i++ {
		args := []environment.Value{environment.IntegerValue(1), environment.IntegerValue(10)}
		result, err := randomIntFunc.NativeImpl(nil, nil, args)
		require.NoError(t, err)

		intResult, ok := result.(environment.IntegerValue)
		require.True(t, ok, "RANDOM_INT should return IntegerValue")

		val := int64(intResult)
		assert.GreaterOrEqual(t, val, int64(1), "RANDOM_INT should be >= min")
		assert.Less(t, val, int64(10), "RANDOM_INT should be < max")
	}

	// Test error case: min >= max
	args := []environment.Value{environment.IntegerValue(10), environment.IntegerValue(5)}
	_, err = randomIntFunc.NativeImpl(nil, nil, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RANDOM_INT: min must be less than max")
}

func TestRANDOMBOOL(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)

	randomBoolFunc, err := env.GetFunction("RANDOM_BOOL")
	require.NoError(t, err)

	// Test that RANDOM_BOOL returns boolean values
	trueCount := 0
	falseCount := 0
	for i := 0; i < 1000; i++ {
		args := []environment.Value{}
		result, err := randomBoolFunc.NativeImpl(nil, nil, args)
		require.NoError(t, err)

		boolResult, ok := result.(environment.BoolValue)
		require.True(t, ok, "RANDOM_BOOL should return BoolValue")

		if boolResult == environment.YEZ {
			trueCount++
		} else {
			falseCount++
		}
	}

	// Both true and false should appear (with high probability)
	assert.Greater(t, trueCount, 0, "Should have some true values")
	assert.Greater(t, falseCount, 0, "Should have some false values")
}

func TestRANDOMCHOICE(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)
	RegisterArraysInEnv(env)

	randomChoiceFunc, err := env.GetFunction("RANDOM_CHOICE")
	require.NoError(t, err)

	// Create a test array
	bukkit := NewBukkitInstance()
	testValues := []environment.Value{
		environment.StringValue("apple"),
		environment.StringValue("banana"),
		environment.StringValue("cherry"),
	}
	bukkit.NativeData = BukkitSlice(testValues)

	// Test RANDOM_CHOICE
	args := []environment.Value{bukkit}
	result, err := randomChoiceFunc.NativeImpl(nil, nil, args)
	require.NoError(t, err)

	// Result should be one of the values in the array
	strResult, ok := result.(environment.StringValue)
	require.True(t, ok, "RANDOM_CHOICE should return a string from the array")

	found := false
	for _, val := range testValues {
		if strVal, ok := val.(environment.StringValue); ok && strVal == strResult {
			found = true
			break
		}
	}
	assert.True(t, found, "RANDOM_CHOICE should return a value from the array")

	// Test error case: empty array
	emptyBukkit := NewBukkitInstance()
	emptyBukkit.NativeData = BukkitSlice{}
	args = []environment.Value{emptyBukkit}
	_, err = randomChoiceFunc.NativeImpl(nil, nil, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RANDOM_CHOICE: empty array")
}

func TestRANDOMSTRING(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)

	randomStringFunc, err := env.GetFunction("RANDOM_STRING")
	require.NoError(t, err)

	// Test generating random string
	args := []environment.Value{
		environment.IntegerValue(10),
		environment.StringValue("abc"),
	}
	result, err := randomStringFunc.NativeImpl(nil, nil, args)
	require.NoError(t, err)

	strResult, ok := result.(environment.StringValue)
	require.True(t, ok, "RANDOM_STRING should return StringValue")

	str := string(strResult)
	assert.Equal(t, 10, len(str), "String should have requested length")

	// All characters should be from the charset
	for _, char := range str {
		assert.Contains(t, "abc", string(char), "All characters should be from charset")
	}

	// Test edge cases
	// Zero length
	args = []environment.Value{
		environment.IntegerValue(0),
		environment.StringValue("abc"),
	}
	result, err = randomStringFunc.NativeImpl(nil, nil, args)
	require.NoError(t, err)
	assert.Equal(t, "", string(result.(environment.StringValue)))

	// Empty charset error
	args = []environment.Value{
		environment.IntegerValue(5),
		environment.StringValue(""),
	}
	_, err = randomStringFunc.NativeImpl(nil, nil, args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RANDOM_STRING: empty charset")
}

func TestUUID(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)

	uuidFunc, err := env.GetFunction("UUID")
	require.NoError(t, err)

	// Generate multiple UUIDs and check format
	uuids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		args := []environment.Value{}
		result, err := uuidFunc.NativeImpl(nil, nil, args)
		require.NoError(t, err)

		uuidResult, ok := result.(environment.StringValue)
		require.True(t, ok, "UUID should return StringValue")

		uuid := string(uuidResult)

		// Check UUID format (8-4-4-4-12 hexadecimal digits)
		assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, uuid)

		// UUIDs should be unique
		assert.False(t, uuids[uuid], "UUIDs should be unique")
		uuids[uuid] = true
	}
}

func TestSHUFFLE(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterRANDOMInEnv(env)
	RegisterArraysInEnv(env)

	shuffleFunc, err := env.GetFunction("SHUFFLE")
	require.NoError(t, err)

	// Create a test array
	bukkit := NewBukkitInstance()
	testValues := []environment.Value{
		environment.IntegerValue(1),
		environment.IntegerValue(2),
		environment.IntegerValue(3),
		environment.IntegerValue(4),
		environment.IntegerValue(5),
	}
	bukkit.NativeData = BukkitSlice(testValues)

	// Seed for reproducible tests
	rand.Seed(42)

	// Test SHUFFLE
	args := []environment.Value{bukkit}
	result, err := shuffleFunc.NativeImpl(nil, nil, args)
	require.NoError(t, err)

	// Result should be a BUKKIT
	shuffledBukkit, ok := result.(*environment.ObjectInstance)
	require.True(t, ok, "SHUFFLE should return ObjectInstance")

	shuffledSlice, ok := shuffledBukkit.NativeData.(BukkitSlice)
	require.True(t, ok, "SHUFFLE should return BukkitSlice")

	// Should have same length
	assert.Equal(t, len(testValues), len(shuffledSlice), "Shuffled array should have same length")

	// Should contain all original values
	for _, originalVal := range testValues {
		found := false
		for _, shuffledVal := range shuffledSlice {
			if originalVal == shuffledVal {
				found = true
				break
			}
		}
		assert.True(t, found, "Shuffled array should contain all original values")
	}

	// Original array should be unchanged
	originalSlice, ok := bukkit.NativeData.(BukkitSlice)
	require.True(t, ok)
	for i, val := range originalSlice {
		assert.Equal(t, testValues[i], val, "Original array should be unchanged")
	}
}

func TestRANDOMSelectiveImport(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test importing only RANDOM_FLOAT
	err := RegisterRANDOMInEnv(env, "RANDOM_FLOAT")
	require.NoError(t, err)

	// Should have RANDOM_FLOAT
	_, err = env.GetFunction("RANDOM_FLOAT")
	assert.NoError(t, err, "RANDOM_FLOAT function should be registered")

	// Should not have UUID
	_, err = env.GetFunction("UUID")
	assert.Error(t, err, "UUID function should not be registered")
}

func TestRANDOMInvalidFunction(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test importing non-existent function
	err := RegisterRANDOMInEnv(env, "NONEXISTENT")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown RANDOM declaration: NONEXISTENT")
}

func TestRANDOMCaseInsensitive(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test case insensitive import
	err := RegisterRANDOMInEnv(env, "random_float", "uuid")
	require.NoError(t, err)

	// Both functions should be available
	_, err = env.GetFunction("RANDOM_FLOAT")
	assert.NoError(t, err, "RANDOM_FLOAT function should be registered with lowercase import")

	_, err = env.GetFunction("UUID")
	assert.NoError(t, err, "UUID function should be registered with lowercase import")
}
