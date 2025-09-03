package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaskitConstructor(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	constructor := baskitClass.PublicFunctions["BASKIT"]

	// Create a new BASKIT instance
	instance := NewBaskitInstance()

	// Call constructor
	_, err := constructor.NativeImpl(nil, instance, []types.Value{})
	require.NoError(t, err)

	// Verify empty map
	baskitMap, ok := instance.NativeData.(BaskitMap)
	require.True(t, ok, "NativeData should be BaskitMap")
	assert.Equal(t, 0, len(baskitMap))

	// Verify SIZ is 0
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists)
	assert.Equal(t, types.IntegerValue(0), sizVar.Value)
}

func TestBaskitPutAndGet(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	getMethod := baskitClass.PublicFunctions["GET"]

	instance := NewBaskitInstance()

	// PUT a value
	_, err := putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
		types.IntegerValue(42),
	})
	require.NoError(t, err)

	// Verify SIZ updated
	sizVar := instance.Variables["SIZ"]
	assert.Equal(t, types.IntegerValue(1), sizVar.Value)

	// GET the value
	result, err := getMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(42), result)

	// PUT another value with different key type (should convert to string)
	_, err = putMethod.NativeImpl(nil, instance, []types.Value{
		types.IntegerValue(123),
		types.StringValue("numeric key"),
	})
	require.NoError(t, err)

	// GET with numeric key (converted to string)
	result, err = getMethod.NativeImpl(nil, instance, []types.Value{
		types.IntegerValue(123),
	})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("numeric key"), result)

	// Verify SIZ updated
	sizVar = instance.Variables["SIZ"]
	assert.Equal(t, types.IntegerValue(2), sizVar.Value)
}

func TestBaskitGetNonexistentKey(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	getMethod := baskitClass.PublicFunctions["GET"]

	instance := NewBaskitInstance()

	// Try to GET non-existent key - should throw exception
	_, err := getMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("nonexistent"),
	})
	require.Error(t, err)

	// Check it's the right type of exception
	exception, ok := err.(ast.Exception)
	require.True(t, ok, "Should be an ast.Exception")
	assert.Contains(t, exception.Message, "Key 'nonexistent' not found in BASKIT")
}

func TestBaskitContains(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	containsMethod := baskitClass.PublicFunctions["CONTAINS"]

	instance := NewBaskitInstance()

	// Check non-existent key
	result, err := containsMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, types.NO, result)

	// PUT a value
	_, err = putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
		types.IntegerValue(42),
	})
	require.NoError(t, err)

	// Check existing key
	result, err = containsMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, types.YEZ, result)
}

func TestBaskitRemove(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	removeMethod := baskitClass.PublicFunctions["REMOVE"]
	containsMethod := baskitClass.PublicFunctions["CONTAINS"]

	instance := NewBaskitInstance()

	// PUT a value
	_, err := putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
		types.StringValue("value1"),
	})
	require.NoError(t, err)

	// REMOVE the value
	result, err := removeMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("value1"), result)

	// Verify SIZ updated
	sizVar := instance.Variables["SIZ"]
	assert.Equal(t, types.IntegerValue(0), sizVar.Value)

	// Verify key no longer exists
	containsResult, err := containsMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, types.NO, containsResult)

	// Try to remove non-existent key - should throw exception
	_, err = removeMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("nonexistent"),
	})
	require.Error(t, err)

	// Check it's the right type of exception
	exception, ok := err.(ast.Exception)
	require.True(t, ok, "Should be an ast.Exception")
	assert.Contains(t, exception.Message, "Key 'nonexistent' not found in BASKIT")
}

func TestBaskitClear(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	clearMethod := baskitClass.PublicFunctions["CLEAR"]

	instance := NewBaskitInstance()

	// PUT multiple values
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
		types.IntegerValue(1),
	})
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key2"),
		types.IntegerValue(2),
	})
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key3"),
		types.IntegerValue(3),
	})

	// Verify SIZ before clear
	sizVar := instance.Variables["SIZ"]
	assert.Equal(t, types.IntegerValue(3), sizVar.Value)

	// CLEAR all values
	_, err := clearMethod.NativeImpl(nil, instance, []types.Value{})
	require.NoError(t, err)

	// Verify SIZ is 0
	sizVar = instance.Variables["SIZ"]
	assert.Equal(t, types.IntegerValue(0), sizVar.Value)

	// Verify map is empty
	baskitMap := instance.NativeData.(BaskitMap)
	assert.Equal(t, 0, len(baskitMap))
}

