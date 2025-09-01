package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterArrays(t *testing.T) {
	env := environment.NewEnvironment(nil)

	RegisterArrays(env)

	// Test that BUKKIT class is registered
	_, err := env.GetClass("BUKKIT")
	assert.NoError(t, err, "BUKKIT class should be registered")
}

func TestBUKKITConstructor(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	bukkitClass, err := env.GetClass("BUKKIT")
	require.NoError(t, err)

	// Create a BUKKIT instance
	instanceInterface, err := env.NewObjectInstance("BUKKIT")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok, "Should return ObjectInstance")

	// Initialize the bukkit (call constructor)
	constructor, exists := bukkitClass.PublicFunctions["BUKKIT"]
	require.True(t, exists, "BUKKIT constructor should exist")

	result, err := constructor.NativeImpl(instance, []types.Value{})
	require.NoError(t, err)
	assert.Equal(t, types.NOTHIN, result)

	// Check that NativeData is initialized as empty slice
	slice, ok := instance.NativeData.([]types.Value)
	require.True(t, ok, "NativeData should be []types.Value")
	assert.Equal(t, 0, len(slice), "Initial slice should be empty")

	// Check that SIZ variable is initialized
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists, "SIZ variable should be initialized")
	assert.Equal(t, types.IntegerValue(0), sizVar.Value)
	assert.True(t, sizVar.IsLocked, "SIZ should be locked")
	assert.True(t, sizVar.IsPublic, "SIZ should be public")
}

func TestBUKKITPUSH(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)

	bukkitClass, _ := env.GetClass("BUKKIT")
	pushFunc := bukkitClass.PublicFunctions["PUSH"]

	// Test pushing different types of values
	tests := []struct {
		name     string
		value    types.Value
		expected int
	}{
		{"Push integer", types.IntegerValue(42), 1},
		{"Push string", types.StringValue("hello"), 2},
		{"Push boolean", types.YEZ, 3},
		{"Push double", types.DoubleValue(3.14), 4},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []types.Value{test.value}
			result, err := pushFunc.NativeImpl(instance, args)
			require.NoError(t, err)

			// PUSH should return new size
			assert.Equal(t, types.IntegerValue(test.expected), result)

			// Check slice was updated
			slice := instance.NativeData.([]types.Value)
			assert.Equal(t, test.expected, len(slice))
			assert.Equal(t, test.value, slice[i])

			// Check SIZ variable was updated
			sizVar := instance.Variables["SIZ"]
			assert.Equal(t, types.IntegerValue(test.expected), sizVar.Value)
		})
	}
}

func TestBUKKITPOP(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	popFunc := bukkitClass.PublicFunctions["POP"]

	// Add some values first
	pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(1)})
	pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(2)})
	pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(3)})

	// Test popping values
	result, err := popFunc.NativeImpl(instance, []types.Value{})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(3), result, "Should pop last element")

	// Check size updated
	slice := instance.NativeData.([]types.Value)
	assert.Equal(t, 2, len(slice))
	assert.Equal(t, types.IntegerValue(2), instance.Variables["SIZ"].Value)

	// Test error case: pop from empty array
	popFunc.NativeImpl(instance, []types.Value{}) // Pop 2
	popFunc.NativeImpl(instance, []types.Value{}) // Pop 1

	// Now should error
	_, err = popFunc.NativeImpl(instance, []types.Value{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot pop from empty array")
}

func TestBUKKITAT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	atFunc := bukkitClass.PublicFunctions["AT"]

	// Add some values
	values := []types.Value{
		types.StringValue("first"),
		types.IntegerValue(42),
		types.DoubleValue(3.14),
	}

	for _, val := range values {
		pushFunc.NativeImpl(instance, []types.Value{val})
	}

	// Test accessing valid indices
	for i, expectedVal := range values {
		result, err := atFunc.NativeImpl(instance, []types.Value{types.IntegerValue(i)})
		require.NoError(t, err)
		assert.Equal(t, expectedVal, result)
	}

	// Test error cases
	tests := []struct {
		name  string
		index types.Value
	}{
		{"Negative index", types.IntegerValue(-1)},
		{"Index too large", types.IntegerValue(10)},
		{"Invalid type", types.StringValue("not a number")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := atFunc.NativeImpl(instance, []types.Value{test.index})
			assert.Error(t, err)
		})
	}
}

