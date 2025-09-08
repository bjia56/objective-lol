package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterArraysInEnv(t *testing.T) {
	env := environment.NewEnvironment(nil)

	RegisterArraysInEnv(env)

	// Test that BUKKIT class is registered
	_, err := env.GetClass("BUKKIT")
	assert.NoError(t, err, "BUKKIT class should be registered")
}

func TestBUKKITConstructor(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	bukkitClass, err := env.GetClass("BUKKIT")
	require.NoError(t, err)

	// Create a BUKKIT instance
	instance, err := env.NewObjectInstance("BUKKIT")
	require.NoError(t, err)

	// Initialize the bukkit (call constructor)
	constructor, exists := bukkitClass.PublicFunctions["BUKKIT"]
	require.True(t, exists, "BUKKIT constructor should exist")

	result, err := constructor.NativeImpl(nil, instance, BukkitSlice{})
	require.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)

	// Check that NativeData is initialized as empty slice
	slice, ok := instance.NativeData.(BukkitSlice)
	require.True(t, ok, "NativeData should be BukkitSlice")
	assert.Equal(t, 0, len(slice), "Initial slice should be empty")

	// Check that SIZ variable is initialized
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists, "SIZ variable should be initialized")
	val, err := sizVar.Get(instance)
	assert.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(0), val)
	assert.True(t, sizVar.IsLocked, "SIZ should be locked")
	assert.True(t, sizVar.IsPublic, "SIZ should be public")
}

func TestBUKKITPUSH(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)

	bukkitClass, _ := env.GetClass("BUKKIT")
	pushFunc := bukkitClass.PublicFunctions["PUSH"]

	// Test pushing different environment of values
	tests := []struct {
		name     string
		value    environment.Value
		expected int
	}{
		{"Push integer", environment.IntegerValue(42), 1},
		{"Push string", environment.StringValue("hello"), 2},
		{"Push boolean", environment.YEZ, 3},
		{"Push double", environment.DoubleValue(3.14), 4},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := BukkitSlice{test.value}
			result, err := pushFunc.NativeImpl(nil, instance, args)
			require.NoError(t, err)

			// PUSH should return new size
			assert.Equal(t, environment.IntegerValue(test.expected), result)

			// Check slice was updated
			slice := instance.NativeData.(BukkitSlice)
			assert.Equal(t, test.expected, len(slice))
			assert.Equal(t, test.value, slice[i])

			// Check SIZ variable was updated
			sizVar := instance.Variables["SIZ"]
			val, err := sizVar.Get(instance)
			assert.NoError(t, err)
			assert.Equal(t, environment.IntegerValue(test.expected), val)
		})
	}
}

func TestBUKKITPOP(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	popFunc := bukkitClass.PublicFunctions["POP"]

	// Add some values first
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(1)})
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(2)})
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(3)})

	// Test popping values
	result, err := popFunc.NativeImpl(nil, instance, BukkitSlice{})
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(3), result, "Should pop last element")

	// Check size updated
	slice := instance.NativeData.(BukkitSlice)
	assert.Equal(t, 2, len(slice))
	val, err := instance.Variables["SIZ"].Get(instance)
	assert.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(2), val)

	// Test error case: pop from empty array
	popFunc.NativeImpl(nil, instance, BukkitSlice{}) // Pop 2
	popFunc.NativeImpl(nil, instance, BukkitSlice{}) // Pop 1

	// Now should error
	_, err = popFunc.NativeImpl(nil, instance, BukkitSlice{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot pop from empty array")
}

func TestBUKKITAT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	atFunc := bukkitClass.PublicFunctions["AT"]

	// Add some values
	values := BukkitSlice{
		environment.StringValue("first"),
		environment.IntegerValue(42),
		environment.DoubleValue(3.14),
	}

	for _, val := range values {
		pushFunc.NativeImpl(nil, instance, BukkitSlice{val})
	}

	// Test accessing valid indices
	for i, expectedVal := range values {
		result, err := atFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(i)})
		require.NoError(t, err)
		assert.Equal(t, expectedVal, result)
	}

	// Test error cases
	tests := []struct {
		name  string
		index environment.Value
	}{
		{"Negative index", environment.IntegerValue(-1)},
		{"Index too large", environment.IntegerValue(10)},
		{"Invalid type", environment.StringValue("not a number")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := atFunc.NativeImpl(nil, instance, BukkitSlice{test.index})
			assert.Error(t, err)
		})
	}
}