func TestBaskitKeys(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	keysMethod := baskitClass.PublicFunctions["KEYS"]

	instance := NewBaskitInstance()

	// PUT multiple values
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key2"),
		types.IntegerValue(2),
	})
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
		types.IntegerValue(1),
	})
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key3"),
		types.IntegerValue(3),
	})

	// Get KEYS
	result, err := keysMethod.NativeImpl(nil, instance, []types.Value{})
	require.NoError(t, err)

	// Should return a BUKKIT
	bukkitValue, ok := result.(types.ObjectValue)
	require.True(t, ok, "KEYS should return ObjectValue")
	assert.Equal(t, "BUKKIT", bukkitValue.ClassName)

	// Get the underlying slice
	bukkitInstance, ok := bukkitValue.Instance.(*environment.ObjectInstance)
	require.True(t, ok)
	slice, ok := bukkitInstance.NativeData.(BukkitSlice)
	require.True(t, ok)

	// Should have 3 keys, sorted alphabetically
	assert.Equal(t, 3, len(slice))
	assert.Equal(t, types.StringValue("key1"), slice[0])
	assert.Equal(t, types.StringValue("key2"), slice[1])
	assert.Equal(t, types.StringValue("key3"), slice[2])
}

func TestBaskitValues(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	valuesMethod := baskitClass.PublicFunctions["VALUES"]

	instance := NewBaskitInstance()

	// PUT multiple values
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key2"),
		types.StringValue("value2"),
	})
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
		types.StringValue("value1"),
	})
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key3"),
		types.StringValue("value3"),
	})

	// Get VALUES
	result, err := valuesMethod.NativeImpl(nil, instance, []types.Value{})
	require.NoError(t, err)

	// Should return a BUKKIT
	bukkitValue, ok := result.(types.ObjectValue)
	require.True(t, ok, "VALUES should return ObjectValue")
	assert.Equal(t, "BUKKIT", bukkitValue.ClassName)

	// Get the underlying slice
	bukkitInstance, ok := bukkitValue.Instance.(*environment.ObjectInstance)
	require.True(t, ok)
	slice, ok := bukkitInstance.NativeData.(BukkitSlice)
	require.True(t, ok)

	// Should have 3 values, in key-sorted order
	assert.Equal(t, 3, len(slice))
	assert.Equal(t, types.StringValue("value1"), slice[0])
	assert.Equal(t, types.StringValue("value2"), slice[1])
	assert.Equal(t, types.StringValue("value3"), slice[2])
}

func TestBaskitPairs(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	pairsMethod := baskitClass.PublicFunctions["PAIRS"]

	instance := NewBaskitInstance()

	// PUT multiple values
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key2"),
		types.IntegerValue(200),
	})
	putMethod.NativeImpl(nil, instance, []types.Value{
		types.StringValue("key1"),
		types.IntegerValue(100),
	})

	// Get PAIRS
	result, err := pairsMethod.NativeImpl(nil, instance, []types.Value{})
	require.NoError(t, err)

	// Should return a BUKKIT
	bukkitValue, ok := result.(types.ObjectValue)
	require.True(t, ok, "PAIRS should return ObjectValue")
	assert.Equal(t, "BUKKIT", bukkitValue.ClassName)

	// Get the underlying slice
	bukkitInstance, ok := bukkitValue.Instance.(*environment.ObjectInstance)
	require.True(t, ok)
	slice, ok := bukkitInstance.NativeData.(BukkitSlice)
	require.True(t, ok)

	// Should have 2 pairs
	assert.Equal(t, 2, len(slice))

	// Each pair should be a BUKKIT with [key, value]
	pair1, ok := slice[0].(types.ObjectValue)
	require.True(t, ok)
	assert.Equal(t, "BUKKIT", pair1.ClassName)

	pair1Instance, ok := pair1.Instance.(*environment.ObjectInstance)
	require.True(t, ok)
	pair1Slice, ok := pair1Instance.NativeData.(BukkitSlice)
	require.True(t, ok)
	assert.Equal(t, 2, len(pair1Slice))
	assert.Equal(t, types.StringValue("key1"), pair1Slice[0])
	assert.Equal(t, types.IntegerValue(100), pair1Slice[1])

	pair2, ok := slice[1].(types.ObjectValue)
	require.True(t, ok)
	pair2Instance, ok := pair2.Instance.(*environment.ObjectInstance)
	require.True(t, ok)
	pair2Slice, ok := pair2Instance.NativeData.(BukkitSlice)
	require.True(t, ok)
	assert.Equal(t, types.StringValue("key2"), pair2Slice[0])
	assert.Equal(t, types.IntegerValue(200), pair2Slice[1])
}