func TestBUKKITSET(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	setFunc := bukkitClass.PublicFunctions["SET"]
	atFunc := bukkitClass.PublicFunctions["AT"]

	// Add some values
	pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(1)})
	pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(2)})
	pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(3)})

	// Test setting value at valid index
	result, err := setFunc.NativeImpl(instance, []types.Value{
		types.IntegerValue(1),
		types.StringValue("changed"),
	})
	require.NoError(t, err)
	assert.Equal(t, types.NOTHIN, result)

	// Verify value was changed
	result, err = atFunc.NativeImpl(instance, []types.Value{types.IntegerValue(1)})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("changed"), result)

	// Test error case: index out of bounds
	_, err = setFunc.NativeImpl(instance, []types.Value{
		types.IntegerValue(10),
		types.StringValue("invalid"),
	})
	assert.Error(t, err)
}

func TestBUKKITSHIFTUNSHIFT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	unshiftFunc := bukkitClass.PublicFunctions["UNSHIFT"]
	shiftFunc := bukkitClass.PublicFunctions["SHIFT"]

	// Test unshift (add to front)
	result, err := unshiftFunc.NativeImpl(instance, []types.Value{types.IntegerValue(1)})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(1), result, "UNSHIFT should return new size")

	result, err = unshiftFunc.NativeImpl(instance, []types.Value{types.IntegerValue(2)})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(2), result)

	// Array should be [2, 1]
	slice := instance.NativeData.([]types.Value)
	assert.Equal(t, types.IntegerValue(2), slice[0])
	assert.Equal(t, types.IntegerValue(1), slice[1])

	// Test shift (remove from front)
	result, err = shiftFunc.NativeImpl(instance, []types.Value{})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(2), result, "Should shift first element")

	// Array should be [1]
	slice = instance.NativeData.([]types.Value)
	assert.Equal(t, 1, len(slice))
	assert.Equal(t, types.IntegerValue(1), slice[0])

	// Test error case: shift from empty array
	shiftFunc.NativeImpl(instance, []types.Value{}) // Remove last element

	_, err = shiftFunc.NativeImpl(instance, []types.Value{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot shift from empty array")
}

func TestBUKKITCLEAR(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	clearFunc := bukkitClass.PublicFunctions["CLEAR"]

	// Add some values
	pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(1)})
	pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(2)})

	// Clear the array
	result, err := clearFunc.NativeImpl(instance, []types.Value{})
	require.NoError(t, err)
	assert.Equal(t, types.NOTHIN, result)

	// Check array is empty
	slice := instance.NativeData.([]types.Value)
	assert.Equal(t, 0, len(slice))
	assert.Equal(t, types.IntegerValue(0), instance.Variables["SIZ"].Value)
}

func TestBUKKITREVERSE(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	reverseFunc := bukkitClass.PublicFunctions["REVERSE"]

	// Add values [1, 2, 3]
	values := []types.Value{
		types.IntegerValue(1),
		types.IntegerValue(2),
		types.IntegerValue(3),
	}
	for _, val := range values {
		pushFunc.NativeImpl(instance, []types.Value{val})
	}

	// Reverse the array
	result, err := reverseFunc.NativeImpl(instance, []types.Value{})
	require.NoError(t, err)
	assert.Equal(t, types.NOTHIN, result)

	// Check array is reversed [3, 2, 1]
	slice := instance.NativeData.([]types.Value)
	assert.Equal(t, types.IntegerValue(3), slice[0])
	assert.Equal(t, types.IntegerValue(2), slice[1])
	assert.Equal(t, types.IntegerValue(1), slice[2])
}

func TestBUKKITSORT(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	sortFunc := bukkitClass.PublicFunctions["SORT"]

	// Add unsorted values
	values := []types.Value{
		types.IntegerValue(3),
		types.IntegerValue(1),
		types.DoubleValue(2.5),
		types.IntegerValue(2),
	}
	for _, val := range values {
		pushFunc.NativeImpl(instance, []types.Value{val})
	}

	// Sort the array
	result, err := sortFunc.NativeImpl(instance, []types.Value{})
	require.NoError(t, err)
	assert.Equal(t, types.NOTHIN, result)

	// Check array is sorted
	slice := instance.NativeData.([]types.Value)
	assert.Equal(t, types.IntegerValue(1), slice[0])
	assert.Equal(t, types.IntegerValue(2), slice[1])
	assert.Equal(t, types.DoubleValue(2.5), slice[2])
	assert.Equal(t, types.IntegerValue(3), slice[3])
}

func TestBUKKITJOIN(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	joinFunc := bukkitClass.PublicFunctions["JOIN"]

	// Test empty array
	result, err := joinFunc.NativeImpl(instance, []types.Value{types.StringValue(",")})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue(""), result)

	// Add values
	values := []types.Value{
		types.StringValue("hello"),
		types.IntegerValue(42),
		types.DoubleValue(3.14),
	}
	for _, val := range values {
		pushFunc.NativeImpl(instance, []types.Value{val})
	}

	// Test join with comma
	result, err = joinFunc.NativeImpl(instance, []types.Value{types.StringValue(",")})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("hello,42,3.14"), result)

	// Test join with space
	result, err = joinFunc.NativeImpl(instance, []types.Value{types.StringValue(" ")})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("hello 42 3.14"), result)
}