func TestBUKKITSET(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	setFunc := bukkitClass.PublicFunctions["SET"]
	atFunc := bukkitClass.PublicFunctions["AT"]

	// Add some values
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(1)})
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(2)})
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(3)})

	// Test setting value at valid index
	result, err := setFunc.NativeImpl(nil, instance, BukkitSlice{
		environment.IntegerValue(1),
		environment.StringValue("changed"),
	})
	require.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)

	// Verify value was changed
	result, err = atFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(1)})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("changed"), result)

	// Test error case: index out of bounds
	_, err = setFunc.NativeImpl(nil, instance, BukkitSlice{
		environment.IntegerValue(10),
		environment.StringValue("invalid"),
	})
	assert.Error(t, err)
}

func TestBUKKITSHIFTUNSHIFT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	unshiftFunc := bukkitClass.PublicFunctions["UNSHIFT"]
	shiftFunc := bukkitClass.PublicFunctions["SHIFT"]

	// Test unshift (add to front)
	result, err := unshiftFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(1)})
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(1), result, "UNSHIFT should return new size")

	result, err = unshiftFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(2)})
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(2), result)

	// Array should be [2, 1]
	slice := instance.NativeData.(BukkitSlice)
	assert.Equal(t, environment.IntegerValue(2), slice[0])
	assert.Equal(t, environment.IntegerValue(1), slice[1])

	// Test shift (remove from front)
	result, err = shiftFunc.NativeImpl(nil, instance, BukkitSlice{})
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(2), result, "Should shift first element")

	// Array should be [1]
	slice = instance.NativeData.(BukkitSlice)
	assert.Equal(t, 1, len(slice))
	assert.Equal(t, environment.IntegerValue(1), slice[0])

	// Test error case: shift from empty array
	shiftFunc.NativeImpl(nil, instance, BukkitSlice{}) // Remove last element

	_, err = shiftFunc.NativeImpl(nil, instance, BukkitSlice{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot shift from empty array")
}

func TestBUKKITCLEAR(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	clearFunc := bukkitClass.PublicFunctions["CLEAR"]

	// Add some values
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(1)})
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(2)})

	// Clear the array
	result, err := clearFunc.NativeImpl(nil, instance, BukkitSlice{})
	require.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)

	// Check array is empty
	slice := instance.NativeData.(BukkitSlice)
	assert.Equal(t, 0, len(slice))

	val, err := instance.Variables["SIZ"].Get(instance)
	assert.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(0), val)
}

func TestBUKKITREVERSE(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	reverseFunc := bukkitClass.PublicFunctions["REVERSE"]

	// Add values [1, 2, 3]
	values := BukkitSlice{
		environment.IntegerValue(1),
		environment.IntegerValue(2),
		environment.IntegerValue(3),
	}
	for _, val := range values {
		pushFunc.NativeImpl(nil, instance, BukkitSlice{val})
	}

	// Reverse the array
	result, err := reverseFunc.NativeImpl(nil, instance, BukkitSlice{})
	require.NoError(t, err)
	assert.Equal(t, instance, result)

	// Check array is reversed [3, 2, 1]
	slice := instance.NativeData.(BukkitSlice)
	assert.Equal(t, environment.IntegerValue(3), slice[0])
	assert.Equal(t, environment.IntegerValue(2), slice[1])
	assert.Equal(t, environment.IntegerValue(1), slice[2])
}

func TestBUKKITSORT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	sortFunc := bukkitClass.PublicFunctions["SORT"]

	// Add unsorted values
	values := BukkitSlice{
		environment.IntegerValue(3),
		environment.IntegerValue(1),
		environment.DoubleValue(2.5),
		environment.IntegerValue(2),
	}
	for _, val := range values {
		pushFunc.NativeImpl(nil, instance, BukkitSlice{val})
	}

	// Sort the array
	result, err := sortFunc.NativeImpl(nil, instance, BukkitSlice{})
	require.NoError(t, err)
	assert.Equal(t, instance, result)

	// Check array is sorted
	slice := instance.NativeData.(BukkitSlice)
	assert.Equal(t, environment.IntegerValue(1), slice[0])
	assert.Equal(t, environment.IntegerValue(2), slice[1])
	assert.Equal(t, environment.DoubleValue(2.5), slice[2])
	assert.Equal(t, environment.IntegerValue(3), slice[3])
}