func TestBaskitMerge(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	mergeMethod := baskitClass.PublicFunctions["MERGE"]

	// Create first BASKIT
	instance1 := NewBaskitInstance()
	putMethod.NativeImpl(nil, instance1, []types.Value{
		types.StringValue("key1"),
		types.StringValue("value1"),
	})
	putMethod.NativeImpl(nil, instance1, []types.Value{
		types.StringValue("key2"),
		types.StringValue("value2"),
	})

	// Create second BASKIT
	instance2 := NewBaskitInstance()
	putMethod.NativeImpl(nil, instance2, []types.Value{
		types.StringValue("key2"),
		types.StringValue("new_value2"), // Override existing key
	})
	putMethod.NativeImpl(nil, instance2, []types.Value{
		types.StringValue("key3"),
		types.StringValue("value3"), // New key
	})

	// Create ObjectValue for second BASKIT
	baskit2Value := types.NewObjectValue(instance2, "BASKIT")

	// MERGE second into first
	_, err := mergeMethod.NativeImpl(nil, instance1, []types.Value{baskit2Value})
	require.NoError(t, err)

	// Verify SIZ updated (should be 3: key1, key2, key3)
	sizVar := instance1.Variables["SIZ"]
	assert.Equal(t, types.IntegerValue(3), sizVar.Value)

	// Verify values
	baskitMap := instance1.NativeData.(BaskitMap)
	assert.Equal(t, types.StringValue("value1"), baskitMap["key1"])
	assert.Equal(t, types.StringValue("new_value2"), baskitMap["key2"]) // Overridden
	assert.Equal(t, types.StringValue("value3"), baskitMap["key3"])     // New key
}

func TestBaskitCopy(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	copyMethod := baskitClass.PublicFunctions["COPY"]

	// Create original BASKIT
	original := NewBaskitInstance()
	putMethod.NativeImpl(nil, original, []types.Value{
		types.StringValue("key1"),
		types.IntegerValue(42),
	})
	putMethod.NativeImpl(nil, original, []types.Value{
		types.StringValue("key2"),
		types.StringValue("hello"),
	})

	// COPY the BASKIT
	result, err := copyMethod.NativeImpl(nil, original, []types.Value{})
	require.NoError(t, err)

	// Should return a BASKIT
	copyValue, ok := result.(types.ObjectValue)
	require.True(t, ok, "COPY should return ObjectValue")
	assert.Equal(t, "BASKIT", copyValue.ClassName)

	// Get the copy instance
	copyInstance, ok := copyValue.Instance.(*environment.ObjectInstance)
	require.True(t, ok)
	copyMap, ok := copyInstance.NativeData.(BaskitMap)
	require.True(t, ok)

	// Should have same contents
	originalMap := original.NativeData.(BaskitMap)
	assert.Equal(t, len(originalMap), len(copyMap))
	assert.Equal(t, originalMap["key1"], copyMap["key1"])
	assert.Equal(t, originalMap["key2"], copyMap["key2"])

	// Should be separate instances (shallow copy)
	assert.NotSame(t, original, copyInstance)
	assert.NotSame(t, originalMap, copyMap)

	// Verify SIZ is copied correctly
	copySizVar := copyInstance.Variables["SIZ"]
	originalSizVar := original.Variables["SIZ"]
	assert.Equal(t, originalSizVar.Value, copySizVar.Value)
}