func TestBUKKITFIND(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	findFunc := bukkitClass.PublicFunctions["FIND"]

	// Add values
	values := []types.Value{
		types.StringValue("hello"),
		types.IntegerValue(42),
		types.DoubleValue(3.14),
	}
	for _, val := range values {
		pushFunc.NativeImpl(instance, []types.Value{val})
	}

	// Test finding existing values
	result, err := findFunc.NativeImpl(instance, []types.Value{types.IntegerValue(42)})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(1), result)

	result, err = findFunc.NativeImpl(instance, []types.Value{types.StringValue("hello")})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(0), result)

	// Test finding non-existing value
	result, err = findFunc.NativeImpl(instance, []types.Value{types.StringValue("not found")})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(-1), result)
}

func TestBUKKITCONTAINS(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	containsFunc := bukkitClass.PublicFunctions["CONTAINS"]

	// Add values
	pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(42)})
	pushFunc.NativeImpl(instance, []types.Value{types.StringValue("test")})

	// Test contains existing value
	result, err := containsFunc.NativeImpl(instance, []types.Value{types.IntegerValue(42)})
	require.NoError(t, err)
	assert.Equal(t, types.YEZ, result)

	// Test contains non-existing value
	result, err = containsFunc.NativeImpl(instance, []types.Value{types.StringValue("not found")})
	require.NoError(t, err)
	assert.Equal(t, types.NO, result)
}

func TestBUKKITSLICE(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterArrays(env)

	instance := createBUKKITInstance(t, env)
	bukkitClass, _ := env.GetClass("BUKKIT")

	pushFunc := bukkitClass.PublicFunctions["PUSH"]
	sliceFunc := bukkitClass.PublicFunctions["SLICE"]

	// Add values [0, 1, 2, 3, 4]
	for i := 0; i < 5; i++ {
		pushFunc.NativeImpl(instance, []types.Value{types.IntegerValue(i)})
	}

	// Test normal slice [1:4] -> [1, 2, 3]
	result, err := sliceFunc.NativeImpl(instance, []types.Value{
		types.IntegerValue(1),
		types.IntegerValue(4),
	})
	require.NoError(t, err)

	// Result should be an ObjectValue containing a new BUKKIT
	objVal, ok := result.(types.ObjectValue)
	require.True(t, ok, "SLICE should return ObjectValue")

	newInstance := objVal.Instance.(*environment.ObjectInstance)
	newSlice := newInstance.NativeData.([]types.Value)

	assert.Equal(t, 3, len(newSlice))
	assert.Equal(t, types.IntegerValue(1), newSlice[0])
	assert.Equal(t, types.IntegerValue(2), newSlice[1])
	assert.Equal(t, types.IntegerValue(3), newSlice[2])

	// Check SIZ variable
	assert.Equal(t, types.IntegerValue(3), newInstance.Variables["SIZ"].Value)

	// Test error case: invalid bounds
	_, err = sliceFunc.NativeImpl(instance, []types.Value{
		types.IntegerValue(10),
		types.IntegerValue(15),
	})
	assert.Error(t, err)
}

func TestUpdateSIZHelper(t *testing.T) {
	// Create a mock object instance with SIZ variable
	instance := &environment.ObjectInstance{
		Variables: map[string]*environment.Variable{
			"SIZ": {
				Name:  "SIZ",
				Type:  "INTEGR",
				Value: types.IntegerValue(0),
			},
		},
	}

	// Test updating SIZ
	slice := []types.Value{
		types.IntegerValue(1),
		types.IntegerValue(2),
		types.IntegerValue(3),
	}

	updateSIZ(instance, slice)

	assert.Equal(t, types.IntegerValue(3), instance.Variables["SIZ"].Value)
}

// Helper function to create and initialize a BUKKIT instance
func createBUKKITInstance(t *testing.T, env *environment.Environment) *environment.ObjectInstance {
	bukkitClass, err := env.GetClass("BUKKIT")
	require.NoError(t, err)

	instanceInterface, err := env.NewObjectInstance("BUKKIT")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok, "Should return ObjectInstance")

	// Initialize with constructor
	constructor := bukkitClass.PublicFunctions["BUKKIT"]
	_, err = constructor.NativeImpl(instance, []types.Value{})
	require.NoError(t, err)

	return instance
}
