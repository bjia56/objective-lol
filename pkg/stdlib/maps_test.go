package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaskitConstructor(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	constructor := baskitClass.PublicFunctions["BASKIT"]

	// Create a new BASKIT instance
	instance := NewBaskitInstance()

	// Call constructor
	_, err := constructor.NativeImpl(nil, instance, []environment.Value{})
	require.NoError(t, err)

	// Verify empty map
	baskitMap, ok := instance.NativeData.(BaskitMap)
	require.True(t, ok, "NativeData should be BaskitMap")
	assert.Equal(t, 0, len(baskitMap))

	// Verify SIZ is 0
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists)
	assert.Equal(t, environment.IntegerValue(0), sizVar.Value)
}

func TestBaskitPutAndGet(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	getMethod := baskitClass.PublicFunctions["GET"]

	instance := NewBaskitInstance()

	// PUT a value
	_, err := putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
		environment.IntegerValue(42),
	})
	require.NoError(t, err)

	// Verify SIZ updated
	sizVar := instance.Variables["SIZ"]
	assert.Equal(t, environment.IntegerValue(1), sizVar.Value)

	// GET the value
	result, err := getMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, environment.IntegerValue(42), result)

	// PUT another value with different key type (should convert to string)
	_, err = putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.IntegerValue(123),
		environment.StringValue("numeric key"),
	})
	require.NoError(t, err)

	// GET with numeric key (converted to string)
	result, err = getMethod.NativeImpl(nil, instance, []environment.Value{
		environment.IntegerValue(123),
	})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("numeric key"), result)

	// Verify SIZ updated
	sizVar = instance.Variables["SIZ"]
	assert.Equal(t, environment.IntegerValue(2), sizVar.Value)
}

func TestBaskitGetNonexistentKey(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	getMethod := baskitClass.PublicFunctions["GET"]

	instance := NewBaskitInstance()

	// Try to GET non-existent key - should throw exception
	_, err := getMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("nonexistent"),
	})
	require.Error(t, err)

	// Check it's the right type of exception
	exception, ok := err.(runtime.Exception)
	require.True(t, ok, "Should be an runtime.Exception")
	assert.Contains(t, exception.Message, "Key 'nonexistent' not found in BASKIT")
}

func TestBaskitContains(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	containsMethod := baskitClass.PublicFunctions["CONTAINS"]

	instance := NewBaskitInstance()

	// Check non-existent key
	result, err := containsMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, environment.NO, result)

	// PUT a value
	_, err = putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
		environment.IntegerValue(42),
	})
	require.NoError(t, err)

	// Check existing key
	result, err = containsMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, environment.YEZ, result)
}

func TestBaskitRemove(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	removeMethod := baskitClass.PublicFunctions["REMOVE"]
	containsMethod := baskitClass.PublicFunctions["CONTAINS"]

	instance := NewBaskitInstance()

	// PUT a value
	_, err := putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
		environment.StringValue("value1"),
	})
	require.NoError(t, err)

	// REMOVE the value
	result, err := removeMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("value1"), result)

	// Verify SIZ updated
	sizVar := instance.Variables["SIZ"]
	assert.Equal(t, environment.IntegerValue(0), sizVar.Value)

	// Verify key no longer exists
	containsResult, err := containsMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
	})
	require.NoError(t, err)
	assert.Equal(t, environment.NO, containsResult)

	// Try to remove non-existent key - should throw exception
	_, err = removeMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("nonexistent"),
	})
	require.Error(t, err)

	// Check it's the right type of exception
	exception, ok := err.(runtime.Exception)
	require.True(t, ok, "Should be an runtime.Exception")
	assert.Contains(t, exception.Message, "Key 'nonexistent' not found in BASKIT")
}