func TestBaskitMixedTypes(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	getMethod := baskitClass.PublicFunctions["GET"]

	instance := NewBaskitInstance()

	// Test different value types
	testValues := []types.Value{
		types.IntegerValue(42),
		types.DoubleValue(3.14),
		types.StringValue("hello"),
		types.YEZ,
		types.NO,
		types.NOTHIN,
	}

	// PUT different types
	for i, value := range testValues {
		key := types.StringValue("key" + string(rune('1'+i)))
		_, err := putMethod.NativeImpl(nil, instance, []types.Value{key, value})
		require.NoError(t, err)
	}

	// GET and verify each type
	for i, expectedValue := range testValues {
		key := types.StringValue("key" + string(rune('1'+i)))
		result, err := getMethod.NativeImpl(nil, instance, []types.Value{key})
		require.NoError(t, err)
		assert.Equal(t, expectedValue, result)
	}
}

func TestBaskitKeyTypeConversion(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	getMethod := baskitClass.PublicFunctions["GET"]
	containsMethod := baskitClass.PublicFunctions["CONTAINS"]

	instance := NewBaskitInstance()

	// Test different key types that should all convert to strings
	testKeys := []types.Value{
		types.StringValue("string_key"),
		types.IntegerValue(123),
		types.DoubleValue(45.67),
		types.YEZ,
		types.NO,
	}

	// PUT values with different key types
	for i, key := range testKeys {
		value := types.IntegerValue(100 + i)
		_, err := putMethod.NativeImpl(nil, instance, []types.Value{key, value})
		require.NoError(t, err)
	}

	// GET values using same key types
	for i, key := range testKeys {
		expectedValue := types.IntegerValue(100 + i)

		result, err := getMethod.NativeImpl(nil, instance, []types.Value{key})
		require.NoError(t, err)
		assert.Equal(t, expectedValue, result)

		// Also test CONTAINS
		containsResult, err := containsMethod.NativeImpl(nil, instance, []types.Value{key})
		require.NoError(t, err)
		assert.Equal(t, types.YEZ, containsResult)
	}
}

func TestNewBaskitInstance(t *testing.T) {
	instance := NewBaskitInstance()

	// Verify basic structure
	assert.NotNil(t, instance)
	assert.Equal(t, []string{"BASKIT"}, instance.Hierarchy)

	// Verify NativeData is initialized
	baskitMap, ok := instance.NativeData.(BaskitMap)
	require.True(t, ok)
	assert.Equal(t, 0, len(baskitMap))

	// Verify SIZ variable
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists)
	assert.Equal(t, "SIZ", sizVar.Name)
	assert.Equal(t, "INTEGR", sizVar.Type)
	assert.Equal(t, types.IntegerValue(0), sizVar.Value)
	assert.True(t, sizVar.IsLocked)
}

func TestRegisterBaskitInEnv(t *testing.T) {
	env := environment.NewEnvironment(nil)

	err := RegisterMapsInEnv(env)
	require.NoError(t, err)

	// Verify BASKIT class is registered
	class, err := env.GetClass("BASKIT")
	require.NoError(t, err)
	require.NotNil(t, class)
	assert.Equal(t, "BASKIT", class.Name)

	// Verify it has expected methods
	expectedMethods := []string{
		"BASKIT", "PUT", "GET", "CONTAINS", "REMOVE", "CLEAR",
		"KEYS", "VALUES", "PAIRS", "MERGE", "COPY",
	}

	for _, methodName := range expectedMethods {
		_, exists := class.PublicFunctions[methodName]
		assert.True(t, exists, "Method %s should exist", methodName)
	}

	// Verify SIZ variable exists
	_, exists := class.PublicVariables["SIZ"]
	assert.True(t, exists, "SIZ variable should exist")
}