func TestBUKKITJOIN(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	joinFunc := bukkitClass.PublicFunctions["JOIN"]

	// Test empty array
	result, err := joinFunc.NativeImpl(nil, instance, BukkitSlice{environment.StringValue(",")})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue(""), result)

	// Add values
	values := BukkitSlice{
		environment.StringValue("hello"),
		environment.IntegerValue(42),
		environment.DoubleValue(3.14),
	}
	for _, val := range values {
		pushFunc.NativeImpl(nil, instance, BukkitSlice{val})
	}

	// Test join with comma
	result, err = joinFunc.NativeImpl(nil, instance, BukkitSlice{environment.StringValue(",")})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("hello,42,3.14"), result)

	// Test join with space
	result, err = joinFunc.NativeImpl(nil, instance, BukkitSlice{environment.StringValue(" ")})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("hello 42 3.14"), result)
}

func TestBUKKITFIND(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	findFunc := bukkitClass.PublicFunctions["FIND"]

	// Add values
	values := BukkitSlice{
		environment.StringValue("hello"),
		environment.IntegerValue(42),
		environment.DoubleValue(3.14),
	}
	for _, val := range values {
		pushFunc.NativeImpl(nil, instance, BukkitSlice{val})
	}

	// Test finding existing values
	result, err := findFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(42)})
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(1), result)

	result, err = findFunc.NativeImpl(nil, instance, BukkitSlice{environment.StringValue("hello")})
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(0), result)

	// Test finding non-existing value
	result, err = findFunc.NativeImpl(nil, instance, BukkitSlice{environment.StringValue("not found")})
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(-1), result)
}

func TestBUKKITCONTAINS(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	containsFunc := bukkitClass.PublicFunctions["CONTAINS"]

	// Add values
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(42)})
	pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.StringValue("test")})

	// Test contains existing value
	result, err := containsFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(42)})
	require.NoError(t, err)
	assert.Equal(t, environment.YEZ, result)

	// Test contains non-existing value
	result, err = containsFunc.NativeImpl(nil, instance, BukkitSlice{environment.StringValue("not found")})
	require.NoError(t, err)
	assert.Equal(t, environment.NO, result)
}

func TestBUKKITSLICE(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArraysInEnv(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	sliceFunc := bukkitClass.PublicFunctions["SLICE"]

	// Add values [0, 1, 2, 3, 4]
	for i := 0; i < 5; i++ {
		pushFunc.NativeImpl(nil, instance, BukkitSlice{environment.IntegerValue(i)})
	}

	// Test normal slice [1:4] -> [1, 2, 3]
	result, err := sliceFunc.NativeImpl(nil, instance, BukkitSlice{
		environment.IntegerValue(1),
		environment.IntegerValue(4),
	})
	require.NoError(t, err)

	// Result should be an ObjectInstance containing a new BUKKIT
	objVal, ok := result.(*environment.ObjectInstance)
	require.True(t, ok, "SLICE should return ObjectInstance")

	newSlice := objVal.NativeData.(BukkitSlice)

	assert.Equal(t, 3, len(newSlice))
	assert.Equal(t, environment.IntegerValue(1), newSlice[0])
	assert.Equal(t, environment.IntegerValue(2), newSlice[1])
	assert.Equal(t, environment.IntegerValue(3), newSlice[2])

	// Check SIZ variable
	val, err := objVal.Variables["SIZ"].Get(objVal)
	assert.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(3), val)

	// Test error case: invalid bounds
	_, err = sliceFunc.NativeImpl(nil, instance, BukkitSlice{
		environment.IntegerValue(10),
		environment.IntegerValue(15),
	})
	assert.Error(t, err)
}

// Helper function to create and initialize a BUKKIT instance
func createBUKKITInstance(t *testing.T, env *environment.Environment) *environment.ObjectInstance {
	bukkitClass, err := env.GetClass("BUKKIT")
	require.NoError(t, err)

	instance, err := env.NewObjectInstance("BUKKIT")
	require.NoError(t, err)

	// Initialize with constructor
	constructor := bukkitClass.PublicFunctions["BUKKIT"]
	_, err = constructor.NativeImpl(nil, instance, BukkitSlice{})
	require.NoError(t, err)

	return instance
}