func TestBaskitClear(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	clearMethod := baskitClass.PublicFunctions["CLEAR"]

	instance := NewBaskitInstance()

	// PUT multiple values
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
		environment.IntegerValue(1),
	})
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key2"),
		environment.IntegerValue(2),
	})
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key3"),
		environment.IntegerValue(3),
	})

	// Verify SIZ before clear
	sizVar := instance.Variables["SIZ"]
	assert.Equal(t, environment.IntegerValue(3), sizVar.Value)

	// CLEAR all values
	_, err := clearMethod.NativeImpl(nil, instance, []environment.Value{})
	require.NoError(t, err)

	// Verify SIZ is 0
	sizVar = instance.Variables["SIZ"]
	assert.Equal(t, environment.IntegerValue(0), sizVar.Value)

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
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key2"),
		environment.IntegerValue(2),
	})
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
		environment.IntegerValue(1),
	})
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key3"),
		environment.IntegerValue(3),
	})

	// Get KEYS
	result, err := keysMethod.NativeImpl(nil, instance, []environment.Value{})
	require.NoError(t, err)

	// Should return a BUKKIT
	bukkitInstance, ok := result.(*environment.ObjectInstance)
	require.True(t, ok, "KEYS should return ObjectInstance")
	assert.Equal(t, "BUKKIT", bukkitInstance.Class.Name)

	// Get the underlying slice
	slice, ok := bukkitInstance.NativeData.(BukkitSlice)
	require.True(t, ok)

	// Should have 3 keys, sorted alphabetically
	assert.Equal(t, 3, len(slice))
	assert.Equal(t, environment.StringValue("key1"), slice[0])
	assert.Equal(t, environment.StringValue("key2"), slice[1])
	assert.Equal(t, environment.StringValue("key3"), slice[2])
}

func TestBaskitValues(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	valuesMethod := baskitClass.PublicFunctions["VALUES"]

	instance := NewBaskitInstance()

	// PUT multiple values
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key2"),
		environment.StringValue("value2"),
	})
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
		environment.StringValue("value1"),
	})
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key3"),
		environment.StringValue("value3"),
	})

	// Get VALUES
	result, err := valuesMethod.NativeImpl(nil, instance, []environment.Value{})
	require.NoError(t, err)

	// Should return a BUKKIT
	bukkitInstance, ok := result.(*environment.ObjectInstance)
	require.True(t, ok, "VALUES should return ObjectInstance")
	assert.Equal(t, "BUKKIT", bukkitInstance.Class.Name)

	// Get the underlying slice
	slice, ok := bukkitInstance.NativeData.(BukkitSlice)
	require.True(t, ok)

	// Should have 3 values, in key-sorted order
	assert.Equal(t, 3, len(slice))
	assert.Equal(t, environment.StringValue("value1"), slice[0])
	assert.Equal(t, environment.StringValue("value2"), slice[1])
	assert.Equal(t, environment.StringValue("value3"), slice[2])
}

func TestBaskitPairs(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	pairsMethod := baskitClass.PublicFunctions["PAIRS"]

	instance := NewBaskitInstance()

	// PUT multiple values
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key2"),
		environment.IntegerValue(200),
	})
	putMethod.NativeImpl(nil, instance, []environment.Value{
		environment.StringValue("key1"),
		environment.IntegerValue(100),
	})

	// Get PAIRS
	result, err := pairsMethod.NativeImpl(nil, instance, []environment.Value{})
	require.NoError(t, err)

	// Should return a BUKKIT
	bukkitInstance, ok := result.(*environment.ObjectInstance)
	require.True(t, ok, "PAIRS should return ObjectInstance")
	assert.Equal(t, "BUKKIT", bukkitInstance.Class.Name)

	// Get the underlying slice
	slice, ok := bukkitInstance.NativeData.(BukkitSlice)
	require.True(t, ok)

	// Should have 2 pairs
	assert.Equal(t, 2, len(slice))

	// Each pair should be a BUKKIT with [key, value]
	pair1Instance, ok := slice[0].(*environment.ObjectInstance)
	require.True(t, ok)
	assert.Equal(t, "BUKKIT", pair1Instance.Class.Name)

	pair1Slice, ok := pair1Instance.NativeData.(BukkitSlice)
	require.True(t, ok)
	assert.Equal(t, 2, len(pair1Slice))
	assert.Equal(t, environment.StringValue("key1"), pair1Slice[0])
	assert.Equal(t, environment.IntegerValue(100), pair1Slice[1])

	pair2Instance, ok := slice[1].(*environment.ObjectInstance)
	require.True(t, ok)
	pair2Slice, ok := pair2Instance.NativeData.(BukkitSlice)
	require.True(t, ok)
	assert.Equal(t, environment.StringValue("key2"), pair2Slice[0])
	assert.Equal(t, environment.IntegerValue(200), pair2Slice[1])
}

func TestBaskitMerge(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	mergeMethod := baskitClass.PublicFunctions["MERGE"]

	// Create first BASKIT
	instance1 := NewBaskitInstance()
	putMethod.NativeImpl(nil, instance1, []environment.Value{
		environment.StringValue("key1"),
		environment.StringValue("value1"),
	})
	putMethod.NativeImpl(nil, instance1, []environment.Value{
		environment.StringValue("key2"),
		environment.StringValue("value2"),
	})

	// Create second BASKIT
	instance2 := NewBaskitInstance()
	putMethod.NativeImpl(nil, instance2, []environment.Value{
		environment.StringValue("key2"),
		environment.StringValue("new_value2"), // Override existing key
	})
	putMethod.NativeImpl(nil, instance2, []environment.Value{
		environment.StringValue("key3"),
		environment.StringValue("value3"), // New key
	})

	// MERGE second into first
	_, err := mergeMethod.NativeImpl(nil, instance1, []environment.Value{instance2})
	require.NoError(t, err)

	// Verify SIZ updated (should be 3: key1, key2, key3)
	sizVar := instance1.Variables["SIZ"]
	assert.Equal(t, environment.IntegerValue(3), sizVar.Value)

	// Verify values
	baskitMap := instance1.NativeData.(BaskitMap)
	assert.Equal(t, environment.StringValue("value1"), baskitMap["key1"])
	assert.Equal(t, environment.StringValue("new_value2"), baskitMap["key2"]) // Overridden
	assert.Equal(t, environment.StringValue("value3"), baskitMap["key3"])     // New key
}

func TestBaskitCopy(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	copyMethod := baskitClass.PublicFunctions["COPY"]

	// Create original BASKIT
	original := NewBaskitInstance()
	putMethod.NativeImpl(nil, original, []environment.Value{
		environment.StringValue("key1"),
		environment.IntegerValue(42),
	})
	putMethod.NativeImpl(nil, original, []environment.Value{
		environment.StringValue("key2"),
		environment.StringValue("hello"),
	})

	// COPY the BASKIT
	result, err := copyMethod.NativeImpl(nil, original, []environment.Value{})
	require.NoError(t, err)

	// Should return a BASKIT
	copyInstance, ok := result.(*environment.ObjectInstance)
	require.True(t, ok, "COPY should return ObjectInstance")
	assert.Equal(t, "BASKIT", copyInstance.Class.Name)

	// Get the copy instance
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

func TestBaskitMixedenvironment(t *testing.T) {
	baskitClass := getMapClasses()["BASKIT"]
	putMethod := baskitClass.PublicFunctions["PUT"]
	getMethod := baskitClass.PublicFunctions["GET"]

	instance := NewBaskitInstance()

	// Test different value environment
	testValues := []environment.Value{
		environment.IntegerValue(42),
		environment.DoubleValue(3.14),
		environment.StringValue("hello"),
		environment.YEZ,
		environment.NO,
		environment.NOTHIN,
	}

	// PUT different environment
	for i, value := range testValues {
		key := environment.StringValue("key" + string(rune('1'+i)))
		_, err := putMethod.NativeImpl(nil, instance, []environment.Value{key, value})
		require.NoError(t, err)
	}

	// GET and verify each type
	for i, expectedValue := range testValues {
		key := environment.StringValue("key" + string(rune('1'+i)))
		result, err := getMethod.NativeImpl(nil, instance, []environment.Value{key})
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

	// Test different key environment that should all convert to strings
	testKeys := []environment.Value{
		environment.StringValue("string_key"),
		environment.IntegerValue(123),
		environment.DoubleValue(45.67),
		environment.YEZ,
		environment.NO,
	}

	// PUT values with different key environment
	for i, key := range testKeys {
		value := environment.IntegerValue(100 + i)
		_, err := putMethod.NativeImpl(nil, instance, []environment.Value{key, value})
		require.NoError(t, err)
	}

	// GET values using same key environment
	for i, key := range testKeys {
		expectedValue := environment.IntegerValue(100 + i)

		result, err := getMethod.NativeImpl(nil, instance, []environment.Value{key})
		require.NoError(t, err)
		assert.Equal(t, expectedValue, result)

		// Also test CONTAINS
		containsResult, err := containsMethod.NativeImpl(nil, instance, []environment.Value{key})
		require.NoError(t, err)
		assert.Equal(t, environment.YEZ, containsResult)
	}
}

func TestNewBaskitInstance(t *testing.T) {
	instance := NewBaskitInstance()

	// Verify basic structure
	assert.NotNil(t, instance)
	assert.Equal(t, []string{"stdlib:MAPS.BASKIT"}, instance.MRO)

	// Verify NativeData is initialized
	baskitMap, ok := instance.NativeData.(BaskitMap)
	require.True(t, ok)
	assert.Equal(t, 0, len(baskitMap))

	// Verify SIZ variable
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists)
	assert.Equal(t, "SIZ", sizVar.Name)
	assert.Equal(t, "INTEGR", sizVar.Type)
	assert.Equal(t, environment.IntegerValue(0), sizVar.Value)
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
